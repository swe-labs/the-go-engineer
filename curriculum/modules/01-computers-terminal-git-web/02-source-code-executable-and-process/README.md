# Source code, executable, and process

## Learning objective

Understand and apply Source code, executable, and process in the context of professional Go software engineering.

## Why this matters

Beginners conflate the static file they write with the dynamic process that runs. This causes confusion about edits, compilation errors, and debugging — understanding the separation clarifies the entire development workflow.

## Mental model

Hello World is the minimal Go program: a package declaration, an import, and a main function calling fmt.Println. go run compiles and executes in one step. Every Go program starts with package main and func main.

## Core idea

Without this distinction, every development action — editing, compiling, running, debugging — becomes muddy. Separating source, executable, and process gives developers precise language to discuss what is happening and why.

## Under the hood

A Go binary uses the ELF (Linux) or PE (Windows) format. The OS loader reads the ELF/PE header, maps segments into virtual memory, resolves dynamic symbols (though Go prefers static linking), and jumps to the runtime entry point. The Go runtime then initializes the garbage collector, scheduler, and finally calls the user's main function.

## How Go uses it

Go's toolchain explicitly separates these three stages: source files (.go) are input, `go build` produces a standalone executable binary, and `go run` combines both steps plus execution in one command. Go binaries are statically linked — the executable contains everything needed to run, including the Go runtime.

## Go example

The example defines a `ProgramLifecycle` struct that holds the source file path, binary path, process PID, and command-line arguments. It uses `os.Executable()` to find the current binary, `os.Getpid()` for the process ID, and `os.Args` for arguments. It also uses `exec.LookPath` to show how the OS finds the binary via PATH resolution.

## Step-by-step execution

1. `main()` gets the current executable path via `os.Executable()` and the working directory via `os.Getwd()`.
2. A `ProgramLifecycle` struct is populated with the source path (main.go in the working directory), the binary path, the current PID, and the command-line arguments.
3. The `Describe()` method formats all fields into a human-readable string.
4. `exec.LookPath` resolves the binary name through the system PATH and prints the result.
5. A final summary restates the three-part model: source → compiled binary → running process.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Believing source code and running program are the same thing | Editing source while the program runs shows no behavior change | Recompile and restart the process after editing source |
| Thinking the executable contains the source code | Compiled binaries contain machine instructions, not human-readable text | Use `strings` on a binary to verify source is not embedded |
| Assuming every running program maps to exactly one file | Daemons and interpreters (Python, Node) can share a single binary | Understand that interpreters are the binary; scripts are data |

## Debugging walkthrough

**Scenario: A developer modifies source while a long-running process is active and expects the change to take effect immediately.**
- **Cause:** The running process loaded the old binary into memory; the new source must be recompiled and the process restarted.
- **Diagnosis:** Check that the binary timestamp is newer than the source file.
- **Resolution:** Stop the process, recompile, and start a new process from the updated binary.

**Scenario: A developer deletes the executable file while the process is still running and wonders why it continues to work.**
- **Cause:** On Unix, the OS keeps the file's inode alive as long as a process holds it open.
- **Diagnosis:** Running `lsof` on the PID will show the deleted file is still open.
- **Resolution:** Restart the process to load the intended binary, or use a deployment strategy that avoids in-use file conflicts.

## Production notes

Production deployments compile code into binaries (or containers), push those artifacts to servers, and restart processes. CI/CD pipelines automate this exact source-to-executable-to-process flow. Understanding each stage is essential for debugging, deployment, and incident response.

## Performance implications

- Source code is human-readable and compressible; executables are larger (2-20 MB for Go) and optimized for the CPU.
- The compilation step adds latency to the edit-run cycle — Go's fast compiler keeps this under 1 second for most projects.
- Statically linked binaries are larger but eliminate runtime dependency resolution, making deployment simpler and faster.

## Tests / verification

1. Pick any compiled program on your system (e.g., `ping`, `notepad`, or a Go-compiled binary). Use `where <name>` (Windows) or `which <name>` (Linux/macOS) to locate the **binary** on disk.
2. Run the program, then open a second terminal and use `ps` (Linux/macOS) or `Get-Process` (PowerShell) to find the running **process**. Note its PID.
3. Modify the source file of a simple program, recompile, and re-run. Confirm that changes appear only after the process is restarted. (If you have Go installed, write a minimal `hello.go`.)
4. The **source** is human-readable text; the **binary** is the compiled machine-code file; the **process** is the loaded program in memory with its PID.

## Review questions

1. Identify which artifact (source, executable, or process) is affected by editing, deleting, or killing.
2. Why can a compiled Go binary run on a machine without Go installed?
3. Trace the lifecycle of a program from keyboard input to terminal output through all three stages.
4. What does `os.Executable()` return in a Go program?
5. What happens when you delete the source file of a currently running process?

## NEXT UP

[Lesson 03: Files, bytes, directories, and paths](../03-files-bytes-directories-and-paths/README.md)
