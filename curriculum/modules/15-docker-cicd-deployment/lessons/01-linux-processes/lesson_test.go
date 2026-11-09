package main

import (
	"testing"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		args     []string
		wantCode int
		wantErr  bool
	}{
		{
			name:     "go version succeeds",
			cmd:      "go",
			args:     []string{"version"},
			wantCode: 0,
			wantErr:  false,
		},
		{
			name:     "non-existent command fails",
			cmd:      "nonexistent-command-12345",
			args:     nil,
			wantCode: -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := runCommand(tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCommand() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if code != tt.wantCode {
				t.Errorf("runCommand() code = %d, want %d", code, tt.wantCode)
			}
		})
	}
}
