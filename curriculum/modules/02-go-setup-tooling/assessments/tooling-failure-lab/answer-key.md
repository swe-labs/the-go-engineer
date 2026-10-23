# Tooling Failure Lab — Answer Key

## Multiple choice (1 point each)

**1. b) `GOROOT` points to the Go standard library and toolchain distribution; `GOPATH` stores downloaded third-party module caches and binaries installed via `go install`**  
*Explanation:* `GOROOT` contains the official Go compiler and standard library packages (e.g. `/usr/local/go` or `C:\Program Files\Go`). In modern Go modules mode, `GOPATH` defaults to `~/go` and houses downloaded third-party packages in `$GOPATH/pkg/mod` and compiled tool binaries in `$GOPATH/bin`.

**2. b) Cryptographic SHA-256 checksums of third-party module zip files and `go.mod` files to verify dependency integrity and prevent supply-chain tampering**  
*Explanation:* `go.sum` contains cryptographic hashes of specific module version contents (`h1:...`). When Go downloads a module, it verifies that the downloaded archive matches the hash in `go.sum` (and the Go Checksum Database `sum.golang.org`), preventing compromised or modified upstream tags.

**3. c) Unreachable code after an unconditional `return` statement or a `Printf` format string mismatch (e.g. `%d` with a string argument)**  
*Explanation:* The Go compiler permits valid syntax that compiles cleanly, such as passing a string to `fmt.Printf("%d", "hello")`. `go vet` uses static analysis heuristics to detect suspicious constructs, format mismatches, mutex copying, and unreachable code.

**4. b) `go run` compiles the binary into a temporary OS directory, executes the binary immediately, and cleans up the binary upon process exit; `go build` emits a permanent executable in the current working directory**  
*Explanation:* `go run` creates a temporary directory in `$TMPDIR`, compiles the executable, runs it with current process arguments, and removes the temp binary after execution. `go build` creates the final binary in the current directory or target path.

**5. a) The very top line stating the panic value, followed immediately by the first application-level stack frame (skipping runtime internal frames) indicating the source file and line number**  
*Explanation:* Go panics print the panic value followed by the stack trace of the active goroutine. The frames immediately below the panic line typically reference internal runtime functions (`runtime.gopanic`). The first frame citing the project's source code path and line number identifies the exact crash origin.

**6. b) Language Server Protocol (LSP) communicating over JSON-RPC**  
*Explanation:* `gopls` is the official Go language server that speaks the Language Server Protocol (LSP), allowing editors (VS Code, Neovim, Goland) to communicate over standard JSON-RPC for hover docs, navigation, and refactoring.

**7. b) Adds module requirements for packages imported in your source code, removes requirements for packages no longer imported, and synchronizes `go.sum`**  
*Explanation:* `go mod tidy` scans the entire project source tree, adds missing module requirements, removes obsolete ones, and cleans up unused entries in `go.sum`.

**8. b) `package main` with `func main()` taking no arguments and returning no values**  
*Explanation:* In Go, a standalone runnable program must belong to package `main` and declare a parameterless, returnless `func main()`. Command-line arguments are accessed via `os.Args` and exit status via `os.Exit()`.

**9. b) By setting environment variables `GOOS=linux GOARCH=amd64` and running `go build`**  
*Explanation:* Because Go's compiler and standard library are written in Go without external C runtime dependencies (unless CGO is enabled), cross-compiling only requires setting `GOOS` (target OS) and `GOARCH` (target CPU architecture) during `go build`.

**10. b) It parses source code into an Abstract Syntax Tree (AST) and prints it with deterministic canonical whitespace, indentation, and alignment, eliminating formatting style debates**  
*Explanation:* `gofmt` re-renders the parsed AST using strict, standardized formatting rules (tabs for indentation, standard brace placement). This ensures universal visual consistency across the entire Go ecosystem.

---

## Short answer (2 points each)

