# Opslane Milestone 06 Assessment — Rubric

## Assessment Dimensions

### 1. Request Hardening & Decoding Safety (30%)
- **Excellent (5 pts)**: `http.MaxBytesReader` caps inbound payloads; `DisallowUnknownFields` prevents silent errors; decodes via streaming `json.Decoder`; handles `MaxBytesError` with 413.
- **Proficient (3 pts)**: Payload decoded properly, but `DisallowUnknownFields` or `MaxBytesReader` omitted.
- **Developing (1 pt)**: Uses `io.ReadAll` without bounds, making server vulnerable to memory exhaustion.
- **Unacceptable (0 pts)**: Request decoding panics or fails on valid JSON.

### 2. Envelope Consistency & Validation Reporting (25%)
- **Excellent (5 pts)**: Standard `{ "data": ... }` and `{ "error": ... }` envelopes consistently returned; field validation collects all errors rather than failing on the first; appropriate HTTP status codes (201, 204, 400, 404, 413).
- **Proficient (3 pts)**: Envelopes used, but status codes inconsistent (e.g. returning 200 for created orders or deletions).
- **Developing (1 pt)**: Raw unstructured strings returned on errors.
- **Unacceptable (0 pts)**: No error envelopes or validation.

### 3. Routing & Correlation Tracing (25%)
- **Excellent (5 pts)**: Routes use Go 1.22+ method + path patterns (`"POST /api/v1/orders"`, `r.PathValue("id")`); `requestIDMiddleware` generates/propagates `X-Request-ID`.
- **Proficient (3 pts)**: Routing works, but `X-Request-ID` is missing or uses predictable random seeds.
- **Developing (1 pt)**: Relies on manual URL parsing instead of pattern matching.
- **Unacceptable (0 pts)**: Router fails to mount endpoints.

### 4. Conceptual & API Architecture Defense (20%)
- **Excellent (5 pts)**: Accurately explains streaming decoding vs buffered reads, `sync.RWMutex` contention dynamics, field-level validation contracts, and distributed correlation tracing.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in `RWMutex` or streaming decoding.
- **Developing (1 pt)**: Cannot articulate why `MaxBytesReader` is necessary or how 204 differs from 200.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
