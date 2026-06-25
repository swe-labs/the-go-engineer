package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLayersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/layers", nil)
	rec := httptest.NewRecorder()

	layersHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var result BuildSuggestion
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result.Layers) == 0 {
		t.Error("expected at least one layer")
	}

	if result.EstimatedSize <= 0 {
		t.Errorf("expected positive estimated size, got %d", result.EstimatedSize)
	}
}

func TestLayersHandlerCacheability(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/layers", nil)
	rec := httptest.NewRecorder()

	layersHandler(rec, req)

	var result BuildSuggestion
	json.NewDecoder(rec.Body).Decode(&result)

	for _, layer := range result.Layers {
		if layer.Command == "" {
			t.Error("found layer with empty command")
		}
	}
}

func TestLayersHasFirstLayer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/layers", nil)
	rec := httptest.NewRecorder()

	layersHandler(rec, req)

	var result BuildSuggestion
	json.NewDecoder(rec.Body).Decode(&result)

	if result.Layers[0].Command != "FROM golang:1.25-alpine" {
		t.Errorf("first layer command = %q, want %q", result.Layers[0].Command, "FROM golang:1.25-alpine")
	}
}
