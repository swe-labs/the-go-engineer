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
		ID:            "core-01-08",
		Title:         "Memory preview: stack vs heap",
		MentalModel:   "The stack is like a neat pile of call frames. The heap is like a larger shared storage room for values that need flexible lifetime.",
		MachineView:   "Go decides where values live using escape analysis. The garbage collector later reclaims heap values that are no longer reachable.",
		CommonMistake: "Thinking stack means fast and heap means bad.",
		Fix:           "Treat stack and heap as lifetime tools; measure before optimizing.",
		Commands:      []string{"go test", "go build"},
		NextStep:      "core-01-09",
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
