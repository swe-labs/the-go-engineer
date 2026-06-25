package main

import "testing"

func TestEvaluateAnswer_PassesWhenAboveThreshold(t *testing.T) {
	rubric := Rubric{
		Categories: []RubricCategory{
			{Name: "Correctness", MaxScore: 4, PassThreshold: 3},
		},
	}
	answers := []Answer{{Category: "Correctness", Score: 3}}
	results := evaluateAnswer(rubric, answers)
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if !results[0].Passed {
		t.Errorf("Correctness should pass (score 3 >= threshold 3)")
	}
}

func TestEvaluateAnswer_FailsBelowThreshold(t *testing.T) {
	rubric := Rubric{
		Categories: []RubricCategory{
			{Name: "Correctness", MaxScore: 4, PassThreshold: 3},
		},
	}
	answers := []Answer{{Category: "Correctness", Score: 1}}
	results := evaluateAnswer(rubric, answers)
	if results[0].Passed {
		t.Errorf("Correctness should fail (score 1 < threshold 3)")
	}
}

func TestEvaluateAnswer_MultipleCategories(t *testing.T) {
	rubric := Rubric{
		Categories: []RubricCategory{
			{Name: "Correctness", MaxScore: 4, PassThreshold: 3},
			{Name: "Clarity", MaxScore: 3, PassThreshold: 2},
			{Name: "Completeness", MaxScore: 3, PassThreshold: 2},
		},
	}
	answers := []Answer{
		{Category: "Correctness", Score: 4},
		{Category: "Clarity", Score: 1},
		{Category: "Completeness", Score: 0},
	}
	results := evaluateAnswer(rubric, answers)
	tests := []struct {
		cat      string
		wantPass bool
	}{
		{"Correctness", true},
		{"Clarity", false},
		{"Completeness", false},
	}
	for _, tc := range tests {
		for _, r := range results {
			if r.Category == tc.cat && r.Passed != tc.wantPass {
				t.Errorf("%s passed = %v, want %v", tc.cat, r.Passed, tc.wantPass)
			}
		}
	}
}

func TestEvaluateAnswer_MissingCategoryGetsZero(t *testing.T) {
	rubric := Rubric{
		Categories: []RubricCategory{
			{Name: "Correctness", MaxScore: 4, PassThreshold: 3},
		},
	}
	results := evaluateAnswer(rubric, nil)
	if results[0].Score != 0 {
		t.Errorf("missing answer should get score 0, got %d", results[0].Score)
	}
	if results[0].Passed {
		t.Errorf("missing answer should not pass")
	}
}

func TestEvaluateAnswer_NilRubricCategories(t *testing.T) {
	rubric := Rubric{}
	results := evaluateAnswer(rubric, nil)
	if len(results) != 0 {
		t.Errorf("empty rubric should produce no results, got %d", len(results))
	}
}
