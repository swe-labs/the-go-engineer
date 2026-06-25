package main

import "testing"

func TestAdd(t *testing.T) {
	if got := add(2, 3); got != 5 {
		t.Errorf("add(2, 3) = %d; want 5", got)
	}
	if got := add(-1, 1); got != 0 {
		t.Errorf("add(-1, 1) = %d; want 0", got)
	}
}

func TestSubtract(t *testing.T) {
	if got := subtract(10, 4); got != 6 {
		t.Errorf("subtract(10, 4) = %d; want 6", got)
	}
}

func TestOperate(t *testing.T) {
	tests := []struct {
		op       string
		a, b, ex int
	}{
		{"add", 6, 2, 8},
		{"sub", 6, 2, 4},
		{"mul", 6, 2, 12},
		{"div", 6, 2, 3},
	}
	for _, tc := range tests {
		got := operate(tc.a, tc.b, tc.op)
		if got != tc.ex {
			t.Errorf("operate(%d, %d, %q) = %d; want %d", tc.a, tc.b, tc.op, got, tc.ex)
		}
	}
}

func TestOperatePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unknown operator")
		}
	}()
	operate(1, 1, "unknown")
}
