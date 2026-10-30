package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Injected via -ldflags during build (e.g. -ldflags="-X main.Version=v1.0.0")
var Version = "dev"

// BinaryInfo holds characteristics of the compiled binary.
type BinaryInfo struct {
	Name      string
	Path      string
	SizeBytes int64
	GOOS      string
	GOARCH    string
	Version   string
}

// GetBinaryInfo inspects the currently running compiled binary on disk.
func GetBinaryInfo() (BinaryInfo, error) {
	execPath, err := os.Executable()
	if err != nil {
		return BinaryInfo{}, err
	}

	info, err := os.Stat(execPath)
	if err != nil {
		return BinaryInfo{}, err
	}

	return BinaryInfo{
		Name:      filepath.Base(execPath),
		Path:      execPath,
		SizeBytes: info.Size(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
		Version:   Version,
	}, nil
}

// FormatBinaryReport returns a human-readable summary of the binary.
func FormatBinaryReport(b BinaryInfo) string {
	return fmt.Sprintf(
		"Binary: %s\nTarget: %s/%s\nVersion: %s\nSize: %.2f KB (%d bytes)",
		b.Name,
		b.GOOS,
		b.GOARCH,
		b.Version,
		float64(b.SizeBytes)/1024.0,
		b.SizeBytes,
	)
}

func main() {
	fmt.Println("=== Deconstructing 'go build' Standalone Binary ===")

	b, err := GetBinaryInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inspecting binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(FormatBinaryReport(b))
	fmt.Println("\nInsight: 'go build' emits a self-contained, statically linked binary executable.")
	fmt.Println("It requires no external Go runtime on the host machine to execute.")
}
