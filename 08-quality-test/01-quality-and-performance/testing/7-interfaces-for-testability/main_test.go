package main

import "testing"

// fixedClock (Struct): groups the state used by the fixed clock example boundary.
type fixedClock struct{}

// fixedClock.Now (Method): applies the now operation to receiver state at a visible boundary.
func (fixedClock) Now() string { return "now" }

// clock (Interface): captures the behavior boundary the clock example depends on.
type clock interface{ Now() string }

// render (Function): runs the render step and keeps its inputs, outputs, or errors visible.
func render(c clock) string { return c.Now() }

func TestTE7Seam(t *testing.T) {
	if got := render(fixedClock{}); got != "now" {
		t.Fatalf("unexpected render output: %s", got)
	}
}
