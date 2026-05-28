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
		ID:            "core-01-02",
		Title:         "Source code, executable, and process",
		MentalModel:   "Source code is a blueprint, an executable is the built machine, and a process is the machine turned on and doing work.",
		MachineView:   "The compiler turns source into machine instructions. The operating system starts a process by loading an executable, assigning memory, and scheduling CPU time.",
		CommonMistake: "Changing source code and assuming a running process changed automatically.",
		Fix:           "Rebuild and restart when the source change must affect a running process.",
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
