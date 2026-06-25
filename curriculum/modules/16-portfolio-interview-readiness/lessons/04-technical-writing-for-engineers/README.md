# Technical writing for engineers

## Learning objective

Write clear, effective technical documentation for Go projects including README files, API documentation, architecture docs, and Go doc comments that render correctly with `godoc`.

## Why this matters

Code tells you what the computer does. Documentation tells humans why and how. The best Go projects fail without good documentation because nobody can understand the purpose, the API contract, or the architecture. In an interview, the quality of your documentation is as visible as the quality of your code. Engineers who write clear docs get faster onboarding, fewer questions, and more credibility. Documentation is a force multiplier for your impact.

## Mental model

Think of documentation as code comments that compile into understanding. The Go compiler ignores them, but the team reads them. Documentation exists on five levels:

1. **Package-level**: What this package does and who should use it.
2. **Symbol-level**: What each exported function, type, and constant does.
3. **Usage-level**: How to install, configure, and run the project (README).
4. **Architecture-level**: How the system is structured and why.
5. **API-level**: The complete contract for external consumers.

Each level serves a different audience. Package-level documentation is for other engineers importing your package. README-level documentation is for users deploying your project. Architecture documentation is for future maintainers.

Go enforces a documentation culture through `godoc`, which automatically renders comments as documentation. A documented exported symbol in Go is the norm, not the exception.

## Core idea

**Go doc comments** follow a strict convention:

- A doc comment is a complete sentence that starts with the symbol name.
- `// Package math provides basic constants and mathematical functions.`
- `// Parse parses a string and returns a Value.`
- The first sentence is the doc comment summary and should stand alone.
- Paragraphs are separated by blank comment lines.
- Code blocks are indented or wrapped in backticks.
- `Deprecated` annotations signal symbols that should not be used.

**README structure** for a portfolio project:

```
# Project name
Short description (one paragraph).

## Features
Bullet list of key capabilities.

## Quick start
One-line install, minimal example to run.

## Architecture
High-level diagram or description of components.

## API
Link to godoc or inline documentation.

## Configuration
Environment variables, config file format.

## Testing
How to run tests, CI badge.

## Contributing
How to contribute (optional for portfolio).
```

**Documentation coverage** is the percentage of exported symbols that have doc comments. A professional Go project targets 100% coverage on exported symbols. Unexported symbols may omit comments if the code is self-documenting, but even they benefit from brief notes.

## Under the hood

The `go/parser` package in Go's standard library parses source code into an AST. The parser can optionally parse comments with `parser.ParseComments`. The AST represents each comment as a `CommentGroup` attached to declarations via the `Doc` field:

- `ast.FuncDecl.Doc` — doc comment for a function or method.
- `ast.GenDecl.Doc` — doc comment for a type, var, or const declaration group.
- `ast.TypeSpec.Doc` — doc comment for a specific type within a group.

The `go/doc` package provides higher-level documentation extraction, organizing symbols by package and respecting the `doc.go` convention.

The documentation coverage checker in this lesson uses `go/parser` at the AST level. This is the same parser the Go compiler and `godoc` use, so the results are authoritative.

## How Go uses it

Go is the only mainstream language with documentation generation built into the toolchain:

- `go doc` prints documentation for any package, symbol, or method from the command line.
- `go doc -http=:8080` launches a godoc HTTP server.
- `pkg.go.dev` hosts documentation for all public Go modules.

The `golang.org/x/tools/cmd/godoc` tool (or `go doc` starting from Go 1.19) renders documentation from source comments. No separate build step or tool is required — documentation is always up to date because it lives with the source.

Go's documentation conventions also enforce good habits:

- Package comments must be in the package clause form (`Package x does y.`).
- Bug annotations (`// BUG(name): description`) are rendered separately in godoc.
- Example functions (`ExampleFoo`) are runnable documentation that is tested.

## Go example

The documentation coverage checker parses a Go source file and reports which exported symbols have doc comments and which are missing them.

```go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type DocStatus int

const (
	Missing DocStatus = iota
	Present
)

func (d DocStatus) String() string {
	switch d {
	case Missing:
		return "Missing"
	case Present:
		return "Present"
	default:
		return "Unknown"
	}
}

type DocItem struct {
	Name    string
	Kind    string
	Status  DocStatus
	Line    int
	DocText string
}

type DocReport struct {
	FilePath string
	Items    []DocItem
	Covered  int
	Total    int
}

func (r DocReport) Coverage() float64 {
	if r.Total == 0 {
		return 100.0
	}
	return float64(r.Covered) / float64(r.Total) * 100.0
}

func AnalyzeFile(filePath string) (DocReport, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return DocReport{}, fmt.Errorf("parse error: %w", err)
	}

	report := DocReport{FilePath: filePath}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE || d.Tok == token.VAR || d.Tok == token.CONST {
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						doc := s.Doc
						if doc == nil {
							doc = d.Doc
						}
						item := DocItem{
							Name:   s.Name.Name,
							Kind:   "type",
							Status: docStatusFromComment(doc),
							Line:   fset.Position(s.Pos()).Line,
						}
						if doc != nil {
							item.DocText = doc.Text()
						}
						report.Items = append(report.Items, item)
						if item.Status == Present {
							report.Covered++
						}
						report.Total++
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.IsExported() {
								item := DocItem{
									Name:   name.Name,
									Kind:   varOrConst(d.Tok),
									Status: docStatusFromComment(s.Doc),
									Line:   fset.Position(s.Pos()).Line,
								}
								if s.Doc != nil {
									item.DocText = s.Doc.Text()
								}
								report.Items = append(report.Items, item)
								if item.Status == Present {
									report.Covered++
								}
								report.Total++
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if d.Name.IsExported() {
				item := DocItem{
					Name:   d.Name.Name,
					Kind:   funcKind(d),
					Status: docStatusFromComment(d.Doc),
					Line:   fset.Position(d.Pos()).Line,
				}
				if d.Doc != nil {
					item.DocText = d.Doc.Text()
				}
				report.Items = append(report.Items, item)
				if item.Status == Present {
					report.Covered++
				}
				report.Total++
			}
		}
	}

	return report, nil
}

func docStatusFromComment(doc *ast.CommentGroup) DocStatus {
	if doc != nil && len(doc.List) > 0 {
		return Present
	}
	return Missing
}

func varOrConst(tok token.Token) string {
	if tok == token.CONST {
		return "const"
	}
	return "var"
}

func funcKind(d *ast.FuncDecl) string {
	if d.Recv != nil {
		return "method"
	}
	return "function"
}

func (r DocReport) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Documentation Coverage Report: %s\n", r.FilePath))
	b.WriteString(fmt.Sprintf("Coverage: %.1f%% (%d/%d exported symbols documented)\n\n", r.Coverage(), r.Covered, r.Total))
	for _, item := range r.Items {
		b.WriteString(fmt.Sprintf("  %s %s (line %d): %s\n", item.Kind, item.Name, item.Line, item.Status))
		if item.DocText != "" {
			b.WriteString(fmt.Sprintf("    %s", item.DocText))
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <file.go>")
		os.Exit(1)
	}

	report, err := AnalyzeFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(report.Summary())
}
```

