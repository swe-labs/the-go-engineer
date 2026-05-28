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
		ID:             "core-02-13",
		Title:          "Go module root basics",
		MentalModel:    "`go.mod` is the address label and dependency contract for a Go workspace.",
		MachineView:    "The Go command searches for `go.mod` to determine the module root, module path, Go version, dependencies, and module cache behavior.",
		CommandPurpose: "Understand `go.mod` as the root declaration for a Go module and learn why commands behave differently inside and outside it.",
		CommonMistake:  "Running Go commands outside the module and blaming packages or imports.",
		Fix:            "Use `go env GOMOD`, `pwd`, and the repository root before debugging imports.",
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
