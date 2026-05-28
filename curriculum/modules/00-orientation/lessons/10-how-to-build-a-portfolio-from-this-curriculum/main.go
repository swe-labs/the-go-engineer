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
		ID:       "core-00-10",
		Title:    "How to build a portfolio from this curriculum",
		Mission:  "Learn how to turn curriculum projects into professional evidence without overselling or underexplaining your work.",
		Proof:    "Create a portfolio plan with three future projects, the proof each should show, and the story each project should tell.",
		NextStep: "module-01",
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
