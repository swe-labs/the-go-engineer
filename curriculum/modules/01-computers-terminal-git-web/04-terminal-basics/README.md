# Terminal basics

## Learning objective

Understand and apply Terminal basics in the context of professional Go software engineering.

## Why this matters

New developers face the terminal with anxiety because every error looks like random noise. Understanding what the shell does — parsing, resolving, forking, waiting — transforms the terminal from a scary black box into a predictable tool.

## Mental model

A development environment pairs an editor with Go support, a terminal, Go tooling, and version control. The integrated terminal runs go build, go test, go fmt, and git without switching contexts.

## Core idea

The terminal predates GUIs and remains the most powerful interface for software development because it is composable, scriptable, remote-friendly, and consistent across machines. Learning the terminal is learning to speak the computer's native language.

## Under the hood

Every terminal command creates at minimum one OS process. The shell uses fork+exec on Unix or CreateProcess on Windows. stdin/stdout/stderr are file descriptors inherited from the parent shell. The shell's job is to orchestrate these processes — piping fd 1 of process A to fd 0 of process B, collecting exit codes, and managing job control (foreground/background).

## How Go uses it

Go developers use the terminal constantly: `go run` to test code, `go test` to verify correctness, `go build` to produce binaries, `go mod tidy` to manage dependencies, and `go vet` to catch bugs. The `go` command is itself a terminal program with flags, arguments, and exit codes.

## Go example

The example implements a minimal shell simulator. It reads input from stdin via `bufio.Scanner`, tokenizes the input by splitting on whitespace, and dispatches commands. It supports `echo` (prints arguments), `hello` (prints a greeting), and `exit` (terminates the process). Unknown commands produce a "command not found" error — mirroring how real shells resolve commands via PATH lookup.

## Step-by-step execution

1. The user opens a terminal emulator which starts a shell process (bash, zsh, or PowerShell).
2. The shell prints a prompt and waits for input — the user types a command and presses Enter.
3. The shell parses the input, resolves the command (built-in, alias, or PATH lookup), forks a child process, and waits for it to complete.
4. In the Go simulation: `bufio.NewScanner` reads lines from stdin, `strings.Fields` tokenizes the input, and `simulateShell` dispatches the command.
5. The child process runs, reading from stdin and writing to stdout/stderr until it exits.
6. The shell collects the exit code, displays the next prompt, and waits for the next command.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Typing commands without understanding the working directory | Running `rm -rf *` in the wrong directory can be catastrophic | Always check `pwd` before running destructive commands |
| Confusing the terminal with the shell | The terminal is the window; the shell is the program running inside it | The terminal emulator hosts the shell; closing the terminal kills the shell |
| Using relative paths without checking the current directory | Leads to "file not found" errors | Use `pwd` to confirm the working directory or use absolute paths |

## Debugging walkthrough

**Scenario: A beginner runs a command that modifies files but gets "permission denied".**
- **Cause:** The user does not have write permission on the target directory or file.
- **Diagnosis:** Run `ls -la` on the target to see ownership and permission bits.
- **Resolution:** Use `sudo` (on Unix) or run as Administrator (on Windows), or change permissions with `chmod`.

**Scenario: A learner runs a command with spaces in a file path and gets unexpected behavior.**
- **Cause:** The shell splits arguments on spaces; the path is broken into multiple arguments.
- **Diagnosis:** Count the arguments the command actually receives (use `echo $#` in bash).
- **Resolution:** Quote the path: `"My Documents/file.txt"` or escape spaces: `My\ Documents/file.txt`.

## Production notes

Professional developers spend hours daily in the terminal running builds, tailing logs, grepping code, managing git, deploying services, and debugging production. Terminal fluency is the difference between fighting the machine and commanding it.

## Performance implications

- Starting a new shell process for every command adds ~1-5ms overhead — negligible for interactive use but significant in tight loops.
- Piping large datasets between commands avoids writing temp files to disk, reducing I/O by 10-100x.
- Using shell built-ins (`cd`, `echo`, `source`) is faster than external commands since no process fork is needed.

## Practice task

Run three terminal commands, redirect output to a file using `>`, pipe data between two programs with `|`, and explain which process owns each step.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/04-terminal-basics/
```

## Review questions

1. Navigate to a specific directory using only `cd` and `ls`, then confirm the absolute path with `pwd`.
2. Predict whether a given command is a shell built-in or an external program.
3. What happens when Ctrl+C is pressed while a command is running?
4. How does the shell find the program to run when you type a command name?
5. What is the difference between stdout and stderr?

## NEXT UP

[Lesson 05: Environment variables](../05-environment-variables/README.md)
