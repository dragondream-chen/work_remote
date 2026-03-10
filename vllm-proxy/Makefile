.PHONY: build run clean test docker-build docker-push vendor-build vendor-deps

BINARY_NAME=vllm-proxy
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd

vendor-build:
	go build -mod=vendor $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd

vendor-deps:
	@echo "克隆依赖包到 vendor 目录..."
	@if [ -f scripts/clone-deps.sh ]; then \
		chmod +x scripts/clone-deps.sh && \
		./scripts/clone-deps.sh; \
	elif [ -f scripts/clone-deps.ps1 ]; then \
		powershell.exe -ExecutionPolicy Bypass -File scripts/clone-deps.ps1; \
	else \
		echo "错误: 找不到依赖克隆脚本"; \
		exit 1; \
	fi
	@echo "生成 modules.txt..."
	@if [ -f scripts/create-modules-txt.sh ]; then \
		chmod +x scripts/create-modules-txt.sh && \
		./scripts/create-modules-txt.sh; \
	elif [ -f scripts/create-modules-txt.ps1 ]; then \
		powershell.exe -ExecutionPolicy Bypass -File scripts/create-modules-txt.ps1; \
	fi

run:
	go run ./cmd --host 0.0.0.0 --port 8000 \
		--prefiller-hosts localhost localhost \
		--prefiller-ports 8100 8101 \
		--decoder-hosts localhost localhost \
		--decoder-ports 8200 8201

run-config:
	go run ./cmd --config configs/config.yaml

clean:
	rm -rf bin/
	go clean

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

docker-build:
	docker build -t vllm-proxy:$(VERSION) -f deployments/docker/Dockerfile .

docker-push:
	docker tag vllm-proxy:$(VERSION) registry.example.com/vllm-proxy:$(VERSION)
	docker push registry.example.com/vllm-proxy:$(VERSION)

benchmark:
	go test -bench=. -benchmem ./...

deps:
	go mod download
	go mod tidy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary (requires network)"
	@echo "  vendor-build   - Build using vendor directory (offline)"
	@echo "  vendor-deps    - Clone dependencies to vendor directory"
	@echo "  run            - Run with command line arguments"
	@echo "  run-config     - Run with config file"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-push    - Push Docker image"
	@echo "  benchmark      - Run benchmarks"
	@echo "  deps           - Download dependencies (requires network)"
