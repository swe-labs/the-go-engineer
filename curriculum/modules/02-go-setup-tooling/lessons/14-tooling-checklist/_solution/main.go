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
		ID:             "core-02-14",
		Title:          "Tooling checklist",
		MentalModel:    "A tooling checklist is a preflight inspection. It catches setup problems before they masquerade as programming problems.",
		MachineView:    "The checklist exercises the shell, Go binary, module root, formatter, compiler, tests, vet, docs, editor, and validation tools.",
		CommandPurpose: "Build a repeatable checklist for verifying Go tooling before writing larger programs.",
		CommonMistake:  "Skipping setup verification and spending hours debugging a broken environment.",
		Fix:            "Run the checklist first when starting on a new machine, branch, module, or CI image.",
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
