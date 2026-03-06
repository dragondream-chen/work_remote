package instance

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/vllm-ascend/vllm-proxy/config"
	"github.com/vllm-ascend/vllm-proxy/internal/loadbalancer"
	"go.uber.org/zap"
)

type InstanceType string

const (
	InstanceTypePrefill InstanceType = "prefill"
	InstanceTypeDecode  InstanceType = "decode"
)

type InstanceEvent struct {
	Type      InstanceEventType
	Instance  *loadbalancer.ServerState
	Timestamp time.Time
}

type InstanceEventType int

const (
	InstanceAdded InstanceEventType = iota
	InstanceRemoved
	InstanceHealthChanged
)

type HealthChecker struct {
	interval   time.Duration
	timeout    time.Duration
	maxRetries int
	client     *http.Client
	logger     *zap.Logger
}

func NewHealthChecker(interval, timeout time.Duration, maxRetries int, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		interval:   interval,
		timeout:    timeout,
		maxRetries: maxRetries,
		client: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (h *HealthChecker) Check(server *loadbalancer.ServerState) bool {
	url := fmt.Sprintf("http://%s:%d/v1/models", server.Host, server.Port)

	for i := 0; i < h.maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			cancel()
			continue
		}

		req.Header.Set("Authorization", "Bearer "+getAPIKey())

		resp, err := h.client.Do(req)
		cancel()

		if err != nil {
			h.logger.Debug("health check failed",
				zap.String("server", server.Address()),
				zap.Int("attempt", i+1),
				zap.Error(err))
			time.Sleep(h.interval)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true
		}

		time.Sleep(h.interval)
	}

	return false
}

func getAPIKey() string {
	return ""
}

type InstanceManager struct {
	prefillers    []*loadbalancer.ServerState
	decoders      []*loadbalancer.ServerState
	healthChecker *HealthChecker
	eventChan     chan InstanceEvent
	mu            sync.RWMutex
	logger        *zap.Logger

	prefillersToAdd    map[string]*loadbalancer.ServerState
	decodersToAdd      map[string]*loadbalancer.ServerState
	prefillersToRemove map[string]*loadbalancer.ServerState
	decodersToRemove   map[string]*loadbalancer.ServerState
}

func NewInstanceManager(cfg *config.Config, logger *zap.Logger) *InstanceManager {
	mgr := &InstanceManager{
		prefillers:          make([]*loadbalancer.ServerState, 0),
		decoders:            make([]*loadbalancer.ServerState, 0),
		eventChan:           make(chan InstanceEvent, 100),
		logger:              logger,
		prefillersToAdd:     make(map[string]*loadbalancer.ServerState),
		decodersToAdd:       make(map[string]*loadbalancer.ServerState),
		prefillersToRemove:  make(map[string]*loadbalancer.ServerState),
		decodersToRemove:    make(map[string]*loadbalancer.ServerState),
	}

	mgr.healthChecker = NewHealthChecker(
		10*time.Second,
		5*time.Second,
		3,
		logger,
	)

	for _, instCfg := range cfg.Prefillers {
		server := loadbalancer.NewServerState(instCfg.Host, instCfg.Port, instCfg.Weight)
		mgr.prefillers = append(mgr.prefillers, server)
	}

	for _, instCfg := range cfg.Decoders {
		server := loadbalancer.NewServerState(instCfg.Host, instCfg.Port, instCfg.Weight)
		mgr.decoders = append(mgr.decoders, server)
	}

	return mgr
}

func (m *InstanceManager) GetPrefillers() []*loadbalancer.ServerState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*loadbalancer.ServerState, len(m.prefillers))
	copy(result, m.prefillers)
	return result
}

func (m *InstanceManager) GetDecoders() []*loadbalancer.ServerState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*loadbalancer.ServerState, len(m.decoders))
	copy(result, m.decoders)
	return result
}

