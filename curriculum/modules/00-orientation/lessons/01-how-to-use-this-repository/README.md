# Lesson 01: How to use this repository

## Learning objective

By the end of this lesson, you will be able to open this repository and immediately identify what each top-level directory holds, how metadata connects to learner-facing content, and which tools validate the connection. You will demonstrate this by running a Go program that models the directory structure as searchable data.

## Why this matters

Every engineering project depends on a shared mental model of its own layout. When a team member asks "Where do I find the lesson files?" or "Which tool checks for missing READMEs?", the answer should come from the structure itself, not from a person. This repository enforces that discipline from day one. Learning the layout now prevents confusion across every subsequent module.

## Mental model

Picture a library. The card catalog (metadata/) describes every book and where it sits. The shelves (curriculum/) hold the books you actually read. The librarians (tools/) check that books are in the right place and that the catalog is accurate. The rules posted on the wall (docs/) tell everyone how things work. The loading dock (dist/) is where packed boxes go before shipment.

When you wonder "Why does this file exist here?", trace it back to one of these five categories. Every file in the root belongs to exactly one.

## Core idea

The repository separates five concerns so that humans and machines each see what they need:

| Directory | Concern | Audience |
|-----------|---------|----------|
| `metadata/` | Curriculum map and lesson graph | Validators, auditors |
| `curriculum/` | Learner-facing lessons, labs, projects | Learners, instructors |
| `tools/` | Validation and generation scripts | CI, maintainers |
| `docs/` | Contribution and style guides | Contributors |
| `dist/` | Build artifacts and exports | Release pipeline |

The most important relationship is between `metadata/` and `curriculum/`. The metadata JSON declares which lessons exist and in what order. The curriculum folder holds the actual README files and code. A validator tool reads the metadata, follows the paths into curriculum, and reports mismatches. If a README is missing, the validator fails. If metadata points to a path that doesn't exist, the validator fails. This creates a contract: if the graph says it exists, the filesystem must agree.

## Under the hood

The `main.go` file in this lesson defines a `DirEntry` struct and a `repoDirs` slice that hard-codes each top-level directory's name, purpose, and typical contents. The `Lookup` function searches the slice by name and returns a pointer to the matching entry or nil if nothing is found.

This is a simple in-memory lookup, but it mirrors how the real curriculum validator works: it reads a metadata graph, follows paths, and checks for file existence. The principle is the same — use data to describe expected structure, then write code to verify it.

## How Go uses it

Go's own source tree follows a similar convention. The `src/` directory holds standard library packages. The `cmd/` directory holds compilers and tools. The `go/types` package is the type checker. Each directory has a clearly documented purpose. When you run `go vet`, Go knows where to find its analysis packages because the structure is predictable.

In this curriculum, we apply the same discipline. Every module has a predictable shape: `lessons/`, `labs/`, `projects/`, and within each, a README, a `main.go`, and a `main_test.go`. This predictability means you can jump to any module and immediately know where to look.

## Go example

```go
package main

import (
	"fmt"
)

type DirEntry struct {
	Name     string
	Purpose  string
	Contents []string
}

var repoDirs = []DirEntry{
	{
		Name:    "metadata/",
		Purpose: "defines the curriculum map: what lessons exist, their order, and how they connect",
		Contents: []string{
			"curriculum.json         -- full lesson graph",
			"modules/                -- module definitions",
		},
	},
	{
		Name:    "curriculum/",
		Purpose: "holds learner-facing content: READMEs, code, starter and solution files",
		Contents: []string{
			"modules/00-orientation/ -- Module 0: getting started",
			"modules/01-go-basics/   -- Module 1: Go fundamentals",
		},
	},
	{
		Name:    "tools/",
		Purpose: "validators, auditors, and generators that check curriculum correctness",
		Contents: []string{
			"validate/   -- checks metadata against filesystem",
			"audit/      -- reviews consistency and completeness",
		},
	},
	{
		Name:    "docs/",
		Purpose: "maintainer documentation: style guide, review process, conventions",
		Contents: []string{
			"CONTRIBUTING.md",
			"style-guide.md",
		},
	},
	{
		Name:    "dist/",
		Purpose: "generated output: exports, bundles, release artifacts (git-ignored)",
		Contents: []string{
			"(generated at build time)",
		},
	},
}

func Lookup(name string) *DirEntry {
	for i := range repoDirs {
		if repoDirs[i].Name == name || repoDirs[i].Name == name+"/" {
			return &repoDirs[i]
		}
	}
	return nil
}

func main() {
	fmt.Println("The Go Engineer — Repository Map")
	fmt.Println()
	for _, d := range repoDirs {
		fmt.Printf("  %-14s %s\n", d.Name, d.Purpose)
	}
	fmt.Println()
	fmt.Println("Use Lookup(name) to get details about any top-level directory.")
}
```

