# Testing, Debugging, and Refactoring

## Goal

Move proof, diagnosis, and safe change before serious backend work.

## Summary

Builds professional feedback loops through unit tests, table tests, debugging, refactoring, coverage, and safe change management.

## Prerequisites

- [Module 05: Types, Interfaces, Packages, and Modules](../05-types-interfaces-packages-modules/README.md) — complete that module's completion checkpoint before starting this one.

## Learning Sessions & Lessons

### Session 1: Unit Testing & Feedback Loops
- `core-06-01` Why testing exists
- `core-06-02` Unit testing
- `core-06-03` Table-driven tests
- `core-06-04` Subtests
- `core-06-05` t.Cleanup
- `core-06-06` Test fixtures
- `core-06-07` Golden files
- `core-06-08` Testable design with io.Writer
- `core-06-09` Fakes before mocks
- `core-06-10` Mocking tradeoffs

### Session 2: Diagnosis Mindset & Stack Traces
- `core-06-11` Debugging mindset
- `core-06-12` Reading stack traces
- `core-06-13` Debugging panics

### Session 3: Interactive Delve Debugging Workflow (Consolidated Workshop)
*A unified four-part debugging workflow resolving a realistic bug under Delve:*
- `core-06-14` **Delve basics** *(Step 1: Launching `dlv debug`, attaching to running processes, navigation)*
- `core-06-15` **Breakpoints** *(Step 2: Line breakpoints, condition triggers, watchpoints)*
- `core-06-16` **Stepping through code** *(Step 3: Execution control: next, step, stepout, frame navigation)*
- `core-06-17` **Inspecting variables** *(Step 4: Runtime state evaluation: print, display, locals, goroutine state)*
- `core-06-18` Debugging tests

### Session 4: Safe Refactoring & Verification
- `core-06-19` Refactoring safely
- `core-06-20` Fuzz testing
- `core-06-21` Race detector preview

## Projects

- [`projects/bug-hunt-pack/README.md`](./projects/bug-hunt-pack/README.md) — **Bug Hunt Pack**
- [`projects/refactor-under-tests/README.md`](./projects/refactor-under-tests/README.md) — **Refactor Under Tests**

## Assessments

- [`assessments/bug-hunt-pack/README.md`](./assessments/bug-hunt-pack/README.md) — **Bug Hunt Pack**
- [`assessments/checkpoint/README.md`](./assessments/checkpoint/README.md) — **Module 06 Checkpoint — Testing and Debugging**
- [`assessments/refactor-under-tests/README.md`](./assessments/refactor-under-tests/README.md) — **Refactor Under Tests**

## Completion

Complete every required lesson, lab, project, and assessment for this module.

## Next module

After passing this module's exit gate, continue to [Module 07: CLI, Files, JSON, Config, and Runtime Coordination](../07-cli-files-json-config/README.md).
