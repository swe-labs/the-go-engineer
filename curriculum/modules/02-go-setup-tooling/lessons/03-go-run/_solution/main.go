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
		ID:             "core-02-03",
		Title:          "go run",
		MentalModel:    "`go run` is a temporary build and launch button.",
		MachineView:    "The Go tool resolves packages, compiles them into a temporary executable, runs that executable, and removes or reuses build cache artifacts afterward.",
		CommandPurpose: "Understand `go run` as compile-then-execute for quick feedback, not as an interpreter.",
		CommonMistake:  "Treating `go run` like a script interpreter and expecting it to skip compilation errors.",
		Fix:            "Remember that Go is compiled; `go run` includes a compile step.",
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
