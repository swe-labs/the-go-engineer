package main

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetCallerInfo(t *testing.T) {
	fn, file, line, ok := GetCallerInfo(1)
	if !ok {
		t.Fatal("expected caller info")
	}
	if !strings.Contains(fn, "TestGetCallerInfo") {
		t.Errorf("expected TestGetCallerInfo, got %s", fn)
	}
	if !strings.HasSuffix(file, "_test.go") {
		t.Errorf("expected *_test.go, got %s", file)
	}
	if line <= 0 {
		t.Errorf("expected positive line number, got %d", line)
	}
}

func TestCallerSkipZero(t *testing.T) {
	fn, _, _, ok := GetCallerInfo(0)
	if !ok {
		t.Fatal("expected caller info")
	}
	if !strings.Contains(fn, "GetCallerInfo") {
		t.Errorf("expected GetCallerInfo, got %s", fn)
	}
}

func TestCallerFile(t *testing.T) {
	_, file, _, ok := GetCallerInfo(0)
	if !ok {
		t.Fatal("expected caller info")
	}
	if !strings.HasSuffix(file, "main.go") {
		t.Errorf("expected main.go, got %s", file)
	}
}

func TestStackTraceBuffer(t *testing.T) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])
	if !strings.Contains(stack, "TestStackTraceBuffer") {
		t.Errorf("expected stack to contain TestStackTraceBuffer, got %s", stack)
	}
}

func TestStackDepth(t *testing.T) {
	depth := StackDepth()
	if depth < 3 {
		t.Errorf("expected depth at least 3 (TestStackDepth + runtime + main), got %d", depth)
	}
}

func TestPanicRecover(t *testing.T) {
	panicked := true
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "boom" {
					t.Errorf("expected 'boom', got %v", r)
				}
			}
		}()
		deep()
		panicked = false
	}()
	if panicked {
		t.Log("panic was recovered as expected")
	}
}
