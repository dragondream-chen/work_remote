.PHONY: build run clean test docker-build docker-push

BINARY_NAME=vllm-proxy
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd

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
	@echo "  build          - Build the binary"
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
	@echo "  deps           - Download dependencies"
