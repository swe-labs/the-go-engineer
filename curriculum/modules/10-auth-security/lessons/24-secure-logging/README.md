# Secure logging

## Learning objective

Implement structured logging in Go that never exposes secrets or PII by using redaction, appropriate log levels, and audit trails for security-relevant events.

## Why this matters

Logs are the primary source of truth for debugging, monitoring, and incident response. But logs are also a common source of data breaches. The 2021 Twilio breach started with logs containing customer credentials. The 2019 Capital One breach involved logs that exposed sensitive configuration data. Go engineers must design logging systems that provide maximum operational insight with minimum data exposure. Secure logging is a regulatory requirement (GDPR, PCI-DSS, HIPAA) and a professional obligation.

## Mental model

Think of logs as a public bulletin board. Everyone on the operations team can read them. You would not pin your password, credit card number, or private conversation to a public board. Logs are the same: they should contain enough context to diagnose problems but never sensitive data that could cause harm if exposed.

A redaction filter is like a privacy stamp: it automatically covers credit card numbers, passwords, and API keys before the log entry is posted. The original data is visible only to the application, never to the log aggregator.

## Core idea

Secure logging principles:

1. **Never log secrets**: passwords, tokens, API keys, credit cards, SSNs, PINs.
2. **Redact PII**: email addresses, phone numbers, names (only when not needed for debugging).
3. **Use structured logging**: JSON-formatted logs are machine-parseable and filterable.
4. **Appropriate log levels**: DEBUG for development, INFO for normal operations, WARN for concerning events, ERROR for failures.
5. **Audit trails**: Security-relevant events (login, logout, permission changes) should be logged to a separate, immutable audit log.

| Log level | When to use | Example |
|---|---|---|
| DEBUG | Development only | "SQL query: SELECT * FROM users" |
| INFO | Normal operations | "User u123 logged in from 192.168.1.1" |
| WARN | Concerning but not failing | "Failed login attempt for u999" |
| ERROR | Failure that requires investigation | "Database connection pool exhausted" |
| FATAL | Unrecoverable error, process exiting | "Cannot open listening port 443" |

## Under the hood

Structured loggers serialize log entries as JSON objects. A typical entry:

```json
{"time":"2026-06-24T10:30:00Z","level":"INFO","msg":"payment processed","order_id":"ord-456","amount":49.99}
```

This format is machine-parseable by log aggregation systems (ELK, Datadog, Loki, Splunk). Each field can be indexed and searched.

Redaction is implemented as a middleware or wrapper that scans log fields for known sensitive key names (case-insensitive: "password", "secret", "token", "api_key", "credit_card", "ssn") and replaces their values with `[REDACTED]`.

Audit logs are typically written to a separate output (file, database table, or stream) that is append-only and immutable. Access to audit logs is restricted and monitored.

## How Go uses it

Go's standard `log` package writes unstructured text. For structured logging, popular packages include:

- `log/slog` (Go 1.21+): Standard library structured logger.
- `github.com/sirupsen/logrus`: Popular structured logger with hooks.
- `go.uber.org/zap`: High-performance structured logger.
- `github.com/rs/zerolog`: Zero-allocation JSON logger.

Redaction can be implemented as a custom handler wrapping `slog.Handler`:

```go
type RedactionHandler struct {
    handler slog.Handler
    redactedKeys map[string]bool
}

func (h *RedactionHandler) Handle(ctx context.Context, r slog.Record) error {
    r.Attrs(func(a slog.Attr) {
        if h.redactedKeys[strings.ToLower(a.Key)] {
            a.Value = slog.StringValue("[REDACTED]")
        }
    })
    return h.handler.Handle(ctx, r)
}
```

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SecureLogger struct {
	MinLevel     LogLevel
	RedactedKeys map[string]bool
	AuditEnabled bool
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
}

func NewSecureLogger() *SecureLogger {
	return &SecureLogger{
		RedactedKeys: map[string]bool{
			"password": true, "secret": true, "token": true,
			"api_key": true, "credit_card": true, "ssn": true,
		},
	}
}

func (sl *SecureLogger) redact(fields map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		if sl.RedactedKeys[strings.ToLower(k)] {
			result[k] = "[REDACTED]"
		} else {
			result[k] = v
		}
	}
	return result
}

func (sl *SecureLogger) Info(msg string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "INFO",
		Message:   msg,
		Fields:    sl.redact(fields),
	}
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}

