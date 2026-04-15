# Advanced Content Roadmap

This document maps advanced overlays to the current top-level architecture.

## Overview

| Folder | Scope | Advanced Overlay |
|--------|-------|------------------|
| `01-foundations` | zero-magic basics | none beyond the local README-first contract |
| `02-engineering-core` | deeper engineering choices | advanced thinking first |
| `03-backend-systems` | service and data behavior | advanced thinking, production notes, backend failures |
| `04-concurrency` | concurrent execution | advanced thinking, failure scenarios, production notes |
| `05-architecture` | design boundaries | advanced thinking |
| `06-production` | operations and reliability | production notes, failure scenarios |
| `07-advanced` | performance and expert trade-offs | all advanced overlays as needed |
| `08-projects` | integration and flagship pressure | all advanced overlays as needed |

## Delivery Map

### 01-foundations

Advanced overlays: none by default.

Keep:
- mission
- mental model
- visual model where appropriate
- machine view where appropriate
- code walkthrough
- try it
- small production relevance

Do not inject:
- custom error taxonomy
- security attack catalogs
- concurrency failure drills
- profiling and escape-analysis detail
- large-scale pressure framing

### 02-engineering-core

Start using [THINKING_SECTIONS_ADVANCED.md](./THINKING_SECTIONS_ADVANCED.md) for topics like:

- custom errors
- wrapping and translation
- interface design trade-offs
- type and boundary decisions

### 03-backend-systems

Use:

- [THINKING_SECTIONS_ADVANCED.md](./THINKING_SECTIONS_ADVANCED.md)
- [PRODUCTION_NOTES_ADVANCED.md](./PRODUCTION_NOTES_ADVANCED.md)
- [FAILURE_SCENARIOS_ADVANCED.md](./FAILURE_SCENARIOS_ADVANCED.md)

Focus areas:

- HTTP timeouts
- request bounds
- database pool behavior
- handler and service failure paths

### 04-concurrency

Use all advanced overlays for:

- goroutine leaks
- races
- worker pools
- backpressure
- shutdown behavior

### 05-architecture

Primarily use [THINKING_SECTIONS_ADVANCED.md](./THINKING_SECTIONS_ADVANCED.md) for:

- boundaries
- dependency direction
- service versus module choices

### 06-production

Use:

- [PRODUCTION_NOTES_ADVANCED.md](./PRODUCTION_NOTES_ADVANCED.md)
- [FAILURE_SCENARIOS_ADVANCED.md](./FAILURE_SCENARIOS_ADVANCED.md)

Focus areas:

- graceful shutdown
- observability
- deployment behavior
- operational failure handling

### 07-advanced

Use all advanced overlays for:

- profiling
- allocation trade-offs
- memory retention
- expert diagnostics

### 08-projects

Use all advanced overlays where the project needs:

- production realism
- system design pressure
- failure diagnosis
- integration reasoning

## Quick Reference

| Situation | Template |
|-----------|----------|
| later-stage design trade-offs | [THINKING_SECTIONS_ADVANCED.md](./THINKING_SECTIONS_ADVANCED.md) |
| production behavior and operational notes | [PRODUCTION_NOTES_ADVANCED.md](./PRODUCTION_NOTES_ADVANCED.md) |
| diagnosis and recovery drills | [FAILURE_SCENARIOS_ADVANCED.md](./FAILURE_SCENARIOS_ADVANCED.md) |
| expert/flagship assessment surfaces | [rubric-checkpoint-template.md](./rubric-checkpoint-template.md) |

## Rule

Advanced content should appear only when:

1. the learner already owns the local syntax and happy path
2. the added complexity helps judgment rather than just sounding impressive
3. the lesson still matches the stage boundary
