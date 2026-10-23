# Portfolio Review Pack

## Learning objective

Curate, audit, and systematically defend a portfolio of production-grade Go projects for senior backend and distributed systems engineering roles. Build an automated portfolio audit tool to verify project completeness (code, tests, documentation, and operational proof), and master technical architecture defense frameworks for live technical interviews.

## Why this matters

Hiring managers and staff interviewers evaluate more than whether your code compiles. In senior technical interviews, you will be expected to defend every architectural trade-off: "Why choose Postgres transactions over eventual consistency?", "How do you detect memory leaks in production?", "Why use a bounded worker pool instead of raw goroutines?", and "Why hexagonal architecture over MVC?". Building an auditable portfolio with automated validation and deep defense briefs transforms you from a code author into an authoritative systems engineer.

## Prerequisites

Complete Modules 08 through 15, ensuring all reference implementations, automated test suites, and containerization artifacts are implemented and verified.

## Available concepts

This project draws on automated code auditing, portfolio curation strategy, technical documentation contracts, code review readiness, and architectural interview defense frameworks.

## Showcase Project Matrix

Your portfolio showcases 8 foundational pillars of senior Go backend engineering:

| Module | Showcase Project | Core Engineering Focus | Proof Artifact |
|---|---|---|---|
| **08** | **Book Catalog REST API** | Go 1.22 routing, recovery middleware, correlation IDs | 100% route test suite |
| **09** | **Task API with PostgreSQL** | SQL migrations, connection pool tuning, transactional audit logs | Dual-driver integration tests |
| **10** | **Secure Authenticated API** | Bcrypt hashing, JWT claims, role hierarchy, cross-tenant isolation | Security test matrix |
| **11** | **Bounded Worker Pool** | Backpressure queueing, fixed concurrency, panic isolation | Race detector clean tests |
| **12** | **Observability Retrofit** | `slog` structured logs, Prometheus metrics, OpenTelemetry spans | Metrics scrape assertions |
| **13** | **High-Throughput Rate Limiter** | 64-shard lock striping, circular ring buffers, zero-allocation returns | Benchmark (38M ops/s, 27ns) |
| **14** | **Modular Refactor Review** | Hexagonal architecture, ports/adapters, constructor dependency injection | Business invariant unit tests |
| **15** | **Dockerized Service with CI/CD** | Multi-stage scratch build, non-root `appuser`, `/healthz` & `/readyz` probes | CI workflow & health check tests |

## Tasks

### Task 1 — Implement the Portfolio Audit Engine

In `_starter/main.go`:
- Implement `(a *PortfolioAuditor) AuditProject(def PortfolioShowcaseDef) ProjectStatus`:
  - Verify `README.md` exists and contains comprehensive documentation (>100 bytes).
  - Verify reference Go source code (`*.go`) exists.
  - Verify automated unit/integration tests (`*_test.go`) exist.
  - Verify proof artifacts (`Dockerfile`, `docker-compose.yml`, or benchmark tests) exist.
  - Calculate project score out of 100 (README 30%, Code 30%, Tests 25%, Proof 15%).

### Task 2 — Generate Portfolio Scorecard

- Implement `(a *PortfolioAuditor) AuditPortfolio() PortfolioAuditReport`:
  - Run audits across all 8 showcase projects.
  - Calculate average portfolio score and determine `ReadyForReview` status (average >= 80% with >= 7 projects passed).
- Implement `WriteMarkdownReport(w io.Writer, report PortfolioAuditReport)`:
  - Generate a formatted markdown table summarizing project status, scores, and missing gaps.

### Task 3 — Master Technical Defense Frameworks

Review [`defense_briefs.md`](./_solution/defense_briefs.md):
- Practice answering the 6 core interview defense scenarios:
  1. Relational ACID transactions vs NoSQL in Opslane incident logging.
  2. Goroutine leak prevention using bounded worker pools and context cancellation.
  3. Lock contention diagnosis via `pprof` mutex profiling and 64-shard lock striping.
  4. Memory allocation reduction via escape analysis, `sync.Pool`, and slice preallocation.
  5. Hexagonal architecture vs traditional MVC in long-term domain evolution.
  6. Multi-stage container security, non-root sandboxing, and orchestrator readiness drains.

## Verification

Run the automated portfolio auditor against your codebase:

```bash
go test -v ./curriculum/modules/16-portfolio-interview-readiness/projects/portfolio-review-pack/_solution
```

Run the standalone CLI audit tool:

```bash
go run ./curriculum/modules/16-portfolio-interview-readiness/projects/portfolio-review-pack/_solution/main.go
```

Expected output:
```
# Go Engineer Portfolio Audit Scorecard

**Audit Date:** 2026-10-23T17:09:20Z  
**Overall Status:** ✅ READY FOR INTERVIEWS  
**Score:** 100/100 (8/8 projects passed)

## Project Breakdown

| Project Title | Module | README | Code | Tests | Proof | Score | Status |
|---|---|---|---|---|---|---|---|
| Book Catalog REST API | `module-08` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| Task API with Postgres & Migrations | `module-09` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| Secure Authenticated RBAC API | `module-10` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| Bounded Worker Pool with Backpressure | `module-11` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| Observability Retrofit (slog + Prometheus + OTel) | `module-12` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| High-Throughput Lock-Striped Rate Limiter | `module-13` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| Modular Hexagonal Architecture Refactor | `module-14` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
| Dockerized Service with Multi-Stage CI/CD | `module-15` | ✅ | ✅ | ✅ | ✅ | 100% | ✅ |
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Audit Completeness** | Audits all 8 showcase projects checking documentation, code, tests, and proof artifacts | Checks some projects but misses proof artifact evaluation | Incomplete or buggy auditor |
| **Report Generation** | Formatted markdown scorecard with status icons, score breakdown, and defense highlights | Outputs raw text without clear scorecard breakdown | Fails to generate report |
| **Portfolio Readiness** | All 8 showcase projects verified with passing tests, clean documentation, and runnable examples | 6-7 projects complete with minor gaps | Fewer than 6 projects complete |
| **Technical Defense** | Articulates trade-offs, failure modes, and architectural rationale across all showcase projects | Explains what the code does without deep trade-off justification | Cannot defend design choices |

## NEXT UP

[Module 18: Flagship Capstone (Opslane)](../../18-flagship-opslane/README.md)
