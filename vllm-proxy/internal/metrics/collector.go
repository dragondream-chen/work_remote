package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsCollector struct {
	RequestsTotal     *prometheus.CounterVec
	RequestDuration   *prometheus.HistogramVec
	ActiveRequests    *prometheus.GaugeVec
	BackendLatency    *prometheus.HistogramVec
	ConnectionPoolSize *prometheus.GaugeVec
	BackendErrors     *prometheus.CounterVec
	LoadBalanceScore  *prometheus.GaugeVec
	KVTransferLatency *prometheus.HistogramVec
	ServerActiveTokens *prometheus.GaugeVec
	ServerActiveKVCache *prometheus.GaugeVec
	InstanceHealth    *prometheus.GaugeVec
}

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vllm_proxy_requests_total",
			Help: "Total number of requests by endpoint and status",
		},
		[]string{"endpoint", "status", "method"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "vllm_proxy_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
		},
		[]string{"endpoint"},
	)

	ActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vllm_proxy_active_requests",
			Help: "Number of active requests",
		},
		[]string{"type"},
	)

	BackendLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "vllm_proxy_backend_latency_seconds",
			Help:    "Backend request latency in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
		},
		[]string{"backend_type", "server"},
	)

	ConnectionPoolSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vllm_proxy_connection_pool_size",
			Help: "Size of connection pool",
		},
		[]string{"server"},
	)

	BackendErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vllm_proxy_backend_errors_total",
			Help: "Total number of backend errors",
		},
		[]string{"backend_type", "server", "error_type"},
	)

	LoadBalanceScore = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vllm_proxy_load_balance_score",
			Help: "Current load balance score for each server",
		},
		[]string{"server_type", "server"},
	)

	KVTransferLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "vllm_proxy_kv_transfer_latency_seconds",
			Help:    "KV transfer latency in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"operation"},
	)

	ServerActiveTokens = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vllm_proxy_server_active_tokens",
			Help: "Number of active tokens on each server",
		},
		[]string{"server_type", "server"},
	)

	ServerActiveKVCache = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vllm_proxy_server_active_kv_cache",
			Help: "Number of active KV cache entries on each server",
		},
		[]string{"server_type", "server"},
	)

	InstanceHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vllm_proxy_instance_health",
			Help: "Health status of each instance (1=healthy, 0=unhealthy)",
		},
		[]string{"server_type", "server"},
	)
)

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		RequestsTotal: RequestsTotal,
		RequestDuration: RequestDuration,
		ActiveRequests: ActiveRequests,
		BackendLatency: BackendLatency,
		ConnectionPoolSize: ConnectionPoolSize,
		BackendErrors: BackendErrors,
		LoadBalanceScore: LoadBalanceScore,
		KVTransferLatency: KVTransferLatency,
		ServerActiveTokens: ServerActiveTokens,
		ServerActiveKVCache: ServerActiveKVCache,
		InstanceHealth: InstanceHealth,
	}
}

func (m *MetricsCollector) RecordRequest(endpoint, status, method string) {
	m.RequestsTotal.WithLabelValues(endpoint, status, method).Inc()
}

func (m *MetricsCollector) RecordDuration(endpoint string, duration float64) {
	m.RequestDuration.WithLabelValues(endpoint).Observe(duration)
}

func (m *MetricsCollector) SetActiveRequests(count float64, reqType string) {
	m.ActiveRequests.WithLabelValues(reqType).Set(count)
}

func (m *MetricsCollector) RecordBackendLatency(backendType, server string, latency float64) {
	m.BackendLatency.WithLabelValues(backendType, server).Observe(latency)
}

func (m *MetricsCollector) RecordBackendError(backendType, server, errorType string) {
	m.BackendErrors.WithLabelValues(backendType, server, errorType).Inc()
}

func (m *MetricsCollector) SetLoadBalanceScore(serverType, server string, score float64) {
	m.LoadBalanceScore.WithLabelValues(serverType, server).Set(score)
}

func (m *MetricsCollector) RecordKVTransferLatency(operation string, latency float64) {
	m.KVTransferLatency.WithLabelValues(operation).Observe(latency)
}

func (m *MetricsCollector) SetServerActiveTokens(serverType, server string, count float64) {
	m.ServerActiveTokens.WithLabelValues(serverType, server).Set(count)
}

func (m *MetricsCollector) SetServerActiveKVCache(serverType, server string, count float64) {
	m.ServerActiveKVCache.WithLabelValues(serverType, server).Set(count)
}

func (m *MetricsCollector) SetInstanceHealth(serverType, server string, healthy bool) {
	value := float64(0)
	if healthy {
		value = 1
	}
	m.InstanceHealth.WithLabelValues(serverType, server).Set(value)
}
