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
		ID:             "core-02-02",
		Title:          "Hello World",
		MentalModel:    "A Hello World program is a doorbell test: it proves the toolchain, source file, package, and output path are connected.",
		MachineView:    "Go parses the file, compiles package `main`, links an executable, starts a process, and writes bytes to standard output.",
		CommandPurpose: "Write and run the smallest useful Go program while understanding package, import, function, and output at a beginner level.",
		CommonMistake:  "Thinking `fmt.Println` is the whole lesson and missing package/import/main structure.",
		Fix:            "Name all four moving parts: package, import, main function, standard output.",
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
