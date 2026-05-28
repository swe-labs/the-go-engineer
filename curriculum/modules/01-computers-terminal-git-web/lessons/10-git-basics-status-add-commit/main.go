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
		ID:            "core-01-10",
		Title:         "Git basics: status, add, commit",
		MentalModel:   "`git status` is the dashboard, `git add` prepares a snapshot, and `git commit` records that snapshot.",
		MachineView:   "Git compares the working tree, index, and last commit to show what changed and what is staged.",
		CommonMistake: "Running `git add .` without checking what changed.",
		Fix:           "Run `git status` and `git diff` before staging.",
		Commands:      []string{"git status", "git diff", "git add", "git commit"},
		NextStep:      "core-01-11",
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
