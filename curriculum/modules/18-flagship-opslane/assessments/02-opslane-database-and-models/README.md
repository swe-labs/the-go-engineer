# Opslane database and models

## Mission

Understand and apply Opslane database and models in the context of professional Go software engineering.

## Prerequisites

- opslane-01

## Mental Model

The database and model layer is the application's memory — it defines what data exists, how it's structured, and how it persists. Just as human memory needs reliable storage and retrieval, so does application data.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, the database driver converts Go types to SQL types and back. Connection pooling uses a channel-based pool of idle connections. Migrations are tracked in a schema_migrations table so each migration runs exactly once.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/02-opslane-database-and-models
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Writing raw SQL everywhere instead of using a structured model layer.
- Ignoring database connection pooling and opening a new connection for each query.
- Not using transactions for operations that modify multiple rows or tables.

## In Production

Production Go services like CockroachDB, TimescaleDB, and Mattermost use Go structs for data modeling and database/sql for querying. Opslane follows the same proven patterns.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-03`.
