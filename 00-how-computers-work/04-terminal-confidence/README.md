# HC.4 Terminal Confidence

## Mission

Build confidence with the terminal as the text-based environment that launches programs, passes arguments, and handles output streams.

## Prerequisites

- `HC.3` memory basics

## Mental Model

The terminal is not "where scary commands live." It is a professional interface for the operating system.

Think of it as a **Conversation with the OS**. You send a request (a command), and the OS provides a response (output).

## Visual Model

```mermaid
graph LR
    User["You type a command"] --> Shell["Shell parses tokens"]
    Shell --> PATH["Finds binary in PATH"]
    PATH --> Process["OS launches Process"]
    Process --> stdout["Standard Output (Default)"]
    Process --> stderr["Standard Error (Issues)"]
    stdout --> Screen["Terminal Screen"]
    stderr --> Screen
```

## Machine View

When you press Enter in the shell:
1. The shell parses the command into a program name and arguments.
2. It looks up the program's location using the `PATH` environment variable.
3. The OS creates a new process for that program.
4. The process is connected to three default **File Descriptors** (indices in the process's open file table):
   - **0 (stdin)**: Input from the keyboard.
   - **1 (stdout)**: Normal output.
   - **2 (stderr)**: Error messages.

> [!NOTE]
> The "Everything is a file" philosophy in Unix-like systems makes these descriptors incredibly powerful for I/O, as you will see in [FS.1 File Basics](../../05-packages-io/02-io-and-cli/03-filesystem/01-files/README.md).

## Run Instructions

```bash
go run ./00-how-computers-work/04-terminal-confidence
```

## Code Walkthrough

- **stdout**: Using `fmt.Println` writes to the standard output by default.
- **stderr**: Using `fmt.Fprintln(os.Stderr, ...)` writes to the standard error stream. This separation allows you to capture logs in a file while still seeing errors on the screen.

## Try It

1. Run the lesson normally and read both lines on your screen.
2. **Redirect stdout (Descriptor 1):** Run `go run ./00-how-computers-work/04-terminal-confidence > output.txt`. The `>` symbol is shorthand for `1>`. Only the error message stays on your screen because standard output is now "piped" into the file.
3. **Redirect stderr (Descriptor 2):** Run `go run ./00-how-computers-work/04-terminal-confidence 2> errors.txt`. Now standard output prints to the screen, but the error message is trapped in `errors.txt`.
4. **Capture Both:** Run `go run ./00-how-computers-work/04-terminal-confidence > all.log 2>&1`. This tells the shell to redirect stdout to a file, and then "copy" stderr (2) to the same place as stdout (1).

## What the Shell is Doing

When you use redirection (`>` or `2>`), the shell performs "plumbing" **before** your program even starts:
1. It sees the redirection token in your command.
2. It opens the target file on your disk.
3. It uses a system call (like `dup2`) to swap the default terminal connection of File Descriptor 1 (or 2) with the newly opened file.
4. Only then does it execute (`exec`) your program.

Because this happens at the OS level, your Go code doesn't even know it's writing to a file—it just keeps writing to "stdout" as usual!

## In Production

Every modern backend is managed through the terminal or terminal-like APIs. Knowing how to filter logs (`grep`), find processes (`ps`), and check exit codes (`$?`) is essential for any backend engineer.

## Thinking Questions

1. Why does redirecting stdout not automatically redirect stderr?
2. Why do shells care about numeric exit codes instead of reading the words in the program output?
3. What would break if the shell could not find programs through the `PATH`?

## Next Step

Next: `HC.5` -> [`00-how-computers-work/05-os-processes`](../05-os-processes/README.md)
