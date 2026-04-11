# Alpha Source Inventory vs Beta Public Architecture

The Go Engineer is in a transition period.

That means two things are true at the same time:

- the repo still contains the alpha-era section inventory as the main source content
- the public learner-facing navigation is now moving to the beta stage model

This document explains how those two layers fit together.

## The Short Version

Use this rule:

- if you want to know **how to navigate the curriculum as a learner**, follow the beta stage docs
- if you want to know **where the current source content physically lives**, follow the alpha source
  inventory

Beta is the navigation truth.
Alpha is still the physical source inventory for much of the repo.

## What Alpha Means Right Now

In this repo, `alpha` refers to the current section-based content inventory:

- `01-core-foundations`
- `02-control-flow`
- `03-data-structures`
- ...
- `15-code-generation`

Alpha proved that the curriculum could be migrated, validated, and cleaned up section by section.

Alpha is still valuable because it contains the lessons, exercises, milestone projects, and source
files learners use today.

## What Beta Means Right Now

`Beta` means the learner-facing curriculum architecture is being reorganized around engineering
stages instead of only section numbers.

The beta stages are:

1. `0 Foundation`
2. `1 Language Fundamentals`
3. `2 Types and Design`
4. `3 Modules and IO`
5. `4 Backend Engineering`
6. `5 Concurrency System`
7. `6 Quality and Performance`
8. `7 Architecture`
9. `8 Production Engineering`
10. `9 Expert Layer`
11. `10 Flagship Project`
12. `11 Code Generation`

These stages are the public learner-routing model.

## Why The Repo Uses Both At Once

We are not discarding alpha content and rewriting the entire repo from scratch in one destructive
move.

Instead, beta does three things:

1. regroup existing alpha content into clearer learner-facing stages
2. split some alpha sections where the learner-facing boundary is better than the current folder
   boundary
3. add a few new beta-only layers where alpha was too thin, especially:
   - `0 Foundation`
   - `8 Production Engineering`
   - `9 Expert Layer`
   - `10 Flagship Project`

This lets the curriculum improve without pretending the current folder tree has already been
perfectly rebuilt.

## What A Split Looks Like

Some alpha sections map cleanly to one beta stage.
Some do not.

Examples:

- `10-web-and-database` mostly feeds `4 Backend Engineering`
- `15-code-generation` maps directly to `11 Code Generation`
- `01-core-foundations` splits across:
  - `0 Foundation`
  - `1 Language Fundamentals`
- `14-application-architecture` splits across:
  - `7 Architecture`
  - `8 Production Engineering`
  - `10 Flagship Project`

So if a learner sees one alpha section feeding more than one beta stage, that is intentional.
It is not duplication by accident.

## What To Trust For What

Use these docs based on the question you are asking:

### If you are a learner

Use:

- [README.md](../README.md)
- [LEARNING-PATH.md](../LEARNING-PATH.md)
- [docs/stages/README.md](./stages/README.md)

These docs tell you how to move through the curriculum now.

### If you want the physical source layout

Use:

- [docs/curriculum/README.md](./curriculum/README.md)
- the current top-level section folders

These show where the source content currently lives.

### If you are a maintainer or contributor

Use:

- the beta shell docs for public navigation truth
- the source inventory docs for where files live today
- `planning/v2` for frozen planning and regrouping direction

## What This Does Not Mean

This transition model does **not** mean:

- the repo has two different public curricula
- learners should ignore the stage model
- the alpha folder structure is the final learner-facing architecture
- every section must be physically moved before beta guidance can be honest

It means the public navigation model is changing first, and the physical regrouping work continues
incrementally.

## Practical Learner Rule

If you are unsure what to do next:

1. start from the beta stage docs
2. choose the stage that matches your current goal
3. follow the linked source sections from there
4. use the alpha inventory only when you need to understand where files live physically

## Bottom Line

The Go Engineer is not trying to keep alpha and beta as competing systems.

The alpha section inventory is still the source content.
The beta stage model is the public learner-facing architecture.

That is the transition model.
