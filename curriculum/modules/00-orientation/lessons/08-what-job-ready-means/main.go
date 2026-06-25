package main

import "fmt"

type SkillCategory struct {
	Name     string
	Evidence []string
}

type JobReadiness struct {
	Categories []SkillCategory
}

func identifyGaps(jr JobReadiness) []string {
	var gaps []string
	for _, cat := range jr.Categories {
		if len(cat.Evidence) == 0 {
			gaps = append(gaps, cat.Name)
		}
	}
	return gaps
}

func readinessSummary(jr JobReadiness) string {
	total := len(jr.Categories)
	filled := 0
	for _, cat := range jr.Categories {
		if len(cat.Evidence) > 0 {
			filled++
		}
	}
	if total == 0 {
		return "No skill categories defined"
	}
	percentage := filled * 100 / total
	if percentage == 100 {
		return "All skill categories have evidence. You are job-ready!"
	}
	return fmt.Sprintf("%d/%d categories have evidence (%d%%). Keep building.", filled, total, percentage)
}

func main() {
	readiness := JobReadiness{
		Categories: []SkillCategory{
			{
				Name: "Implementation",
				Evidence: []string{
					"Built a CLI tool in Go",
					"Implemented REST API with tests",
				},
			},
			{
				Name:     "Testing",
				Evidence: []string{},
			},
			{
				Name: "Debugging",
				Evidence: []string{
					"Debugged a production race condition",
				},
			},
			{
				Name:     "Communication",
				Evidence: []string{},
			},
			{
				Name:     "Operations",
				Evidence: []string{},
			},
		},
	}

	fmt.Println("Job Readiness Report")
	fmt.Println("====================")
	for _, cat := range readiness.Categories {
		status := "has evidence"
		if len(cat.Evidence) == 0 {
			status = "GAP - no evidence"
		}
		fmt.Printf("%-16s %s\n", cat.Name+":", status)
	}

	fmt.Println()
	fmt.Println(readinessSummary(readiness))

	gaps := identifyGaps(readiness)
	if len(gaps) > 0 {
		fmt.Println("\nFocus on these gaps:")
		for _, g := range gaps {
			fmt.Printf("  - %s\n", g)
		}
	}
}
