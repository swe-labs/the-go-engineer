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
		ID:             "core-02-11",
		Title:          "Reading runtime errors",
		MentalModel:    "A runtime error is a machine saying: the instructions were valid enough to start, but something went wrong while doing them.",
		MachineView:    "Runtime failures can come from panics, nil pointer use, missing files, invalid inputs, network failures, permissions, or explicit error returns.",
		CommandPurpose: "Understand runtime errors as failures that happen after a program starts executing.",
		CommonMistake:  "Treating every error as a compiler problem.",
		Fix:            "Ask when the failure happened: before build, during test, or while running.",
		Commands:       []string{"go run .", "panic", "stack trace"},
		NextStep:       "core-02-12",
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
