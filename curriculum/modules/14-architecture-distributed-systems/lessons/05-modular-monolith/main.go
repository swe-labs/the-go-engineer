package main

import (
	"errors"
	"fmt"
	"time"
)

type Order struct {
	ID        string
	Total     float64
	Status    string
	CreatedAt time.Time
}

type Invoice struct {
	OrderID string
	Amount  float64
	Paid    bool
}

type OrderStore interface {
	Save(order Order) error
	FindByID(id string) (Order, error)
}

type InvoiceStore interface {
	Save(invoice Invoice) error
	FindByOrderID(orderID string) (Invoice, error)
}

type BillingService struct {
	invoices InvoiceStore
}

func NewBillingService(invoices InvoiceStore) *BillingService {
	return &BillingService{invoices: invoices}
}

func (s *BillingService) CreateInvoice(order Order) (Invoice, error) {
	inv := Invoice{OrderID: order.ID, Amount: order.Total, Paid: false}
	if err := s.invoices.Save(inv); err != nil {
		return Invoice{}, err
	}
	return inv, nil
}

type OrderingBoundedContext struct {
	orders  OrderStore
	billing *BillingService
}

func NewOrderingBoundedContext(orders OrderStore, billing *BillingService) *OrderingBoundedContext {
	return &OrderingBoundedContext{orders: orders, billing: billing}
}

func (ctx *OrderingBoundedContext) PlaceOrder(id string, total float64) (Order, error) {
	if total <= 0 {
		return Order{}, errors.New("total must be positive")
	}
	order := Order{ID: id, Total: total, Status: "placed", CreatedAt: time.Now()}
	if err := ctx.orders.Save(order); err != nil {
		return Order{}, err
	}
	inv, err := ctx.billing.CreateInvoice(order)
	if err != nil {
		return Order{}, fmt.Errorf("billing failed: %w", err)
	}
	fmt.Printf("Invoice %s created for order %s ($%.2f)\n", inv.OrderID, order.ID, inv.Amount)
	return order, nil
}

type inMemoryOrderStore struct {
	data map[string]Order
}

func newInMemoryOrderStore() *inMemoryOrderStore {
	return &inMemoryOrderStore{data: make(map[string]Order)}
}

func (s *inMemoryOrderStore) Save(order Order) error {
	s.data[order.ID] = order
	return nil
}

func (s *inMemoryOrderStore) FindByID(id string) (Order, error) {
	o, ok := s.data[id]
	if !ok {
		return Order{}, errors.New("order not found")
	}
	return o, nil
}

type inMemoryInvoiceStore struct {
	data map[string]Invoice
}

func newInMemoryInvoiceStore() *inMemoryInvoiceStore {
	return &inMemoryInvoiceStore{data: make(map[string]Invoice)}
}

func (s *inMemoryInvoiceStore) Save(invoice Invoice) error {
	s.data[invoice.OrderID] = invoice
	return nil
}

func (s *inMemoryInvoiceStore) FindByOrderID(orderID string) (Invoice, error) {
	inv, ok := s.data[orderID]
	if !ok {
		return Invoice{}, errors.New("invoice not found")
	}
	return inv, nil
}

func main() {
	orders := newInMemoryOrderStore()
	invoices := newInMemoryInvoiceStore()
	billing := NewBillingService(invoices)
	ordering := NewOrderingBoundedContext(orders, billing)

	order, err := ordering.PlaceOrder("ORD-001", 129.99)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Order placed: %s (status: %s)\n", order.ID, order.Status)
}
