package main

import "fmt"

type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
}

func (p Person) Greet() string {
	return "Hello, I'm " + p.Name
}

type Robot struct {
	Model string
}

func (r *Robot) Greet() string {
	return "Beep boop, model " + r.Model
}

type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + " says woof"
}

type Cat struct {
	Name string
}

func (c *Cat) Speak() string {
	return c.Name + " says meow"
}

func main() {
	var g Greeter

	g = Person{Name: "Alice"}
	fmt.Println(g.Greet())

	g = &Person{Name: "Bob"}
	fmt.Println(g.Greet())

	g = &Robot{Model: "R2"}
	fmt.Println(g.Greet())

	var s Speaker
	s = Dog{Name: "Rex"}
	fmt.Println(s.Speak())

	s = &Dog{Name: "Max"}
	fmt.Println(s.Speak())

	s = &Cat{Name: "Whiskers"}
	fmt.Println(s.Speak())
}
