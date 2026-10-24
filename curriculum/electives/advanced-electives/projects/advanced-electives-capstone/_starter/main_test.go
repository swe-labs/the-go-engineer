package main

import "testing"

func TestNewRingBuffer(t *testing.T) {
	rb := NewRingBuffer[int](5)
	if rb == nil {
		t.Fatal("expected non-nil buffer")
	}
	if !rb.Push(42) {
		t.Fatal("expected push to succeed")
	}
}
