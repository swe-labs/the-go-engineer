package main

import (
	"runtime"
	"testing"
)

func TestGetEnvironmentInfo(t *testing.T) {
	info := GetEnvironmentInfo()

	if info.Version != runtime.Version() {
		t.Errorf("version mismatch: got %s, want %s", info.Version, runtime.Version())
	}
	if info.GOOS != runtime.GOOS {
		t.Errorf("GOOS mismatch: got %s, want %s", info.GOOS, runtime.GOOS)
	}
	if info.GOARCH != runtime.GOARCH {
		t.Errorf("GOARCH mismatch: got %s, want %s", info.GOARCH, runtime.GOARCH)
	}
	if info.NumCPU <= 0 {
		t.Errorf("logical CPUs must be positive, got %d", info.NumCPU)
	}
}

func TestVerifyInstallation(t *testing.T) {
	// 1. Valid info
	validInfo := GetEnvironmentInfo()
	ok, issues := VerifyInstallation(validInfo)
	if !ok {
		t.Fatalf("expected active environment to verify successfully, got issues: %v", issues)
	}

	// 2. Invalid info checks
	badInfo := EnvironmentInfo{
		Version: "unknown-lang",
		GOOS:    "",
		GOARCH:  "",
		NumCPU:  0,
	}
	badOk, badIssues := VerifyInstallation(badInfo)
	if badOk {
		t.Fatalf("expected bad info to fail verification")
	}
	if len(badIssues) != 4 {
		t.Errorf("expected 4 issues, got %d: %v", len(badIssues), badIssues)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

