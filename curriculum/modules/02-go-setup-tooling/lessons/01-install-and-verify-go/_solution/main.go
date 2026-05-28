package main

import (
	"fmt"
	"strings"
)

type toolCard struct {
	ID             string
	Title          string
	MentalModel    string
	MachineView    string
	CommandPurpose string
	CommonMistake  string
	Fix            string
}

func card() toolCard {
	return toolCard{
		ID:             "core-02-01",
		Title:          "Install and verify Go",
		MentalModel:    "Installing Go is like putting a workshop on your machine: the `go` command is the front door, and `go env` tells you where the workshop keeps its tools.",
		MachineView:    "The shell finds the `go` executable through `PATH`. The toolchain then reads environment values such as `GOROOT`, `GOPATH`, `GOMODCACHE`, and `GOCACHE` to decide where tools, modules, and build artifacts live.",
		CommandPurpose: "Install Go and verify that the `go` command, version, environment, and workspace assumptions are visible.",
		CommonMistake:  "Installing Go but using a terminal session whose PATH does not include the Go binary.",
		Fix:            "Open a new terminal, verify `which go` or `command -v go`, and inspect `go env` before changing code.",
	}
}

func render(c toolCard) string {
	return strings.Join([]string{
		"ID: " + c.ID,
		"Title: " + c.Title,
		"Mental model: " + c.MentalModel,
		"Machine view: " + c.MachineView,
		"Purpose: " + c.CommandPurpose,
		"Common mistake: " + c.CommonMistake,
		"Fix: " + c.Fix,
	}, "\n")
}

func main() {
	fmt.Println(render(card()))
}
