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
		ID:            "core-01-11",
		Title:         "Branching and merging",
		MentalModel:   "A branch is a movable label. Merging joins two lines of work into one shared history.",
		MachineView:   "Git finds a common ancestor, compares changes on both sides, and creates a merge commit or fast-forward when possible.",
		CommonMistake: "Panicking when a merge conflict appears.",
		Fix:           "A conflict is a request for a human decision, not data loss by default.",
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
