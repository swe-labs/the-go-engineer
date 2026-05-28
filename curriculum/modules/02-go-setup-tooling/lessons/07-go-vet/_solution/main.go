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
	}
}

func render(c toolCard) string {
	return strings.Join([]string{
		"ID: " + c.ID,
		"Title: " + c.Title,
		"Mental model: " + c.MentalModel,
		"Machine view: " + c.MachineView,
		"Purpose: " + c.CommandPurpose,
		"Common mistake: " + c.CommonMistake,
		"Fix: " + c.Fix,
	}, "\n")
}

func main() {
	fmt.Println(render(card()))
}
