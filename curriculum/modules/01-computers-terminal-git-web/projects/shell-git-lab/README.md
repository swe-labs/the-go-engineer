# Shell + Git Lab

## Goal

Practice shell navigation and Git commands by building a directory structure and managing code changes through a complete Git workflow — all in your terminal. No Go code required.

## Learning objectives

- Navigate the filesystem using shell commands (`cd`, `ls`, `pwd`, `mkdir`, `rm`)
- Initialize and manage a Git repository (`init`, `add`, `commit`, `branch`, `merge`)
- Resolve merge conflicts
- Understand shell and Git mental models from lessons 1–15

## Tasks

### 1. Shell fluency

Open your terminal and complete the following:

- Create a directory tree: `project/src/`, `project/docs/`, `project/tests/`
- Create files: `project/README.md`, `project/src/main.go`, `project/docs/design.md`
- Use `ls`, `pwd`, `tree` (or `Get-ChildItem` on PowerShell) to inspect the structure
- Copy, move, and delete files using `cp`, `mv`, `rm`
- Use `>` to write output to a file and `>>` to append

### 2. Git workflow

Inside `project/` initialize a Git repo and complete a full workflow:

- `git init`
- Create an initial commit with `README.md`
- Create a `feature` branch, add `src/main.go`, commit
- Switch back to `main`, add `docs/design.md`, commit
- Merge `feature` into `main` (fast-forward or three-way)
- Create a deliberate merge conflict: edit the same line on two branches, attempt merge, resolve it
- Run `git log --oneline --graph` and verify the DAG

### 3. Conceptual understanding

Write a brief explanation (100–200 words) answering:

- What is the difference between the working directory, the staging area, and the commit history?
- Why does Git use SHA-1 hashes for commits?
- What happens when you merge two branches with diverged histories?

## Verification

```powershell
# PowerShell (Windows / pwsh)
powershell -ExecutionPolicy Bypass -File ./verify.ps1
```

```bash
# Bash (Linux / macOS / Git Bash on Windows)
./verify.sh
```

On Windows, you can run `verify.ps1` directly from PowerShell. On Linux/macOS, or if you use Git Bash on Windows, run `verify.sh` instead.

## Deliverables

- Working Git repository with branches, commits, and a resolved merge conflict
- All verification checks passing (14/14)
- Written explanation of the three prompts above
