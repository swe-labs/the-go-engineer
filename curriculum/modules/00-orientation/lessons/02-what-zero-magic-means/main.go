package main

import (
	"fmt"
	"strings"
)

type lessonCard struct {
	ID       string
	Title    string
	Mission  string
	Proof    string
	NextStep string
}

func card() lessonCard {
	return lessonCard{
		ID:       "core-00-02",
		Title:    "What zero magic means",
		Mission:  "Understand the curriculum promise that no required idea should appear without explanation, proof, and context.",
		Proof:    "Take one compact expression or command and expand every assumption it hides.",
		NextStep: "core-00-03",
	}
}

func (c lessonCard) summary() string {
	lines := []string{
		"Go Engineer Orientation",
		"ID: " + c.ID,
		"Title: " + c.Title,
		"Mission: " + c.Mission,
		"Proof: " + c.Proof,
		"Next: " + c.NextStep,
	}
	return strings.Join(lines, "\n")
}

func main() {
	fmt.Println(card().summary())
}
