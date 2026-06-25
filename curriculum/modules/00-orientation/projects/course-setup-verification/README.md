# Course Setup Verification

## Goal

Verify that your development environment is ready for this curriculum. You will write a Go program that checks whether Go, Git, and a text editor are installed on your system. This is the first project you will complete and the first piece of evidence in your portfolio.

## Prerequisites

- Lesson 01 through Lesson 10 of Module 00 (Orientation)
- Go installed (you should have done this during Lesson 01)
- A terminal open
- A text editor

## Task

Write a Go program that does the following:

1. Define a function `checkTool(name string, path string) string` that uses `os/exec.LookPath` to check whether a tool is installed. If the tool is found, return a string starting with "OK". If not, return a string starting with "NOT FOUND".
2. In `main()`, call `checkTool` for `go` and `git`.
3. Add a check for a text editor. Popular choices are `code` (VS Code), `vim`, or `nano`. Pick the one you use.
4. Print all results to the terminal.
5. Write table-driven tests for `checkTool` that cover at least:
   - A known installed tool (go, git)
   - A tool that definitely does not exist on any system
   - An empty path

## Starter

The `_starter/main.go` file contains the scaffolding:

```go
package main

import (
	"fmt"
	"os/exec"
)

func checkTool(name string, path string) string {
	_, err := exec.LookPath(path)
	if err != nil {
		return fmt.Sprintf("NOT FOUND: %s is not installed or not in PATH", name)
	}
	return fmt.Sprintf("OK: %s is installed", name)
}

func main() {
	fmt.Println("Course Setup Verification")
	fmt.Println(checkTool("Go", "go"))
	fmt.Println(checkTool("Git", "git"))

	// TODO: Check for a text editor (e.g., "code" for VS Code, "vim", "nano", etc.)
}
```

The `checkTool` function is implemented for you. The `main()` function calls it for Go and Git but leaves the editor check as a TODO. Your job is to add the editor check and write the tests.

The `_starter/main_test.go` file contains the test skeleton:

```go
package main

import (
	"strings"
	"testing"
)

func TestCheckTool(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		path     string
		wantOK   bool
	}{
		{name: "go is installed", toolName: "Go", path: "go", wantOK: true},
		{name: "nonexistent tool", toolName: "Nonexistent", path: "this-tool-does-not-exist-12345", wantOK: false},
		{name: "git is installed", toolName: "Git", path: "git", wantOK: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkTool(tc.toolName, tc.path)
			if tc.wantOK && !strings.HasPrefix(result, "OK") {
				t.Errorf("expected OK, got: %s", result)
			}
			if !tc.wantOK && !strings.HasPrefix(result, "NOT FOUND") {
				t.Errorf("expected NOT FOUND, got: %s", result)
			}
		})
	}
}
```

## Solution

The `_solution/` directory contains a complete implementation and passing tests. Use it to check your work, but try to complete the project on your own first.

The solution adds checks for multiple editors (`code`, `vim`, `nano`) and prints a summary message. The tests include an extra case for an empty path.

## Submission

To verify your project:

```bash
go run ./curriculum/modules/00-orientation/projects/course-setup-verification/_starter
go test ./curriculum/modules/00-orientation/projects/course-setup-verification/_starter
```

Expected `go run` output (your results will depend on what is installed):

```
Course Setup Verification
=========================
OK: Go is installed
OK: Git is installed
NOT FOUND: code is not installed or not in PATH
```

Expected `go test` output:

```
ok  	github.com/rasel9t6/the-go-engineer/curriculum/modules/00-orientation/projects/course-setup-verification/_starter
```

When you are done, update your Portfolio Plan (from Lesson 10) to include this project with status "complete" and the three evidence items: `main.go`, `main_test.go`, and this README.
