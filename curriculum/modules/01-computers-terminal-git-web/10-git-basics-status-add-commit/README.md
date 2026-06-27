# Git basics: status, add, commit

## Learning objective

Understand and apply Git basics: status, add, commit in the context of professional Go software engineering.

## Why this matters

Developers need a repeatable workflow to save work-in-progress, review changes before committing, and build a clean history. The status-add-commit cycle provides a safe, structured way to record progress.

## Mental model

The Git workflow is three baskets: the working directory (your current edits), the staging area (what you plan to commit), and the commit history (what has been saved). Changes flow left to right: working → staged → committed.

## Core idea

Without staging, every file change would be automatically included in the next commit — making partial commits, code review, and selective history impossible. The status-add-commit workflow gives developers fine-grained control over what enters the permanent record.

## Under the hood

The Git index (staging area) is a binary file at .git/index. It contains a sorted list of paths with their metadata (mode, SHA-1 hash, stage number). git add updates the hash in the index. git commit reads the index, builds a tree object from it, creates a commit object pointing to that tree, and writes the commit hash into the current branch ref.

## How Go uses it

Go developers use Git for every project: 'git init' creates a module, 'git add' and 'git commit' save work, 'git push' publishes code. Go's 'go mod vendor' creates a vendor directory that should be committed. GitHub Actions CI runs tests on every git push.

## Go example

```go
package main

import (
	"fmt"
	"sort"
)

// Repo simulates a basic Git repository with working directory, staging, and history.
type Repo struct {
	Working map[string]string
	staging map[string]string
	History []Commit
}

// Commit represents a saved snapshot.
type Commit struct {
	Message string
	Files   map[string]string
}

// NewRepo creates an empty repository.
func NewRepo() *Repo {
	return &Repo{
		Working: make(map[string]string),
		staging: make(map[string]string),
		History: []Commit{},
	}
}

// StageNames returns sorted names of staged files.
func (r *Repo) StageNames() []string {
	names := make([]string, 0, len(r.staging))
	for name := range r.staging {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Stage moves a file from working directory to staging area.
func (r *Repo) Stage(name string) error {
	content, ok := r.Working[name]
	if !ok {
		return fmt.Errorf("file not found: %s", name)
	}
	r.staging[name] = content
	return nil
}

// Unstage removes a file from staging area.
func (r *Repo) Unstage(name string) {
	delete(r.staging, name)
}

// Commit saves all staged files to history.
func (r *Repo) Commit(msg string) error {
	if len(r.staging) == 0 {
		return fmt.Errorf("nothing staged")
	}
	files := make(map[string]string)
	for k, v := range r.staging {
		files[k] = v
	}
	r.History = append(r.History, Commit{Message: msg, Files: files})
	r.staging = make(map[string]string)
	return nil
}

func main() {
	repo := NewRepo()
	repo.Working["hello.txt"] = "Hello, Git!"

	fmt.Println("Staged:", repo.StageNames())

	repo.Stage("hello.txt")
	fmt.Println("After stage:", repo.StageNames())

	repo.Commit("feat: add hello.txt")
	fmt.Printf("History: %d commits\n", len(repo.History))

	repo.Working["main.go"] = "package main"
	repo.Stage("main.go")
	repo.Commit("feat: add main.go")
	fmt.Printf("History: %d commits\n", len(repo.History))
}
```

## Step-by-step execution

1. Run 'git status' to see the current state — modified tracked files, staged files, and untracked files are listed separately.
2. Run 'git add <file>' or 'git add .' to stage changes — Git hashes file contents and updates the index.
3. Run 'git diff --staged' to review exactly what will be included in the next commit.
4. Run 'git commit -m "type: description"' to create a commit from the staged snapshot.
5. Run 'git log --oneline' to verify the commit appears at the top of the history.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Running git add . blindly without checking what files are being staged | Convenience over caution — you accidentally commit build artifacts or secrets | Use 'git status' and 'git add <file>' selectively, or use .gitignore |
| Writing meaningless commit messages like 'fix' or 'update' | Rush to save work without considering future readers | Use conventional commits: 'feat: add login endpoint', 'fix: handle nil pointer in parser' |
| Committing large binary files or dependency directories | Not understanding that Git stores every version permanently | Use .gitignore for binaries and vendored deps, or use Git LFS for large files |

## Debugging walkthrough

### Scenario: A developer runs 'git commit' without configuring user.name and user.email and gets an error.

**Cause:** Git requires identity information for every commit; it reads from git config user.name and user.email.

**Fix:** Set global config: 'git config --global user.name "Your Name"' and 'git config --global user.email "you@example.com"'.

### Scenario: A developer stages a file, edits it again, then commits — the second edit is not included in the commit.

**Cause:** git add captures the file's content at the moment of staging; subsequent edits are unstaged modifications.

**Fix:** Always run 'git status' before committing, and re-stage changed files or use 'git commit -a' for tracked files.

## Production notes

The status-add-commit cycle is performed dozens of times daily by every professional developer. It is the fundamental rhythm of Git-based development, used in every codebase from small personal projects to monorepos at Google and Microsoft.

## Performance implications

- git add and git commit are local operations — no network, no latency. A typical commit takes under 100ms.
- Large repositories (100k+ files) slow down git status — use .gitignore to exclude build artifacts and dependencies.
- Git's index file scales with the number of tracked files — performance degrades linearly.

## Practice task

The learner must create a repository, modify a file, stage it, make another modification, use git diff to see unstaged changes, re-stage, commit with a proper message, and confirm the commit with git log.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/10-git-basics-status-add-commit/
```

## Review questions

1. Learner must use git status to identify which changes are staged, unstaged, and untracked in a given scenario.
2. Learner must write a meaningful commit message following the conventional commit format (type: description).
3. Learner must use git diff --staged to review what will be committed before running git commit.

## NEXT UP

Lesson 11: [Branching and merging](../11-branching-and-merging/README.md)
