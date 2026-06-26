# Lesson 05: How starter folders work

## Learning objective

By the end of this lesson, you will understand how `_starter/` and `_solution/` directories work together in this curriculum. You will know why they exist, how to use them without cheating yourself, and how to programmatically compare a starter against a solution. You will demonstrate this by running a Go program that detects missing fields in incomplete work.

## Why this matters

Practice is how you learn to code. But practice only works if you struggle first. The starter folder gives you a scaffold — enough to start, but not enough to finish without thinking. The solution folder exists so you can verify your work, not so you can skip the work. This lesson teaches you the workflow so that every exercise in this curriculum maximizes your learning.

## Mental model

Imagine you are assembling furniture. The starter folder is the instruction sheet with some steps missing. You have the tools, you have the parts, but you need to figure out steps 4-7 yourself. The solution folder is a video of someone who already assembled it — you watch it only after you try, to see if you did it right.

If you watch the video first, you learn nothing. If you never watch the video, you never confirm your approach. The right sequence is: try from the starter, then check against the solution.

## Core idea

Every lab in this curriculum that requires writing code has two additional directories:

| Directory | Purpose | What is inside |
|-----------|---------|---------------|
| `_starter/` | A scaffold for the learner | Incomplete `main.go`, a README with the task, sometimes incomplete tests |
| `_solution/` | A reference implementation | Complete `main.go`, passing `main_test.go`, a README with the walkthrough |

The underscore prefix (`_starter/` not `starter/`) is intentional. In Go's toolchain, directories starting with `_` or `.` are ignored by `go build` and `go test`. This means:

- The starter will not compile on its own (if it imports from the parent package).
- The solution will not accidentally be compiled as part of the lesson's package.
- The learner must copy or work inside the starter directory to make progress.

To work on a lab, you copy `_starter/` to a new directory (or edit files in place if the lesson instructions say so), complete the missing code, and compare your result against `_solution/`.

## Under the hood

The Go example below defines a `Task` struct with three fields: `Name`, `Desc`, and `Complete`. The `CompareTasks` function takes two tasks — a starter and a solution — and returns a `Comparison` struct that reports which fields differ.

The `Comparison.Report()` method produces a human-readable string. If the starter's `Name` or `Desc` differs from the solution, the report says what is missing. If the starter matches the solution, the report says the starter is complete.

This is a simplified model of what the curriculum validator does. The real validator reads files from `_starter/` and `_solution/`, compares them structurally, and reports gaps.

## How Go uses it

The Go standard library uses a similar pattern in its own tests. The `go/types` package has test data directories with "golden" files (expected outputs) that tests compare against. The `go/importer` package has test fixtures. The `encoding/json` package has testdata directories with JSON files for encoding and decoding tests.

In each case, the test compares actual output against expected output. The starter/solution pattern is the same idea applied to learning: the starter is your "actual," the solution is the "expected," and you are the test runner checking whether they match.

## Go example

```go
package main

import (
	"fmt"
)

type Task struct {
	Name     string
	Desc     string
	Complete bool
}

type Comparison struct {
	MissingName bool
	MissingDesc bool
	IsComplete  bool
}

func CompareTasks(starter, solution Task) Comparison {
	return Comparison{
		MissingName: starter.Name != solution.Name,
		MissingDesc: starter.Desc != solution.Desc,
		IsComplete:  starter.Complete,
	}
}

func (c Comparison) Report() string {
	if c.MissingName || c.MissingDesc {
		return fmt.Sprintf("Starter is incomplete. Missing name: %v, missing desc: %v", c.MissingName, c.MissingDesc)
	}
	if c.IsComplete {
		return "Starter is complete and matches the solution."
	}
	return "Starter has all fields but is marked incomplete."
}

func main() {
	starter := Task{Name: "", Desc: "", Complete: false}
	solution := Task{Name: "Greet the user", Desc: "Write a function that prints a greeting", Complete: true}

	result := CompareTasks(starter, solution)
	fmt.Println("Starter:", starter)
	fmt.Println("Solution:", solution)
	fmt.Println("Comparison:", result.Report())

	fmt.Println()

	completed := Task{Name: "Greet the user", Desc: "Write a function that prints a greeting", Complete: true}
	result2 := CompareTasks(completed, solution)
	fmt.Println("Completed starter:", completed)
	fmt.Println("Comparison:", result2.Report())
}
```

## Step-by-step execution

