# Project Assessment — PII Redactor

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-pii-redactor`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Assessment Overview & Industrial Context

Under compliance mandates including GDPR, HIPAA, and PCI-DSS, logging unmasked Personally Identifiable Information (PII), authentication credentials, or financial account identifiers to centralized log collectors (e.g., Datadog, Splunk, Elastic) constitutes a major security failure and legal liability. 

Standard application loggers often inadvertently capture passwords, API tokens, credit card numbers, social security numbers, and email addresses in debug statements or error stack traces.

In this assessment, you will implement a high-throughput, privacy-preserving log sanitizer using Go's `log/slog` custom handler architecture. Your handler will intercept log records, recursively inspect attribute trees and message strings, mask sensitive data according to industry-standard redaction rules, track telemetry metrics via atomic counters, and forward sanitized records to downstream handlers without degrading logging throughput.

---

## 2. Core Competencies Tested

- **Custom `slog.Handler` Implementation**: Implementing `Enabled`, `Handle`, `WithAttrs`, and `WithGroup` conforming to Go 1.21's structured logging specification.
- **Recursive Attribute Sanitization**: Walking arbitrary trees of `slog.Attr`, including `slog.KindGroup` nested groups, without dropping attributes or mutating shared memory.
- **Regex & Pattern-Based Masking**: Applying precompiled regular expressions to scrub credit card numbers, SSNs, and email addresses while preserving auditability (e.g. keeping the last 4 digits of credit cards).
- **Concurrency & Atomic Telemetry**: Recording processed logs and redaction event metrics using lock-free atomic primitives (`sync/atomic`).
- **Zero-Allocation Log Fast Paths**: Bypassing regex evaluation for benign numeric and boolean attributes.

---

## 3. Zero-Magic Architectural Breakdown

```
[ Application Code ]
         │ slog.InfoContext(ctx, "charge complete", "user", "alice", "cc", "4111-2222-3333-1234")
         ▼
