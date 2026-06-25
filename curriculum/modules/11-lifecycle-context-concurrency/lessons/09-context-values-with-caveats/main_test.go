package main

import (
	"context"
	"testing"
)

func TestContextValues(t *testing.T) {
	key := contextKey("color")
	ctx := context.WithValue(context.Background(), key, "blue")

	val, ok := ctx.Value(key).(string)
	if !ok {
		t.Fatal("Value returned non-string")
	}
	if val != "blue" {
		t.Errorf("Value = %q, want %q", val, "blue")
	}
}

func TestContextValueMissing(t *testing.T) {
	key := contextKey("missing")
	ctx := context.Background()

	if ctx.Value(key) != nil {
		t.Error("expected nil for missing key")
	}
}

func TestContextValueDifferentKeys(t *testing.T) {
	key1 := contextKey("a")
	key2 := contextKey("b")
	key3 := contextKey("c")

	ctx := context.WithValue(context.Background(), key1, "x")
	ctx = context.WithValue(ctx, key2, "y")

	if v := ctx.Value(key1).(string); v != "x" {
		t.Errorf("key1 = %q, want %q", v, "x")
	}
	if v := ctx.Value(key2).(string); v != "y" {
		t.Errorf("key2 = %q, want %q", v, "y")
	}
	if ctx.Value(key3) != nil {
		t.Error("key3 should be nil")
	}
}

func TestContextValueTypeSafety(t *testing.T) {
	type stringKey string
	type intKey int

	sk := stringKey("key")
	ik := intKey(42)

	ctx := context.WithValue(context.Background(), sk, "string val")
	ctx = context.WithValue(ctx, ik, 100)

	if v := ctx.Value(sk).(string); v != "string val" {
		t.Errorf("string key = %q", v)
	}
	if v := ctx.Value(ik).(int); v != 100 {
		t.Errorf("int key = %d", v)
	}
}

func TestContextValueImmutability(t *testing.T) {
	key := contextKey("items")
	original := []string{"a", "b"}
	ctx := context.WithValue(context.Background(), key, original)

	// Modify the original slice outside context
	original[0] = "mutated"

	retrieved := ctx.Value(key).([]string)
	if retrieved[0] != "mutated" {
		t.Error("context value was not mutated - this test demonstrates the caveat")
	}
}
