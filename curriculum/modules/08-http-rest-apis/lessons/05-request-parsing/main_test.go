package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProcessHandler_Success(t *testing.T) {
	body := `{"name":"Alice","count":5}`
	req := httptest.NewRequest(http.MethodPost, "/process/u1?format=json", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Trace-ID", "trace-123")
	req.SetPathValue("id", "u1")
	rec := httptest.NewRecorder()
	processHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["id"] != "u1" {
		t.Errorf("expected id u1, got %v", resp["id"])
	}
	if resp["trace_id"] != "trace-123" {
		t.Errorf("expected trace_id trace-123, got %v", resp["trace_id"])
	}
}

func TestProcessHandler_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/process/", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	processHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestProcessHandler_MissingFields(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/process/u1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "u1")
	rec := httptest.NewRecorder()
	processHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if _, ok := errResp.Errors["name"]; !ok {
		t.Error("expected name error")
	}
	if _, ok := errResp.Errors["count"]; !ok {
		t.Error("expected count error")
	}
}

func TestProcessHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/process/u1", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "u1")
	rec := httptest.NewRecorder()
	processHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestProcessHandler_DefaultFormat(t *testing.T) {
	body := `{"name":"Bob","count":3}`
	req := httptest.NewRequest(http.MethodPost, "/process/u2", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "u2")
	rec := httptest.NewRecorder()
	processHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProcessHandler_UnsupportedFormat(t *testing.T) {
	body := `{"name":"Bob","count":3}`
	req := httptest.NewRequest(http.MethodPost, "/process/u2?format=xml", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "u2")
	rec := httptest.NewRecorder()
	processHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
