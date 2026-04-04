// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ============================================================================
// Section 13: Form Processing
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Custom form validation library
//   - Server-side validation (NEVER trust client-side validation alone)
//   - Error collection and display
//   - Reusable validation methods: Required, MinLength, MaxLength, MatchesPattern
//
// ENGINEERING DEPTH:
//   A massive pitfall for juniors is trusting frontend `<input required>` and
//   regex attributes. The golden rule of Web Security is: "The Client is in the
//   hands of the Enemy." Anyone can open Chrome DevTools, delete the `required`
//   tag, and submit a malformed HTTP `application/x-www-form-urlencoded` request.
//   Your server MUST strictly validate the payload. We encapsulate `url.Values`
//   in a custom `Form` struct to attach an `Errors map[string][]string` directly
//   to the execution lifecycle, cleanly isolating business validation from transport.
//
// RUN: go run ./13-web-masterclass/7-forms
// ============================================================================

// Form wraps url.Values and adds validation capabilities.
// This is a common pattern in Go web apps — a thin validation layer
// over the standard library's form parsing.
type Form struct {
	url.Values
	Errors map[string][]string
}

// NewForm creates a Form from parsed form data.
func NewForm(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(map[string][]string),
	}
}

// Required checks that the specified fields are not blank.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors[field] = append(f.Errors[field], "This field is required")
		}
	}
}

// MinLength checks that a field has at least `d` characters.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if len(value) < d {
		f.Errors[field] = append(f.Errors[field],
			fmt.Sprintf("This field must be at least %d characters", d))
	}
}

// MaxLength checks that a field has at most `d` characters.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if len(value) > d {
		f.Errors[field] = append(f.Errors[field],
			fmt.Sprintf("This field must be at most %d characters", d))
	}
}

// MatchesField checks that two fields have the same value (e.g., password confirmation).
func (f *Form) MatchesField(field, otherField string) {
	if f.Get(field) != f.Get(otherField) {
		f.Errors[field] = append(f.Errors[field],
			fmt.Sprintf("This field must match %s", otherField))
	}
}

// Valid returns true if there are no validation errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Show a simple HTML form
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<h1>Registration Form</h1>
			<form method="POST" action="/register">
				<p>Email: <input name="email" type="email"></p>
				<p>Username: <input name="username"></p>
				<p>Password: <input name="password" type="password"></p>
				<p>Confirm: <input name="password_confirm" type="password"></p>
				<button type="submit">Register</button>
			</form>
		`)
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		// Parse form data from the request body
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Create a Form and run validations
		form := NewForm(r.PostForm)
		form.Required("email", "username", "password", "password_confirm")
		form.MinLength("username", 3)
		form.MinLength("password", 8)
		form.MaxLength("username", 50)
		form.MatchesField("password_confirm", "password")

		// Check if validation passed
		if !form.Valid() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Validation errors:\n")
			for field, errs := range form.Errors {
				for _, e := range errs {
					fmt.Fprintf(w, "  %s: %s\n", field, e)
				}
			}
			return
		}

		fmt.Fprintf(w, "Registration successful for: %s\n", form.Get("username"))
	})

	fmt.Println("Form demo server starting on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
