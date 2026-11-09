package main

import (
	"strings"
	"testing"
)

func TestSafeExecute_NoPanic(t *testing.T) {
	recovered, stack := SafeExecute(func() {
		_ = 2 + 2
	})
	if recovered != nil {
		t.Errorf("expected nil for non-panicking code, got %v", recovered)
	}
	if stack != "" {
		t.Errorf("expected empty stack trace for clean execution")
	}
}

func TestSafeExecute_InterceptPanic(t *testing.T) {
	slice := []int{1, 2, 3}
	recovered, stack := SafeExecute(func() {
		_ = AccessElement(slice, 10)
	})

	if recovered == nil {
		t.Fatalf("expected panic to be intercepted")
	}

	errMsg := recovered.(error).Error()
	if !strings.Contains(errMsg, "out of range") {
		t.Errorf("expected 'out of range' error message, got %q", errMsg)
	}
	if !strings.Contains(stack, "goroutine") {
		t.Errorf("expected stack trace to contain goroutine header, got: %s", stack)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
