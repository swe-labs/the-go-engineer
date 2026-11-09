package main

import "testing"

func TestPersonValueSatisfiesGreeter(t *testing.T) {
	var g Greeter = Person{Name: "Alice"}
	if g.Greet() != "Hello, I'm Alice" {
		t.Errorf("unexpected greet: %s", g.Greet())
	}
}

func TestPersonPointerSatisfiesGreeter(t *testing.T) {
	var g Greeter = &Person{Name: "Bob"}
	if g.Greet() != "Hello, I'm Bob" {
		t.Errorf("unexpected greet: %s", g.Greet())
	}
}

func TestRobotPointerSatisfiesGreeter(t *testing.T) {
	var g Greeter = &Robot{Model: "R2"}
	if g.Greet() != "Beep boop, model R2" {
		t.Errorf("unexpected greet: %s", g.Greet())
	}
}

func TestDogValueSatisfiesSpeaker(t *testing.T) {
	var s Speaker = Dog{Name: "Rex"}
	if s.Speak() != "Rex says woof" {
		t.Errorf("unexpected speak: %s", s.Speak())
	}
}

func TestDogPointerSatisfiesSpeaker(t *testing.T) {
	var s Speaker = &Dog{Name: "Max"}
	if s.Speak() != "Max says woof" {
		t.Errorf("unexpected speak: %s", s.Speak())
	}
}

func TestCatPointerSatisfiesSpeaker(t *testing.T) {
	var s Speaker = &Cat{Name: "Whiskers"}
	if s.Speak() != "Whiskers says meow" {
		t.Errorf("unexpected speak: %s", s.Speak())
	}
}
