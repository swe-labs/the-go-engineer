package main

import "testing"

func TestBookStringer(t *testing.T) {
	b := Book{Title: "1984", Author: "Orwell"}
	if b.String() != "1984 by Orwell" {
		t.Errorf("unexpected String(): %s", b.String())
	}
}

func TestStringerInterface(t *testing.T) {
	var s Stringer = Book{Title: "Test", Author: "Author"}
	if s.String() != "Test by Author" {
		t.Errorf("unexpected String() via interface: %s", s.String())
	}
}

func TestNilInterface(t *testing.T) {
	var s Stringer
	if s != nil {
		t.Error("uninitialized Stringer should be nil")
	}
}

func TestEmailNotifier(t *testing.T) {
	e := EmailNotifier{Address: "a@b.com"}
	err := e.Notify("hello")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmailNotifierEmptyAddress(t *testing.T) {
	e := EmailNotifier{}
	err := e.Notify("hello")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestSMSNotifier(t *testing.T) {
	s := SMSNotifier{Phone: "+123"}
	err := s.Notify("hello")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSMSNotifierEmptyPhone(t *testing.T) {
	s := SMSNotifier{}
	err := s.Notify("hello")
	if err == nil {
		t.Fatal("expected error for empty phone")
	}
}

func TestSendAlerts(t *testing.T) {
	notifiers := []Notifier{
		EmailNotifier{Address: "a@b.com"},
		SMSNotifier{Phone: "+123"},
	}
	SendAlerts(notifiers, "test message")
}

func TestEmptyInterface(t *testing.T) {
	var v interface{}
	v = 42
	if v != 42 {
		t.Errorf("expected 42, got %v", v)
	}
	v = "hello"
	if v != "hello" {
		t.Errorf("expected hello, got %v", v)
	}
}
