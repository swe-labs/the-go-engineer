package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExecutionMetadata stores details about how the binary was executed.
type ExecutionMetadata struct {
	ExecutablePath string
	IsTemporary    bool
	Directory      string
	Filename       string
}

// InspectExecution determines if the program is running from a temporary build directory.
func InspectExecution() (ExecutionMetadata, error) {
	execPath, err := os.Executable()
	if err != nil {
		return ExecutionMetadata{}, err
	}

	cleanPath := filepath.Clean(execPath)
	dir := filepath.Dir(cleanPath)
	file := filepath.Base(cleanPath)

	// "go run" places compiled binaries in a temporary build directory (e.g. go-build or OS Temp)
	tempDir := os.TempDir()
	isTemp := strings.Contains(cleanPath, "go-build") || strings.HasPrefix(cleanPath, tempDir)

	return ExecutionMetadata{
		ExecutablePath: cleanPath,
		IsTemporary:    isTemp,
		Directory:      dir,
		Filename:       file,
	}, nil
}

func main() {
	fmt.Println("=== Deconstructing 'go run' Execution ===")

	meta, err := InspectExecution()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to inspect executable: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Active Executable Path: %s\n", meta.ExecutablePath)
	fmt.Printf("Binary Filename:        %s\n", meta.Filename)
	fmt.Printf("Directory Location:     %s\n", meta.Directory)
	fmt.Printf("Is Temporary Build Dir: %t\n", meta.IsTemporary)

	if meta.IsTemporary {
		fmt.Println("\nInsight: 'go run' compiled this code into a temporary binary and executed it directly.")
		fmt.Println("When this process exits, the temporary directory will be cleaned up.")
	} else {
		fmt.Println("\nInsight: Running from a persistent compiled binary on disk.")
	}
}
