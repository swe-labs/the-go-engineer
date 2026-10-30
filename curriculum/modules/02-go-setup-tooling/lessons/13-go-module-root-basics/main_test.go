package main

import (
	"testing"
)

func TestParseGoMod_Valid(t *testing.T) {
	input := `
// Auto-generated module file
module github.com/user/demo-service

go 1.22.4
`
	mod, err := ParseGoMod(input)
	if err != nil {
		t.Fatalf("ParseGoMod failed: %v", err)
	}

	if mod.ModulePath != "github.com/user/demo-service" {
		t.Errorf("expected module path 'github.com/user/demo-service', got %q", mod.ModulePath)
	}
	if mod.GoVersion != "1.22.4" {
		t.Errorf("expected Go version '1.22.4', got %q", mod.GoVersion)
	}
}

func TestParseGoMod_MissingModule(t *testing.T) {
	input := `go 1.22`
	_, err := ParseGoMod(input)
	if err == nil {
		t.Errorf("expected error when 'module' directive is missing")
	}
}

func TestResolveImportPath(t *testing.T) {
	tests := []struct {
		module   string
		sub      string
		expected string
	}{
		{"example.com/app", "pkg/auth", "example.com/app/pkg/auth"},
		{"example.com/app/", "/pkg/auth", "example.com/app/pkg/auth"},
		{"example.com/app", "", "example.com/app"},
	}

	for _, tt := range tests {
		got := ResolveImportPath(tt.module, tt.sub)
		if got != tt.expected {
			t.Errorf("ResolveImportPath(%q, %q) = %q, want %q", tt.module, tt.sub, got, tt.expected)
		}
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

