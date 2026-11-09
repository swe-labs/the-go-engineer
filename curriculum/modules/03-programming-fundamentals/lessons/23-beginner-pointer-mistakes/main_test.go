package main

import "testing"

func TestMakeIncrementer(t *testing.T) {
	inc := makeIncrementer()
	for i := 1; i <= 5; i++ {
		got := inc()
		if got != i {
			t.Errorf("inc() = %d; want %d (iteration %d)", got, i, i)
		}
	}
}

func TestMakeIncrementerIsolated(t *testing.T) {
	inc1 := makeIncrementer()
	inc2 := makeIncrementer()
	if inc1() != 1 || inc1() != 2 {
		t.Error("inc1 should be independent")
	}
	if inc2() != 1 || inc2() != 2 {
		t.Error("inc2 should be independent")
	}
}

func TestCaptureCorrect(t *testing.T) {
	items := []int{10, 20, 30}
	ptrs := captureCorrect(items)
	for i, p := range ptrs {
		if *p != items[i] {
			t.Errorf("ptrs[%d] = %d; want %d", i, *p, items[i])
		}
	}
	if ptrs[0] == ptrs[1] {
		t.Error("all pointers should be distinct addresses")
	}
}

func TestCaptureBuggy(t *testing.T) {
	items := []int{10, 20, 30}
	ptrs := captureBuggy(items)
	// In the buggy version, all pointers point to the same address.
	if len(ptrs) < 2 {
		t.Fatal("need at least 2 items")
	}
	// The repeated value should be the last item (30) if the bug manifests.
	for i := 1; i < len(ptrs); i++ {
		if *ptrs[i] != items[len(items)-1] {
			t.Log("buggy capture may not reproduce on Go 1.22+")
			return
		}
	}
}
