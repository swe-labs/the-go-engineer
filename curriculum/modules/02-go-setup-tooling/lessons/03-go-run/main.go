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
		ID:             "core-02-03",
		Title:          "go run",
		MentalModel:    "`go run` is a temporary build and launch button.",
		MachineView:    "The Go tool resolves packages, compiles them into a temporary executable, runs that executable, and removes or reuses build cache artifacts afterward.",
		CommandPurpose: "Understand `go run` as compile-then-execute for quick feedback, not as an interpreter.",
		CommonMistake:  "Treating `go run` like a script interpreter and expecting it to skip compilation errors.",
		Fix:            "Remember that Go is compiled; `go run` includes a compile step.",
		Commands:       []string{"go run .", "go run ./path/to/package"},
		NextStep:       "core-02-04",
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
