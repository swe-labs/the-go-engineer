# The Go Engineer

[![CI](https://github.com/swe-labs/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/swe-labs/the-go-engineer/actions)
[![License: TGE License v1.0](https://img.shields.io/badge/License-TGE_v1.0-red.svg)](#license)

The Go Engineer is a repository-first Go software engineering curriculum. It teaches Go by combining runnable lessons, production-shaped examples, tests, validation, and a final integrated backend project.

The stable v2.1 line is organized as a 5-phase, 12-section learning system with 215 registered curriculum items. The public structure is locked in [ARCHITECTURE.md](./ARCHITECTURE.md), and the machine-readable registry is [curriculum.v2.json](./curriculum.v2.json).

## Status

Current stable line: `v2.1.x`

Current stable release: `v2.1.1`

Supported branches:

| Branch | Purpose |
| --- | --- |
| `main` | active post-v2.1 implementation and integration line |
| `release/v2` | stable v2.1.x maintenance line |
| `release/v1` | stable v1 maintenance line |

Architecture v2.1 is locked. Normal work may improve lessons, tests, documentation, validators, or the Opslane flagship implementation, but must not add, remove, rename, or reorder public root sections without explicit maintainer approval.

## Quick Start

Requirements:

- Go version declared in [go.mod](./go.mod)
- CGO-capable C compiler for `go test -race ./...` and SQLite-backed paths

```bash
git clone https://github.com/swe-labs/the-go-engineer.git
cd the-go-engineer
go mod download
go run ./00-how-computers-work/1-what-is-a-program
```

Run the curriculum validator:

```bash
go run ./scripts/validate_curriculum.go
```

Expected stable output:

```text
Success! 597 files with run commands validated, and 12 v2 sections plus 215 v2 items checked.
```

## Curriculum Map

| Phase | Sections | Focus | Progress |
| --- | --- | --- | --- |
| 0 | s00 | machine foundation | 0% to 5% |
| 1 | s01-s04 | language foundation | 5% to 52% |
| 2 | s05-s08 | engineering core | 52% to 87% |
| 3 | s09-s10 | systems engineering | 87% to 96% |
| 4 | s11 | integrated flagship project | 96% to 100% |

## Sections

| Section | Folder | Primary outcome |
| --- | --- | --- |
| s00 | [00-how-computers-work](./00-how-computers-work) | understand execution, memory, terminal use, and processes |
| s01 | [01-getting-started](./01-getting-started) | install Go, run programs, and read tool output |
| s02 | [02-language-basics](./02-language-basics) | use values, control flow, slices, maps, and pointers intentionally |
| s03 | [03-functions-errors](./03-functions-errors) | compose functions and handle failure explicitly |
| s04 | [04-types-design](./04-types-design) | model behavior with structs, methods, interfaces, composition, and generics |
| s05 | [05-packages-io](./05-packages-io) | work with modules, packages, CLI input, encoding, and filesystems |
| s06 | [06-backend-db](./06-backend-db) | build HTTP servers, API contracts, gRPC services, and database-backed code |
| s07 | [07-concurrency](./07-concurrency) | coordinate goroutines, channels, context, synchronization, and worker patterns |
| s08 | [08-quality-test](./08-quality-test) | test, benchmark, profile, and verify behavior with evidence |
| s09 | [09-architecture](./09-architecture) | apply package design, architecture patterns, and security controls |
| s10 | [10-production](./10-production) | add logging, shutdown, configuration, observability, deployment, and code generation |
| s11 | [11-flagship](./11-flagship) | integrate the curriculum through the Opslane SaaS backend capstone |

## Learning Model

Each lesson is designed around four surfaces:

- a README-first explanation with prerequisites, mental model, machine view, and production notes
- runnable Go code that demonstrates the concept directly
- tests or verification surfaces where behavior should be provable
- forward and backward references when a concept depends on another lesson

The curriculum avoids unexplained jumps. If a later idea appears early, the lesson must name it as a forward reference. If a current idea depends on previous work, the lesson must point back to the lesson that established it.

## Validation

Run the CI-equivalent bundle before release or review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile coverage.out ./...
go run ./scripts/validate_curriculum.go
```

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/
```

## Documentation

| Document | Purpose |
| --- | --- |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | locked public curriculum structure |
| [curriculum.v2.json](./curriculum.v2.json) | machine-readable curriculum registry |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md) | teaching contract and lesson authoring rules |
| [LEARNING-PATH.md](./LEARNING-PATH.md) | recommended paths through the curriculum |
| [docs/PROGRESSION.md](./docs/PROGRESSION.md) | phase and milestone progression |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md) | Go source, comments, and teaching-code standards |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md) | test and verification expectations |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md) | common Go mistakes and where they are taught |
| [docs/ENGINEERING_ERROR_FRAMEWORK.md](./docs/ENGINEERING_ERROR_FRAMEWORK.md) | backend and production error model |
| [docs/flagship/OPSLANE_SAAS_BACKEND.md](./docs/flagship/OPSLANE_SAAS_BACKEND.md) | Opslane flagship specification |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | contribution and pull request workflow |
| [RELEASE.md](./RELEASE.md) | release and branch maintenance process |
| [MAINTAINER-CHECKLIST.md](./MAINTAINER-CHECKLIST.md) | maintainer review and release checklist |
| [CHANGELOG.md](./CHANGELOG.md) | release history |

## Contribution Workflow

Contributions follow [CONTRIBUTING.md](./CONTRIBUTING.md):

1. Create or confirm an issue.
2. Branch from the correct long-lived line.
3. Open a draft PR early.
4. Keep commits focused and use the bracketed commit format.
5. Run local validation.
6. Address review findings and CI failures.
7. Hand off for maintainer review and squash merge.

Commit format:

```text
[TYPE] short imperative description
```

Allowed types:

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

## Windows CGO Note

Some paths use CGO through SQLite or the Go race detector. On Windows, install a GCC-compatible toolchain such as MinGW-w64 or use WSL2, then verify:

```bash
go env CGO_ENABLED CC CXX
go test -race ./...
```

## License

This project is licensed under The Go Engineer License (TGE License) v1.0.

- Free for personal, educational, and non-commercial use
- Commercial use requires written permission
