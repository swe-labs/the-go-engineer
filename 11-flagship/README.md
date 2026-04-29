# 11 Flagship

## Mission

This final stage is where the curriculum stops being a collection of lessons and becomes integrated
system proof.

By the end of this stage, a learner should be able to:

- navigate a larger integrated codebase
- understand how multiple engineering layers meet in one project
- extend a production-style service without learning a brand-new tool category mid-capstone

## Required Flagship Path

| Project | Module Path | Core Job |
| --- | --- | --- |
| `Opslane` | `OPSL.1 -> OPSL.10` | integrated SaaS backend with explicit module-by-module proof |

## Suggested Learning Flow

1. Start with [Opslane](./01-opslane).
2. Use [`MODULES.md`](./01-opslane/MODULES.md) and the Opslane progress checker to verify each
   module boundary before moving forward.
3. Bring Section 10's code-generation lessons with you as supporting tooling, not new flagship
   scope.

## Additional Flagship Projects

Stage 11 is now designed to hold more than one flagship project.

Opslane is the canonical required path today. Future flagship projects should:

- live under `11-flagship/NN-project-name/`
- use their own unique project prefix, such as `CRMX.*` or `BILL.*`
- keep their own `MODULES.md`, module READMEs, and `scripts/progress.go`

They can be added as optional flagship expansions without changing the required Opslane path.

## Finish This Stage When

- you can complete the required Opslane path through `OPSL.10`
- you can follow code across layers without losing the big picture
- you can explain how earlier production lessons show up inside one integrated system

## Next Step

This is the end of the current public learner path.
