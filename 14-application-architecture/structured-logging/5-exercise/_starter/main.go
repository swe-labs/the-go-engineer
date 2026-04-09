// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"log/slog"
	"os"
)

// ============================================================================
// Section 14: Application Architecture - Structured Logging: Exercise Starter
// Level: Intermediate
// ============================================================================
//
// THE TASK:
// You are building an enterprise service. A common requirement is ensuring
// that Personally Identifiable Information (PII) like passwords,
// credit cards, and social security numbers are NEVER logged in plain text.
//
// Your goal is to configure a custom slog.JSONHandler using
// slog.HandlerOptions and its ReplaceAttr function.
//
// REQUIREMENTS:
// 1. If any logged attribute has the key "password", "credit_card", or "ssn",
//    replace its value with the string "[REDACTED]".
// 2. The logger must output in JSON format to os.Stdout.
// 3. DO NOT hardcode the reduction in the slog.Info() calls; it must happen
//    automatically inside the logger's Handler.
//
// RUN: go run ./14-application-architecture/structured-logging/5-exercise/_starter
// ============================================================================

func main() {
	// TODO 1: Create our custom options with the ReplaceAttr function
	// HINT: ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	opts := &slog.HandlerOptions{
		// ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr { ... }
	}

	// TODO 2: Create a JSON handler using os.Stdout and your opts
	_ = opts      // remove this once implemented
	_ = os.Stdout // remove this once implemented
	// handler := ...

	// TODO 3: Create the logger and set it as the default
	// logger := slog.New(handler)
	// slog.SetDefault(logger)

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
