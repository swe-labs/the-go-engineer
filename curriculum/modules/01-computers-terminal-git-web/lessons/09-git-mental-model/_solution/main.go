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
		ID:            "core-01-09",
		Title:         "Git mental model",
		MentalModel:   "Git is a timeline of labeled snapshots. Each commit points backward to previous history.",
		MachineView:   "Git stores objects by content hash, tracks names with references, and lets branches move to different commits.",
		CommonMistake: "Thinking a branch contains files rather than pointing to a commit.",
		Fix:           "Picture branches as movable labels on commits.",
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
