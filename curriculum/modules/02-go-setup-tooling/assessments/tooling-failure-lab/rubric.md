# Tooling Failure Lab — Rubric

## Criteria

### 1. Toolchain and environment fluency (30%)

| Level | Score | Description |
|-------|-------|-------------|
| Excellent | 5 | Deeply understands `GOROOT`, `GOPATH`, module cache mechanics, `go.mod`/`go.sum` lifecycle, and cross-compilation (`GOOS`/`GOARCH`). Configures tools without trial-and-error. |
| Proficient | 3 | Accurately manages standard Go toolchain commands (`build`, `test`, `vet`, `run`, `mod tidy`); minor uncertainty regarding module cache internals or path overrides. |
| Developing | 1 | Relies on guessing paths or commands; struggles to distinguish between module mode and legacy GOPATH mode. |
| Not attempted | 0 | No toolchain configuration demonstrated. |

### 2. Diagnostic and debugging competency (35%)

| Level | Score | Description |
|-------|-------|-------------|
| Excellent | 5 | Quickly and methodically diagnoses compiler errors, syntax issues, `go vet` warnings, missing dependency resolution, and panic tracebacks. Pinpoints exact file, line, and root cause. |
| Proficient | 3 | Resolves standard compilation errors and missing modules, but takes multiple attempts or misinterprets runtime panic call stacks. |
| Developing | 1 | Cannot effectively read error messages or panic traces; changes code haphazardly without diagnosing the root cause. |
| Not attempted | 0 | No debugging or failure remediation demonstrated. |

### 3. Verification and proof of execution (20%)

| Level | Score | Description |
|-------|-------|-------------|
| Excellent | 5 | All verification commands pass cleanly (`go test`, `go vet`, `go build`). Provides clean terminal capture proving correct environment configuration and dependency verification. |
| Proficient | 3 | Primary verification checks pass; minor warnings or incomplete terminal environment capture. |
| Developing | 1 | Multiple verification checks fail; test suite or compiler reports outstanding errors. |
| Not attempted | 0 | No verification output provided. |

### 4. Conceptual clarity and trade-off analysis (15%)

| Level | Score | Description |
|-------|-------|-------------|
| Excellent | 5 | Clearly articulates the compiler pipeline (lexing -> parsing -> type checking -> SSA -> codegen/link), the role of `go.sum` in supply chain security, and LSP architecture with `gopls`. |
| Proficient | 3 | Can explain core compiler and module concepts clearly with minor omissions in internal mechanics. |
| Developing | 1 | Can repeat buzzwords but cannot explain why Go requires `package main`, how `go vet` differs from the compiler, or what `gopls` does. |
| Not attempted | 0 | No conceptual explanation provided. |

## Scoring guide

| Total points | Grade |
|--------------|-------|
| 18–20 | Exceeds expectations |
| 14–17 | Meets expectations (Passing) |
| 8–13 | Needs improvement |
| 0–7 | Incomplete — retake recommended |