**11.** The Go compilation pipeline (`gc`) operates in the following stages:
1. **Lexical Analysis (Scanning/Lexing):** Converts raw `.go` UTF-8 source characters into a stream of syntax tokens (keywords, identifiers, literals, operators).
2. **Parsing:** Builds an Abstract Syntax Tree (AST) from the token stream, validating grammatical rules and structure.
3. **Type-Checking:** Traverses the AST to check type compatibility, identifier declarations, interface fulfillment, and constant evaluation.
4. **Intermediate Representation & SSA (Static Single Assignment) Optimization:** Transforms the AST into SSA form. In SSA, every variable is assigned exactly once, allowing aggressive optimizations such as dead-code elimination, nil-check removal, bounds-check elimination, and function inlining.
5. **Machine Code Generation & Linking:** Translates optimized SSA instructions into architecture-specific machine code (e.g. x86-64 or ARM64) and links packages, runtime routines, and data sections into a single executable binary.

*Grading: 1 point for identifying the key stages (lexing, parsing, type-checking, SSA, codegen/linking); 1 point for explaining the transformation from source text to AST to SSA to machine binary.*

**12.** Three common root causes and resolutions:
1. **Module Not Downloaded / Missing from `go.mod`:** The code imports `github.com/google/uuid` but it hasn't been fetched into the local module cache.  
   *Resolution:* Run `go get github.com/google/uuid` or `go mod tidy`.
2. **Project Missing `go.mod` (Legacy GOPATH mode):** The project is located outside `$GOPATH/src` without a `go.mod` file initialized, so Go does not know how to resolve third-party modules.  
   *Resolution:* Initialize the module with `go mod init <module-path>` and run `go mod tidy`.
3. **Typo in Import Path or Package Identifier:** The import path in the `.go` file is misspelled (e.g., `github.com/google/uuidv4` instead of `github.com/google/uuid`).  
   *Resolution:* Inspect the exact import statement against upstream documentation and correct the import string.

*Grading: 1 point for identifying at least two distinct root causes; 1 point for providing accurate diagnostic/remediation commands (`go mod tidy`, `go get`, `go mod init`).*

**13.** Reading and interpreting a Go panic traceback:
1. **Root Panic Message:** The first line identifies the panic reason (e.g., `panic: runtime error: invalid memory address or nil pointer dereference` or `panic: index out of range [5] with length 3`).
2. **Goroutine State:** The traceback indicates which goroutine caused the panic (e.g., `goroutine 1 [running]:`).
3. **Stack Frames Hierarchy:** The call stack lists frames from top (most recent call) to bottom (earliest call). The topmost frames often include internal Go runtime functions like `runtime.gopanic` or `runtime.panicmem`, which should be bypassed when locating application bugs.
4. **Pinpointing the Application Origin:** Look for the highest frame referencing your project's package and file path (e.g., `main.processOrder(0x0)\n\t/home/user/app/main.go:42 +0x2b`). The line number (`main.go:42`) indicates the exact source line where the invalid operation was attempted.

*Grading: 1 point for explaining how to identify the root cause and distinguish internal runtime frames from application frames; 1 point for explaining how to navigate the file and line number coordinates.*

**14.** Contrast between `go.mod` and `go.sum`:
- **`go.mod` (Intent and Declarative Dependencies):** Defines the module's identity path, the minimum required Go language version, and the minimum required semantic versions of direct and indirect dependencies (`require` directives).
- **`go.sum` (Cryptographic Verification and Tamper Proofing):** Contains cryptographic SHA-256 checksums (`h1:...`) for every downloaded dependency module archive and its `go.mod` file.
- **Supply-Chain Security:** Without `go.sum`, an attacker who compromises an upstream GitHub account or proxy could replace an existing release tag (e.g., `v1.2.0`) with malicious code. When a teammate or CI server builds the project, Go checks the downloaded content against the recorded checksum in `go.sum`. If the hashes mismatch, Go halts compilation immediately with a security verification error. Both files must be committed to ensure reproducible, secure builds across all developer workstations and CI runners.

*Grading: 1 point for articulating the distinct purposes of `go.mod` (declarative versioning) vs `go.sum` (cryptographic integrity); 1 point for explaining how `go.sum` protects against upstream tampering and supply-chain attacks.*
