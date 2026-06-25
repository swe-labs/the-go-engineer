package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateOrder_Success(t *testing.T) {
	body := `{"item":"Widget"}`
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	createOrder(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var order Order
	json.NewDecoder(rec.Body).Decode(&order)
	if order.Item != "Widget" {
		t.Errorf("expected Widget, got %s", order.Item)
	}
	if order.Status != "created" {
		t.Errorf("expected status 'created', got %s", order.Status)
	}
}

func TestCreateOrder_MissingItem(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	createOrder(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rec.Code)
	}
}

func TestCreateOrder_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	createOrder(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetOrder_Found(t *testing.T) {
	// Seed an order.
	orders["ord-1"] = Order{ID: "ord-1", Item: "Gadget", Status: "created"}
	req := httptest.NewRequest(http.MethodGet, "/orders/ord-1", nil)
	req.SetPathValue("id", "ord-1")
	rec := httptest.NewRecorder()
	getOrder(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var order Order
	json.NewDecoder(rec.Body).Decode(&order)
	if order.ID != "ord-1" {
		t.Errorf("expected ord-1, got %s", order.ID)
	}
	delete(orders, "ord-1")
}

func TestGetOrder_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/ord-none", nil)
	req.SetPathValue("id", "ord-none")
	rec := httptest.NewRecorder()
	getOrder(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestCancelOrder_Success(t *testing.T) {
	orders["ord-5"] = Order{ID: "ord-5", Item: "Test", Status: "created"}
	req := httptest.NewRequest(http.MethodPut, "/orders/ord-5/cancel", nil)
	req.SetPathValue("id", "ord-5")
	rec := httptest.NewRecorder()
	cancelOrder(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var order Order
	json.NewDecoder(rec.Body).Decode(&order)
	if order.Status != "cancelled" {
		t.Errorf("expected status 'cancelled', got %s", order.Status)
	}
	delete(orders, "ord-5")
}

func TestCancelOrder_AlreadyCancelled(t *testing.T) {
	orders["ord-6"] = Order{ID: "ord-6", Item: "Test", Status: "cancelled"}
	req := httptest.NewRequest(http.MethodPut, "/orders/ord-6/cancel", nil)
	req.SetPathValue("id", "ord-6")
	rec := httptest.NewRecorder()
	cancelOrder(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rec.Code)
	}
	delete(orders, "ord-6")
}

func TestCancelOrder_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/orders/ord-missing/cancel", nil)
	req.SetPathValue("id", "ord-missing")
	rec := httptest.NewRecorder()
	cancelOrder(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
