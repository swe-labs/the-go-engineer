# V2 Bible

Status: Draft 0.2  
Audience: Maintainers and curriculum designers  
Scope: Planning only. No live curriculum migration starts until the prototype gate is approved.

## Purpose

The v2 Bible is the canonical planning set for the next major version of The Go Engineer.
It exists to stop strategy drift before implementation begins.

v2 is not just "more lessons." It is a curriculum system redesign with:

- a clearer learner model
- a standard lesson and exercise architecture
- a stronger project ladder
- a more explicit folder and metadata contract
- a migration map from v1 to v2
- contributor and validation rules that reduce future drift

## Why v2 Exists

v1 already has strong content depth, real examples, and serious engineering ambition.
Its main weakness is systems consistency: lesson structure, docs, curriculum metadata, exercises,
and repo navigation do not yet feel like one unified training system.

v2 should preserve the strengths of v1 while fixing the delivery system around them.

## Non-Negotiables

The following rules are in force for v2 planning:

1. `main` remains the v2 development line.
2. `release/v1` remains the stable v1 support line.
3. No broad content migration begins before the v2 prototype is approved.
4. Every v1 lesson, exercise, and major doc must have a migration decision:
   keep, revise, split, merge, remove, or add.
5. Exercises must become a first-class system, not an afterthought.
6. v2 must stay repo-first. We are building a training system before considering a platform.

## Planning Gates

v2 work moves through the following gates:

### Gate 0: Thesis Freeze

- Agree on what v2 is and is not.
- Agree on the target learner model.
- Agree that v2 is an incremental migration, not a hidden rewrite.

### Gate 1: Bible Draft

- Publish the first complete planning docs.
- Define the section map, content model, exercise model, and migration rules.
- Capture open questions in the decision log instead of resolving them ad hoc in chat.

### Gate 2: Structural Prototype

- Build one canonical section outline.
- Build one canonical lesson.
- Build one canonical exercise.
- Build one canonical mini-project.
- Build one canonical curriculum metadata example.

No live section migration starts before this gate is approved.

### Gate 3: Issue Breakdown

- Create the GitHub issues, milestones, and project views needed for execution.
- Break the approved system into implementation waves.

### Gate 4: Incremental Migration

- Move section by section on `main`.
- Keep `release/v1` stable.
- Tag alpha releases from `main` at meaningful checkpoints.

## Canonical Document Set

The v2 Bible is split into focused docs so the planning system stays readable:

| File | Role | Current Focus |
| ---- | ---- | ------------- |
| `docs/v2/BIBLE.md` | Index, gates, rules | Live |
| `docs/v2/00-VISION.md` | Product thesis and non-goals | Draft |
| `docs/v2/01-CURRICULUM-PHILOSOPHY.md` | Teaching model and instructional rules | Draft |
| `docs/v2/02-LEARNER-MODEL.md` | Primary learners, pacing, and support rules | Draft |
| `docs/v2/03-CONTENT-TYPE-SYSTEM.md` | Canonical content roles and boundaries | Draft |
| `docs/v2/04-LESSON-SPEC.md` | Canonical lesson contract and definition of done | Draft |
| `docs/v2/05-EXERCISE-BANK-SYSTEM.md` | Exercise architecture and repo layout | Draft |
| `docs/v2/06-PROJECT-LADDER.md` | Project ladder, capstones, and reuse strategy | Draft |
| `docs/v2/07-LEARNING-PATHS.md` | Canonical learner routes, skip rules, and validation floors | Draft |
| `docs/v2/08-CURRICULUM-MAP.md` | Proposed v2 phases and sections | Draft |
| `docs/v2/09-FOLDER-STRUCTURE.md` | Directory rules and migration-safe repo layout | Draft |
| `docs/v2/10-DOCS-NAVIGATION-SYSTEM.md` | Public documentation surfaces and navigation ownership rules | Draft |
| `docs/v2/11-CURRICULUM-SCHEMA.md` | Metadata model and validator contract | Draft |
| `docs/v2/12-MIGRATION-MAP-V1-TO-V2.md` | Content and docs migration rules | Draft |
| `docs/v2/14-QA-VALIDATION.md` | Quality layers, validator scope, CI responsibilities, and review gates | Draft |
| `docs/v2/15-IMPLEMENTATION-ROADMAP.md` | Workstreams, issues, milestones, rollout | Draft |
| `docs/v2/appendix/DECISION-LOG.md` | Open decisions and review outcomes | Live |

The remaining planned docs can be added after the first review pass:

- `13-CONTRIBUTOR-WORKFLOW.md`
- `16-RELEASE-ROLLOUT.md`

## Current Working Assumptions

These assumptions guide the first draft:

- v2 keeps the general topic progression of v1 because the sequence is already strong.
- v2 improves structure, exercises, checkpoints, projects, and docs more than it changes subject matter.
- breaking change value will come from consistency, not from renaming everything for novelty.
- migration should be wave-based and reversible.
- stable users should not be forced onto v2 until the line is clearly better and documented.

## What This Bible Must Produce

Before v2 implementation begins, this planning set must answer:

- What is the promise of v2?
- Who is it for?
- What are the required content types?
- What does a complete section contain?
- How do exercises, checkpoints, and projects fit together?
- How does each v1 asset move into v2?
- What is the execution order for GitHub issues and milestones?

## Current Review Priorities

The next review cycle should focus on five alignment questions:

1. Does the learner model match the actual people this repo serves?
2. Is the canonical lesson spec realistic for maintainers to follow consistently?
3. Does the project ladder feel like a curriculum ladder instead of a pile of ideas?
4. Is the folder structure conservative enough to protect migration stability?
5. Is the curriculum schema rich enough to support validation without becoming bureaucracy?

## Open Questions To Resolve Next

These are the first major questions still open:

- Should v2 keep the exact v1 directory numbering for all sections?
- How much lesson granularity should move from single-file examples to package-sized exercises?
- Should `curriculum.json` be extended or replaced by a richer schema?
- How much assessment structure should exist before we consider a platform or site layer?

Until those questions are answered, maintainers should treat this draft as the source of direction,
not as a frozen spec.

## Working Rule For The Next Drafts

When a future v2 planning doc adds detail, it should prefer:

- extending the current strong v1 sequence over renaming it
- clarifying standards over inventing extra abstractions
- making prototype decisions explicit over leaving them implied

That keeps the Bible practical instead of aspirational-only.