┌─────────────────────────────────────────────────────────────────┐
│ PIIRedactingHandler (implements slog.Handler)                   │
│   1. atomic.AddUint64(&stats.LogsProcessed, 1)                  │
│   2. Mask record message: maskString(r.Message)                 │
│   3. Clone & sanitize attributes:                               │
│      - If key in {password, token, secret, api_key}:            │
│        value -> "[REDACTED_SECRET]"                             │
│      - If kind == KindGroup:                                    │
│        recursively sanitize child attributes                    │
│      - If kind == KindString:                                   │
│        apply CC, SSN, and Email regex masking                   │
│   4. Forward sanitized record to downstream next.Handle()       │
└─────────────────────────────────┬───────────────────────────────┘
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│ Downstream Handler (e.g., slog.JSONHandler -> os.Stdout)        │
└─────────────────────────────────────────────────────────────────┘
```

### The Inherent Risk of Log Sanitization
Sanitizing logs at the format handler level rather than requiring developers to manually sanitize strings ensures defense-in-depth: even if a developer writes `slog.Error("auth failed", "password", rawPassword)`, the redaction handler guarantees the secret is replaced with `[REDACTED_SECRET]` before hitting disk, stdout, or the network.

### Why Atomic Primitives Matter in Handlers
Logging handlers are invoked concurrently from hundreds of active goroutines. Using a `sync.Mutex` inside `Handle()` to track counters introduces severe lock contention. Using `atomic.AddUint64` ensures lock-free telemetry accounting with nanosecond execution overhead.

---

## 4. Redaction Specification & Rules

| Category | Input Example | Expected Redacted Output | Rule Description |
|---|---|---|---|
| **Sensitive Keys** | `password: "secret123"` | `password: "[REDACTED_SECRET]"` | Match keys: `password`, `token`, `secret`, `authorization`, `api_key` (case-insensitive) |
| **Credit Card** | `4111-2222-3333-1234` | `****-****-****-1234` | Mask first 12 digits, retain last 4 digits and hyphens |
| **Social Security** | `123-45-6789` | `***-**-6789` | Mask first 5 digits, retain last 4 digits and hyphens |
| **Email Address** | `john.doe@example.com` | `j***e@example.com` | Retain first and last characters of local-part + `@domain` |

---

## 5. Implementation Requirements & Task Breakdown

### Task 1: Key-Based Attribute Sanitization
- In `sanitizeAttr(a slog.Attr) slog.Attr`:
  - Check lowercased `a.Key`.
  - If the key is in the sensitive key set, replace `a.Value` with `slog.StringValue("[REDACTED_SECRET]")`.
  - Increment the atomic `PIIRedacted` counter.

### Task 2: Value-Based Pattern Masking
- Implement regex-based masking in `maskString(text string) string`:
  - Credit Card pattern: `\b(?:\d{4}-){3}\d{4}\b` -> `****-****-****-$1`.
  - SSN pattern: `\b\d{3}-\d{2}-(\d{4})\b` -> `***-**-$1`.
  - Email pattern: mask local part leaving first and last characters.
  - Increment `PIIRedacted` for each pattern matched and replaced.

### Task 3: Group Recursion and Full Handler Implementation
- Handle `slog.KindGroup`:
  - Recursively iterate through child attributes and return a new group attribute with sanitized children.
- Implement the complete `slog.Handler` interface:
  - `Enabled(ctx context.Context, level slog.Level) bool`: delegate to `h.next.Enabled(ctx, level)`.
  - `WithAttrs(attrs []slog.Attr) slog.Handler`: return a new `PIIRedactingHandler` wrapping `h.next.WithAttrs(...)` with pre-sanitized attributes.
  - `WithGroup(name string) slog.Handler`: return a new handler wrapping `h.next.WithGroup(name)`.
  - `Handle(ctx context.Context, r slog.Record) error`: construct a fresh `slog.Record`, apply `maskString` to `r.Message`, sanitize each attribute via `r.Attrs()`, and forward to `h.next.Handle(ctx, sanitizedRecord)`.

---

## 6. Verification & Acceptance Criteria

Run all unit tests in the project directory:

```bash
go test -v -race ./curriculum/modules/12-observability-diagnostics/projects/pii-redactor/_starter
```

### Acceptance Checklist
- [ ] `TestPIIRedaction_SensitiveKeys`: verifies password, token, and authorization headers are masked.
- [ ] `TestPIIRedaction_CreditCardAndSSN`: verifies credit card and SSN patterns are masked to last 4 digits.
- [ ] `TestPIIRedaction_Email`: confirms email addresses are obfuscated.
- [ ] `TestPIIRedaction_NestedGroups`: validates that nested `slog.Group` attributes are traversed and sanitized.
- [ ] `TestPIIRedaction_BenignLogUnchanged`: ensures timestamps, numbers, and benign text are unaltered.
- [ ] Code passes under `-race` without data race reports.

---

## 7. Common Pitfalls & Anti-Patterns

1. **Mutating Inbound `slog.Record` Directly**:
   *Anti-pattern*: Attempting to overwrite slices inside `r` in place, which can cause race conditions if the caller retains references.
   *Correction*: Create a new `slog.Record` using `slog.NewRecord(r.Time, r.Level, maskedMessage, r.PC)` and add sanitized attributes.

2. **Compiling Regex on Every Log Call**:
   *Anti-pattern*: Calling `regexp.MustCompile` inside `Handle()` or `maskString()`.
   *Correction*: Compile regular expressions once at package initialization into global `var`s or handler struct fields.

3. **Breaking Handlers on `WithAttrs` and `WithGroup`**:
   *Anti-pattern*: Returning `h` unmodified in `WithAttrs`, causing attributes passed via `logger.With(...)` to bypass the redaction pipeline.
   *Correction*: Sanitize attributes passed to `WithAttrs` and return a new child handler instance wrapping `next.WithAttrs(...)`.

4. **Corrupting Non-String Attribute Types**:
   *Anti-pattern*: Converting int, float, or bool attributes to string representations during sanitization.
   *Correction*: Check `a.Value.Kind()` and only apply string pattern matching to `slog.KindString`.

---

## 8. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Code fails to compile or panics on nested groups. Sensitive keys are not sanitized or leaked in error messages. Regex expressions are recompiled on every invocation. |
| **Level 2: Advanced Beginner** | Sensitive keys are masked, but credit cards or emails in log messages leak into output. Handler does not implement `WithGroup` or drops group names. |
| **Level 3: Competent** | All sensitive keys, credit cards, SSNs, and emails are correctly redacted. Nested groups are sanitized recursively. Atomic stats counters track events accurately. All unit tests pass. |
| **Level 4: Proficient** | Flawless custom handler implementation. Benign types bypass regex paths with zero unnecessary allocations. High-performance string manipulation preserves auditability. Full compliance with `slog.Handler` contract. |

---

## 9. Submission & Reflection

1. Submit output from `go test -v -race ./...`.
2. Submit a sample JSON log output before and after passing through your redactor.
3. Answer: *Why is it critical that `WithAttrs` and `WithGroup` return a new handler instance rather than mutating the receiver?*
