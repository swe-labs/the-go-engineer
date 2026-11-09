# Opslane Milestone 06 Assessment — HTTP API Layer: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** What is the purpose of `r.Body = http.MaxBytesReader(w, r.Body, 1048576)` in `handleCreateOrder`?
- a) To limit the execution duration of the request to 1 millisecond.
- b) To cap the inbound request payload to 1MB, preventing unbounded memory allocation and denial-of-service (OOM) attacks from maliciously oversized bodies.
- c) To compress the incoming JSON payload.
- d) To limit the database row count to 1,048,576 rows.

**2.** What happens when a client sends an unrecognized field (e.g. `"titel": "Widget"`) to an endpoint configured with `decoder.DisallowUnknownFields()`?
- a) Go silently ignores the unrecognized field and decodes the rest.
- b) The decoder returns an error (`json: unknown field "titel"`), triggering a 400 Bad Request response.
- c) The Go program panics.
- d) The unknown field is automatically stored in a separate table.

**3.** In Go 1.22+, how are path parameters defined in `http.ServeMux` routes?
- a) Using colon syntax like `:id`.
- b) Using curly brace wildcard syntax like `"GET /api/v1/orders/{id}"`, accessed via `r.PathValue("id")`.
- c) Using regular expressions inside string concatenation.
- d) Path parameters require third-party packages like gorilla/mux.

**4.** What is the HTTP status code specified by REST standards for a successful `DELETE` request that returns no response body?
- a) HTTP 200 OK.
- b) HTTP 201 Created.
- c) HTTP 204 No Content.
- d) HTTP 202 Accepted.

**5.** What is the role of `X-Request-ID` in `requestIDMiddleware`?
- a) It authenticates the user's password.
- b) It provides an end-to-end correlation ID attached to request and response headers, allowing operations teams to trace a specific request across distributed microservices and log aggregators.
- c) It prevents caching in web browsers.
- d) It encrypts the HTTP body.

**6.** When request validation fails in `handleCreateOrder`, what does the server return?
- a) HTTP 500 with a stack trace.
- b) HTTP 400 with an `APIErrorEnvelope` containing code `"VALIDATION_FAILED"` and an array of field-level issues.
- c) HTTP 404 Not Found.
- d) An empty JSON object `{}`.

---

## Short Answer (2 points each)

**7.** Contrast streaming JSON decoding via `json.NewDecoder(r.Body).Decode(&req)` with `io.ReadAll(r.Body)` followed by `json.Unmarshal`. Why is `json.NewDecoder` combined with `MaxBytesReader` superior for production servers?

**8.** Explain why production APIs standardize on consistent response envelopes (e.g. `{ "data": ... }` and `{ "error": ... }`). How does this contract simplify frontend client development?

**9.** Why is `DisallowUnknownFields()` considered an important API security and reliability defense against subtle client bugs?

**10.** Explain the thread-safety design of `APIServer`. Why is `sync.RWMutex` used instead of `sync.Mutex`, and where are read locks (`RLock`) versus write locks (`Lock`) acquired?
