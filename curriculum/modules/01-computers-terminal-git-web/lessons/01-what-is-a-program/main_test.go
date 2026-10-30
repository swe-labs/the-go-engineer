package main

import (
	"strings"
	"testing"
)

func TestProgramStage_String(t *testing.T) {
	tests := []struct {
		stage    ProgramStage
		expected string
	}{
		{StageSource, "Source Code"},
		{StageBinary, "Executable Binary"},
		{StageProcess, "Operating System Process"},
		{ProgramStage(99), "Unknown Stage"},
	}

	for _, tt := range tests {
		got := tt.stage.String()
		if !strings.Contains(got, tt.expected) {
			t.Errorf("stage %v.String() = %q; expected to contain %q", tt.stage, got, tt.expected)
		}
	}
}

func TestMainExecution(t *testing.T) {
	// Verify main executes cleanly without panic
	main()
}
