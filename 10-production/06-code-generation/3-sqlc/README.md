# CG.3 sqlc Workflow

## Mission

Generate typed Go query code from real SQL schema and query files so the database layer stays
explicit, reviewable, and checked at build time.

This surface is the Stage 10 code-generation milestone output.

## Files

- [main.go](./main.go): lesson walkthrough for the sqlc workflow
- [schema/schema.sql](./schema/schema.sql): schema input
- [queries/query.sql](./queries/query.sql): query definitions
- [sqlc.yaml](./sqlc.yaml): generation config
- `internal/db`: generated output target configured in `sqlc.yaml`

## Run Instructions

```bash
go run ./10-production/06-code-generation/3-sqlc
```

## Success Criteria

You should be able to:

- explain why sqlc keeps SQL explicit while still generating typed Go code
- describe how schema files, query files, and `sqlc.yaml` work together
- explain why generated query code is safer than a reflection-heavy ORM in many Go services

## Next Step

After `CG.3`, continue to [11 Flagship](../../../11-flagship) and use the enterprise capstone to
apply the tooling you now understand inside a larger integrated service.
