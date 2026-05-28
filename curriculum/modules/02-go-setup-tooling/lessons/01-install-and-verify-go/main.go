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
	Commands       []string
	NextStep       string
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
		Commands:       []string{"go version", "go env", "command -v go"},
		NextStep:       "core-02-02",
	}
}

func (c toolCard) summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "ID: %s\n", c.ID)
	fmt.Fprintf(&b, "Title: %s\n", c.Title)
	fmt.Fprintf(&b, "Mental model: %s\n", c.MentalModel)
	fmt.Fprintf(&b, "Machine view: %s\n", c.MachineView)
	fmt.Fprintf(&b, "Purpose: %s\n", c.CommandPurpose)
	fmt.Fprintf(&b, "Common mistake: %s\n", c.CommonMistake)
	fmt.Fprintf(&b, "Fix: %s\n", c.Fix)
	fmt.Fprintf(&b, "Commands: %s\n", strings.Join(c.Commands, ", "))
	fmt.Fprintf(&b, "Next: %s\n", c.NextStep)
	return b.String()
}

func main() {
	fmt.Print(card().summary())
}
