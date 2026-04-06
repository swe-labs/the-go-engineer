// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./15-code-generation/3-sqlc
package main

import (
	"fmt"
	"log/slog"
	"os"
)

// ============================================================================
// Section 15: Code Generation — SQLC
// Level: Expert
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Generating type-safe Go code from raw SQL queries
//   - Using sqlc.yaml for configuration
//   - Schema-driven development vs. ORMs
//   - Interface generation for easier testing
//
// ENGINEERING DEPTH:
//   ORMs (like GORM) reflect on your structs at runtime to generate SQL. This
//   is flexible but slow and prone to "magic" errors. SQLC takes the opposite
//   approach: you write the SQL first, and it generates the Go structs and
//   query methods at compile-time.
//
//   BENEFITS:
//     1. No Reflection: Zero runtime overhead for query mapping.
//     2. Type Safety: If you rename a column in SQL, sqlc fails at build time.
//     3. Full SQL Power: No need to learn a DSL; just write standard Postgres/MySQL.
//
// RUN: go run ./15-code-generation/3-sqlc
// ============================================================================

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("SQLC Lesson")
	fmt.Println("\nTo generate Go code from your SQL schema and queries, run:")
	fmt.Println("  sqlc generate")

	fmt.Println("\nSQLC parses the files in:")
	fmt.Println("  - schema/schema.sql (The data definition)")
	fmt.Println("  - queries/query.sql (The CRUD operations)")

	fmt.Println("\nIt then outputs type-safe Go code into internal/db/ folder.")

	fmt.Println("\nOnce generated, you can use the Querier interface like this:")
	fmt.Println(`
    func GetUserHandler(q db.Querier, id int64) {
        user, err := q.GetUser(context.Background(), id)
        if err != nil {
            return
        }
        fmt.Printf("User: [name]\n", user.Name)
    }`)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 CURRICULUM COMPLETE: You've reached the end!")
	fmt.Println("   Explore the Enterpise Capstone (Section 14) to see all lessons combined.")
	fmt.Println("---------------------------------------------------")
}
