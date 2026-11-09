# How to ask good debugging questions

> *"Everyone knows that debugging is twice as hard as writing a program in the first place. So if you're as clever as you can be when you write it, how will you ever debug it?"* — Brian Kernighan

## Learning objective

By the end of this lesson, you will be able to structure technical debugging inquiries using the 5 Evidence Fields, isolate root causes systematically, and turn vague confusion into high-signal problem reports.

## Why this matters

In professional engineering environments, asking "Why doesn't this work?" without evidence wastes team time, delays incident resolution, and signals a lack of engineering rigor. The best engineers communicate problems with exact evidence: reproduction steps, configuration states, and verbatim error outputs. Learning to formulate high-signal debugging questions often leads you directly to the solution before you even have to ask someone else.

## Mental model

Whenever you encounter a broken test, compilation failure, or unexpected behavior, construct a **5-Field Debugging Case File**:

```text
The 5 Evidence Fields of Structured Debugging
1. EXPECTED BEHAVIOR    ──> What should happen according to the specification?
2. ACTUAL BEHAVIOR      ──> What happened instead? (Symptoms, observed output)
3. EXACT COMMAND & PATH ──> From which directory did you run what exact command?
4. FULL ERROR OUTPUT    ──> Verbatim terminal text and error traces (never cropped)
5. ENVIRONMENT CONTEXT  ──> Operating system, tool versions, and recent file changes
```

When you document all five fields, you replace subjective guesswork with empirical evidence. In many cases, writing down the exact command and reading the full error output reveals the root cause immediately.

## Visual model

```mermaid
flowchart TD
    classDef frontDoor fill:#1a233a,stroke:#3b82f6,stroke-width:2px,color:#ffffff;
    classDef exitGate fill:#1e3a1e,stroke:#4caf50,stroke-width:2px,color:#ffffff;
    classDef standardNode fill:#22272e,stroke:#444c56,stroke-width:1px,color:#ffffff;
    classDef capstoneNode fill:#2d2215,stroke:#f59e0b,stroke-width:2px,color:#ffffff;

    subgraph LowSignal ["Low-Signal Anti-Pattern"]
        L1["Error Encountered"]:::frontDoor --> L2["Vague Message: 'It broke'"]:::standardNode
        L2 --> L3["Frustration & Guessing"]:::capstoneNode
    end

    subgraph HighSignal ["High-Signal Engineering Protocol"]
        H1["Error Encountered"]:::frontDoor --> H2["Collect 5 Evidence Fields<br/>(Expected, Actual, Command, Error, Env)"]:::standardNode
        H2 --> H3["Formulate Falsifiable Hypothesis"]:::standardNode
        H3 --> H4["Test One Variable at a Time"]:::standardNode
        H4 ==> H5["Root Cause Isolated & Resolved"]:::exitGate
    end
```

## Core idea

1. **Verbatim Text, Never Screenshots**: Always copy-paste the exact error text. Screenshots cannot be searched, copied, or grepped, and they often cut off essential line numbers or stack frames.
2. **One Variable at a Time**: When trying to fix an issue, change one thing, re-test, and observe the result. Changing multiple parameters simultaneously obscures which change resolved or worsened the problem.
3. **Trace Back to the Origin**: Read error messages starting from the first line of output, not just the last line. The first line usually reports the primary failure; subsequent lines are often cascading side effects.
4. **Reproducibility**: If you cannot reproduce an error consistently, you cannot be sure you have fixed it. Always verify that you can trigger the failure predictably before implementing a solution.

## Practice task

1. Create a structured debugging case file template in a text file on your local computer containing the 5 Evidence Fields:
   - Expected behavior
   - Actual observed behavior
   - Exact command and working directory
   - Full verbatim error output
   - Environment context (OS, toolchain versions, terminal shell)
2. Execute an intentional invalid command in your terminal (such as attempting to change into a non-existent directory: `cd /non/existent/path/xyz`).
3. Complete the 5 Evidence Fields in your case file using the real terminal output produced by this failure.
4. Formulate a single falsifiable hypothesis explaining why the error occurred, and state the exact verification command that proves or disproves your hypothesis.

## Common mistakes

| Mistake | Consequence | Correction |
|---|---|---|
| Paraphrasing error messages (for example: "I got a permissions error") | Omitting the exact file path or error code needed for diagnosis | Copy and paste the entire verbatim error output from the terminal |
| Changing code randomly in hopes that it starts working | Introducing new bugs while failing to understand the original issue | Read the documentation, state a hypothesis, and test one change at a time |
| Asking for help without stating what you have already tried | Forcing others to suggest things you already ruled out | Include a summary of your attempted fixes and their observed outcomes |

## Review questions

1. What are the 5 Evidence Fields required for a structured debugging report?
2. Why is copying verbatim error text superior to taking screenshots or paraphrasing errors?
3. Why should you change only one variable at a time when diagnosing an issue?
4. Why does reading the first line of an error trace matter more than reading the final line?

## NEXT UP

Advance to understanding what job readiness looks like in software engineering:

- [Lesson 08: What job ready means](../08-what-job-ready-means/README.md)
