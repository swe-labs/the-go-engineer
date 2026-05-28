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
		ID:             "core-02-10",
		Title:          "Reading compiler errors",
		MentalModel:    "A compiler error is a map pin plus a complaint: here is where I stopped, and here is what did not make sense.",
		MachineView:    "The compiler parses and type-checks source files. When it cannot build a valid program, it reports file, line, column, and error message.",
		CommandPurpose: "Learn to read compiler errors as structured location-and-cause reports, not as scary walls of text.",
		CommonMistake:  "Reading only the final words of the error and ignoring file/line context.",
		Fix:            "Start with the first compiler error and inspect the exact source location.",
		Commands:       []string{"go test ./...", "go build ."},
		NextStep:       "core-02-11",
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
