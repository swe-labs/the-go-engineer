package main

import "fmt"

type folderPurpose struct {
	Name    string
	Purpose string
}

func purposes() []folderPurpose {
	return []folderPurpose{
		{Name: "metadata/", Purpose: "TODO: explain the source of truth"},
		{Name: "curriculum/", Purpose: "TODO: explain learner-facing content"},
		{Name: "tools/", Purpose: "TODO: explain validation and automation"},
		{Name: "docs/", Purpose: "TODO: explain maintainer documentation"},
		{Name: "dist/", Purpose: "TODO: explain generated artifacts"},
	}
}

func main() {
	for _, purpose := range purposes() {
		fmt.Printf("%-12s %s\n", purpose.Name, purpose.Purpose)
	}
}
