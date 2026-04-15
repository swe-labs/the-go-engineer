// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

type commandInfo struct {
	name        string
	description string
}

type toolInfo struct {
	name string
}

func main() {
	commands := []commandInfo{
		{name: "go fmt ./...", description: "format Go code into the standard style"},
		{name: "go build ./...", description: "compile packages to verify they are valid"},
		{name: "go run ./01-foundations/01-getting-started/2-hello-world", description: "compile and execute one target"},
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
	fmt.Println("NEXT UP: LB.1 variables")
	fmt.Println("Current: GT.4 (dev-environment)")
	fmt.Println("---------------------------------------------------")
}
