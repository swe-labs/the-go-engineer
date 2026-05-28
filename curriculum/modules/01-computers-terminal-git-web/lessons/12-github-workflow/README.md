# Lesson: GitHub workflow

## Mission

Learn how local Git work connects to a remote repository and collaboration workflow.

By the end of this lesson, you should be able to explain the concept, run the example, and identify the mistake this lesson is designed to prevent.

## Prerequisites

You should have completed the previous lesson:

- core-01-11 (Branching and merging)

Understand branches as safe lines of work.

## Mental Model

Your local repo is your workshop. GitHub is a shared copy where teammates can see, discuss, and integrate work.

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

`git push` sends commits to a remote. `git fetch` downloads remote references. `git pull` fetches then integrates.

A professional engineer eventually learns to read the system from the machine's point of view. That does not mean memorizing internals immediately. It means asking what object exists, where it lives, who owns it, and what can observe it.

## Run Instructions

From the repository root:

```bash
go run ./curriculum/modules/01-computers-terminal-git-web/lessons/12-github-workflow
go test ./curriculum/modules/01-computers-terminal-git-web/lessons/12-github-workflow
```

You may also inspect or try these related commands:

- `git clone`
- `git remote -v`
- `git fetch`
- `git push`

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

Explain clone, remote, fetch, pull, and push without memorizing commands.

Then change one line in `main.go`, rerun the program, and explain what changed.

## Common Mistakes

| Mistake | Why it happens | Correction |
|---|---|---|
| Using `pull` and `push` as magic sync buttons. | The visible surface hides an important machine boundary. | Know which direction commits are moving and whether integration is happening. |
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

Most professional teams use hosted Git platforms for review, CI, issue tracking, release notes, and collaboration.

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

- explain `GitHub workflow` in your own words
- run the example
- run the test
- describe the common mistake
- explain the fix without reading this page

## Next Step

Continue to:

```text
core-01-13 — ../13-pull-requests-and-code-review/README.md
```
