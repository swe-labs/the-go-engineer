package main

import "testing"

func TestGetPID(t *testing.T) {
	if got := GetPID(); got <= 0 {
		t.Errorf("GetPID() = %d, want > 0", got)
	}
}

func TestGetPPID(t *testing.T) {
	if got := GetPPID(); got <= 0 {
		t.Errorf("GetPPID() = %d, want > 0", got)
	}
}

func TestRunCommandValid(t *testing.T) {
	got, err := RunCommand("go", "version")
	if err != nil {
		t.Errorf("RunCommand() error = %v", err)
	}
	if got != 0 {
		t.Errorf("RunCommand() = %d, want 0", got)
	}
}

func TestRunCommandInvalid(t *testing.T) {
	got, err := RunCommand("nonexistent-command-xyz")
	if err == nil {
		t.Error("RunCommand() expected error")
	}
	if got != -1 {
		t.Errorf("RunCommand() = %d, want -1", got)
	}
}

func TestRunCommandFailsWithNonZero(t *testing.T) {
	got, err := RunCommand("go", "build", "-invalid-flag")
	if err != nil {
		t.Errorf("RunCommand() unexpected error = %v", err)
	}
	if got == 0 {
		t.Error("RunCommand() expected non-zero exit code")
	}
}

func TestPIDandPPIDAreDifferent(t *testing.T) {
	pid := GetPID()
	ppid := GetPPID()
	if pid == ppid {
		t.Errorf("PID (%d) and PPID (%d) should be different", pid, ppid)
	}
}
