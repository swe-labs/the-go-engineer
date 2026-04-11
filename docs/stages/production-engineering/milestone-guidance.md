# Production Engineering Milestone Guidance

## What Counts As Stage Completion

You should be able to explain how a service is observed, how it stops safely, and how packaging
choices affect runtime behavior.

## Milestones

### `SL.5` PII redactor exercise

This proves that learners understand logging as policy-enforced operational data, not as a stream
of ad hoc debug strings.

### `GS.3` shutdown capstone

This proves that lifecycle behavior can be designed deliberately instead of being left to process
termination luck.

## Bridge Path Check

If you are moving quickly through this stage, make sure you can still explain:

- why handler-based logging and context propagation matter operationally
- why shutdown order affects correctness during deploys
- why deployment packaging changes operational behavior instead of merely changing delivery format

## Ready To Move On

Move to [9 Expert Layer](../09-expert-layer.md) or [10 Flagship Project](../10-flagship-project.md)
once the current stage path feels operationally understandable instead of mysterious.
