# How to ask good debugging questions

## Learning objective

Turn vague confusion into a structured question that contains evidence and a testable hypothesis. You will learn the five fields every debugging question should include and how to score your own questions for completeness before you ask for help.

## Why this matters

The most common reason engineers get stuck is not lack of skill. It is asking a question that cannot be answered. "It doesn't work" tells the helper nothing. "I expected X, I got Y, I ran Z, here is the error, here is my environment" gives them everything they need to help you in minutes. This skill separates junior engineers who stay stuck for hours from engineers who unblock themselves or get quick help. It also makes you a better debugger of your own problems: writing down the fields forces you to think clearly.

## Mental model

Think of a debugging question as a detective case file. The file has five slots: what you expected to happen, what actually happened, what command you ran, what error you saw, and what environment you were in. A complete case file lets another detective pick it up and start investigating immediately. A file with missing slots forces them to track you down and ask basic questions, wasting everyone's time.

## Core idea

Every debugging question can be scored on completeness from 1 to 5 based on five fields: Expected, Actual, Command, Error, Environment. A score of 5 means the question contains everything needed to reproduce and diagnose the problem. A score of 1 or 2 means the question will require back-and-forth before anyone can help.

## Under the hood

The five fields serve specific purposes:
- Expected: defines the correct behavior. Without this, the helper does not know what "working" looks like.
- Actual: defines the observed behavior. The gap between expected and actual is the bug.
- Command: defines exactly how to reproduce. "I ran `curl -v http://localhost:8080/health`" is reproducible. "I tried to start the server" is not.
- Error: defines the specific error signal. Paste the exact error message, not a paraphrase.
- Environment: defines the context. OS version, Go version, running locally or in Docker, etc.

When you fill all five, you have usually solved the problem yourself because writing forces clarity.

## How Go uses it

The Go program defines a `DebugQuestion` struct with string fields for each of the five categories. The `scoreQuestion` function counts how many fields are non-empty (after trimming whitespace) and returns a score from 0 to 5. The `main` function demonstrates two questions: a vague one (score 2) and a high-quality one (score 5), showing the learner what each level looks like.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

type DebugQuestion struct {
	Expected    string
	Actual      string
	Command     string
	Error       string
	Environment string
}

func scoreQuestion(q DebugQuestion) int {
	score := 0
	if strings.TrimSpace(q.Expected) != "" {
		score++
	}
	if strings.TrimSpace(q.Actual) != "" {
		score++
	}
	if strings.TrimSpace(q.Command) != "" {
		score++
	}
	if strings.TrimSpace(q.Error) != "" {
		score++
	}
	if strings.TrimSpace(q.Environment) != "" {
		score++
	}
	return score
}

func main() {
	vague := DebugQuestion{
		Expected:    "",
		Actual:      "it doesn't work",
		Command:     "",
		Error:       "some error",
		Environment: "",
	}

	good := DebugQuestion{
		Expected:    "the server should return 200 OK on GET /health",
		Actual:      "the server returns 503 Service Unavailable",
		Command:     "curl -v http://localhost:8080/health",
		Error:       "HTTP/1.1 503 Service Unavailable",
		Environment: "Windows 11, Go 1.25, running locally with go run",
	}

	fmt.Println("Vague question score:", scoreQuestion(vague), "/ 5")
	fmt.Println("Good question score: ", scoreQuestion(good), "/ 5")
}
```

## Step-by-step execution

1. The program creates two `DebugQuestion` values. The `vague` one fills only Actual ("it doesn't work") and Error ("some error"), leaving the other three fields empty.
2. The `good` one fills all five fields with specific, concrete information.
3. `scoreQuestion` is called for each. It checks each field with `strings.TrimSpace` to handle whitespace-only inputs.
4. For `vague`: only Actual and Error are non-empty, so the score is 2.
5. For `good`: all five fields are non-empty, so the score is 5.
6. `main` prints both scores side by side.

The lesson is visual: you can see exactly which fields the vague question is missing and how much more useful the good question is.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Describing the error from memory instead of pasting it | You think you remember it exactly | Always copy-paste the exact error. One character difference can matter. |
| Leaving out the expected behavior | You assume the helper knows what should happen | State the expected outcome explicitly. Do not make assumptions. |
| Writing commands from memory with typos | You abbreviate or approximate | Copy-paste the actual command you ran from your terminal history. |
| Omitting the environment | You think the problem is universal | Include OS, language version, and runtime context. Many bugs are environment-specific. |
| Writing a paragraph instead of structured fields | You dump everything you know in one block | Use the five-field structure explicitly. Your helper will thank you. |

## Debugging walkthrough

Imagine you run this program and it prints a score of 0 for a question where you filled in some fields.

1. Check for whitespace. A field containing `" "` (a single space) will be treated as empty because `strings.TrimSpace(" ")` returns `""`.
2. Verify the field names in your `DebugQuestion` literal match the struct exactly. A typo like `Environent` will compile but the field will stay empty.
3. Add a debug print of each field with its length before scoring: `fmt.Printf("Expected: %q (len=%d)\n", q.Expected, len(q.Expected))`.
4. If scoreQuestion is returning unexpected results, write a small test that passes a known value and check the return.
5. If the program does not compile, check that you imported `"strings"` and that all struct fields are correctly typed as `string`.

For debugging your own learning: if you are stuck and cannot formulate a good question, try writing down the five fields in a text file. The act of writing often reveals the answer.

## Production notes

Many engineering teams use debugging templates in their issue trackers. A bug report template with fields for Expected, Actual, Steps to Reproduce, Logs, and Environment is essentially the same idea as this five-field struct. Some teams automate the scoring: if a bug report is missing required fields, a bot labels it "needs more info."

When you are helping someone else debug, ask them for the five fields before reading their error. This trains them to provide structured information and reduces the time you spend extracting basic facts.

## Performance implications

The scoring function runs in O(1) time with five field checks. There is no performance concern. The meaningful performance metric is human: a well-structured question can reduce debugging time from hours to minutes. The cost of writing a structured question is about 30 seconds. The cost of a vague question is often several hours of back-and-forth.

## Practice task

Take a debugging question you have asked recently (or invent one). Write it as a `DebugQuestion` struct. Score it. Identify which fields are missing. Rewrite the question to include all five fields. Run both through the program and compare the scores.

## Tests / verification

```bash
go test ./curriculum/modules/00-orientation/lessons/07-how-to-ask-good-debugging-questions
```

The tests verify:
- All five fields filled scores 5
- No fields filled scores 0
- Partial fields score correctly (2 for Expected+Actual, 1 for only Error, 4 for all but Environment)
- Whitespace-only fields are treated as empty
- The `formatField` helper returns "(empty)" for blank input

## Review questions

1. What are the five fields of a structured debugging question?
2. Why does the function use `strings.TrimSpace` before checking if a field is empty?
3. A question scores 3/5. What information is likely missing?
4. How does writing a structured question help you debug even before asking for help?
5. What would you change in the struct to add a "Steps to Reproduce" field?

## NEXT UP

[What job-ready means](../08-what-job-ready-means/README.md) — Define job readiness as demonstrated ability across five skill categories.
