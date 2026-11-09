package main

import (
	"testing"
)

func TestNewAddress_Valid(t *testing.T) {
	addr, err := NewAddress("123 Main St", "Springfield", "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.City != "Springfield" {
		t.Errorf("expected Springfield, got %s", addr.City)
	}
}

func TestNewAddress_Invalid(t *testing.T) {
	_, err := NewAddress("", "City", "12345")
	if err == nil {
		t.Fatal("expected error for empty street")
	}
}

func TestNewOrderItem_Valid(t *testing.T) {
	item, err := NewOrderItem("p1", 2, 9.99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Total() != 19.98 {
		t.Errorf("expected 19.98, got %.2f", item.Total())
	}
}

func TestNewOrderItem_InvalidQuantity(t *testing.T) {
	_, err := NewOrderItem("p1", 0, 9.99)
	if err == nil {
		t.Fatal("expected error for zero quantity")
	}
}

func TestNewOrderItem_InvalidPrice(t *testing.T) {
	_, err := NewOrderItem("p1", 1, 0)
	if err == nil {
		t.Fatal("expected error for zero price")
	}
}

func TestNewOrder_Valid(t *testing.T) {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	item, _ := NewOrderItem("p1", 1, 10.00)
	order, err := NewOrder("ORD-001", "Alice", []OrderItem{item}, addr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.Status != OrderPending {
		t.Errorf("expected pending, got %s", order.Status)
	}
	if order.Total() != 10.00 {
		t.Errorf("expected 10.00, got %.2f", order.Total())
	}
	if len(order.Events) != 1 || order.Events[0].EventType != "order.created" {
		t.Errorf("expected order.created event")
	}
}

func TestNewOrder_MissingCustomer(t *testing.T) {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	item, _ := NewOrderItem("p1", 1, 10.00)
	_, err := NewOrder("ORD-001", "", []OrderItem{item}, addr)
	if err == nil {
		t.Fatal("expected error for empty customer")
	}
}

func TestNewOrder_NoItems(t *testing.T) {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	_, err := NewOrder("ORD-001", "Alice", []OrderItem{}, addr)
	if err == nil {
		t.Fatal("expected error for empty items")
	}
}

func TestOrder_Ship(t *testing.T) {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	item, _ := NewOrderItem("p1", 1, 10.00)
	order, _ := NewOrder("ORD-001", "Alice", []OrderItem{item}, addr)

	err := order.Ship()
	if err != nil {
		t.Fatalf("Ship failed: %v", err)
	}
	if order.Status != OrderShipped {
		t.Errorf("expected shipped, got %s", order.Status)
	}
	if len(order.Events) != 2 || order.Events[1].EventType != "order.shipped" {
		t.Errorf("expected order.shipped event")
	}
}

func TestOrder_Ship_Twice(t *testing.T) {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	item, _ := NewOrderItem("p1", 1, 10.00)
	order, _ := NewOrder("ORD-001", "Alice", []OrderItem{item}, addr)
	order.Ship()
	err := order.Ship()
	if err == nil {
		t.Fatal("expected error when shipping already shipped order")
	}
}
