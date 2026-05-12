// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

// Package config provides configuration loading and validation for the Opslane server.
// It reads configuration from environment variables with sensible defaults and validates
// all settings before returning to ensure safe startup.
package config

import (
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"strings"
	"time"
)

// Config (Struct): aggregates all configuration sub-sections needed to run the server.
// It is loaded via Load() which validates all settings before returning.
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Auth     AuthConfig
	OTEL     OTELConfig
}

// AppConfig (Struct): defines application-level settings including service identity,
// runtime environment, and logging verbosity.
type AppConfig struct {
	Name     string
	Env      string
	LogLevel slog.Level
}

// HTTPConfig (Struct): defines HTTP server parameters including address, timeouts,
// and trusted proxy CIDRs for proper remote IP detection behind load balancers.
type HTTPConfig struct {
	Address           string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	TrustedProxyCIDRs []netip.Prefix
}

// DatabaseConfig (Struct): defines PostgreSQL connection parameters including DSN,
// connection pool sizing, and connection lifetime settings.
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

// AuthConfig (Struct): defines authentication token settings including secret key,
// issuer identifier, and token time-to-live duration.
type AuthConfig struct {
	TokenSecret string
	TokenIssuer string
	TokenTTL    time.Duration
}

// OTELConfig (Struct): defines OpenTelemetry tracer settings including endpoint,
// security mode, timeout, and sampling rate. These values are passed to the
// otel.Tracer constructor at startup.
type OTELConfig struct {
	Endpoint   string
	Insecure   bool
	Enabled    bool
	Timeout    time.Duration
	SampleRate float64
}

const defaultDevelopmentTokenSecret = "development-only-opslane-secret-change-me"

// Load (Function): reads configuration from environment variables using os.LookupEnv.
// It is a convenience wrapper around LoadFromLookup for standard usage.
func Load() (Config, error) {
	return LoadFromLookup(os.LookupEnv)
}

