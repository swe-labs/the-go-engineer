package main

import (
	"strings"
	"testing"
)

func TestSource(t *testing.T) {
	got := stageString(StageSource)
	want := "source code: human-readable text written by the programmer"
	if got != want {
		t.Errorf("stageString(StageSource) = %q; want %q", got, want)
	}
}

func TestBinary(t *testing.T) {
	got := stageString(StageBinary)
	want := "executable binary: machine instructions produced by the compiler"
	if got != want {
		t.Errorf("stageString(StageBinary) = %q; want %q", got, want)
	}
}

func TestProcess(t *testing.T) {
	got := stageString(StageProcess)
	want := "running process: binary loaded into memory by the OS"
	if got != want {
		t.Errorf("stageString(StageProcess) = %q; want %q", got, want)
	}
}

func TestUnknown(t *testing.T) {
	got := stageString(99)
	if !strings.Contains(got, "unknown") {
		t.Errorf("expected unknown for invalid stage, got %q", got)
	}
}

func TestStageOrder(t *testing.T) {
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
