package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteProducts(t *testing.T) {
	products := []Product{
		{"A100", "Widget", 9.99},
		{"A200", "Gadget", 24.99},
	}
	var buf bytes.Buffer
	err := WriteProducts(products, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for i, line := range lines {
		var p Product
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}

func TestWriteProductsEmpty(t *testing.T) {
	var buf bytes.Buffer
	err := WriteProducts(nil, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil slice, got %q", buf.String())
	}
}

func TestWriteProductsRoundTrip(t *testing.T) {
	original := []Product{{"Z999", "Test", 1.23}}
	var buf bytes.Buffer
	if err := WriteProducts(original, &buf); err != nil {
		t.Fatal(err)
	}
	var decoded Product
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.SKU != "Z999" || decoded.Price != 1.23 {
		t.Errorf("round-trip mismatch: %+v", decoded)
	}
}

func TestCompiles(t *testing.T) {
}
