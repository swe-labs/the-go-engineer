package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshalUser(t *testing.T) {
	result, err := MarshalUser("Alice", 30, "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, `"name": "Alice"`) {
		t.Errorf("expected name field, got %s", result)
	}
	if !strings.Contains(result, `"email": "alice@example.com"`) {
		t.Errorf("expected email field, got %s", result)
	}
}

func TestMarshalUserOmitempty(t *testing.T) {
	result, err := MarshalUser("Bob", 25, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, `"email"`) {
		t.Errorf("expected email to be omitted, got %s", result)
	}
}

func TestMarshalUserNegativeAge(t *testing.T) {
	_, err := MarshalUser("Carol", -1, "c@t.com")
	if err == nil {
		t.Fatal("expected error for negative age")
	}
}

func TestMarshalUserValidJSON(t *testing.T) {
	result, err := MarshalUser("Dave", 42, "dave@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var u User
	if err := json.Unmarshal([]byte(result), &u); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if u.Name != "Dave" || u.Age != 42 || u.Email != "dave@example.com" {
		t.Errorf("round-trip mismatch: %+v", u)
	}
}

func TestCompiles(t *testing.T) {
}
