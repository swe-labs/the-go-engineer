package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderService_SuccessfulCheckoutAndDiscount(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	payment := NewMockPaymentGateway()
	notifier := NewMockNotificationService()
	service := NewOrderService(repo, payment, notifier)

	ctx := context.Background()

	// Order >= $100: 2 items @ $60 each = $120 (12,000 cents). Discount: 10% = $12 (1,200 cents). Total = $108 (10,800 cents)
	dto := CreateOrderDTO{
		Customer: "Acme Corp",
		Email:    "orders@acme.com",
		Items: []OrderItem{
			{SKU: "SERVER-RACK-01", Quantity: 2, PriceCents: 6000},
		},
		PaymentToken: "tok_valid_test",
	}

	order, err := service.Checkout(ctx, dto)
	if err != nil {
		t.Fatalf("unexpected checkout error: %v", err)
	}

	if order.Status != StatusPaid {
		t.Errorf("expected status PAID, got %s", order.Status)
	}
	if order.Subtotal != 12000 {
		t.Errorf("expected subtotal 12000, got %d", order.Subtotal)
	}
	if order.Discount != 1200 {
		t.Errorf("expected discount 1200, got %d", order.Discount)
	}
	if order.Total != 10800 {
		t.Errorf("expected total 10800, got %d", order.Total)
	}
	if order.PaymentRef == "" {
		t.Errorf("expected non-empty payment reference")
	}

	// Verify Payment gateway recorded exact charged amount
	if len(payment.RecordedCalls) != 1 || payment.RecordedCalls[0] != 10800 {
		t.Errorf("expected payment charge of 10800 cents, got: %+v", payment.RecordedCalls)
	}

	// Verify notification sent
	if len(notifier.SentRecipients) != 1 || notifier.SentRecipients[0] != "orders@acme.com" {
		t.Errorf("expected notification to orders@acme.com, got: %+v", notifier.SentRecipients)
	}

	// Verify saved in repository
	saved, err := repo.FindByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("order not found in repository: %v", err)
	}
	if saved.Total != 10800 {
		t.Errorf("saved order total mismatch: %d", saved.Total)
	}
}

func TestOrderService_PaymentFailureRollback(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	payment := NewMockPaymentGateway()
	notifier := NewMockNotificationService()
	service := NewOrderService(repo, payment, notifier)

	ctx := context.Background()

	dto := CreateOrderDTO{
		Customer:     "Failing Corp",
		Email:        "fail@corp.com",
		Items:        []OrderItem{{SKU: "ITEM-1", Quantity: 1, PriceCents: 5000}},
		PaymentToken: "fail_token",
	}

	_, err := service.Checkout(ctx, dto)
	if err == nil {
		t.Fatalf("expected error on failing payment token, got nil")
	}

	// Order must NOT be saved in repository
	orders, _ := repo.List(ctx)
	if len(orders) != 0 {
		t.Errorf("expected 0 saved orders after failed payment, got %d", len(orders))
	}

	// Notification must NOT be sent
	if len(notifier.SentRecipients) != 0 {
		t.Errorf("expected 0 notifications sent, got %d", len(notifier.SentRecipients))
	}
}

func TestHTTPHandler_EndToEnd(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	payment := NewMockPaymentGateway()
	notifier := NewMockNotificationService()
	service := NewOrderService(repo, payment, notifier)
	handler := NewHTTPHandler(service)
	router := handler.Routes()

	// 1. POST /orders
	payload := `{
		"customer": "Bob",
		"email": "bob@example.com",
		"items": [{"sku": "GADGET-1", "quantity": 1, "price_cents": 2500}],
		"payment_token": "tok_ok"
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var created Order
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == "" || created.Total != 2500 || created.Status != StatusPaid {
		t.Fatalf("unexpected created order: %+v", created)
	}

	// 2. GET /orders/{id}
	reqGet := httptest.NewRequest(http.MethodGet, "/orders/"+created.ID, nil)
	recGet := httptest.NewRecorder()
	router.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", recGet.Code)
	}

	// 3. GET /orders
	reqList := httptest.NewRequest(http.MethodGet, "/orders", nil)
	recList := httptest.NewRecorder()
	router.ServeHTTP(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from list, got %d", recList.Code)
	}
	var list []Order
	if err := json.Unmarshal(recList.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Fatalf("unexpected order list: %+v", list)
	}

	// 4. POST /orders with failing token -> 402 Payment Required
	failPayload := `{
		"customer": "Eve",
		"email": "eve@example.com",
		"items": [{"sku": "GADGET-2", "quantity": 1, "price_cents": 1000}],
		"payment_token": "fail_token"
	}`
	reqFail := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(failPayload))
	recFail := httptest.NewRecorder()
	router.ServeHTTP(recFail, reqFail)

	if recFail.Code != http.StatusPaymentRequired {
		t.Fatalf("expected 402 Payment Required, got %d", recFail.Code)
	}
}
