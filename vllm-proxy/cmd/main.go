package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vllm-ascend/vllm-proxy/config"
	"github.com/vllm-ascend/vllm-proxy/internal/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	host           string
	port           int
	configPath     string
	prefillerHosts arrayFlags
	prefillerPorts arrayFlags
	decoderHosts   arrayFlags
	decoderPorts   arrayFlags
	maxRetries     int
	retryDelay     float64
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "Host address to bind the proxy server")
	flag.IntVar(&port, "port", 8000, "Port to bind the proxy server")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Var(&prefillerHosts, "prefiller-hosts", "Prefiller server hosts")
	flag.Var(&prefillerPorts, "prefiller-ports", "Prefiller server ports")
	flag.Var(&decoderHosts, "decoder-hosts", "Decoder server hosts")
	flag.Var(&decoderPorts, "decoder-ports", "Decoder server ports")
	flag.IntVar(&maxRetries, "max-retries", 3, "Maximum number of retries for HTTP requests")
	flag.Float64Var(&retryDelay, "retry-delay", 0.2, "Base delay (seconds) for exponential backoff retries")
}

func main() {
	flag.Parse()

	logger := initLogger()
	defer logger.Sync()

	var cfg *config.Config
	var err error

	if configPath != "" {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			logger.Fatal("Failed to load config", zap.Error(err))
		}
	} else if len(prefillerHosts) > 0 && len(decoderHosts) > 0 {
		cfg, err = config.LoadConfigFromArgs(host, port, prefillerHosts, prefillerPorts, decoderHosts, decoderPorts)
		if err != nil {
			logger.Fatal("Failed to parse arguments", zap.Error(err))
		}
	} else {
		cfg = config.DefaultConfig()
		cfg.Server.Host = host
		cfg.Server.Port = port
	}

	cfg.Retry.MaxRetries = maxRetries
	cfg.Retry.BaseDelay = time.Duration(retryDelay * float64(time.Second))

	logger.Info("Configuration loaded",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
		zap.Int("prefillers", len(cfg.Prefillers)),
		zap.Int("decoders", len(cfg.Decoders)),
	)

	proxyServer := server.NewProxyServer(cfg, logger)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		if err := proxyServer.Start(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		if err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	case sig := <-stopChan:
		logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
		if err := proxyServer.Shutdown(); err != nil {
			logger.Error("Error during shutdown", zap.Error(err))
		}
	}

	logger.Info("Server stopped")
}

func initLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
