package main

import (
	"testing"
)

func TestUnmarshalEventValid(t *testing.T) {
	data := []byte(`{"id": 1, "title": "deploy", "timestamp": "2025-06-01T10:00:00Z", "payload": {"env":"prod"}}`)
	e, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.ID != 1 || e.Title != "deploy" || e.Timestamp != "2025-06-01T10:00:00Z" {
		t.Errorf("unexpected fields: %+v", e)
	}
	if len(e.Payload) == 0 {
		t.Error("expected non-empty payload")
	}
}

func TestUnmarshalEventWrongType(t *testing.T) {
	data := []byte(`{"id": "not-a-number", "title": "x", "timestamp": "2025-01-01T00:00:00Z"}`)
	_, err := UnmarshalEvent(data)
	if err == nil {
		t.Error("expected error for wrong type")
	}
}

func TestUnmarshalEventNoTimestamp(t *testing.T) {
	data := []byte(`{"id": 2, "title": "test"}`)
	e, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.ID != 2 || e.Title != "test" {
		t.Errorf("unexpected fields: %+v", e)
	}
}

func TestUnmarshalEventExtraFields(t *testing.T) {
	data := []byte(`{"id": 3, "title": "w", "timestamp": "2025-01-01T00:00:00Z", "extra": "ignored"}`)
	_, err := UnmarshalEvent(data)
	if err != nil {
		t.Errorf("extra fields should be ignored, got error: %v", err)
	}
}

func TestCompiles(t *testing.T) {
}
