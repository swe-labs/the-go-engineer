package main

import "testing"

func TestSafeSet(t *testing.T) {
	t.Run("nil map creates new map", func(t *testing.T) {
		var m map[string]int
		m = safeSet(m, "key", 42)
		if m == nil {
			t.Fatal("safeSet should not return nil map")
		}
		if m["key"] != 42 {
			t.Errorf("safeSet nil map result: m[key] = %d; want 42", m["key"])
		}
	})

	t.Run("existing map is mutated", func(t *testing.T) {
		m := make(map[string]int)
		m = safeSet(m, "a", 1)
		m = safeSet(m, "b", 2)
		if m["a"] != 1 || m["b"] != 2 {
			t.Errorf("expected a=1, b=2; got a=%d, b=%d", m["a"], m["b"])
		}
	})

	t.Run("multiple inserts", func(t *testing.T) {
		var m map[string]int
		for i := 0; i < 10; i++ {
			key := string(rune('a' + i))
			m = safeSet(m, key, i)
		}
		if len(m) != 10 {
			t.Errorf("expected 10 entries; got %d", len(m))
		}
	})
}

func TestSafeGet(t *testing.T) {
	t.Run("nil map returns zero, false", func(t *testing.T) {
		var m map[string]int
		val, ok := safeGet(m, "anything")
		if val != 0 {
			t.Errorf("expected 0; got %d", val)
		}
		if ok {
			t.Errorf("expected false; got true")
		}
	})

	t.Run("existing key returns value", func(t *testing.T) {
		m := map[string]int{"existing": 99}
		val, ok := safeGet(m, "existing")
		if val != 99 {
			t.Errorf("expected 99; got %d", val)
		}
		if !ok {
			t.Errorf("expected true; got false")
		}
	})

	t.Run("missing key returns zero, false", func(t *testing.T) {
		m := map[string]int{"a": 1}
		val, ok := safeGet(m, "missing")
		if val != 0 {
			t.Errorf("expected 0; got %d", val)
		}
		if ok {
			t.Errorf("expected false; got true")
		}
	})
}
