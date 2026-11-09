package main

import (
	"fmt"
	"testing"
)

// FakeUserStore is an in-memory fake for UserStore.
type FakeUserStore struct {
	users map[string]int
}

func NewFakeUserStore() *FakeUserStore {
	return &FakeUserStore{users: make(map[string]int)}
}

func (f *FakeUserStore) Save(name string, age int) error {
	f.users[name] = age
	return nil
}

func (f *FakeUserStore) Find(name string) (int, error) {
	age, ok := f.users[name]
	if !ok {
		return 0, fmt.Errorf("user %q not found", name)
	}
	return age, nil
}

func TestGreeter(t *testing.T) {
	store := NewFakeUserStore()
	store.Save("Alice", 30)

	greeter := NewGreeter(store)
	msg, err := greeter.Greet("Alice")
	if err != nil {
		t.Fatalf("Greet failed: %v", err)
	}
	want := "Hello, Alice! You are 30 years old."
	if msg != want {
		t.Errorf("Greet = %q; want %q", msg, want)
	}
}

func TestGreeterNotFound(t *testing.T) {
	store := NewFakeUserStore()
	greeter := NewGreeter(store)
	_, err := greeter.Greet("Bob")
	if err == nil {
		t.Fatal("expected error for unknown user")
	}
}

func TestFakeStoreRoundTrip(t *testing.T) {
	store := NewFakeUserStore()
	if err := store.Save("Bob", 25); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	age, err := store.Find("Bob")
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}
	if age != 25 {
		t.Errorf("age = %d; want 25", age)
	}
}

// FakeCalculator stores the last operation.
type FakeCalculator struct {
	LastOp  string
	LastA   int
	LastB   int
	AddRes  int
	MultRes int
}

func (f *FakeCalculator) Add(a, b int) int {
	f.LastOp = "Add"
	f.LastA = a
	f.LastB = b
	return f.AddRes
}

func (f *FakeCalculator) Multiply(a, b int) int {
	f.LastOp = "Multiply"
	f.LastA = a
	f.LastB = b
	return f.MultRes
}

func TestDoubleWithFakeCalculator(t *testing.T) {
	fc := &FakeCalculator{AddRes: 10}
	got := Double(fc, 5)
	if got != 10 {
		t.Errorf("Double(5) = %d; want 10", got)
	}
	if fc.LastOp != "Add" {
		t.Errorf("LastOp = %s; want Add", fc.LastOp)
	}
}

func TestDoubleWithRealCalculator(t *testing.T) {
	got := Double(RealCalculator{}, 5)
	if got != 10 {
		t.Errorf("Double(5) = %d; want 10", got)
	}
}
