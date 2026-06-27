package main

import (
	"testing"
)

func TestGetPID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"returns positive PID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPID(); got <= 0 {
				t.Errorf("GetPID() = %d, want > 0", got)
			}
		})
	}
}

func TestGetPPID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"returns positive PPID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPPID(); got <= 0 {
				t.Errorf("GetPPID() = %d, want > 0", got)
			}
		})
	}
}

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		args    []string
		want    int
		wantErr bool
	}{
		{"valid command returns 0", "go", []string{"version"}, 0, false},
		{"invalid command returns error", "nonexistent-command-xyz", []string{}, -1, true},
		{"failing command returns non-zero", "go", []string{"build", "-invalid-flag"}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunCommand(tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunCommand() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPIDandPPIDAreDifferent(t *testing.T) {
	pid := GetPID()
	ppid := GetPPID()
	if pid == ppid {
		t.Errorf("PID (%d) and PPID (%d) should be different", pid, ppid)
	}
}
