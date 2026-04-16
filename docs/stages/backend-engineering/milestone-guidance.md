# Backend Engineering Milestone Guidance

The `06 Backend & DB` stage currently has one promoted live proof surface.

That narrower scope is intentional for now.
The beta route is promoting the databases path first instead of pretending the entire Section `10`
inventory is already equally mature.

## Live Milestone Backbone

| ID | Surface | What It Proves |
| --- | --- | --- |
| `DB.6` | repository pattern project | you can connect persistence, transactions, and repository boundaries inside a service-shaped backend flow |

## How To Use This Milestone

### Full Path

Complete the databases track in order and use `DB.6` as the main proof surface for the stage.

### Bridge Path

If you move faster through the lessons, use `DB.6` as the non-skippable proof surface.
Do not claim the stage is done if you only skimmed the data-access lessons and skipped the project.

### Targeted Path

If you enter this stage late because your main gap is backend persistence and service boundaries,
start with the databases path first.
Use the other Section `10` directories as reference reading rather than as a substitute for the
live milestone path.

## Honest Readiness Signals

You are likely ready to leave this stage when:

- you can complete `DB.6` without copying the solution line by line
- you can explain why the repository exists instead of only saying that it works
- you can reason about transaction boundaries and context-aware execution clearly
- you can describe which Section `10` surfaces are live beta path versus reference-only

## Pressure Follow-Through

After `DB.6`, use the next pressure surface instead of treating the project as a dead end:

- [Review DB.6 repository boundary](../expert-layer/tasks/review-db6-repository-boundary.md)

## If The Milestone Feels Too Hard

That is usually a routing signal, not a failure signal.

Go back one layer:

- `DB.6` too hard: revisit connecting to the database, query execution, scanning rows, prepared
  statements, and transaction flow
