// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: sqlc Workflow
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Generate typed Go query code from SQL schema and query files
//   - Using sqlc.yaml for configuration
//   - Schema-driven development vs ORMs
//   - Interface generation for easier testing
//
// WHY THIS MATTERS:
//   - ORMs reflect at runtime (flexible but slow, magic errors).
//   - sqlc generates at compile time: zero runtime overhead, fails fast.
//
// RUN:
//   go run ./10-production/06-code-generation/3-sqlc
// KEY TAKEAWAY:
//   - Write SQL first, sqlc generates type-safe Go at build time.
//   - Column rename fails at build, not in production.
// ============================================================================

package main

import (
	"fmt"
	"log/slog"
	"os"
)

// Stage 10: Code Generation - sqlc Workflow
//
//   - Generating type-safe Go code from raw SQL queries
//   - Using sqlc.yaml for configuration
//   - Schema-driven development vs. ORMs
//   - Interface generation for easier testing
//
// ENGINEERING DEPTH:
//   ORMs (like GORM) reflect on your structs at runtime to generate SQL. This
//   is flexible but slow and prone to magic errors. sqlc takes the opposite
//   approach: you write the SQL first, and it generates the Go structs and
//   query methods at compile time.
//
//   BENEFITS:
//     1. No Reflection: Zero runtime overhead for query mapping.
//     2. Type Safety: If you rename a column in SQL, sqlc fails at build time.
//     3. Full SQL Power: No need to learn a DSL; just write standard SQL.
//

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("sqlc Workflow")
	fmt.Println("\nTo generate Go code from your SQL schema and queries, run:")
	fmt.Println("  sqlc generate")

	fmt.Println("\nsqlc parses the files in:")
	fmt.Println("  - schema/schema.sql (the data definition)")
	fmt.Println("  - queries/query.sql (the CRUD operations)")

	fmt.Println("\nIt then outputs type-safe Go code into internal/db/.")

	fmt.Println("\nOnce generated, you can use the Querier interface like this:")
	_, _ = os.Stdout.WriteString(`
    func GetUserHandler(q db.Querier, id int64) {
        user, err := q.GetUser(context.Background(), id)
        if err != nil {
            return
        }
        fmt.Printf("User: %s\n", user.Name)
    }`)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: OPSL.1 -> 11-flagship/01-opslane/modules/01-foundation")
	fmt.Println("   Current: CG.3 (sqlc workflow)")
	fmt.Println("---------------------------------------------------")
}
