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
		ID:            "core-01-14",
		Title:         "Web preview: client, server, DNS, and ports",
		MentalModel:   "DNS is the address book, an IP is the building, and a port is the specific door where a service listens.",
		MachineView:   "A client resolves a hostname to an IP address, opens a network connection to a port, and sends protocol-specific bytes.",
		CommonMistake: "Assuming a server is broken when the DNS or port route is wrong.",
		Fix:           "Debug the path in layers: name, address, port, protocol, application.",
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
