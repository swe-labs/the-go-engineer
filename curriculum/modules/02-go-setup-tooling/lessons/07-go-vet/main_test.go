package main

import (
	"testing"
)

func TestDiagnoseCommonVetIssues(t *testing.T) {
	diagnostics := DiagnoseCommonVetIssues()
	if len(diagnostics) < 4 {
		t.Fatalf("expected at least 4 diagnostic patterns, got %d", len(diagnostics))
	}

	foundPrintf := false
	for _, d := range diagnostics {
		if d.CheckName == "printf" {
			foundPrintf = true
			if d.Message == "" || d.Example == "" {
				t.Errorf("diagnostic fields must be populated")
			}
		}
	}

	if !foundPrintf {
		t.Errorf("expected 'printf' checker diagnostic")
	}
}

func TestCheckFormatArgs(t *testing.T) {
	tests := []struct {
		format   string
		argCount int
		wantErr  bool
	}{
		{"Hello, %s!", 1, false},
		{"Values: %d, %d, %s", 3, false},
		{"Percent: 100%% complete", 0, false},
		{"Percent: 100%% complete for %s", 1, false},
		{"Missing arg: %s %d", 1, true},
		{"Extra arg: %s", 2, true},
	}

	for _, tt := range tests {
		err := CheckFormatArgs(tt.format, tt.argCount)
		if (err != nil) != tt.wantErr {
			t.Errorf("CheckFormatArgs(%q, %d) error = %v, wantErr %v", tt.format, tt.argCount, err, tt.wantErr)
		}
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
