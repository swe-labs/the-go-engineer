# Templates

This folder contains reusable curriculum-reference templates.

The key rule is scope:

- early learner stages (`01` through `04`) keep the small README-first lesson contract
- advanced thinking, failure, and production overlays are references for later stages
- expert and flagship assessment surfaces still use the shared [rubric-checkpoint-template.md](./rubric-checkpoint-template.md)

## Template Overview

| Template | Purpose | Use It In |
|----------|---------|-----------|
| [rubric-checkpoint-template.md](./rubric-checkpoint-template.md) | Expert and flagship assessment surfaces | Expert layer and flagship project |
| [THINKING_SECTIONS_ADVANCED.md](./THINKING_SECTIONS_ADVANCED.md) | Later-stage design and trade-off prompts | `02-engineering-core` and above |
| [PRODUCTION_NOTES_ADVANCED.md](./PRODUCTION_NOTES_ADVANCED.md) | Production-hardening overlays | `03-backend-systems`, `04-concurrency`, `06-production`, `07-advanced`, `08-projects` |
| [FAILURE_SCENARIOS_ADVANCED.md](./FAILURE_SCENARIOS_ADVANCED.md) | Diagnosis and recovery drills | `04-concurrency` and above |
| [ADVANCED_CONTENT_ROADMAP.md](./ADVANCED_CONTENT_ROADMAP.md) | Delivery map for advanced overlays | Maintainer reference |

## Foundations Rule

For the early learner stages (`01` through `04`), use the local lesson contract only:

- mission
- mental model where appropriate
- visual model where appropriate
- machine view where appropriate
- code walkthrough
- try it
- small production relevance

Do not force in:

- custom error taxonomies
- SQL injection or XSS exercises
- concurrency failure drills
- profiling and escape-analysis detail
- "10k requests per second" pressure framing

## Quick Decision Guide

```text
Is this lesson in stages 01-04?
  Yes -> keep the basic README-first lesson contract only
  No  -> is the lesson mainly about design trade-offs?
           Yes -> use THINKING_SECTIONS_ADVANCED
           No  -> is it mainly about production behavior?
                    Yes -> use PRODUCTION_NOTES_ADVANCED
                    No  -> is it mainly about failures and diagnosis?
                             Yes -> use FAILURE_SCENARIOS_ADVANCED
                             No  -> stay with the local lesson contract
```

## Guidance

1. Match the template to the stage.
2. Treat advanced templates as references, not copy-paste obligations.
3. Keep early lessons readable before trying to make them exhaustive.
4. Use the roadmap when a topic could fit in more than one place.

## Next Reference

Use [ADVANCED_CONTENT_ROADMAP.md](./ADVANCED_CONTENT_ROADMAP.md) when deciding when a later-stage concept should appear.
