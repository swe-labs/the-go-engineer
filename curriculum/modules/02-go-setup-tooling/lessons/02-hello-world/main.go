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
		ID:             "core-02-02",
		Title:          "Hello World",
		MentalModel:    "A Hello World program is a doorbell test: it proves the toolchain, source file, package, and output path are connected.",
		MachineView:    "Go parses the file, compiles package `main`, links an executable, starts a process, and writes bytes to standard output.",
		CommandPurpose: "Write and run the smallest useful Go program while understanding package, import, function, and output at a beginner level.",
		CommonMistake:  "Thinking `fmt.Println` is the whole lesson and missing package/import/main structure.",
		Fix:            "Name all four moving parts: package, import, main function, standard output.",
		Commands:       []string{"go run .", "go test ."},
		NextStep:       "core-02-03",
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
