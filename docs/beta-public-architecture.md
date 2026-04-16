# Source Inventory vs Public 11-Stage Architecture

The Go Engineer is in a transition period, but the public navigation rule is now simple:

- **the 11-stage model is the learner-facing truth**
- **the current source folders are the physical implementation layer**

This document explains how those two layers fit together.

## The Short Version

If you are a learner:

- start from [README.md](../README.md)
- then follow [docs/stages/README.md](./stages/README.md)

If you are trying to find files:

- use the current top-level folders
- use [docs/curriculum/README.md](./curriculum/README.md)

The stage model tells learners **where to go next**.
The folder tree tells contributors **where content lives today**.

## The Public 11 Stages

1. `01 Getting Started`
2. `02 Language Basics`
3. `03 Functions & Errors`
4. `04 Types & Design`
5. `05 Packages & IO`
6. `06 Backend & DB`
7. `07 Concurrency`
8. `08 Quality & Test`
9. `09 Architecture`
10. `10 Production`
11. `11 Flagship`

These 11 stage pages are the public routing surfaces for the curriculum.

## How The Current Source Maps

The current source inventory feeds the 11-stage model like this:

- `01-foundations/01-getting-started` -> `01 Getting Started`
- `01-foundations/02-language-basics`, `01-foundations/03-control-flow`, `01-foundations/04-data-structures` -> `02 Language Basics`
- `01-foundations/05-functions-and-errors` -> `03 Functions & Errors`
- `01-foundations/06-types-and-interfaces`, `05-composition`, `06-strings-and-text` -> `04 Types & Design`
- `07-modules-and-packages`, `08-io-and-cli` -> `05 Packages & IO`
- `09-web-and-database` -> `06 Backend & DB`
- `10-concurrency`, `11-concurrency-patterns` -> `07 Concurrency`
- `12-quality-and-performance` -> `08 Quality & Test`
- `13-application-architecture/package-design`, `13-application-architecture/grpc` -> `09 Architecture`
- `13-application-architecture/structured-logging`, `13-application-architecture/graceful-shutdown`, `13-application-architecture/docker-and-deployment` -> `10 Production`
- `13-application-architecture/enterprise-capstone`, `14-code-generation` -> `11 Flagship`

## Why We Use Both Layers

We are not rewriting the repo from scratch in one destructive move.

Instead, we are doing two things deliberately:

1. keeping the learner-facing 11-stage architecture stable
2. finishing the physical folder migration behind that public surface

That lets the curriculum become clearer for learners without pretending every folder move is already
finished.

## What To Trust For What

### If you are a learner

Use:

- [README.md](../README.md)
- [LEARNING-PATH.md](../LEARNING-PATH.md)
- [docs/stages/README.md](./stages/README.md)

### If you want the physical source layout

Use:

- [docs/curriculum/README.md](./curriculum/README.md)
- the current top-level folders

### If you are a maintainer or contributor

Use:

- the stage pages for public routing truth
- the current folder tree for implementation truth
- `ARCHITECTURE.md` and `CURRICULUM-BLUEPRINT.md` for the longer-range design

## What This Does Not Mean

This transition model does **not** mean:

- learners should navigate by old stage names
- the repo has competing public architectures
- every old planning surface is still public truth
- physical regrouping must finish before the stage model can be honest

## Bottom Line

The 11-stage architecture is the public curriculum.
The current source tree is the implementation surface that is still being finished underneath it.
