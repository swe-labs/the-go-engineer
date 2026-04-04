# Roadmap

This document tracks what is built, what is in progress, and what is planned.

## Status Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Complete and tested |
| 🚧 | In progress |
| 📋 | Planned |
| 💡 | Under consideration |

---

## 01 — Core Foundations ✅

- `getting-started`: installation, hello world, compilation model, dev tools
- `language-basics`: variables, types, zero values, constants, iota

## 02 — Control Flow ✅

- `control-flow`: for, if/else, switch, line-of-sight principle

## 03 — Data Structures ✅

- `data-structures`: arrays, slices, maps, pointers, escape analysis

## 04 — Functions and Errors ✅

- `functions-and-errors`: multi-return, closures, defer, panic/recover, error wrapping

## 05 — Types and Interfaces ✅

- `types-and-interfaces`: structs, methods, interfaces, generics

## 06 — Composition ✅

- `composition`: composition vs inheritance

## 07 — Strings and Text ✅

- `strings-and-text`: string internals, regex, templates, Builder

## 08 — Modules and Packages ✅

- `modules-and-packages`: go.mod, versioning, replace

## 09 — IO and CLI ✅

- `filesystem`: file I/O, paths, directories, embed, io.Reader/Writer patterns
- `cli-tools`: os.Args, flag, subcommands, file organiser
- `encoding`: JSON marshal/unmarshal, streaming, base64

## 10 — Web and Database ✅

- `databases`: sql.DB, connection pooling, transactions, repository pattern
- `database-migrations`: golang-migrate, embed, schema evolution
- `web-masterclass`: routing, DI, templates, middleware, sessions, auth, forms, CRUD, pagination, comments
- `http-client`: http.Client configuration, calling APIs

## 11 — Concurrency ✅

- `concurrency`: goroutines, WaitGroup, channels, select, race conditions, sync primitives
- `context`: Background, WithCancel, WithTimeout, WithValue
- `time-and-scheduling`: time, scheduling, timers, tickers

## 12 — Concurrency Patterns ✅

- `concurrency-patterns`: errgroup.Group, fan-out pipelines, sync.Pool zero-allocation

## 13 — Quality and Performance 🚧

- `testing`: unit tests, table-driven, HTTP handler tests, benchmarking
- `http-client-testing`: manual mocks, testify/mocking abstractions
- `profiling`: CPU profiles, memory profiles, flame graphs, pprof 🚧

## 14 — Application Architecture 🚧

- `package-design`: naming, visibility, internal/, project layout
- `docker-and-deployment`: multi-stage Dockerfile, layer caching
- `enterprise-capstone`: full REST API, PostgreSQL, Docker Compose
- `structured-logging`: slog basics, context-keyed logger
- `grpc`: proto3 service definition, unary and streaming RPC 🚧
- `graceful-shutdown`: signal.NotifyContext, full graceful drain shutdown

## 15 — Code Generation 📋

- `go-generate`: //go:generate directive, tools as part of build
- `mockery` & `stringer` automation

---

## 29 — Cloud Native 💡

- OpenTelemetry traces and metrics
- Kubernetes health probes (liveness vs readiness)
- Graceful shutdown for gRPC servers
- Config management with environment variables and Viper

## 30 — Performance Patterns 💡

- `atomic.Value` for lock-free config hot-reload
- SIMD-friendly data layout patterns

---

## Windows Setup Note

The `10-web-and-database/databases` and `14-application-architecture/enterprise-capstone` sections depend on `github.com/mattn/go-sqlite3` which requires CGO and a C compiler. On Windows without WSL:

1. Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
2. Ensure `gcc` is in your PATH: `gcc --version`
3. Set `CGO_ENABLED=1` in your environment

Alternatively, use WSL2 (recommended). All other sections work on Windows without a C compiler.

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how to add new lessons, exercises, and sections.
