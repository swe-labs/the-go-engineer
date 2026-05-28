# Lesson: Branching and merging

## Mission

Understand branches as safe lines of work and merges as ways to combine histories.

By the end of this lesson, you should be able to explain the concept, run the example, and identify the mistake this lesson is designed to prevent.

## Prerequisites

You should have completed the previous lesson:

- core-01-10 (Git basics: status, add, commit)

Be comfortable with the basic Git workflow loop.

## Mental Model

A branch is a movable label. Merging joins two lines of work into one shared history.

The mental model is deliberately simple. Later modules will add details, but this version is enough to stop common beginner confusion.

## Visual Model

```text
learner action
    ↓
machine boundary
    ↓
observable result
    ↓
proof or debugging signal
```

For this lesson, focus on the boundary: what does the learner ask the machine to do, and what does the machine actually receive or produce?

## Machine View

Git finds a common ancestor, compares changes on both sides, and creates a merge commit or fast-forward when possible.

A professional engineer eventually learns to read the system from the machine's point of view. That does not mean memorizing internals immediately. It means asking what object exists, where it lives, who owns it, and what can observe it.

## Run Instructions

From the repository root:

```bash
go run ./curriculum/modules/01-computers-terminal-git-web/lessons/11-branching-and-merging
go test ./curriculum/modules/01-computers-terminal-git-web/lessons/11-branching-and-merging
```

You may also inspect or try these related commands:

- `git switch -c`
- `git merge`
- `git log --oneline --graph`

## Code Walkthrough

Open `main.go`.

The program models the lesson as a small `conceptCard`:

1. `ID` ties the code back to metadata.
2. `Title` names the concept.
3. `MentalModel` gives the human model.
4. `MachineView` gives the operational model.
5. `CommonMistake` names the trap.
6. `Fix` gives the correction.
7. `Commands` lists commands worth recognizing.

The code is intentionally small because this module is about foundations, not language complexity.

## Try It

Draw a branch before and after a merge.

Then change one line in `main.go`, rerun the program, and explain what changed.

## Common Mistakes

| Mistake | Why it happens | Correction |
|---|---|---|
| Panicking when a merge conflict appears. | The visible surface hides an important machine boundary. | A conflict is a request for a human decision, not data loss by default. |
| Copying a command without knowing its working directory | The shell accepts the command, but paths resolve somewhere else. | Run `pwd` and explain the current directory first. |
| Treating a word as magic vocabulary | Terms like process, branch, port, or request get memorized but not understood. | Define the term using a concrete example. |

## Debugging Signals

Watch for these signals:

- the command fails only in one directory
- output mentions a missing file or unknown path
- a process keeps running after you expected it to stop
- Git says files are modified but you do not know why
- a network command fails before reaching application code

The first debugging move is to write down the exact command, exact output, and current directory.

## In Production

Teams use branches to isolate work, review safely, and keep mainline stable.

The professional habit is to connect vocabulary to operational evidence. If you cannot observe it, test it, log it, inspect it, or explain where it lives, the concept is still too magical.

## Performance Notes

This lesson is not about optimizing code. Its performance value is avoiding wasted time. Clear mental models reduce repeated debugging loops.

## Security Notes

Do not paste secrets into terminal commands, screenshots, issue descriptions, or pull requests. Environment variables, Git history, logs, and shell history can all leak sensitive values.

## Thinking Questions

1. What object is this lesson really about?
2. Where does that object live: source file, filesystem, process, Git history, network, or browser?
3. What command or output would prove you understand it?
4. What mistake would show confusion?
5. How will this idea appear again in later Go work?

## Proof of Understanding

You are complete when you can:

- explain `Branching and merging` in your own words
- run the example
- run the test
- describe the common mistake
- explain the fix without reading this page

## Next Step

Continue to:

```text
core-01-12 — ../12-github-workflow/README.md
```