func main() {
	logger := NewSecureLogger()
	logger.Info("Payment processed", map[string]interface{}{
		"order_id":    "ord-456",
		"amount":      49.99,
		"credit_card": "4111-1111-1111-1111",
	})
}
```

## Step-by-step execution

When a payment is processed and the application logs the event:

1. Developer writes: `logger.Info("payment processed", fields)` where fields includes `credit_card`.
2. The `logger.Info` method calls `redact(fields)`.
3. Redaction scans each field key: `"order_id"` is not in the redacted keys set, so it passes through. `"credit_card"` is in the redacted keys set (case-insensitive match), so its value is replaced with `"[REDACTED]"`.
4. The `LogEntry` struct is serialized to JSON.
5. The JSON line is written to stdout (or file, or network).
6. The log aggregator ingests: `{"fields":{"credit_card":"[REDACTED]","order_id":"ord-456"}}`.
7. No sensitive data has left the application process.

If redaction were not in place:

1. `{"fields":{"credit_card":"4111-1111-1111-1111","order_id":"ord-456"}}` is written to the log.
2. An attacker who gains access to the log aggregation system can steal the credit card number.
3. Even without a malicious actor, a developer debugging with `tail -f logs/app.log` sees the full credit card number.

## Common mistakes

- Mistake: Logging full HTTP request bodies or environment variables without redaction.
  - Why it happens: Developers log `r.Body` or `os.Environ()` for debugging convenience.
  - Fix: Never log raw bodies or environment variables. Log specific fields that are known to be safe.

- Mistake: Logging error messages that contain sensitive input.
  - Why it happens: `fmt.Sprintf("failed for user %s", userInput)` or `log.Printf("query failed: %s", query)`.
  - Fix: Log only the error type and ID, not the user input or full query. Use error codes that can be looked up in documentation.

- Mistake: Using DEBUG level in production.
  - Why it happens: Developers set `log.SetLevel(log.DebugLevel)` to diagnose a production issue and forget to remove it.
  - Fix: Log level should be configurable via environment variable. Never commit a production config with DEBUG level.

- Mistake: Logging stack traces that contain function arguments with sensitive data.
  - Why it happens: Developers use `debug.Stack()` or `runtime.Caller()` to capture stack traces.
  - Fix: Strip sensitive arguments from stack traces, or log only the function name and file location.

- Mistake: Writing audit logs to the same stream as application logs.
  - Why it happens: It is easier to use one logger for everything.
  - Fix: Audit logs should be separate, append-only, and immutable. Application logs can be rotated and deleted; audit logs must be retained.

## Debugging walkthrough

Consider this insecure logging code:

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)
	log.Printf("Login attempt: email=%s password=%s", creds.Email, creds.Password)
	// ... authenticate
}
```

Symptom: The application logs contain plaintext passwords. A developer sees them in `kubectl logs`.

Investigation: Check the log output:

```
2026/06/24 10:30:00 Login attempt: email=alice@example.com password=hunter2
```

Root cause: The developer used `log.Printf` with a format string that includes the password field directly.

Fix: Never log passwords. Remove the password from the log message and add redaction:

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)
	log.Printf("Login attempt: email=%s", creds.Email)
	// Consider logging IP, user agent, and timestamp only
}
```

Production-grade fix using structured logging:

```go
logger.Info("login attempt", map[string]interface{}{
	"email": creds.Email,
	"ip":    r.RemoteAddr,
	"agent": r.UserAgent(),
})
```

## Production notes

- Use `log/slog` (Go 1.21+) or a mature structured logging library (zap, zerolog, logrus) with a custom redaction handler.
- Configure log level via environment variable: `LOG_LEVEL=INFO`. Never default to DEBUG in production.
- Implement redaction at the logger level, not at the call site. This ensures all logs are redacted consistently.
- For audit trails, use a separate logger that writes to an immutable store (e.g., a database table with append-only permissions, or a cloud audit log service).
- Monitor for missing redaction: use a log scanner that alerts when sensitive patterns (credit card regex, etc.) appear in log output.
- In Go, use `slog.Record`'s attribute API to implement redaction as a `slog.Handler` wrapper.
- Retain logs according to compliance requirements (typically 30-90 days for app logs, 1-7 years for audit logs).

## Performance implications

- Structured logging with JSON serialization adds overhead compared to unstructured text logging. For high-throughput services (10,000+ logs/sec), use zero-allocation loggers like `zerolog` or `zap`.
- Redaction adds O(n) overhead per log entry (scanning n fields). This is negligible (nanoseconds to microseconds).
- Audit log writes should be asynchronous (buffered) to avoid blocking the request handler.
- The main performance cost of logging is I/O (writing to disk or network). Use buffered writers and consider sampling in high-throughput environments.

## Practice task

Write a function `NewAuditLogger(auditWriter io.Writer) *SecureLogger` that creates a logger with redaction and writes all WARN and ERROR entries to an audit trail.

The audit trail should include:
- Timestamp
- Event type (from message)
- User ID (from fields)
- Action (from fields)
- Result (from fields: "success" or "failure")

Then write a `main()` that:
1. Creates an audit logger with redaction.
2. Logs several events: login success, login failure (with password redacted), payment processed (with credit card redacted).
3. Prints the audit trail entries to a buffer.
4. Verifies that no sensitive data appears in the output.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/24-secure-logging
go test ./curriculum/modules/10-auth-security/lessons/24-secure-logging
```

The existing tests verify redaction of sensitive keys (case-insensitive), log level filtering, audit trail capture, and audit trail summary generation.

## Review questions

1. What is the difference between structured logging and unstructured logging, and why is structured logging preferred for security?
2. Why should redaction happen at the logger level rather than at each call site?
3. What fields should be redacted by default in a web application's logs?
4. What is the purpose of an audit trail, and how does it differ from application logs?
5. How would you implement a redaction handler for Go's `log/slog` package?

## NEXT UP

Refresh tokens -- implementing secure refresh token rotation with short-lived access tokens, long-lived refresh tokens, secure storage, and revocation in Go APIs.
