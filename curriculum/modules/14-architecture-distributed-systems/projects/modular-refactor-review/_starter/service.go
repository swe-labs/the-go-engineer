package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type CreateOrderDTO struct {
	Customer     string      `json:"customer"`
	Email        string      `json:"email"`
	Items        []OrderItem `json:"items"`
	PaymentToken string      `json:"payment_token"`
}

type OrderService struct {
	repo     OrderRepository
	payment  PaymentGateway
	notifier NotificationService
}

func NewOrderService(repo OrderRepository, payment PaymentGateway, notifier NotificationService) *OrderService {
	return &OrderService{
		repo:     repo,
		payment:  payment,
		notifier: notifier,
	}
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *OrderService) Checkout(ctx context.Context, dto CreateOrderDTO) (Order, error) {
	// TODO: Task 3 - Implement checkout orchestration:
	// 1. Construct Order with StatusPending
	// 2. CalculateTotals()
	// 3. Charge via payment gateway
	// 4. Update status to StatusPaid, record paymentRef
	// 5. Save to repository
	// 6. Send order confirmation notification
	_ = time.Now()
	return Order{}, errors.New("checkout unimplemented in starter")
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (Order, error) {
	if s.repo == nil {
		return Order{}, ErrOrderNotFound
	}
	return s.repo.FindByID(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]Order, error) {
	if s.repo == nil {
		return nil, nil
	}
	return s.repo.List(ctx)
}
