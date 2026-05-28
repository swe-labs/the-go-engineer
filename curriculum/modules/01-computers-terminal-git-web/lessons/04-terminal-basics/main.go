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
		ID:            "core-01-04",
		Title:         "Terminal basics",
		MentalModel:   "The terminal is a conversation with the computer: you type a command, the shell interprets it, and the program writes back.",
		MachineView:   "A shell parses your command line, resolves the executable, passes arguments and environment variables, then waits for an exit code.",
		CommonMistake: "Typing commands from the wrong directory and blaming the tool.",
		Fix:           "Run `pwd` first and verify the path before running a command.",
		Commands:      []string{"pwd", "cd", "ls", "echo"},
		NextStep:      "core-01-05",
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
