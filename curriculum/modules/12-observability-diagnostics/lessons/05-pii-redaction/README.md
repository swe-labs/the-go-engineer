# PII redaction

## Learning objective

Implement a custom `slog.Handler` that redacts personally identifiable information (PII) from log output by intercepting and masking attribute values before serialization.

## Why this matters

Log files are copied to log aggregators (Datadog, Grafana Loki, CloudWatch), indexed for search, and often retained for months or years. If a log line contains an email address, password, credit card number, or SSN, that PII is now stored in a system that may not have the same access controls as the application database. GDPR (EU) fines can reach 4% of global revenue for PII exposure. CCPA (California) allows private lawsuits for data breaches. HIPAA (healthcare) violations carry criminal penalties. PII redaction is not optional — it is a legal requirement.

## Mental model

PII redaction is like a security guard at the logging factory gate. Every log record (the item on the assembly line) passes through the guard before being packaged for shipment. The guard checks each field against a list of "sensitive" labels (email, password, ssn, credit_card). If a field's label matches, the guard replaces the value with a masked version before it leaves the factory. The guard never lets a raw sensitive value reach the shipping dock.

In Go terms, the security guard is a custom `slog.Handler` that wraps the real handler (text or JSON). The guard's `Handle` method inspects every attribute, redacts matching keys, and passes the cleaned record to the inner handler.

## Core idea

Redact at the handler level, not at the call site. A handler-level redaction intercepts _every_ log call, including calls from third-party libraries and future code. Call-site redaction (manually masking values before passing to `slog.Info`) is fragile because a developer can forget to redact a new log line.

The handler checks each attribute's key against a set of sensitive keys. If the key matches, the attribute value is replaced with a redacted string (`[REDACTED]`). Non-sensitive attributes pass through unchanged. The redacted handler wraps the real output handler, so the output format (text, JSON, custom) is preserved.

## Under the hood

A custom `slog.Handler` must implement three methods:

- `Enabled(ctx, level) bool` — delegate to the inner handler.
- `Handle(ctx, record) error` — iterate the record's attrs, redact sensitive ones, create a new record with cleaned attrs, pass to inner handler.
- `WithAttrs(attrs []Attr) Handler` — return a new redaction handler wrapping the inner handler's WithAttrs result.
- `WithGroup(name string) Handler` — return a new redaction handler wrapping the inner handler's WithGroup result.

The `Handle` method rebuilds the record because `slog.Record` does not expose a way to modify attrs in place. The pattern is:

```go
func (h *RedactHandler) Handle(ctx context.Context, rec slog.Record) error {
    var attrs []slog.Attr
    rec.Attrs(func(a slog.Attr) bool {
        if h.redact[strings.ToLower(a.Key)] {
            a.Value = slog.StringValue("[REDACTED]")
        }
        attrs = append(attrs, a)
        return true
    })
    newRec := slog.NewRecord(rec.Time, rec.Level, rec.Message, rec.PC)
    for _, a := range attrs {
        newRec.AddAttrs(a)
    }
    return h.inner.Handle(ctx, newRec)
}
```

This rebuilds the record with the same time, level, message, PC, and all attrs (with sensitive ones redacted). The inner handler serializes the cleaned record as usual.

## How Go uses it

PII redaction is applied at the handler level so it covers all log output:

```go
func main() {
    baseHandler := slog.NewTextHandler(os.Stdout, nil)
    redactHandler := NewRedactHandler(baseHandler, "email", "password", "ssn", "credit_card")
    logger := slog.New(redactHandler)

    // These values are redacted before reaching the output
    logger.Info("user login",
        "user_id", "u-42",
        "email", "alice@example.com",  // → [REDACTED]
        "password", "s3cret!",          // → [REDACTED]
    )
}
```

The output shows `email=[REDACTED] password=[REDACTED]` while `user_id=u-42` passes through unchanged. The redaction is transparent to the calling code — developers do not need to remember to redact at each call site.

