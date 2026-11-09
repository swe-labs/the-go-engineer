package main

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestCreateAndReadFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "nested", "dir", "test.txt")
	testData := []byte("Hello, Go Filesystem!\r\nLine 2\n")

	// 1. Create file in nested directory
	err := CreateFile(filePath, testData, 0644)
	if err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}

	// 2. Read file back
	info, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if info.Size != int64(len(testData)) {
		t.Errorf("expected size %d, got %d", len(testData), info.Size)
	}
	if info.IsDir {
		t.Errorf("expected IsDir to be false")
	}
	if !bytes.Equal(info.RawBytes, testData) {
		t.Errorf("raw bytes mismatch: got %v, expected %v", info.RawBytes, testData)
	}
	if !bytes.Contains(info.RawBytes, []byte("\r\n")) {
		t.Errorf("expected CRLF bytes in file")
	}

	// 3. Test non-existent file
	_, err = ReadFile(filepath.Join(tempDir, "missing.txt"))
	if err == nil {
		t.Errorf("expected error reading missing file")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
