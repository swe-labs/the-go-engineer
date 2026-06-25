package main

import "testing"

func TestPersonPromotedField(t *testing.T) {
	p := Person{
		Name:    "Alice",
		Address: Address{City: "Portland", State: "OR"},
	}
	if p.City != "Portland" {
		t.Errorf("expected Portland, got %s", p.City)
	}
}

func TestPersonPromotedMethod(t *testing.T) {
	p := Person{
		Name:    "Alice",
		Address: Address{City: "Portland", State: "OR"},
	}
	if p.Full() != "Portland, OR" {
		t.Errorf("expected Portland, OR, got %s", p.Full())
	}
}

func TestPersonExplicitField(t *testing.T) {
	p := Person{
		Name:    "Alice",
		Address: Address{City: "Portland", State: "OR"},
	}
	if p.Address.City != "Portland" {
		t.Errorf("expected Portland via explicit path, got %s", p.Address.City)
	}
}

func TestEmployeeDoublePromotion(t *testing.T) {
	e := Employee{
		Person:   Person{Name: "Bob", Address: Address{City: "Seattle", State: "WA"}},
		Position: "Engineer",
	}
	if e.Name != "Bob" || e.Position != "Engineer" {
		t.Errorf("unexpected Employee fields: %+v", e)
	}
	if e.Full() != "Seattle, WA" {
		t.Errorf("expected Seattle, WA, got %s", e.Full())
	}
}

func TestContactInfo(t *testing.T) {
	c := Contact{
		Address: Address{City: "Austin", State: "TX"},
		Email:   "test@example.com",
	}
	info := c.Info()
	expected := "test@example.com - Austin, TX"
	if info != expected {
		t.Errorf("expected %q, got %q", expected, info)
	}
}
