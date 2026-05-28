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
		ID:            "core-01-03",
		Title:         "Files, bytes, directories, and paths",
		MentalModel:   "A filesystem is a labeled warehouse. Directories are shelves, files are boxes, and paths are written directions to a box.",
		MachineView:   "The operating system stores bytes and metadata. Programs ask the OS to open, read, write, and close files using paths.",
		CommonMistake: "Assuming a relative path starts from the file's directory.",
		Fix:           "Remember that relative paths are resolved from the current working directory of the process.",
		Commands:      []string{"pwd", "ls", "mkdir", "cat"},
		NextStep:      "core-01-04",
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
