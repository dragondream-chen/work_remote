package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vllm-ascend/vllm-proxy/config"
	"github.com/vllm-ascend/vllm-proxy/internal/instance"
	"github.com/vllm-ascend/vllm-proxy/internal/kvtransfer"
	"github.com/vllm-ascend/vllm-proxy/internal/loadbalancer"
	"github.com/vllm-ascend/vllm-proxy/internal/metrics"
	"go.uber.org/zap"
)

type ProxyServer struct {
	config        *config.Config
	httpServer    *http.Server
	instanceMgr   *instance.InstanceManager
	loadBalancer  *loadbalancer.LoadBalancer
	kvHandler     *kvtransfer.KVTransferHandler
	metrics       *metrics.MetricsCollector
	logger        *zap.Logger
	reqDataMap    sync.Map
}

func NewProxyServer(cfg *config.Config, logger *zap.Logger) *ProxyServer {
	prefillers := make([]*loadbalancer.ServerState, 0, len(cfg.Prefillers))
	for _, inst := range cfg.Prefillers {
		prefillers = append(prefillers, loadbalancer.NewServerState(inst.Host, inst.Port, inst.Weight))
	}

	decoders := make([]*loadbalancer.ServerState, 0, len(cfg.Decoders))
	for _, inst := range cfg.Decoders {
		decoders = append(decoders, loadbalancer.NewServerState(inst.Host, inst.Port, inst.Weight))
	}

	instMgr := instance.NewInstanceManager(cfg, logger)
	lb := loadbalancer.NewLoadBalancer(prefillers, decoders)
	kvHandler := kvtransfer.NewKVTransferHandler(cfg, lb, logger)

	return &ProxyServer{
		config:       cfg,
		instanceMgr:  instMgr,
		loadBalancer: lb,
		kvHandler:    kvHandler,
		metrics:      metrics.NewMetricsCollector(),
		logger:       logger,
		reqDataMap:   sync.Map{},
	}
}

func (s *ProxyServer) SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(s.loggingMiddleware())
	router.Use(s.metricsMiddleware())

	router.POST("/v1/completions", s.handleCompletions)
	router.POST("/v1/chat/completions", s.handleChatCompletions)
	router.GET("/healthcheck", s.handleHealthCheck)
	router.POST("/instances/add", s.handleAddInstances)
	router.POST("/instances/remove", s.handleRemoveInstances)
	router.POST("/v1/metaserver", s.handleMetaServer)

	if s.config.Metrics.Enabled {
		router.GET(s.config.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	return router
}

func (s *ProxyServer) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		s.logger.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func (s *ProxyServer) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := fmt.Sprintf("%d", c.Writer.Status())

		s.metrics.RecordRequest(path, status, c.Request.Method)
		s.metrics.RecordDuration(path, duration)
	}
}

func (s *ProxyServer) handleCompletions(c *gin.Context) {
	s.handleCompletionsInternal("/completions", c)
}

func (s *ProxyServer) handleChatCompletions(c *gin.Context) {
	s.handleCompletionsInternal("/chat/completions", c)
}

