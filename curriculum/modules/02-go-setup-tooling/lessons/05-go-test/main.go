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
		ID:             "core-02-05",
		Title:          "go test",
		MentalModel:    "`go test` is a referee. It builds your code and asks named test functions whether behavior matches expectations.",
		MachineView:    "The Go tool compiles the package and `_test.go` files into a temporary test binary, runs it, collects results, and reports pass or fail.",
		CommandPurpose: "Learn how `go test` compiles packages with test files and runs `Test...` functions as proof.",
		CommonMistake:  "Writing tests that only execute code without checking behavior.",
		Fix:            "Assert outcomes and include failure messages that explain what went wrong.",
		Commands:       []string{"go test .", "go test ./...", "go test -run TestName"},
		NextStep:       "core-02-06",
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
