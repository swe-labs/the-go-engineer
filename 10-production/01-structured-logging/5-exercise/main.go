// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: PII Redactor
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Use slog.HandlerOptions.ReplaceAttr to build a logger that automatically redacts sensitive attributes before they reach the output handler.
//
// WHY THIS MATTERS:
//   - GDPR, PCI-DSS, and HIPAA require specific handling of PII.
//   - ReplaceAttr is the hook to centrally enforce redaction.
//
// KEY TAKEAWAY:
//   - One ReplaceAttr function protects all downstream handlers.
// ============================================================================

package main

import (
	"log/slog"
	"os"
)

// Stage 10: Application Architecture - Structured Logging: Exercise Solution

func main() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Check if the key matches any of our PII fields
			switch a.Key {
			case "password", "credit_card", "ssn":
				// Return a new Attr with the same key, but redacted value
				return slog.String(a.Key, "[REDACTED]")
			}
			return a // Return original attribute if no match
		},
	}

	// Create a JSON handler using our custom options
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Set it as the default logger for the standard library
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// --- DO NOT MODIFY BELOW THIS LINE ---
	slog.Info("user registration attempt",
		slog.String("username", "jdoe"),
		slog.String("password", "supersecret123"),
		slog.String("ip", "192.168.1.50"),
	)

	slog.Error("payment failed",
		slog.String("user_id", "u_999"),
		slog.String("credit_card", "4111-1111-1111-1111"),
		slog.String("error", "insufficient funds"),
	)
}
