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
		ID:            "core-01-03",
		Title:         "Files, bytes, directories, and paths",
		MentalModel:   "A filesystem is a labeled warehouse. Directories are shelves, files are boxes, and paths are written directions to a box.",
		MachineView:   "The operating system stores bytes and metadata. Programs ask the OS to open, read, write, and close files using paths.",
		CommonMistake: "Assuming a relative path starts from the file's directory.",
		Fix:           "Remember that relative paths are resolved from the current working directory of the process.",
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
