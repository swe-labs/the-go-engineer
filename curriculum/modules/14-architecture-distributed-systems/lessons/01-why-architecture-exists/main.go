package main

import (
	"errors"
	"fmt"
)

type Order struct {
	ID     string
	Amount float64
	Status string
}

type BadOrderService struct{}

func (BadOrderService) Process(orderID string, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	fmt.Printf("BadOrderService: processing order %s for $%.2f\n", orderID, amount)
	return nil
}

type Logger interface {
	Log(message string)
}

type PaymentGateway interface {
	Charge(amount float64) error
}

type OrderValidator interface {
	Validate(amount float64) error
}

type OrderProcessor struct {
	logger    Logger
	payment   PaymentGateway
	validator OrderValidator
}

func NewOrderProcessor(logger Logger, payment PaymentGateway, validator OrderValidator) *OrderProcessor {
	return &OrderProcessor{logger: logger, payment: payment, validator: validator}
}

func (p *OrderProcessor) Process(orderID string, amount float64) error {
	p.logger.Log("processing order " + orderID)
	if err := p.validator.Validate(amount); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if err := p.payment.Charge(amount); err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}
	p.logger.Log("order " + orderID + " processed successfully")
	return nil
}

type ConsoleLogger struct{}

func (ConsoleLogger) Log(message string) {
	fmt.Println("LOG:", message)
}

type StripeGateway struct{}

func (StripeGateway) Charge(amount float64) error {
	fmt.Printf("Stripe: charging $%.2f\n", amount)
	return nil
}

type AmountValidator struct{}

func (AmountValidator) Validate(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return nil
}

func main() {
	tight := BadOrderService{}
	if err := tight.Process("ORD-001", 49.99); err != nil {
		fmt.Println("Error:", err)
	}

	loose := NewOrderProcessor(ConsoleLogger{}, StripeGateway{}, AmountValidator{})
	if err := loose.Process("ORD-002", 99.99); err != nil {
		fmt.Println("Error:", err)
	}
}
