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

func NewConfigParserWithLookup(lookup func(string) (string, bool)) *ConfigParser {
	return &ConfigParser{
		envLookup: lookup,
	}
}

// Load loads configuration from defaults, merges file content, and applies env overrides.
func (p *ConfigParser) Load(filePath string) (AppConfig, error) {
	// TODO: Start with DefaultConfig()
	// TODO: If filePath is provided, read file and unmarshal based on extension (.json or .env)
	// TODO: Return ErrUnsupportedFormat for unrecognized extensions
	// TODO: Apply environment variable overrides via applyEnvOverrides
	// TODO: Run Validate(cfg) and return validation errors if any invariants fail
	return DefaultConfig(), errors.New("not implemented")
}

func (p *ConfigParser) parseKeyValueConfig(content string, cfg *AppConfig) error {
	// TODO: Split lines, parse KEY=VALUE, call applySingleKey
	return nil
}

func (p *ConfigParser) applyEnvOverrides(cfg *AppConfig) {
	// TODO: Map APP_* env variables (APP_SERVER_PORT, APP_DB_NAME, etc.) to config fields
}

func (p *ConfigParser) applySingleKey(key, val string, cfg *AppConfig) {
	// TODO: Set matching struct field based on key
}

func (p *ConfigParser) Validate(cfg AppConfig) error {
	// TODO: Check server.port is between 1 and 65535
	// TODO: Check database.name is non-empty
	// TODO: Check database.max_conns > 0
	return errors.New("not implemented")
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
