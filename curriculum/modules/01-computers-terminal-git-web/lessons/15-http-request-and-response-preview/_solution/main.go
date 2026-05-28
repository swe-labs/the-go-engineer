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
		ID:            "core-01-15",
		Title:         "HTTP request and response preview",
		MentalModel:   "HTTP is a written exchange: the client asks with a method, path, headers, and optional body; the server replies with a status, headers, and body.",
		MachineView:   "HTTP bytes travel over a network connection. Servers parse the request, route it to handler logic, and serialize a response.",
		CommonMistake: "Thinking HTTP status codes are only display messages for users.",
		Fix:           "Treat status codes as machine-readable signals between clients and servers.",
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
