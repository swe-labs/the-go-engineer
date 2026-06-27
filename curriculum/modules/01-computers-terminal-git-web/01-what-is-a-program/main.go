package main

import (
	"fmt"
	"os"
)

const (
	StageSource int = iota
	StageBinary
	StageProcess
)

func stageString(s int) string {
	switch s {
	case StageSource:
		return "source code: human-readable text written by the programmer"
	case StageBinary:
		return "executable binary: machine instructions produced by the compiler"
	case StageProcess:
		return "running process: binary loaded into memory by the OS"
	default:
		return "unknown"
	}
}

func main() {
	fmt.Println("=== What is a program? ===")
	fmt.Println()

	for stage := StageSource; stage <= StageProcess; stage++ {
		fmt.Printf("Stage %d: %s\n", stage+1, stageString(stage))
	}

	fmt.Println()
	fmt.Printf("This running program's source file:  %s\n", os.Args[0])
	fmt.Println()
	fmt.Println("A program is source code compiled into a binary that the OS loads as a process.")
}
