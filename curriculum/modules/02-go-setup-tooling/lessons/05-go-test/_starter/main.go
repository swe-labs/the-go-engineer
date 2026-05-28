package main

import "fmt"

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
		ID:             "core-02-05",
		Title:          "go test",
		MentalModel:    "TODO: explain the learner-friendly model",
		MachineView:    "TODO: explain what the Go tool or machine does",
		CommandPurpose: "TODO: explain what proof this command gives",
		CommonMistake:  "TODO: name a likely setup/tooling mistake",
		Fix:            "TODO: explain the first recovery step",
	}
}

func main() {
	c := card()
	fmt.Println(c.ID)
	fmt.Println(c.Title)
	fmt.Println(c.MentalModel)
	fmt.Println(c.MachineView)
	fmt.Println(c.CommandPurpose)
	fmt.Println(c.CommonMistake)
	fmt.Println(c.Fix)
}
