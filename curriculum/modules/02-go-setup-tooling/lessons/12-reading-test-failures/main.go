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
	Commands       []string
	NextStep       string
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
		Commands:       []string{"go test -run TestName -v", "go test ./..."},
		NextStep:       "core-02-13",
	}
}

func (c toolCard) summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "ID: %s\n", c.ID)
	fmt.Fprintf(&b, "Title: %s\n", c.Title)
	fmt.Fprintf(&b, "Mental model: %s\n", c.MentalModel)
	fmt.Fprintf(&b, "Machine view: %s\n", c.MachineView)
	fmt.Fprintf(&b, "Purpose: %s\n", c.CommandPurpose)
	fmt.Fprintf(&b, "Common mistake: %s\n", c.CommonMistake)
	fmt.Fprintf(&b, "Fix: %s\n", c.Fix)
	fmt.Fprintf(&b, "Commands: %s\n", strings.Join(c.Commands, ", "))
	fmt.Fprintf(&b, "Next: %s\n", c.NextStep)
	return b.String()
}

func main() {
	fmt.Print(card().summary())
}
