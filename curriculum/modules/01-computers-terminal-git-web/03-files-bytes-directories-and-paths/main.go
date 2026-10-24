package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// FileInfo holds metadata and content for a file.
type FileInfo struct {
	Path        string
	Size        int64
	IsDir       bool
	Content     string
	RawBytes    []byte
	Permissions os.FileMode
}

// CreateFile writes bytes to disk, creating parent directories if necessary.
func CreateFile(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}
	if err := os.WriteFile(path, content, perm); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}
	return nil
}

// ReadFile retrieves content, metadata, and raw byte slice for a given path.
func ReadFile(path string) (*FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat failed: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}
	return &FileInfo{
		Path:        path,
		Size:        stat.Size(),
		IsDir:       stat.IsDir(),
		Content:     string(data),
		RawBytes:    data,
		Permissions: stat.Mode().Perm(),
	}, nil
}

func main() {
	tempDir, err := os.MkdirTemp("", "the-go-engineer-file-io-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // cleanup on exit

	samplePath := filepath.Join(tempDir, "data", "records.txt")
	// Notice explicit CRLF vs LF:
	sampleContent := []byte("Hello, byte-level file I/O!\r\nLine 2: 42 bytes.\n")

	if err := CreateFile(samplePath, sampleContent, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file: %v\n", err)
		os.Exit(1)
	}

	info, err := ReadFile(samplePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== File Metadata and Permissions ===")
	fmt.Printf("Path:        %s\n", info.Path)
	fmt.Printf("Size:        %d bytes\n", info.Size)
	fmt.Printf("Permissions: %o (%s)\n", info.Permissions, info.Permissions)
	fmt.Printf("Is Dir:      %t\n", info.IsDir)

	fmt.Println("\n=== Byte Inspection & Line Endings ===")
	fmt.Printf("Has Windows CRLF (\\r\\n): %t\n", bytes.Contains(info.RawBytes, []byte("\r\n")))
	fmt.Printf("Has Unix LF (\\n):         %t\n", bytes.Contains(info.RawBytes, []byte("\n")))
	fmt.Printf("First 10 Bytes (Hex):     % x\n", info.RawBytes[:10])

	fmt.Println("\n=== Standard Process File Descriptors ===")
	fmt.Printf("os.Stdin.Fd():  %d\n", os.Stdin.Fd())
	fmt.Printf("os.Stdout.Fd(): %d\n", os.Stdout.Fd())
	fmt.Printf("os.Stderr.Fd(): %d\n", os.Stderr.Fd())

	fmt.Println("\n=== Cross-Platform Path Manipulation ===")
	abs, _ := filepath.Abs(samplePath)
	fmt.Printf("filepath.Abs:  %s\n", abs)
	fmt.Printf("filepath.Base: %s\n", filepath.Base(samplePath))
	fmt.Printf("filepath.Dir:  %s\n", filepath.Dir(samplePath))
	fmt.Printf("filepath.Ext:  %s\n", filepath.Ext(samplePath))
}
