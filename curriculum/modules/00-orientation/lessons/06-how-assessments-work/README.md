# How assessments work

## Learning objective

Explain how checkpoints, rubrics, and answer keys measure mastery. You will be able to read a rubric, evaluate an answer against it, and describe how self-grading helps you identify what to study next.

## Why this matters

Without assessments, you cannot tell whether you learned something or just read it. Reading feels productive, but it does not produce evidence. Assessments close the gap between "I read the page" and "I can do the thing." In professional engineering, you will encounter code reviews, interview loops, on-call drills, and performance evaluations. All of them are assessments. Learning how they work now makes them less intimidating later.

## Mental model

Think of an assessment as a measuring tape. The tape does not make you taller. It tells you how tall you are. A rubric is the set of markings on the tape. Each category — correctness, clarity, completeness — is a different measurement. Your answer is the height. The pass threshold is the line you need to reach.

You are not being judged. You are being measured against a standard you can see ahead of time.

## Core idea

An assessment has three parts: a rubric that defines the criteria, an answer that provides the evidence, and an evaluation that compares the answer against the rubric. If you know the rubric before you start, you know exactly what to aim for.

## Under the hood

A rubric is a table. Each row is a category with a name, a maximum score, and a pass threshold. The evaluator looks at the answer, assigns a score per category, and checks whether the score meets or exceeds the threshold. Categories that fall short tell you where to focus. This is called criterion-referenced assessment: you are measured against fixed criteria, not against other people.

## How Go uses it

The Go program in this lesson models exactly that: a `Rubric` struct holds a list of `RubricCategory` values. Each category has a `Name`, `MaxScore`, `PassThreshold`, and `Description`. An `Answer` struct holds the score for a category. The `evaluateAnswer` function compares each answer against the rubric and returns an `Evaluation` with pass/fail status and feedback text.

## Go example

```go
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
			{Name: "Correctness", MaxScore: 4, PassThreshold: 3,
				Description: "Answer is factually correct"},
			{Name: "Clarity", MaxScore: 3, PassThreshold: 2,
				Description: "Answer is clearly explained"},
			{Name: "Completeness", MaxScore: 3, PassThreshold: 2,
				Description: "Answer covers all parts of the question"},
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
		fmt.Printf("%-14s %d/%d  %s  (%s)\n",
			r.Category, r.Score, r.MaxScore, status, r.Feedback)
	}
}
```

## Step-by-step execution

1. The `main` function creates a `Rubric` with three categories: Correctness (max 4, pass at 3), Clarity (max 3, pass at 2), Completeness (max 3, pass at 2).
2. It creates three `Answer` values: Correctness gets 4/4, Clarity gets 1/3, Completeness gets 0/3.
3. `evaluateAnswer` builds a map from category name to score so it can look up each category quickly.
4. For each rubric category, it looks up the score (defaults to 0 if missing) and checks `score >= PassThreshold`.
5. It builds an `Evaluation` with the pass/fail result and descriptive feedback.
6. `main` prints each category, the score fraction, PASS/FAIL, and feedback.

The output shows Correctness passes, Clarity fails, Completeness fails. The learner can see exactly which categories need work.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Confusing MaxScore with PassThreshold | They sound similar | Remember: MaxScore is the ceiling, PassThreshold is the passing line. They can be different numbers. |
| Forgetting to handle missing answers | If no answer exists for a category, the score defaults to 0 | Always initialize your answer map so missing entries get score 0. |
| Treating all categories as equally important | Rubrics often weight categories differently | Look at MaxScore and PassThreshold together. A category with max 8 and pass at 6 is more significant than one with max 2 and pass at 1. |
| Only counting total score | A total of 7/10 might hide that one category scored 0 | Check each category individually. A zero in a critical category matters more than a high total. |
| Skipping the feedback field | PASS/FAIL alone does not tell you what to improve | Always include a feedback message that explains the result. |

## Debugging walkthrough

Suppose you run this program and every category shows PASS even when you expect some to fail.

1. Check your `PassThreshold` values. If you set `PassThreshold: 0`, every answer will pass.
2. Check the `Score` values in your `Answer` slice. If you accidentally set all scores to 4, they will all pass.
3. Verify that the `Category` string in each answer matches exactly the `Name` in the rubric. "correctness" does not match "Correctness".
4. Print the `answerMap` after it is built to confirm it has the right keys and values.
5. Add a debug print of the threshold comparison: `fmt.Printf("%s: score=%d >= threshold=%d? %v\n", cat.Name, score, cat.PassThreshold, score >= cat.PassThreshold)`.

If instead you get no output, make sure `main` is defined and you are running `go run` on the correct directory. Run `go run .` from the lesson directory.

## Production notes

In real engineering teams, assessments take many forms: code review feedback forms, interview scorecards, performance review rubrics. The same principle applies: define clear criteria in advance, evaluate each criterion independently, and provide actionable feedback. Never change the rubric after seeing the answer. That invalidates the measurement.

When building automated assessment tools, store rubrics as data (JSON, YAML, or a database table) rather than hard-coding them in Go structs. This lets non-engineers update criteria without touching code.

## Performance implications

This evaluation runs in O(n + m) where n is the number of rubric categories and m is the number of answers. The map lookup per category is O(1). For any realistic assessment (fewer than 100 categories), performance is irrelevant. The important performance consideration is human: a clear rubric reduces the time a reviewer spends deciding whether an answer passes.

## Practice task

Create a new rubric with four categories: `Syntax`, `Logic`, `Efficiency`, `Style`. Assign each a MaxScore and PassThreshold. Write five answers with different score patterns. Run the evaluation and identify which categories pass and which fail. Then change one threshold and see how the results change.

## Tests / verification

```bash
go test ./curriculum/modules/00-orientation/lessons/06-how-assessments-work
```

The tests verify:
- Answers above threshold pass
- Answers below threshold fail
- Multiple categories are evaluated independently
- Missing answers get score 0 and fail
- Empty rubric produces no evaluations

## Review questions

1. What is the difference between MaxScore and PassThreshold in a rubric?
2. What happens to an answer category that has no matching rubric category?
3. Why does the program use a map to store answers before evaluating?
4. If a learner scores 2/3 on Clarity and the PassThreshold is 2, do they pass?
5. How would you modify the program to support weighted categories?

## NEXT UP

[How to ask good debugging questions](../07-how-to-ask-good-debugging-questions/README.md) — Learn to turn vague confusion into structured, answerable questions.
