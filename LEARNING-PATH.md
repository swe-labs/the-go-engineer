# The Go Engineer Learning Path Guide

This guide explains how learners should move through the beta public curriculum.

The repo is now organized around a beta stage model, even though most source content still lives in
the current alpha-era section folders.

For the public explanation of how that works, read
[docs/beta-public-architecture.md](./docs/beta-public-architecture.md).

Start with the stage pages:

- [docs/stages/README.md](./docs/stages/README.md)
- [docs/stages/00-foundation.md](./docs/stages/00-foundation.md)
- [docs/stages/01-language-fundamentals.md](./docs/stages/01-language-fundamentals.md)
- [docs/stages/02-types-and-design.md](./docs/stages/02-types-and-design.md)
- [docs/stages/03-modules-and-io.md](./docs/stages/03-modules-and-io.md)
- [docs/stages/04-backend-engineering.md](./docs/stages/04-backend-engineering.md)
- [docs/stages/05-concurrency-system.md](./docs/stages/05-concurrency-system.md)
- [docs/stages/06-quality-and-performance.md](./docs/stages/06-quality-and-performance.md)
- [docs/stages/07-architecture.md](./docs/stages/07-architecture.md)
- [docs/stages/08-production-engineering.md](./docs/stages/08-production-engineering.md)
- [docs/stages/09-expert-layer.md](./docs/stages/09-expert-layer.md)
- [docs/stages/10-flagship-project.md](./docs/stages/10-flagship-project.md)
- [docs/stages/11-code-generation.md](./docs/stages/11-code-generation.md)

## Current Repo Reality

During beta:

- the stage model is the public navigation truth
- the current section folders remain the source inventory
- some stages regroup multiple alpha sections
- some alpha sections split across more than one beta stage

Use this rule while navigating:

1. choose a path
2. open the stage page for that path
3. follow the linked source tracks and sections from there

If you want the raw source inventory, use [docs/curriculum/README.md](./docs/curriculum/README.md).

## The Three Canonical Paths

The beta learning system uses one curriculum with three routes:

1. `Full Path`
2. `Bridge Path`
3. `Targeted Path`

The routes do not create separate curricula.
They change pacing and entry strategy while keeping one shared system.

## Full Path

### Best For

- complete beginners
- learners new to programming
- learners new to Go who want the most support

### Route Rule

Follow the stages in order:

`0 -> 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8 -> 9 -> 10 -> 11`

### What To Do

- study the stage entry page before entering each stage
- follow the linked source sections in order
- do the required milestone exercises, checkpoints, and mini-projects
- do not skip repetition by default

### Why This Path Exists

This is the honest zero-to-engineer route.
It optimizes for confidence, reinforcement, and fewer hidden gaps.

### Stage Flow

- [0 Foundation](./docs/stages/00-foundation.md)
- [1 Language Fundamentals](./docs/stages/01-language-fundamentals.md)
- [2 Types and Design](./docs/stages/02-types-and-design.md)
- [3 Modules and IO](./docs/stages/03-modules-and-io.md)
- [4 Backend Engineering](./docs/stages/04-backend-engineering.md)
- [5 Concurrency System](./docs/stages/05-concurrency-system.md)
- [6 Quality and Performance](./docs/stages/06-quality-and-performance.md)
- [7 Architecture](./docs/stages/07-architecture.md)
- [8 Production Engineering](./docs/stages/08-production-engineering.md)
- [9 Expert Layer](./docs/stages/09-expert-layer.md)
- [10 Flagship Project](./docs/stages/10-flagship-project.md)
- [11 Code Generation](./docs/stages/11-code-generation.md)

## Bridge Path

### Best For

- experienced programmers who are new to Go
- learners who know programming already but need Go-specific instincts

### Route Rule

Keep the same stage order, but skip repetition rather than proof.

Recommended route:

`0 (skim) -> 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8 -> 10 -> 11`

The `9 Expert Layer` remains valuable, but it is not the first stop for most bridge learners.

### What To Skim

Bridge learners may skim:

- some beginner setup detail in `0 Foundation`
- repetition-heavy early examples once the toolchain and basic execution model are clear
- selected drills that only repeat concepts they already command

### What Not To Skip

Bridge learners should not skip:

- Go-specific function and error handling patterns
- structs, interfaces, and composition decisions
- package and module boundaries
- concurrency checkpoints if they plan to enter the concurrency stages
- milestone exercises that prove real Go understanding

### Best Starting Point

Start here:

- [0 Foundation](./docs/stages/00-foundation.md) for a fast setup and tooling sanity pass
- then move quickly into [1 Language Fundamentals](./docs/stages/01-language-fundamentals.md)

### Why This Path Exists

This route is for transfer, not for ego.
If you already know programming, you should move faster, but you still need proof that your prior
mental models fit Go's style.

## Targeted Path

### Best For

