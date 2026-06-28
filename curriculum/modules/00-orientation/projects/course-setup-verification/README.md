# Course Setup Verification

## Goal

Verify that your development environment is ready for this curriculum. You will check that Go, Git, and a text editor are installed on your system using your terminal.

## Prerequisites

- Lesson 01 through Lesson 10 of Module 00 (Orientation)
- Go installed (you should have done this during Lesson 01)
- A terminal open
- A text editor

## Task

Open your terminal and run the following checks:

1. **Go**: Run `go version`. You should see output like `go version go1.22.0 windows/amd64`.
2. **Git**: Run `git version`. You should see output like `git version 2.43.0`.
3. **Text editor**: Run the command for your editor — `code --version` (VS Code), `vim --version`, or `nano --version`.

Record the output of each command in a file called `setup.txt`:

```bash
go version > setup.txt
git version >> setup.txt
code --version >> setup.txt   # or vim/nano
```

## Verification

Run each command below and confirm the output matches what you expect:

```bash
go version
git version
code --version   # or your editor
```

The expected output for Go is a version string (e.g., `go version go1.22.0`). If you see `command not found` or similar, the tool is not installed or not in your PATH.

When you are done, update your Portfolio Plan (from Lesson 10) to include this project with status "complete" and the evidence file `setup.txt`.

## NEXT UP

[Assessment: Course Setup](../assessments/course-setup/README.md) — confirm your environment passes the verification.
