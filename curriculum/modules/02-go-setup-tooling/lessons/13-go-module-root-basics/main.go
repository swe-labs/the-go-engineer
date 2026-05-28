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
		ID:             "core-02-13",
		Title:          "Go module root basics",
		MentalModel:    "`go.mod` is the address label and dependency contract for a Go workspace.",
		MachineView:    "The Go command searches for `go.mod` to determine the module root, module path, Go version, dependencies, and module cache behavior.",
		CommandPurpose: "Understand `go.mod` as the root declaration for a Go module and learn why commands behave differently inside and outside it.",
		CommonMistake:  "Running Go commands outside the module and blaming packages or imports.",
		Fix:            "Use `go env GOMOD`, `pwd`, and the repository root before debugging imports.",
		Commands:       []string{"go env GOMOD", "go list ./...", "go mod tidy"},
		NextStep:       "core-02-14",
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
