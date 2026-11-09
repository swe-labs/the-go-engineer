package main

import (
	"errors"
	"testing"
)

type mockLogger struct {
	lastMessage string
}

func (m *mockLogger) Log(message string) {
	m.lastMessage = message
}

type mockPayment struct {
	shouldFail bool
}

func (m *mockPayment) Charge(amount float64) error {
	if m.shouldFail {
		return errors.New("payment declined")
	}
	return nil
}

type mockValidator struct {
	shouldFail bool
}

func (m *mockValidator) Validate(amount float64) error {
	if m.shouldFail {
		return errors.New("invalid amount")
	}
	return nil
}

func TestOrderProcessor_ValidOrder(t *testing.T) {
	logger := &mockLogger{}
	processor := NewOrderProcessor(logger, &mockPayment{}, &mockValidator{})
	err := processor.Process("ORD-001", 50.00)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if logger.lastMessage != "order ORD-001 processed successfully" {
		t.Errorf("unexpected log: %s", logger.lastMessage)
	}
}

func TestOrderProcessor_ValidationFails(t *testing.T) {
	processor := NewOrderProcessor(&mockLogger{}, &mockPayment{}, &mockValidator{shouldFail: true})
	err := processor.Process("ORD-002", 50.00)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderProcessor_PaymentFails(t *testing.T) {
	processor := NewOrderProcessor(&mockLogger{}, &mockPayment{shouldFail: true}, &mockValidator{})
	err := processor.Process("ORD-003", 50.00)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBadOrderService_ValidOrder(t *testing.T) {
	s := BadOrderService{}
	err := s.Process("ORD-004", 25.00)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBadOrderService_InvalidAmount(t *testing.T) {
	s := BadOrderService{}
	err := s.Process("ORD-005", 0)
	if err == nil {
		t.Fatal("expected error for zero amount")
	}
}
