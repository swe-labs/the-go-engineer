package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	ErrValidation        = errors.New("configuration validation failed")
	ErrUnsupportedFormat = errors.New("unsupported configuration format")
)

type ServerConfig struct {
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	ReadTimeout time.Duration `json:"read_timeout"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	MaxConns int    `json:"max_conns"`
}

type MetricsConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type AppConfig struct {
	Environment string         `json:"environment"`
	Server      ServerConfig   `json:"server"`
	Database    DatabaseConfig `json:"database"`
	Metrics     MetricsConfig  `json:"metrics"`
}

// DefaultConfig returns baseline configuration values.
func DefaultConfig() AppConfig {
	return AppConfig{
		Environment: "development",
		Server: ServerConfig{
			Host:        "0.0.0.0",
			Port:        8080,
			ReadTimeout: 5 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "app_dev",
			User:     "postgres",
			MaxConns: 25,
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
		},
	}
}

type ConfigParser struct {
	envLookup func(string) (string, bool)
}

func NewConfigParser() *ConfigParser {
	return &ConfigParser{
		envLookup: os.LookupEnv,
	}
}

// WithEnvLookup allows mocking environment variables in tests.
func NewConfigParserWithLookup(lookup func(string) (string, bool)) *ConfigParser {
	return &ConfigParser{
		envLookup: lookup,
	}
}

// Load loads configuration from defaults, merges file content, and applies env overrides.
func (p *ConfigParser) Load(filePath string) (AppConfig, error) {
	cfg := DefaultConfig()

	// 1. Merge File if specified
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return cfg, fmt.Errorf("failed to read config file: %w", err)
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".json":
			if err := json.Unmarshal(content, &cfg); err != nil {
				return cfg, fmt.Errorf("failed to parse json config: %w", err)
			}
		case ".env", ".kv":
			if err := p.parseKeyValueConfig(string(content), &cfg); err != nil {
				return cfg, fmt.Errorf("failed to parse kv config: %w", err)
			}
		default:
			return cfg, fmt.Errorf("%w: %s (must be .json or .env)", ErrUnsupportedFormat, ext)
		}
	}

	// 2. Apply Environment Variable Overrides
	p.applyEnvOverrides(&cfg)

	// 3. Validate Configuration Invariants
	if err := p.Validate(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (p *ConfigParser) parseKeyValueConfig(content string, cfg *AppConfig) error {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed kv line %d: %q", i+1, line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		p.applySingleKey(key, val, cfg)
	}
	return nil
}

func (p *ConfigParser) applyEnvOverrides(cfg *AppConfig) {
	overrides := map[string]string{
		"APP_ENV":             "environment",
		"APP_SERVER_HOST":     "server.host",
		"APP_SERVER_PORT":     "server.port",
		"APP_DB_HOST":         "database.host",
		"APP_DB_PORT":         "database.port",
		"APP_DB_NAME":         "database.name",
		"APP_DB_USER":         "database.user",
		"APP_DB_PASSWORD":     "database.password",
		"APP_DB_MAX_CONNS":    "database.max_conns",
		"APP_METRICS_ENABLED": "metrics.enabled",
		"APP_METRICS_PATH":    "metrics.path",
	}

	for envKey, configKey := range overrides {
		if val, exists := p.envLookup(envKey); exists && val != "" {
			p.applySingleKey(configKey, val, cfg)
		}
	}
}

func (p *ConfigParser) applySingleKey(key, val string, cfg *AppConfig) {
	switch strings.ToLower(key) {
	case "environment":
		cfg.Environment = val
	case "server.host":
		cfg.Server.Host = val
	case "server.port":
		if p, err := strconv.Atoi(val); err == nil {
			cfg.Server.Port = p
		}
	case "database.host":
		cfg.Database.Host = val
	case "database.port":
		if p, err := strconv.Atoi(val); err == nil {
			cfg.Database.Port = p
		}
	case "database.name":
		cfg.Database.Name = val
	case "database.user":
		cfg.Database.User = val
	case "database.password":
		cfg.Database.Password = val
	case "database.max_conns":
		if m, err := strconv.Atoi(val); err == nil {
			cfg.Database.MaxConns = m
		}
	case "metrics.enabled":
		cfg.Metrics.Enabled = strings.EqualFold(val, "true") || val == "1"
	case "metrics.path":
		cfg.Metrics.Path = val
	}
}

func (p *ConfigParser) Validate(cfg AppConfig) error {
	var errs []string

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		errs = append(errs, fmt.Sprintf("server.port must be between 1 and 65535 (got %d)", cfg.Server.Port))
	}
	if strings.TrimSpace(cfg.Database.Name) == "" {
		errs = append(errs, "database.name cannot be empty")
	}
	if cfg.Database.MaxConns <= 0 {
		errs = append(errs, "database.max_conns must be positive")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidation, strings.Join(errs, "; "))
	}
	return nil
}

func main() {
	parser := NewConfigParser()
	cfg, err := parser.Load("")
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded Environment: %s\n", cfg.Environment)
	fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: %s@%s:%d/%s (max conns: %d)\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.MaxConns)
}
