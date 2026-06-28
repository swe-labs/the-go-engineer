# Module 02 Assessment — Go Setup and Tooling

## Overview

This assessment verifies that you can independently apply everything from Module 02. You will demonstrate that your Go toolchain is correctly installed, configured, and usable — and that you can diagnose and fix common tooling failures.

## Evidence required

Submit the following artifacts:

1. **Tooling Failure Lab project** — your completed `_starter/` directory with all failures fixed, all verification commands passing, and `TROUBLESHOOTING.md` documenting each failure
2. **Tooling configuration output** — the output of `go version`, `go env GOROOT GOPATH GOOS GOARCH`, and `gopls version`
3. **Module verification** — the output of `go list -m all` and `go mod verify` from the curriculum root

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| Conceptual clarity | 30% | Clearly explains Go installation, workspace setup, build toolchain, and debugging tools — including mechanics, purpose, and common mistakes | Explains basics with some gaps | Cannot explain core concepts |
| Implementation | 35% | Complete, clean fix of all tooling failures with proper Go idioms | Partial fix with minor gaps | No working fix |
| Proof | 25% | Strong proof with tests or verification covering Go installation and workspace setup, including edge cases | Some proof provided but lacks depth | No proof of understanding |
| Communication | 10% | Professional explanation linking Go tooling to broader development workflow and trade-offs | Basic explanation with some clarity | No explanation provided |

**Passing score**: 80%

## Review questions

Be prepared to answer these questions in a review session:

1. What does `go build` actually do? Explain the compilation pipeline from source to binary.
2. Why does Go enforce that import cycles are forbidden? What problem does the module cache solve?
3. Show me a compiler error you have seen. Explain exactly why the compiler rejected the code.
4. What is the difference between `go run` and `go build`? When would you use each?
5. What does `gofmt` guarantee about your code? Why does Go have a single official format?
6. What kind of bugs does `go vet` catch that the compiler does not? Give a concrete example.
7. How do you read a panic traceback? Where do you look first and why?
8. What is `gopls` and what does it provide that `go build` alone does not?
9. What is `go.mod` and what problem does the module system solve?
10. Walk through your troubleshooting document: what tooling failures did you encounter and how did you diagnose each one?

## Submission process

1. Complete the Tooling Failure Lab project in `projects/tooling-failure-lab/`.
2. Ensure all verification commands pass.
3. Write `TROUBLESHOOTING.md` documenting each failure.
4. Capture the output of `go version`, `go env GOROOT GOPATH GOOS GOARCH`, `gopls version`, `go list -m all`, and `go mod verify`.
5. Submit all artifacts for review.

## Retake policy

If you score below 80%, review the lessons corresponding to your rubric gaps, fix the issues, and resubmit with a short correction note explaining what you changed and why.

## NEXT UP

[Module 03 — Programming Fundamentals with Go](../../../03-programming-fundamentals/README.md)
