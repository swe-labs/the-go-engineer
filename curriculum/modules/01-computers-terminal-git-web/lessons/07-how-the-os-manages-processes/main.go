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
	Commands      []string
	NextStep      string
}

func card() conceptCard {
	return conceptCard{
		ID:            "core-01-07",
		Title:         "How the OS manages processes",
		MentalModel:   "The operating system is a traffic controller. It lets many processes share CPU, memory, files, and network devices without crashing into each other.",
		MachineView:   "A process has an ID, memory, open file descriptors, environment, arguments, and scheduling state. The OS switches between runnable processes.",
		CommonMistake: "Assuming a program can run forever without OS supervision or resource limits.",
		Fix:           "Learn to inspect processes, handle shutdown signals, and design cleanup paths.",
		Commands:      []string{"ps", "top", "kill"},
		NextStep:      "core-01-08",
	}
}

func (c conceptCard) summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "ID: %s\n", c.ID)
	fmt.Fprintf(&b, "Title: %s\n", c.Title)
	fmt.Fprintf(&b, "Mental model: %s\n", c.MentalModel)
	fmt.Fprintf(&b, "Machine view: %s\n", c.MachineView)
	fmt.Fprintf(&b, "Common mistake: %s\n", c.CommonMistake)
	fmt.Fprintf(&b, "Fix: %s\n", c.Fix)
	fmt.Fprintf(&b, "Try commands: %s\n", strings.Join(c.Commands, ", "))
	fmt.Fprintf(&b, "Next: %s\n", c.NextStep)
	return b.String()
}

func main() {
	fmt.Print(card().summary())
}
