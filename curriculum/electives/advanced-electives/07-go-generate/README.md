# go generate

## When you need this

Go code generation tasks (mocking interfaces, generating protobuf stubs, embedding assets, creating stringer methods) are manual, error-prone, and must be re-run when source changes. Without a standard tool, each project invents its own generation script. go generate solves this by providing a standard, project-wide mechanism to document and invoke code generators via //go:generate directives in Go source files.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-06

## Mental model

go generate is a project-wide command runner. A //go:generate directive in a Go source file tells go generate to run a specific command. Running go generate ./... finds all directives in the project and executes them in order. Each directive is associated with a Go source file, making it clear which generators produce which files. The generated files are checked into version control so consumers do not need to run go generate.

## How it works

go generate scans .go files for lines matching //go:generate followed by a command and arguments. It executes the command in the file's directory. The command's stdout and stderr are printed to the terminal. go generate does not parse Go syntax — it only scans for the //go:generate prefix. Generated files typically have a _gen.go suffix or are written to a specified output directory.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: go generate")
}
```

## Tradeoffs

- go generate adds zero runtime overhead — generated code is regular Go code compiled into the binary
- Build time increases because generated files add more code to compile — but generation itself runs before compilation
- Large generated files (e.g., from protocol buffers) can significantly increase compile time and binary size
- Running go generate as part of CI adds latency — each CI run regenerates, even if nothing changed; use make or CI caching to skip unchanged generators

## Common mistakes

- Running go generate from the wrong directory — commands are executed relative to the file's directory, not the project root.
- Forgetting to check generated files into version control — consumers (CI, teammates) cannot build without running go generate first.
- Not making go generate idempotent — running go generate twice should produce the same output and not fail on the second run.
- Putting complex logic in the go:generate line — the directive should be a simple command invocation; complex logic belongs in a script.
- Not documenting which go:generate directives to run — some generates depend on others; document the execution order in a Makefile or README.

## Practice task

The learner must create a Go file with a //go:generate directive that runs stringer on an enum type, run go generate, verify the generated String() method compiles and works, and document the generation command in the README.

## Production notes

go generate is used in every Go project that uses code generation: protobuf (protoc --go_out), OpenAPI (oapi-codegen), mocking (mockgen, counterfeiter), stringer (stringer), and database (sqlc, ent). Production Go projects use go generate as a standard build step documented in CONTRIBUTING.md. CI runs go generate ./... and then git diff --exit-code to ensure generated files are up to date.

## Related core lessons

- Module 16.06 — Build tags
- Module 17.05 — Method values
