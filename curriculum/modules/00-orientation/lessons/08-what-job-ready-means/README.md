# What job-ready means

## Learning objective

Define job readiness as demonstrable evidence across five skill categories: Implementation, Testing, Debugging, Communication, and Operations. You will be able to audit your own portfolio, identify gaps, and describe what evidence proves readiness in each category.

## Why this matters

"Job-ready" is a vague phrase that means different things to different people. Some think it means knowing a language syntax. Others think it means having a degree. Employers think it means something else: can you deliver working software, verify it works, fix it when it breaks, explain your decisions, and run it in production. If you cannot point to evidence for each of these, you are not job-ready yet. The good news is that this curriculum is designed to produce that evidence. You just need to know what to collect and where the gaps are.

## Mental model

Think of job readiness as a set of five buckets. Each bucket holds evidence: code you wrote, tests you passed, bugs you fixed, documents you wrote, deployments you did. A bucket with no evidence is a gap. Your goal is to fill every bucket before you apply for jobs. You do not need the same amount in every bucket, but every bucket must have at least one piece of evidence.

## Core idea

Job readiness is not a feeling. It is a portfolio of repeatable evidence across five dimensions. If you cannot show a project, a test, a debugging story, a clear README, and a deployment, you are not ready. If you can show all five, you are ready even if you feel nervous.

## Under the hood

The five categories map to real engineering activities:
- Implementation: writing code that compiles, runs, and solves a problem. Evidence includes working projects, pull requests, and code samples.
- Testing: verifying that your code works correctly. Evidence includes test files, test coverage reports, and CI pipelines.
- Debugging: finding and fixing problems. Evidence includes bug fix commits, debugging notes, and postmortems.
- Communication: explaining technical decisions clearly. Evidence includes README files, design documents, and code review comments.
- Operations: running software reliably. Evidence includes Dockerfiles, deployment scripts, monitoring dashboards, and incident reports.

A candidate who has evidence in all five categories is hireable. A candidate who has evidence only in implementation will struggle in interviews that ask about testing or debugging.

## How Go uses it

The Go program defines a `SkillCategory` struct with a `Name` and an `Evidence` string slice. A `JobReadiness` struct holds a slice of categories. The `identifyGaps` function scans all categories and returns the names of those with zero evidence. The `readinessSummary` function calculates a percentage of categories with evidence and returns a human-readable summary.

## Go example

```go
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
	return fmt.Sprintf("%d/%d categories have evidence (%d%%). Keep building.",
		filled, total, percentage)
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
			{Name: "Testing", Evidence: []string{}},
			{
				Name: "Debugging",
				Evidence: []string{
					"Debugged a production race condition",
				},
			},
			{Name: "Communication", Evidence: []string{}},
			{Name: "Operations", Evidence: []string{}},
		},
	}

	fmt.Println("Job Readiness Report")
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
		fmt.Println("Focus on these gaps:")
		for _, g := range gaps {
			fmt.Printf("  - %s\n", g)
		}
	}
}
```

## Step-by-step execution

1. `main` creates a `JobReadiness` with five categories. Implementation has two pieces of evidence, Debugging has one, and Testing, Communication, and Operations have none.
2. The first loop prints each category with its status: "has evidence" or "GAP - no evidence".
3. `identifyGaps` iterates each category, checks `len(Evidence)`, and collects names of empty ones.
4. `readinessSummary` counts filled vs total categories and computes a percentage.
5. The output shows 2/5 categories have evidence (40%) and lists the three gaps.

The learner can see that they have implementation and debugging stories but need to work on testing, communication, and operations.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Counting years of experience as evidence | Time spent does not equal skill | Evidence must be specific and verifiable: a link, a file, a description. |
| Treating coursework as job-ready evidence | Course grades show you passed a class, not that you can deliver software | Course projects count if they are complete, tested, and documented. But they need to be real, not toy examples. |
| Ignoring communication and operations | Most learners focus only on coding | A candidate who codes well but cannot write a README or deploy an app will not pass most interviews. |
| Thinking one category is enough | Being great at implementation but weak at testing is a red flag | Employers want balanced readiness. A gap in any category weakens your overall signal. |
| Not updating evidence as you progress | You build skills but forget to record them | Review your readiness every module and add new evidence as you complete projects. |

## Debugging walkthrough

Suppose you add evidence to a category but `identifyGaps` still reports it as a gap.

1. Check the `Name` field spelling. `"Testing"` vs `"testing"` are different strings.
2. Verify you are appending to the `Evidence` slice correctly: `Evidence: []string{"my evidence"}` not `Evidence: "my evidence"` (which would be a type error).
3. Print `len(cat.Evidence)` inside the loop to confirm the data is there.
4. If the summary percentage seems wrong, check the division: `filled * 100 / total` uses integer division. If filled=2 and total=5, this is 40, which is correct. If you use `filled / total` without multiplying by 100, you get 0.
5. If you are modifying the struct after creating it and the changes do not appear, remember Go passes by value. Modify the slice directly or use pointers.

If the program crashes with a nil pointer, check that you initialized the `Categories` slice: `Categories: []SkillCategory{...}` not just declaring the struct.

## Production notes

Engineering teams use readiness assessments in various forms: skill matrices for promotions, hiring rubrics for interviews, and individual development plans. The five-category model here is simplified. Real engineering ladders at companies like Google, Stripe, and GitHub have 10-15 categories with multiple levels each. But the principle is identical: define the categories, collect evidence, identify gaps, and work on them systematically.

When you prepare for job interviews, map each interview round to a category. A system design interview tests Communication and Operations. A coding interview tests Implementation. A debugging interview tests Debugging. A take-home project tests Implementation and Testing. Prepare evidence for each category before the interview, not after.

## Performance implications

The gap analysis runs in O(n) where n is the number of categories. For the five categories we use, this is trivial. For an engineering ladder with 15 categories, it is still trivial. The real performance consideration is human: how often do you audit yourself? Monthly is good. Weekly is better. Set a reminder.

## Practice task

Create your own `JobReadiness` with the five categories. Add at least one piece of evidence you already have to at least two categories. Run the program and note your gaps. Then write down one concrete action you can take in the next week to fill one gap. For example, if Testing is a gap, write a test for a project you have already built.

## Tests / verification

```bash
go test ./curriculum/modules/00-orientation/lessons/08-what-job-ready-means
```

The tests verify:
- No gaps when all categories have evidence
- All gaps when no categories have evidence
- Specific gap names are returned correctly
- Summary says "job-ready" when 100% filled
- Summary shows fraction when partially filled
- Empty readiness returns appropriate message

## Review questions

1. What are the five skill categories of job readiness?
2. Why is "I have 2 years of experience" not sufficient evidence?
3. What percentage of categories have evidence if 3 out of 5 are filled?
4. How would you add evidence to the Operations category?
5. Why does the program use integer division for the percentage calculation?

## NEXT UP

[How to use the roadmap](../09-how-to-use-the-roadmap/README.md) — Learn how module dependencies, prerequisites, and pacing guide your learning path.
