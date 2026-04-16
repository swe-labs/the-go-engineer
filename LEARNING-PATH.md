# The Go Engineer Learning Path

This guide explains how to move through the current public curriculum.

The repo now uses an **11-stage learner-facing architecture**.
Those stage pages are the navigation truth.
Some source content is still being regrouped physically, but learners should navigate by stage
first.

## Start With The Stage Index

- [docs/stages/README.md](./docs/stages/README.md)

## The 11 Stages

1. [01 Getting Started](./docs/stages/01-getting-started.md)
2. [02 Language Basics](./docs/stages/02-language-basics.md)
3. [03 Functions & Errors](./docs/stages/03-functions-errors.md)
4. [04 Types & Design](./docs/stages/04-types-design.md)
5. [05 Packages & IO](./docs/stages/05-packages-io.md)
6. [06 Backend & DB](./docs/stages/06-backend-db.md)
7. [07 Concurrency](./docs/stages/07-concurrency.md)
8. [08 Quality & Test](./docs/stages/08-quality-test.md)
9. [09 Architecture](./docs/stages/09-architecture.md)
10. [10 Production](./docs/stages/10-production.md)
11. [11 Flagship](./docs/stages/11-flagship.md)

## Three Ways To Move

### Full Path

Best for:

- complete beginners
- learners new to programming
- learners new to Go who want the most support

Rule:

- follow the 11 stages in order
- complete the required milestone or proof surface in each stage
- do not skip repetition by default

### Bridge Path

Best for:

- experienced programmers who are new to Go
- learners who know programming already but need Go-specific instincts

Rule:

- keep the same stage order
- skim setup repetition where appropriate
- do not skip proof surfaces just because the syntax feels familiar

Suggested route:

1. `01 Getting Started` as a short sanity pass
2. `02 Language Basics`
3. `03 Functions & Errors`
4. continue through the remaining stages in order

### Targeted Path

Best for:

- working Go developers
- learners strengthening one weak area
- learners returning to improve a specific skill without replaying everything

Rule:

- choose one stage intentionally
- check its prerequisites honestly
- complete that stage’s proof surface before claiming mastery

## Recommended Entry Points

| Goal | Start Here | Before You Start |
| --- | --- | --- |
| first Go setup and execution | [01 Getting Started](./docs/stages/01-getting-started.md) | none |
| strengthen fundamentals | [02 Language Basics](./docs/stages/02-language-basics.md) | skim Stage 01 if tooling still feels shaky |
| improve function boundaries and error flow | [03 Functions & Errors](./docs/stages/03-functions-errors.md) | be comfortable with values, control flow, and collections |
| improve modeling and design | [04 Types & Design](./docs/stages/04-types-design.md) | be solid on functions and explicit error handling |
| improve modules, files, encoding, or CLI work | [05 Packages & IO](./docs/stages/05-packages-io.md) | be comfortable with types and design basics |
| build backend applications | [06 Backend & DB](./docs/stages/06-backend-db.md) | be solid on packages, IO, structs, interfaces, and error handling |
| improve concurrency | [07 Concurrency](./docs/stages/07-concurrency.md) | be solid on backend flow, errors, and context-aware thinking |
| improve testing or profiling | [08 Quality & Test](./docs/stages/08-quality-test.md) | be able to build meaningful programs first |
| improve structure and boundaries | [09 Architecture](./docs/stages/09-architecture.md) | be solid on packages, backend boundaries, and testing |
| improve runtime and deployment instincts | [10 Production](./docs/stages/10-production.md) | be comfortable with backend and architecture concerns |
| build a portfolio-level integrated system | [11 Flagship](./docs/stages/11-flagship.md) | have backend, concurrency, testing, architecture, and production depth |

## Validation Floors

Fast paths are allowed.
Proof is not optional.

- `Full Path`: complete every required milestone and stage proof surface
- `Bridge Path`: complete the important stage proof surfaces even if you skim repetition
- `Targeted Path`: complete the chosen stage’s proof surface before claiming mastery

## Current Repo Reality

- the 11-stage model is the public navigation truth
- some source content is already physically migrated into `01-foundations/`
- later-stage source content is still regrouping while stage pages remain stable

If you want the transition explanation, read
[docs/beta-public-architecture.md](./docs/beta-public-architecture.md).

## Companion Docs

- [README.md](./README.md)
- [docs/stages/README.md](./docs/stages/README.md)
- [docs/beta-public-architecture.md](./docs/beta-public-architecture.md)
- [docs/curriculum/README.md](./docs/curriculum/README.md)
- [COMMON-MISTAKES.md](./COMMON-MISTAKES.md)

## Bottom Line

Choose the stage that matches your real background, follow the stage page first, then use the
linked source content from there.