1. The program defines `Task` with fields `Name`, `Desc`, and `Complete`. A starter task is typically incomplete — missing name, missing description, and not marked complete.
2. `CompareTasks` compares a starter against a solution field by field. It returns a `Comparison` struct with boolean flags for each possible difference.
3. The first comparison: starter has empty fields, solution has values. `CompareTasks` returns `MissingName: true, MissingDesc: true, IsComplete: false`. The report says "Starter is incomplete."
4. The second comparison: the starter has been edited to match the solution. `CompareTasks` returns all false for missing fields and `IsComplete: true`. The report says "Starter is complete and matches the solution."
5. The tests confirm all four field-combination scenarios: both missing, name only, desc only, and complete match.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Looking at the solution before trying the starter | The `_solution/` folder is visible in the same directory. It is tempting to peek when stuck. | Discipline: write down what you think the answer should be before checking. If you must peek, set a 15-minute timer first. |
| Editing the solution instead of the starter | Both folders look similar. You might accidentally modify `_solution/` and then wonder why the starter is still broken. | Always check the directory name before editing. The starter is the one you copy; the solution is read-only. |
| Copying the entire solution into the starter | If you copy the solution file over the starter, you learn nothing. The tests pass but your brain does not. | Compare the diff: `diff _starter/main.go _solution/main.go` to see what you needed to add. |
| Forgetting that underscore-directories are ignored by Go | If you run `go run ./curriculum/...` from the root, `_starter/` and `_solution/` are not compiled. If your code imports from them, the import will fail. | Run commands from inside the specific lesson directory, or use full relative paths that exclude the underscore directory. |

## Debugging walkthrough

Scenario: You copied `_starter/` to a work directory, completed the task, but `go test .` fails.

Step 1: Run `go test -v .` to see which specific test failed and what the expected vs actual values were. The `-v` flag prints the full test output including any messages from `t.Errorf`.

Step 2: In a lab that provides a solution, open `_solution/main_test.go` and read the test cases. The tests define the expected behavior. Your starter code must match these expectations.

Step 3: Compare your solution against the provided solution using a diff tool. On the command line: `diff -u your_work/main.go _solution/main.go`. Each line starting with `-` is in your file but not in the solution. Each line starting with `+` is in the solution but not in your file.

Step 4: If the diff is large, focus on the function signatures and return types first. A mismatch in the function signature (e.g., `func Add(a int, b int) int` vs `func Add(a, b int) int`) will cause a compilation error, not just a test failure.

Step 5: If the diff is small (one or two lines), your logic is likely correct but there is a subtle difference — an off-by-one or a missing edge case. Read the failing test case again and trace through your code with those specific inputs.

Scenario: `go build ./_starter/` says "no non-test Go files."

Step 1: This is expected. Directories starting with `_` are ignored by the Go toolchain. The starter is not meant to be built directly from its underscore directory.

Step 2: Copy the starter to a directory without the underscore prefix, or work inside the lesson directory and use `go run .` there.

Step 3: If you want to verify the starter compiles on its own, rename the directory temporarily: copy `_starter/` to `starter/` (no underscore) and build from there.

## Production notes

- In production code reviews, the reviewer's job is similar to the solution comparison: they compare your code against an expected standard (style guide, best practices, architecture). The starter/solution workflow teaches you to self-review before submitting.
- Underscore-prefixed directories in Go are a convention, not a security measure. A determined learner can still read the solution. The workflow relies on trust and good habits.
- Some teams use a "golden file" pattern in tests: the expected output is stored in a `testdata/` directory (also ignored by Go) and tests compare generated output against goldens. This is the same idea as starter/solution.

## Performance implications

- The `CompareTasks` function is O(1) — it compares three fields, no loops, no allocations. Performance is irrelevant for this lesson.
- The `Report()` method uses `fmt.Sprintf` which allocates a string. For a one-time report, this is fine. In a hot loop (thousands of reports per second), you would use a `strings.Builder`.
- The `Task` struct uses strings. In Go, string comparison is O(n) where n is the length of the shorter string. For the tiny strings in this curriculum (< 200 chars each), this is instantaneous.

## Practice task

Add a new field `Difficulty string` to the `Task` struct. Update `CompareTasks` to detect differences in `Difficulty`. Add a `MissingDifficulty` field to `Comparison`. Update `Report()` to mention difficulty: "Missing name, desc, or difficulty." Add table-driven test cases for all three combinations (difficulty missing, difficulty matches, difficulty differs).

## Tests / verification

The inline code example is for reading and understanding. To verify your understanding, complete the practice task above and check your answers against the description. You can also copy the inline code into a local `.go` file and run `go run .` and `go test .` in that directory to experiment with the output.

## Review questions

1. Why does the starter folder use an underscore prefix (`_starter/` instead of `starter/`)?
2. What does the `CompareTasks` function return when the starter's name field is empty but the solution's is not?
3. Why is it important to try the starter before looking at the solution?
4. How would you modify `CompareTasks` to return a list of missing field names instead of booleans?
5. What is the purpose of the `IsComplete` field in the `Comparison` struct?

## NEXT UP

[Lesson 06: How assessments work](../06-how-assessments-work/README.md) — where you learn how checkpoints and assessments verify your mastery before you move to the next module.
