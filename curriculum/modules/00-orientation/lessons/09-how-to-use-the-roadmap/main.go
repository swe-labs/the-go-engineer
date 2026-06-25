package main

import "fmt"

type Module struct {
	ID            string
	Title         string
	Prerequisites []string
}

type ValidationResult struct {
	Valid bool
	Step  int
	Path  []string
	Error string
}

func validatePath(path []Module) ValidationResult {
	completed := make(map[string]bool)
	for i, m := range path {
		for _, prereq := range m.Prerequisites {
			if !completed[prereq] {
				pathNames := make([]string, len(path))
				for j, p := range path {
					pathNames[j] = p.ID
				}
				return ValidationResult{
					Valid: false,
					Step:  i,
					Path:  pathNames,
					Error: fmt.Sprintf("module %q requires %q, but %q has not been completed", m.ID, prereq, prereq),
				}
			}
		}
		completed[m.ID] = true
	}
	pathNames := make([]string, len(path))
	for i, p := range path {
		pathNames[i] = p.ID
	}
	return ValidationResult{
		Valid: true,
		Step:  len(path) - 1,
		Path:  pathNames,
		Error: "",
	}
}

func main() {
	modules := map[string]Module{
		"00": {ID: "00", Title: "Orientation", Prerequisites: nil},
		"01": {ID: "01", Title: "Computers, Terminal, Git, Web", Prerequisites: []string{"00"}},
		"02": {ID: "02", Title: "Go Basics", Prerequisites: []string{"00", "01"}},
		"03": {ID: "03", Title: "Testing in Go", Prerequisites: []string{"00", "01", "02"}},
	}

	validPath := []Module{modules["00"], modules["01"], modules["02"], modules["03"]}
	invalidPath := []Module{modules["02"], modules["01"], modules["00"], modules["03"]}

	fmt.Println("Roadmap Path Validation")
	fmt.Println("=======================")
	fmt.Println()

	result1 := validatePath(validPath)
	fmt.Println("Path: 00 -> 01 -> 02 -> 03")
	fmt.Printf("Valid: %t\n", result1.Valid)
	fmt.Println()

	result2 := validatePath(invalidPath)
	fmt.Println("Path: 02 -> 01 -> 00 -> 03")
	fmt.Printf("Valid: %t\n", result2.Valid)
	if !result2.Valid {
		fmt.Printf("Error at step %d: %s\n", result2.Step, result2.Error)
	}
}
