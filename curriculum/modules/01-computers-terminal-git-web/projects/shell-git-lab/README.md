# Shell + Git Lab

## Goal

Practice shell navigation and Git commands by building and navigating a simulated filesystem and managing code changes through a complete Git workflow.

## Learning objectives

- Navigate the filesystem using shell commands (cd, ls, pwd, mkdir, rm)
- Initialize and manage a Git repository (init, add, commit, branch, merge)
- Resolve merge conflicts
- Understand the shell and Git mental models from lessons 1-15

## Tasks

### 1. Set up and navigate the terminal

Open your terminal, verify shell access, navigate directories, create and delete files, and configure basic shell aliases.

### 2. Initialize and manage a Git repository

Initialize a Git repository, stage and commit changes, create branches, merge, and resolve a simulated merge conflict.

### 3. Complete the Go simulation

The `_starter` directory contains a Go program that simulates checking shell and Git knowledge. Complete the TODO markers to make all tests pass.

## Starter / solution workflow

1. Start by reading this README and understanding the project structure.
2. Open the `_starter` directory and read `main.go` and `main_test.go`.
3. Complete the TODO markers in `_starter/main.go` — do not modify the test file.
4. Run `go test -v ./_starter/` to verify your implementation.
5. Once all tests pass, compare your solution with `_solution/main.go` to see alternative approaches.

## Verification

```bash
# Run tests
go test -v ./curriculum/modules/01-computers-terminal-git-web/projects/shell-git-lab/_starter/

# Build the project
go build ./curriculum/modules/01-computers-terminal-git-web/projects/shell-git-lab/_starter/
```

## Deliverables

- Working implementation in `_starter/main.go`
- All tests passing
- Explanation of your design decisions
