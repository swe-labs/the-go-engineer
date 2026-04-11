# 8 Production Engineering

## Purpose

`8 Production Engineering` turns built systems into operated systems.

## Who This Is For

- learners who can already build and structure meaningful applications
- developers who want stronger runtime and operations instincts

## Mental Model

Software is not finished when it runs once on your machine.
This stage teaches what changes when software must be observed, configured, shut down safely, and
deployed with less guesswork.

## Why This Stage Exists

This stage exists because "it works on my machine" is not the same thing as "this system can be
operated safely."

The goal is to help learners think about logs, shutdown behavior, packaging, and deployment as part
of correctness instead of late operational cleanup.

## What You Should Learn Here

- structured logging
- graceful shutdown
- deployment packaging
- configuration boundaries
- early observability and runtime support thinking

## Stage Shape

This stage currently has two live public paths plus one deployment-oriented reference surface:

1. `structured-logging`
   - the live beta path for log shape, request context, handlers, and redaction
2. `graceful-shutdown`
   - the live beta path for service lifecycle handling and drain order
3. `docker-and-deployment`
   - a production reference surface for packaging and deployment workflow

That means the stage is honest about where proof work happens now while still exposing deployment
material that matters for later flagship and production growth.

## Current Source Content

- [14-application-architecture/structured-logging](../../14-application-architecture/structured-logging/)
- [14-application-architecture/graceful-shutdown](../../14-application-architecture/graceful-shutdown/)
- [14-application-architecture/docker-and-deployment](../../14-application-architecture/docker-and-deployment/)

## Stage Support Docs

Use these support docs when you want the beta-stage view without digging through all of Section
`14`:

- [Production Engineering support index](./production-engineering/README.md)
- [Stage map](./production-engineering/stage-map.md)
- [Milestone guidance](./production-engineering/milestone-guidance.md)

## Where This Stage Starts

Start with [14-application-architecture/structured-logging](../../14-application-architecture/structured-logging/).

That is the most accessible public entry because it teaches learners how to think about runtime
signals from a running system before they move into shutdown and deployment flow.

## Recommended Order

Use this order for the current beta-facing path:

1. complete the structured-logging track from `SL.1` through `SL.5`
2. complete the graceful-shutdown track from `GS.1` through `GS.3`
3. use `docker-and-deployment` as the next production reference surface for packaging and runtime
   rollout thinking

## Path Guidance

### Full Path

Complete structured logging first, then graceful shutdown, then deployment reading.
This keeps the stage ordered around runtime visibility, safe stop behavior, and finally deployment
packaging.

### Bridge Path

You can move faster if observability and service lifecycle concepts already feel familiar, but do
not skip:

- `SL.1`
- `SL.2`
- `SL.5`
- `GS.1`
- `GS.3`

Those are the main proof surfaces that show you can reason about operated systems instead of only
running programs locally.

### Targeted Path

If you enter this stage with a narrow goal:

- start with `structured-logging` if your gap is runtime visibility
- start with `graceful-shutdown` if your gap is lifecycle safety during deploys
- use `docker-and-deployment` when your gap is packaging and container workflow

## Stage Milestones

The current live milestone backbone is:

- `SL.5` PII redactor exercise
- `GS.3` shutdown capstone

`docker-and-deployment` is an important production surface, but it is currently a reference layer
instead of the main public beta proof path.

## Finish This Stage When

- you can explain what good logs are for
- you know how to stop services safely
- you can package and run a deployment-oriented application flow
- you can reason about runtime behavior, not just implementation details

More concretely:

- you can explain why structured logs are operational data instead of prettier strings
- you can coordinate readiness, request drain, worker drain, and shutdown order
- you understand why deployment packaging changes the operational shape of a program
- you know which production surfaces are the current beta path and which remain reference material

## Next Stage

Move to [9 Expert Layer](./09-expert-layer.md) or [10 Flagship Project](./10-flagship-project.md),
depending on whether you want more pressure first or a longer integrated build path first.
