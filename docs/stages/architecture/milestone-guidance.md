# Architecture Milestone Guidance

## What Counts As Stage Completion

You should be able to explain why a system is split the way it is, not just say that the folders
look organized.

## Milestone

### `PD.3` project layout

This proves that you can move from naming and visibility rules into a layout that matches the real
shape of a growing Go system.

## Bridge Path Check

If you are moving quickly through this stage, make sure you can still explain:

- why package names should represent a domain or responsibility
- why `internal/` is a real boundary and not a decorative folder
- why a project layout should grow with the system instead of being over-designed on day one
- why service contracts and code generation still need architectural judgment

## Ready To Move On

Move to [10 Production](../10-production.md) when package boundaries and
service seams feel understandable instead of magical.
