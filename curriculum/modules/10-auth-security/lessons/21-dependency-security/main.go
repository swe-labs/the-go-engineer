package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func hashModule(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:]), nil
}

type ModuleEntry struct {
	Path    string
	Version string
	Hash    string
}

func parseGoSum(sumFile string) ([]ModuleEntry, error) {
	data, err := os.ReadFile(sumFile)
	if err != nil {
		return nil, err
	}
	var entries []ModuleEntry
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			entries = append(entries, ModuleEntry{
				Path:    parts[0],
				Version: parts[1],
				Hash:    parts[2],
			})
		}
	}
	return entries, nil
}

func minifyDeps(deps []string) []string {
	// In a real scenario, this would analyze imports.
	// For this lesson, we simulate by removing deps that start with "unused-".
	var result []string
	for _, d := range deps {
		if !strings.HasPrefix(d, "unused-") {
			result = append(result, d)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Dependency Security Demo ===")

	fmt.Println("\n1. Go's checksum database (simulated)")
	entries := []ModuleEntry{
		{Path: "github.com/gorilla/mux", Version: "v1.8.1", Hash: "abc123def456..."},
		{Path: "github.com/lib/pq", Version: "v1.10.9", Hash: "789012ghi345..."},
		{Path: "golang.org/x/crypto", Version: "v0.17.0", Hash: "567890jkl123..."},
	}
	for _, e := range entries {
		fmt.Printf("  %-35s %-12s %s\n", e.Path, e.Version, e.Hash)
	}

	fmt.Println("\n2. Module hash verification")
	tmpFile := "go.sum.test"
	content := "github.com/example/pkg v1.0.0 h1:abcdef\n"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	h, _ := hashModule(tmpFile)
	fmt.Printf("  go.sum hash: %s...\n", h[:16])

	fmt.Println("\n3. Dependency minification (unused deps removal)")
	deps := []string{
		"github.com/gorilla/mux",
		"github.com/lib/pq",
		"unused-old-package",
		"github.com/sirupsen/logrus",
		"unused-debug-tool",
	}
	fmt.Println("  Before:", deps)
	minified := minifyDeps(deps)
	fmt.Println("  After: ", minified)

	fmt.Println("\n4. go.sum parsing")
	fmt.Println("  Total modules tracked:", len(entries))

	fmt.Println("\n=== Software Bill of Materials ===")
	fmt.Println("Module: example-app v1.0.0")
	fmt.Println("Dependencies:")
	for _, e := range entries {
		fmt.Printf("  - %s@%s\n", e.Path, e.Version)
	}

	fmt.Println("\nNote: Run 'go mod verify' to check go.sum matches cached modules.")
}
