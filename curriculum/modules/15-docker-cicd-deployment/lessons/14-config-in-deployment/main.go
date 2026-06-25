package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           int
	DatabaseURL    string
	LogLevel       string
	MaxConnections int
	Environment    string
	FeatureFlags   []string
}

type ConfigError struct {
	Field string
	Err   error
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config %s: %v", e.Field, e.Err)
}

func LoadConfigFromEnv() (*Config, []error) {
	var errs []error
	cfg := &Config{
		Port:           8080,
		LogLevel:       "info",
		MaxConnections: 100,
		Environment:    "development",
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, ConfigError{"PORT", err})
		} else {
			cfg.Port = p
		}
	}

	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("MAX_CONNECTIONS"); v != "" {
		m, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, ConfigError{"MAX_CONNECTIONS", err})
		} else {
			cfg.MaxConnections = m
		}
	}
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		cfg.Environment = v
	}
	if v := os.Getenv("FEATURE_FLAGS"); v != "" {
		cfg.FeatureFlags = strings.Split(v, ",")
	}

	return cfg, errs
}

func ValidateConfig(cfg *Config) []error {
	var errs []error
	if cfg.Port < 1 || cfg.Port > 65535 {
		errs = append(errs, ConfigError{"Port", fmt.Errorf("must be between 1 and 65535, got %d", cfg.Port)})
	}
	if cfg.DatabaseURL == "" {
		errs = append(errs, ConfigError{"DatabaseURL", fmt.Errorf("database URL is required")})
	}
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[cfg.LogLevel] {
		errs = append(errs, ConfigError{"LogLevel", fmt.Errorf("invalid log level: %s", cfg.LogLevel)})
	}
	if cfg.MaxConnections < 1 {
		errs = append(errs, ConfigError{"MaxConnections", fmt.Errorf("must be at least 1")})
	}
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[cfg.Environment] {
		errs = append(errs, ConfigError{"Environment", fmt.Errorf("invalid environment: %s", cfg.Environment)})
	}
	return errs
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{Port=%d, DatabaseURL=%q, LogLevel=%q, MaxConnections=%d, Environment=%q, FeatureFlags=%v}",
		c.Port, maskDSN(c.DatabaseURL), c.LogLevel, c.MaxConnections, c.Environment, c.FeatureFlags,
	)
}

func maskDSN(dsn string) string {
	if dsn == "" {
		return ""
	}
	atIdx := strings.LastIndex(dsn, "@")
	if atIdx < 0 {
		return dsn
	}
	protoEnd := strings.Index(dsn, "://")
	if protoEnd < 0 {
		protoEnd = -3
	}
	credentials := dsn[protoEnd+3 : atIdx]
	colonIdx := strings.IndexByte(credentials, ':')
	if colonIdx >= 0 {
		return dsn[:protoEnd+3+colonIdx+1] + "****" + dsn[atIdx:]
	}
	return dsn
}

func main() {
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/app?sslmode=disable")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("MAX_CONNECTIONS", "50")
	os.Setenv("ENVIRONMENT", "development")
	os.Setenv("FEATURE_FLAGS", "new-checkout,dark-mode")

	cfg, errs := LoadConfigFromEnv()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "config error: %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Println("Loaded config:")
	fmt.Println(cfg)

	valErrs := ValidateConfig(cfg)
	if len(valErrs) > 0 {
		for _, e := range valErrs {
			fmt.Fprintf(os.Stderr, "validation error: %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Println("Config validation: OK")
}
