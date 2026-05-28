package main

import (
	"fmt"
	"strings"
)

type lessonCard struct {
	ID       string
	Title    string
	Mission  string
	Proof    string
	NextStep string
}

func card() lessonCard {
	return lessonCard{
		ID:       "core-00-01",
		Title:    "How to use this repository",
		Mission:  "Learn how the repository is organized and how metadata connects to learner-facing files.",
		Proof:    "Find three lesson IDs in metadata, locate their README paths, and explain the relationship between the JSON entry and the folder on disk.",
		NextStep: "core-00-02",
	}
}

func (c lessonCard) summary() string {
	lines := []string{
		"Go Engineer Orientation",
		"ID: " + c.ID,
		"Title: " + c.Title,
		"Mission: " + c.Mission,
		"Proof: " + c.Proof,
		"Next: " + c.NextStep,
	}
	return strings.Join(lines, "\n")
}

func main() {
	fmt.Println(card().summary())
}
