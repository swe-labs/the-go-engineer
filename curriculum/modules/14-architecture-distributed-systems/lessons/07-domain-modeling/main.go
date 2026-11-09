package main

import (
	"errors"
	"fmt"
	"time"
)

type Address struct {
	Street  string
	City    string
	ZipCode string
}

func NewAddress(street, city, zipCode string) (Address, error) {
	if street == "" || city == "" || zipCode == "" {
		return Address{}, errors.New("all address fields are required")
	}
	return Address{Street: street, City: city, ZipCode: zipCode}, nil
}

type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice float64
}

func NewOrderItem(productID string, quantity int, unitPrice float64) (OrderItem, error) {
	if productID == "" {
		return OrderItem{}, errors.New("product ID is required")
	}
	if quantity <= 0 {
		return OrderItem{}, errors.New("quantity must be positive")
	}
	if unitPrice <= 0 {
		return OrderItem{}, errors.New("unit price must be positive")
	}
	return OrderItem{ProductID: productID, Quantity: quantity, UnitPrice: unitPrice}, nil
}

func (i OrderItem) Total() float64 {
	return float64(i.Quantity) * i.UnitPrice
}

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderShipped   OrderStatus = "shipped"
	OrderDelivered OrderStatus = "delivered"
)

type OrderEvent struct {
	OrderID   string
	EventType string
	Timestamp time.Time
	Data      interface{}
}

type Order struct {
	ID        string
	Customer  string
	Items     []OrderItem
	Shipping  Address
	Status    OrderStatus
	Events    []OrderEvent
	CreatedAt time.Time
}

func NewOrder(id, customer string, items []OrderItem, shipping Address) (*Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}
	if customer == "" {
		return nil, errors.New("customer is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	now := time.Now()
	order := &Order{
		ID:        id,
		Customer:  customer,
		Items:     items,
		Shipping:  shipping,
		Status:    OrderPending,
		CreatedAt: now,
	}
	order.raiseEvent("order.created", nil)
	return order, nil
}

func (o *Order) Total() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Total()
	}
	return total
}

func (o *Order) Ship() error {
	if o.Status != OrderPending {
		return errors.New("only pending orders can be shipped")
	}
	o.Status = OrderShipped
	o.raiseEvent("order.shipped", nil)
	return nil
}

func (o *Order) raiseEvent(eventType string, data interface{}) {
	o.Events = append(o.Events, OrderEvent{
		OrderID:   o.ID,
		EventType: eventType,
		Timestamp: time.Now(),
		Data:      data,
	})
}

func main() {
	addr, _ := NewAddress("123 Main St", "Springfield", "12345")
	item, _ := NewOrderItem("prod-1", 2, 19.99)
	order, _ := NewOrder("ORD-001", "Alice", []OrderItem{item}, addr)

	fmt.Printf("Order %s: $%.2f (status: %s)\n", order.ID, order.Total(), order.Status)
	_ = order.Ship()
	fmt.Printf("Order shipped: %s\n", order.Status)
	fmt.Printf("Events: %d raised\n", len(order.Events))
}
