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
		ID:            "core-01-07",
		Title:         "How the OS manages processes",
		MentalModel:   "The operating system is a traffic controller. It lets many processes share CPU, memory, files, and network devices without crashing into each other.",
		MachineView:   "A process has an ID, memory, open file descriptors, environment, arguments, and scheduling state. The OS switches between runnable processes.",
		CommonMistake: "Assuming a program can run forever without OS supervision or resource limits.",
		Fix:           "Learn to inspect processes, handle shutdown signals, and design cleanup paths.",
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
