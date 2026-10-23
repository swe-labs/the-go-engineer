# Opslane HTTP API layer

## Mission

Build an enterprise-grade RESTful HTTP API layer for Opslane using Go 1.22+ `net/http` enhanced routing, payload size hardening (`http.MaxBytesReader`), strict JSON decoding (`DisallowUnknownFields`), field-level validation error envelopes, and request tracing via `X-Request-ID`.

## Prerequisites

- `opslane-05`

## Mental Model

An HTTP API layer serves as the public translation boundary between external clients and internal domain services. Production APIs require uniform contracts: successful responses are wrapped in a standard `{ "data": ... }` envelope, while failures emit a consistent `{ "error": { "code": "...", "message": "...", "details": [...] } }` envelope. To protect the service from malicious or malformed input, every inbound request enforces a 1MB payload ceiling via `http.MaxBytesReader`, rejects unknown fields to prevent silent parameter ignoring, and tags every request with a correlated `X-Request-ID` header.

## Visual Model

```text
Inbound HTTP Request
  │ (POST /api/v1/orders)
  ▼
┌─────────────────────────────────────────────────────────────┐
│ requestIDMiddleware                                         │
│   Extracts or generates X-Request-ID -> sets response header│
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ handleCreateOrder Pipeline                                  │
│   ├─► http.MaxBytesReader(w, r.Body, 1MB) (Memory guard)    │
│   ├─► json.Decoder.DisallowUnknownFields()                  │
│   ├─► req.Validate() -> []FieldError                        │
└──────────────┬───────────────────────────────┬──────────────┘
               │ Valid                         │ Validation Error (400)
               ▼                               ▼
┌──────────────────────────────┐ ┌────────────────────────────┐
│ 201 Created                  │ │ 400 Bad Request            │
│ { "data": { "id": "ord-1" } }│ │ { "error": { "code":... } }│
└──────────────────────────────┘ └────────────────────────────┘
```

## Machine View

Under the hood:
1. `http.MaxBytesReader` wraps the request body reader. If a client attempts to stream more than 1MB, the reader returns `http.MaxBytesError`, which the handler catches to return HTTP 413 Payload Too Large before allocating excessive heap memory.
2. `decoder.DisallowUnknownFields()` detects typos in client request payloads (e.g. sending `"titel"` instead of `"title"`), returning an immediate decode error rather than silently populating default zero values.
3. Path routing leverages Go 1.22+ pattern matching (`"GET /api/v1/orders/{id}"`), using `r.PathValue("id")` to extract path parameters cleanly without external router dependencies.
4. `DELETE /api/v1/orders/{id}` returns HTTP 204 No Content upon successful deletion, conforming to REST standards.

## Run Instructions

```bash
# Run API layer test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/06-opslane-http-api-layer

# Launch HTTP server
go run ./curriculum/modules/18-flagship-opslane/assessments/06-opslane-http-api-layer
```

## Code Walkthrough

1. **API Envelopes (lines 18–33)**: Defines `APIResponse` and `APIErrorEnvelope` with field-level `FieldError` tracking.
2. **DTO & Validation (lines 36–70)**: `CreateOrderRequest.Validate()` inspects string length, positive amounts, and currency whitelisting (`USD`, `EUR`, `GBP`).
3. **`APIServer` & Router (lines 73–94)**: In-memory order map protected by `sync.RWMutex`, wired to Go 1.22+ ServeMux patterns.
4. **`requestIDMiddleware` (lines 96–107)**: Checks for `X-Request-ID`; generates 8-byte hex string from `crypto/rand` if absent; sets response header.
5. **Handlers (lines 109–202)**:
   - `handleCreateOrder`: Enforces 1MB limit, decodes without unknown fields, validates, and stores order with mutex protection.
   - `handleGetOrder`: Reads `r.PathValue("id")`; returns 404 or 200 with data envelope.
   - `handleListOrders`: Supports `?limit=` query parameter; defaults to 10 orders.
   - `handleDeleteOrder`: Deletes order; returns 204 No Content.
6. **Response Helpers (lines 204–220)**: `respondJSON` and `respondError` ensure standard envelope formatting and content headers.

## Try It

1. Run the test suite: verify that valid creation returns 201, invalid fields return 400 with details, and unknown IDs return 404.
2. Test unknown fields: send `{ "title": "Test", "amount": 10, "foo": "bar" }` and verify `MALFORMED_JSON` rejection.
3. Send a payload exceeding 1MB: verify `PAYLOAD_TOO_LARGE` (413) is returned.

Common mistakes:
- Forgetting `decoder.DisallowUnknownFields()`, allowing client typos to go unnoticed.
- Reading the entire request body into memory with `io.ReadAll` without `MaxBytesReader`, making the server vulnerable to OOM crashes.
- Returning inconsistent error payloads across different endpoints.

## In Production

Enterprise REST APIs must maintain predictable error shapes so frontend and mobile clients can programmatically map server field validation failures directly to UI form inputs. `X-Request-ID` propagation across the entire call stack enables distributed log tracing in APM tools (Datadog, Honeycomb).

## Thinking Questions

1. Why should APIs use `http.MaxBytesReader` rather than relying on load balancer payload size limits alone?
2. What are the client usability benefits of returning an array of field-level errors rather than failing on the first invalid field?
3. What is the semantic difference between HTTP 200 OK with an empty body versus HTTP 204 No Content?

## Next Step

Proceed to [Milestone 07 — Order Processing](../07-opslane-order-processing/README.md) to implement the order lifecycle state machine, transitions, and business calculations.
