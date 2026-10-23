# Tooling Failure Lab — Questions

Total: 14 questions (10 multiple choice × 1 pt, 4 short answer × 2 pt = 18 points). Passing: 10/14 correct (70%).

---

## Multiple choice (1 point each)

**1.** What is the distinct purpose of `GOROOT` versus `GOPATH` in modern Go?

a) `GOROOT` is where project source code lives; `GOPATH` is where the Go compiler binary is installed  
b) `GOROOT` points to the Go standard library and toolchain distribution; `GOPATH` stores downloaded third-party module caches and binaries installed via `go install`  
c) `GOROOT` is for compiling tests; `GOPATH` is for compiling production binaries  
d) Modern Go requires setting both `GOROOT` and `GOPATH` manually before any project can compile

**2.** What does `go.sum` contain and what is its primary security purpose?

a) A summary list of all exported functions in the root module  
b) Cryptographic SHA-256 checksums of third-party module zip files and `go.mod` files to verify dependency integrity and prevent supply-chain tampering  
c) A generated build script for Make and Docker  
d) The cumulative test coverage score across all sub-packages

**3.** Which bug pattern will `go vet` flag that the Go compiler successfully compiles without error?

a) Calling a non-existent method on a struct  
b) An unused local variable  
c) Unreachable code after an unconditional `return` statement or a `Printf` format string mismatch (e.g. `%d` with a string argument)  
d) A syntax error with an unclosed curly brace

**4.** How does `go run main.go` differ mechanically from `go build`?

a) `go run` compiles code into an in-memory interpreter without generating machine code  
b) `go run` compiles the binary into a temporary OS directory, executes the binary immediately, and cleans up the binary upon process exit; `go build` emits a permanent executable in the current working directory  
c) `go run` automatically disables race detection and garbage collection  
d) `go run` skips dependency resolution and ignores `go.mod`

**5.** When analyzing a panic traceback printed to `stderr`, where should an engineer look first to find the root application crash site?

a) The very top line stating the panic value, followed immediately by the first application-level stack frame (skipping runtime internal frames) indicating the source file and line number  
b) The very last line of the console output regardless of goroutine  
c) The list of OS environment variables printed before the stack  
d) The Git commit hash at the bottom of the trace

**6.** What language protocol server does `gopls` implement to deliver IDE features like jump-to-definition, autocomplete, and inline diagnostics?

a) RESTful OpenAPI 3.0  
b) Language Server Protocol (LSP) communicating over JSON-RPC  
c) gRPC over Unix Domain Sockets  
d) GraphQL WebSocket subscriptions

**7.** What does running `go mod tidy` do to your project?

a) Formats all `.go` files according to `gofmt` standards  
b) Adds module requirements for packages imported in your source code, removes requirements for packages no longer imported, and synchronizes `go.sum`  
c) Deletes the `go.mod` file and rebuilds the project from scratch  
d) Runs unit tests and lints across all subdirectories

**8.** What package and function signature are strictly required by the Go compiler to build a standalone runnable executable binary?

a) `package app` with `func Start()`  
b) `package main` with `func main()` taking no arguments and returning no values  
c) `package root` with `func Main(args []string) int`  
d) Any package name as long as it exports an `Init()` function

**9.** How can a Go engineer cross-compile a Linux 64-bit binary from a Windows or macOS workstation without installing third-party toolchains?

a) By running `go build --target=linux-amd64`  
b) By setting environment variables `GOOS=linux GOARCH=amd64` and running `go build`  
c) By installing a virtual machine running Ubuntu  
d) Cross-compilation is unsupported in standard Go

**10.** What guarantee does the official `gofmt` tool provide across Go codebases?

a) It guarantees that all code has 100% test coverage  
b) It parses source code into an Abstract Syntax Tree (AST) and prints it with deterministic canonical whitespace, indentation, and alignment, eliminating formatting style debates  
c) It automatically refactors imperative loops into functional pipelines  
d) It renames identifiers to follow camelCase conventions

---

## Short answer (2 points each)

**11.** Trace the compilation pipeline executed by the Go compiler (`gc`) when turning source code text into a runnable binary. Name the key phases (lexing, parsing, type-checking, SSA optimization, and code generation/linking) and briefly describe what each phase produces.

**12.** A junior developer runs `go build` and encounters: `build command-line-arguments: cannot find package "github.com/google/uuid" in any of...`. Explain three distinct root causes that trigger this error and the exact commands to diagnose and resolve each.

**13.** Explain how to read and interpret a Go panic traceback when a goroutine crashes. Detail the distinction between runtime frames (e.g., `runtime.gopanic`) and application frames, and explain how to locate the exact file, line number, and function where the crash originated.

**14.** Contrast the distinct roles of `go.mod` versus `go.sum`. Explain why both files must be committed to version control and how `go.sum` protects teams against malicious dependency modifications (supply chain attacks).
