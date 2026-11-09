package main

import "fmt"

type Address struct {
	City, State string
}

func (a Address) Full() string {
	return a.City + ", " + a.State
}

type Person struct {
	Name string
	Address
}

type Employee struct {
	Person
	Position string
}

type Contact struct {
	Address
	Email string
}

func (c Contact) Info() string {
	return c.Email + " - " + c.Full()
}

func main() {
	p := Person{
		Name:    "Alice",
		Address: Address{City: "Portland", State: "OR"},
	}
	fmt.Println(p.Name)
	fmt.Println(p.City)
	fmt.Println(p.Address.City)
	fmt.Println(p.Full())

	e := Employee{
		Person:   Person{Name: "Bob", Address: Address{City: "Seattle", State: "WA"}},
		Position: "Engineer",
	}
	fmt.Println(e.Name, e.Position)
	fmt.Println(e.Full())

	c := Contact{
		Address: Address{City: "Austin", State: "TX"},
		Email:   "alice@example.com",
	}
	fmt.Println(c.Info())
}