func (m *InstanceManager) AddInstance(instanceType InstanceType, host string, port int, weight int) error {
	server := loadbalancer.NewServerState(host, port, weight)

	if !m.healthChecker.Check(server) {
		m.mu.Lock()
		if instanceType == InstanceTypePrefill {
			m.prefillersToAdd[server.Address()] = server
		} else {
			m.decodersToAdd[server.Address()] = server
		}
		m.mu.Unlock()

		go m.waitForInstance(instanceType, server)
		return fmt.Errorf("instance not ready, added to waiting queue")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if instanceType == InstanceTypePrefill {
		m.prefillers = append(m.prefillers, server)
	} else {
		m.decoders = append(m.decoders, server)
	}

	m.eventChan <- InstanceEvent{
		Type:      InstanceAdded,
		Instance:  server,
		Timestamp: time.Now(),
	}

	m.logger.Info("instance added",
		zap.String("type", string(instanceType)),
		zap.String("address", server.Address()))

	return nil
}

func (m *InstanceManager) RemoveInstance(instanceType InstanceType, host string, port int) error {
	address := fmt.Sprintf("%s:%d", host, port)

	m.mu.Lock()
	defer m.mu.Unlock()

	var server *loadbalancer.ServerState
	var found bool

	if instanceType == InstanceTypePrefill {
		for i, s := range m.prefillers {
			if s.Address() == address {
				server = s
				m.prefillers = append(m.prefillers[:i], m.prefillers[i+1:]...)
				found = true
				break
			}
		}
	} else {
		for i, s := range m.decoders {
			if s.Address() == address {
				server = s
				m.decoders = append(m.decoders[:i], m.decoders[i+1:]...)
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("instance not found: %s", address)
	}

	m.eventChan <- InstanceEvent{
		Type:      InstanceRemoved,
		Instance:  server,
		Timestamp: time.Now(),
	}

	m.logger.Info("instance removed",
		zap.String("type", string(instanceType)),
		zap.String("address", address))

	return nil
}

func (m *InstanceManager) TaintInstance(instanceType InstanceType, host string, port int) error {
	address := fmt.Sprintf("%s:%d", host, port)

	m.mu.Lock()
	defer m.mu.Unlock()

	var server *loadbalancer.ServerState
	var found bool

	if instanceType == InstanceTypePrefill {
		for _, s := range m.prefillers {
			if s.Address() == address {
				server = s
				found = true
				break
			}
		}
	} else {
		for _, s := range m.decoders {
			if s.Address() == address {
				server = s
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("instance not found: %s", address)
	}

	server.Tainted = true
	m.logger.Info("instance tainted",
		zap.String("type", string(instanceType)),
		zap.String("address", address))

	return nil
}

func (m *InstanceManager) waitForInstance(instanceType InstanceType, server *loadbalancer.ServerState) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	attempts := 0
	maxAttempts := 30

	for {
		select {
		case <-ticker.C:
			attempts++
			if attempts > maxAttempts {
				m.logger.Warn("instance wait timeout, removing from queue",
					zap.String("type", string(instanceType)),
					zap.String("address", server.Address()))

				m.mu.Lock()
				if instanceType == InstanceTypePrefill {
					delete(m.prefillersToAdd, server.Address())
				} else {
					delete(m.decodersToAdd, server.Address())
				}
				m.mu.Unlock()
				return
			}

			if m.healthChecker.Check(server) {
				m.mu.Lock()
				if instanceType == InstanceTypePrefill {
					delete(m.prefillersToAdd, server.Address())
					m.prefillers = append(m.prefillers, server)
				} else {
					delete(m.decodersToAdd, server.Address())
					m.decoders = append(m.decoders, server)
				}
				m.mu.Unlock()

				m.eventChan <- InstanceEvent{
					Type:      InstanceAdded,
					Instance:  server,
					Timestamp: time.Now(),
				}

				m.logger.Info("instance ready and added",
					zap.String("type", string(instanceType)),
					zap.String("address", server.Address()))
				return
			}
		}
	}
}

func (m *InstanceManager) StartHealthCheck() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			m.checkAllInstances()
		}
	}()
}

func (m *InstanceManager) checkAllInstances() {
	m.mu.RLock()
	prefillers := make([]*loadbalancer.ServerState, len(m.prefillers))
	decoders := make([]*loadbalancer.ServerState, len(m.decoders))
	copy(prefillers, m.prefillers)
	copy(decoders, m.decoders)
	m.mu.RUnlock()

	for _, server := range prefillers {
		healthy := m.healthChecker.Check(server)
		if healthy != server.Healthy {
			server.Healthy = healthy
			m.eventChan <- InstanceEvent{
				Type:      InstanceHealthChanged,
				Instance:  server,
				Timestamp: time.Now(),
			}
		}
	}

	for _, server := range decoders {
		healthy := m.healthChecker.Check(server)
		if healthy != server.Healthy {
			server.Healthy = healthy
			m.eventChan <- InstanceEvent{
				Type:      InstanceHealthChanged,
				Instance:  server,
				Timestamp: time.Now(),
			}
		}
	}
}

func (m *InstanceManager) Events() <-chan InstanceEvent {
	return m.eventChan
}

func (m *InstanceManager) PrefillerCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.prefillers)
}

func (m *InstanceManager) DecoderCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.decoders)
}
