# Changelog

All notable changes to The Go Engineer are documented here.

Format: **[Date] — Description**. Sections: `Added`, `Fixed`, `Changed`, `Removed`.

---

## [2026-04-06] — Curriculum Dependency Integration

### Added

- **Curriculum Dependency Graph:** Automated granular mapping of 15 sections across 28 sub-graphs using `curriculum.json`.
- **Premium Console Footers:** Integrated "Next Up" context-aware navigation footers into all 115 `main.go` files.
- **Visual Breadcrumbs:** Added previous/next/exercise navigation links to all section READMEs.
- **Root Metadata:** Updated root `README.md` and `ROADMAP.md` with the new ID-based learning path (GS, LB, CF, etc.).
- **Validation Script:** Created `scripts/validate_curriculum.go` to ensure filesystem integrity for all future additions.

### Fixed

- **ID Collision:** Resolved ID collisions for "Graceful Shutdown" vs "Getting Started" by implementing section-specific context namespaces.
- **Bulk Refactor Cleanup:** Fixed widespread Go syntax errors and markdown lint issues introduced during the integration phase.
- **[BUG] Chapter 1 Lesson 1:** Added missing instructions to `1-installation` (GS.1) boilerplate (#4).

---

## [2026-04-05] — The v1.0 Curriculum Migration

### Changed

- **Curriculum Architecture:** Completely refactored the legacy 28 flat-folder structure into an elite, chronologically ordered 15-chapter book format.
- **Cross-Domain Separation:** Decoupled concepts by their primary domain, cleanly separating runtime operations (`11-concurrency/time-and-scheduling`) and external operations (`09-io-and-cli/encoding`).
- **HTTP Domain Split:** Amputated the HTTP Clients/Mocking directory exactly down the middle, separating true HTTP networking (`10-web-and-database/http-client`) from test abstraction limits (`13-quality-and-performance/http-client-testing`).
- **Internal Imports:** Ran a mass-migration script converting all legacy `internal/` pointers uniformly.
- **Docs Update:** Rewrote README tables, ROADMAP tracking lists, and CONTRIBUTING guides natively reflecting the elite `v1` path routing.

### Removed

- DELETED the legacy compatibility `symlinks` from the repository root entirely!

---

## [2026-04-04] — Phase 7 additions + bug fixes

### Added

- **Section 23: Structured Logging** — slog basics, context-keyed logger, custom Handler, zerolog comparison
- **Section 24: errgroup & sync.Pool** — errgroup, errgroup+context pipeline, sync.Pool, bounded pipeline resizer exercise, URL checker exercise
- **Section 25: Profiling** — CPU profile, live `net/http/pprof` endpoint, go tool pprof workflow
- **Section 26: gRPC** — proto3 definition, unary server with interceptors, typed client stub
- **Section 27: Graceful Shutdown** — signal.NotifyContext, http.Server.Shutdown, complete production capstone
- **Section 10 supplement** — `io/fs` as a testing seam with `fstest.MapFS`
- **`COMMON-MISTAKES.md`** — 15 most common Go bugs with fixes and section cross-references
- **`ROADMAP.md`** — tracks current progress and planned future sections
- **`CHANGELOG.md`** — this file

### Fixed

- **Bug:** `23-structured-logging/3-custom-handler` — `slog.DiscardHandler` instantiation type error fixed by removing struct brackets.
- **Bug:** `24-errgroup-and-pools/3-sync-pool` — duplicate `BenchmarkWithPool` declaration causing compilation failure removed from `main.go`.
- **Bug:** `25-profiling/1-cpu-profile` and `10-filesystem/8-fs-testing-seam` — redundant newline formatting in `fmt.Println` fixed to standard `fmt.Print`.
- **Bug:** `26-grpc` — missing generated Protocol Buffer Go files added using `protoc` and properly imported into `gen` package.
- **Bug:** Global repository string literals — fixed widespread compilation errors caused by improper line breaks across several past modules.

### Changed

- `go.mod` and Dockerfiles now correctly declare `go 1.26` (was previously correct — this entry confirms it intentionally targets Go 1.26 stable)
- README exercise table updated to include all new sections 23–27
- README Windows CGO note added for `go-sqlite3` dependency

---

## [2026-01-15] — Enterprise Capstone

### Added

- Section 22: Enterprise Capstone — multi-package PostgreSQL REST API with Docker Compose
- Section 21: Database Migrations — golang-migrate with embedded SQL files
- Section 20: Docker & Deployment — multi-stage Dockerfiles

---

## [2025-11-01] — Web Masterclass

### Added

- Section 13: Web Masterclass — complete 11-part progression from routing to full production app
- Section 16: HTTP Clients & Mocking — testify/mock, manual mocking, table-driven tests

---

## [2025-09-01] — Initial Release

### Added

- Sections 00–19: complete Beginner → Expert learning path
- 17 exercises with `_starter/` stubs and full solutions
- GitHub Actions CI pipeline
- CONTRIBUTING.md with file templates and quality checklist
