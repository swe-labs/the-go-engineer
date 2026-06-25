package main

import "fmt"

type cycleStep struct {
	Name    string
	Purpose string
}

func steps() []cycleStep {
	return []cycleStep{
		{Name: "read", Purpose: "TODO"},
		{Name: "run", Purpose: "TODO"},
		{Name: "try", Purpose: "TODO"},
		{Name: "test", Purpose: "TODO"},
		{Name: "reflect", Purpose: "TODO"},
		{Name: "checkpoint", Purpose: "TODO"},
	}
}

func main() {
	for i, step := range steps() {
		fmt.Printf("%d. %s — %s\n", i+1, step.Name, step.Purpose)
	}
}
