# CG.3 sqlc Workflow

## Mission

Generate typed Go query code from real SQL schema and query files so the database layer stays
explicit, reviewable, and checked at build time.

This surface is the Section 15 milestone output.

## Files

- [main.go](./main.go): lesson walkthrough for the sqlc workflow
- [schema/schema.sql](./schema/schema.sql): schema input
- [queries/query.sql](./queries/query.sql): query definitions
- [sqlc.yaml](./sqlc.yaml): generation config
- [internal/db](./internal/db): generated output target

## Run Instructions

```bash
go run ./14-code-generation/3-sqlc
```

## Success Criteria

You should be able to:

- explain why sqlc keeps SQL explicit while still generating typed Go code
- describe how schema files, query files, and `sqlc.yaml` work together
- explain why generated query code is safer than a reflection-heavy ORM in many Go services

## Next Step

After `CG.3`, return to the [Section 15 overview](../README.md) and use the remaining time to
revisit any weak milestone or rebuild a larger project with the tooling you now understand.
