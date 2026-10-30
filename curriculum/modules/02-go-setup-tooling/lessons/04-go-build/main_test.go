package main

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetBinaryInfo(t *testing.T) {
	info, err := GetBinaryInfo()
	if err != nil {
		t.Fatalf("GetBinaryInfo failed: %v", err)
	}

	if info.SizeBytes <= 0 {
		t.Errorf("binary size must be positive, got %d", info.SizeBytes)
	}
	if info.GOOS != runtime.GOOS {
		t.Errorf("expected GOOS %s, got %s", runtime.GOOS, info.GOOS)
	}
	if info.GOARCH != runtime.GOARCH {
		t.Errorf("expected GOARCH %s, got %s", runtime.GOARCH, info.GOARCH)
	}
	if info.Version != "dev" {
		t.Errorf("expected default version 'dev', got %s", info.Version)
	}
}

func TestFormatBinaryReport(t *testing.T) {
	dummy := BinaryInfo{
		Name:      "my-service",
		SizeBytes: 2048,
		GOOS:      "linux",
		GOARCH:    "amd64",
		Version:   "v1.2.3",
	}

	report := FormatBinaryReport(dummy)
	if !strings.Contains(report, "my-service") {
		t.Errorf("missing name in report: %s", report)
	}
	if !strings.Contains(report, "linux/amd64") {
		t.Errorf("missing target architecture in report: %s", report)
	}
	if !strings.Contains(report, "v1.2.3") {
		t.Errorf("missing version in report: %s", report)
	}
	if !strings.Contains(report, "2.00 KB") {
		t.Errorf("expected size formatting: %s", report)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

