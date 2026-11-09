package main

import "testing"

func TestGreeting(t *testing.T) {
	got := greeting("World")
	want := "Hello, World!"
	if got != want {
		t.Errorf("greeting(%q) = %q, want %q", "World", got, want)
	}
}

func TestGreetingEmpty(t *testing.T) {
	got := greeting("")
	want := "Hello, !"
	if got != want {
		t.Errorf("greeting(%q) = %q, want %q", "", got, want)
	}
}

func TestCompute(t *testing.T) {
	tests := []struct {
		a, b int
	}{
		{10, 5},
		{0, 0},
		{3, 3},
		{-2, 5},
		{100, 1},
	}
	for _, tc := range tests {
		got := compute(tc.a, tc.b)
		want := (tc.a + tc.b) + (tc.a - tc.b) + (tc.a * tc.b)
		if got != want {
			t.Errorf("compute(%d,%d) = %d, want %d", tc.a, tc.b, got, want)
		}
	}
}
