package main

import (
	"strings"
	"testing"
)

func TestInstructionStep_Properties(t *testing.T) {
	step := InstructionStep{
		Address:     "0x0045A120",
		Opcodes:     "48 c7 c0 2a 00 00 00",
		Assembly:    "MOVQ $0x2a, %rax",
		Description: "Load immediate integer 42 into register RAX",
	}

	if step.Address != "0x0045A120" {
		t.Errorf("unexpected address: %s", step.Address)
	}
	if !strings.Contains(step.Assembly, "MOVQ") {
		t.Errorf("expected MOVQ in assembly, got: %s", step.Assembly)
	}
}

func TestMainExecution(t *testing.T) {
	// Verify main executes cleanly without panic
	main()
}
