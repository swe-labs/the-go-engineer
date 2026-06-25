# HTTP as a protocol

## Learning objective

Explain the HTTP request/response model, identify the purpose of each HTTP method (GET, POST, PUT, DELETE), read and write headers, and classify HTTP status codes by class (1xx-5xx).

## Why this matters

HTTP is the foundation of every web API, microservice communication, and browser-server interaction. Every REST API, GraphQL endpoint, and gRPC-to-HTTP proxy speaks HTTP. If you cannot reason about methods, status codes, and headers, you cannot debug a misbehaving API, design a clean endpoint, or secure a service. For Go engineers, understanding HTTP on the wire is essential because Go's `net/http` mirrors the protocol directly — there is no magic framework layer to hide behind.

## Mental model

Think of HTTP as a conversation between a client and a server. The client sends an envelope (the request) containing a verb (method), an address (URL), some labels (headers), and optionally a package (body). The server opens it, does work, and sends back an envelope (the response) with a status code (did it work?), labels (response headers), and optionally a package (body). The conversation is stateless — each request-response pair is independent, with no memory of previous ones.

## Core idea

HTTP is a text-based request-response protocol running over TCP (typically port 80) or TLS (port 443). Every HTTP message consists of:

- **A start line**: for requests, `METHOD path HTTP/version` (e.g., `GET /api/users HTTP/1.1`). For responses, `HTTP/version status-code reason` (e.g., `HTTP/1.1 200 OK`).
- **Headers**: key-value pairs like `Content-Type: application/json`.
- **An empty line**: `\r\n` separating headers from body.
- **An optional body**: arbitrary bytes.

The four primary HTTP methods map to CRUD operations:

| Method | CRUD | Idempotent | Safe | Body expected |
|--------|------|------------|------|---------------|
| GET | Read | Yes | Yes | No |
| POST | Create | No | No | Yes |
| PUT | Replace | Yes | No | Yes |
| DELETE | Delete | Yes | No | Maybe |

Status codes fall into five classes:

| Class | Range | Meaning | Example |
|-------|-------|---------|---------|
| 1xx | 100-199 | Informational | 101 Switching Protocols |
| 2xx | 200-299 | Success | 200 OK, 201 Created |
| 3xx | 300-399 | Redirection | 301 Moved Permanently |
| 4xx | 400-499 | Client Error | 404 Not Found, 400 Bad Request |
| 5xx | 500-599 | Server Error | 500 Internal Server Error |

## Under the hood

When a TCP connection is established, the client writes the HTTP request as plain text (or TLS-encrypted text) to the socket. The server reads the request line, parses headers, and streams the body. Go's `net/http` handles this parsing internally, converting raw bytes into an `*http.Request` struct. After HTTP/1.1, the connection is kept alive by default (keep-alive) so multiple requests reuse the same TCP socket. HTTP/2 multiplexes multiple concurrent streams over a single connection, reducing head-of-line blocking. HTTP/2 is binary, not text, and supports server push, but the request-response semantics remain identical from the application layer.

## How Go uses it

The `net/http` package exposes HTTP concepts as Go types:

- `http.Request` models the incoming request with fields for `Method`, `URL`, `Header`, `Body`, and more.
- `http.ResponseWriter` is the interface for building the response.
- `http.StatusOK`, `http.StatusNotFound` etc. are named constants for status codes.
- `http.Header` is a `map[string][]string` for multi-valued headers.

Go also provides `http.Get`, `http.Post`, `http.NewRequest` on the client side, all of which construct and send HTTP messages over the wire.

## Go example

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
)

// requestSummary builds a human-readable summary of an HTTP request.
func requestSummary(method, url, body string) string {
	var b strings.Builder
	b.WriteString(method + " " + url + " HTTP/1.1\n")
	b.WriteString("Host: api.example.com\n")
	b.WriteString("Content-Type: application/json\n")
	if body != "" {
		b.WriteString(fmt.Sprintf("Content-Length: %d\n", len(body)))
		b.WriteString("\n" + body)
	}
	return b.String()
}

// parseStatusClass returns the human-readable class of an HTTP status code.
func parseStatusClass(code int) string {
	switch {
	case code >= 100 && code < 200:
		return "Informational"
	case code >= 200 && code < 300:
		return "Success"
	case code >= 300 && code < 400:
		return "Redirection"
	case code >= 400 && code < 500:
		return "Client Error"
	case code >= 500 && code < 600:
		return "Server Error"
	default:
		return "Unknown"
	}
}

