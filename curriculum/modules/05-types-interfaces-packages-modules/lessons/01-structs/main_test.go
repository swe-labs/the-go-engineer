package main

import (
	"encoding/json"
	"testing"
)

func TestProductZeroValue(t *testing.T) {
	var p Product
	if p.ID != 0 || p.Name != "" || p.Price != 0.0 {
		t.Errorf("expected zero-value Product, got %+v", p)
	}
}

func TestProductLiteral(t *testing.T) {
	p := Product{ID: 42, Name: "Test", Price: 1.99}
	if p.ID != 42 || p.Name != "Test" || p.Price != 1.99 {
		t.Errorf("unexpected Product: %+v", p)
	}
}

func TestProductComparison(t *testing.T) {
	a := Product{ID: 1, Name: "A"}
	b := Product{ID: 1, Name: "A"}
	if a != b {
		t.Error("identical structs should be equal")
	}
}

func TestProductJSON(t *testing.T) {
	p := Product{ID: 5, Name: "Widget", Price: 0}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Product
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ID != p.ID || decoded.Name != p.Name {
		t.Errorf("round-trip mismatch: got %+v", decoded)
	}
}

func TestBookJSONOmitsPages(t *testing.T) {
	b := Book{Title: "T", Author: "A", Pages: 100}
	data, _ := json.Marshal(b)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	if _, ok := result["pages"]; ok {
		t.Error("expected Pages field to be omitted from JSON")
	}
}

func TestAnonymousStruct(t *testing.T) {
	pt := struct{ X, Y int }{X: 3, Y: 4}
	if pt.X != 3 || pt.Y != 4 {
		t.Errorf("unexpected anonymous struct: %+v", pt)
	}
}

func TestCopySemantics(t *testing.T) {
	original := Product{ID: 1, Name: "Original", Price: 10.0}
	copied := original
	copied.Price = 99.0
	if original.Price != 10.0 {
		t.Error("struct assignment should copy, not share")
	}
}
