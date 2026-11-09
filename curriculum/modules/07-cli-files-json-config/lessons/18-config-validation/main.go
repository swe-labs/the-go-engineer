package main

import (
	"errors"
	"fmt"
	"net/url"
)

type ServerConfig struct {
	Port     int
	Host     string
	DBURL    string
	LogLevel string
	Workers  int
}

func ValidateConfig(cfg ServerConfig) error {
	var errs []error

	if cfg.Port < 1 || cfg.Port > 65535 {
		errs = append(errs, fmt.Errorf("port %d out of range [1, 65535]", cfg.Port))
	}

	if cfg.Host == "" {
		errs = append(errs, errors.New("host is required"))
	}

	if cfg.DBURL == "" {
		errs = append(errs, errors.New("database_url is required"))
	} else {
		if _, err := url.Parse(cfg.DBURL); err != nil {
			errs = append(errs, fmt.Errorf("database_url is not a valid URL: %w", err))
		}
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[cfg.LogLevel] {
		errs = append(errs, fmt.Errorf("log_level %q must be one of: debug, info, warn, error", cfg.LogLevel))
	}

	if cfg.Workers < 1 {
		errs = append(errs, fmt.Errorf("workers must be at least 1, got %d", cfg.Workers))
	}

	return errors.Join(errs...)
}

func main() {
	// Invalid config with multiple errors
	bad := ServerConfig{
		Port:     99999,
		Host:     "",
		DBURL:    "not-a-url",
		LogLevel: "critical",
		Workers:  0,
	}
	if err := ValidateConfig(bad); err != nil {
		fmt.Println("Invalid config errors:")
		fmt.Println(err)
	}

	fmt.Println()

	// Valid config
	good := ServerConfig{
		Port:     8080,
		Host:     "0.0.0.0",
		DBURL:    "postgres://localhost:5432/db?sslmode=disable",
		LogLevel: "info",
		Workers:  4,
	}
	if err := ValidateConfig(good); err != nil {
		fmt.Println("Unexpected error:", err)
	} else {
		fmt.Println("Valid config: no errors")
	}
}
