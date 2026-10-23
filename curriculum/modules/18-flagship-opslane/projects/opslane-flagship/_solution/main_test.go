package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestServer(t *testing.T) (*OpslaneServer, *sql.DB) {
	t.Helper()
	dbName := fmt.Sprintf("file:test_opslane_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	cfg := &Config{
		Port:           8080,
		Environment:    "test",
		WorkerCount:    2,
		QueueCapacity:  16,
		RequestTimeout: 2 * time.Second,
	}

	srv, err := NewOpslaneServer(cfg, db, nil)
	if err != nil {
		t.Fatalf("failed to create OpslaneServer: %v", err)
	}

	return srv, db
}

func TestHealthAndReadinessEndpoints(t *testing.T) {
	srv, db := setupTestServer(t)
	defer db.Close()
	defer srv.Close()

	handler := srv.Routes()

	// Test /healthz
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 from /healthz, got %d", rec.Code)
	}

	// Test /readyz
	reqReady := httptest.NewRequest("GET", "/readyz", nil)
	recReady := httptest.NewRecorder()
	handler.ServeHTTP(recReady, reqReady)

	if recReady.Code != http.StatusOK {
		t.Errorf("expected 200 from /readyz, got %d", recReady.Code)
	}
}

func TestMultiTenantIsolation(t *testing.T) {
	srv, db := setupTestServer(t)
	defer db.Close()
	defer srv.Close()

	ctx := context.Background()

	// 1. Tenant Acme creates an order
	orderAcme, err := srv.CreateOrder(ctx, "tenant-acme", "Acme Retailer", 4500)
	if err != nil {
		t.Fatalf("failed to create order for tenant-acme: %v", err)
	}

	// 2. Tenant Globex creates an order
	orderGlobex, err := srv.CreateOrder(ctx, "tenant-globex", "Globex Corp", 9000)
	if err != nil {
		t.Fatalf("failed to create order for tenant-globex: %v", err)
	}

	// 3. Globex attempts to read Acme's order -> must fail with ErrNotFound (isolation boundary)
	_, err = srv.GetOrder(ctx, "tenant-globex", orderAcme.ID)
	if err == nil {
		t.Fatalf("CRITICAL SECURITY VIOLATION: tenant-globex accessed tenant-acme order %s", orderAcme.ID)
	}

	// 4. Listing orders for Acme returns only Acme orders
	acmeOrders, err := srv.ListOrders(ctx, "tenant-acme")
	if err != nil {
		t.Fatalf("failed to list orders for acme: %v", err)
	}
	if len(acmeOrders) != 1 || acmeOrders[0].ID != orderAcme.ID {
		t.Errorf("expected 1 order for acme, got %d", len(acmeOrders))
	}

	// 5. Listing orders for Globex returns only Globex orders
	globexOrders, err := srv.ListOrders(ctx, "tenant-globex")
	if err != nil {
		t.Fatalf("failed to list orders for globex: %v", err)
	}
	if len(globexOrders) != 1 || globexOrders[0].ID != orderGlobex.ID {
		t.Errorf("expected 1 order for globex, got %d", len(globexOrders))
	}
}

