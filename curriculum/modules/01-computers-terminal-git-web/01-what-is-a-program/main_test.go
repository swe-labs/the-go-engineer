package main

import (
	"strings"
	"testing"
)

func TestProgramStageString(t *testing.T) {
	tests := []struct {
		stage ProgramStage
		want  string
	}{
		{StageSource, "source code: human-readable text written by the programmer"},
		{StageBinary, "executable binary: machine instructions produced by the compiler"},
		{StageProcess, "running process: binary loaded into memory by the OS"},
	}

	for _, tt := range tests {
		got := tt.stage.String()
		if got != tt.want {
			t.Errorf("ProgramStage(%d).String() = %q; want %q", tt.stage, got, tt.want)
		}
	}
}

func TestProgramStageOrder(t *testing.T) {
	if StageSource != 0 {
		t.Errorf("StageSource should be 0, got %d", StageSource)
	}
	if StageBinary != 1 {
		t.Errorf("StageBinary should be 1, got %d", StageBinary)
	}
	if StageProcess != 2 {
		t.Errorf("StageProcess should be 2, got %d", StageProcess)
	}
}

func TestProgramStageUnknown(t *testing.T) {
	var bad ProgramStage = 99
	got := bad.String()
	if !strings.Contains(got, "unknown") {
		t.Errorf("expected unknown for invalid stage, got %q", got)
	}
}
