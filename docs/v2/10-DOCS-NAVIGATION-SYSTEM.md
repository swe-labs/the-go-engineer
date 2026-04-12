# V2 Docs Navigation System

## Purpose

v2 needs a navigation system that helps learners move through the repo without forcing them to infer
the curriculum structure from scattered READMEs, issue threads, or directory names alone.

This document defines the canonical public documentation surfaces and the job each one should do.

## Core Rule

Each public doc surface should have one primary navigation job.

If the root README, the learning-path guide, the curriculum map, and section READMEs all try to be
full curriculum overviews, they will drift and confuse learners.

## Current Problems To Fix

The current repo already has useful navigation material, but it overlaps too much.

Main pain points:

- the root `README.md` tries to be landing page, curriculum map, learner guide, and release note
- `LEARNING-PATH.md` and `docs/curriculum/README.md` repeat large parts of the same story
- section-level README coverage is inconsistent or missing
- some sub-track READMEs are strong, but they do not fit into a consistent repo-wide contract
- learners often need to infer "what next" from folder order rather than explicit navigation

v2 should reduce duplication and make route logic explicit.

## Navigation Goals

The v2 docs navigation system should help learners answer these questions quickly:

- what is this repo
- where should I start
- which path fits me
- what does this section contain
- what should I do next
- what can I skip
- what proves I am ready to move on

The system should also help contributors answer:

- which doc owns which explanation
- where path guidance belongs
- where section outcomes and milestone rules should live

## Canonical Navigation Surfaces

The first v2 draft should use these public navigation surfaces.

### 1. Root README

Primary job:

- serve as the front door to the repo

It should answer:

- what The Go Engineer is
- who it is for
- which branch or release channel a learner should use
- how to choose a learning path
- where to find the curriculum map and contribution docs

It should not try to hold:

- the full curriculum explanation
- detailed section breakdowns for every phase
- deep path logic
- contributor implementation details

### 2. Learning Path Guide

Primary job:

- help learners choose and follow the right route through the curriculum

It should answer:

- who each path is for
- where each path starts
- what each path must complete
- what each path may skim
- what checkpoints or projects cannot be skipped

It should not try to duplicate:

- the full section inventory
- every lesson dependency
- deep implementation details of each section

### 3. Curriculum Map

Primary job:

- present the structural map of phases, sections, and high-level outputs

It should answer:

- how the curriculum is organized
- what each phase and section is trying to accomplish
- how sections relate to one another
- where checkpoints, mini-projects, and capstones sit at a high level

It should not become:

- the main path-selection guide
- a giant learner handbook
- a duplicate of the root README

### 4. Section README

Primary job:

- act as the entry guide for one section

Every v2 section should eventually have a top-level `README.md`.

It should answer:

- what this section is for
- who should start here
- what earlier knowledge is assumed
- what sub-tracks or local directories exist
- which items are lessons, drills, exercises, checkpoints, and mini-projects
- what outputs prove the learner is ready to leave the section

This is the most important missing public navigation layer in the current repo.

### 5. Lesson Surface

Primary job:

- teach one item and point clearly to the next step

The lesson README should be the primary explanation surface for learner-facing lessons.
The runnable code file remains the primary execution surface.

Every lesson should still expose:

- what it teaches
- how to run or verify it
- why it matters
- what to do next

Practical rule:

- `README.md` should carry the deeper walkthrough first
- `main.go` should stay readable, runnable, and required

### 6. Exercise, Checkpoint, And Project README

Primary job:

- explain the rules of the artifact the learner is expected to build or complete

These items should be more explicit than ordinary lessons because they act as path validation and
milestone surfaces.

They should answer:

- what the learner must do
- prerequisites
- success criteria
- how to verify completion
- whether the item is required for all paths or only some

## Ownership Rules By Surface

Use these ownership boundaries to reduce duplication.

