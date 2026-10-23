# Module 06 Checkpoint — Testing, Debugging, and Refactoring

## Purpose

This checkpoint verifies that you understand Go's testing framework, debugging tools, and refactoring discipline before advancing to CLI applications and file I/O. You should be able to write table-driven tests, use the race detector, navigate Delve, and apply characterization testing methodology.

## Scope

This checkpoint covers:

- `testing.T` basics: test functions, assertions, and failure reporting
- Table-driven tests and subtests with `t.Run`
- Test coverage analysis with `go test -cover` and `-coverprofile`
- Benchmarking with `testing.B`
- The Go race detector (`go test -race`)
- `sync.Mutex` and `sync.RWMutex` for concurrent safety
- Delve debugger: setting breakpoints, stepping, inspecting variables
- Error tracing and stack trace reading
- Characterization testing methodology
- Safe refactoring under test coverage
- Code organization: separation of concerns

## Time expectation

Plan for 45 to 60 minutes.

## Allowed resources

You may use:

- Module 06 lesson READMEs
- Your terminal, editor, and Delve
- Your own notes

You may not copy from the answer key before attempting the questions.

## Evidence required

Submit:

1. Answers to `questions.md`
2. A table-driven test you wrote for a simple function (at least 4 test cases including an error case)
3. A brief paragraph explaining the difference between a race condition and a data race in Go

## Rubric

| Dimension | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Testing mechanics** | Writes clean table-driven tests with subtests; understands `t.Fatal` vs `t.Error`; can read coverage output | Writes basic tests but not table-driven | Cannot write a test function |
| **Race detection** | Explains what `go test -race` instruments; understands false-negative rate; uses mutex correctly | Knows `-race` flag but cannot explain what it checks | Unaware of race detection |
| **Debugging** | Uses Delve to set breakpoints, step through code, and inspect variables; reads panic stack traces | Knows Delve exists but relies on print debugging | Cannot use any debugger |
| **Refactoring discipline** | Explains characterization testing workflow; understands why tests must precede refactoring | Knows refactoring exists but no systematic method | Refactors without tests |

**Passing threshold**: 3/5 average across all dimensions

## Retake policy

If you miss a category, return to the relevant lesson, revise your answer, and write one sentence explaining what changed in your understanding.

## NEXT UP

[Module 07 — CLI, Files, JSON, and Config](../../../07-cli-files-json-config/README.md)
