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
		ID:             "core-02-06",
		Title:          "gofmt",
		MentalModel:    "`gofmt` removes style debates by giving every Go file the same shape.",
		MachineView:    "`gofmt` parses Go source into syntax and prints it back in the canonical format.",
		CommandPurpose: "Use `gofmt` to make Go formatting automatic and non-negotiable.",
		CommonMistake:  "Manually formatting Go code and creating inconsistent diffs.",
		Fix:            "Run `gofmt` or configure your editor to format on save.",
		Commands:       []string{"gofmt -w main.go", "gofmt -l ."},
		NextStep:       "core-02-07",
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
