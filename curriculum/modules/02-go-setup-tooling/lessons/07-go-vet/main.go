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
		ID:             "core-02-07",
		Title:          "go vet",
		MentalModel:    "`go vet` is a cautious reviewer that looks for code smells with likely bugs.",
		MachineView:    "`go vet` analyzes Go syntax and type information to detect risky constructs, such as malformed format strings or unreachable patterns.",
		CommandPurpose: "Use `go vet` as a static analyzer that catches suspicious code patterns tests may miss.",
		CommonMistake:  "Thinking `go vet` proves the program is correct.",
		Fix:            "Use vet with tests, review, and runtime validation; each catches different failure modes.",
		Commands:       []string{"go vet .", "go vet ./..."},
		NextStep:       "core-02-08",
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
