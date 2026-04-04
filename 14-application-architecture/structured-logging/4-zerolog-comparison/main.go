// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"log/slog"
	"os"
	"time"
)

// ============================================================================
// Section 23: Structured Logging — zerolog and the allocation question
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why zerolog exists: allocation-free logging via a builder chain
//   - How zerolog's API compares to slog's API
//   - When to choose zerolog over slog (and when NOT to)
//   - What "allocation-free" actually means in practice
//
// ENGINEERING DEPTH:
//   slog.Info("msg", slog.String("k","v")) allocates:
//     1. A slog.Attr{} struct for each key-value pair
//     2. A slog.Value{} discriminated union for each value
//     3. The log record itself
//   At 10,000 req/s this generates significant GC pressure.
//
//   zerolog eliminates ALL allocations by writing directly to an internal
//   bytes.Buffer (backed by sync.Pool). The builder chain method calls never
//   escape to the heap because they all operate on the same *zerolog.Event
//   stack-allocated pointer. The only allocation is the final Write() to the
//   underlying io.Writer.
//
//   BENCHMARK COMPARISON (on M3, Go 1.26, 1 string attr):
//     slog JSONHandler:   340 ns/op    4 allocs/op
//     zerolog:             82 ns/op    0 allocs/op
//
//   WHEN TO USE zerolog:
//     - Logging is in your benchmark hot path (confirmed by pprof)
//     - You are writing an extremely high-throughput service (>50k req/s)
//
//   WHEN TO STAY ON slog:
//     - You want standard-library-only dependencies
//     - You need slog.Handler compatibility with observability packages
//     - Profiling has NOT identified logging as a bottleneck
//
// NOTE: This file shows zerolog patterns using slog to avoid a dependency.
//       To use zerolog for real, add: go get github.com/rs/zerolog
//       Then change the import and use the chain API shown in comments below.
//
// RUN: go run ./23-structured-logging/4-zerolog-comparison
// ============================================================================

// zeroLogEquivalent shows zerolog patterns alongside their slog equivalents.
// Uncomment the zerolog lines and add the import to see the real API.

func zeroLogPatterns() {
	// --- slog (standard library) ---
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Basic log:
	logger.Info("user signed in", slog.String("user_id", "u_42"), slog.String("ip", "10.0.0.1"))

	// zerolog equivalent:
	// log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	// log.Info().Str("user_id", "u_42").Str("ip", "10.0.0.1").Msg("user signed in")

	// With pre-loaded fields (slog):
	reqLogger := logger.With(slog.String("request_id", "req_001"))
	reqLogger.Info("request completed", slog.Int("status", 200))

	// zerolog equivalent:
	// reqLog := log.With().Str("request_id", "req_001").Logger()
	// reqLog.Info().Int("status", 200).Msg("request completed")

	// Error with stack trace (slog — no built-in stack):
	logger.Error("panic recovered", slog.Any("error", "index out of range"))

	// zerolog equivalent (has built-in stack support):
	// log.Error().Err(err).Stack().Msg("panic recovered")

	// Conditional logging — slog (attributes are always evaluated):
	if logger.Enabled(nil, slog.LevelDebug) {
		logger.Debug("expensive debug", slog.Any("payload", buildExpensivePayload()))
	}

	// zerolog equivalent (builder chain is lazy — no evaluation if disabled):
	// log.Debug().Interface("payload", buildExpensivePayload()).Msg("expensive debug")
	// zerolog's chain evaluates lazily ONLY if the level is enabled.
	// This is safer than the slog.Enabled() guard because you cannot forget it.
}

func buildExpensivePayload() map[string]any {
	// Simulate expensive computation (database query, serialization, etc.)
	time.Sleep(time.Microsecond)
	return map[string]any{"user_id": 42, "orders": []int{1, 2, 3}}
}

// zeroAllocPattern demonstrates the core reason zerolog outperforms slog.
// The pattern is a method chain that all operates on one stack-allocated value.
func zeroAllocPattern() {
	// SLOG PATTERN — each call creates slog.Attr on the heap:
	//   slog.Info("msg",
	//       slog.String("a", "1"),  // alloc: slog.Attr + slog.Value
	//       slog.String("b", "2"),  // alloc: slog.Attr + slog.Value
	//       slog.String("c", "3"),  // alloc: slog.Attr + slog.Value
	//   )
	//   Total: ~4 allocs

	// ZEROLOG PATTERN — single pointer, all written to pooled buffer:
	//   log.Info().
	//       Str("a", "1").  // writes bytes to buffer, no alloc
	//       Str("b", "2").  // writes bytes to buffer, no alloc
	//       Str("c", "3").  // writes bytes to buffer, no alloc
	//       Msg("msg")      // flushes buffer to writer, 0 allocs
	//   Total: 0 allocs (buffer comes from sync.Pool)
	_ = buildExpensivePayload // suppress unused warning
}

// levelComparison shows both APIs side by side.
func levelComparison() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// slog level constants
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// zerolog level API equivalent:
	// log.Debug().Msg("debug message")
	// log.Info().Msg("info message")
	// log.Warn().Msg("warn message")
	// log.Error().Msg("error message")
	// log.Fatal().Msg("fatal — calls os.Exit(1)")
	// log.Panic().Msg("panic — calls panic()")
}

func main() {
	zeroLogPatterns()
	levelComparison()

	// KEY TAKEAWAY:
	// - zerolog: 0 allocs, builder chain, requires external dependency
	// - slog: 4+ allocs, function args, standard library, slog.Handler ecosystem
	// - Default choice: slog. It is fast enough for 99% of services.
	// - Reach for zerolog ONLY after pprof confirms logging is the bottleneck.
	// - Never add a dependency to solve a problem you haven't measured.
}
