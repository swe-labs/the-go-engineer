# The Go Engineer - Architecture v2.1 Final (Locked)

## Locked Structure

This file defines the stable public curriculum structure for the v2.1 line.

Architecture v2.1 is final and locked. Normal post-release work should improve lesson depth, proof surfaces, docs, tests, validation, and the Opslane implementation without changing the public 12-section spine.

## Source-of-Truth Rule

If any document, script, template, lesson README, section README, or metadata file disagrees with this file on public curriculum structure, this file wins.

The only exception is an explicit maintainer-approved architecture change. Without that approval, do not add, remove, rename, or reorder public sections.

## Section Map

| Section | Title | Core Tracks | Outputs |
| --- | --- | --- | --- |
| s00 | How Computers Work | HC | HC.5 |
| s01 | Getting Started | GT | GT.6 |
| s02 | Language Basics | LB, CF, DS | LB.4, CF.7, DS.6 |
| s03 | Functions and Errors | FE | FE.10 |
| s04 | Types and Design | TI, CO, ST | TI.15, CO.3, ST.6 |
| s05 | Packages and I/O | MP, CL, EN, FS | MP.4, CL.4, EN.6, FS.8 |
| s06 | Backend, APIs & Databases | HS, API, DB | HS.10, API.9, DB.8 |
| s07 | Concurrency | GC, SY, CT, TM, CP | GC.7, SY.6, CT.5, TM.7, CP.5 |
| s08 | Quality & Testing | TE, PR | TE.10, PR.6 |
| s09 | Architecture & Security | PD, ARCH, SEC | PD.3, ARCH.9, SEC.11 |
| s10 | Production Operations | SL, GS, CFG, OPS, DOCKER, DEPLOY, CG | SL.5, GS.3, CFG.5, OPS.5, DOCKER.3, DEPLOY.3, CG.3 |
| s11 | Flagship | OPSL (+ future flagship prefixes) | OPSL.10 |

## Canonical Root Folders

| Section | Folder |
| --- | --- |
| s00 | `00-how-computers-work` |
| s01 | `01-getting-started` |
| s02 | `02-language-basics` |
| s03 | `03-functions-errors` |
| s04 | `04-types-design` |
| s05 | `05-packages-io` |
| s06 | `06-backend-db` |
| s07 | `07-concurrency` |
| s08 | `08-quality-test` |
| s09 | `09-architecture` |
| s10 | `10-production` |
| s11 | `11-flagship` |

## Structural Notes

- Legacy folders were removed after the active lesson surfaces were relocated.
- Section 04 owns composition and strings/text under `04-types-design/`.
- HC.8 moved into Section 07 as GC.0, and HC.6 moved into Section 08 as PR.6.
- TI.12 and TI.13 were merged into TI.11, and the later advanced TI lessons were renumbered to TI.12-TI.15.
- Section 06 includes HTTP server and API tracks alongside the existing database work.
- Section 09 includes both architecture-pattern and security tracks.
- Section 10 includes configuration, observability, deployment, and code-generation tracks alongside logging and shutdown.
- Code generation is locked into Section 10 because learners usually first meet `go generate` inside Docker and CI workflows before applying it inside the flagship capstone.
- `FG` is a human-facing umbrella for Section 11, while individual flagship projects use their own prefixes such as `OPSL`.

## Change Guardrails

A normal PR may:

- add or revise lessons inside existing sections
- update proof surfaces, tests, and starter code
- strengthen section READMEs
- tighten `curriculum.v2.json`
- improve validators and maintainer scripts
- deepen the Opslane implementation

A normal PR must not:

- create a new public root section
- rename canonical section folders
- move a lesson prefix to a different section without maintainer approval
- revive removed legacy folder structures
- update architecture docs without matching curriculum metadata and validation

## Validation

Any structure-affecting change must pass:

```bash
go run ./scripts/validate_curriculum.go
go test ./...
go vet ./...
```

For full PR readiness, use the complete local verification bundle in `CODE-STANDARDS.md`.
