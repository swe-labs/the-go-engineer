# Status codes

## Learning objective

Choose the correct HTTP status code for every response type, distinguish between 4xx and 5xx classes, and avoid common misclassification mistakes like returning 200 for errors or 500 for validation failures.

## Why this matters

HTTP status codes are the first thing a client sees in a response. A correct status code tells the client whether to retry (429, 503), redirect (301, 302), fix the request (400, 422), or celebrate (200, 201). An incorrect status code causes clients to misbehave: retrying a 400 that will never succeed, caching a 200 error response, or treating a 500 validation error as a server bug. Getting status codes right is the cheapest way to make your API robust.

## Mental model

Status codes are traffic lights. Green (2xx): go ahead, everything worked. Yellow (3xx): look this way, the resource moved. Red (4xx): you made a mistake, fix it. Flashing red (5xx): the server is broken, wait and retry. A client obeys the traffic light without reading the sign (response body). If you show green when the road is closed (200 for an error), the client crashes.

## Core idea

HTTP status codes are three-digit integers grouped into five classes:

| Code | Meaning | When to use |
|------|---------|-------------|
| 200 OK | Success | GET, PUT, PATCH, DELETE (with body) |
| 201 Created | Resource created | POST |
| 204 No Content | Success, no body | DELETE, PUT (no body needed) |
| 301 Moved Permanently | Resource relocated | Legacy URL redirects |
| 302 Found | Temporary redirect | Post-login redirect |
| 304 Not Modified | Cached version is fresh | Conditional GET with `If-None-Match` |
| 400 Bad Request | Malformed syntax | Invalid JSON, missing required headers |
| 401 Unauthorized | Authentication required | Missing or invalid credentials |
| 403 Forbidden | Authenticated but not allowed | Insufficient permissions |
| 404 Not Found | Resource doesn't exist | Unknown ID, wrong path |
| 405 Method Not Allowed | Wrong HTTP method | POST on a GET-only endpoint |
| 409 Conflict | State conflict | Duplicate resource, version conflict |
| 422 Unprocessable Entity | Semantic validation failure | Invalid field values |
| 429 Too Many Requests | Rate limit exceeded | Client exceeded quota |
| 500 Internal Server Error | Unexpected server failure | Unhandled panic, DB connection lost |
| 502 Bad Gateway | Upstream failed | Proxy/gateway upstream error |
| 503 Service Unavailable | Server overloaded | Maintenance, rate limiting at LB level |

## Under the hood

`net/http` provides named constants for every standard status code: `http.StatusOK` (200), `http.StatusNotFound` (404), `http.StatusInternalServerError` (500), etc. These are defined in `net/http/status.go`. The `http.StatusText(code)` function returns the human-readable reason phrase (e.g., "Not Found" for 404). Go's HTTP server automatically includes the reason phrase in the status line. Status codes outside the 100-599 range are invalid and cause `WriteHeader` to silently ignore them (in Go versions < 1.20) or panic (in Go 1.20+ with `-d=checkptr`).

## How Go uses it

```go
// Success
w.WriteHeader(http.StatusOK)
w.WriteHeader(http.StatusCreated)
w.WriteHeader(http.StatusNoContent)

// Client errors
w.WriteHeader(http.StatusBadRequest)
w.WriteHeader(http.StatusNotFound)
w.WriteHeader(http.StatusUnprocessableEntity)

// Server errors
w.WriteHeader(http.StatusInternalServerError)
w.WriteHeader(http.StatusServiceUnavailable)
```

Always use the named constants, never bare integers. `w.WriteHeader(200)` is a code smell.

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	items   = map[string]Item{}
	nextID  int64
)

func createItem(w http.ResponseWriter, r *http.Request) {
	var input Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		http.Error(w, "name is required", http.StatusUnprocessableEntity)
		return
	}

	id := fmt.Sprintf("item-%d", atomic.AddInt64(&nextID, 1))
	input.ID = id
	items[id] = input

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	item, ok := items[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if _, ok := items[id]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	delete(items, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /items", createItem)
	mux.HandleFunc("GET /items/{id}", getItem)
	mux.HandleFunc("DELETE /items/{id}", deleteItem)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. POST `/items` with `{"name":"Widget"}` → status 201 Created, body includes generated ID.
2. POST `/items` with `{}` → status 422 Unprocessable Entity, `"name is required"`.
3. GET `/items/item-1` → status 200 OK, body with the item.
4. GET `/items/nonexistent` → status 404 Not Found.
5. DELETE `/items/item-1` → status 204 No Content, empty body.
6. DELETE `/items/item-1` again → status 404 Not Found (item already deleted).

## Common mistakes

- Returning 500 for validation errors. Validation errors are the client's fault → 4xx. Only 5xx when the server is broken.
- Returning 200 with `{"error":"..."}`. If there's an error, use the appropriate 4xx/5xx code. The client checks the status first.
- Returning 201 for GET requests. GET is idempotent and should return 200. 201 is for successful resource creation.
- Returning 204 with a body. 204 explicitly means "no content". If you have content, use 200. If you write a body after 204, Go silently drops it.
- Using 403 instead of 401 for missing authentication. 401 means "authenticate first". 403 means "you're authenticated but not allowed".
- Forgetting 429 for rate limiting. Without 429, clients don't know to back off and may retry aggressively.

## Debugging walkthrough

A mobile client keeps retrying a request that always fails:

**Symptom**: The client retries 10 times with exponential backoff, but the server always returns the same error.

**Root cause**: The server returns `500 Internal Server Error` for a validation failure (missing field). The client sees 5xx and retries because 5xx means "temporary server failure". The request will never succeed because the client sends the same invalid data.

**Fix**: Change the status code to `422 Unprocessable Entity`. The client receives 4xx, knows not to retry, and displays the error message to the user.

## Production notes

- Use `http.StatusXX` constants everywhere. Bare integer status codes are prone to typos and make code harder to review.
- Document every status code each endpoint can return in your API reference. Clients need to handle all of them.
- For 3xx redirects, always include a `Location` header with the target URL.
- For 429, include a `Retry-After` header so the client knows how long to wait.
- In microservices, translate downstream 5xx errors to 502 or 503 depending on whether the downstream is the origin (502) or temporarily unavailable (503).

## Performance implications

- Status codes are a single integer in the response. Zero performance impact.
- 304 Not Modified responses should have an empty body. Combined with ETags, they save bandwidth and server CPU by letting clients use cached responses.
- 204 No Content responses save bandwidth by omitting the body. Use them for DELETE and update operations where the client doesn't need a response body.
- 3xx redirects add an extra round trip. Avoid redirect chains (>3 hops) in production APIs.

## Practice task

Build a simple order management API with:
- `POST /orders` — creates an order, returns 201 with the order.
- `GET /orders/{id}` — returns the order or 404.
- `PUT /orders/{id}/cancel` — cancels an order. Return 200 if cancelled, 409 if already cancelled, 404 if not found.
- Use appropriate status codes for each case. Test with httptest, verifying both status codes and response bodies.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/09-status-codes
go test ./curriculum/modules/08-http-rest-apis/lessons/09-status-codes
```

## Review questions

1. What is the difference between 401 Unauthorized and 403 Forbidden?
2. When should you use 422 Unprocessable Entity instead of 400 Bad Request?
3. What status code should a POST endpoint return on success? What about a DELETE?
4. Why is returning 200 with an error body considered an anti-pattern?
5. What header should accompany a 429 Too Many Requests response?

## NEXT UP

Middleware pattern — wrapping `http.Handler` to add cross-cutting concerns like logging, recovery, and CORS.
