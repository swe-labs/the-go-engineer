package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOkHandler(t *testing.T) {
	handler := withRecovery(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %s", body["status"])
	}
}

func TestPanicHandler_ReturnsJSON(t *testing.T) {
	handler := withRecovery(http.HandlerFunc(panicHandler))
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("expected JSON content type, got %s", ct)
	}
	var apiErr APIError
	if err := json.NewDecoder(rec.Body).Decode(&apiErr); err != nil {
		t.Fatalf("failed to decode error JSON: %v", err)
	}
	if apiErr.Code != "PANIC" {
		t.Errorf("expected code PANIC, got %s", apiErr.Code)
	}
	if apiErr.Message != "internal error" {
		t.Errorf("expected message 'internal error', got %s", apiErr.Message)
	}
}

func TestNotFoundHandler(t *testing.T) {
	handler := withRecovery(http.HandlerFunc(notFoundHandler))
	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
	var apiErr APIError
	json.NewDecoder(rec.Body).Decode(&apiErr)
	if apiErr.Code != "NOT_FOUND" {
		t.Errorf("expected code NOT_FOUND, got %s", apiErr.Code)
	}
}

func TestRecoveryWriter_WriteHeaderOnce(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &recoveryWriter{ResponseWriter: rec, status: http.StatusOK}
	rw.WriteHeader(http.StatusNotFound)
	rw.WriteHeader(http.StatusOK) // second call should be ignored

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
	if rw.status != http.StatusNotFound {
		t.Errorf("expected tracked status 404, got %d", rw.status)
	}
}

func TestPanicDoesNotCrash(t *testing.T) {
	// This test verifies that the middleware catches the panic and the test
	// process does not crash.
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("panic was not caught by middleware")
		}
	}()

	handler := withRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 after panic, got %d", rec.Code)
	}
}

func TestPlainErrorHandler(t *testing.T) {
	handler := withRecovery(http.HandlerFunc(plainErrorHandler))
	req := httptest.NewRequest(http.MethodGet, "/plain-error", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	// The plain error handler already wrote a text body.
	// The middleware should not overwrite it.
	body := rec.Body.String()
	if !strings.Contains(body, "something went wrong") {
		t.Errorf("expected original error message, got %s", body)
	}
}
