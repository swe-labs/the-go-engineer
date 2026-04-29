package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileOrganizer(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "test-organizer-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files
	files := []string{"test1.txt", "test2.png", "test3.txt", "noextension"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", f, err)
		}
	}

	// Run the tool in normal mode
	cmd := exec.Command("go", "run", "main.go", "--dir="+tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	if !strings.Contains(outStr, "Done! 3 files moved.") {
		t.Errorf("Expected output to report 3 files moved, got:\n%s", outStr)
	}

	// Verify the files were moved
	if _, err := os.Stat(filepath.Join(tmpDir, "txt", "test1.txt")); os.IsNotExist(err) {
		t.Errorf("Expected txt/test1.txt to exist")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "png", "test2.png")); os.IsNotExist(err) {
		t.Errorf("Expected png/test2.png to exist")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "noextension")); os.IsNotExist(err) {
		t.Errorf("Expected 'noextension' to remain in root")
	}
}
