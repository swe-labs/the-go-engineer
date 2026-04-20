package main

import "testing"

type fixedClock struct{}

func (fixedClock) Now() string { return "now" }

type clock interface{ Now() string }

func render(c clock) string { return c.Now() }

func TestTE7Seam(t *testing.T) {
	if got := render(fixedClock{}); got != "now" {
		t.Fatalf("unexpected render output: %s", got)
	}
}
