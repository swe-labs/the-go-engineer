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
		ID:             "core-02-14",
		Title:          "Tooling checklist",
		MentalModel:    "A tooling checklist is a preflight inspection. It catches setup problems before they masquerade as programming problems.",
		MachineView:    "The checklist exercises the shell, Go binary, module root, formatter, compiler, tests, vet, docs, editor, and validation tools.",
		CommandPurpose: "Build a repeatable checklist for verifying Go tooling before writing larger programs.",
		CommonMistake:  "Skipping setup verification and spending hours debugging a broken environment.",
		Fix:            "Run the checklist first when starting on a new machine, branch, module, or CI image.",
		Commands:       []string{"go version", "go env GOMOD", "gofmt -l .", "go test ./...", "go vet ./..."},
		NextStep:       "module-03",
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
