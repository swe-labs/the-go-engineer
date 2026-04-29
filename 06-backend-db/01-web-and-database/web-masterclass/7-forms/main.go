// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Working with Forms
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to parse and process HTML forms using 'r.ParseForm'.
//   - How to perform server-side validation.
//   - Why you should never trust client-side validation alone.
//   - How to handle multi-value form inputs.
//
// WHY THIS MATTERS:
//   - Forms are the primary way users interact with your application.
//     Learning how to safely extract, clean, and validate this data is
//     essential for building secure and reliable web services.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/7-forms
//
// KEY TAKEAWAY:
//   - Trust but verify. Always validate everything on the server.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Stage 06: Web Masterclass - Working with Forms
//
//   - r.ParseForm(): Reading the request body
//   - r.PostFormValue: Helper for single values
//   - Server-Side Validation: The Golden Rule
//
// ENGINEERING DEPTH:
//   In Go, calling `r.ParseForm()` is mandatory before you can
//   access form values. This function reads the request body (which
//   is a stream) and parses it into a map. Note that this consumes
//   the body stream-you can only parse it once! Go's form parsing
//   also supports multi-value fields (like checkboxes), which are
//   stored as a `[]string` in the `r.PostForm` map.

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /signup", handleSignup)

	fmt.Println("=== Web Masterclass: Working with Forms ===")
	fmt.Println("  🚀 Server starting on http://localhost:8086")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":8086", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.8 posts-crud")
	fmt.Println("Current: MC.7 (forms)")
	fmt.Println("Previous: MC.6 (authentication)")
	fmt.Println("---------------------------------------------------")
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<h2>Sign Up</h2>
		<form action="/signup" method="POST">
			<label>Email: <input type="email" name="email"></label><br>
			<label>Password: <input type="password" name="password"></label><br>
			<button type="submit">Register</button>
		</form>
	`)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// 2. Extract values
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// 3. Server-side validation (NEVER TRUST THE CLIENT)
	errors := make(map[string]string)

	if strings.TrimSpace(email) == "" {
		errors["email"] = "Email is required."
	}
	if len(password) < 8 {
		errors["password"] = "Password must be at least 8 characters."
	}

	// 4. Handle validation failure
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, "<h3>Validation Errors:</h3><ul>")
		for field, msg := range errors {
			fmt.Fprintf(w, "<li>%s: %s</li>", field, msg)
		}
		fmt.Fprintln(w, "</ul><a href='/'>Go Back</a>")
		return
	}

	// 5. Success!
	fmt.Fprintf(w, "Registration successful for %s!", email)
}
