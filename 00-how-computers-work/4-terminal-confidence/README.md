# HC.4 Terminal Confidence

## Mission

Build confidence with the terminal as the text-based environment that launches programs, passes arguments, and handles output streams.

## Prerequisites

- `HC.3` memory basics

## Mental Model

The terminal is not “where scary commands live.”
It is a process-launching and output-reading interface for the operating system.

## Visual Model

```mermaid
graph LR
    A["You type a command"] --> B["Shell parses it"]
    B --> C["Shell finds program in PATH"]
    C --> D["OS launches the program"]
    D --> E["Program writes stdout/stderr"]
    E --> F["Terminal displays or redirects output"]
```

## Machine View

When you press Enter in the shell:

1. the shell parses the command
2. it finds the program to run
3. the OS starts a new process
4. that process writes output to file descriptors like stdout and stderr
5. the shell receives an exit code when the process ends

That is why redirecting output and chaining commands works.

## Run Instructions

```bash
go run ./00-how-computers-work/4-terminal-confidence
```

## Code Walkthrough

In `main.go`, the lesson writes one line to stdout and one line to stderr.
That small demo makes the terminal's two most common output channels visible.

## Try It

1. Run the lesson normally and read both lines.
2. Run `go run ./00-how-computers-work/4-terminal-confidence > output.txt` and notice which line stays in the terminal.
3. Run a failing command and inspect its exit code in your shell.

## ⚠️ In Production

When production systems fail, you often have a shell, logs, and process output before you have anything else.
Terminal confidence becomes operational confidence.

## 🤔 Thinking Questions

1. Why does redirecting stdout not automatically redirect stderr?
2. Why do shells care about exit codes instead of reading the English words in program output?
3. What would break if the shell could not resolve programs through `PATH`?

## Next Step

Continue to [HC.5 How the OS Manages Processes](../5-os-processes).
