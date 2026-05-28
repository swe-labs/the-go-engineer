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
		ID:            "core-01-12",
		Title:         "GitHub workflow",
		MentalModel:   "Your local repo is your workshop. GitHub is a shared copy where teammates can see, discuss, and integrate work.",
		MachineView:   "`git push` sends commits to a remote. `git fetch` downloads remote references. `git pull` fetches then integrates.",
		CommonMistake: "Using `pull` and `push` as magic sync buttons.",
		Fix:           "Know which direction commits are moving and whether integration is happening.",
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
