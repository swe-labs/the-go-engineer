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
		ID:            "core-01-04",
		Title:         "Terminal basics",
		MentalModel:   "The terminal is a conversation with the computer: you type a command, the shell interprets it, and the program writes back.",
		MachineView:   "A shell parses your command line, resolves the executable, passes arguments and environment variables, then waits for an exit code.",
		CommonMistake: "Typing commands from the wrong directory and blaming the tool.",
		Fix:           "Run `pwd` first and verify the path before running a command.",
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
