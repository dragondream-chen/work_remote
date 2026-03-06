package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	Prefillers     []InstanceConfig     `mapstructure:"prefillers"`
	Decoders       []InstanceConfig     `mapstructure:"decoders"`
	ConnectionPool ConnectionPoolConfig `mapstructure:"connection_pool"`
	Retry          RetryConfig          `mapstructure:"retry"`
	Logging        LoggingConfig        `mapstructure:"logging"`
	Metrics        MetricsConfig        `mapstructure:"metrics"`
}

type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	MaxConnections int           `mapstructure:"max_connections"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
}

type InstanceConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Weight int    `mapstructure:"weight"`
}

func (i InstanceConfig) Address() string {
	return fmt.Sprintf("%s:%d", i.Host, i.Port)
}

func (i InstanceConfig) URL() string {
	return fmt.Sprintf("http://%s:%d/v1", i.Host, i.Port)
}

type ConnectionPoolConfig struct {
	MaxIdleConns      int           `mapstructure:"max_idle_conns"`
	MaxConnsPerHost   int           `mapstructure:"max_conns_per_host"`
	IdleConnTimeout   time.Duration `mapstructure:"idle_conn_timeout"`
	HandshakeTimeout  time.Duration `mapstructure:"handshake_timeout"`
	ResponseHeaderTimeout time.Duration `mapstructure:"response_header_timeout"`
}

type RetryConfig struct {
	MaxRetries int           `mapstructure:"max_retries"`
	BaseDelay  time.Duration `mapstructure:"base_delay"`
	MaxDelay   time.Duration `mapstructure:"max_delay"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:           "0.0.0.0",
			Port:           8000,
			MaxConnections: 100000,
			RequestTimeout: 30 * time.Second,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			IdleTimeout:    120 * time.Second,
		},
		Prefillers: []InstanceConfig{},
		Decoders:   []InstanceConfig{},
		ConnectionPool: ConnectionPoolConfig{
			MaxIdleConns:        10000,
			MaxConnsPerHost:     1000,
			IdleConnTimeout:     90 * time.Second,
			HandshakeTimeout:    10 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
		Retry: RetryConfig{
			MaxRetries: 3,
			BaseDelay:  200 * time.Millisecond,
			MaxDelay:   5 * time.Second,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
			Port:    9090,
		},
	}
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("/etc/vllm-proxy")
	}

	viper.SetEnvPrefix("VLLM_PROXY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return cfg, nil
}

func LoadConfigFromArgs(host string, port int, prefillerHosts, prefillerPorts, decoderHosts, decoderPorts []string) (*Config, error) {
	cfg := DefaultConfig()

	cfg.Server.Host = host
	cfg.Server.Port = port

	if len(prefillerHosts) != len(prefillerPorts) {
		return nil, fmt.Errorf("number of prefiller hosts (%d) must match number of prefiller ports (%d)",
			len(prefillerHosts), len(prefillerPorts))
	}

	if len(decoderHosts) != len(decoderPorts) {
		return nil, fmt.Errorf("number of decoder hosts (%d) must match number of decoder ports (%d)",
			len(decoderHosts), len(decoderPorts))
	}

	for i := range prefillerHosts {
		var port int
		fmt.Sscanf(prefillerPorts[i], "%d", &port)
		cfg.Prefillers = append(cfg.Prefillers, InstanceConfig{
			Host:   prefillerHosts[i],
			Port:   port,
			Weight: 1,
		})
	}

	for i := range decoderHosts {
		var port int
		fmt.Sscanf(decoderPorts[i], "%d", &port)
		cfg.Decoders = append(cfg.Decoders, InstanceConfig{
			Host:   decoderHosts[i],
			Port:   port,
			Weight: 1,
		})
	}

	return cfg, nil
}
