package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// InstructionStep represents a CPU-level Fetch-Decode-Execute operation.
type InstructionStep struct {
	Address     string
	Opcodes     string
	Assembly    string
	Description string
}

func main() {
	fmt.Println("=== CPU Architecture & Execution Environment ===")
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("CPU Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Logical CPU Cores: %d\n", runtime.NumCPU())
	fmt.Printf("Current Process PID: %d\n\n", os.Getpid())

	// Simulated trace of compiled machine instructions for a simple addition:
	// Go: a := 42; b := 8; return a + b
	trace := []InstructionStep{
		{
			Address:     "0x0045A120",
			Opcodes:     "48 c7 c0 2a 00 00 00",
			Assembly:    "MOVQ $0x2a, %rax",
			Description: "Load immediate integer 42 into register RAX",
		},
		{
			Address:     "0x0045A127",
			Opcodes:     "48 c7 c3 08 00 00 00",
			Assembly:    "MOVQ $0x08, %rbx",
			Description: "Load immediate integer 8 into register RBX",
		},
		{
			Address:     "0x0045A12E",
			Opcodes:     "48 01 d8",
			Assembly:    "ADDQ %rbx, %rax",
			Description: "ALU adds RAX + RBX; result (50) stored in RAX",
		},
		{
			Address:     "0x0045A131",
			Opcodes:     "c3",
			Assembly:    "RETQ",
			Description: "Pop return address from stack into RIP; jump to caller",
		},
	}

	fmt.Println("=== Tracing CPU Fetch-Decode-Execute Cycle ===")
	for i, step := range trace {
		fmt.Printf("[Step %d] PC/RIP -> %s\n", i+1, step.Address)
		fmt.Printf("   Opcodes (Raw Bytes): %s\n", step.Opcodes)
		fmt.Printf("   Instruction:         %-20s\n", step.Assembly)
		fmt.Printf("   Action:              %s\n\n", step.Description)
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("Conclusion: The CPU executes byte sequences by sequentially bumping the RIP register.")
}
