# Content Quality Standard

This document defines the quality bar for learner-facing curriculum content.

The curriculum should be useful to a beginner without becoming shallow for an experienced engineer.

## Core principle

No magic.

Every important idea must be explained through:

1. the problem it solves
2. the mental model
3. the machine/runtime view
4. runnable examples
5. failure modes
6. debugging signals
7. production relevance
8. practice
9. proof of understanding

## Module README standard

A module README must include:

```text
# Module NN — Title
## Mission
## Who this module is for
## What you will build
## Prerequisites
## Concept map
## Lessons
## Labs
## Projects
## Assessments
## Common failure modes
## Completion checklist
## Next module
```

A module README should help the learner answer:

- Why does this module exist?
- What will I be able to do after finishing it?
- What should I already know?
- Which lessons are required?
- What proof do I need before moving on?

## Lesson README standard

A lesson README must include:

```text
# Lesson Title
## Mission
## Prerequisites
## Mental Model
## Visual Model
## Machine View
## Run Instructions
## Code Walkthrough
## Try It
## Common Mistakes
## Debugging Signals
## In Production
## Performance Notes
## Security Notes
## Thinking Questions
## Proof of Understanding
## Next Step
```

For conceptual non-code lessons, `Run Instructions` may explain that there is no code execution, but the lesson still needs an observable proof task.

## Lab README standard

A lab is guided practice. It should include:

```text
# Lab Title
## Goal
## Scenario
## Starting Point
## Tasks
## Hints
## Verification
## Reflection
## Extension
```

Labs must not become vague prompts. They should be specific enough that a learner can start without guessing.

## Project README standard

A project proves integrated ability.

A project README must include:

```text
# Project Title
## Mission
## Real-world scenario
## Requirements
## Constraints
## Architecture
## Deliverables
## Starter instructions
## Implementation milestones
## Testing requirements
## Operational requirements
## Security requirements
## Rubric
## Submission checklist
## Portfolio guidance
```

Projects must produce a meaningful artifact, not just "write some code."

## Assessment standard

An assessment must include:

```text
# Assessment Title
## Purpose
## Scope
## Time expectation
## Allowed resources
## Questions or tasks
## Evidence required
## Scoring rubric
## Passing standard
## Retake policy
## Feedback guide
```

Assessments must measure understanding, not memorization.

## README writing style

Use:

- short paragraphs
- concrete examples
- clear headings
- second person for learner instructions
- precise terms after they are introduced
- analogies only when they are technically accurate
- callouts for traps, not decorative notes

Avoid:

- unexplained jargon
- vague encouragement
- long walls of prose
- "just", "simply", or "obviously"
- code snippets without explaining why they exist
- forward references without naming where the concept is taught

## Code quality standard

Lesson code must be:

- runnable
- formatted with `gofmt`
- small enough to study
- realistic enough to matter
- commented around non-obvious behavior
- tested when behavior should be proven

Avoid code that teaches bad habits for convenience.

Examples of invalid shortcuts:

- HTTP servers without timeouts in production-shaped lessons
- SQL string concatenation with user input
- goroutines without lifecycle control
- tests that only check that code runs
- ignoring errors to keep examples short

## Starter and solution standard

For implementation lessons, labs, and projects:

- `_starter/` must compile unless the exercise intentionally asks the learner to fix compile errors.
- `_solution/` must include a complete reference implementation.
- starter and solution must use the same package/module shape.
- tests should be runnable against both when appropriate.
- README instructions must tell learners which folder to use.

## Visual model standard

A visual model may be:

- SVG diagram
- Mermaid diagram
- ASCII diagram
- table
- execution timeline
- request lifecycle diagram
- memory layout diagram

It must clarify a concept, not decorate the page.

Every diagram should answer one of these:

- What moves where?
- What owns what?
- What happens first, second, third?
- What boundary is crossed?
- What fails if the design is wrong?

## Common mistakes standard

Each technical lesson should list common mistakes with fixes.

Good format:

```text
Mistake: forgetting to close rows.
Why it happens: the query succeeds, so cleanup feels optional.
Symptom: connection pool exhaustion under load.
Fix: defer rows.Close() immediately after checking the query error.
```

## Production context standard

Every lesson from backend modules onward must explain:

- where the idea appears in real systems
- what can fail
- how teams debug it
- what tradeoff the design introduces

## Completion checklist

A content item is complete when:

- metadata is stable
- README exists
- required sections exist in order
- no placeholder language remains
- code compiles when present
- tests pass when present
- starter/solution exist when required
- assets exist when referenced
- assessment/project proof exists
- validator passes
