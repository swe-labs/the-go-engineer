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
		ID:            "core-01-01",
		Title:         "What is a program?",
		MentalModel:   "A program is a recipe plus a kitchen. The recipe is the instructions; the kitchen is the computer that follows them.",
		MachineView:   "At runtime, a program is loaded by the operating system, given memory and handles to input/output, and executed instruction by instruction by the CPU.",
		CommonMistake: "Thinking code is the same thing as a running program.",
		Fix:           "Separate source code, executable artifact, and running process in your vocabulary.",
		Commands:      []string{"pwd", "ls", "go run ./curriculum/modules/01-computers-terminal-git-web/lessons/01-what-is-a-program"},
		NextStep:      "core-01-02",
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
