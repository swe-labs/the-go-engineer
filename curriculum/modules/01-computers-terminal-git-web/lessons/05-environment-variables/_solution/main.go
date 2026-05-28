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
		ID:            "core-01-05",
		Title:         "Environment variables",
		MentalModel:   "Environment variables are labeled notes handed to a program at startup.",
		MachineView:   "When a process starts, it receives a map of string keys and values. The program can read those values but changes do not automatically flow backward to the parent shell.",
		CommonMistake: "Treating environment variables as secure storage.",
		Fix:           "Use secrets management for sensitive values and avoid logging environment values blindly.",
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
