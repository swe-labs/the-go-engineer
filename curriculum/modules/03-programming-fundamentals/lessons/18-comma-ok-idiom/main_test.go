package main

import "testing"

func TestLookupSafe(t *testing.T) {
	data := map[string]interface{}{
		"name":   "Alice",
		"age":    30,
		"score":  95.5,
		"active": true,
	}

	t.Run("key exists and type matches string", func(t *testing.T) {
		val, ok, err := LookupSafe(data, "name", "string")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if s, _ := val.(string); s != "Alice" {
			t.Fatalf("expected Alice, got %v", val)
		}
	})

	t.Run("key exists and type matches int", func(t *testing.T) {
		val, ok, err := LookupSafe(data, "age", "int")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if v, _ := val.(int); v != 30 {
			t.Fatalf("expected 30, got %v", val)
		}
	})

	t.Run("key exists but type mismatches", func(t *testing.T) {
		_, ok, err := LookupSafe(data, "age", "string")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for type mismatch")
		}
	})

	t.Run("key does not exist", func(t *testing.T) {
		val, ok, err := LookupSafe(data, "missing", "string")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected ok=false for missing key")
		}
		if val != nil {
			t.Fatalf("expected nil, got %v", val)
		}
	})

	t.Run("unknown type returns error", func(t *testing.T) {
		_, _, err := LookupSafe(data, "name", "bytes")
		if err == nil {
			t.Fatal("expected error for unknown type")
		}
	})

	t.Run("float64 type", func(t *testing.T) {
		val, ok, err := LookupSafe(data, "score", "float64")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if v, _ := val.(float64); v != 95.5 {
			t.Fatalf("expected 95.5, got %v", val)
		}
	})

	t.Run("bool type", func(t *testing.T) {
		val, ok, err := LookupSafe(data, "active", "bool")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected ok=true")
		}
		if v, _ := val.(bool); v != true {
			t.Fatalf("expected true, got %v", val)
		}
	})
}
