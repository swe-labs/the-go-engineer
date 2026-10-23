package main

import (
	"context"
)

// OrderRepository abstracts data persistence for Orders.
type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	FindByID(ctx context.Context, id string) (Order, error)
	List(ctx context.Context) ([]Order, error)
}

// PaymentGateway abstracts external merchant payment charge operations.
type PaymentGateway interface {
	Charge(ctx context.Context, amountCents int, paymentToken string) (string, error)
}

// NotificationService abstracts customer message delivery (e.g. Email/SMS/Webhook).
type NotificationService interface {
	SendOrderConfirmation(ctx context.Context, order Order) error
}
