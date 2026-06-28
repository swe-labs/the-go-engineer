# How to use the roadmap

## Learning objective

Understand how the curriculum roadmap works as a dependency graph, how prerequisites determine valid learning paths, and how to validate your own progression through the modules. You will be able to read the module dependency graph, identify the correct order, and explain why some modules must come before others.

## Why this matters

A curriculum without structure is a playlist. You watch videos in any order and hope it works out. A curriculum with a dependency graph is a building. You cannot build the second floor before the first floor. If you skip Orientation (module 00) and jump straight to Testing (module 03), you will not have Go installed, you will not know the terminal, and you will not understand the testing tools. The roadmap tells you what depends on what so you never find yourself in a module that assumes knowledge you do not have.

## Mental model

Think of each module as a node in a graph. An arrow points from a module to another module that depends on it. Module 00 (Orientation) points to Module 01 (Computers, Terminal, Git, Web), which points to Module 02 (Go Basics), which points to Module 03 (Testing in Go). You cannot arrive at Module 03 without passing through 00, 01, and 02. The graph is acyclic, meaning there are no loops. Every path has a start and an end.

## Core idea

A roadmap is a dependency graph, not a playlist. Each module lists its prerequisites. A learning path is valid if every prerequisite for each module has been completed before that module is attempted. If you attempt a module whose prerequisites are not met, the path is invalid. You must go back and complete the prerequisite first.

## Under the hood

Each module in the curriculum has an ID, a title, and a list of prerequisite module IDs. Module 00 is the root: it has no prerequisites. Every other module depends on module 00 directly or indirectly. The validation algorithm keeps a set of completed modules. As it walks the path, it checks that every prerequisite for the current module exists in the completed set. If a prerequisite is missing, the path is invalid and the algorithm reports which step failed and why.

## How Go uses it

The Go program defines a `Module` struct with `ID`, `Title`, and `Prerequisites` (a string slice). The `validatePath` function takes a slice of `Module` values and walks them in order. It maintains a `completed` set (a `map[string]bool`). For each module, it checks all prerequisites. If any prerequisite is not in the completed set, it returns a `ValidationResult` with `Valid: false` and a descriptive error. If the entire path is valid, it returns `Valid: true`.

## Go example

```go
package main

import "fmt"

type Module struct {
	ID            string
	Title         string
	Prerequisites []string
}

type ValidationResult struct {
	Valid bool
	Step  int
	Path  []string
	Error string
}

func validatePath(path []Module) ValidationResult {
	completed := make(map[string]bool)
	for i, m := range path {
		for _, prereq := range m.Prerequisites {
			if !completed[prereq] {
				pathNames := make([]string, len(path))
				for j, p := range path {
					pathNames[j] = p.ID
				}
				return ValidationResult{
					Valid: false,
					Step:  i,
					Path:  pathNames,
					Error: fmt.Sprintf("module %q requires %q, but %q has not been completed",
						m.ID, prereq, prereq),
				}
			}
		}
		completed[m.ID] = true
	}
	pathNames := make([]string, len(path))
	for i, p := range path {
		pathNames[i] = p.ID
	}
	return ValidationResult{Valid: true, Path: pathNames}
}

func main() {
	modules := map[string]Module{
		"00": {ID: "00", Title: "Orientation"},
		"01": {ID: "01", Title: "Computers, Terminal, Git, Web",
			Prerequisites: []string{"00"}},
		"02": {ID: "02", Title: "Go Basics",
			Prerequisites: []string{"00", "01"}},
		"03": {ID: "03", Title: "Testing in Go",
			Prerequisites: []string{"00", "01", "02"}},
	}

	validPath := []Module{modules["00"], modules["01"], modules["02"], modules["03"]}
	invalidPath := []Module{modules["02"], modules["01"], modules["00"], modules["03"]}

	fmt.Println("Valid path (00->01->02->03):", validatePath(validPath).Valid)

	result := validatePath(invalidPath)
	fmt.Println("Invalid path (02->01->00->03):", result.Valid)
	fmt.Println("Error:", result.Error)
}
```

## Step-by-step execution

