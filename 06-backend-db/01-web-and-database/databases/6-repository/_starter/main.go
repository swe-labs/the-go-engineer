// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Stage 06: Databases - Repository Pattern Project (Exercise Starter)
// Level: Advanced
// ============================================================================
//
// EXERCISE: Build a SQLite-backed User Repository
//
// REQUIREMENTS:
//  1. [ ] Define `User` and `Profile` domain models for the repository contract
//  2. [ ] Define a `UserRepository` interface with create/read methods
//  3. [ ] Implement a concrete repository backed by `*sql.DB`
//  4. [ ] Use a transaction to create the user row and profile row together
//  5. [ ] Read users back with `QueryRow` / `Query` and `Scan`
//  6. [ ] Keep SQL parameterized instead of building query strings manually
//  7. [ ] Prove the contract with tests before comparing against the solution
//
// HINTS:
//   - A repository boundary should hide SQL details from callers, not hide the data model
//   - `defer tx.Rollback()` is the safety net for multi-step writes
//   - `rows.Close()` and `rows.Err()` are part of the read-side contract
//   - Keep the starter focused on one repository, not a full web app
//
// RUN: go run ./06-backend-db/01-web-and-database/databases/6-repository/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define your domain models here.

// TODO: Define the UserRepository interface.

// TODO: Define the concrete SQL repository and constructor.

// TODO: Implement create/read methods with transactions and scans.

func main() {
	fmt.Println("=== Repository Pattern Project Starter ===")
	fmt.Println()
	fmt.Println("TODO: implement the repository contract and database workflow.")
	fmt.Println("Use the REQUIREMENTS above as your checklist.")
	fmt.Println()
	fmt.Println("When finished, compare your work with ../main.go and run the tests.")
}
