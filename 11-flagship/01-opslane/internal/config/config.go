// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
}

type AppConfig struct {
	Name     string
	Env      string
	LogLevel slog.Level
}

type HTTPConfig struct {
	Address           string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (Config, error) {
	return LoadFromLookup(os.LookupEnv)
}

func LoadFromLookup(lookup LookupFunc) (Config, error) {
	readHeaderTimeout, err := durationFromEnv(lookup, "OPSLANE_HTTP_READ_HEADER_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_HTTP_READ_HEADER_TIMEOUT: %w", err)
	}

	readTimeout, err := durationFromEnv(lookup, "OPSLANE_HTTP_READ_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_HTTP_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := durationFromEnv(lookup, "OPSLANE_HTTP_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_HTTP_WRITE_TIMEOUT: %w", err)
	}

	idleTimeout, err := durationFromEnv(lookup, "OPSLANE_HTTP_IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_HTTP_IDLE_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := durationFromEnv(lookup, "OPSLANE_HTTP_SHUTDOWN_TIMEOUT", 20*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_HTTP_SHUTDOWN_TIMEOUT: %w", err)
	}

	logLevel, err := parseLogLevel(stringFromEnv(lookup, "OPSLANE_LOG_LEVEL", "info"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		App: AppConfig{
			Name:     "opslane",
			Env:      stringFromEnv(lookup, "OPSLANE_ENV", "development"),
			LogLevel: logLevel,
		},
		HTTP: HTTPConfig{
			Address:           stringFromEnv(lookup, "OPSLANE_HTTP_ADDR", ":8080"),
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			ShutdownTimeout:   shutdownTimeout,
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	switch c.App.Env {
	case "development", "staging", "production", "test":
	default:
		return fmt.Errorf("invalid OPSLANE_ENV %q", c.App.Env)
	}

	if strings.TrimSpace(c.HTTP.Address) == "" {
		return fmt.Errorf("OPSLANE_HTTP_ADDR must not be empty")
	}

	if c.HTTP.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("OPSLANE_HTTP_READ_HEADER_TIMEOUT must be positive")
	}

	if c.HTTP.ReadTimeout <= 0 {
		return fmt.Errorf("OPSLANE_HTTP_READ_TIMEOUT must be positive")
	}

	if c.HTTP.WriteTimeout <= 0 {
		return fmt.Errorf("OPSLANE_HTTP_WRITE_TIMEOUT must be positive")
	}

	if c.HTTP.IdleTimeout <= 0 {
		return fmt.Errorf("OPSLANE_HTTP_IDLE_TIMEOUT must be positive")
	}

	if c.HTTP.ShutdownTimeout <= 0 {
		return fmt.Errorf("OPSLANE_HTTP_SHUTDOWN_TIMEOUT must be positive")
	}

	return nil
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid OPSLANE_LOG_LEVEL %q", value)
	}
}
