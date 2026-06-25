package main

import "fmt"

type PortfolioProject struct {
	Name        string
	Description string
	Evidence    []string
	Status      string
}

type PortfolioPlan struct {
	Owner    string
	Projects []PortfolioProject
}

func scoreReadiness(projects []PortfolioProject) map[string]int {
	scores := make(map[string]int)
	for _, p := range projects {
		score := 0
		if p.Description != "" {
			score++
		}
		if len(p.Evidence) > 0 {
			score += len(p.Evidence)
			if score > 3 {
				score = 3
			}
		}
		if p.Status == "complete" {
			score += 2
		} else if p.Status == "in-progress" {
			score++
		}
		if score > 5 {
			score = 5
		}
		scores[p.Name] = score
	}
	return scores
}

func main() {
	plan := PortfolioPlan{
		Owner: "Learner",
		Projects: []PortfolioProject{
			{
				Name:        "Course Setup Verification",
				Description: "Verify Go, git, and editor installation",
				Evidence:    []string{"main.go", "main_test.go", "README.md"},
				Status:      "complete",
			},
			{
				Name:        "CLI Task Tracker",
				Description: "A command-line task manager in Go",
				Evidence:    []string{},
				Status:      "planned",
			},
			{
				Name:        "REST API Server",
				Description: "A REST API with tests and middleware",
				Evidence:    []string{"design doc"},
				Status:      "in-progress",
			},
		},
	}

	scores := scoreReadiness(plan.Projects)

	fmt.Println("Portfolio Readiness Report")
	fmt.Println("==========================")
	fmt.Printf("Owner: %s\n\n", plan.Owner)
	for _, p := range plan.Projects {
		fmt.Printf("Project: %s\n", p.Name)
		fmt.Printf("  Description: %s\n", p.Description)
		fmt.Printf("  Evidence:    %d items\n", len(p.Evidence))
		fmt.Printf("  Status:      %s\n", p.Status)
		fmt.Printf("  Readiness:   %d/5\n\n", scores[p.Name])
	}

	fmt.Println("Recommendations:")
	for _, p := range plan.Projects {
		if scores[p.Name] < 3 {
			fmt.Printf("  - %s: add evidence and complete implementation\n", p.Name)
		}
	}
}
