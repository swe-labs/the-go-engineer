# PII Redactor

## Learning objective

Build a high-performance privacy-preserving log redactor in Go using `log/slog`. Master custom `slog.Handler` middleware architecture, sensitive attribute sanitization, pattern-based redaction (emails, credit card numbers, SSNs), recursive group inspection, and diagnostic metric collection.

## Why this matters

Under regulations like GDPR, HIPAA, and PCI-DSS, logging raw Personally Identifiable Information (PII) or credentials to centralized log aggregators (e.g. Datadog, Elasticsearch) constitutes a major compliance violation and security liability. Building automatic redaction directly into Go's standard structured logging handler guarantees sensitive data is sanitized before leaving the application process.

## Prerequisites

Complete all Module 12 lessons (01–21) on structured logging with `log/slog`, custom handlers, metrics instrumentation, and secure data handling.

## Available concepts

This project draws on `log/slog`, `slog.Handler` interface, `slog.Record`, `slog.Attr`, recursive `slog.Group` processing, regular expressions with `regexp`, atomic metric counters with `sync/atomic`, and JSON log parsing.

## Project overview

The `_starter/` directory provides type signatures, regex definitions, and stub methods. Your job is to implement the custom `slog.Handler` middleware so that:

- Sensitive key names (`password`, `token`, `secret`, `authorization`, `api_key`) are masked with `[REDACTED_SECRET]`.
- Credit card numbers are masked with leading asterisks while preserving the last 4 digits (`****-****-****-1234`).
- Social Security Numbers are masked as `***-**-1234`.
- Email addresses are partially masked (`j***e@domain.com`).
- Nested `slog.Group` hierarchies are sanitized recursively.
- Diagnostic metrics track `LogsProcessed` and `PIIRedacted`.

## Redaction Rules

| Pattern / Key | Input Example | Redacted Output |
|---|---|---|
| Sensitive Keys | `password: "secret123"` | `password: "[REDACTED_SECRET]"` |
| Credit Cards | `4111-2222-3333-1234` | `****-****-****-1234` |
| Social Security | `123-45-6789` | `***-**-6789` |
| Email Addresses | `john.doe@example.com` | `j***e@example.com` |

## Tasks

### Task 1 — Implement key-based sanitization

In `_starter/main.go`:
- Implement `sanitizeAttr(a slog.Attr) slog.Attr`:
  - Check if `a.Key` matches sensitive words (`password`, `token`, `secret`, `authorization`, `api_key`).
  - If matched, replace with `[REDACTED_SECRET]` and increment `h.stats.PIIRedacted`.

### Task 2 — Implement pattern-based value masking

- Implement `maskString(text string) string`:
  - Replace credit card numbers with `****-****-****-<last4>`.
  - Replace SSNs with `***-**-<last4>`.
  - Mask email addresses while retaining the first character, last character of the username, and the domain.
  - Increment `h.stats.PIIRedacted` for each replacement.

### Task 3 — Handle nested groups & slog delegation

- Implement `slog.KindGroup` recursive sanitization in `sanitizeAttr`.
- Implement `Handle(ctx context.Context, r slog.Record) error`:
  - Increment `h.stats.LogsProcessed`.
  - Construct a sanitized record with masked message and sanitized attributes.
  - Forward to `h.next.Handle(ctx, sanitized)`.

## Verification

Run all automated unit tests:

```bash
go test -v ./curriculum/modules/12-observability-diagnostics/projects/pii-redactor/_starter
```

Expected output:
```
=== RUN   TestPIIRedaction_SensitiveKeys
--- PASS: TestPIIRedaction_SensitiveKeys (0.00s)
=== RUN   TestPIIRedaction_CreditCardAndSSN
--- PASS: TestPIIRedaction_CreditCardAndSSN (0.00s)
=== RUN   TestPIIRedaction_Email
--- PASS: TestPIIRedaction_Email (0.00s)
=== RUN   TestPIIRedaction_NestedGroups
--- PASS: TestPIIRedaction_NestedGroups (0.00s)
=== RUN   TestPIIRedaction_BenignLogUnchanged
--- PASS: TestPIIRedaction_BenignLogUnchanged (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Redaction Thoroughness** | Sensitive keys, credit cards, SSNs, and emails sanitized in both attributes and messages | Keys sanitized but message string unmasked | Sensitive data leaks into final output |
| **Handler Architecture** | Implements full `slog.Handler` interface (`Handle`, `WithAttrs`, `WithGroup`, `Enabled`) cleanly | Implements `Handle` but ignores `WithGroup` | Custom handler crashes on complex attribute types |
| **Non-PII Integrity** | Benign log entries, timestamps, and numbers pass through without corruption | Numbers or booleans accidentally coerced to strings | Benign attributes stripped or corrupted |
| **Observability Accounting** | Atomic counters track processed logs and redaction events with zero race conditions | Counters implemented without atomic primitives | No metric instrumentation |

## NEXT UP

[Observability Retrofit](../observability-retrofit/README.md)
