package main

import (
	"testing"
)

func TestStackAllocated(t *testing.T) {
	u := stackAllocated(1, "Alice")
	if u.ID != 1 || u.Name != "Alice" {
		t.Errorf("stackAllocated = %+v", u)
	}
}

func TestHeapAllocated(t *testing.T) {
	u := heapAllocated(2, "Bob")
	if u.ID != 2 || u.Name != "Bob" {
		t.Errorf("heapAllocated = %+v", u)
	}
}

func TestEscapeAnalysisComparison(t *testing.T) {
	// Both should produce identical values
	sa := stackAllocated(10, "same")
	ha := heapAllocated(10, "same")

	if sa != *ha {
		t.Errorf("values should match: %+v vs %+v", sa, *ha)
	}
}

func TestNoEscapeOnSmallStruct(t *testing.T) {
	// A small struct created and used locally should not escape
	type small struct {
		a, b, c byte
	}
	s := small{1, 2, 3}
	if s.a+s.b+s.c != 6 {
		t.Errorf("unexpected sum")
	}
}

func TestEscapeThroughClosure(t *testing.T) {
	fn := func() *int {
		x := 42
		return &x
	}
	if got := *fn(); got != 42 {
		t.Errorf("closure escape = %d; want 42", got)
	}
}

func TestNoEscapeThroughClosure(t *testing.T) {
	// Value receiver closure does not cause escape of the struct
	type point struct{ x, y int }
	p := point{3, 4}
	dist := func() int {
		return p.x*p.x + p.y*p.y
	}
	if dist() != 25 {
		t.Errorf("unexpected distance")
	}
}
