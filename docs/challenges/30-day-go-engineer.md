# 30-Day Go Engineer Challenge

This challenge gives learners a paced route through the stable v2.1 curriculum.

Use it when you want structure, public accountability, or a team study plan. It does not replace [LEARNING-PATH.md](../../LEARNING-PATH.md); it turns the same curriculum into a 30-day schedule.

## Before You Start

Set up the repository:

```bash
git clone https://github.com/swe-labs/the-go-engineer.git
cd the-go-engineer
go mod download
go run ./scripts/validate_curriculum.go
```

Expected validator output:

```text
Success! 601 files with run commands validated, and 12 v2 sections plus 215 v2 items checked.
```

## Week 1: Machine And Language Foundations

Goal: understand what Go programs are doing before moving into larger code.

- Day 1: s00, machine basics and terminal confidence
- Day 2: s01, installation, toolchain, and compiler output
- Day 3: s02, variables, constants, and control flow
- Day 4: s02, arrays, slices, maps, and pointers
- Day 5: s03, functions and multiple returns
- Day 6: s03, errors as values and validation
- Day 7: review, rerun examples, and write down unclear concepts

## Week 2: Types, Packages, I/O, And Backend Basics

Goal: move from small programs to maintainable application boundaries.

- Day 8: s04, structs, methods, and interfaces
- Day 9: s04, composition, generics, and text handling
- Day 10: s05, modules and package boundaries
- Day 11: s05, CLI input, encoding, and filesystem work
- Day 12: s06, HTTP server fundamentals
- Day 13: s06, APIs, databases, and repositories
- Day 14: review by running tests for sections s04-s06

## Week 3: Concurrency, Quality, Architecture, And Security

Goal: learn how production Go code handles coordination, evidence, and risk.

- Day 15: s07, goroutines, wait groups, and channels
- Day 16: s07, context, cancellation, and worker patterns
- Day 17: s08, unit tests, subtests, mocks, and golden files
- Day 18: s08, benchmarks, profiling, and fuzzing
- Day 19: s09, package design and architecture patterns
- Day 20: s09, security lessons and secure API exercise
- Day 21: review with `go test ./...` and `go run ./scripts/validate_curriculum.go`

## Week 4: Production And Opslane

Goal: connect the curriculum to an integrated SaaS backend.

- Day 22: s10, structured logging and configuration
- Day 23: s10, graceful shutdown and deployment
- Day 24: s10, observability and code generation
- Day 25: s11, Opslane foundation, database, auth, and HTTP API
- Day 26: s11, order processing, payment pipeline, and event workers
- Day 27: s11, caching, observability, shutdown, and deployment
- Day 28: run Opslane tests and read the threat model
- Day 29: revisit weak areas and improve one note or test locally
- Day 30: share a completion note with the command output you verified

## Completion Proof

A good completion note includes:

- the final section or module reached
- one concept that changed how you write Go
- one test, validator, or Opslane command that passed locally
- one next improvement you want to make

Suggested final checks:

```bash
go test ./...
go run ./scripts/validate_curriculum.go
```
