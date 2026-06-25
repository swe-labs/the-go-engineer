package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create temp file error: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	fmt.Printf("Temp file: %s\n", tmpFile.Name())
	tmpFile.Write([]byte("temporary data"))

	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "example-dir-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create temp dir error: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("Temp dir: %s\n", tmpDir)

	// Use temp dir
	subFile := filepath.Join(tmpDir, "data.txt")
	os.WriteFile(subFile, []byte("data in temp dir"), 0644)
	fmt.Printf("Wrote %s\n", subFile)

	// Check the system temp directory
	fmt.Printf("System temp dir: %s\n", os.TempDir())
}