## Go example

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type RedactHandler struct {
	inner  slog.Handler
	redact map[string]bool
}

func NewRedactHandler(inner slog.Handler, sensitiveKeys ...string) *RedactHandler {
	r := make(map[string]bool, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		r[strings.ToLower(k)] = true
	}
	return &RedactHandler{inner: inner, redact: r}
}

func (h *RedactHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *RedactHandler) Handle(ctx context.Context, rec slog.Record) error {
	var attrs []slog.Attr
	rec.Attrs(func(a slog.Attr) bool {
		if h.redact[strings.ToLower(a.Key)] {
			a.Value = slog.StringValue("[REDACTED]")
		}
		attrs = append(attrs, a)
		return true
	})
	newRec := slog.NewRecord(rec.Time, rec.Level, rec.Message, rec.PC)
	for _, a := range attrs {
		newRec.AddAttrs(a)
	}
	return h.inner.Handle(ctx, newRec)
}

func (h *RedactHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &RedactHandler{inner: h.inner.WithAttrs(attrs), redact: h.redact}
}

func (h *RedactHandler) WithGroup(name string) slog.Handler {
	return &RedactHandler{inner: h.inner.WithGroup(name), redact: h.redact}
}

func main() {
	baseHandler := slog.NewTextHandler(os.Stdout, nil)
	handler := NewRedactHandler(baseHandler, "email", "password", "ssn", "credit_card")
	logger := slog.New(handler)

	logger.Info("user operation",
		"user_id", "u-42",
		"email", "alice@example.com",
		"name", "Alice Smith",
		"role", "admin",
	)
}
```

Run with `go run .` to see `email=[REDACTED]` while `user_id=u-42` and `role=admin` pass through unchanged.

## Step-by-step execution

For `logger.Info("user operation", "email", "alice@example.com", "role", "admin")` with the redact handler:

1. `logger.Info` creates a `slog.Record` with the message and attrs: `[email=alice@example.com, role=admin]`.
2. The redact handler's `Handle` is called with the record.
3. `rec.Attrs` iterates the attrs. For `email` (sensitive): the value is replaced with `[REDACTED]`. For `role` (not sensitive): the value passes through unchanged.
4. A new record is created with the redacted attrs: `[email=[REDACTED], role=admin]`.
5. The new record is passed to `h.inner.Handle(ctx, newRec)`.
6. The inner (text) handler serializes: `time=... level=INFO msg="user operation" email=[REDACTED] role=admin`.
7. The output is written to stderr.

## Common mistakes

- **Redacting at the presentation layer only**: Applying redaction only in the HTTP response (masking email in JSON body) but not in logs. The full email appears in the log aggregator where operators and third-party tools can see it. Fix: redact at every output boundary — response, log, error message, metrics tag.

- **Using regex redaction on serialized JSON**: Applying a regex like `s/(\w+)@/***@/` on the serialized log string can match keys instead of values (the JSON key `"email"` also contains `@`). Fix: apply redaction on the structured attribute values before serialization, using the key-matching approach above.

- **Assuming redaction once means redacted forever**: A code change that adds a new log line with a new user field bypasses existing redaction logic. Fix: enforce redaction at the handler level (custom `slog.Handler`), not at individual log call sites. Handler-level redaction covers all current and future log calls.

- **Case-insensitive matching only on keys, not values**: Some developers try to match PII patterns in values (e.g., `\b[\w.]+@\w+\.\w+\b`). This is expensive and error-prone. Fix: match on keys, which are controlled by the developer, rather than values, which are user data.

- **Not redacting in error strings**: Error messages that include user input (`fmt.Errorf("user %s not found: invalid email %s", userID, email)`) leak PII if the error is logged. Fix: wrap errors without user data in the message, or redact error string attributes at the handler level.

## Debugging walkthrough

Consider a user registration handler that logs the full request body:

```go
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    slog.Info("registering user",
        "user", u, // logs all fields including password!
    )
}
```

**Symptom**: The log aggregator contains plaintext passwords for all registered users. GDPR violation.

**Investigation**: An audit of the log aggregator reveals `password=plaintext` in the `user` attribute for every registration log line.

**Root cause**: The `User` struct has a `Password` field, and `slog.Any("user", u)` serializes all exported fields. No redaction was applied.

**Fix**: Add a redact handler that redacts the `"password"` key at the handler level:

```go
baseHandler := slog.NewJSONHandler(logFile, nil)
redactHandler := NewRedactHandler(baseHandler, "password", "email", "ssn")
logger := slog.New(redactHandler)
slog.SetDefault(logger)
```

Alternatively, add a `LogValue` method to the `User` struct that returns a `slog.Value` with the password field redacted, but this requires modifying every sensitive type. The handler-level approach requires no type changes and covers all code paths.

Post-fix, the log aggregator shows `password=[REDACTED]` for all registration log lines.

## Production notes

PII redaction is legally required for GDPR (EU), CCPA (California), SOC 2, and HIPAA (healthcare) compliance. Every production Go service that logs user data needs PII redaction.

Production redaction checklist:

- **Handler-level redaction**: Wrap the output handler with a redact handler that covers all log output.
- **Key-based matching**: Match on attribute keys, not values. Keys are developer-controlled; values are user-controlled.
- **Case-insensitive keys**: `Email`, `EMAIL`, `email` should all match.
- **Regular audit**: Periodically scan log samples for unredacted PII.
- **Redact error messages too**: Error strings that contain user input are a common PII leak path.
- **Use `slog.Record` rebuilding**: The clean record preserves time, level, message, and non-sensitive attrs.

Stripe's logging pipeline redacts credit card numbers before they reach the centralized log store. GitHub's audit log redacts email addresses. In healthcare systems, every PHI (Protected Health Information) field is redacted at the handler level.

## Performance implications

The redact handler adds overhead proportional to the number of attributes per log line:

- **Attribute iteration**: `rec.Attrs` iterates all attributes. For typical log lines with 5–10 attrs, this is a loop of 5–10 iterations.
- **Key matching**: Each attribute key is lowercased and checked against a map lookup — O(1) per attribute.
- **Record rebuilding**: `slog.NewRecord` creates a new record, and `AddAttrs` copies the attribute slice. For a record with 10 attrs, this allocates ~320 bytes.
- **Total cost**: ~200–500ns per log line, dominated by the new record allocation.

For the critical path of a handler (1–2 log lines), the cost is ~1μs total. If the handler processes in 50ms, this is 0.002% overhead — negligible.

If allocation is a concern (extremely hot path), the redact handler can be optimized to reuse a pre-allocated attribute slice with sync.Pool, but this is rarely necessary.

## Practice task

Implement the `RedactHandler` type with the four handler methods (`Enabled`, `Handle`, `WithAttrs`, `WithGroup`), then use it to create a logger that redacts `email`, `password`, and `ssn` fields. The tests verify:

1. Sensitive keys are replaced with `[REDACTED]`.
2. Non-sensitive keys pass through unchanged.
3. Multiple sensitive keys can be redacted simultaneously.
4. When no sensitive keys are present in a log line, no redaction occurs.

## Tests / verification

```bash
go run ./curriculum/modules/12-observability-diagnostics/lessons/05-pii-redaction
go test ./curriculum/modules/12-observability-diagnostics/lessons/05-pii-redaction
```

## Review questions

1. Why is handler-level redaction safer than call-site redaction?
2. What is the risk of applying regex-based redaction to serialized JSON log output instead of structured attributes?
3. Name three types of PII that should be redacted in a typical user-facing Go service.
4. How would you handle case-insensitive key matching in a redaction handler?
5. A developer writes `slog.Error("failed", "err", fmt.Errorf("user %s not found", email))`. Is the email redacted by a handler-level redaction? Why or why not?

## NEXT UP

Metrics — numeric aggregations (counters, gauges, histograms) that answer "how many times did it happen?" for request rate, error rate, and latency.
