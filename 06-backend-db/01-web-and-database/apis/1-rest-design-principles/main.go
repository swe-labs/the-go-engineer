// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: REST Design Principles
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The philosophy of Representational State Transfer (REST).
//   - How to use HTTP verbs (GET, POST, PUT, DELETE) as actions.
//   - The discipline of resource-oriented naming.
//
// WHY THIS MATTERS:
//   - REST isn't a framework; it's a set of constraints that makes APIs
//     predictable, scalable, and easy to integrate for humans and machines.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/1-rest-design-principles
//
// KEY TAKEAWAY:
//   - Resources are Nouns. HTTP Verbs are Actions.
// ============================================================================

package main

import "fmt"

// Stage 06: APIs - REST Design Principles
//
//   - Resources as Nouns: /users, /orders, /books
//   - Uniform Interface: Consistent use of HTTP semantics
//   - Statelessness: Each request contains all information needed
//
// ENGINEERING DEPTH:
//   REST is about "Representational State Transfer". We don't call
//   remote functions; we request a "Representation" of a resource
//   (like a JSON file) and we modify the "State" by sending back
//   a new representation. This architectural style allows for
//   caching, load balancing, and independent evolution of
//   client and server.

func main() {
	fmt.Println("=== REST Design Principles ===")
	fmt.Println()

	fmt.Println("  1. RESOURCE NAMING (Nouns)")
	fmt.Println("     ❌ /getUsers         (Verb in URL)")
	fmt.Println("     ❌ /create_new_user  (Verb in URL)")
	fmt.Println("     ✅ /users            (Clean noun)")
	fmt.Println()

	fmt.Println("  2. HTTP VERBS AS ACTIONS")
	fmt.Println("     ✅ GET /users        -> List users")
	fmt.Println("     ✅ POST /users       -> Create user")
	fmt.Println("     ✅ GET /users/42     -> Get user 42")
	fmt.Println("     ✅ PUT /users/42     -> Update user 42")
	fmt.Println("     ✅ DELETE /users/42  -> Remove user 42")
	fmt.Println()

	fmt.Println("  3. CONSISTENCY IS KING")
	fmt.Println("     - Use standard HTTP status codes (200, 201, 400, 404, 500).")
	fmt.Println("     - Use consistent JSON payload structures.")
	fmt.Println("     - Treat your API as a UI for other developers.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.2 api-versioning-strategies")
	fmt.Println("Current: API.1 (rest-design-principles)")
	fmt.Println("Previous: HS.10 (rest-api-exercise)")
	fmt.Println("---------------------------------------------------")
}
