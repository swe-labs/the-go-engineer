# Roadmap

> This document tracks the current release state of the v2.1 curriculum.
> If this file and [ARCHITECTURE.md](./ARCHITECTURE.md) disagree on structure, `ARCHITECTURE.md` wins.

## Current Status

The v2 beta migration is complete on `main`.

- the 12-section architecture is live
- root-stage source folders are aligned to the public section map
- `curriculum.v2.json` is the active curriculum source
- the single validator is the required repo health check

The next release step is **RC hardening**, not more beta migration.

## Branch Model

- `main`: active v2 development and prerelease integration
- `release/v1`: stable v1 maintenance line
- `release/v2`: created when v2 RC stabilization begins

## Beta Completion Snapshot

| Area | Status | Notes |
| --- | --- | --- |
| Public architecture | Complete | 12 sections aligned across root folders and metadata |
| Curriculum metadata | Complete | `curriculum.v2.json` is current and validated |
| Learner-facing section roots | Complete | section entry points exist from `s00` to `s11` |
| Validator | Complete | single Go validator is the required repo health check |
| Beta migration work | Complete | no remaining beta-architecture migration blocker |

## RC Focus

RC work should now concentrate on:

1. learner-path polish and consistency
2. validator strictness and repo hygiene
3. missing test surfaces where behavior should be provable
4. documentation cleanup and release metadata
5. flagship depth and integration confidence

## Section Status

| Section | Status | RC Focus |
| --- | --- | --- |
| s00 How Computers Work | Beta-ready | polish diagrams and machine explanations where useful |
| s01 Getting Started | Beta-ready | keep setup surfaces crisp and trustworthy |
| s02 Language Basics | Beta-ready | preserve zero-magic flow across language, control flow, and data structures |
| s03 Functions & Errors | Beta-ready | continue tightening proof surfaces and orchestration clarity |
| s04 Types & Design | Beta-ready | keep canonical path coherent while stretch content stays clearly optional |
| s05 Packages, I/O & CLI | Beta-ready | strengthen milestone polish and cross-track navigation |
| s06 Backend, APIs & Databases | Beta-ready | keep HTTP, API, and DB tracks aligned and explicit |
| s07 Concurrency | Beta-ready | preserve clear path from fundamentals to patterns |
| s08 Quality & Testing | Beta-ready | keep testing and profiling proof surfaces honest |
| s09 Architecture & Security | Beta-ready | sharpen trade-off teaching and security progression |
| s10 Production Operations | Beta-ready | keep config, observability, deployment, and code generation coherent |
| s11 GoScale Flagship | Beta-ready | deepen integrated proof and release-readiness |

## RC Exit Criteria

Before cutting `release/v2`, we should be able to say:

- the public docs, metadata, and validator agree
- section navigation has no dead internal links
- the repo validates cleanly with one command
- the release-facing docs reflect the real state of the curriculum
- the beta learner path is coherent enough to stabilize instead of restructure

## Version Plan

| Version | Target | Criteria |
| --- | --- | --- |
| v2.1-beta | current on `main` | beta migration complete and validator green |
| v2.1-rc | next | stabilization, polish, and release prep on `release/v2` |
| v2.1 | release | RC passes and release docs are complete |
