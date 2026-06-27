package main

import "testing"

func TestExitCodeByFileStatusExists(t *testing.T) {
	if got := ExitCodeByFileStatus(true); got != 0 {
		t.Errorf("ExitCodeByFileStatus(true) = %d, want 0", got)
	}
}

func TestExitCodeByFileStatusNotExists(t *testing.T) {
	if got := ExitCodeByFileStatus(false); got != 1 {
		t.Errorf("ExitCodeByFileStatus(false) = %d, want 1", got)
	}
}
