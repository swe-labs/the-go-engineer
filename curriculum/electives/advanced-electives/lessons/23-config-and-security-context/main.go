package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is the validated runtime configuration for the service.
type Config struct {
	Addr        string
	LogLevel    string
	MaxConns    int
	DatabaseURL string
}

// LoadConfig builds a Config from environment variables and a file-mounted
// secret. It returns the first validation error it finds so the process can
// fail fast at startup instead of serving with missing or invalid settings.
func LoadConfig(getenv func(string) string, readFile func(string) ([]byte, error)) (Config, error) {
	c := Config{
		Addr:     getenv("APP_ADDR"),
		LogLevel: getenv("APP_LOG_LEVEL"),
	}
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return Config{}, fmt.Errorf("APP_LOG_LEVEL %q is not one of debug, info, warn, error", c.LogLevel)
	}

	raw := getenv("APP_MAX_CONNS")
	if raw == "" {
		return Config{}, errors.New("APP_MAX_CONNS is required")
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return Config{}, fmt.Errorf("APP_MAX_CONNS %q must be a positive integer", raw)
	}
	c.MaxConns = n

	// A Secret mounted as a file is read at startup. The Deployment supplies the
	// path, so the same binary works whether the value comes from an environment
	// variable or a projected volume.
	path := getenv("APP_DATABASE_URL_FILE")
	if path == "" {
		return Config{}, errors.New("APP_DATABASE_URL_FILE is required")
	}
	data, err := readFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read secret %s: %w", path, err)
	}
	c.DatabaseURL = strings.TrimSpace(string(data))
	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("secret %s is empty", path)
	}
	return c, nil
}

func main() {
	cfg, err := LoadConfig(os.Getenv, os.ReadFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("listening on %s log=%s max_conns=%d\n", cfg.Addr, cfg.LogLevel, cfg.MaxConns)
}
