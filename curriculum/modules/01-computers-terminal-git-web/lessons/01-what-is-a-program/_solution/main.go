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
		ID:            "core-01-01",
		Title:         "What is a program?",
		MentalModel:   "A program is a recipe plus a kitchen. The recipe is the instructions; the kitchen is the computer that follows them.",
		MachineView:   "At runtime, a program is loaded by the operating system, given memory and handles to input/output, and executed instruction by instruction by the CPU.",
		CommonMistake: "Thinking code is the same thing as a running program.",
		Fix:           "Separate source code, executable artifact, and running process in your vocabulary.",
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
