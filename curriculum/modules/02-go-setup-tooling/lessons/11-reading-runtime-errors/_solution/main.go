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
		ID:             "core-02-11",
		Title:          "Reading runtime errors",
		MentalModel:    "A runtime error is a machine saying: the instructions were valid enough to start, but something went wrong while doing them.",
		MachineView:    "Runtime failures can come from panics, nil pointer use, missing files, invalid inputs, network failures, permissions, or explicit error returns.",
		CommandPurpose: "Understand runtime errors as failures that happen after a program starts executing.",
		CommonMistake:  "Treating every error as a compiler problem.",
		Fix:            "Ask when the failure happened: before build, during test, or while running.",
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
