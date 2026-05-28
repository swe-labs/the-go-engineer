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
		ID:             "core-02-04",
		Title:          "go build",
		MentalModel:    "`go build` turns the recipe into a packaged tool you can run later.",
		MachineView:    "The Go tool compiles packages, reuses the build cache when safe, links dependencies, and writes an executable for command packages.",
		CommandPurpose: "Understand `go build` as the command that produces reusable executable artifacts.",
		CommonMistake:  "Building from the wrong directory and wondering why the binary name or output path is surprising.",
		Fix:            "Know the package path and use `-o` when you need a predictable output filename.",
		Commands:       []string{"go build .", "go build -o ./tmp/tooling-demo ."},
		NextStep:       "core-02-05",
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
