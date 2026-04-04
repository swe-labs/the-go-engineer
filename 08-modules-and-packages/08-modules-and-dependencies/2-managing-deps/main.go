// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 8: Modules & Dependencies — Managing Dependencies
// Level: Intermediate
// ============================================================================
//
// RUN: go run ./08-modules-and-dependencies/2-managing-deps
// ============================================================================

import (
	"fmt"
	"os/exec"
	"strings"
)

// ============================================================================
// Section 8: Managing Dependencies
// Level: Intermediate
// ============================================================================
//
// This file demonstrates dependency management workflows.
// Run the commands below in a terminal to see them in action.
//
// ADDING DEPENDENCIES:
//   go get github.com/stretchr/testify@latest       — latest version
//   go get github.com/stretchr/testify@v1.11.1       — specific version
//   go get github.com/stretchr/testify@v1.10.0       — downgrade
//
// REMOVING DEPENDENCIES:
//   go get github.com/some/pkg@none                  — remove from go.mod
//   go mod tidy                                       — clean up unused
//
// INSPECTING DEPENDENCIES:
//   go list -m all                                    — all modules
//   go list -m -versions github.com/stretchr/testify  — available versions
//   go mod why github.com/stretchr/objx              — why is this needed?
//   go mod graph                                      — full dependency tree
//
// SECURITY:
//   go install golang.org/x/vuln/cmd/govulncheck@latest
//   govulncheck ./...                                — scan for known vulnerabilities
// ============================================================================

func main() {
	fmt.Println("=== Managing Dependencies ===")
	fmt.Println()

	// Run "go list -m all" to show current dependencies
	fmt.Println("Current module dependencies:")
	fmt.Println(strings.Repeat("-", 60))

	out, err := exec.Command("go", "list", "-m", "all").Output()
	if err != nil {
		fmt.Printf("Error running 'go list -m all': %v\n", err)
		fmt.Println("(This is expected if CGO is not available for sqlite3)")
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Classify as direct or indirect
		if strings.Contains(line, "the-go-engineer") {
			fmt.Printf("  [ROOT]     %s\n", line)
		} else {
			fmt.Printf("  [DEP]      %s\n", line)
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))

	// Show why a specific indirect dependency exists
	fmt.Println("\nWhy do we have github.com/stretchr/objx?")
	whyOut, err := exec.Command("go", "mod", "why", "github.com/stretchr/objx").Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(whyOut))
}
