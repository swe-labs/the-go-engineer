// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// ============================================================================
// Section 13: Web Masterclass — HTML Templates
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - html/template for safe HTML rendering (auto-escapes XSS)
//   - Template caching — parse once, execute many times
//   - Template composition — base layout + page content + partials
//   - Passing dynamic data to templates via structs
//
// ENGINEERING DEPTH:
//   The `html/template` package does not just do dumb string replacement `{{.Data}}`.
//   It fully parses your HTML into an Abstract Syntax Tree (AST) to understand
//   the DOM context. If you inject data inside a `<script>` tag, it automatically
//   JSON-encodes it. If you inject it inside an `href`, it URL-encodes it. This
//   makes Go templates completely immune to Cross-Site Scripting (XSS) by default.
//   However, parsing this AST is incredibly CPU-expensive. This is why you MUST
//   parse them ONCE on server startup and store them in a synchronized Cache map,
//   rather than parsing the files from disk on every single HTTP request.
//
// RUN: go run ./13-web-masterclass/3-templates
// ============================================================================

// templateCache stores pre-parsed templates keyed by page name.
// Parsing on every request is expensive — we parse ONCE at startup.
type templateCache map[string]*template.Template

// newTemplateCache parses all template files and returns a cache.
// Each page template is composed of: base layout + partials + page content.
func newTemplateCache(dir string) (templateCache, error) {
	cache := templateCache{}

	// Find all page templates
	pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the page filename (e.g., "home.html")
		name := filepath.Base(page)

		// Parse template files in order: base → partials → page
		// The first template parsed defines the root template name.
		ts, err := template.New(name).ParseFiles(
			filepath.Join(dir, "base.html"), // Layout wrapper
		)
		if err != nil {
			return nil, err
		}

		// Add partials (nav, footer, etc.)
		ts, err = ts.ParseGlob(filepath.Join(dir, "partials", "*.html"))
		if err != nil {
			return nil, err
		}

		// Add the specific page content
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// templateData holds all data passed to a template.
// Every handler creates a templateData and passes it to render().
type templateData struct {
	Title   string
	Content string
	Year    int
	Items   []string
}

type application struct {
	templates templateCache
	mu        sync.RWMutex // protects template cache in dev mode
}

// render executes a named template with the given data.
// It writes the result to the http.ResponseWriter.
func (app *application) render(w http.ResponseWriter, name string, data templateData) {
	app.mu.RLock()
	ts, ok := app.templates[name]
	app.mu.RUnlock()

	if !ok {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}

	// Execute the "base" template (defined in base.html).
	// This renders the full page: layout + nav + page content.
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	cache, err := newTemplateCache("./13-web-masterclass/3-templates/templates")
	if err != nil {
		log.Fatalf("Failed to create template cache: %v", err)
	}

	app := &application{templates: cache}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		app.render(w, "home.html", templateData{
			Title:   "Home",
			Content: "Welcome to the Go Web Masterclass!",
			Year:    2025,
			Items:   []string{"Routing", "Templates", "Middleware", "Auth"},
		})
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		app.render(w, "about.html", templateData{
			Title:   "About",
			Content: "Learning Go web development step by step.",
			Year:    2025,
		})
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
