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
		ID:             "core-02-09",
		Title:          "Editor setup and gopls",
		MentalModel:    "Your editor is a cockpit. `gopls` provides instruments, but the actual aircraft is still the Go toolchain.",
		MachineView:    "`gopls` reads your module, parses packages, tracks symbols, reports diagnostics, and powers completion, navigation, rename, and quick fixes.",
		CommandPurpose: "Understand what your editor and Go language server do, and how to verify they are helping rather than hiding problems.",
		CommonMistake:  "Believing an editor warning is the compiler, or ignoring terminal commands because the editor looks green.",
		Fix:            "Use editor diagnostics as fast feedback, then confirm with `go test`, `go vet`, and validators.",
		Commands:       []string{"gopls version", "go env GOMOD", "go test ./..."},
		NextStep:       "core-02-10",
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
