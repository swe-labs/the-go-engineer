package main

import (
	"errors"
	"fmt"
)

type Stringer interface {
	String() string
}

type Book struct {
	Title  string
	Author string
}

func (b Book) String() string {
	return b.Title + " by " + b.Author
}

func printAny(v interface{}) {
	fmt.Printf("value=%v type=%T\n", v, v)
}

type Notifier interface {
	Notify(message string) error
}

type EmailNotifier struct {
	Address string
}

func (e EmailNotifier) Notify(message string) error {
	if e.Address == "" {
		return errors.New("email address is empty")
	}
	fmt.Printf("Email to %s: %s\n", e.Address, message)
	return nil
}

type SMSNotifier struct {
	Phone string
}

func (s SMSNotifier) Notify(message string) error {
	if s.Phone == "" {
		return errors.New("phone number is empty")
	}
	fmt.Printf("SMS to %s: %s\n", s.Phone, message)
	return nil
}

func SendAlerts(notifiers []Notifier, message string) {
	for _, n := range notifiers {
		if err := n.Notify(message); err != nil {
			fmt.Println("notify error:", err)
		}
	}
}

func main() {
	var s Stringer
	fmt.Printf("nil? %v, type=%T\n", s == nil, s)

	s = Book{Title: "1984", Author: "Orwell"}
	fmt.Println(s.String())
	fmt.Printf("type=%T\n", s)

	var p *Book
	s = p
	fmt.Printf("typed nil assignment: nil? %v, type=%T\n", s == nil, s)

	var empty interface{}
	empty = 42
	printAny(empty)
	empty = "hello"
	printAny(empty)
	empty = Book{Title: "Brave New World", Author: "Huxley"}
	printAny(empty)

	notifiers := []Notifier{
		EmailNotifier{Address: "alice@example.com"},
		SMSNotifier{Phone: "+1234567890"},
	}
	SendAlerts(notifiers, "System maintenance at midnight")
}
