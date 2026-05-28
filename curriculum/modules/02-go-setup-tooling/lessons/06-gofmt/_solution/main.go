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
		ID:             "core-02-06",
		Title:          "gofmt",
		MentalModel:    "`gofmt` removes style debates by giving every Go file the same shape.",
		MachineView:    "`gofmt` parses Go source into syntax and prints it back in the canonical format.",
		CommandPurpose: "Use `gofmt` to make Go formatting automatic and non-negotiable.",
		CommonMistake:  "Manually formatting Go code and creating inconsistent diffs.",
		Fix:            "Run `gofmt` or configure your editor to format on save.",
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
