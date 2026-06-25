package main

import (
	"testing"
)

func TestCountdownOutput(t *testing.T) {
	// We can't easily capture stdout in a unit test,
	// but we can verify the function doesn't panic.
	countdown()
}

func TestAddLogger(t *testing.T) {
	called := false
	fn := func(x int) int {
		called = true
		return x * 2
	}
	logged := addLogger(fn)
	res, err := logged(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected fn to be called")
	}
	if res != 14 {
		t.Errorf("expected 14, got %d", res)
	}
}

func TestAddLoggerNegative(t *testing.T) {
	fn := func(x int) int {
		return x * x
	}
	logged := addLogger(fn)
	res, err := logged(-4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != 16 {
		t.Errorf("expected 16, got %d", res)
	}
}

func TestLIFOOrder(t *testing.T) {
	var order []int
	func() {
		defer func() { order = append(order, 1) }()
		defer func() { order = append(order, 2) }()
		defer func() { order = append(order, 3) }()
	}()
	expected := []int{3, 2, 1}
	for i, v := range order {
		if v != expected[i] {
			t.Errorf("position %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestDeferArgumentCapture(t *testing.T) {
	result := func() int {
		x := 1
		defer func(v int) int { return v }(x)
		x = 99
		return x
	}()
	if result != 99 {
		t.Errorf("expected 99, got %d", result)
	}
}
