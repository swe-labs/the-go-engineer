package main

import (
	"fmt"
)

type Step struct {
	Index     int
	Name      string
	Action    string
	Artifact  string
}

type LearningCycle struct {
	Title    string
	Steps    []Step
	Complete bool
}

func NewLearningCycle(title string) LearningCycle {
	return LearningCycle{
		Title: title,
		Steps: []Step{
			{Index: 1, Name: "Read", Action: "Read the lesson README and example code", Artifact: "notes or highlights"},
			{Index: 2, Name: "Run", Action: "Run the Go program and observe output", Artifact: "terminal output"},
			{Index: 3, Name: "Try", Action: "Complete the practice task in _starter/", Artifact: "modified code"},
			{Index: 4, Name: "Test", Action: "Run go test to verify your solution", Artifact: "PASS result"},
			{Index: 5, Name: "Reflect", Action: "Answer review questions from the README", Artifact: "written answers"},
			{Index: 6, Name: "Checkpoint", Action: "Run the module checkpoint or assessment", Artifact: "completion proof"},
		},
		Complete: false,
	}
}

func (lc *LearningCycle) Advance() bool {
	if len(lc.Steps) == 0 {
		return false
	}
	lc.Steps = lc.Steps[1:]
	if len(lc.Steps) == 0 {
		lc.Complete = true
	}
	return true
}

func main() {
	cycle := NewLearningCycle("How lessons, exercises, projects, and checkpoints work")
	fmt.Println("Learning Cycle:", cycle.Title)
	fmt.Println()
	for _, s := range cycle.Steps {
		fmt.Printf("  Step %d: %s\n", s.Index, s.Name)
		fmt.Printf("          %s\n", s.Action)
		fmt.Printf("          Artifact: %s\n", s.Artifact)
		fmt.Println()
	}
	fmt.Println("Complete:", cycle.Complete)
}
