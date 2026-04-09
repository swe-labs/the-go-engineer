// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

// ============================================================================
// Section 14: Application Architecture - Structured Logging: slog Basics
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why structured logging beats fmt.Printf for production systems
//   - slog.TextHandler vs slog.JSONHandler
//   - Log levels: Debug, Info, Warn, Error
//   - Attributes: typed key-value pairs attached to log records
//   - Groups: namespacing related attributes together
//   - The default logger vs creating your own
//
// ENGINEERING DEPTH:
//   Every log statement you write is a contract with your future self at 3am.
//   `fmt.Println("user logged in: alice")` gives you exactly one string to grep.
//   `logger.Info("user logged in", slog.String("user", "alice"))` gives you a
//   queryable field in Datadog, Loki, or CloudWatch — you can filter a million
//   events down to `user=alice` in milliseconds.
//
//   Internally, slog separates the RECORD (what happened) from the HANDLER
//   (where it goes and how it is formatted). This lets you swap output formats
//   without changing a single call site. In development, use TextHandler
//   for readability. In production, use JSONHandler for machine processing.
//
// RUN: go run ./14-application-architecture/structured-logging/1-slog-basics
// ============================================================================

func main() {
	// =========================================================================
	// 1. Text Handler — human-readable output
	// =========================================================================
	// slog.NewTextHandler writes: time=... level=INFO msg="..." key=value
	// The first arg is the destination (os.Stdout, a file, bytes.Buffer).
	// The second arg is options: set minimum level, add source location, etc.
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Show ALL levels (default drops Debug)
	}))

	textLogger.Info("server started", slog.String("addr", ":8080"), slog.String("env", "development"))
	textLogger.Debug("cache miss", slog.String("key", "user:42"))
	textLogger.Warn("high memory usage", slog.Int("mb", 3814))
	textLogger.Error("database timeout", slog.Duration("elapsed", 5*time.Second))

	// =========================================================================
	// 2. JSON Handler — machine-readable output
	// =========================================================================
	// slog.NewJSONHandler writes one JSON object per line (NDJSON format).
	// This is the format your observability platform expects in production.
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	jsonLogger.Info("user authenticated",
		slog.String("user_id", "usr_01HX4K"),
		slog.String("method", "oauth2"),
		slog.Bool("mfa_used", true),
		slog.Duration("latency", 42*time.Millisecond),
	)

	// =========================================================================
	// 3. The 4 typed attribute constructors
	// =========================================================================
	// Use specific constructors (slog.String, slog.Int, etc.) rather than
	// slog.Any when you know the type. Typed constructors are ~2x faster
	// because they avoid reflect.ValueOf() calls inside slog.Any.
	jsonLogger.Info("payment processed",
		slog.String("payment_id", "pay_789"),                  // string
		slog.Int("amount_cents", 4999),                        // int (any int width)
		slog.Float64("exchange_rate", 1.0824),                 // float
		slog.Bool("requires_3ds", false),                      // bool
		slog.Time("processed_at", time.Now()),                 // time.Time
		slog.Duration("gateway_latency", 88*time.Millisecond), // duration
	)

	// =========================================================================
	// 4. Groups — namespace related attributes
	// =========================================================================
	// slog.Group("request", ...) produces: request.method=GET request.path=/api/v1
	// This prevents key collisions when logging from multiple subsystems.
	jsonLogger.Info("http request",
		slog.Group("request",
			slog.String("method", "POST"),
			slog.String("path", "/api/v1/orders"),
			slog.Int("status", 201),
		),
		slog.Group("timing",
			slog.Duration("db", 12*time.Millisecond),
			slog.Duration("total", 28*time.Millisecond),
		),
	)

	// =========================================================================
	// 5. With — pre-loading common fields (logger.With)
	// =========================================================================
	// Use With() to create a child logger that includes shared fields on every
	// subsequent log call. This is how you attach request IDs, service names,
	// and user IDs without repeating them on every line.
	requestLogger := jsonLogger.With(
		slog.String("request_id", "req_abc123"),
		slog.String("user_id", "usr_01HX4K"),
	)

	// Both of these logs automatically include request_id and user_id:
	requestLogger.Info("order created", slog.String("order_id", "ord_001"))
	requestLogger.Info("inventory reserved", slog.Int("units", 3))

	// =========================================================================
	// 6. slog.Default — the package-level logger
	// =========================================================================
	// Replace the default logger so all slog.Info() calls use your handler.
	// This is useful when libraries use slog.Default() internally.
	slog.SetDefault(jsonLogger.With(slog.String("service", "api-gateway")))
	slog.Info("default logger replaced") // Now emits JSON + service field

	// KEY TAKEAWAY:
	// - slog separates record creation (call sites) from output (handler)
	// - Use slog.TextHandler in dev, slog.JSONHandler in production
	// - Use typed constructors (slog.String, slog.Int) not slog.Any
	// - slog.Group namespaces related fields
	// - logger.With() pre-loads shared fields (request IDs, service name)
	// - slog.SetDefault() makes every slog.Info() call use your handler
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: SL.2 context-keyed logger")
	fmt.Println("   Current: SL.1 (slog basics)")
	fmt.Println("---------------------------------------------------")
}
