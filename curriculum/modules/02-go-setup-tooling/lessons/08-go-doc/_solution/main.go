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
		ID:             "core-02-08",
		Title:          "go doc",
		MentalModel:    "`go doc` is a local manual attached to your toolchain.",
		MachineView:    "The Go tool reads documentation comments and exported declarations from packages and prints the selected symbol or package documentation.",
		CommandPurpose: "Learn to inspect package documentation from the terminal instead of guessing API behavior.",
		CommonMistake:  "Searching randomly online before checking the documentation installed with the toolchain.",
		Fix:            "Use `go doc` for quick API shape, then use official docs for deeper examples.",
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
