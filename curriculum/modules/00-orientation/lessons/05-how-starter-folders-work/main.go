package main

import (
	"fmt"
)

type Task struct {
	Name     string
	Desc     string
	Complete bool
}

type Comparison struct {
	MissingName bool
	MissingDesc bool
	IsComplete  bool
}

func CompareTasks(starter, solution Task) Comparison {
	return Comparison{
		MissingName: starter.Name != solution.Name,
		MissingDesc: starter.Desc != solution.Desc,
		IsComplete:  starter.Complete,
	}
}

func (c Comparison) Report() string {
	if c.MissingName || c.MissingDesc {
		return fmt.Sprintf("Starter is incomplete. Missing name: %v, missing desc: %v", c.MissingName, c.MissingDesc)
	}
	if c.IsComplete {
		return "Starter is complete and matches the solution."
	}
	return "Starter has all fields but is marked incomplete."
}

func main() {
	starter := Task{Name: "", Desc: "", Complete: false}
	solution := Task{Name: "Greet the user", Desc: "Write a function that prints a greeting", Complete: true}

	result := CompareTasks(starter, solution)
	fmt.Println("Starter:", starter)
	fmt.Println("Solution:", solution)
	fmt.Println("Comparison:", result.Report())

	fmt.Println()

	completed := Task{Name: "Greet the user", Desc: "Write a function that prints a greeting", Complete: true}
	result2 := CompareTasks(completed, solution)
	fmt.Println("Completed starter:", completed)
	fmt.Println("Comparison:", result2.Report())
}
