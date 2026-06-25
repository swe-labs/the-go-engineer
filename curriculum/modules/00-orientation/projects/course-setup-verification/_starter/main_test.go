package main

import (
	"strings"
	"testing"
)

func TestCheckTool(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		path     string
		wantOK   bool
	}{
		{
			name:     "go is installed",
			toolName: "Go",
			path:     "go",
			wantOK:   true,
		},
		{
			name:     "nonexistent tool not found",
			toolName: "Nonexistent",
			path:     "this-tool-does-not-exist-12345",
			wantOK:   false,
		},
		{
			name:     "git is installed",
			toolName: "Git",
			path:     "git",
			wantOK:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkTool(tc.toolName, tc.path)
			if tc.wantOK && !strings.HasPrefix(result, "OK") {
				t.Errorf("expected OK, got: %s", result)
			}
			if !tc.wantOK && !strings.HasPrefix(result, "NOT FOUND") {
				t.Errorf("expected NOT FOUND, got: %s", result)
			}
		})
	}
}
