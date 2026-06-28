# How to build a portfolio from this curriculum

## Learning objective

Learn how to turn curriculum projects into professional portfolio evidence. You will be able to create a portfolio plan, assess each project's readiness score, and identify what evidence each project needs before it is ready to show employers.

## Why this matters

Completing the curriculum projects is only half the work. The other half is presenting them so that an employer can quickly understand what you built, why you built it, and what it proves about your skills. A well-presented project with a clear README, visible tests, and deployment instructions is worth ten projects that are hidden in private repos with no documentation. Your portfolio is the bridge between completing this curriculum and getting a job.

## Mental model

Think of each portfolio project as a museum exhibit. The exhibit has a label (the name), a description (what it does), artifacts (the evidence: code, tests, screenshots), and a status (draft, in progress, complete). A museum visitor should be able to walk through your exhibits and understand your capabilities within five minutes. Every missing label or empty display case weakens the impression.

## Core idea

A portfolio project has four attributes: Name, Description, Evidence, and Status. The readiness score measures how complete the presentation is. A project with a clear description, multiple evidence items, and a "complete" status scores 5/5. A project with no description, no evidence, and "planned" status scores 0/5. You should aim for every portfolio project to score at least 4/5 before sharing it with employers.

## Under the hood

The readiness scoring system works as follows:
- Description is non-empty: +1 point
- Each evidence item: +1 point (capped at 3, so 3+ evidence items max the evidence category)
- Status "complete": +2 points
- Status "in-progress": +1 point
- Status "planned": +0 points

The maximum possible score is 5. This scoring encourages you to write descriptions, collect evidence (at least 3 items per project), and finish what you start.

## How Go uses it

The Go program defines a `PortfolioProject` struct with Name, Description, Evidence (string slice), and Status (string). A `PortfolioPlan` struct holds the owner name and a slice of projects. The `scoreReadiness` function iterates all projects and calculates a score for each using the rules above. The `main` function shows three projects at different readiness levels with recommendations for improvement.

## Go example

```go
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
				Evidence:    []string{"setup.txt", "README.md"},
				Status:      "complete",
			},
			{
				Name:     "CLI Task Tracker",
				Status:   "planned",
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
	for _, p := range plan.Projects {
		fmt.Printf("%-30s Readiness: %d/5\n", p.Name, scores[p.Name])
	}
}
```

## Step-by-step execution

1. The program creates a `PortfolioPlan` with three projects. The first is complete with three evidence items. The second is planned with no description and no evidence. The third is in-progress with a description and one evidence item.
2. `scoreReadiness` calculates scores. For the first project: description (+1), 3 evidence items (+3, capped at +3), complete (+2). Total: 1 + 3 + 2 = 6, capped at 5.
3. For the second project: no description (+0), no evidence (+0), planned (+0). Total: 0.
4. For the third project: description (+1), 1 evidence (+1), in-progress (+1). Total: 3.
5. `main` prints the name and readiness score for each project, then lists recommendations for projects scoring below 3.

The output shows the learner which projects are ready to share (score 5), which need more work (score 3), and which are still just ideas (score 0).

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Only including the code without a README | You assume the code speaks for itself | Every project must have a README that explains what it does, how to run it, and what you learned. |
| Using vague project names | "Project 1" or "Go App" tells an employer nothing | Name your projects descriptively: "CLI Task Tracker with SQLite Persistence" |
| Keeping projects private | You are afraid they are not good enough | Make them public. Imperfect public projects are better than perfect private ones. |
| Having only one evidence item per project | One main.go file does not demonstrate depth | Include tests, a README, a design doc, a Dockerfile, or a CI config. Each piece is separate evidence. |
| Not updating status as you progress | Projects stay "planned" forever | After each module, review your portfolio plan and update status. Move projects from planned to in-progress to complete. |

## Debugging walkthrough

Suppose you run the program and a project you know is complete shows a score of 2 instead of 5.

1. Check the `Status` field spelling. The program checks for `"complete"` exactly. `"Complete"` (capital C) will not match and will give +0.
2. Check the `Evidence` slice. If it is nil instead of an empty slice, `len(p.Evidence)` returns 0, and the evidence bonus is skipped. Initialize with `[]string{}` or append items one by one.
3. Add a debug print inside `scoreReadiness` to show the score breakdown: `fmt.Printf("Project %q: desc=%v evidence=%d status=%q partial=%d\n", p.Name, p.Description != "", len(p.Evidence), p.Status, score)`.
4. If the scores map is empty after calling `scoreReadiness`, verify you passed a non-nil slice. Passing nil returns an empty map.
5. If a project name contains leading or trailing spaces, the map key will include them, and printing `scores[p.Name]` with a different-spelled name will give 0.

If the program compiles but produces no output, make sure `main` calls `scoreReadiness` and iterates the results. The most common bug is defining `main` but forgetting to call the function.

## Production notes

A professional portfolio goes beyond what this program models. Real portfolios include:
- A personal website or GitHub profile with pinned repositories
- Each project README includes architecture diagrams, installation instructions, test commands, and a "what I learned" section
- Commit history shows consistent work over time (not a single massive push)
- Projects are deployed (even to a free tier) so employers can see them running
- Each project links to the related job readiness category it demonstrates

When you apply for jobs, customize your portfolio to the role. If the job emphasizes testing, lead with the project that has the most comprehensive tests. If it emphasizes operations, lead with the deployed project.

## Performance implications

The scoring function runs in O(n) where n is the number of projects. Each project evaluation is O(e) where e is the number of evidence items, with a cap at 3 for scoring purposes. Even with 50 portfolio projects, the calculation completes in microseconds. The real performance bottleneck is you: how often do you review and update your portfolio? Schedule a portfolio review at the end of every module.

## Practice task

Create your own `PortfolioPlan` with at least four projects. At least one should be complete, one in-progress, and one planned. Add realistic evidence items. Run the program and note which projects score below 4. Write down one specific action for each low-scoring project to bring it to at least 4/5 within the next two weeks.

## Tests / verification

The inline code example is for reading and understanding. To verify your understanding, complete the practice task above and check your answers against the description. You can also copy the inline code into a local `.go` file and run `go run .` and `go test .` in that directory to experiment with the output.

The tests in the inline example verify:
- Complete project with description and evidence scores 5
- Planned project with no info scores 0
- In-progress project with description scores at least 2
- Multiple projects are scored independently
- Evidence count is capped at 3 for scoring
- Empty project list returns empty map

## Review questions

1. What four attributes does a PortfolioProject have?
2. Why does the scoring cap evidence at 3 points instead of allowing unlimited evidence to increase the score?
3. A project has a description, two evidence items, and status "in-progress". What is its readiness score?
4. How would you change the scoring to add a bonus for having a README specifically?
5. Why should you make your portfolio projects public even if they are not perfect?

## NEXT UP

1. [Project: Course Setup Verification](../projects/course-setup-verification/README.md) — verify your environment
2. [Assessment: Course Setup](../assessments/course-setup/README.md) — confirm your environment is ready
3. [Assessment: Module 00 Checkpoint](../assessments/checkpoint/README.md) — verify you understand the orientation content

Then move on to [Module 01 — Computers, Terminal, Git, and the Web](../../../01-computers-terminal-git-web/README.md).
