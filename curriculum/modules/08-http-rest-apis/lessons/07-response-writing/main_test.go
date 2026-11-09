package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	pingHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("expected application/json, got %s", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if body["message"] != "pong" {
		t.Errorf("expected message 'pong', got %s", body["message"])
	}
}

func TestReportHandler_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/report", nil)
	rec := httptest.NewRecorder()
	reportHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestReportHandler_JSONStructure(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/report", nil)
	rec := httptest.NewRecorder()
	reportHandler(rec, req)

	// The output is a JSON sequence (not a single JSON value).
	// We check that the body contains key structural elements.
	body := rec.Body.String()
	if !strings.Contains(body, `"header":"report"`) {
		t.Error("expected header in report output")
	}
	if !strings.Contains(body, `"type":"summary"`) {
		t.Error("expected type in report output")
	}
	if !strings.Contains(body, `"items":"start"`) {
		t.Error("expected items start marker")
	}
	if !strings.Contains(body, `"items":"end"`) {
		t.Error("expected items end marker")
	}
	if !strings.Contains(body, `"status":"complete"`) {
		t.Error("expected footer")
	}
	// Check that all 100 items are present.
	for i := 1; i <= 100; i++ {
		expected := fmt.Sprintf(`"id":%d`, i)
		if !strings.Contains(body, expected) {
			t.Errorf("missing item %d in report", i)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusCreated, map[string]int{"id": 42})

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("expected application/json, got %s", ct)
	}
	var body map[string]int
	json.NewDecoder(rec.Body).Decode(&body)
	if body["id"] != 42 {
		t.Errorf("expected id 42, got %d", body["id"])
	}
}
