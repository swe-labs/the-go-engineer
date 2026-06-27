# What is a program?

## Learning objective

Understand and apply What is a program? in the context of professional Go software engineering.

## Why this matters

New programmers conflate source code, executables, and running processes. This confusion makes debugging, compilation errors, and deployment seem like magic. A clear three-part model removes the mystery.

## Mental model

Installing Go means downloading the distribution, setting GOROOT and GOPATH, and verifying with go version. The toolchain is self-contained — one download includes the compiler, linker, formatter, and package manager.

## Core idea

Computers execute machine instructions, not human language. The concept of a "program" bridges the gap — it is the unit of work that the OS understands. Without a clear definition, every compile error, segfault, or deployment failure feels inexplicable.

## Under the hood

The Go compiler (gc) translates Go source through several stages: lexing splits text into tokens, parsing builds an AST, type-checking enforces correctness, SSA (static single assignment) optimization rewrites the IR, and code generation emits machine instructions. The linker resolves package imports into a single binary with a predetermined entry point.

## How Go uses it

Go compiles source code into a single static binary with no external dependencies. The `go build` command handles the entire pipeline: lexing, parsing, type-checking, code generation, and linking. Go programs start at package main's func main() — there is no hidden runtime initialization beyond what the Go runtime manages.

## Go example

The example program defines three stages of a program (source, binary, process) as an enum and prints each stage's description. It also displays runtime information about the current process: the binary path, compiler version, OS/architecture, and CPU count.

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

## Practice task

Write a Go program that prints "Hello, program!", compile it with `go build`, run the resulting binary, then explain the role of each artifact (source, binary, process).

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/01-what-is-a-program/
```

## Review questions

1. Distinguish between source code, compiled binary, and running process from a list of descriptions.
2. What happens when a program with a syntax error is submitted to the compiler?
3. Why does deleting the .go source file not stop a running program?
4. What is the entry point for every Go program?
5. Why does editing source code not change the behavior of an already-running process?

## NEXT UP

[Lesson 02: Source code, executable, and process](../02-source-code-executable-and-process/README.md)
