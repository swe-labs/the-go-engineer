# Section 21: Database Migrations

Welcome to **Schema Evolution**.

When interacting with a database locally (Stage 07), it's easy to just run `CREATE TABLE` inside your SQL IDE. But what happens when you deploy to production? What happens when you add a new `age` column to the `users` table next week?

You cannot manually run SQL queries across production databases. You need **Migrations**.

## What is golang-migrate?

`golang-migrate` is the industry-standard CLI and Go library for tracking database schemas.
It uses `UP` scripts (to apply changes) and `DOWN` scripts (to revert changes).

## In This Section

Look inside `1-embedded-migrations`. We will do something incredible:
We will take raw `.sql` migration files and use Go's `//go:embed` directive (which we learned in Stage 06) to compile the SQL files natively into the binary.

Then, when your Go server boots up, it will execute the migrations **automatically** before opening the HTTP port.

### How to use golang-migrate CLI locally (Optional)

To create new migration files locally (like `000001_create_users.up.sql`), you install the CLI:
`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

Usage:
`migrate create -ext sql -dir db/migrations -seq create_users_table`
