// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Dependency Injection
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use the "Application Struct" pattern to share dependencies.
//   - Why global variables are an anti-pattern in Go web services.
//   - How to use methods as HTTP handlers.
//
// WHY THIS MATTERS:
//   - As your application grows, your handlers will need access to database
//     pools, loggers, and configuration settings. Dependency Injection (DI)
//     ensures your code remains clean, modular, and easy to unit test.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/2-dependency-injection
//
// KEY TAKEAWAY:
//   - Explicit is better than implicit. No magic, just clean wiring.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Stage 06: Web Masterclass - Dependency Injection
//
//   - The App Struct: Container for shared resources
//   - Method Handlers: Giving functions context
//   - Logger Injection: Centralized observability
//
// ENGINEERING DEPTH:
//   In Go, we avoid "Magic" DI frameworks that use reflection.
//   Instead, we manually "wire" our application by creating a
//   struct (often called `application` or `server`) that holds
//   pointers to our dependencies. Because our handlers are methods
//   on this struct, they have direct, type-safe access to everything
//   they need without resorting to global variables.

// application holds all shared dependencies.
type application struct {
	logger *slog.Logger
	env    string
}

func main() {
	// 1. Initialize dependencies
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
		env:    "development",
	}

	mux := http.NewServeMux()

	// 2. Register handlers as methods
	// This is the "Magic Sauce". Because handleHome is a method on 'app',
	// it can access app.logger and app.env!
	mux.HandleFunc("GET /", app.handleHome)
	mux.HandleFunc("GET /health", app.handleHealth)

	fmt.Println("=== Web Masterclass: Dependency Injection ===")
	fmt.Println("  🚀 Server starting on http://localhost:8081")
	fmt.Println()

	// 3. Start the server
	log.Fatal(http.ListenAndServe(":8081", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.3 -> 06-backend-db/01-web-and-database/web-masterclass/3-templates")
	fmt.Println("Current: MC.2 (dependency-injection)")
	fmt.Println("Previous: MC.1 (routing)")
	fmt.Println("---------------------------------------------------")
}

// handleHome is a method on the *application struct.
func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("home page requested", "path", r.URL.Path)
	fmt.Fprintf(w, "Hello! You are running in the %s environment.", app.env)
}

func (app *application) handleHealth(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("health check pinged")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}
