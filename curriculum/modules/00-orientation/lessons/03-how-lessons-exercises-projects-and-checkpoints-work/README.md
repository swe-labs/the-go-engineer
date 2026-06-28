# Lab 03: How lessons, exercises, projects, and checkpoints work

## Learning objective

By the end of this lab, you will understand the four content types used in this curriculum — lesson, exercise, project, and checkpoint — and how they connect to form a complete learning cycle. You will demonstrate this by running a Go program that models the cycle as a series of steps, advancing through each one programmatically.

## Why this matters

A curriculum without a consistent structure is a pile of files. When every module follows the same pattern — read, run, try, test, reflect, checkpoint — you never have to guess what to do next. This lab teaches you that pattern so that from Module 01 onward, you focus on Go, not on navigation.

## Mental model

Think of learning like learning a martial art. A lesson is the instructor showing you a move. An exercise is practicing that move on a training dummy. A project is sparring with a partner. A checkpoint is the belt test where a judge verifies your skill.

Each phase has a purpose. Skipping exercises to get to the project is like trying to spar before you can throw a punch. The learning cycle ensures you build skill in the right order.

## Core idea

The curriculum defines four content types:

| Type | Purpose | Length | Verification |
|------|---------|--------|-------------|
| **Lesson** | Teach a concept with explanation and example code | 1 session | Tests pass, review questions answered |
| **Exercise** | Apply the concept in a focused, well-scoped task | 1-2 sessions | Modified starter passes tests |
| **Project** | Combine multiple concepts into a real-world task | 2-4 sessions | Rubric-based assessment |
| **Checkpoint** | Verify mastery before moving to the next module | 1 session | All tests pass, self-assessment complete |

These types appear inside each module in the same order: lessons first, then exercises, then a project, then a checkpoint.

## Under the hood

The Go program in this lab defines a `LearningCycle` struct that holds a title and a slice of `Step` values. Each `Step` has an index, name, action description, and expected artifact. The `NewLearningCycle` constructor initializes the six-step cycle: Read, Run, Try, Test, Reflect, Checkpoint.

The `Advance()` method removes the first step from the slice. When the slice is empty, `Complete` becomes true. This simulates progressing through the cycle.

In the real curriculum, progression is not automatic. You decide when to move from one step to the next by verifying your own understanding. The model here is a simplified version that shows the structure.

## How Go uses it

Go's own learning resources follow a similar progression. The Go tour (`go.dev/tour`) introduces concepts sequentially (lessons), then asks you to modify code in the browser (exercises). The standard library documentation provides examples (lessons), and the "codewalk" articles demonstrate complete programs (projects). The Go blog's "maps" tutorial is a lesson; the "writing web applications" article is a project. This curriculum formalizes that progression.

## Go example

```go
package main

import (
	"fmt"
)

type Step struct {
	Index    int
	Name     string
	Action   string
	Artifact string
}

type LearningCycle struct {
	Title    string
	Steps    []Step
	Complete bool
}

func NewLearningCycle(title string) LearningCycle {
	return LearningCycle{
		Title: title,
		Steps: []Step{
			{Index: 1, Name: "Read", Action: "Read the lesson README and example code", Artifact: "notes or highlights"},
			{Index: 2, Name: "Run", Action: "Run the Go program and observe output", Artifact: "terminal output"},
			{Index: 3, Name: "Try", Action: "Complete the practice task in _starter/", Artifact: "modified code"},
			{Index: 4, Name: "Test", Action: "Run go test to verify your solution", Artifact: "PASS result"},
			{Index: 5, Name: "Reflect", Action: "Answer review questions from the README", Artifact: "written answers"},
			{Index: 6, Name: "Checkpoint", Action: "Run the module checkpoint or assessment", Artifact: "completion proof"},
		},
		Complete: false,
	}
}

func (lc *LearningCycle) Advance() bool {
	if len(lc.Steps) == 0 {
		return false
	}
	lc.Steps = lc.Steps[1:]
	if len(lc.Steps) == 0 {
		lc.Complete = true
	}
	return true
}

func main() {
	cycle := NewLearningCycle("How lessons, exercises, projects, and checkpoints work")
	fmt.Println("Learning Cycle:", cycle.Title)
	for _, s := range cycle.Steps {
		fmt.Printf("  Step %d: %s\n", s.Index, s.Name)
		fmt.Printf("          %s\n", s.Action)
		fmt.Printf("          Artifact: %s\n", s.Artifact)
	}
	fmt.Println("Complete:", cycle.Complete)
}
```

