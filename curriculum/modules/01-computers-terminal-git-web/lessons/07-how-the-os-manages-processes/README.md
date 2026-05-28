# Lesson: How the OS manages processes

## Mission

Preview process lifecycle, scheduling, standard streams, and how the operating system supervises running programs.

By the end of this lesson, you should be able to explain the concept, run the example, and identify the mistake this lesson is designed to prevent.

## Prerequisites

You should have completed the previous lesson:

- core-01-06 (Exit codes)

Understand how programs signal success or failure.

## Mental Model

The operating system is a traffic controller. It lets many processes share CPU, memory, files, and network devices without crashing into each other.

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

A process has an ID, memory, open file descriptors, environment, arguments, and scheduling state. The OS switches between runnable processes.

A professional engineer eventually learns to read the system from the machine's point of view. That does not mean memorizing internals immediately. It means asking what object exists, where it lives, who owns it, and what can observe it.

## Run Instructions

From the repository root:

```bash
go run ./curriculum/modules/01-computers-terminal-git-web/lessons/07-how-the-os-manages-processes
go test ./curriculum/modules/01-computers-terminal-git-web/lessons/07-how-the-os-manages-processes
```

You may also inspect or try these related commands:

- `ps`
- `top`
- `kill`

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

Identify a process ID and explain why a long-running server is still just a process.

Then change one line in `main.go`, rerun the program, and explain what changed.

## Common Mistakes

| Mistake | Why it happens | Correction |
|---|---|---|
| Assuming a program can run forever without OS supervision or resource limits. | The visible surface hides an important machine boundary. | Learn to inspect processes, handle shutdown signals, and design cleanup paths. |
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

Debugging servers requires understanding process status, signals, ports, logs, and resource limits.

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

- explain `How the OS manages processes` in your own words
- run the example
- run the test
- describe the common mistake
- explain the fix without reading this page

## Next Step

Continue to:

```text
core-01-08 — ../08-memory-preview-stack-vs-heap/README.md
```
