package main

import (
	"fmt"
	"os"
	"runtime"
)

// ProgramStage represents one of the three fundamental states of a program.
type ProgramStage int

const (
	StageSource ProgramStage = iota
	StageBinary
	StageProcess
)

func (s ProgramStage) String() string {
	switch s {
	case StageSource:
		return "Source Code (inert human-readable text on disk)"
	case StageBinary:
		return "Executable Binary (compiled machine code on disk)"
	case StageProcess:
		return "Operating System Process (active execution in memory)"
	default:
		return "Unknown Stage"
	}
}

func main() {
	fmt.Println("=== Three Stages of a Program ===")
	stages := []ProgramStage{StageSource, StageBinary, StageProcess}
	for i, stage := range stages {
		fmt.Printf("Stage %d: %s\n", i+1, stage.String())
	}

	fmt.Println("\n=== Active Process Runtime Metadata ===")
	fmt.Printf("Executable Path (os.Args[0]): %s\n", os.Args[0])
	fmt.Printf("Compiler:                     %s\n", runtime.Compiler)
	fmt.Printf("Target OS / Architecture:     %s / %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Available Logical CPUs:       %d\n", runtime.NumCPU())
	fmt.Printf("Process PID:                  %d\n", os.Getpid())

	fmt.Println("\nSummary: Source code compiles into a binary executable, which the OS kernel executes as a process.")
}
