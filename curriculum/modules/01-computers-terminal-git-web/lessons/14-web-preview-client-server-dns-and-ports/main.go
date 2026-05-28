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
		ID:            "core-01-14",
		Title:         "Web preview: client, server, DNS, and ports",
		MentalModel:   "DNS is the address book, an IP is the building, and a port is the specific door where a service listens.",
		MachineView:   "A client resolves a hostname to an IP address, opens a network connection to a port, and sends protocol-specific bytes.",
		CommonMistake: "Assuming a server is broken when the DNS or port route is wrong.",
		Fix:           "Debug the path in layers: name, address, port, protocol, application.",
		Commands:      []string{"curl", "dig", "nc"},
		NextStep:      "core-01-15",
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
