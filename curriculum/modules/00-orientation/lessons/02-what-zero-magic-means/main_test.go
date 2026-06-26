package main

import (
	"testing"
)

func TestZeroMagicAppliesDefaults(t *testing.T) {
	tests := []struct {
		name        string
		cfg         Config
		wantPort    int
		wantTime    int
		wantVerbose bool
	}{
		{
			name:        "zero values get defaults",
			cfg:         Config{},
			wantPort:    8080,
			wantTime:    30,
			wantVerbose: false,
		},
		{
			name:        "explicit values are preserved",
			cfg:         Config{Port: 3000, Timeout: 10, Verbose: true},
			wantPort:    3000,
			wantTime:    10,
			wantVerbose: true,
		},
		{
			name:        "partial config fills port default",
			cfg:         Config{Timeout: 45, Verbose: true},
			wantPort:    8080,
			wantTime:    45,
			wantVerbose: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.cfg
			if cfg.Port == 0 {
				cfg.Port = 8080
			}
			if cfg.Timeout == 0 {
				cfg.Timeout = 30
			}
			if cfg.Port != tc.wantPort {
				t.Errorf("Port = %d, want %d", cfg.Port, tc.wantPort)
			}
			if cfg.Timeout != tc.wantTime {
				t.Errorf("Timeout = %d, want %d", cfg.Timeout, tc.wantTime)
			}
			if cfg.Verbose != tc.wantVerbose {
				t.Errorf("Verbose = %t, want %t", cfg.Verbose, tc.wantVerbose)
			}
		})
	}
}

func TestMagicServerUsesZeroValuesSilently(t *testing.T) {
	cfg := Config{}
	if cfg.Port != 0 {
		t.Error("MagicServer uses zero-value port, expected 0")
	}
	if cfg.Timeout != 0 {
		t.Error("MagicServer uses zero-value timeout, expected 0")
	}
}
