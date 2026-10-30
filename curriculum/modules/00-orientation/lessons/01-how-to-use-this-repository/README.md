# How to use this repository

## Learning objective

By the end of this lesson, you will be able to open this repository and immediately identify what each top-level directory holds, how metadata connects to learner-facing content, and which tools validate the connection. You will demonstrate this by running a Go program that models the directory structure as searchable data.

## Why this matters

Every engineering project depends on a shared mental model of its own layout. When a team member asks "Where do I find the lesson files?" or "Which tool checks for missing READMEs?", the answer should come from the structure itself, not from a person. This repository enforces that discipline from day one. Learning the layout now prevents confusion across every subsequent module.

## Mental model

Picture a modern software engineering system. The curriculum core (`curriculum/`) holds both the authoritative machine-readable blueprints (such as `path.core.json` and `projects.json`) and the physical instructional modules (`modules/` and `electives/`). The internal toolchain (`internal/tools/curriculum/`) serves as the automated quality compiler, executing 14 strict quality gates. The learner web platform (`web/`) visualizes your progression graph. The agent governance framework (`.agents/`) maintains pedagogical standards.

When you wonder "Why does this file exist here?", trace it back to one of these core system layers. Every file in the root belongs to a clearly defined architectural concern.

## Core idea

The repository separates concerns so that learners, maintainers, and automated validators each have clear boundaries:

| Directory | Concern | Audience |
|-----------|---------|----------|
| `curriculum/` | Manifest blueprints (`*.json`), lesson modules (`modules/`), electives (`electives/`), and companion frontend (`companion-frontend/`) | Learners, validators, instructors |
| `internal/tools/` | Automated quality gates, graph validation, and verification CLI | CI pipelines, maintainers |
| `web/` | Learner portal, curriculum graph visualizer, and dashboard | Learners, reviewers |
| `.agents/` | Pedagogical skills, governance rubrics, and authoring workflows | Engineering mentors, AI assistants |

The most important relationship is between the **manifest blueprints** inside `curriculum/` (such as `path.core.json`, `projects.json`, and `schema.v3.json`) and the **physical filesystem** (`curriculum/modules/`). The metadata JSON declares which lessons exist, their prerequisites, and their contracts. The module directories hold the actual README files, code, tests, and starter scaffolds. The validator tool (`internal/tools/curriculum`) reads the metadata, walks the filesystem, and verifies strict alignment. If an artifact is missing or invalid, the quality gates fail.

## Under the hood

The Go example below defines a `DirEntry` struct and a `repoDirs` slice that models the v3 repository layout. The `Lookup` function searches the slice by name and returns a pointer to the matching entry or nil if nothing is found.

This in-memory lookup mirrors how the curriculum validator inspects repository paths and enforces integrity contracts across the codebase.

## How Go uses it

Go's own source tree and production Go services follow clear package boundaries. Standard packages live in `src/`, compiler tools in `cmd/`, and internal unexported utilities in `internal/`. By placing validation tools under `internal/tools/curriculum/`, the repository prevents external packages from importing internal CLI commands while maintaining clean dependency isolation.

In this curriculum, every module follows a predictable shape: `lessons/` (or direct topic folders), `projects/` (with `_starter/` and `_solution/`), and `assessments/` (with explicit rubrics). This structure ensures you always know where to find what you need.

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
		Name:    "curriculum/",
		Purpose: "authoritative metadata blueprints (path.core.json, projects.json) and instructional modules",
		Contents: []string{
			"path.core.json          -- authoritative core learning graph",
			"projects.json           -- portfolio project registry and rubrics",
			"modules/00-orientation/ -- Module 0: orientation and setup",
			"modules/01-...          -- technical modules and portfolio code",
			"electives/              -- advanced elective tracks",
			"companion-frontend/     -- 4-stage progressive server-driven hypermedia track",
		},
	},
	{
		Name:    "internal/tools/",
		Purpose: "curriculum quality gates and automated validation tooling",
		Contents: []string{
			"curriculum/main.go      -- 14 automated quality validators (schema, graph, zero-magic)",
		},
	},
	{
		Name:    "web/",
		Purpose: "learner web application, graph visualizer, and progress tracking dashboard",
		Contents: []string{
			"app/                    -- frontend interface for navigating curriculum paths",
		},
	},
	{
		Name:    ".agents/",
		Purpose: "curriculum skills, quality rubrics, and automated authoring assistants",
		Contents: []string{
			"skills/                 -- domain workflows (curriculum-validator, project-mentor, etc.)",
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
	fmt.Println("The Go Engineer v3 — Repository Map")
	fmt.Println()
	for _, d := range repoDirs {
		fmt.Printf("  %-16s %s\n", d.Name, d.Purpose)
	}
	fmt.Println()
	fmt.Println("Use Lookup(name) to get details about any repository subsystem.")
}
```

## Step-by-step execution

1. The program defines `DirEntry`, a struct with three fields: `Name`, `Purpose`, and `Contents`.
2. `repoDirs` is a slice holding one `DirEntry` for each architectural subsystem in the v3 repository.
3. `Lookup("curriculum/")` iterates over `repoDirs` and compares each entry's `Name` field against the argument, returning a pointer or `nil`.
4. `main()` prints a formatted table of all subsystems and their responsibilities.
5. Verification tests execute `Lookup` with canonical subsystem paths and assert that all returned entries contain non-empty metadata.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Hand-editing manifests without running validators | Assuming JSON changes are valid without checking schema constraints | Run `go run ./internal/tools/curriculum validate-all` after any edit |
| Confusing `curriculum/` manifests with lesson READMEs | Manifests declare the schema and graph; lesson READMEs contain the instructional prose | Keep metadata in `path.*.json` and human-facing teaching in the module folders |
| Looking for obsolete pre-v3 directories | Historical v2 folders were consolidated into `curriculum/` and `internal/tools/` in v3 | Reference the canonical v3 structure documented in this lesson |
| Forgetting that `Lookup` returns `nil` | Dereferencing a nil pointer causes an immediate runtime panic in Go | Always check `if entry != nil` before accessing fields |

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

The inline code example is for reading and understanding. To verify your understanding, complete the practice task above and check your answers against the description. You can also copy the inline code into a local `.go` file and run `go run .` and `go test .` in that directory to experiment with the output.

## Review questions

1. What does the `Lookup` function return when given `"metadata"` (without a trailing slash)?
2. Why does the curriculum separate metadata from curriculum content instead of combining them into one directory?
3. What would break if you deleted `curriculum/modules/00-orientation/` but left the metadata JSON unchanged?
4. In the `DirEntry` struct, which field stores the list of files inside the directory?
5. How would you change the Go example to use a map instead of a slice for constant-time lookups?

## NEXT UP

[Lesson 02: What zero magic means](../02-what-zero-magic-means/README.md) — where you learn why every concept in this curriculum is explained explicitly, with no hidden surprises.
