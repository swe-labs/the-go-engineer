package main

import (
	"strings"
	"testing"
)

func TestDescribeProgram(t *testing.T) {
	got := describeProgram("/home/user/main.go", "/home/user/myapp", 1234, []string{"./myapp", "--flag"})
	wantParts := []string{"Source:", "/home/user/main.go", "Binary:", "/home/user/myapp", "PID:", "1234", "Args:", "./myapp", "--flag"}
	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Errorf("describeProgram() missing %q\nFull output:\n%s", part, got)
		}
	}
}

func TestDescribeProgramEmptyArgs(t *testing.T) {
	got := describeProgram("main.go", "./app", 0, []string{})
	if !strings.Contains(got, "Args:") {
		t.Errorf("expected Args in describeProgram(), got %q", got)
	}
}

func TestDescribeProgramPIDZero(t *testing.T) {
	got := describeProgram("src/main.go", "/bin/app", 0, []string{"app"})
	if !strings.Contains(got, "PID:    0") {
		t.Errorf("expected PID 0, got %q", got)
	}
}
