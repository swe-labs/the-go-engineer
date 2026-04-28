package main

import "testing"

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{"Trace", LevelTrace, "Trace"},
		{"Debug", LevelDebug, "Debug"},
		{"Info", LevelInfo, "Info"},
		{"Warning", LevelWarning, "Warning"},
		{"Error", LevelError, "Error"},
		{"Unknown (too high)", 10, "Unknown"},
		{"Unknown (negative)", -1, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel(%d).String() = %q, want %q", tt.level, got, tt.expected)
			}
		})
	}
}
