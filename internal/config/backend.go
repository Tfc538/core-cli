package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultBackendHost            = "127.0.0.1"
	defaultBackendPort            = 8080
	defaultBackendShutdownTimeout = 5 * time.Second
	defaultBackendReadTimeout     = 15 * time.Second
	defaultBackendWriteTimeout    = 15 * time.Second
	defaultBackendIdleTimeout     = 60 * time.Second
	defaultBackendReadHeader      = 5 * time.Second
)

// BackendConfig holds runtime configuration for core-backend.
type BackendConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ReadHeader      time.Duration
}

// Addr returns host:port for net/http server.
func (c BackendConfig) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// LoadBackend reads backend configuration from environment variables.
func LoadBackend() (BackendConfig, error) {
	cfg := BackendConfig{
		Host:            defaultBackendHost,
		Port:            defaultBackendPort,
		ShutdownTimeout: defaultBackendShutdownTimeout,
		ReadTimeout:     defaultBackendReadTimeout,
		WriteTimeout:    defaultBackendWriteTimeout,
		IdleTimeout:     defaultBackendIdleTimeout,
		ReadHeader:      defaultBackendReadHeader,
	}

	if host := os.Getenv("CORE_BACKEND_HOST"); host != "" {
		cfg.Host = host
	}

	if portStr := os.Getenv("CORE_BACKEND_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return BackendConfig{}, fmt.Errorf("invalid CORE_BACKEND_PORT: %w", err)
		}
		cfg.Port = port
	}

	if timeoutStr := os.Getenv("CORE_BACKEND_SHUTDOWN_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return BackendConfig{}, fmt.Errorf("invalid CORE_BACKEND_SHUTDOWN_TIMEOUT: %w", err)
		}
		cfg.ShutdownTimeout = timeout
	}

	if timeoutStr := os.Getenv("CORE_BACKEND_READ_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return BackendConfig{}, fmt.Errorf("invalid CORE_BACKEND_READ_TIMEOUT: %w", err)
		}
		cfg.ReadTimeout = timeout
	}

	if timeoutStr := os.Getenv("CORE_BACKEND_WRITE_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return BackendConfig{}, fmt.Errorf("invalid CORE_BACKEND_WRITE_TIMEOUT: %w", err)
		}
		cfg.WriteTimeout = timeout
	}

	if timeoutStr := os.Getenv("CORE_BACKEND_IDLE_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return BackendConfig{}, fmt.Errorf("invalid CORE_BACKEND_IDLE_TIMEOUT: %w", err)
		}
		cfg.IdleTimeout = timeout
	}

	if timeoutStr := os.Getenv("CORE_BACKEND_READ_HEADER_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return BackendConfig{}, fmt.Errorf("invalid CORE_BACKEND_READ_HEADER_TIMEOUT: %w", err)
		}
		cfg.ReadHeader = timeout
	}

	return cfg, nil
}
