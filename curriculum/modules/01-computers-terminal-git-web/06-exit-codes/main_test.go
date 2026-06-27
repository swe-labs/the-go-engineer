package main

import (
	"testing"
)

func TestExitCodeByFileStatus(t *testing.T) {
	tests := []struct {
		name       string
		fileExists bool
		want       int
	}{
		{"file exists", true, 0},
		{"file does not exist", false, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCodeByFileStatus(tt.fileExists); got != tt.want {
				t.Errorf("ExitCodeByFileStatus(%v) = %d, want %d", tt.fileExists, got, tt.want)
			}
		})
	}
}

func TestExitCodeByFileStatusEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		fileExists bool
		want       int
	}{
		{"empty struct false", false, 1},
		{"empty struct true", true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCodeByFileStatus(tt.fileExists); got != tt.want {
				t.Errorf("ExitCodeByFileStatus(%v) = %d, want %d", tt.fileExists, got, tt.want)
			}
		})
	}
}
