package main

import (
	"fmt"
	"testing"
)

// MockEmailSender implements EmailSender for testing without testify.
type MockEmailSender struct {
	calls  []sendCall
	onSend func(to, subject, body string) error
}

type sendCall struct {
	to, subject, body string
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	m.calls = append(m.calls, sendCall{to, subject, body})
	if m.onSend != nil {
		return m.onSend(to, subject, body)
	}
	return nil
}

func (m *MockEmailSender) AssertCalled(t *testing.T, to, subject, body string) {
	t.Helper()
	for _, c := range m.calls {
		if c.to == to && c.subject == subject && c.body == body {
			return
		}
	}
	t.Errorf("expected Send(%q, %q, %q) to be called", to, subject, body)
}

// MockPaymentProcessor implements PaymentProcessor for testing.
type MockPaymentProcessor struct {
	onCharge func(amount int, currency string) (string, error)
}

func (m *MockPaymentProcessor) Charge(amount int, currency string) (string, error) {
	if m.onCharge != nil {
		return m.onCharge(amount, currency)
	}
	return "txn_123", nil
}

func TestNotifyUser(t *testing.T) {
	mockSender := &MockEmailSender{}

	notifier := NewNotifier(mockSender, "noreply@example.com")
	err := notifier.NotifyUser("alice@example.com", "Welcome!")

	if err != nil {
		t.Errorf("NotifyUser returned error: %v", err)
	}
	mockSender.AssertCalled(t, "alice@example.com", "Notification", "Welcome!")
}

func TestNotifyUserFails(t *testing.T) {
	mockSender := &MockEmailSender{
		onSend: func(to, subject, body string) error {
			return fmt.Errorf("send failed")
		},
	}

	notifier := NewNotifier(mockSender, "noreply@example.com")
	err := notifier.NotifyUser("bad@example.com", "Test")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	mockSender.AssertCalled(t, "bad@example.com", "Notification", "Test")
}

func TestPurchaseSuccess(t *testing.T) {
	processor := &MockPaymentProcessor{
		onCharge: func(amount int, currency string) (string, error) {
			if amount != 1999 {
				t.Errorf("Charge amount = %d; want 1999", amount)
			}
			if currency != "USD" {
				t.Errorf("Charge currency = %s; want USD", currency)
			}
			return "txn_abc123", nil
		},
	}

	txnID, err := Purchase(processor, "widget", 1999)
	if err != nil {
		t.Fatalf("Purchase failed: %v", err)
	}
	if txnID != "txn_abc123" {
		t.Errorf("txnID = %q; want %q", txnID, "txn_abc123")
	}
}

func TestPurchaseFailure(t *testing.T) {
	processor := &MockPaymentProcessor{
		onCharge: func(amount int, currency string) (string, error) {
			return "", fmt.Errorf("card declined")
		},
	}

	_, err := Purchase(processor, "widget", 1999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
