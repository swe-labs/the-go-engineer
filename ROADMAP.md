# Roadmap

> This document tracks the current release state of the v2.1 curriculum.
> If this file and [ARCHITECTURE.md](./ARCHITECTURE.md) disagree on structure, `ARCHITECTURE.md` wins.

## Current Status

The v2.1 stable release is shipped.

- the 12-section architecture is live
- root-stage source folders are aligned to the public section map
- `curriculum.v2.json` is the active curriculum source
- the single validator is the required repo health check

The next major stream is **post-release implementation work on `main`**, starting with the
Opslane flagship build.

## Branch Model

- `main`: active post-v2.1 development line
- `release/v1`: stable v1 maintenance line
- `release/v2`: stable v2.1.x maintenance line

## Stable Snapshot

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
| s11 Opslane Flagship | Stable-ready | begin post-release flagship implementation |

## Post-Release Focus

After shipping `v2.1.0`, we should be able to say:

- the public docs, metadata, and validator still agree
- stable maintenance stays bounded and low-noise
- flagship implementation work advances on `main` without destabilizing the release line
- the next major engineering gains come from integrated system work, not architecture churn

## Version Plan

| Version | Target | Criteria |
| --- | --- | --- |
| v2.1.0 | released | stable curriculum release is published |
| v2.1.x | maintenance | stable fixes and low-risk corrections on `release/v2` |
| post-v2.1 | current on `main` | flagship implementation and deeper engineering content |
