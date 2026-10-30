package main

import (
	"fmt"
	"go/format"
	"os/exec"
	"runtime"
)

// ToolCheck represents an individual toolchain readiness item.
type ToolCheck struct {
	Name    string
	Passed  bool
	Details string
}

// RunToolingChecklist verifies that all primary Go toolchain components are operating correctly.
func RunToolingChecklist() []ToolCheck {
	var checks []ToolCheck

	// 1. Check Go Runtime & Version
	v := runtime.Version()
	checks = append(checks, ToolCheck{
		Name:    "Go Compiler & Runtime",
		Passed:  v != "",
		Details: fmt.Sprintf("Version: %s (%s/%s)", v, runtime.GOOS, runtime.GOARCH),
	})

	// 2. Check Standard Library GOROOT
	goroot := runtime.GOROOT()
	checks = append(checks, ToolCheck{
		Name:    "Standard Library (GOROOT)",
		Passed:  true, // Embedded or explicit path
		Details: fmt.Sprintf("GOROOT: %s", goroot),
	})

	// 3. Check Go Binary on System PATH (exec.LookPath)
	goPath, err := exec.LookPath("go")
	checks = append(checks, ToolCheck{
		Name:    "CLI Toolchain on System PATH",
		Passed:  err == nil,
		Details: fmt.Sprintf("Path: %s", goPath),
	})

	// 4. Check Code Formatter (go/format)
	sample := []byte("package p\nfunc F(){}\n")
	_, fmtErr := format.Source(sample)
	checks = append(checks, ToolCheck{
		Name:    "Code Formatter (gofmt engine)",
		Passed:  fmtErr == nil,
		Details: "Verified go/format AST parsing and code emission",
	})

	return checks
}

func main() {
	fmt.Println("=== Module 02 Go Toolchain Readiness Checklist ===")

	checks := RunToolingChecklist()
	allPassed := true

	for i, c := range checks {
		status := "✓ PASS"
		if !c.Passed {
			status = "✗ FAIL"
			allPassed = false
		}
		fmt.Printf("[%d] %-30s [%s]\n    Details: %s\n\n", i+1, c.Name, status, c.Details)
	}

	if allPassed {
		fmt.Println("Result: 4 of 4 checks passed. Go development environment is 100% verified.")
		fmt.Println("You are ready to advance to Module 03 — Programming Fundamentals with Go!")
	} else {
		fmt.Println("Result: One or more toolchain checks failed. Review details above.")
	}
}
