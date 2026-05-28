# Lesson: How assessments work

## Mission

Learn how checkpoints measure mastery and how to use rubrics as feedback.

By the end, you should be able to explain the idea in your own words, point to the files involved, and describe what proof shows completion.

## Prerequisites

You need:

- the repository open in an editor
- a terminal
- basic ability to read Markdown
- willingness to run small commands and inspect output

No prior Go expertise is required for this orientation module.

## Mental Model

An assessment is a mirror. It does not make you better by itself; it shows where to practice next.

When something feels confusing, ask: "Which surface am I looking at?"

- metadata: the map
- curriculum: the learning experience
- tools: the checkers and generators
- docs: the maintainer rules
- dist: generated release output

## Visual Model

```text
metadata entry
    |
    | points to
    v
curriculum file or folder
    |
    | learner reads/runs/completes
    v
proof artifact
    |
    | validator/auditor checks
    v
completion signal
```

The important idea is direction: metadata describes the intended curriculum; learner-facing files provide the experience; proof and validation close the loop.

## Machine View

Assessment metadata points to questions, answer keys, rubrics, target lessons, and evidence requirements.

The filesystem is not just storage. In this curriculum, folder structure is part of the contract. A wrong path is not cosmetic; it can break validation, confuse learners, and hide incomplete work.

## Run Instructions

From the repository root, run:

```bash
go run ./curriculum/modules/00-orientation/lessons/06-how-assessments-work
go test ./curriculum/modules/00-orientation/lessons/06-how-assessments-work
```

Expected result:

```text
PASS
```

The program prints an orientation card. The tests verify that the lesson metadata and guidance are not empty or miswired.

## Code Walkthrough

Open `main.go`.

You will see a small program that models this orientation concept as structured data. The goal is not to learn advanced Go yet. The goal is to see that even orientation topics can have observable behavior and tests.

Important parts:

1. `lessonCard` stores the lesson ID, title, mission, and proof task.
2. `card()` creates the lesson's data.
3. `summary()` formats the data for terminal output.
4. `main()` prints the summary.

This pattern appears throughout the curriculum: small concept, clear representation, runnable proof.

## Try It

Read the module checkpoint rubric and identify what evidence would prove each category.

Then make a tiny change to `main.go`, run the program again, and explain what changed.

## Common Mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Reading without proving | Reading feels productive, but it does not show whether you can apply the idea. | Complete the practice task and run the test. |
| Skipping file paths | Paths look like boring details. | Treat paths as the map between metadata and learning. |
| Copying output without understanding | Terminal output can look like success even when the concept is unclear. | Explain the output in one paragraph. |

## Debugging Signals

If something fails:

- `no such file or directory` usually means you ran the command from the wrong root.
- `go: cannot find main module` means the repository root or Go module setup is missing.
- `undefined` means code references a name that does not exist in that package.
- test failures usually tell you which expected field is missing.

Debugging starts by copying the exact command and exact error.

## In Production

Professional engineering teams also use repository contracts. Build systems, CI, deployment pipelines, and code review tools all depend on stable paths and clear ownership.

A curriculum repository is smaller than a production system, but the habit is the same: make structure explicit so humans and tools can trust it.

## Performance Notes

This lesson has no meaningful runtime performance concern. The performance lesson here is cognitive: good structure reduces search time and prevents repeated confusion.

## Security Notes

Do not commit secrets such as `.env` files, tokens, SSH keys, or credentials. Repository structure should make it easy to keep secrets out and generated artifacts separate.

## Thinking Questions

1. What file or folder proves that this lesson exists?
2. What would break if metadata pointed to the wrong README path?
3. How would you explain this lesson to a learner who has never used GitHub?
4. What is one mistake you are likely to make in this module?
5. What proof would convince you that you are ready to move on?

## Proof of Understanding

You are done when you can:

- explain the lesson's mission without reading this page
- run the example command
- run the test command
- complete the practice task
- identify the next lesson or module

## Next Step

Continue to:

```text
core-00-07 — ../07-how-to-ask-good-debugging-questions/README.md
```
