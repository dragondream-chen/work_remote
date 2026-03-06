# vLLM High Performance Proxy

A high-performance data routing and distribution proxy for vLLM disaggregated prefill scenarios, written in Go.

## Features

- **High Performance**: Built with Go for optimal concurrency and low latency
- **Load Balancing**: Priority-based load balancing with KV cache awareness
- **KV Transfer Support**: Full support for disaggregated prefill/decode architecture
- **Dynamic Instance Management**: Add/remove prefiller and decoder instances at runtime
- **Streaming Support**: Full SSE streaming for real-time responses
- **Prometheus Metrics**: Built-in metrics for monitoring and observability
- **Multi-node Deployment**: Support for horizontal scaling with external load balancers

## Quick Start

### Build

```bash
make build
```

### Run with Command Line Arguments

```bash
./bin/vllm-proxy \
  --host 0.0.0.0 \
  --port 8000 \
  --prefiller-hosts 10.0.0.1 10.0.0.2 \
  --prefiller-ports 8100 8101 \
  --decoder-hosts 10.0.0.3 10.0.0.4 \
  --decoder-ports 8200 8201
```

### Run with Configuration File

```bash
./bin/vllm-proxy --config configs/config.yaml
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/completions` | POST | OpenAI-compatible text completions |
| `/v1/chat/completions` | POST | OpenAI-compatible chat completions |
| `/healthcheck` | GET | Health check endpoint |
| `/instances/add` | POST | Add prefiller or decoder instances |
| `/instances/remove` | POST | Remove prefiller or decoder instances |
| `/v1/metaserver` | POST | KV transfer metadata server |
| `/metrics` | GET | Prometheus metrics |

## Usage Examples

### Send a Completion Request

```bash
curl -X POST http://localhost:8000/v1/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "DeepSeek-V2-Lite-Chat",
    "prompt": "Hello, world!",
    "max_tokens": 16
  }'
```

### Send a Chat Completion Request

```bash
curl -X POST http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "DeepSeek-V2-Lite-Chat",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 16
  }'
```

### Add Instances Dynamically

```bash
curl -X POST http://localhost:8000/instances/add \
  -H "Content-Type: application/json" \
  -d '{
    "type": "prefill",
    "instances": ["10.0.0.5:8102"]
  }'
```

### Remove Instances

```bash
curl -X POST http://localhost:8000/instances/remove \
  -H "Content-Type: application/json" \
  -d '{
    "type": "decode",
    "instances": ["10.0.0.4:8201"]
  }'
```

## Configuration

### Configuration File (YAML)

```yaml
server:
  host: 0.0.0.0
  port: 8000
  max_connections: 100000
  request_timeout: 30s

prefillers:
  - host: 10.0.0.1
    port: 8100
    weight: 1

decoders:
  - host: 10.0.0.3
    port: 8200
    weight: 1

connection_pool:
  max_idle_conns: 10000
  max_conns_per_host: 1000
  idle_conn_timeout: 90s

retry:
  max_retries: 3
  base_delay: 200ms
  max_delay: 5s

logging:
  level: info
  format: json

metrics:
  enabled: true
  path: /metrics
```

## Deployment

### Docker

```bash
docker build -t vllm-proxy:latest -f deployments/docker/Dockerfile .
docker run -p 8000:8000 -p 9090:9090 vllm-proxy:latest
```

### Docker Compose

```bash
cd deployments/docker
docker-compose up -d
```

### Kubernetes

```bash
kubectl apply -f deployments/kubernetes/deployment.yaml
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Proxy Server (Go)                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ HTTP Server  │  │Load Balancer │  │Instance Mgr  │          │
│  │   (Gin)      │  │ (Priority)   │  │ (Health)     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │KV Transfer   │  │ Client Pool  │  │   Metrics    │          │
│  │  Handler     │  │  (HTTP/2)    │  │ (Prometheus) │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└────────────────────────────────┬────────────────────────────────┘
                                 │
        ┌────────────────────────┴────────────────────────┐
        │                                                 │
        ▼                                                 ▼
┌──────────────────┐                             ┌──────────────────┐
│ Prefiller Pool   │                             │  Decoder Pool    │
│  (kv_producer)   │                             │  (kv_consumer)   │
└──────────────────┘                             └──────────────────┘
```

## Performance

| Metric | Target |
|--------|--------|
| Single Node QPS | >50,000 |
| Average Latency | <10ms |
| Connection Support | >100,000 |
| CPU Utilization | 30-50% |

## Project Structure

```
vllm-proxy/
├── cmd/
│   └── main.go              # Entry point
├── config/
│   └── config.go            # Configuration management
├── internal/
│   ├── server/
│   │   └── server.go        # HTTP server and handlers
│   ├── loadbalancer/
│   │   └── balancer.go      # Load balancing logic
│   ├── instance/
│   │   └── manager.go       # Instance management
│   ├── kvtransfer/
│   │   └── handler.go       # KV transfer handling
│   └── metrics/
│       └── collector.go     # Prometheus metrics
├── deployments/
│   ├── docker/
│   └── kubernetes/
├── configs/
│   └── config.yaml          # Example configuration
├── Makefile
└── go.mod
```

## License

Apache-2.0
