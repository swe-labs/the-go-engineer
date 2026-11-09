package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteResourceHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/resources/42", nil)
	req.SetPathValue("id", "42")
	rec := httptest.NewRecorder()
	deleteResourceHandler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body, got %s", rec.Body.String())
	}
}

func TestDeleteResourceHandler_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/resources/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	deleteResourceHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
	var errResp APIError
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Code != "NOT_FOUND" {
		t.Errorf("expected error code NOT_FOUND, got %s", errResp.Code)
	}
}

func TestDeleteResourceHandler_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/resources/", nil)
	rec := httptest.NewRecorder()
	deleteResourceHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp APIError
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Code != "BAD_REQUEST" {
		t.Errorf("expected error code BAD_REQUEST, got %s", errResp.Code)
	}
}

func TestDeleteResourceHandler_Panics(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/resources/panic", nil)
	req.SetPathValue("id", "panic")
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic but none occurred")
		}
	}()

	deleteResourceHandler(rec, req)
}

func TestWriteError_ContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, http.StatusTeapot, APIError{
		Code:    "TEAPOT",
		Message: "I'm a teapot",
	})

	if rec.Code != http.StatusTeapot {
		t.Errorf("expected 418, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("expected application/json, got %s", ct)
	}
	var errResp APIError
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Code != "TEAPOT" {
		t.Errorf("expected TEAPOT, got %s", errResp.Code)
	}
}

func TestWriteError_WithDetails(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, http.StatusUnprocessableEntity, APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid fields",
		Details: []FieldError{
			{Field: "name", Message: "required"},
		},
	})

	var errResp APIError
	json.NewDecoder(rec.Body).Decode(&errResp)
	if len(errResp.Details) != 1 {
		t.Errorf("expected 1 detail, got %d", len(errResp.Details))
	}
	if errResp.Details[0].Field != "name" {
		t.Errorf("expected field 'name', got %s", errResp.Details[0].Field)
	}
}
