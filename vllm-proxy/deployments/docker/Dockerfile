# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /vllm-proxy ./cmd

# Runtime stage
FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /vllm-proxy /app/vllm-proxy
COPY --from=builder /app/configs /app/configs

EXPOSE 8000 9090

ENTRYPOINT ["/app/vllm-proxy"]
CMD ["--config", "/app/configs/config.yaml"]
