package main

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	doubled := Map(nums, func(n int) int { return n * 2 })
	want := []int{2, 4, 6, 8}

	if !reflect.DeepEqual(doubled, want) {
		t.Errorf("Map() = %v, want %v", doubled, want)
	}

	// Test cross-type mapping (int to string)
	strs := Map(nums, func(n int) string { return string(rune('A' + n - 1)) })
	wantStrs := []string{"A", "B", "C", "D"}
	if !reflect.DeepEqual(strs, wantStrs) {
		t.Errorf("Map() cross-type = %v, want %v", strs, wantStrs)
	}
}

func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	odds := Filter(nums, func(n int) bool { return n%2 != 0 })
	want := []int{1, 3, 5}

	if !reflect.DeepEqual(odds, want) {
		t.Errorf("Filter() = %v, want %v", odds, want)
	}
}

func TestContains(t *testing.T) {
	items := []int{10, 20, 30}
	if !Contains(items, 20) {
		t.Errorf("expected Contains to find 20")
	}
	if Contains(items, 99) {
		t.Errorf("expected Contains not to find 99")
	}
}

func TestMin(t *testing.T) {
	if got := Min(10, 20); got != 10 {
		t.Errorf("Min(10, 20) = %d, want 10", got)
	}
	if got := Min("zebra", "apple"); got != "apple" {
		t.Errorf("Min('zebra', 'apple') = %s, want 'apple'", got)
	}
}

func TestStack(t *testing.T) {
	stack := &Stack[int]{}

	if _, ok := stack.Pop(); ok {
		t.Errorf("expected Pop on empty stack to return false")
	}

	stack.Push(100)
	stack.Push(200)

	if stack.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", stack.Len())
	}

	val, ok := stack.Pop()
	if !ok || val != 200 {
		t.Errorf("expected Pop 200, got %d (ok=%t)", val, ok)
	}

	val, ok = stack.Pop()
	if !ok || val != 100 {
		t.Errorf("expected Pop 100, got %d (ok=%t)", val, ok)
	}

	if stack.Len() != 0 {
		t.Errorf("expected empty stack, got len %d", stack.Len())
	}
}
