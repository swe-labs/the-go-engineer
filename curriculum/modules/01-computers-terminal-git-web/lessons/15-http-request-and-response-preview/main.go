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
		ID:            "core-01-15",
		Title:         "HTTP request and response preview",
		MentalModel:   "HTTP is a written exchange: the client asks with a method, path, headers, and optional body; the server replies with a status, headers, and body.",
		MachineView:   "HTTP bytes travel over a network connection. Servers parse the request, route it to handler logic, and serialize a response.",
		CommonMistake: "Thinking HTTP status codes are only display messages for users.",
		Fix:           "Treat status codes as machine-readable signals between clients and servers.",
		Commands:      []string{"curl -i", "GET /path HTTP/1.1", "HTTP/1.1 200 OK"},
		NextStep:      "module-02",
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
