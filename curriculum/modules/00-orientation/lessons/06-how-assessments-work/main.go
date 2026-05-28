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
		ID:       "core-00-06",
		Title:    "How assessments work",
		Mission:  "Learn how checkpoints measure mastery and how to use rubrics as feedback.",
		Proof:    "Read the module checkpoint rubric and identify what evidence would prove each category.",
		NextStep: "core-00-07",
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
