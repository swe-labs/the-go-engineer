package main

import (
	"testing"
)

func TestOrderingBoundedContext_PlaceOrder(t *testing.T) {
	orders := newInMemoryOrderStore()
	invoices := newInMemoryInvoiceStore()
	billing := NewBillingService(invoices)
	ctx := NewOrderingBoundedContext(orders, billing)

	order, err := ctx.PlaceOrder("ORD-001", 129.99)
	if err != nil {
		t.Fatalf("PlaceOrder failed: %v", err)
	}
	if order.ID != "ORD-001" {
		t.Errorf("expected ORD-001, got %s", order.ID)
	}
	if order.Status != "placed" {
		t.Errorf("expected placed, got %s", order.Status)
	}
}

func TestOrderingBoundedContext_PlaceOrder_ZeroTotal(t *testing.T) {
	orders := newInMemoryOrderStore()
	invoices := newInMemoryInvoiceStore()
	billing := NewBillingService(invoices)
	ctx := NewOrderingBoundedContext(orders, billing)

	_, err := ctx.PlaceOrder("ORD-002", 0)
	if err == nil {
		t.Fatal("expected error for zero total")
	}
}

func TestBillingService_CreateInvoice(t *testing.T) {
	invoices := newInMemoryInvoiceStore()
	billing := NewBillingService(invoices)

	inv, err := billing.CreateInvoice(Order{ID: "ORD-001", Total: 99.99})
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}
	if inv.Amount != 99.99 {
		t.Errorf("expected 99.99, got %.2f", inv.Amount)
	}
	if inv.Paid {
		t.Error("expected invoice to be unpaid")
	}
}

func TestBillingService_CreateInvoice_StoreInvoice(t *testing.T) {
	invoices := newInMemoryInvoiceStore()
	billing := NewBillingService(invoices)
	billing.CreateInvoice(Order{ID: "ORD-001", Total: 49.99})

	found, err := invoices.FindByOrderID("ORD-001")
	if err != nil {
		t.Fatalf("FindByOrderID failed: %v", err)
	}
	if found.Amount != 49.99 {
		t.Errorf("expected 49.99, got %.2f", found.Amount)
	}
}
