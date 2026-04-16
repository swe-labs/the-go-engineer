# Roadmap

This document tracks what is built, what is in progress, and what is planned.

## Branch Model

- `main` tracks active v2 development and prereleases.
- `release/v1` tracks stable v1 maintenance for current users.
- `release/v2` will be cut from `main` when v2 reaches beta and feature freeze.

## Current V2 Focus

As of April 7, 2026, the v2 transition is in the public planning and prototype-design phase.

- Stable learners should keep using `release/v1`.
- Public v2 planning now lives on [`planning/v2`](https://github.com/rasel9t6/the-go-engineer/tree/planning/v2).
- Design work is tracked in the [V2 project board](https://github.com/users/rasel9t6/projects/2).
- Initial public planning issues are [#79](https://github.com/rasel9t6/the-go-engineer/issues/79) through [#88](https://github.com/rasel9t6/the-go-engineer/issues/88).
- Broad content migration will not begin until the structural prototype is approved.

## Status Legend

| Symbol | Meaning |
|--------|---------|
| âœ… | Complete and tested |
| ðŸš§ | In progress |
| ðŸ“‹ | Planned |
| ðŸ’¡ | Under consideration |

---

## 01 â€” Core Foundations âœ…

- `getting-started`: installation, hello world, compilation model, dev tools (GS.1â€“4)
- `language-basics`: variables, types, zero values, constants, iota (LB.1â€“4)

## 02 â€” Control Flow âœ…

- `control-flow`: for, if/else, switch, line-of-sight principle (CF.1â€“4)

## 03 â€” Data Structures âœ…

- `data-structures`: arrays, slices, maps, pointers, escape analysis (DS.1â€“6)

## 04 â€” Functions and Errors âœ…

- `functions-and-errors`: multi-return, closures, defer, panic/recover, error wrapping (FE.1â€“9)

## 05 â€” Types and Interfaces âœ…

- `types-and-interfaces`: structs, methods, interfaces, generics (TI.1â€“6)

## 06 â€” Composition âœ…

- `composition`: composition vs inheritance (CO.1â€“3)

## 07 â€” Strings and Text âœ…

- `strings-and-text`: string internals, regex, templates, Builder (ST.1â€“6)

## 08 â€” Modules and Packages âœ…

- `modules-and-packages`: go.mod, versioning, replace (MP.1â€“3)

## 09 â€” IO and CLI âœ…

- `filesystem`: file I/O, paths, directories, embed, io.Reader/Writer patterns (FS.1â€“8)
- `cli-tools`: os.Args, flag, subcommands, file organiser (CL.1â€“4)
- `encoding`: JSON marshal/unmarshal, streaming, base64 (EN.1â€“6)

## 10 â€” Web and Database âœ…

- `databases`: sql.DB, connection pooling, transactions, repository pattern (DB.1â€“6)
- `database-migrations`: golang-migrate, embed, schema evolution
- `web-masterclass`: routing, DI, templates, middleware, sessions, auth, forms, CRUD, pagination, comments (WM.1â€“11)
- `http-client`: http.Client configuration, calling APIs (HC.1â€“2)

## 11 â€” Concurrency âœ…

- `concurrency`: goroutines, WaitGroup, channels, select, race conditions, sync primitives (GC.1â€“10)
- `context`: Background, WithCancel, WithTimeout, WithValue (CT.1â€“5)
- `time-and-scheduling`: time, scheduling, timers, tickers (TM.1â€“7)

## 12 â€” Concurrency Patterns âœ…

- `concurrency-patterns`: errgroup.Group, fan-out pipelines, sync.Pool zero-allocation (CP.1â€“5)

## 13 â€” Quality and Performance âœ…

- `testing`: unit tests, table-driven, HTTP handler tests, benchmarking (TE.1â€“4)
- `http-client-testing`: manual mocks, testify/mocking abstractions (HM.1â€“4)
- `profiling`: CPU profiles, memory profiles, flame graphs, pprof (PR.1â€“2)

## 14 â€” Application Architecture âœ…

- `package-design`: naming, visibility, internal/, project layout (PD.1â€“3)
- `docker-and-deployment`: single-stage Dockerfile, multi-stage builds, layer caching (DO.1â€“3)
- `enterprise-capstone`: full REST API, PostgreSQL, Docker Compose (EC.1)
- `structured-logging`: slog basics, context-keyed logger, custom handlers (SL.1â€“4)
- `grpc`: proto3 service definition, unary and streaming RPC (GR.1â€“3)
- `graceful-shutdown`: signal.NotifyContext, HTTP graceful drain (GS.1â€“2)

## 15 â€” Code Generation ðŸ“‹

- `go-generate`: //go:generate directive, tools as part of build (CG.1)
- `mockery` & `stringer` automation

---

## 29 â€” Cloud Native ðŸ’¡

- OpenTelemetry traces and metrics
- Kubernetes health probes (liveness vs readiness)
- Graceful shutdown for gRPC servers
- Config management with environment variables and Viper

## 30 â€” Performance Patterns ðŸ’¡

- `atomic.Value` for lock-free config hot-reload
- SIMD-friendly data layout patterns

---

## Windows Setup Note

The `09-web-and-database/databases` and `13-application-architecture/enterprise-capstone` sections depend on `github.com/mattn/go-sqlite3` which requires CGO and a C compiler. On Windows without WSL:

1. Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
2. Ensure `gcc` is in your PATH: `gcc --version`
3. Set `CGO_ENABLED=1` in your environment

Alternatively, use WSL2 (recommended). All other sections work on Windows without a C compiler.

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how to add new lessons, exercises, and sections.