- working Go developers
- learners fixing a specific weakness
- learners returning to strengthen one area without replaying the whole curriculum

### Route Rule

Choose a stage intentionally, review the prerequisites honestly, then complete the stage's required
practice and proof surfaces.

### Entry Rule

Before entering a later stage, ask:

1. Do I really understand the earlier concepts this stage depends on?
2. Can I complete the relevant checkpoint or milestone work, not just read the lessons?
3. Am I skipping repetition, or am I skipping proof?

### Recommended Entry Points

| Goal | Start Here | Before You Start |
| --- | --- | --- |
| tighten Go fundamentals | [1 Language Fundamentals](./docs/stages/01-language-fundamentals.md) | skim [0 Foundation](./docs/stages/00-foundation.md) if tooling or execution still feels shaky |
| improve data modeling and design | [2 Types and Design](./docs/stages/02-types-and-design.md) | be solid on functions, errors, and collections |
| strengthen modules, files, encoding, or CLI work | [3 Modules and IO](./docs/stages/03-modules-and-io.md) | be comfortable with structs, errors, and basic package flow |
| build backend applications | [4 Backend Engineering](./docs/stages/04-backend-engineering.md) | be solid on modules, IO, structs, interfaces, and error handling |
| improve concurrency | [5 Concurrency System](./docs/stages/05-concurrency-system.md) | be solid on backend flow, errors, and context-aware thinking |
| improve testing or profiling | [6 Quality and Performance](./docs/stages/06-quality-and-performance.md) | be able to build meaningful programs first |
| improve architecture | [7 Architecture](./docs/stages/07-architecture.md) | be solid on packages, backend boundaries, and testing |
| improve operations and runtime practices | [8 Production Engineering](./docs/stages/08-production-engineering.md) | be comfortable with backend and architecture concerns |
| focus on review and diagnosis | [9 Expert Layer](./docs/stages/09-expert-layer.md) | come in with real implementation experience |
| build a portfolio-level integrated system | [10 Flagship Project](./docs/stages/10-flagship-project.md) | have backend, concurrency, testing, and architecture depth |
| use generation responsibly | [11 Code Generation](./docs/stages/11-code-generation.md) | understand the systems and code shapes being generated |

### Why This Path Exists

Not every serious learner should restart from zero.
But targeted entry only works when prerequisites are checked honestly.

## Validation Floors

Fast paths are allowed.
Proof is not optional.

Use these validation floors:

- `Full Path`: complete every required checkpoint and milestone mini-project
- `Bridge Path`: complete every checkpoint plus the important stage milestones
- `Targeted Path`: complete the chosen stage's checkpoint or equivalent proof artifact before
  claiming mastery

## Mental-Model Rule

Every beta stage includes a mental-model surface.
Use the stage pages for that.

That means:

- beginners should not skip the stage README and mental-model framing
- experienced learners may move faster, but should still read the stage framing before jumping into
  source files
- advanced learners should use the stage page to verify they are entering the right layer of the
  curriculum

## Suggested Starting Points By Learner Type

### Complete Beginner

Start with [0 Foundation](./docs/stages/00-foundation.md), then follow the Full Path in order.

### Experienced Programmer New To Go

Start with [0 Foundation](./docs/stages/00-foundation.md) for a quick environment and execution
pass, then move into [1 Language Fundamentals](./docs/stages/01-language-fundamentals.md) on the
Bridge Path.

### Working Go Developer

Choose your target stage from the table above and use the Targeted Path honestly.

## How This Relates To The Existing Sections

The beta stages regroup the current source sections like this:

- `0 Foundation`: `01-foundations/01-getting-started`
- `1 Language Fundamentals`: `01-core-foundations/language-basics`, `02`, `03`, `04`
- `2 Types and Design`: `05`, `06`, `07`
- `3 Modules and IO`: `08`, `09`
- `4 Backend Engineering`: `10`
- `5 Concurrency System`: `11`, `12`
- `6 Quality and Performance`: `13`
- `7 Architecture`: architecture-focused tracks from `14`
- `8 Production Engineering`: operations-focused tracks from `14`
- `9 Expert Layer`: beta additions
- `10 Flagship Project`: `enterprise-capstone` plus later flagship regrouping
- `11 Code Generation`: `15`

## Recommended Companion Docs

- [README.md](./README.md)
- [docs/stages/README.md](./docs/stages/README.md)
- [docs/beta-public-architecture.md](./docs/beta-public-architecture.md)
- [docs/curriculum/README.md](./docs/curriculum/README.md)
- [COMMON-MISTAKES.md](./COMMON-MISTAKES.md)
- [CONTRIBUTING.md](./CONTRIBUTING.md)

## Bottom Line

The beta learning-path rule is simple:

- choose the path that matches your background honestly
- navigate by stage first
- follow the linked source content from the stage page
- move fast if you want, but do not skip proof
