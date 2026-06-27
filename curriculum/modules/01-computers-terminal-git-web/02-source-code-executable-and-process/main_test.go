package main

import (
	"strings"
	"testing"
)

func TestProgramLifecycleDescribe(t *testing.T) {
	p := ProgramLifecycle{
		SourceFile: "/home/user/main.go",
		BinaryPath: "/home/user/myapp",
		PID:        1234,
		Args:       []string{"./myapp", "--flag"},
	}

	got := p.Describe()
	wantParts := []string{"Source:", "/home/user/main.go", "Binary:", "/home/user/myapp", "PID:", "1234", "Args:", "./myapp", "--flag"}
	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Errorf("Describe() missing %q\nFull output:\n%s", part, got)
		}
	}
}

func TestProgramLifecycleEmptyArgs(t *testing.T) {
	p := ProgramLifecycle{
		SourceFile: "main.go",
		BinaryPath: "./app",
		PID:        0,
		Args:       []string{},
	}
	got := p.Describe()
	if !strings.Contains(got, "Args:") {
		t.Errorf("expected Args in Describe(), got %q", got)
	}
}

func TestProgramLifecyclePIDZero(t *testing.T) {
	p := ProgramLifecycle{
		SourceFile: "src/main.go",
		BinaryPath: "/bin/app",
		PID:        0,
		Args:       []string{"app"},
	}
	got := p.Describe()
	if !strings.Contains(got, "PID:    0") {
		t.Errorf("expected PID 0, got %q", got)
	}
}
