# Opslane Milestone 01 Assessment — Rubric

## Assessment Dimensions

### 1. Configuration & Fail-Fast Validation (30%)
- **Excellent (5 pts)**: Configuration parameters strictly typed; defaults provided for development; thorough bounds validation on numeric and duration types; production guardrails prevent default secret leakage.
- **Proficient (3 pts)**: Environment variables loaded and typed correctly, but minor bounds checks omitted.
- **Developing (1 pt)**: Relies on `os.Getenv` string lookups without typing or validation.
- **Unacceptable (0 pts)**: Configuration fails to load or panics on missing keys.

### 2. Structured Logging & Observability (25%)
- **Excellent (5 pts)**: `log/slog` handlers cleanly configured with JSON output for production and text for development; log levels correctly parsed and applied.
- **Proficient (3 pts)**: `log/slog` initialized, but environment-specific handler switching is missing or defaults to raw text.
- **Developing (1 pt)**: Uses unstructured `log.Print` alongside `slog`.
- **Unacceptable (0 pts)**: No logging instrumentation.

### 3. HTTP Server Reliability & Middleware (25%)
- **Excellent (5 pts)**: Health and readiness endpoints separated; panic recovery middleware cleanly returns HTTP 500 JSON payloads; server declares explicit read/write timeouts.
- **Proficient (3 pts)**: Health routes present, but server lacks timeouts or recovery middleware leaves client hanging.
- **Developing (1 pt)**: Unhandled panics crash the process; endpoints lack proper status codes.
- **Unacceptable (0 pts)**: Server does not compile or bind to port.

### 4. Conceptual & Architecture Defense (20%)
- **Excellent (5 pts)**: Accurately explains runtime panic recovery, port privilege boundaries, fail-fast mechanics, and `ErrServerClosed` lifecycle.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in runtime mechanics.
- **Developing (1 pt)**: Cannot articulate why production requires non-default secrets or how `recover()` works.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