## Step-by-step execution

1. The program defines `DirEntry`, a struct with three fields: `Name`, `Purpose`, and `Contents`.
2. `repoDirs` is a slice that holds one `DirEntry` for each top-level directory in the repository.
3. `Lookup("curriculum/")` iterates over `repoDirs` and compares each entry's `Name` field against the argument. It returns the matching entry or `nil`.
4. `main()` prints a formatted table of all directories and their purposes.
5. When you run `go run .`, you see the full map. When you run `go test .`, the tests call `Lookup` with known names and verify it returns non-nil entries with non-empty `Purpose` fields.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Adding a trailing slash inconsistently | The `Lookup` function accepts both `"metadata"` and `"metadata/"`, but if you write code that always appends a slash before calling `Lookup`, you might match an entry that the caller didn't intend. | Always strip trailing slashes before lookup, or always include them; pick one convention and stick to it. |
| Confusing `metadata/` with `curriculum/` | Both contain JSON-like structures. Metadata describes the intended curriculum; curriculum holds the actual content. | Remember: metadata is the map, curriculum is the territory. |
| Assuming `dist/` is hand-written | Generated directories should not be edited directly. If you need to change a dist artifact, change the source that generates it. | Check `tools/` for the generator; never edit `dist/` by hand. |
| Forgetting that `Lookup` returns `nil` | Go programs that dereference a nil pointer panic. | Always check `Lookup` result for nil before using it. |

## Debugging walkthrough

Scenario: You run `go run .` and see no output at all.

Step 1: Check that the file compiles. Run `go build .` — if it succeeds, the code is syntactically correct. If it fails, the compiler will tell you the exact line and error.

Step 2: If the build succeeds but `go run .` produces nothing, suspect that `main()` is not printing. Add a `fmt.Println("start")` as the first line of `main()` and run again. If "start" appears but the rest doesn't, the error is after that line.

Step 3: Check for a runtime panic. Go panics produce a stack trace. If the program crashes silently, redirect stderr: `go run . 2>&1`.

Step 4: If `go test .` fails with "nil pointer dereference", the `Lookup` function returned `nil` and the test tried to access a field on it. Look at which name was passed — it might have a typo (e.g., `"curriculm/"` instead of `"curriculum/"`).

## Production notes

- The `Lookup` function is O(n) where n is the number of directories. For this curriculum (5 top-level dirs), that is instant. A production service with thousands of entries would use a `map[string]DirEntry` for O(1) lookups.
- Hard-coding directory data in a Go slice is fine for this lesson. In the real curriculum, directory data comes from the filesystem and metadata JSON, not from hard-coded strings.
- When adding a new top-level directory, you must update `repoDirs` and the metadata JSON. The validator will fail if only one is updated.

## Performance implications

- `Lookup` performs a linear scan over 5 entries. That is negligible.
- Memory usage: each `DirEntry` stores strings. Strings in Go are immutable headers (pointer + length), so the total memory is roughly the size of the string data plus small struct overhead.
- If this lookup were called millions of times per second, you would switch to a map. For a CLI tool or curriculum validator, linear scan is perfectly adequate.

## Practice task

Add a new top-level directory to `repoDirs` called `"scripts/"` with purpose "helper scripts for common development tasks" and contents `["setup.sh", "lint.sh"]`. Run `go test .` to verify that `Lookup("scripts/")` now returns the new entry, and update the tests to include the new directory in the expected lookup results.

## Tests / verification

Run the tests from the repository root:

```bash
go test ./curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/
```

Expected output:

```
ok      github.com/rasel9t6/the-go-engineer/curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository
```

## Review questions

1. What does the `Lookup` function return when given `"metadata"` (without a trailing slash)?
2. Why does the curriculum separate metadata from curriculum content instead of combining them into one directory?
3. What would break if you deleted `curriculum/modules/00-orientation/` but left the metadata JSON unchanged?
4. In the `DirEntry` struct, which field stores the list of files inside the directory?
5. How would you change the Go example to use a map instead of a slice for constant-time lookups?

## NEXT UP

Lesson 02: What zero magic means — where you learn why every concept in this curriculum is explained explicitly, with no hidden surprises.
