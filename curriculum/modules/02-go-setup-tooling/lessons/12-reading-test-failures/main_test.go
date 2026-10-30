package main

import (
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"standard addition", 2, 2, 4},
		{"identity with zero", 0, 0, 0},
		{"negative integers", -5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("%s", FormatFailureMessage(tt.name, got, tt.want))
			}
		})
	}
}

func TestFormatFailureMessage(t *testing.T) {
	msg := FormatFailureMessage("Sqrt(4)", 3, 2)
	expected := "Sqrt(4) = 3; want 2"
	if msg != expected {
		t.Errorf("FormatFailureMessage() = %q, want %q", msg, expected)
	}
	if !strings.Contains(msg, "want 2") {
		t.Errorf("message should clearly specify desired want value")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
