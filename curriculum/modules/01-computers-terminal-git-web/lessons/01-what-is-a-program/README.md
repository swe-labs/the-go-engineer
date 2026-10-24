# What is a program?

## Learning objective

Distinguish between source code text, compiled machine binaries, and active operating system processes in execution.

## Why this matters

New programmers conflate source code, executables, and running processes. This confusion makes debugging, compilation errors, and deployment seem like magic. A clear three-part model removes the mystery.

## Mental model

A computer program exists in three distinct states:
1. **Source Code (Passive Text):** Human-readable text files sitting inert on disk (e.g., `.go` files).
2. **Binary Executable (Passive Machine Code):** A formatted file of compiled machine instructions and static data on disk (ELF on Linux, Mach-O on macOS, PE on Windows) with a designated entry point.
3. **Process (Active Execution):** An instantiated instance of that binary loaded into system memory by the OS kernel, assigned its own private address space, PID, file descriptors, and execution threads.

Conflating these three states causes confusion when debugging builds versus runtime crashes.

## Visual model

```mermaid
flowchart LR
    SRC["1. Source Code
main.go
(Passive text on disk)"] -->|go build| BIN["2. Executable Binary
main (ELF / PE)
(Passive machine code on disk)"]
    BIN -->|OS Loader (execve)| PROC["3. Active Process
PID: 1234
(Allocated virtual memory & CPU execution)"]
```

## Core idea

Computers execute machine instructions, not human language. The concept of a "program" bridges the gap — it is the unit of work that the OS understands. Without a clear definition, every compile error, segfault, or deployment failure feels inexplicable.

## Under the hood

The Go compiler (gc) translates Go source through several stages: lexing splits text into tokens, parsing builds an AST, type-checking enforces correctness, SSA (static single assignment) optimization rewrites the IR, and code generation emits machine instructions. The linker resolves package imports into a single binary with a predetermined entry point.

## How Go uses it

Go compiles source code into a single static binary with no external dependencies. The `go build` command handles the entire pipeline: lexing, parsing, type-checking, code generation, and linking. Go programs start at package main's func main() — there is no hidden runtime initialization beyond what the Go runtime manages.

## Go example

The example program defines three stages of a program (source, binary, process) as an enum and prints each stage's description. It also displays runtime information about the current process: the binary path, compiler version, OS/architecture, and CPU count.

```go
package main

import (
	"fmt"
	"os"
	"runtime"
)

// ProgramStage represents one of the three fundamental states of a program.
type ProgramStage int

const (
	StageSource ProgramStage = iota
	StageBinary
	StageProcess
)

func (s ProgramStage) String() string {
	switch s {
	case StageSource:
		return "Source Code (inert human-readable text on disk)"
	case StageBinary:
		return "Executable Binary (compiled machine code on disk)"
	case StageProcess:
		return "Operating System Process (active execution in memory)"
	default:
		return "Unknown Stage"
	}
}

func main() {
	fmt.Println("=== Three Stages of a Program ===")
	stages := []ProgramStage{StageSource, StageBinary, StageProcess}
	for i, stage := range stages {
		fmt.Printf("Stage %d: %s\n", i+1, stage.String())
	}

	fmt.Println("\n=== Active Process Runtime Metadata ===")
	fmt.Printf("Executable Path (os.Args[0]): %s\n", os.Args[0])
	fmt.Printf("Compiler:                     %s\n", runtime.Compiler)
	fmt.Printf("Target OS / Architecture:     %s / %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Available Logical CPUs:       %d\n", runtime.NumCPU())
	fmt.Printf("Process PID:                  %d\n", os.Getpid())

	fmt.Println("\nSummary: Source code compiles into a binary executable, which the OS kernel executes as a process.")
}
```

## Step-by-step execution

1. The program defines a `ProgramStage` type with three constants: `StageSource`, `StageBinary`, and `StageProcess`.
2. Each stage has a `String()` method that returns a human-readable description.
3. `main()` iterates over all three stages and prints each one with its index.
4. The program then prints runtime info — the binary path (`os.Args[0]`), the Go compiler name, the OS and architecture, and the number of CPUs.
5. A final summary line reinforces that a program is source compiled into a binary loaded as a process.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Thinking a program is just a file you double-click | Ignoring the compilation and loading steps that turn source into a running process | Trace each step: write source, compile, load, execute |
| Confusing source code with the executable binary | Human-readable text looks like code; machine code is invisible | Remember: source is for humans, binaries are for CPUs |
| Assuming all programs are .exe files | Not recognizing scripts, bytecode, or JIT-compiled programs as valid programs | A program is any file the OS can execute, regardless of extension |

## Debugging walkthrough

**Scenario: A Go program compiles but does nothing when run.**
- **Cause:** The program has no `main` function, or `main()` returns immediately without calling any logic.
- **Diagnosis:** Check the file for `package main` and `func main()`. Add a `fmt.Println` call at the start of `main()`.
- **Resolution:** Ensure the entry point exists and contains the intended logic.

**Scenario: Editing source code but the old behavior persists.**
- **Cause:** The old compiled binary was run instead of recompiling — the source change was never built.
- **Diagnosis:** Check the binary's modification time vs the source file.
- **Resolution:** Always recompile after editing: run `go build` or `go run .` to produce a fresh binary.

## Production notes

Every day, developers write code, compile it, and run it. When a build fails, they read compiler errors. When a binary crashes, they inspect the process. Understanding the source-binary-process pipeline is the foundation of all software engineering.

## Performance implications

- Interpreted programs (Python, Ruby) start instantly but run 10-50x slower than compiled code.
- Compiled programs (Go, C, Rust) have a build step but run at near-native CPU speed.
- JIT-compiled programs (Java, JavaScript V8) warm up over time, starting slow then approaching compiled speed.

## Tests / verification

1. Open your terminal and run a command like `ping localhost` (or `Get-Process -Name "explorer"` on PowerShell). Note that every running command is a **process** with a PID.
2. Locate the executable on disk for a program you use daily (e.g., `where ping` / `Get-Command ping`). This is the **binary** — the compiled file stored on disk.
3. The source code that produced that binary is human-readable text. Edit a simple text file, compile it, and observe: changing the source requires recompilation to affect the running process. (If you have Go installed, write a one-line `hello.go` and compile with `go build`.)
4. Open Task Manager (Windows) or run `ps` (Linux/macOS) while a program runs. Confirm the process name matches a binary on disk and note its PID. Kill it with `kill <PID>` or `Stop-Process -Id <PID>` — the binary stays on disk but the process disappears.

## Review questions

1. Distinguish between source code, compiled binary, and running process from a list of descriptions.
2. What happens when a program with a syntax error is submitted to the compiler?
3. Why does deleting the .go source file not stop a running program?
4. What is the entry point for every Go program?
5. Why does editing source code not change the behavior of an already-running process?

## NEXT UP

[Lesson 02: Source code, executable, and process](../02-source-code-executable-and-process/README.md)
