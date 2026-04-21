# The Go Engineer - Architecture v2.1-final (locked)

## Locked Structure

This file is the curriculum shape after the restructure plan has been applied.

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
| s10 | Production Operations | SL, GS, CFG, OPS, DOCKER, DEPLOY, CG | SL.5, GS.3, CFG.5, OPS.5, DEPLOY.3, CG.3 |
| s11 | Flagship | FG | FG.1 |

## Structural Notes

- Legacy folders were removed after the active lesson surfaces were relocated.
- Section 04 now owns composition and strings/text under 04-types-design/.
- HC.8 moved into Section 07 as GC.0 and HC.6 moved into Section 08 as PR.6.
- TI.12 and TI.13 were merged into TI.11, and the later advanced TI lessons were renumbered to TI.12-TI.15.
- Section 06 now includes HTTP server and API tracks alongside the existing database work.
- Section 09 now includes both architecture-pattern and security tracks.
- Section 10 now includes configuration, observability, deployment, and code-generation tracks alongside logging and shutdown.
- Code generation is locked into Section 10 because learners usually first meet `go generate` inside Docker and CI workflows before they apply it inside the flagship capstone.
