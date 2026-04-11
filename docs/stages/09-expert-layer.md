# 9 Expert Layer

## Purpose

`9 Expert Layer` creates engineering pressure through review, diagnosis, and trade-off analysis.

## Who This Is For

- learners who can already build, test, structure, and operate meaningful systems
- developers who want stronger review and judgment skills instead of more syntax

## Mental Model

This stage is about judgment under pressure.
The work shifts from "can I implement this?" to "can I critique, diagnose, and improve this under
constraints?"

## Why This Stage Exists

This stage exists because strong engineers are not defined only by the code they can write from a
blank file.

They are also defined by the quality of their judgment when reviewing, debugging, redesigning, and
making trade-offs under pressure.

## What You Should Learn Here

- code review and critique
- anti-pattern detection
- failure analysis
- redesign under constraints
- trade-off reasoning across correctness, simplicity, speed, and operability

## Stage Shape

This stage is intentionally built from beta shells and source seeds rather than from one dedicated
alpha section.

The current public shape is:

1. review and diagnosis tasks derived from later-stage milestone surfaces
2. anti-pattern and redesign tasks derived from backend, concurrency, and production work
3. rubric-heavy pressure work that prepares learners for the flagship project

That means this stage is about engineering pressure and judgment, not about introducing another
large topic inventory.

## Current Source Content

This stage does not yet have a dedicated public source folder on `main`.

## Stage Support Docs

Use these support docs when you want the beta-stage view of Expert Layer:

- [Expert Layer support index](./expert-layer/README.md)
- [Source seeds and ownership](./expert-layer/source-seeds.md)
- [Stage map](./expert-layer/stage-map.md)
- [Pressure guidance](./expert-layer/pressure-guidance.md)

## Where This Stage Starts

This stage starts after the learner can already build, test, structure, and operate meaningful Go
systems.

In practice, the first entry surfaces should be diagnosis and review tasks derived from the later
beta milestones, not new syntax lessons.

## Recommended Order

Use this order for the current beta-facing shell:

1. start with review and diagnosis tasks derived from completed milestone work
2. move into redesign and anti-pattern tasks that force trade-off explanation
3. use rubric-heavy pressure work as the handoff into the flagship project

## Path Guidance

### Full Path

Enter this stage only after the earlier build/test/operate stages feel solid.
The goal here is not new breadth.
It is stronger judgment.

### Bridge Path

If you already work professionally, you may move into this stage earlier, but only if you can
already explain:

- why a design is weak, not just that it feels odd
- how to diagnose a failure with evidence
- how to defend a trade-off between correctness, simplicity, performance, and operability

### Targeted Path

If you enter this stage with a narrow goal:

- choose review tasks if your gap is critique quality
- choose diagnosis tasks if your gap is debugging and failure analysis
- choose redesign tasks if your gap is trade-off reasoning

## Current Public Output Shape

This stage currently publishes guidance and source seeds instead of one fixed milestone backbone.

That is intentional.
The public beta goal is to define the pressure model clearly before it gets expanded into a larger
task bank.

## Finish This Stage When

- you can explain why code is weak, not just say that it feels wrong
- you can diagnose failures with evidence instead of guesswork
- you can defend trade-offs and redesign choices clearly
- you can review unfamiliar code and still identify the real risks

More concretely:

- you can review a real code surface and identify the most meaningful risks first
- you can explain how you would redesign a weak surface under explicit constraints
- you can connect review comments back to correctness, operability, and maintainability
- you are ready to carry that judgment into a longer flagship build path

## Next Stage

Move to [10 Flagship Project](./10-flagship-project.md).