func TestIdempotentPaymentWorkflow(t *testing.T) {
	srv, db := setupTestServer(t)
	defer db.Close()
	defer srv.Close()

	ctx := context.Background()
	tenantID := "tenant-alpha"

	// Create pending order
	order, err := srv.CreateOrder(ctx, tenantID, "Customer Alpha", 7500)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	idemKey := "idem-key-payment-uuid-123"

	// 1. Initial payment
	p1, err := srv.ProcessPayment(ctx, tenantID, order.ID, idemKey)
	if err != nil {
		t.Fatalf("initial payment failed: %v", err)
	}
	if p1.Status != "SUCCESS" || p1.AmountCents != 7500 {
		t.Errorf("unexpected payment response: %+v", p1)
	}

	// Verify order status is COMPLETED
	updatedOrder, err := srv.GetOrder(ctx, tenantID, order.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated order: %v", err)
	}
	if updatedOrder.Status != StatusCompleted {
		t.Errorf("expected order status COMPLETED, got %s", updatedOrder.Status)
	}

	// 2. Duplicate payment submission with SAME idempotency key
	p2, err := srv.ProcessPayment(ctx, tenantID, order.ID, idemKey)
	if err != nil {
		t.Fatalf("duplicate payment call failed: %v", err)
	}
	if p2.ID != p1.ID {
		t.Errorf("expected idempotent response with payment ID %s, got %s", p1.ID, p2.ID)
	}

	// 3. Duplicate payment submission with DIFFERENT key on completed order must fail
	_, err = srv.ProcessPayment(ctx, tenantID, order.ID, "new-different-key")
	if err == nil {
		t.Fatalf("expected error attempting to repay completed order with new idempotency key")
	}
}

func TestCrossTenantPaymentPrevention(t *testing.T) {
	srv, db := setupTestServer(t)
	defer db.Close()
	defer srv.Close()

	ctx := context.Background()

	// Tenant A creates an order
	orderA, err := srv.CreateOrder(ctx, "tenant-a", "Customer A", 1000)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Tenant B tries to pay for Tenant A's order -> must be denied
	_, err = srv.ProcessPayment(ctx, "tenant-b", orderA.ID, "pay-attempt-1")
	if err == nil {
		t.Fatalf("CRITICAL SECURITY VIOLATION: tenant-b processed payment for tenant-a order")
	}
}

func TestHTTPAPI_RBACAndFullLifecycle(t *testing.T) {
	srv, db := setupTestServer(t)
	defer db.Close()
	defer srv.Close()

	handler := srv.Routes()

	// 1. Viewer attempts to create an order -> 403 Forbidden
	payload := []byte(`{"customer":"Big Corp","amount_cents":5000}`)
	reqViewer := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewReader(payload))
	reqViewer.Header.Set("X-Tenant-ID", "tenant-corp")
	reqViewer.Header.Set("X-Role", "viewer")
	recViewer := httptest.NewRecorder()
	handler.ServeHTTP(recViewer, reqViewer)

	if recViewer.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for viewer role, got %d", recViewer.Code)
	}

	// 2. Operator creates order -> 201 Created
	reqOp := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewReader(payload))
	reqOp.Header.Set("X-Tenant-ID", "tenant-corp")
	reqOp.Header.Set("X-Role", "operator")
	recOp := httptest.NewRecorder()
	handler.ServeHTTP(recOp, reqOp)

	if recOp.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for operator role, got %d", recOp.Code)
	}

	var createdOrder Order
	if err := json.Unmarshal(recOp.Body.Bytes(), &createdOrder); err != nil {
		t.Fatalf("failed to decode created order: %v", err)
	}

	// 3. Operator pays order -> 200 OK
	payPayload := []byte(`{"idempotency_key":"idem-http-001"}`)
	reqPay := httptest.NewRequest("POST", "/api/v1/orders/"+createdOrder.ID+"/pay", bytes.NewReader(payPayload))
	reqPay.Header.Set("X-Tenant-ID", "tenant-corp")
	reqPay.Header.Set("X-Role", "operator")
	recPay := httptest.NewRecorder()
	handler.ServeHTTP(recPay, reqPay)

	if recPay.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on payment, got %d (body: %s)", recPay.Code, recPay.Body.String())
	}

	// 4. Metrics endpoint reflects operations
	reqMetrics := httptest.NewRequest("GET", "/metrics", nil)
	recMetrics := httptest.NewRecorder()
	handler.ServeHTTP(recMetrics, reqMetrics)

	body, _ := io.ReadAll(recMetrics.Body)
	if !bytes.Contains(body, []byte(`"orders_created":1`)) || !bytes.Contains(body, []byte(`"payments_success":1`)) {
		t.Errorf("metrics response does not reflect operations: %s", string(body))
	}
}
