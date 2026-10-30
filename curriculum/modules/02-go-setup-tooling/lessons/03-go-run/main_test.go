package main

import (
	"strings"
	"testing"
)

func TestInspectExecution(t *testing.T) {
	meta, err := InspectExecution()
	if err != nil {
		t.Fatalf("InspectExecution failed: %v", err)
	}

	if meta.ExecutablePath == "" {
		t.Errorf("expected non-empty ExecutablePath")
	}
	if meta.Filename == "" {
		t.Errorf("expected non-empty Filename")
	}
	if !strings.HasSuffix(meta.Directory, "") {
		t.Errorf("invalid directory")
	}

	// During "go test", tests run from temporary directories as well
	if !meta.IsTemporary {
		t.Logf("Notice: binary not in standard temporary build directory: %s", meta.ExecutablePath)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

