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
		ID:       "core-00-07",
		Title:    "How to ask good debugging questions",
		Mission:  "Turn vague confusion into a structured question that contains evidence and a testable hypothesis.",
		Proof:    "Rewrite a vague question into a high-signal debugging question with reproduction steps.",
		NextStep: "core-00-08",
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
