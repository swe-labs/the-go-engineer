# Code Generation Milestone Guidance

## What Counts As Stage Completion

You should be able to explain why a generation workflow exists, how it fits into the build and
review process, and when it should not be used.

## Milestone

### `CG.3` sqlc workflow

This proves that generation can support production data-access code without turning the workflow
into hidden runtime magic.

## Bridge Path Check

If you are moving quickly through this stage, make sure you can still explain:

- why `go generate` is a build-time workflow
- why generated mocks and query code still need review
- why generation should reduce manual repetition without hiding the real system model

## Ready To Use

At this point, generation should feel like leverage you understand, not magic you tolerate.
