package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_CreateOrderSuccess(t *testing.T) {
	server := NewAPIServer()
	router := server.Router()

	payload := `{"title":"Enterprise License Subscription","amount":499.00,"currency":"USD"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Data OrderDTO `json:"data"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode API envelope: %v", err)
	}

	if resp.Data.ID == "" || resp.Data.Title != "Enterprise License Subscription" || resp.Data.Amount != 499.00 {
		t.Errorf("unexpected created order payload: %+v", resp.Data)
	}
}

func TestAPI_CreateOrderValidationFailure(t *testing.T) {
	// Core practice task requirement:
	// Verify request payload validation, error mapping, and field details.
	server := NewAPIServer()
	router := server.Router()

	// Missing title, negative amount, bad currency
	payload := `{"title":"","amount":-10.0,"currency":"BITCOIN"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", rr.Code)
	}

	var errEnv APIErrorEnvelope
	if err := json.NewDecoder(rr.Body).Decode(&errEnv); err != nil {
		t.Fatalf("failed to decode error envelope: %v", err)
	}

	if errEnv.Error.Code != "VALIDATION_FAILED" {
		t.Errorf("expected error code 'VALIDATION_FAILED', got %q", errEnv.Error.Code)
	}

	if len(errEnv.Error.Details) != 3 {
		t.Errorf("expected 3 field validation errors, got %d: %+v", len(errEnv.Error.Details), errEnv.Error.Details)
	}
}

func TestAPI_GetOrder(t *testing.T) {
	server := NewAPIServer()
	router := server.Router()

	// Pre-seed an order
	server.orders["ord-test-42"] = &OrderDTO{
		ID:     "ord-test-42",
		Title:  "Seeded Order",
		Amount: 75.0,
	}

	// 1. Existing order
	reqFound := httptest.NewRequest(http.MethodGet, "/api/v1/orders/ord-test-42", nil)
	rrFound := httptest.NewRecorder()
	router.ServeHTTP(rrFound, reqFound)

	if rrFound.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rrFound.Code)
	}

	// 2. Non-existent order
	reqMissing := httptest.NewRequest(http.MethodGet, "/api/v1/orders/ord-missing-99", nil)
	rrMissing := httptest.NewRecorder()
	router.ServeHTTP(rrMissing, reqMissing)

	if rrMissing.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", rrMissing.Code)
	}

	var errEnv APIErrorEnvelope
	_ = json.NewDecoder(rrMissing.Body).Decode(&errEnv)
	if errEnv.Error.Code != "ORDER_NOT_FOUND" {
		t.Errorf("expected code 'ORDER_NOT_FOUND', got %q", errEnv.Error.Code)
	}
}

func TestAPI_DeleteOrder(t *testing.T) {
	server := NewAPIServer()
	router := server.Router()

	server.orders["ord-del-1"] = &OrderDTO{ID: "ord-del-1", Title: "To Delete"}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/ord-del-1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on delete, got %d", rr.Code)
	}

	// Delete again -> 404
	req2 := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/ord-del-1", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on subsequent delete, got %d", rr2.Code)
	}
}

func TestAPI_RequestIDPropagation(t *testing.T) {
	server := NewAPIServer()
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	req.Header.Set("X-Request-ID", "custom-req-id-12345")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if headerVal := rr.Header().Get("X-Request-ID"); headerVal != "custom-req-id-12345" {
		t.Errorf("expected X-Request-ID header to be propagated, got %q", headerVal)
	}
}

func TestAPI_PayloadTooLarge(t *testing.T) {
	server := NewAPIServer()
	router := server.Router()

	// 1.5MB payload (exceeds 1MB limit)
	hugeString := strings.Repeat("A", 1500000)
	payload := `{"title":"` + hugeString + `","amount":10.0}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 Payload Too Large, got %d", rr.Code)
	}
}
