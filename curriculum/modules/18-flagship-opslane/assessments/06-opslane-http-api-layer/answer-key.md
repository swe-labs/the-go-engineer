# Opslane Milestone 06 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) To cap the inbound request payload to 1MB, preventing unbounded memory allocation and denial-of-service (OOM) attacks from maliciously oversized bodies.**  
*Explanation:* If an attacker streams a 2GB file to an unconstrained JSON endpoint, the server will buffer the data into RAM until memory is exhausted and the process is killed. `http.MaxBytesReader` cuts off the stream and returns an error as soon as 1MB is exceeded.

**2. b) The decoder returns an error (`json: unknown field "titel"`), triggering a 400 Bad Request response.**  
*Explanation:* By default, Go's JSON decoder discards unmatched keys silently. Calling `decoder.DisallowUnknownFields()` causes decoding to fail if any field in the JSON payload does not map to a struct field, immediately alerting clients to misspelled attributes.

**3. b) Using curly brace wildcard syntax like `"GET /api/v1/orders/{id}"`, accessed via `r.PathValue("id")`.**  
*Explanation:* Go 1.22 introduced enhanced pattern routing to `net/http.ServeMux`. Path parameters enclosed in curly braces are extracted using `r.PathValue(paramName)`.

**4. c) HTTP 204 No Content.**  
*Explanation:* Under RFC 9110, HTTP 204 indicates that the server successfully fulfilled the request and there is no additional content to send in the response payload.

**5. b) It provides an end-to-end correlation ID attached to request and response headers, allowing operations teams to trace a specific request across distributed microservices and log aggregators.**  
*Explanation:* In distributed architectures, tracking a single user interaction across gateways, services, and databases requires a common identifier. The `X-Request-ID` is passed along all network hops and logged in every log line.

**6. b) HTTP 400 with an `APIErrorEnvelope` containing code `"VALIDATION_FAILED"` and an array of field-level issues.**  
*Explanation:* Well-designed REST APIs return structured error envelopes specifying the error code and an array of individual field errors (`field` + `issue`) so API clients know exactly which fields to fix.

---

## Short Answer (2 points each)

**7. Streaming `json.NewDecoder` vs. `io.ReadAll` + `json.Unmarshal`:**
- `io.ReadAll`: Allocates a contiguous byte buffer in memory to hold the entire request body before unmarshaling. For large requests, this causes double allocation (raw bytes + unmarshaled Go structs) and high GC pressure.
- `json.NewDecoder`: Streams data directly from the network socket into struct fields incrementally. When paired with `http.MaxBytesReader`, it prevents unbounded memory buffering while failing immediately if size limits are breached.
*Grading: 1 pt for explaining memory buffering / allocation overhead; 1 pt for explaining streaming efficiency and early termination.*

**8. Architectural benefits of standardized API envelopes:**
1. **Predictable Client Contracts:** API consumers (web apps, mobile apps, SDKs) can implement a single generic HTTP client wrapper. Every response checks for `data` or `error` at the root, streamlining error handling across hundreds of endpoints.
2. **Field-Level Error Mapping:** The `{ "error": { "details": [...] } }` envelope provides field-specific errors, allowing frontend forms to highlight invalid input fields directly rather than displaying a generic alert.
*Grading: 1 pt for predictable client deserialization; 1 pt for UI form error mapping.*

**9. Defense value of `DisallowUnknownFields()`:**
Without `DisallowUnknownFields()`, if a client sends `{ "role": "admin", "titl": "Order" }` to an endpoint where `titl` was misspelled, the server ignores `titl`, sets the order title to empty string `""`, and proceeds. In more serious scenarios, clients might send outdated, deprecated, or misspelled parameters (e.g. `"auto_renew": false` vs `"autorenew": false`), causing critical business logic to default to unintended values silently. Disallowing unknown fields forces clients to adhere strictly to the schema.
*Grading: 1 pt for explaining prevention of silent field dropping; 1 pt for highlighting business/security risks of defaulting.*

**10. Thread safety with `sync.RWMutex` in `APIServer`:**
- **Why RWMutex:** Read operations (fetching an order, listing orders) heavily outnumber write operations (creating or deleting orders). A standard `sync.Mutex` serializes all reads, creating a bottleneck. `sync.RWMutex` allows multiple concurrent readers to access the map simultaneously while guaranteeing exclusive access to writers.
- **Where Applied:** `RLock()` is acquired in `handleGetOrder` and `handleListOrders`; `Lock()` is acquired in `handleCreateOrder` and `handleDeleteOrder`.
*Grading: 1 pt for explaining concurrent read optimization; 1 pt for identifying which handlers use RLock vs Lock.*
