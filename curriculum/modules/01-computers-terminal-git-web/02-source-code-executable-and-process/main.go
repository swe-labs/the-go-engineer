package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func describeProgram(sourceFile, binaryPath string, pid int, args []string) string {
	return fmt.Sprintf(
		"Source: %s\nBinary: %s\nPID:    %d\nArgs:   %v",
		sourceFile, binaryPath, pid, args,
	)
}

func main() {
	exe, _ := os.Executable()
	wd, _ := os.Getwd()

	sourceFile := filepath.Join(wd, "main.go")
	binaryPath := exe
	pid := os.Getpid()
	args := os.Args

	fmt.Println("=== Source Code, Executable, and Process ===")
	fmt.Println()
	fmt.Println(describeProgram(sourceFile, binaryPath, pid, args))
	fmt.Println()

	path, _ := exec.LookPath(os.Args[0])
	fmt.Printf("Binary resolved via PATH: %s\n", path)
	fmt.Println()
	fmt.Println("The source file is read by the compiler to produce the binary.")
	fmt.Println("The OS loader creates a process from the binary at runtime.")
}
