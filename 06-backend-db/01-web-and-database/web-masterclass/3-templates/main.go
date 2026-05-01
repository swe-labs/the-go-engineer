// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - HTML Templates
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use the 'html/template' package for safe HTML rendering.
//   - How to define layouts and partials for reusable UI components.
//   - The importance of template caching for production performance.
//   - How Go protects you from Cross-Site Scripting (XSS) automatically.
//
// WHY THIS MATTERS:
//   - While APIs are great, many applications still need to render
//     dynamic HTML on the server. Go's template engine is fast,
//     secure, and built directly into the standard library.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/3-templates
//
// KEY TAKEAWAY:
//   - Context-aware escaping makes Go templates one of the most secure
//     engines in the industry.
// ============================================================================

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Stage 06: Web Masterclass - HTML Templates
//
//   - html/template: Context-aware security
//   - Template Caching: Parse once, execute many
//   - Template Execution: Merging data with UI
//
// ENGINEERING DEPTH:
//   The `html/template` package is not just a string replacer.
//   It understands HTML, CSS, and JavaScript. If you try to inject
//   a malicious `<script>` tag into a user's name field, Go will
//   automatically escape it (e.g., `&lt;script&gt;`) so the browser
//   treats it as plain text. This "Context-Aware Escaping" is a
//   powerful defense against the most common web vulnerability: XSS.

// application holds the pre-parsed templates.
type application struct {
	templates *template.Template
}

func main() {
	// 1. Parse templates on startup (The "Parse Once" pattern)
	// In a real app, you would load these from files or use //go:embed.
	// For this demo, we use a string-based template.
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head><title>{{.Title}}</title></head>
		<body>
			<h1>{{.Header}}</h1>
			<ul>
				{{range .Items}}
					<li>{{.}}</li>
				{{end}}
			</ul>
			<p>Environment: {{.Env}}</p>
		</body>
		</html>
	`)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	app := &application{templates: tmpl}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handleHome)

	fmt.Println("=== Web Masterclass: HTML Templates ===")
	fmt.Println("  🚀 Server starting on http://localhost:8082")
	fmt.Println()

	// 2. Start the server
	log.Fatal(http.ListenAndServe(":8082", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.4 -> 06-backend-db/01-web-and-database/web-masterclass/4-middleware")
	fmt.Println("Current: MC.3 (templates)")
	fmt.Println("Previous: MC.2 (dependency-injection)")
	fmt.Println("---------------------------------------------------")
}

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	// 3. Define the data for the template
	data := struct {
		Title  string
		Header string
		Items  []string
		Env    string
	}{
		Title:  "Learning Templates",
		Header: "Welcome to the Template Masterclass",
		Items:  []string{"Safe Escaping", "Loops & Ranges", "Conditionals"},
		Env:    "Go 1.22+",
	}

	// 4. Execute the template with the data
	// Note: html/template automatically handles escaping 'data'!
	err := app.templates.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
