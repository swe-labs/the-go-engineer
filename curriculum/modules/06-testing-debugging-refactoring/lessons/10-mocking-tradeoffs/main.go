package main

import "fmt"

// EmailSender sends emails.
type EmailSender interface {
	Send(to, subject, body string) error
}

// Notifier sends notifications via email.
type Notifier struct {
	sender EmailSender
	from   string
}

func NewNotifier(sender EmailSender, from string) *Notifier {
	return &Notifier{sender: sender, from: from}
}

func (n *Notifier) NotifyUser(email, message string) error {
	return n.sender.Send(email, "Notification", message)
}

// PaymentProcessor processes payments.
type PaymentProcessor interface {
	Charge(amount int, currency string) (string, error)
}

// Purchase charges for an item and returns the transaction ID.
func Purchase(processor PaymentProcessor, item string, price int) (string, error) {
	txnID, err := processor.Charge(price, "USD")
	if err != nil {
		return "", fmt.Errorf("purchase of %s failed: %w", item, err)
	}
	return txnID, nil
}

func main() {
	fmt.Println("See tests for mocking examples")
}
