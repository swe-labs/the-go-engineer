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
		ID:             "core-02-05",
		Title:          "go test",
		MentalModel:    "`go test` is a referee. It builds your code and asks named test functions whether behavior matches expectations.",
		MachineView:    "The Go tool compiles the package and `_test.go` files into a temporary test binary, runs it, collects results, and reports pass or fail.",
		CommandPurpose: "Learn how `go test` compiles packages with test files and runs `Test...` functions as proof.",
		CommonMistake:  "Writing tests that only execute code without checking behavior.",
		Fix:            "Assert outcomes and include failure messages that explain what went wrong.",
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
