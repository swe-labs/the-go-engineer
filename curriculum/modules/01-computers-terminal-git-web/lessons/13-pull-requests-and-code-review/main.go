package main

import (
	"fmt"
	"strings"
)

type conceptCard struct {
	ID            string
	Title         string
	MentalModel   string
	MachineView   string
	CommonMistake string
	Fix           string
	Commands      []string
	NextStep      string
}

func card() conceptCard {
	return conceptCard{
		ID:            "core-01-13",
		Title:         "Pull requests and code review",
		MentalModel:   "A pull request is a case file: it shows what changed, why it changed, how it was tested, and what reviewers should examine.",
		MachineView:   "A PR compares two Git references, runs checks, collects comments, and records review decisions before merge.",
		CommonMistake: "Opening a PR with no explanation and expecting reviewers to infer intent.",
		Fix:           "Give reviewers context, evidence, and a clear validation plan.",
		Commands:      []string{"git diff", "gh pr create", "gh pr status"},
		NextStep:      "core-01-14",
	}
}

func (c conceptCard) summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "ID: %s\n", c.ID)
	fmt.Fprintf(&b, "Title: %s\n", c.Title)
	fmt.Fprintf(&b, "Mental model: %s\n", c.MentalModel)
	fmt.Fprintf(&b, "Machine view: %s\n", c.MachineView)
	fmt.Fprintf(&b, "Common mistake: %s\n", c.CommonMistake)
	fmt.Fprintf(&b, "Fix: %s\n", c.Fix)
	fmt.Fprintf(&b, "Try commands: %s\n", strings.Join(c.Commands, ", "))
	fmt.Fprintf(&b, "Next: %s\n", c.NextStep)
	return b.String()
}

func main() {
	fmt.Print(card().summary())
}
