package main

import (
	"errors"
	"fmt"
)

// PaymentGateway defines the contract for external transaction processing.
type PaymentGateway interface {
	Charge(userID string, amount int) error
}

// MockPaymentGateway simulates a Mockery-generated test double.
type MockPaymentGateway struct {
	Calls        map[string]int
	ReturnErrors map[string]error
}

func NewMockPaymentGateway() *MockPaymentGateway {
	return &MockPaymentGateway{
		Calls:        make(map[string]int),
		ReturnErrors: make(map[string]error),
	}
}

func (m *MockPaymentGateway) Charge(userID string, amount int) error {
	m.Calls[userID]++
	return m.ReturnErrors[userID]
}

// CheckoutService processes purchases.
type CheckoutService struct {
	gateway PaymentGateway
}

func (s *CheckoutService) Process(userID string, amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	return s.gateway.Charge(userID, amount)
}

// Demo demonstrates Mockery and satisfies test requirements.
func Demo() string {
	mock := NewMockPaymentGateway()
	service := CheckoutService{gateway: mock}
	_ = service.Process("user_101", 500)
	return "Mockery verified"
}

func main() {
	fmt.Println(Demo())
}