## Step-by-step execution

1. `NewLearningCycle` creates a `LearningCycle` with six steps in order. Each step has a name that describes the phase, an action that tells you what to do, and an artifact that proves you did it.
2. The `main` function prints the title and every step. This is your reference for the learning cycle.
3. `Advance()` removes the front step. After six calls, `Steps` is empty and `Complete` is true.
4. The tests verify that `NewLearningCycle` creates exactly six steps with correct names, that `Advance` works correctly, and that `Complete` is false until all steps are exhausted.
5. In the real curriculum, you do not call `Advance()` — you physically move from README to code to test. But the mental model is the same: one step at a time.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Skipping the Reflect step | Learners rush to the next lesson without answering review questions. Reflection solidifies memory. | Treat review questions as mandatory. Write answers in a journal file. |
| Going to the solution before trying the starter | The `_solution/` folder exists, but looking at it before attempting the task steals the learning. | Follow the cycle: Try first, then check the solution only to verify. |
| Treating all four content types the same | A lesson is read-only. An exercise expects you to write code. A project combines multiple skills. A checkpoint is a gate. | Check the metadata to see which content type you are in. Each type has a different expectation. |
| Thinking Complete means "done forever" | The learning cycle is a loop, not a line. You will revisit concepts in later modules. | Use the checkpoint as a signal to move forward, not a claim of permanent mastery. |

## Debugging walkthrough

Scenario: You run `go run .` and see an empty steps list.

Step 1: Check `NewLearningCycle`. The function should return a `LearningCycle` with a non-nil `Steps` slice. Verify the function body.

Step 2: If `Steps` is empty, the slice literal might have been removed or the function might return a zero-value `LearningCycle`. Add `fmt.Printf("%+v", lc)` after the call to inspect.

Step 3: If `Advance()` returns `false` unexpectedly, check whether `len(lc.Steps) == 0` is true. The `Advance` method modifies the receiver — ensure you called it on a pointer (`*LearningCycle`) not a value. If you call `lc.Advance()` where `lc` is a value (not a pointer), Go passes a copy, and the original `Steps` slice is never modified.

Step 4: Check that `Complete` is set to `true` only after the last step is removed. If `Complete` is `true` before all steps are exhausted, look at the condition `if len(lc.Steps) == 0 { lc.Complete = true }`. It should be inside `Advance()` after removing the step.

## Production notes

- In a production learning platform, the learning cycle would be persisted to a database. Each user would have a progress record tracking which step they completed.
- The six-step cycle is prescriptive. In practice, engineers loop back: you might Test, find a bug, go back to Try, fix it, and Test again. Allow iteration.
- Content types (lesson, exercise, project, checkpoint) should be discoverable from the filesystem layout. Module 01's structure should mirror Module 02's so that navigation is predictable.

## Performance implications

- The `LearningCycle` struct stores a slice of `Step` values. Slices in Go are backed by arrays; removing from the front with `Steps[1:]` is O(1) because it just moves the slice header, not the underlying data.
- The cycle model here uses < 200 bytes of memory. Even with thousands of users, the data structure is tiny.
- If you were building a production system, you would use an integer `currentStep` index instead of slicing. This avoids any allocation at all.

## Practice task

Add a seventh step called "Share" between Reflect and Checkpoint with Index 6 (renumber Checkpoint to 7). The Share step should have Action "Explain the concept to a peer or write a summary" and Artifact "explanation or summary paragraph." Update the tests to expect 7 steps and verify the new step's name appears in the correct position.

## Tests / verification

1. Read the inline Go example and trace the learning cycle from lesson to checkpoint.
2. Copy the code into a local `.go` file and run `go run .` to see the full cycle in action.
3. Complete the practice task above — add a seventh "Share" step and confirm the cycle is correct by reading the output.

## Review questions

1. What are the four content types used in this curriculum, and what is the purpose of each?
2. How does the `Advance()` method in the Go example simulate progression through the learning cycle?
3. Why should you complete the Reflec step before moving to the next lesson?
4. What would happen if a module had two projects but no exercises?
5. In the Go example, what change would you make to allow a learner to go back to a previous step?

## NEXT UP

[Lesson 04: How to run code in this repository](../04-how-to-run-code-in-this-repository/README.md) — where you learn the difference between `go run`, `go test`, `go build`, and `go vet`.
