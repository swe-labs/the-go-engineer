# Backend Engineering Stage Map

This stage turns language, type, and I/O fluency into service-shaped application work.

It starts with the live databases path, while keeping the rest of the current Section `10` source
tree available as reference material.

## Stage Flow

1. `databases`
   - source: [09-web-and-database/databases](../../../09-web-and-database/databases/)
   - core job: teach `database/sql`, transactions, repository boundaries, and context-aware
     persistence
   - milestone: `DB.6` repository pattern project
2. `http-client`, `database-migrations`, and `web-masterclass`
   - source surfaces:
     - [09-web-and-database/http-client](../../../09-web-and-database/http-client/)
     - [09-web-and-database/database-migrations](../../../09-web-and-database/database-migrations/)
     - [09-web-and-database/web-masterclass](../../../09-web-and-database/web-masterclass/)
   - core job: provide broader backend reference material while the live beta route is still
     focused on the databases path first

## What Each Part Adds

### `databases`

This is where the learner stops treating persistence as hidden plumbing and starts reasoning about
queries, repositories, transactions, and data safety directly.

### `reference surfaces`

These directories widen the backend picture around HTTP and migrations without pretending they are
already the main public beta learner route.

## Recommended Full-Path Order

1. Finish the databases track from `DB.1` through `DB.6`.
2. Use the other Section `10` surfaces as reinforcement after the databases milestone path is
   stable.

## Bridge-Path Reminder

If SQL basics or service-style code already feel somewhat familiar, you can move faster through the
early repetition.
What you should not skip is proof:

- `DB.1`
- `DB.3`
- `DB.5`
- `DB.6`

## Exit Condition

You are ready for `07 Concurrency` when you can finish the live databases milestone path
honestly and explain how request flow, persistence boundaries, and transaction safety connect in a
real service-style program.
