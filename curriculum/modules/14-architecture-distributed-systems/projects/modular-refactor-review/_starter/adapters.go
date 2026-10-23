package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// InMemoryOrderRepository implements OrderRepository in memory with thread safety.
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]Order),
	}
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id string) (Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return Order{}, ErrOrderNotFound
	}
	return order, nil
}

func (r *InMemoryOrderRepository) List(ctx context.Context) ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]Order, 0, len(r.orders))
	for _, o := range r.orders {
		res = append(res, o)
	}
	return res, nil
}

// MockPaymentGateway simulates an external payment gateway.
type MockPaymentGateway struct {
	mu            sync.Mutex
	RecordedCalls []int
}

func NewMockPaymentGateway() *MockPaymentGateway {
	return &MockPaymentGateway{
		RecordedCalls: make([]int, 0),
	}
}

func (p *MockPaymentGateway) Charge(ctx context.Context, amountCents int, paymentToken string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.RecordedCalls = append(p.RecordedCalls, amountCents)

	if paymentToken == "fail_token" || paymentToken == "" {
		return "", errors.New("card declined or invalid token")
	}

	return fmt.Sprintf("ch_%s", generateID()[:8]), nil
}

// MockNotificationService simulates sending customer notifications.
type MockNotificationService struct {
	mu             sync.Mutex
	SentRecipients []string
}

func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{
		SentRecipients: make([]string, 0),
	}
}

func (n *MockNotificationService) SendOrderConfirmation(ctx context.Context, order Order) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.SentRecipients = append(n.SentRecipients, order.Email)
	return nil
}
