package main

import (
	"fmt"
	"strings"
)

type toolCard struct {
	ID             string
	Title          string
	MentalModel    string
	MachineView    string
	CommandPurpose string
	CommonMistake  string
	Fix            string
}

func card() toolCard {
	return toolCard{
		ID:             "core-02-12",
		Title:          "Reading test failures",
		MentalModel:    "A failing test is not a punishment. It is a precise disagreement between expectation and reality.",
		MachineView:    "The test binary reports which test failed, which assertion failed, and any messages written with `t.Fatal`, `t.Fatalf`, `t.Error`, or `t.Errorf`.",
		CommandPurpose: "Learn to use test failure output to identify expected behavior, actual behavior, and the smallest broken assumption.",
		CommonMistake:  "Fixing code randomly until tests pass.",
		Fix:            "Read the failing assertion and explain what behavior the test is protecting.",
	}
}

func render(c toolCard) string {
	return strings.Join([]string{
		"ID: " + c.ID,
		"Title: " + c.Title,
		"Mental model: " + c.MentalModel,
		"Machine view: " + c.MachineView,
		"Purpose: " + c.CommandPurpose,
		"Common mistake: " + c.CommonMistake,
		"Fix: " + c.Fix,
	}, "\n")
}

func main() {
	fmt.Println(render(card()))
}
