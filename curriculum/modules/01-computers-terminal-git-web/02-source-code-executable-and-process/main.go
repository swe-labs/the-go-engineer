package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ProgramLifecycle struct {
	SourceFile string
	BinaryPath string
	PID        int
	Args       []string
}

func (p ProgramLifecycle) Describe() string {
	return fmt.Sprintf(
		"Source: %s\nBinary: %s\nPID:    %d\nArgs:   %v",
		p.SourceFile, p.BinaryPath, p.PID, p.Args,
	)
}

func main() {
	exe, _ := os.Executable()
	wd, _ := os.Getwd()

	lifecycle := ProgramLifecycle{
		SourceFile: filepath.Join(wd, "main.go"),
		BinaryPath: exe,
		PID:        os.Getpid(),
		Args:       os.Args,
	}

	fmt.Println("=== Source Code, Executable, and Process ===")
	fmt.Println()
	fmt.Println(lifecycle.Describe())
	fmt.Println()

	path, _ := exec.LookPath(os.Args[0])
	fmt.Printf("Binary resolved via PATH: %s\n", path)
	fmt.Println()
	fmt.Println("The source file is read by the compiler to produce the binary.")
	fmt.Println("The OS loader creates a process from the binary at runtime.")
}
