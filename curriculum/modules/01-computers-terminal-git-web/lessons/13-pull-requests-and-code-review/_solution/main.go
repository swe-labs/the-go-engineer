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
		ID:            "core-01-13",
		Title:         "Pull requests and code review",
		MentalModel:   "A pull request is a case file: it shows what changed, why it changed, how it was tested, and what reviewers should examine.",
		MachineView:   "A PR compares two Git references, runs checks, collects comments, and records review decisions before merge.",
		CommonMistake: "Opening a PR with no explanation and expecting reviewers to infer intent.",
		Fix:           "Give reviewers context, evidence, and a clear validation plan.",
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