func (s *ProxyServer) handleCompletionsInternal(api string, c *gin.Context) {
	s.kvHandler.IncrementActiveRequests()
	defer s.kvHandler.DecrementActiveRequests()

	var reqData map[string]interface{}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	body, _ := json.Marshal(reqData)
	requestLength := len(body)
	decoderScore := loadbalancer.CalculateDecodeScore(requestLength)

	decoder, decoderIdx := s.loadBalancer.SelectDecoder(decoderScore)
	if decoder == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no decoder available"})
		return
	}

	requestID := s.kvHandler.GenerateRequestID()
	requestIDAPI := s.kvHandler.GetAPIRequestID(api, requestID)

	streamFlag := false
	if stream, ok := reqData["stream"].(bool); ok {
		streamFlag = stream
	}

	originMaxTokens := 16
	if mt, ok := reqData["max_tokens"].(float64); ok {
		originMaxTokens = int(mt)
	}

	originPrompt := ""
	chatFlag := false
	if prompt, ok := reqData["prompt"]; ok {
		if p, ok := prompt.(string); ok {
			originPrompt = p
		}
	} else if messages, ok := reqData["messages"].([]interface{}); ok && len(messages) > 0 {
		chatFlag = true
		if msg, ok := messages[0].(map[string]interface{}); ok {
			if content, ok := msg["content"]; ok {
				switch c := content.(type) {
				case string:
					originPrompt = c
				case []interface{}:
					if len(c) > 0 {
						if item, ok := c[0].(map[string]interface{}); ok {
							if text, ok := item["text"].(string); ok {
								originPrompt = text
							}
						}
					}
				}
			}
		}
	}

	info := &kvtransfer.RequestInfo{
		RequestID:       requestID,
		DecoderIdx:      decoderIdx,
		DecoderScore:    decoderScore,
		ReqData:         reqData,
		OriginPrompt:    originPrompt,
		OriginMaxTokens: originMaxTokens,
	}

	reqData["kv_transfer_params"] = map[string]interface{}{
		"do_remote_decode":  false,
		"do_remote_prefill": true,
		"metaserver":        fmt.Sprintf("http://%s:%d/v1/metaserver", s.config.Server.Host, s.config.Server.Port),
	}

	s.reqDataMap.Store(requestIDAPI, info)

	resp, err := s.kvHandler.SendDecoderRequest(decoder, api, reqData, requestID)
	if err != nil {
		s.loadBalancer.ReleaseDecoder(decoderIdx, decoderScore)
		s.logger.Error("failed to send request to decoder",
			zap.String("decoder", decoder.Address()),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if streamFlag {
		contentType = "text/event-stream"
	}
	c.Header("Content-Type", contentType)

	s.streamResponse(c, resp, decoderIdx, decoderScore, info, streamFlag, chatFlag, api, reqData)
}

func (s *ProxyServer) streamResponse(
	c *gin.Context,
	resp *http.Response,
	decoderIdx int,
	decoderScore float64,
	info *kvtransfer.RequestInfo,
	streamFlag bool,
	chatFlag bool,
	api string,
	reqData map[string]interface{},
) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	reader := bufio.NewReader(resp.Body)
	generatedToken := ""
	completionTokens := 0
	retryCount := 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			s.logger.Error("error reading stream", zap.Error(err))
			break
		}

		lineStr := strings.TrimSpace(string(line))
		if lineStr == "" {
			continue
		}

		if strings.HasPrefix(lineStr, "data: ") {
			lineStr = strings.TrimPrefix(lineStr, "data: ")
		}

		if lineStr == "[done]" {
			c.Write([]byte("data: [done]\n\n"))
			flusher.Flush()
			break
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal([]byte(lineStr), &chunk); err != nil {
			c.Write(line)
			flusher.Flush()
			continue
		}

		choices, ok := chunk["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			data, _ := json.Marshal(chunk)
			c.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()
			continue
		}

		choice, ok := choices[0].(map[string]interface{})
		if !ok {
			data, _ := json.Marshal(chunk)
			c.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()
			continue
		}

		var content string
		if delta, ok := choice["delta"].(map[string]interface{}); ok {
			if c, ok := delta["content"].(string); ok {
				content = c
			}
		} else if message, ok := choice["message"].(map[string]interface{}); ok {
			if c, ok := message["content"].(string); ok {
				content = c
			}
		} else if text, ok := choice["text"].(string); ok {
			content = text
		}

		generatedToken += content

		if streamFlag {
			completionTokens++
		} else {
			if usage, ok := chunk["usage"].(map[string]interface{}); ok {
				if ct, ok := usage["completion_tokens"].(float64); ok {
					completionTokens = int(ct)
				}
			}
		}

		stopReason, _ := choice["finish_reason"].(string)
		if stopReason == "recomputed" {
			retryCount++
			if chatFlag {
				if messages, ok := reqData["messages"].([]interface{}); ok && len(messages) > 0 {
					if msg, ok := messages[0].(map[string]interface{}); ok {
						msg["content"] = info.OriginPrompt + generatedToken
					}
				}
			} else {
				reqData["prompt"] = info.OriginPrompt + generatedToken
			}
			reqData["max_tokens"] = info.OriginMaxTokens - completionTokens + retryCount

			s.loadBalancer.ReleaseDecoder(decoderIdx, decoderScore)

			newBody, _ := json.Marshal(reqData)
			newDecoderScore := loadbalancer.CalculateDecodeScore(len(newBody))
			newDecoder, newDecoderIdx := s.loadBalancer.SelectDecoder(newDecoderScore)
			if newDecoder == nil {
				s.logger.Error("no decoder available for retry")
				break
			}

			decoderIdx = newDecoderIdx
			decoderScore = newDecoderScore
			decoder = newDecoder

			newResp, err := s.kvHandler.SendDecoderRequest(decoder, api, reqData, info.RequestID)
			if err != nil {
				s.loadBalancer.ReleaseDecoder(decoderIdx, decoderScore)
				break
			}
			defer newResp.Body.Close()
			resp = newResp
			reader = bufio.NewReader(resp.Body)
			continue
		}

		if retryCount > 0 && !streamFlag {
			if chatFlag {
				choice["message"] = map[string]interface{}{"content": generatedToken}
			} else {
				choice["text"] = generatedToken
			}
		}

		data, _ := json.Marshal(chunk)
		c.Write([]byte("data: " + string(data) + "\n\n"))
		flusher.Flush()
	}

	s.loadBalancer.ReleaseDecoder(decoderIdx, decoderScore)
}

func (s *ProxyServer) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":            "ok",
		"prefill_instances": s.loadBalancer.PrefillerCount(),
		"decode_instances":  s.loadBalancer.DecoderCount(),
		"active_requests":   s.kvHandler.ActiveRequests(),
	})
}

