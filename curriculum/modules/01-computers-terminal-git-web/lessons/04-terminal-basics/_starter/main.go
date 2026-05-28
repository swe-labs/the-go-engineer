package main

import "fmt"

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
		ID:            "core-01-04",
		Title:         "Terminal basics",
		MentalModel:   "TODO: describe the human mental model",
		MachineView:   "TODO: describe what the machine or system sees",
		CommonMistake: "TODO: name a likely beginner mistake",
		Fix:           "TODO: explain the correction",
	}
}

func main() {
	c := card()
	fmt.Println(c.ID)
	fmt.Println(c.Title)
	fmt.Println(c.MentalModel)
	fmt.Println(c.MachineView)
	fmt.Println(c.CommonMistake)
	fmt.Println(c.Fix)
}
