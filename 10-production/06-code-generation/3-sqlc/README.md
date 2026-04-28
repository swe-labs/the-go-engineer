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


## 



## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Code 



## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.


## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

## Code Walkthrough

We step through the code sequentially, examining how the interfaces are satisfied, where the errors are checked, and how the core loop manages control flow.

## In Production

Database interaction in Go usually forces a painful choice: use an ORM that relies heavily on reflection and hides performance characteristics, or write raw SQL strings and risk runtime crashes when query results don't match Go structs. Tools like `sqlc` represent the idiomatic Go compromise for production services. By writing raw SQL and generating typed Go code at build time, teams get the exact performance they expect from their queries, complete visibility into the executed SQL, and compile-time safety. If a developer renames a database column, `sqlc` will fail to generate the Go code (or the generated code will fail to compile against the rest of the app), catching the regression before it ever reaches a staging environment. In large codebases, this build-time safety net is essential for refactoring database schemas confidently.

## Thinking Questions

1. Why does `sqlc` parse the schema files along with the query files when generating code?
2. What are the trade-offs between a code-generation tool like `sqlc` and a reflection-based ORM like `gorm` when it comes to performance and debuggability?
3. If you need to build a highly dynamic query (e.g., an advanced search filter where 10 different columns are optional), why might `sqlc` be less suitable than a query builder?
4. How does having the generated database access code in its own isolated package (`internal/db`) help enforce a clean architecture?

## Next Step

After `CG.3`, continue to [11 Flagship](../../../11-flagship) and use Opslane to
apply the tooling you now understand inside a larger integrated service.