## Step-by-step execution

1. `AnalyzeFile` takes a file path and uses `go/parser` to parse the source into an AST with comments.
2. It iterates over all top-level declarations (`f.Decls`), distinguishing between `GenDecl` (type, var, const) and `FuncDecl` (function, method).
3. For `FuncDecl`, it checks `d.Name.IsExported()` and records the doc comment status from `d.Doc`.
4. For `GenDecl`, it iterates individual specs. TypeSpecs get doc from `s.Doc` (falling back to `d.Doc`), and ValueSpecs check each exported variable or constant name.
5. Each item records the name, kind, line number, whether documented, and the comment text.
6. Coverage is calculated as `Covered / Total * 100`, with edge case for zero exported symbols.
7. `Summary()` formats the report with per-symbol status and the overall coverage percentage.

## Common mistakes

- Mistake: Writing `// This function parses a string` instead of `// Parse parses a string`.
  - Why it happens: It feels natural to write "This function" as the subject.
  - Fix: Always start with the symbol name. `godoc` renders the name automatically, so "Parse parses" becomes redundant but is the convention.

- Mistake: Leaving exported symbols undocumented because the code is "obvious."
  - Why it happens: What is obvious to the author is not obvious to a reader seeing the code for the first time.
  - Fix: Document every exported symbol. If you cannot write a one-sentence doc comment, the function might have too many responsibilities.

- Mistake: Writing a README that only says "See code for details."
  - Why it happens: The code exists and works; the README feels secondary.
  - Fix: The README is the front door. Write it first. It forces you to clarify the project's purpose before writing code.

- Mistake: Updating code but forgetting to update comments.
  - Why it happens: Code changes are visible; comment drift is silent.
  - Fix: Review documentation changes alongside code changes in every PR.

## Debugging walkthrough

Running the coverage checker on a file with:

```go
package config
// DefaultPort is the default server port.
const DefaultPort = 8080
func Load(path string) (*Config, error) { return nil, nil }
```

Reports:

```
Coverage: 50.0% (1/2 exported symbols documented)

  function Load (line 4): Missing
  const DefaultPort (line 2): Present
    DefaultPort is the default server port.
```

The developer adds a doc comment for `Load` and re-runs, achieving 100%. The checker also reveals that `Config` is unexported (returns `*Config` but `Config` is defined elsewhere or is a type in the same file that is not documented). The developer must document `Config` separately.

## Production notes

In a production project:

- Integrate the coverage checker into CI. Require 100% documentation coverage on exported symbols for library packages.
- For application code (main packages), a lower threshold (80%) may be acceptable for internal functions.
- Use `golangci-lint` with the `godox` linter to catch outstanding work-in-progress comments that should be resolved.
- Write example functions for every public API. Examples are tested documentation.
- Consider a `doc.go` file for package-level documentation that spans multiple files.

## Performance implications

The AST-based parser is efficient. Parsing a file of 1000 lines takes under a millisecond. For large projects, parse concurrently using goroutines — each file is independent. The checker does not run the compiler, so there is no type-checking overhead. If you need type information, use `go/types` which adds significant computation.

## Practice task

Extend the documentation checker to also detect deprecated symbols. Add a `Deprecated` field to `DocItem`. A symbol is deprecated if its doc comment contains `Deprecated:` on a line by itself. Update the report to show deprecated symbols. Add a test that verifies a deprecated function is detected.

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/04-technical-writing-for-engineers
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/04-technical-writing-for-engineers main.go
```

The tests verify detection of undocumented and documented exported symbols, correct coverage calculation, and handling of files with no exported symbols. After the practice task, the deprecated detection tests pass.

## Review questions

1. What are the five levels of documentation, and who is the audience for each?
2. What is the Go convention for doc comments on exported functions?
3. Why should you write the README before writing the code?
4. What does `parser.ParseComments` do, and why is it necessary for the coverage checker?
5. How do you mark a function as deprecated in Go, and how does `godoc` render it?

## NEXT UP

Resume and LinkedIn optimization — learn to craft a resume and LinkedIn profile that passes keyword screens and tells a compelling career story.
