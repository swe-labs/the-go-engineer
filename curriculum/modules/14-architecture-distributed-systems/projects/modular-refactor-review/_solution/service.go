package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
	order := Order{
		ID:        generateID(),
		Customer:  dto.Customer,
		Email:     dto.Email,
		Items:     dto.Items,
		Status:    StatusPending,
		CreatedAt: time.Now().UTC(),
	}

	if err := order.CalculateTotals(); err != nil {
		return Order{}, err
	}

	paymentRef, err := s.payment.Charge(ctx, order.Total, dto.PaymentToken)
	if err != nil {
		return Order{}, fmt.Errorf("%w: %v", ErrPaymentFailed, err)
	}

	order.Status = StatusPaid
	order.PaymentRef = paymentRef

	if err := s.repo.Save(ctx, order); err != nil {
		return Order{}, err
	}

	_ = s.notifier.SendOrderConfirmation(ctx, order)
	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]Order, error) {
	return s.repo.List(ctx)
}
