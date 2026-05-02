// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started
// Title: Development Environment
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the small command loop that makes day-to-day Go work predictable.
//
// WHY THIS MATTERS:
//   - The Go toolchain is a workflow, not a single command: 1. Edit code. 2. Format it.
//     3. Build or run it. 4. Test it when tests exist.
//   - Following this loop ensures your code is clean and valid before you share it.
//
// RUN:
//   go run ./01-getting-started/4-dev-environment
//
// KEY TAKEAWAY:
//   - The "Standard Loop" (fmt -> build -> test) is the heartbeat of a Go engineer's day.
// ============================================================================

package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// commandInfo (Struct): groups the state used by the command info example boundary.
type commandInfo struct {
	name        string
	description string
}

// toolInfo (Struct): groups the state used by the tool info example boundary.
type toolInfo struct {
	name string
}

func main() {
	commands := []commandInfo{
		{name: "go fmt ./...", description: "format Go code into the standard style"},
		{name: "go build ./...", description: "compile packages to verify they are valid"},
		{name: "go run [path]", description: "compile and execute one target"},
		{name: "go test ./...", description: "run tests across the repo"},
	}

	tools := []toolInfo{
		{name: "gofmt"},
		{name: "gopls"},
	}

	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
	fmt.Println("Core Go commands:")
	for _, command := range commands {
		fmt.Printf("- %-16s %s\n", command.name, command.description)
	}

	fmt.Println()
	fmt.Println("Tool check:")
	for _, tool := range tools {
		path, err := exec.LookPath(tool.name)
		if err != nil {
			fmt.Printf("- %s: not found on PATH\n", tool.name)
			continue
		}
		fmt.Printf("- %s: %s\n", tool.name, path)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.5 -> 01-getting-started/5-go-tools")
	fmt.Println("Run    : go run ./01-getting-started/5-go-tools")
	fmt.Println("Current: GT.4 (dev-environment)")
	fmt.Println("---------------------------------------------------")
}
