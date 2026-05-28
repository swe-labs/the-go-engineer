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
		ID:            "core-01-10",
		Title:         "Git basics: status, add, commit",
		MentalModel:   "`git status` is the dashboard, `git add` prepares a snapshot, and `git commit` records that snapshot.",
		MachineView:   "Git compares the working tree, index, and last commit to show what changed and what is staged.",
		CommonMistake: "Running `git add .` without checking what changed.",
		Fix:           "Run `git status` and `git diff` before staging.",
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
