// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// ============================================================================
// Stage 06: Web Masterclass — Dependency Injection
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The "application struct" pattern — Go's standard DI approach
//   - Why globals are bad and how to avoid them
//   - Handler methods on a struct vs standalone functions
//   - Constructor functions for wiring dependencies
//
// ENGINEERING DEPTH:
//   Many frameworks in Java, Python, and TS use "Global Magic" where they scan
//   your code using heavy Reflection to inject Singletons automatically on boot.
//   This causes massive startup latency and hides control flow. Go strongly
//   rejects "Magic". Instead, we explicitly bind dependencies directly to an
//   `application` struct. Because Go statically compiles your binary, if you
//   misconfigure your dependencies, it won't even compile. You catch
//   infrastructure failures during `go build` rather than mid-flight at runtime.
//
// RUN: go run ./06-backend-db/01-web-and-database/web-masterclass/2-dependency-injection
// ============================================================================

// config holds application configuration
type config struct {
	port string
	env  string
}

// application holds all dependencies shared across handlers.
// This is the standard Go pattern for dependency injection in web apps.
// Instead of global variables, dependencies are explicit and testable.
type application struct {
	config config
	logger *slog.Logger
}

// newApplication is a constructor that wires all dependencies together.
// In production, this would also accept repositories, cache clients, etc.
func newApplication() *application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg := config{
		port: ":8080",
		env:  "development",
	}

	return &application{
		config: cfg,
		logger: logger,
	}
}

func main() {
	app := newApplication()

	mux := http.NewServeMux()

	// Handlers are methods on *application — they have access to
	// app.logger, app.config, and any other injected dependencies.
	mux.HandleFunc("GET /", app.handleHome)
	mux.HandleFunc("GET /health", app.handleHealth)

	app.logger.Info("server starting",
		slog.String("port", app.config.port),
		slog.String("env", app.config.env),
	)

	log.Fatal(http.ListenAndServe(app.config.port, mux))
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: WM.3 templates")
	fmt.Println("   Current: WM.2 (dependency injection)")
	fmt.Println("---------------------------------------------------")
}

// handleHome is a method on *application, NOT a standalone function.
// This gives it access to all injected dependencies via `app`.
func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("home page requested",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
	fmt.Fprintf(w, "Hello from %s environment!", app.config.env)
}

func (app *application) handleHealth(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("health check")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","env":"%s"}`, app.config.env)
}
