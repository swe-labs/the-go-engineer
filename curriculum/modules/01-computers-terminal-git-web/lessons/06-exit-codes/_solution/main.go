package main

import (
	"fmt"
	"strings"
)

type conceptCard struct {
	ID            string
	Title         string
	MentalModel   string
	MachineView   string
	CommonMistake string
	Fix           string
}

func card() conceptCard {
	return conceptCard{
		ID:            "core-01-06",
		Title:         "Exit codes",
		MentalModel:   "An exit code is a program's final vote: zero means success, non-zero means something failed.",
		MachineView:   "When a process exits, the OS records a small integer status. The parent process can read it and decide what to do next.",
		CommonMistake: "Printing an error message but exiting with success.",
		Fix:           "Return or exit with a non-zero code when the command failed.",
	}
}

func render(c conceptCard) string {
	return strings.Join([]string{
		"ID: " + c.ID,
		"Title: " + c.Title,
		"Mental model: " + c.MentalModel,
		"Machine view: " + c.MachineView,
		"Common mistake: " + c.CommonMistake,
		"Fix: " + c.Fix,
	}, "\n")
}

func main() {
	fmt.Println(render(card()))
}
