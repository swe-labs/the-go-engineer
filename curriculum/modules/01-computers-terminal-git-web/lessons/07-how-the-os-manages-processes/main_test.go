package main

import (
	"os"
	"runtime"
	"testing"
)

func TestProcessMetadata(t *testing.T) {
	meta := ProcessMetadata{
		PID:  os.Getpid(),
		PPID: os.Getppid(),
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	if meta.PID <= 0 {
		t.Errorf("expected positive PID, got %d", meta.PID)
	}
	if meta.OS == "" {
		t.Errorf("expected non-empty OS string")
	}
	if meta.Arch == "" {
		t.Errorf("expected non-empty Arch string")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