func (s *ProxyServer) handleAddInstances(c *gin.Context) {
	var req struct {
		Type      string   `json:"type"`
		Instances []string `json:"instances"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != "prefill" && req.Type != "decode" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "type must be 'prefill' or 'decode'",
		})
		return
	}

	var added []string
	var waiting []string

	for _, inst := range req.Instances {
		parts := strings.Split(inst, ":")
		if len(parts) != 2 {
			continue
		}

		host := parts[0]
		var port int
		fmt.Sscanf(parts[1], "%d", &port)

		var instanceType instance.InstanceType
		if req.Type == "prefill" {
			instanceType = instance.InstanceTypePrefill
		} else {
			instanceType = instance.InstanceTypeDecode
		}

		err := s.instanceMgr.AddInstance(instanceType, host, port, 1)
		if err != nil {
			waiting = append(waiting, inst)
		} else {
			added = append(added, inst)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                   fmt.Sprintf("Added %d instances, %d waiting", len(added), len(waiting)),
		"added_instances":           added,
		"waiting_instances":         waiting,
		"current_prefill_instances": s.instanceMgr.PrefillerCount(),
		"current_decode_instances":  s.instanceMgr.DecoderCount(),
	})
}

func (s *ProxyServer) handleRemoveInstances(c *gin.Context) {
	var req struct {
		Type      interface{} `json:"type"`
		Instances interface{} `json:"instances"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	typeStr, ok := req.Type.(string)
	if !ok || (typeStr != "prefill" && typeStr != "decode") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "type must be 'prefill' or 'decode'",
		})
		return
	}

	var instances []string
	switch v := req.Instances.(type) {
	case string:
		instances = []string{v}
	case []interface{}:
		for _, inst := range v {
			if str, ok := inst.(string); ok {
				instances = append(instances, str)
			}
		}
	case []string:
		instances = v
	}

	var removed []string
	var failed []string

	for _, inst := range instances {
		parts := strings.Split(inst, ":")
		if len(parts) != 2 {
			failed = append(failed, inst)
			continue
		}

		host := parts[0]
		var port int
		fmt.Sscanf(parts[1], "%d", &port)

		var instanceType instance.InstanceType
		if typeStr == "prefill" {
			instanceType = instance.InstanceTypePrefill
		} else {
			instanceType = instance.InstanceTypeDecode
		}

		err := s.instanceMgr.RemoveInstance(instanceType, host, port)
		if err != nil {
			failed = append(failed, inst)
		} else {
			removed = append(removed, inst)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                   fmt.Sprintf("Removed %d instances, %d failed", len(removed), len(failed)),
		"removed_instances":         removed,
		"failed_instances":          failed,
		"current_prefill_instances": s.instanceMgr.PrefillerCount(),
		"current_decode_instances":  s.instanceMgr.DecoderCount(),
	})
}

func (s *ProxyServer) handleMetaServer(c *gin.Context) {
	var kvParams map[string]interface{}
	if err := c.ShouldBindJSON(&kvParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID, ok := kvParams["request_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing request_id"})
		return
	}

	var api string
	if strings.HasPrefix(requestID, "cmpl-") {
		api = "/completions"
	} else if strings.HasPrefix(requestID, "chatcmpl-") {
		api = "/chat/completions"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request_id format"})
		return
	}

	value, ok := s.reqDataMap.Load(requestID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	info := value.(*kvtransfer.RequestInfo)

	reqData := make(map[string]interface{})
	for k, v := range info.ReqData {
		reqData[k] = v
	}
	reqData["kv_transfer_params"] = kvParams

	reqBody, _ := json.Marshal(reqData)
	requestLength := len(reqBody)
	prefillerScore := loadbalancer.CalculatePrefillScore(requestLength)

	prefiller, prefillerIdx := s.loadBalancer.SelectPrefiller(prefillerScore)
	if prefiller == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no prefiller available"})
		return
	}

	originReqID := s.kvHandler.GetOriginRequestID(api, requestID)

	prefillReqData := make(map[string]interface{})
	for k, v := range reqData {
		prefillReqData[k] = v
	}
	prefillReqData["stream"] = false
	prefillReqData["max_tokens"] = 1
	prefillReqData["min_tokens"] = 1
	if _, ok := prefillReqData["max_completion_tokens"]; ok {
		prefillReqData["max_completion_tokens"] = 1
	}
	if _, ok := prefillReqData["stream_options"]; ok {
		delete(prefillReqData, "stream_options")
	}

	_, err := s.kvHandler.SendPrefillRequest(prefiller, prefillerIdx, api, prefillReqData, originReqID)
	if err != nil {
		s.loadBalancer.ReleasePrefiller(prefillerIdx, prefillerScore)
		s.loadBalancer.ReleasePrefillerKV(prefillerIdx, prefillerScore)
		s.logger.Error("prefill request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.loadBalancer.ReleasePrefiller(prefillerIdx, prefillerScore)
	s.loadBalancer.ReleasePrefillerKV(prefillerIdx, prefillerScore)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *ProxyServer) Start() error {
	router := s.SetupRouter()

	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	s.instanceMgr.StartHealthCheck()

	s.logger.Info("starting proxy server",
		zap.String("address", addr),
		zap.Int("prefillers", s.loadBalancer.PrefillerCount()),
		zap.Int("decoders", s.loadBalancer.DecoderCount()))

	return s.httpServer.ListenAndServe()
}

func (s *ProxyServer) Shutdown() error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(nil)
	}
	return nil
}
