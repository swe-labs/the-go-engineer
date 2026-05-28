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
		ID:             "core-02-08",
		Title:          "go doc",
		MentalModel:    "`go doc` is a local manual attached to your toolchain.",
		MachineView:    "The Go tool reads documentation comments and exported declarations from packages and prints the selected symbol or package documentation.",
		CommandPurpose: "Learn to inspect package documentation from the terminal instead of guessing API behavior.",
		CommonMistake:  "Searching randomly online before checking the documentation installed with the toolchain.",
		Fix:            "Use `go doc` for quick API shape, then use official docs for deeper examples.",
		Commands:       []string{"go doc fmt", "go doc fmt.Println", "go doc testing.T"},
		NextStep:       "core-02-09",
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
