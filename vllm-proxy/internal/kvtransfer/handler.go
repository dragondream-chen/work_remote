package kvtransfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vllm-ascend/vllm-proxy/config"
	"github.com/vllm-ascend/vllm-proxy/internal/loadbalancer"
	"go.uber.org/zap"
)

type KVTransferParams struct {
	DoRemoteDecode   bool                   `json:"do_remote_decode"`
	DoRemotePrefill  bool                   `json:"do_remote_prefill"`
	RemoteEngineID   string                 `json:"remote_engine_id,omitempty"`
	RemoteBlockIDs   []int                  `json:"remote_block_ids,omitempty"`
	RemoteHost       string                 `json:"remote_host,omitempty"`
	RemotePort       int                    `json:"remote_port,omitempty"`
	RemoteTPSize     int                    `json:"remote_tp_size,omitempty"`
	RemotePCPSize    int                    `json:"remote_pcp_size,omitempty"`
	RemoteDCPSize    int                    `json:"remote_dcp_size,omitempty"`
	AbortedRequests  []string               `json:"aborted_requests,omitempty"`
	RequestID        string                 `json:"request_id,omitempty"`
	MetaServer       string                 `json:"metaserver,omitempty"`
	AdditionalParams map[string]interface{} `json:"-"`
}

type RequestInfo struct {
	RequestID        string
	PrefillerIdx     int
	PrefillerScore   float64
	DecoderIdx       int
	DecoderScore     float64
	ReqData          map[string]interface{}
	OriginPrompt     string
	OriginMaxTokens  int
	GeneratedToken   string
	CompletionTokens int
	RetryCount       int
}

type ClientPool struct {
	clients sync.Map
	config  *config.ConnectionPoolConfig
	logger  *zap.Logger
}

func NewClientPool(cfg *config.ConnectionPoolConfig, logger *zap.Logger) *ClientPool {
	return &ClientPool{
		config: cfg,
		logger: logger,
	}
}

func (p *ClientPool) GetClient(host string, port int) *http.Client {
	key := fmt.Sprintf("%s:%d", host, port)

	if client, ok := p.clients.Load(key); ok {
		return client.(*http.Client)
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          p.config.MaxIdleConns,
		MaxIdleConnsPerHost:   p.config.MaxConnsPerHost,
		IdleConnTimeout:       p.config.IdleConnTimeout,
		TLSHandshakeTimeout:   p.config.HandshakeTimeout,
		ResponseHeaderTimeout: p.config.ResponseHeaderTimeout,
		DisableCompression:    false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   0,
	}

	actual, _ := p.clients.LoadOrStore(key, client)
	return actual.(*http.Client)
}

type KVTransferHandler struct {
	clientPool     *ClientPool
	loadBalancer   *loadbalancer.LoadBalancer
	config         *config.Config
	logger         *zap.Logger
	activeRequests int64
	abortedReqs    sync.Map
}

func NewKVTransferHandler(
	cfg *config.Config,
	lb *loadbalancer.LoadBalancer,
	logger *zap.Logger,
) *KVTransferHandler {
	return &KVTransferHandler{
		clientPool:   NewClientPool(&cfg.ConnectionPool, logger),
		loadBalancer: lb,
		config:       cfg,
		logger:       logger,
	}
}

func (h *KVTransferHandler) IncrementActiveRequests() {
	atomic.AddInt64(&h.activeRequests, 1)
}

func (h *KVTransferHandler) DecrementActiveRequests() {
	atomic.AddInt64(&h.activeRequests, -1)
}

func (h *KVTransferHandler) ActiveRequests() int64 {
	return atomic.LoadInt64(&h.activeRequests)
}

func (h *KVTransferHandler) GenerateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), atomic.AddInt64(&h.activeRequests, 1))
}

func (h *KVTransferHandler) GetAPIRequestID(api, reqID string) string {
	if api == "/completions" {
		return "cmpl-" + reqID + "-0"
	} else if api == "/chat/completions" {
		return "chatcmpl-" + reqID
	}
	return reqID
}

func (h *KVTransferHandler) GetOriginRequestID(api, reqID string) string {
	if api == "/completions" {
		if len(reqID) > len("cmpl-")+2 {
			return reqID[len("cmpl-") : len(reqID)-2]
		}
		return reqID
	} else if api == "/chat/completions" {
		if len(reqID) > len("chatcmpl-") {
			return reqID[len("chatcmpl-"):]
		}
		return reqID
	}
	return reqID
}

func (h *KVTransferHandler) SendDecoderRequest(
	server *loadbalancer.ServerState,
	api string,
	reqData map[string]interface{},
	requestID string,
) (*http.Response, error) {
	client := h.clientPool.GetClient(server.Host, server.Port)

	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%d/v1%s", server.Host, server.Port, api)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", requestID)

	return client.Do(req)
}

func (h *KVTransferHandler) SendPrefillRequest(
	server *loadbalancer.ServerState,
	prefillerIdx int,
	api string,
	reqData map[string]interface{},
	requestID string,
) ([]byte, error) {
	client := h.clientPool.GetClient(server.Host, server.Port)

	abortedReqs := h.acquireAbortedRequests(prefillerIdx)
	if len(abortedReqs) > 0 {
		if kvParams, ok := reqData["kv_transfer_params"].(map[string]interface{}); ok {
			kvParams["aborted_requests"] = abortedReqs
		}
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%d/v1%s", server.Host, server.Port, api)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", requestID)

	var lastErr error
	for attempt := 1; attempt <= h.config.Retry.MaxRetries; attempt++ {
		reqBody, _ := json.Marshal(reqData)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			h.logger.Warn("prefill request failed",
				zap.Int("attempt", attempt),
				zap.String("server", server.Address()),
				zap.Error(err))

			if attempt < h.config.Retry.MaxRetries {
				time.Sleep(h.calculateRetryDelay(attempt))
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("prefill request failed with status: %d", resp.StatusCode)
			if attempt < h.config.Retry.MaxRetries {
				time.Sleep(h.calculateRetryDelay(attempt))
			}
			continue
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return buf.Bytes(), nil
	}

	return nil, lastErr
}

func (h *KVTransferHandler) AbortPrefillerRequest(prefillerIdx int, requestID string) {
	key := fmt.Sprintf("prefiller_%d", prefillerIdx)
	value, _ := h.abortedReqs.LoadOrStore(key, &sync.Map{})
	abortedSet := value.(*sync.Map)
	abortedSet.Store(requestID, true)
}

func (h *KVTransferHandler) acquireAbortedRequests(prefillerIdx int) []string {
	key := fmt.Sprintf("prefiller_%d", prefillerIdx)
	value, ok := h.abortedReqs.Load(key)
	if !ok {
		return nil
	}

	abortedSet := value.(*sync.Map)
	var result []string

	abortedSet.Range(func(key, value interface{}) bool {
		result = append(result, key.(string))
		abortedSet.Delete(key)
		return true
	})

	return result
}

func (h *KVTransferHandler) calculateRetryDelay(attempt int) time.Duration {
	delay := h.config.Retry.BaseDelay * time.Duration(1<<(attempt-1))
	if delay > h.config.Retry.MaxDelay {
		delay = h.config.Retry.MaxDelay
	}
	return delay
}