1. The program creates four modules: 00 (no prerequisites), 01 (needs 00), 02 (needs 00, 01), 03 (needs 00, 01, 02).
2. Two paths are constructed. The valid path goes 00, 01, 02, 03 in order. The invalid path starts with 02.
3. `validatePath(validPath)` starts with an empty `completed` set. It adds 00 first. Then 01 checks that 00 is completed (yes). Then 02 checks 00 and 01 (both yes). Then 03 checks 00, 01, 02 (all yes). Returns Valid: true.
4. `validatePath(invalidPath)` starts with module 02. It checks 02's prerequisites: 00 and 01. Neither is in `completed`. It immediately returns Valid: false with an error: "module '02' requires '00', but '00' has not been completed".
5. The output shows the valid path passes and the invalid path fails with a clear error message.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Skipping module 00 | It is labeled "Orientation" so it sounds optional | Module 00 is a hard prerequisite for every other module. Do not skip it. |
| Confusing suggested order with dependency order | Some modules can be taken in parallel if you meet their prerequisites | Check the Prerequisites field in each module. That is the hard constraint. Everything else is suggestion. |
| Thinking you can "catch up" later | If module 01 requires module 00 skills, you will struggle in 01 without 00 | Do not start a module until you have completed its prerequisites. Struggling does not mean learning faster. |
| Ignoring indirect dependencies | Module 03 requires 00 and 01 and 02 | You must complete all indirect prerequisites too. The validatePath function checks all prerequisites recursively. |
| Trying to learn everything linearly | Some modules may not depend on each other and can be done in any order | Use the dependency graph to find parallel paths. Work on multiple independent modules to stay engaged. |

## Debugging walkthrough

Suppose you define a path that you think is valid but the program says it is invalid.

1. Print the `completed` map after each step to trace what is being tracked: `fmt.Printf("Completed after %s: %v\n", m.ID, completed)`.
2. Check that the `Prerequisites` field uses the exact string IDs. Module 01 must list "00", not "0" or "orientation".
3. If you are missing a module definition, the map lookup `modules["00"]` will return a zero-value `Module` with an empty ID and nil prerequisites. This might accidentally pass validation. Verify with a `_, ok := modules["00"]` check.
4. If the error message says module X requires Y but you thought Y was completed in an earlier step, verify the order of modules in the path slice. The first element is step 0.
5. If the path is empty, `validatePath` returns `Valid: true` (an empty path has no violations). If this surprises you, add an explicit check for empty paths.

If the program panics with a nil map, make sure you used `make(map[string]bool)` to initialize `completed` inside `validatePath`.

## Production notes

Real engineering roadmaps use the same dependency-graph principle. Build systems like Make, Bazel, and Gradle track file dependencies. Deployment pipelines track environment dependencies (staging before production). Microservice architectures track API dependencies. The same idea appears everywhere: know what depends on what, validate the order, and fail fast when a prerequisite is missing.

When planning your learning schedule, use the dependency graph to estimate time. Each module takes roughly 1-2 weeks. Modules with many prerequisites take longer because you must complete the prerequisites first. Module 03 (Testing in Go) requires 00, 01, and 02, so it is effectively week 4-6 of the curriculum.

## Performance implications

The validation algorithm runs in O(p * m) where p is the average number of prerequisites per module and m is the number of modules in the path. For this curriculum with fewer than 15 modules and at most 4 prerequisites each, performance is irrelevant. The important metric is correctness: the algorithm must catch every missing prerequisite and report it with a clear message.

## Practice task

Extend the program to add two more modules: Module 04 (Databases) with prerequisites ["00", "01", "02"] and Module 05 (Web Servers) with prerequisites ["00", "01", "02", "03", "04"]. Create a valid path through all six modules and an invalid path that violates the dependencies. Run validation on both.

## Tests / verification

1. Copy the inline Go code into a local `.go` file and run `go run .` to see the valid and invalid path results.
2. Create a path with a single module that has no prerequisites — confirm `validatePath` returns valid.
3. Swap two modules in the valid path so they are out of order — confirm the validation fails with a clear error.

## Review questions

1. What is the difference between a dependency graph and a playlist?
2. Why does module 02 require both module 00 and module 01?
3. What does the program return when a prerequisite is missing?
4. How does the `completed` map prevent duplicate prerequisites from being checked twice?
5. If you have completed modules 00 and 01, can you start module 03? Why or why not?

## NEXT UP

[How to build a portfolio from this curriculum](../10-how-to-build-a-portfolio-from-this-curriculum/README.md) — Turn your curriculum projects into professional portfolio evidence.