func main() {
	req := requestSummary("POST", "/api/users", `{"name":"Alice"}`)
	fmt.Println("=== Request ===")
	fmt.Println(req)

	codes := []int{200, 201, 301, 400, 404, 500, 503}
	fmt.Println("\n=== Status Code Classes ===")
	for _, c := range codes {
		fmt.Printf("%d -> %s\n", c, parseStatusClass(c))
	}
}
```

## Step-by-step execution

1. `requestSummary` is called with method `"POST"`, path `"/api/users"`, body `{"name":"Alice"}`.
2. A `strings.Builder` accumulates the request line: `POST /api/users HTTP/1.1`.
3. Headers `Host` and `Content-Type` are written.
4. If a body is present, `Content-Length` is computed from `len(body)` and written.
5. An empty line `\n` separates headers from body, then the body is appended.
6. Back in `main`, `codes` is iterated and each status code is classified by `parseStatusClass`.
7. `200` falls into the `code >= 200 && code < 300` branch → `"Success"`.
8. `500` falls into `code >= 500 && code < 600` → `"Server Error"`.
9. Each result is printed.

## Common mistakes

- Using the wrong method: GET with a body for a mutation, or POST for a pure read. GET requests with bodies are technically allowed but most servers and proxies ignore or reject them.
- Confusing PUT with PATCH: PUT replaces the entire resource; PATCH applies a partial update. Sending a partial JSON object via PUT may delete omitted fields.
- Assuming 200 means "everything is fine": a 200 response can contain an application-level error in the body. Always check both the status code and the response payload.
- Forgetting that headers are case-insensitive: `Content-Type` and `content-type` are the same header. Go's `http.Header.Get` handles this, but raw comparisons will fail.
- Misclassifying 4xx vs 5xx: 4xx means the client sent a bad request (fix the client). 5xx means the server failed (fix the server). Returning 500 for validation errors is a common anti-pattern.

## Debugging walkthrough

A developer reports that their Go HTTP client gets a 403 Forbidden but the API key is correct.

```go
resp, err := http.Get("https://api.example.com/users")
fmt.Println(resp.StatusCode) // 403
```

**Investigation**: Use a tool like `curl -v` to see the full request/response:

```bash
$ curl -v https://api.example.com/users
> GET /users HTTP/2
> Host: api.example.com
> Authorization: Bearer <redacted>
>
< HTTP/2 403
< content-type: application/json
< {"error":"missing header X-API-Version"}
```

**Root cause**: The API requires an `X-API-Version` header. The client didn't set it.

**Fix**: Add the header to the request:

```go
req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)
req.Header.Set("X-API-Version", "2")
req.Header.Set("Authorization", "Bearer "+key)
resp, _ := http.DefaultClient.Do(req)
```

The 403 becomes a 200.

## Production notes

- Always set a `Content-Type` header when sending a body. Without it, clients and proxies may misinterpret bytes.
- Use `http.Client` with timeouts; the zero-value client has no timeout and can hang forever.
- HTTP/2 is negotiated during the TLS handshake (ALPN). Go's `http.Server` supports HTTP/2 automatically when using TLS.
- Rate limiters and API gateways often inspect method + path pairs. A `GET /api/orders` and a `POST /api/orders` may need different rate limits.

## Performance implications

- HTTP/1.1 keep-alive reuses TCP connections, avoiding three-way handshake overhead on subsequent requests.
- HTTP/2 multiplexing allows multiple concurrent requests on one connection, reducing latency in high-throughput systems.
- Large headers (e.g., oversized JWTs, cookie blobs) add round-trip overhead since they are sent with every request. Keep headers under 8 KB where possible.
- Go's `http.Client` reuses connections via the built-in `http.Transport` connection pool. Creating a new `http.Client` per request defeats this pool.

## Practice task

Write a function `buildHTTPRequest(method, path string, headers map[string]string, body string) string` that returns the full raw HTTP/1.1 request text (start line, headers, empty line, body). Then write a function `classifyStatus(code int) string` that returns the status class name. In `main`, call both with at least three different scenarios and print the results. Run and verify the output.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/01-http-as-a-protocol
go test ./curriculum/modules/08-http-rest-apis/lessons/01-http-as-a-protocol
```

## Review questions

1. What is the difference between a 401 Unauthorized and a 403 Forbidden response?
2. Why is PUT considered idempotent but POST is not? Give an example.
3. What happens if you send a `Content-Length` header that is smaller than the actual body size?
4. List the five HTTP status code classes and give one example code for each.
5. How does HTTP/2 improve performance over HTTP/1.1?

## NEXT UP

net/http basics — building your first Go HTTP server with `http.HandleFunc` and `http.ListenAndServe`.
