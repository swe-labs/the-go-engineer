# Course Setup Verification

## Goal

Verify that your development environment is ready for this curriculum. You will check that Go, Git, and a text editor are installed on your system using your terminal, and confirm that the automated quality toolchain runs successfully.

## Prerequisites

- Lesson 01 through Lesson 10 of Module 00 (Orientation)
- Go 1.24+ (or Go 1.26+) installed
- A terminal open
- A text editor (VS Code, GoLand, Neovim, etc.)

## Task

Open your terminal and run the following environment checks:

1. **Go**: Run `go version`. You should see output indicating Go 1.24 or higher, such as `go version go1.26.0 windows/amd64` (or `darwin/arm64`, `linux/amd64`).
2. **Git**: Run `git version`. You should see output like `git version 2.43.0` or higher.
3. **Text editor**: Run the command for your editor — `code --version` (VS Code), `vim --version`, or `nano --version`.
4. **Toolchain verification**: Run the automated quality gate runner from repository root:
   ```bash
   go run ./scripts/validate_curriculum.go
   ```
   Confirm that all 17 quality gates pass.

Record the output of each command in a file called `setup.txt`:

```bash
go version > setup.txt
git version >> setup.txt
code --version >> setup.txt   # or vim/nano
go run ./scripts/validate_curriculum.go >> setup.txt
```

## Verification

Run each command below and confirm the output matches what you expect:

```bash
go version
git version
code --version   # or your editor
go run ./scripts/validate_curriculum.go
```

The expected output for Go is a valid version string (e.g., `go version go1.26.0 ...`). If you see `command not found` or similar, the tool is not installed or not in your system `PATH`.

When you are done, update your Portfolio Plan (from Lesson 10) to include this project with status "complete" and the evidence file `setup.txt`.

## NEXT UP

[Assessment: Course Setup](../../assessments/course-setup/README.md) — confirm your environment passes the verification.
