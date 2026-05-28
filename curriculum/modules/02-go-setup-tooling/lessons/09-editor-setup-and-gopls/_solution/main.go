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
		ID:             "core-02-09",
		Title:          "Editor setup and gopls",
		MentalModel:    "Your editor is a cockpit. `gopls` provides instruments, but the actual aircraft is still the Go toolchain.",
		MachineView:    "`gopls` reads your module, parses packages, tracks symbols, reports diagnostics, and powers completion, navigation, rename, and quick fixes.",
		CommandPurpose: "Understand what your editor and Go language server do, and how to verify they are helping rather than hiding problems.",
		CommonMistake:  "Believing an editor warning is the compiler, or ignoring terminal commands because the editor looks green.",
		Fix:            "Use editor diagnostics as fast feedback, then confirm with `go test`, `go vet`, and validators.",
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
