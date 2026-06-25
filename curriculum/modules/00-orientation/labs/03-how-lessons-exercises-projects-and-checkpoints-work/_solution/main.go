package main

import (
	"fmt"
	"strings"
)

type cycleStep struct {
	Name    string
	Purpose string
}

func steps() []cycleStep {
	return []cycleStep{
		{Name: "read", Purpose: "build the mental model before touching code"},
		{Name: "run", Purpose: "observe the provided behavior"},
		{Name: "try", Purpose: "change something deliberately"},
		{Name: "test", Purpose: "prove behavior still matches expectations"},
		{Name: "reflect", Purpose: "explain what changed and why"},
		{Name: "checkpoint", Purpose: "verify the module-level skill"},
	}
}

func render(steps []cycleStep) string {
	var b strings.Builder
	for i, step := range steps {
		fmt.Fprintf(&b, "%d. %s — %s\n", i+1, step.Name, step.Purpose)
	}
	return b.String()
}

func main() {
	fmt.Print(render(steps()))
}
