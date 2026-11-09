# Tooling Failure Lab — Rubric

This rubric scores the 14-question knowledge assessment (18 points: 10 multiple-choice
at 1 point, 4 short-answer at 2 points). Question-to-dimension mapping:

| Dimension | Questions | Points | Approx. weight |
|---|---|---|---|
| Toolchain & module system (GOROOT/GOPATH, tidy, cross-compile, dependency errors) | 1, 7, 9, 12 | 5 | ~28% |
| Build, run & static diagnostics (run vs build, vet, gofmt, pipeline) | 3, 4, 6, 10, 11 | 6 | ~33% |
| Runtime failures & traceback triage | 5, 13 | 3 | ~17% |
| Module integrity & supply chain (main entrypoint, go.mod vs go.sum) | 2, 8, 14 | 4 | ~22% |

## Dimension standards

### Toolchain & module system (5 points)

| Level | Standard |
|-------|----------|
| Excellent (5) | Correctly distinguishes GOROOT vs GOPATH, states exactly what `go mod tidy` synchronizes, gives the `GOOS`/`GOARCH` cross-compilation recipe, and names distinct root causes for missing-package build errors with their fix commands. |
| Proficient (3-4) | Commands correct; conflates module cache internals or names fewer than three distinct root causes. |
| Developing (1-2) | Guesses paths or confuses module mode with legacy GOPATH mode. |
| Not demonstrated (0) | Toolchain answers incorrect or absent. |

### Build, run & static diagnostics (6 points)

| Level | Standard |
|-------|----------|
| Excellent (6) | Precisely contrasts the `go run` temp-binary lifecycle with `go build`'s persistent artifact; names the vet analyzer classes (printf, unreachable) the compiler accepts; explains gofmt's AST-based determinism; walks the compilation pipeline stages in order. |
| Proficient (4-5) | Lifecycle and pipeline correct with minor omissions (e.g. skips the SSA stage). |
| Developing (2-3) | Believes `go run` interprets code or that vet finds compile-time errors. |
| Not demonstrated (0-1) | Build/run answers incorrect or absent. |

### Runtime failures & traceback triage (3 points)

| Level | Standard |
|-------|----------|
| Excellent (3) | Reads the panic value first, skips `runtime.gopanic` frames, and pinpoints the first application frame's file and line. |
| Proficient (2) | Finds the crash site but cannot distinguish runtime frames from application frames. |
| Not demonstrated (0-1) | Traceback reading incorrect or absent. |

### Module integrity & supply chain (4 points)

| Level | Standard |
|-------|----------|
| Excellent (4) | States the exact `package main` / `func main()` entrypoint contract and explains `go.sum`'s SHA-256 checksum verification against upstream tampering. |
| Proficient (2-3) | Entry-point contract correct; supply-chain explanation generic. |
| Not demonstrated (0-1) | Module integrity answers incorrect or absent. |

## Scoring guide

| Total points | Result |
|--------------|--------|
| 16-18 | Exceeds expectations |
| 13-15 | Meets expectations — module passed (72%+) |
| 9-12 | Needs improvement — retake recommended |
| 0-8 | Incomplete — retake required |

A passing score is at least 13 of 18 points (72%). The Tooling Failure Lab project
itself is graded separately by the project rubric in `projects/tooling-failure-lab/README.md`
(passing 80% with all failures fixed).
