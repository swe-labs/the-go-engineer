// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 14: Application Architecture - Package Design: Naming Conventions
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Go's package naming conventions
//   - How names compose: package.Function reads as natural English
//   - Common naming mistakes and how to avoid them
//   - Real examples from the standard library
//
// ENGINEERING DEPTH:
//   Go takes an extreme minimalist approach to naming to reduce mental load.
//   There are no "namespaces" in Go - the package name is the namespace.
//   When compiling, the Go linker builds symbol tables using the exact package name.
//   This is why `utils` or `helpers` are architectural failures: a symbol like
//   `utils.Format()` tells the compiler and the human nothing about the domain
//   boundary of the code it operates on.
//
// RUN: go run ./14-application-architecture/package-design/1-naming
// ============================================================================

func main() {
	fmt.Println("=== Package Naming Conventions ===")
	fmt.Println()

	rules := []struct {
		rule string
		good string
		bad  string
		why  string
	}{
		{
			rule: "1. Use short, lowercase, one-word names",
			good: "auth, store, email, user",
			bad:  "authorization, dataStore, emailService, userManager",
			why:  "Short names compose better: auth.Login() vs authorization.LoginUser()",
		},
		{
			rule: "2. No stuttering (don't repeat the package name)",
			good: "http.Client (not http.HTTPClient)",
			bad:  "http.HTTPClient, user.UserService",
			why:  "The package name already provides the namespace",
		},
		{
			rule: "3. No utility/helper/common packages",
			good: "strings.Split(), path.Join()",
			bad:  "utils.SplitString(), helpers.JoinPath()",
			why:  "'utils' says nothing about what's inside - it becomes a junk drawer",
		},
		{
			rule: "4. Name by what it provides, not how it works",
			good: "store (provides storage), auth (provides authentication)",
			bad:  "postgres (implementation), oauth2 (protocol)",
			why:  "Implementation names leak details and make swapping harder",
		},
		{
			rule: "5. Prefer singular package names",
			good: "model, handler, middleware",
			bad:  "models, handlers, middlewares",
			why:  "Go convention: model.User reads better than models.User",
		},
	}

	for _, r := range rules {
		fmt.Printf("  📏 %s\n", r.rule)
		fmt.Printf("     ✅ Good: %s\n", r.good)
		fmt.Printf("     ❌ Bad:  %s\n", r.bad)
		fmt.Printf("     💡 Why:  %s\n\n", r.why)
	}

	fmt.Println("=== Standard Library Naming Excellence ===")
	examples := []struct {
		call  string
		reads string
	}{
		{"fmt.Println()", "format print-line"},
		{"os.Open()", "OS open"},
		{"io.Copy()", "I/O copy"},
		{"sync.Mutex", "sync mutex"},
		{"net.Dial()", "network dial"},
		{"bytes.Buffer", "bytes buffer"},
		{"sort.Slice()", "sort slice"},
	}

	for _, e := range examples {
		fmt.Printf("  %-20s -> reads as: %q\n", e.call, e.reads)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. Package names are short, lowercase, and usually one word")
	fmt.Println("  2. Names compose: package.Export should read naturally")
	fmt.Println("  3. No stutter: http.Client, not http.HTTPClient")
	fmt.Println("  4. No utils/helpers - name by responsibility")
	fmt.Println("  5. Name by what it provides, not how it works")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: PD.2 visibility")
	fmt.Println("   Current: PD.1 (naming)")
	fmt.Println("---------------------------------------------------")
}
