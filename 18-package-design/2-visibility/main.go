// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 18: Package Design — Visibility & Export Rules
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Uppercase = exported (public), lowercase = unexported (private)
//   - The internal/ directory: compiler-enforced visibility boundaries
//   - Designing minimal public APIs
//   - Why Go doesn't have private/public/protected keywords
//
// ENGINEERING DEPTH:
//   In JVM languages, keywords like `public` and `private` are evaluated during
//   an expensive AST parsing phase just to determine visibility. Go's genius was
//   pushing visibility directly into the Lexer! By checking the Unicode uppercase
//   flag of the very first byte of a symbol's AST node, Go's compiler instantly
//   knows if a symbol is exported without running a single line of access-modifier
//   logic. This is one of the hundreds of micro-optimizations that makes the Go
//   compiler legendary for its compilation speed.
//
// RUN: go run ./18-package-design/2-visibility
// ============================================================================

// --- EXPORT RULES ---
// Go uses the simplest visibility system in any mainstream language:
//
//   Uppercase first letter → EXPORTED (visible outside the package)
//     func Connect()    ← any package can call this
//     type User struct  ← any package can reference this type
//     const MaxRetries  ← any package can read this constant
//
//   Lowercase first letter → UNEXPORTED (private to the package)
//     func validate()   ← only this package can call this
//     type config struct ← only this package can use this type
//     const maxBuffer   ← only this package can read this
//
// There are NO keywords like "public", "private", or "protected".
// The case of the first letter IS the access modifier.

// User is an EXPORTED type — other packages can use it.
// The struct fields follow the same rule:
//
//	Name  → exported (accessible from other packages)
//	email → unexported (only accessible within this package)
type User struct {
	Name  string // Exported — other packages can read/write this
	email string // Unexported — hidden from other packages
}

// NewUser is an EXPORTED constructor function.
// This is Go's alternative to constructors in OOP languages.
// By convention, functions that create types are named New<Type>.
func NewUser(name, email string) User {
	return User{
		Name:  name,
		email: sanitizeEmail(email), // Uses unexported helper internally
	}
}

// Email is an EXPORTED accessor (getter) for the unexported email field.
// This pattern gives you control over how the field is accessed:
//   - You can add validation
//   - You can add logging
//   - You can change the internal representation without breaking callers
func (u User) Email() string {
	return u.email
}

// sanitizeEmail is an UNEXPORTED helper — implementation detail.
// Other packages cannot call this directly. If we later change
// the sanitization logic, no external code breaks.
func sanitizeEmail(email string) string {
	// In production, you'd do actual validation here
	return email
}

func main() {
	fmt.Println("=== Visibility & Export Rules ===")
	fmt.Println()

	// Create a user via the exported constructor
	user := NewUser("Alice", "alice@example.com")

	// Accessing exported field directly
	fmt.Printf("  Name (exported):  %s\n", user.Name)

	// Accessing unexported field via exported getter
	fmt.Printf("  Email (via getter): %s\n", user.Email())

	// This would NOT compile from another package:
	// fmt.Println(user.email)     ← COMPILE ERROR: email is unexported
	// sanitizeEmail("test")       ← COMPILE ERROR: sanitizeEmail is unexported

	fmt.Println()

	// --- THE internal/ DIRECTORY ---
	// Go has a special compiler-enforced rule for directories named "internal/":
	//
	//   myapp/
	//   ├── cmd/server/         ← Can import myapp/internal/...
	//   ├── internal/
	//   │   ├── auth/           ← Only importable by myapp/* and children
	//   │   └── db/             ← NOT importable by external packages
	//   └── pkg/api/            ← Can import myapp/internal/...
	//
	// Any package outside the parent of "internal/" gets a COMPILE ERROR
	// if it tries to import an internal package.
	//
	// Use internal/ for:
	//   - Database queries and schema details
	//   - Business logic that shouldn't be a public API
	//   - Helper utilities specific to your application
	fmt.Println("=== The internal/ Directory ===")
	fmt.Println("  myapp/internal/auth → importable by myapp/* only")
	fmt.Println("  myapp/internal/db   → NOT importable by other modules")
	fmt.Println("  This is COMPILER-ENFORCED, not just convention!")
	fmt.Println()

	// --- MINIMAL API SURFACE ---
	fmt.Println("=== Design Principle: Export the Minimum ===")
	fmt.Println("  1. Start with everything unexported")
	fmt.Println("  2. Export ONLY what external packages need")
	fmt.Println("  3. Use constructors (NewX) to validate on creation")
	fmt.Println("  4. Use getters for controlled access to internal state")
	fmt.Println("  5. Once exported, it's part of your API — removing it is a breaking change")
}
