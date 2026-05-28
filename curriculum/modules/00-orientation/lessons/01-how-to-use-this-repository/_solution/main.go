package main

import (
	"fmt"
	"strings"
)

type folderPurpose struct {
	Name    string
	Purpose string
}

func purposes() []folderPurpose {
	return []folderPurpose{
		{Name: "metadata/", Purpose: "source of truth for graph, concepts, projects, assessments, contracts, and migration"},
		{Name: "curriculum/", Purpose: "learner-facing READMEs, code, tests, labs, projects, assessments, diagrams, and assets"},
		{Name: "tools/", Purpose: "validation, generation, audit, migration, and authoring automation"},
		{Name: "docs/", Purpose: "maintainer documentation, standards, governance, and release process"},
		{Name: "dist/", Purpose: "generated release artifacts only; never hand-edit"},
	}
}

func render(purposes []folderPurpose) string {
	var b strings.Builder
	for _, purpose := range purposes {
		fmt.Fprintf(&b, "%-12s %s\n", purpose.Name, purpose.Purpose)
	}
	return b.String()
}

func main() {
	fmt.Print(render(purposes()))
}
