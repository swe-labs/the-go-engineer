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
		ID:             "core-02-04",
		Title:          "go build",
		MentalModel:    "`go build` turns the recipe into a packaged tool you can run later.",
		MachineView:    "The Go tool compiles packages, reuses the build cache when safe, links dependencies, and writes an executable for command packages.",
		CommandPurpose: "Understand `go build` as the command that produces reusable executable artifacts.",
		CommonMistake:  "Building from the wrong directory and wondering why the binary name or output path is surprising.",
		Fix:            "Know the package path and use `-o` when you need a predictable output filename.",
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