// LoadFromLookup (Function): reads configuration from the provided lookup function,
// allowing dependency injection for testing. It validates all settings before returning.
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

	trustedProxyCIDRs, err := cidrPrefixesFromEnv(lookup, "OPSLANE_HTTP_TRUSTED_PROXY_CIDRS")
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_HTTP_TRUSTED_PROXY_CIDRS: %w", err)
	}

	maxOpenConns, err := intFromEnv(lookup, "OPSLANE_DB_MAX_OPEN_CONNS", 4)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_DB_MAX_OPEN_CONNS: %w", err)
	}

	maxIdleConns, err := intFromEnv(lookup, "OPSLANE_DB_MAX_IDLE_CONNS", 2)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_DB_MAX_IDLE_CONNS: %w", err)
	}

	connMaxIdleTime, err := durationFromEnv(lookup, "OPSLANE_DB_CONN_MAX_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_DB_CONN_MAX_IDLE_TIME: %w", err)
	}

	connMaxLifetime, err := durationFromEnv(lookup, "OPSLANE_DB_CONN_MAX_LIFETIME", 30*time.Minute)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_DB_CONN_MAX_LIFETIME: %w", err)
	}

	tokenTTL, err := durationFromEnv(lookup, "OPSLANE_AUTH_TOKEN_TTL", time.Hour)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_AUTH_TOKEN_TTL: %w", err)
	}

	logLevel, err := parseLogLevel(stringFromEnv(lookup, "OPSLANE_LOG_LEVEL", "info"))
	if err != nil {
		return Config{}, err
	}

	otlpTimeout, err := durationFromEnv(lookup, "OPSLANE_OTEL_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_OTEL_TIMEOUT: %w", err)
	}

	otlpSampleRate, err := floatFromEnv(lookup, "OPSLANE_OTEL_SAMPLE_RATE", 1.0)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPSLANE_OTEL_SAMPLE_RATE: %w", err)
	}

	otelEndpoint := stringFromEnv(lookup, "OPSLANE_OTEL_ENDPOINT", "")

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
			TrustedProxyCIDRs: trustedProxyCIDRs,
		},
		Database: DatabaseConfig{
			DSN:             stringFromEnv(lookup, "OPSLANE_DB_DSN", "postgres://opslane:secretpassword@localhost:5432/opslane?sslmode=disable"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxIdleTime: connMaxIdleTime,
			ConnMaxLifetime: connMaxLifetime,
		},
		Auth: AuthConfig{
			TokenSecret: stringFromEnv(lookup, "OPSLANE_AUTH_TOKEN_SECRET", defaultDevelopmentTokenSecret),
			TokenIssuer: stringFromEnv(lookup, "OPSLANE_AUTH_TOKEN_ISSUER", "opslane"),
			TokenTTL:    tokenTTL,
		},
		OTEL: OTELConfig{
			Endpoint:   otelEndpoint,
			Insecure:   stringFromEnv(lookup, "OPSLANE_OTEL_INSECURE", "") == "true",
			Enabled:    otelEndpoint != "",
			Timeout:    otlpTimeout,
			SampleRate: otlpSampleRate,
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate (Method): checks that all configuration values are within acceptable ranges
// and meet security requirements (e.g., production token secret must be different from
// the default development secret). Returns an error if any validation fails.
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

	for _, prefix := range c.HTTP.TrustedProxyCIDRs {
		if !prefix.IsValid() {
			return fmt.Errorf("OPSLANE_HTTP_TRUSTED_PROXY_CIDRS contains an invalid prefix")
		}
	}

	if strings.TrimSpace(c.Database.DSN) == "" {
		return fmt.Errorf("OPSLANE_DB_DSN must not be empty")
	}

	if c.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("OPSLANE_DB_MAX_OPEN_CONNS must be positive")
	}

	if c.Database.MaxIdleConns < 0 {
		return fmt.Errorf("OPSLANE_DB_MAX_IDLE_CONNS must be zero or positive")
	}

	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("OPSLANE_DB_MAX_IDLE_CONNS must not exceed OPSLANE_DB_MAX_OPEN_CONNS")
	}

	if c.Database.ConnMaxIdleTime <= 0 {
		return fmt.Errorf("OPSLANE_DB_CONN_MAX_IDLE_TIME must be positive")
	}

	if c.Database.ConnMaxLifetime <= 0 {
		return fmt.Errorf("OPSLANE_DB_CONN_MAX_LIFETIME must be positive")
	}

	if strings.TrimSpace(c.Auth.TokenSecret) == "" {
		return fmt.Errorf("OPSLANE_AUTH_TOKEN_SECRET must not be empty")
	}

	if len(c.Auth.TokenSecret) < 32 {
		return fmt.Errorf("OPSLANE_AUTH_TOKEN_SECRET must be at least 32 characters")
	}

	if c.App.Env == "production" && c.Auth.TokenSecret == defaultDevelopmentTokenSecret {
		return fmt.Errorf("OPSLANE_AUTH_TOKEN_SECRET must be changed in production")
	}

	if strings.TrimSpace(c.Auth.TokenIssuer) == "" {
		return fmt.Errorf("OPSLANE_AUTH_TOKEN_ISSUER must not be empty")
	}

	if c.Auth.TokenTTL <= 0 {
		return fmt.Errorf("OPSLANE_AUTH_TOKEN_TTL must be positive")
	}

	if c.OTEL.SampleRate < 0 || c.OTEL.SampleRate > 1 {
		return fmt.Errorf("OPSLANE_OTEL_SAMPLE_RATE must be between 0 and 1, got %f", c.OTEL.SampleRate)
	}

	if c.OTEL.Timeout <= 0 {
		return fmt.Errorf("OPSLANE_OTEL_TIMEOUT must be positive")
	}

	return nil
}

// parseLogLevel (Function): converts a string log level (debug, info, warn, error)
// into an slog.Level value. It is an unexported helper used by LoadFromLookup.
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
