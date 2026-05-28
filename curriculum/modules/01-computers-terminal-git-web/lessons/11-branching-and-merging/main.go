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
		ID:            "core-01-11",
		Title:         "Branching and merging",
		MentalModel:   "A branch is a movable label. Merging joins two lines of work into one shared history.",
		MachineView:   "Git finds a common ancestor, compares changes on both sides, and creates a merge commit or fast-forward when possible.",
		CommonMistake: "Panicking when a merge conflict appears.",
		Fix:           "A conflict is a request for a human decision, not data loss by default.",
		Commands:      []string{"git switch -c", "git merge", "git log --oneline --graph"},
		NextStep:      "core-01-12",
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
