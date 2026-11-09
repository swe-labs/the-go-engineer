package main

import (
	"strings"
	"testing"
)

func TestSumJSON(t *testing.T) {
	data := `1
2
3
4
5
`
	sum, err := SumJSON(strings.NewReader(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != 15 {
		t.Errorf("expected 15, got %d", sum)
	}
}

func TestSumJSONEmpty(t *testing.T) {
	sum, err := SumJSON(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != 0 {
		t.Errorf("expected 0, got %d", sum)
	}
}

func TestSumJSONInvalid(t *testing.T) {
	data := `42
abc
`
	_, err := SumJSON(strings.NewReader(data))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSumJSONSingle(t *testing.T) {
	sum, err := SumJSON(strings.NewReader("99\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != 99 {
		t.Errorf("expected 99, got %d", sum)
	}
}

func TestCompiles(t *testing.T) {
}
