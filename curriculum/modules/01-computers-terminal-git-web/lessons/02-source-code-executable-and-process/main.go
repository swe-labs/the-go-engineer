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
		ID:            "core-01-02",
		Title:         "Source code, executable, and process",
		MentalModel:   "Source code is a blueprint, an executable is the built machine, and a process is the machine turned on and doing work.",
		MachineView:   "The compiler turns source into machine instructions. The operating system starts a process by loading an executable, assigning memory, and scheduling CPU time.",
		CommonMistake: "Changing source code and assuming a running process changed automatically.",
		Fix:           "Rebuild and restart when the source change must affect a running process.",
		Commands:      []string{"go build", "ps", "kill"},
		NextStep:      "core-01-03",
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