| Surface | Owns | Must Not Own |
| ------- | ---- | ------------ |
| `README.md` | repo purpose, quick start, branch/release channel guidance, path entry links | full curriculum explanation, detailed section breakdowns |
| `LEARNING-PATH.md` or v2 equivalent | route selection, path rules, skip rules, validation floors | full curriculum inventory, contributor internals |
| `docs/curriculum/README.md` or v2 equivalent | phase and section map, high-level structure | path-selection advice in detail |
| `NN-section/README.md` | section entry guidance, outputs, section flow | full repo overview |
| lesson README + lesson code file | item-specific teaching, walkthrough, and next step | full section map |
| exercise/checkpoint/project README | requirements, verification, milestone meaning | unrelated section-wide background |

## Required Section README Contract

Every section README should eventually include:

- section title and mission
- target learners and entry guidance
- prerequisites
- path-aware guidance for Full Path, Bridge Path, and Targeted Path
- section structure or local track map
- content inventory grouped by type
- checkpoint and mini-project expectations
- suggested next step after the section

## Path-Aware Navigation Rules

The docs system should expose path logic without duplicating content trees.

### Root README

Should only say enough to route learners:

- "new here"
- "coming from another language"
- "already working in Go"

Then link them to the learning-path guide and the relevant stable or active branch.

### Section README

Should carry the practical path-aware entry rules:

- Full Path: what to do in order
- Bridge Path: what may be skimmed and what must still be validated
- Targeted Path: what to review before entering here and what output proves readiness

### Checkpoints And Projects

Should state explicitly whether the item is:

- required for all learners
- required for faster paths as validation
- optional enrichment

This is how v2 keeps fast routes honest.

## Navigation Flow Model

The intended learner navigation flow should be:

1. root README
2. choose path
3. open section README
4. enter lesson or practice item
5. complete checkpoint or mini-project
6. move to next section or milestone

The repo should not require learners to jump back and forth between multiple overlapping index docs
to understand the next step.

## Relationship To Folder Structure

The folder structure and the docs system should match.

Because v2 section directories will eventually include a top-level `README.md`, docs navigation can
follow the same top-level section structure as the filesystem.

This means:

- section READMEs are mandatory navigation objects
- sub-track directories may use local README files when they represent a real internal track
- numbered item directories remain the canonical teaching/practice objects

## Relationship To Lesson And Exercise Specs

The docs navigation system depends on the item contracts already defined elsewhere.

Specifically:

- lessons must expose a next step
- exercises, checkpoints, and projects must expose verification rules
- path-aware rules should be declared at the section surface, not reinvented inside every lesson

## Relationship To Curriculum Schema

The metadata system should eventually support navigation with enough structure to avoid drift.

Useful schema-supported navigation fields may include:

- section mission and summary
- item type
- next items
- milestone tags
- path-critical tags
- section entry notes

The docs should remain the learner-facing explanation layer.
The schema should support validation and navigation consistency, not replace readable docs.

## Doc Style Rules

The public navigation docs should follow these style rules:

- favor short, direct language
- answer "what next" explicitly
- avoid giant prose walls when a table or short list is clearer
- keep path terminology consistent across all surfaces
- do not hide critical learner guidance inside maintainer-only docs

## First Migration Guidance

During early v2 work, do not try to rewrite every public doc at once.

Recommended order:

1. define the navigation system
2. create the section README contract
3. simplify the root README into a real front door
4. reduce overlap between the learning-path guide and curriculum map
5. update active sections incrementally during migration waves

## Definition Of Done

The docs navigation system is ready when:

- each public doc surface has a clear job
- the root README acts like a front door instead of a giant handbook
- section READMEs are defined as required navigation objects
- path-aware guidance has an explicit home
- checkpoints and projects are treated as milestone docs, not only folders
- maintainers can explain where new navigation guidance belongs without guessing

## Deferred Questions Beyond The First Prototype

These points can stay deferred beyond the first prototype:

- how much section inventory detail should live in section READMEs versus generated curriculum views
- whether the curriculum map should eventually move to a more generated navigation surface after the
  metadata system matures

## Working Recommendation

For the first v2 implementation:

- keep the root README short and routing-focused
- keep one canonical learning-path guide at the repo root during the prototype stage
- keep one canonical curriculum map and use `docs/curriculum/README.md` as that public map during
  the prototype stage
- require a section README for every section
- require local sub-track README files only when a section has multiple learner-facing internal
  tracks or milestone surfaces
- make lesson-level next-step guidance and milestone verification explicit
- reduce duplication before adding more navigation surfaces
