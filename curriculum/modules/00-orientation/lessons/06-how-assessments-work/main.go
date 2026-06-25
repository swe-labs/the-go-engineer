package main

import "fmt"

type RubricCategory struct {
	Name          string
	MaxScore      int
	PassThreshold int
	Description   string
}

type Rubric struct {
	Categories []RubricCategory
}

type Answer struct {
	Category string
	Score    int
}

type Evaluation struct {
	Category string
	Score    int
	MaxScore int
	Passed   bool
	Feedback string
}

func evaluateAnswer(rubric Rubric, answers []Answer) []Evaluation {
	answerMap := make(map[string]int)
	for _, a := range answers {
		answerMap[a.Category] = a.Score
	}

	var results []Evaluation
	for _, cat := range rubric.Categories {
		score := answerMap[cat.Name]
		passed := score >= cat.PassThreshold
		feedback := "needs improvement"
		if passed {
			feedback = "meets standard"
		}
		results = append(results, Evaluation{
			Category: cat.Name,
			Score:    score,
			MaxScore: cat.MaxScore,
			Passed:   passed,
			Feedback: feedback,
		})
	}
	return results
}

func main() {
	rubric := Rubric{
		Categories: []RubricCategory{
			{Name: "Correctness", MaxScore: 4, PassThreshold: 3, Description: "Answer is factually correct"},
			{Name: "Clarity", MaxScore: 3, PassThreshold: 2, Description: "Answer is clearly explained"},
			{Name: "Completeness", MaxScore: 3, PassThreshold: 2, Description: "Answer covers all parts of the question"},
		},
	}

	answers := []Answer{
		{Category: "Correctness", Score: 4},
		{Category: "Clarity", Score: 1},
		{Category: "Completeness", Score: 0},
	}

	results := evaluateAnswer(rubric, answers)
	fmt.Println("Assessment Results:")
	fmt.Println("-------------------")
	for _, r := range results {
		status := "FAIL"
		if r.Passed {
			status = "PASS"
		}
		fmt.Printf("%-14s %d/%d  %s  (%s)\n", r.Category, r.Score, r.MaxScore, status, r.Feedback)
	}
}
