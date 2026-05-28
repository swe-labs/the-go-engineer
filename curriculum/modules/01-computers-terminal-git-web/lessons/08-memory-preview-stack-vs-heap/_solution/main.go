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
		ID:            "core-01-08",
		Title:         "Memory preview: stack vs heap",
		MentalModel:   "The stack is like a neat pile of call frames. The heap is like a larger shared storage room for values that need flexible lifetime.",
		MachineView:   "Go decides where values live using escape analysis. The garbage collector later reclaims heap values that are no longer reachable.",
		CommonMistake: "Thinking stack means fast and heap means bad.",
		Fix:           "Treat stack and heap as lifetime tools; measure before optimizing.",
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
