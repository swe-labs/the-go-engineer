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
		ID:            "core-01-12",
		Title:         "GitHub workflow",
		MentalModel:   "Your local repo is your workshop. GitHub is a shared copy where teammates can see, discuss, and integrate work.",
		MachineView:   "`git push` sends commits to a remote. `git fetch` downloads remote references. `git pull` fetches then integrates.",
		CommonMistake: "Using `pull` and `push` as magic sync buttons.",
		Fix:           "Know which direction commits are moving and whether integration is happening.",
		Commands:      []string{"git clone", "git remote -v", "git fetch", "git push"},
		NextStep:      "core-01-13",
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
