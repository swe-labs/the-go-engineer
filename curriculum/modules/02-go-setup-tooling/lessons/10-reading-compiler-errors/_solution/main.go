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
		ID:             "core-02-10",
		Title:          "Reading compiler errors",
		MentalModel:    "A compiler error is a map pin plus a complaint: here is where I stopped, and here is what did not make sense.",
		MachineView:    "The compiler parses and type-checks source files. When it cannot build a valid program, it reports file, line, column, and error message.",
		CommandPurpose: "Learn to read compiler errors as structured location-and-cause reports, not as scary walls of text.",
		CommonMistake:  "Reading only the final words of the error and ignoring file/line context.",
		Fix:            "Start with the first compiler error and inspect the exact source location.",
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
