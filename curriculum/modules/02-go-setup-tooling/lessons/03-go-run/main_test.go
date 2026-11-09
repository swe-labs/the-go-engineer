package main

import (
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
	if meta.Directory == "" || meta.Directory == "." {
		t.Errorf("expected an absolute directory, got %q", meta.Directory)
	}

	// During "go test", tests run from temporary directories as well
	if !meta.IsTemporary {
		t.Logf("Notice: binary not in standard temporary build directory: %s", meta.ExecutablePath)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
