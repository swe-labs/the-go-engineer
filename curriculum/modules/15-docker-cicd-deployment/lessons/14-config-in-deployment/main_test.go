package main

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfigFromEnvDefaults(t *testing.T) {
	os.Clearenv()
	cfg, errs := LoadConfigFromEnv()
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log level info, got %s", cfg.LogLevel)
	}
	if cfg.MaxConnections != 100 {
		t.Errorf("expected default max connections 100, got %d", cfg.MaxConnections)
	}
	if cfg.Environment != "development" {
		t.Errorf("expected default environment development, got %s", cfg.Environment)
	}
}

func TestLoadConfigFromEnvCustom(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("MAX_CONNECTIONS", "200")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("FEATURE_FLAGS", "a,b,c")

	cfg, errs := LoadConfigFromEnv()
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://localhost/mydb" {
		t.Errorf("unexpected database URL: %s", cfg.DatabaseURL)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected debug, got %s", cfg.LogLevel)
	}
	if cfg.MaxConnections != 200 {
		t.Errorf("expected 200, got %d", cfg.MaxConnections)
	}
	if cfg.Environment != "production" {
		t.Errorf("expected production, got %s", cfg.Environment)
	}
	if len(cfg.FeatureFlags) != 3 {
		t.Errorf("expected 3 feature flags, got %d", len(cfg.FeatureFlags))
	}
}

func TestLoadConfigInvalidPort(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "notanumber")
	_, errs := LoadConfigFromEnv()
	found := false
	for _, e := range errs {
		if ce, ok := e.(ConfigError); ok && ce.Field == "PORT" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ConfigError for PORT, got %v", errs)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr int
	}{
		{
			name:    "valid",
			cfg:     &Config{Port: 8080, DatabaseURL: "postgres://localhost/db", LogLevel: "info", MaxConnections: 100, Environment: "development"},
			wantErr: 0,
		},
		{
			name:    "invalid port",
			cfg:     &Config{Port: 0, DatabaseURL: "postgres://localhost/db", LogLevel: "info", MaxConnections: 100, Environment: "development"},
			wantErr: 1,
		},
		{
			name:    "missing database",
			cfg:     &Config{Port: 8080, DatabaseURL: "", LogLevel: "info", MaxConnections: 100, Environment: "development"},
			wantErr: 1,
		},
		{
			name:    "invalid log level",
			cfg:     &Config{Port: 8080, DatabaseURL: "postgres://localhost/db", LogLevel: "trace", MaxConnections: 100, Environment: "development"},
			wantErr: 1,
		},
		{
			name:    "invalid max connections",
			cfg:     &Config{Port: 8080, DatabaseURL: "postgres://localhost/db", LogLevel: "info", MaxConnections: 0, Environment: "development"},
			wantErr: 1,
		},
		{
			name:    "invalid environment",
			cfg:     &Config{Port: 8080, DatabaseURL: "postgres://localhost/db", LogLevel: "info", MaxConnections: 100, Environment: "prod"},
			wantErr: 1,
		},
		{
			name:    "multiple errors",
			cfg:     &Config{Port: 0, DatabaseURL: "", LogLevel: "trace", MaxConnections: 0, Environment: "prod"},
			wantErr: 5,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateConfig(tc.cfg)
			if len(errs) != tc.wantErr {
				t.Errorf("expected %d errors, got %d: %v", tc.wantErr, len(errs), errs)
			}
		})
	}
}

func TestConfigStringMaskDSN(t *testing.T) {
	cfg := &Config{Port: 8080, DatabaseURL: "postgres://user:secret@localhost:5432/app", LogLevel: "info", MaxConnections: 100, Environment: "production"}
	s := cfg.String()
	if strings.Contains(s, "secret") {
		t.Errorf("config string should not contain password: %s", s)
	}
}

func TestMaskDSN(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"postgres://user:pass@localhost/db", "postgres://user:****@localhost/db"},
		{"postgres://localhost/db", "postgres://localhost/db"},
	}
	for _, tc := range tests {
		got := maskDSN(tc.input)
		if got != tc.want {
			t.Errorf("maskDSN(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
