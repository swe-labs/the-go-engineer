# HTTP request and response preview

## Learning objective

Deconstruct HTTP requests and responses, inspecting request methods, URL paths, headers, status codes, and message bodies.

## Why this matters

Almost every software system today communicates over HTTP: mobile apps communicating with cloud servers, microservices exchanging data within a cluster, web browsers loading interactive sites, and third-party webhooks delivering events. Understanding the raw wire format of HTTP — request lines, status codes, headers, and payloads — demystifies web development and makes debugging APIs second nature.

> [!NOTE]
> **Zero-Magic Reassurance:**
> The Go example below uses `http.NewServeMux()` and a background server goroutine to demonstrate a live, functional HTTP exchange. Full Go web server architecture and middleware design will be taught in depth in Module 14. For now, focus on the **HTTP wire protocol**: how plain-text request lines, headers, status codes, and JSON bodies travel over the network.

## Mental model

HTTP (Hypertext Transfer Protocol) is an application-level, text-based, request-response protocol built on top of a TCP stream:
- **Stateless:** The server treats each incoming request independently. The protocol itself has no built-in memory of previous requests.
- **Client Request:** The client opens a TCP connection and sends a formatted text message asking the server to perform an action on a specific resource URL.
- **Server Response:** The server processes the request and sends back a formatted text message with a numerical status code indicating whether the action succeeded or failed, followed by the response content.

## Visual model

```mermaid
flowchart TD
    subgraph RequestWire["1. HTTP Request (Wire Format Over TCP)"]
        direction TB
        REQ_LINE["Request Line: GET /api/v1/users HTTP/1.1"]
        REQ_HEADERS["Headers:
Host: api.example.com
User-Agent: Go-http-client/1.1
Accept: application/json
Authorization: Bearer secret-token"]
        REQ_EMPTY["Empty Line delimiter: \r\n"]
        REQ_BODY["Body: (Empty for GET requests)"]

        REQ_LINE --- REQ_HEADERS --- REQ_EMPTY --- REQ_BODY
    end

    subgraph ResponseWire["2. HTTP Response (Wire Format Over TCP)"]
        direction TB
        RES_LINE["Status Line: HTTP/1.1 200 OK"]
        RES_HEADERS["Headers:
Content-Type: application/json; charset=utf-8
Content-Length: 52
Date: Sat, 24 Oct 2026 17:00:00 GMT"]
        RES_EMPTY["Empty Line delimiter: \r\n"]
        RES_BODY["Body (JSON Payload):
[{\"id\": 1, \"name\": \"Gopher\", \"role\": \"engineer\"}]"]

        RES_LINE --- RES_HEADERS --- RES_EMPTY --- RES_BODY
    end

    RequestWire -->|Client sends over TCP Socket| ResponseWire
```

## Core idea

Beneath all web frameworks, REST APIs, and GraphQL endpoints, HTTP is fundamentally just formatted text sent over a TCP socket. A client sends a message starting with a method and URL, followed by key-value headers and an optional body. The server replies with a status code, headers, and a body. Understanding this raw wire format allows you to debug any API with simple tools like `curl`.

## Under the hood

### 1. The Wire Format Structure
Every HTTP/1.1 message adheres to a strict four-part structure separated by carriage-return line-feeds (`\r\n`):
1. **Start Line:**
   - For Requests: `METHOD /path/to/resource HTTP/1.1`
   - For Responses: `HTTP/1.1 STATUS_CODE REASON_PHRASE` (e.g., `HTTP/1.1 200 OK`)
2. **Headers:** Key-value metadata lines in the format `Name: Value\r\n`.
3. **Empty Line (`\r\n`):** A single blank line signaling the end of headers and the beginning of the body.
4. **Body (Payload):** Raw bytes of data (HTML, JSON, binary image) whose byte length matches the `Content-Length` header.

### 2. Standard HTTP Methods & REST Semantics
- `GET`: Retrieve representation of a resource. Safe (does not modify server state) and idempotent (repeating produces same result).
- `POST`: Submit data to be processed or create a new subordinate resource. Neither safe nor idempotent.
- `PUT`: Completely replace the target resource with request payload. Idempotent.
- `PATCH`: Apply partial modifications to a resource.
- `DELETE`: Remove the specified resource. Idempotent.

### 3. HTTP Status Code Categorization
Status codes are three-digit integers grouped into five functional classes:
- **1xx (Informational):** Request received, continuing process (e.g., `101 Switching Protocols` for WebSockets).
- **2xx (Success):** Action successfully received and processed:
  - `200 OK`: Standard success response.
  - `201 Created`: Resource successfully created (common after POST).
  - `204 No Content`: Action succeeded, no response body returned (common after DELETE).
- **3xx (Redirection):** Further action needed to complete request:
  - `301 Moved Permanently`: Resource permanently relocated to a new URL.
  - `304 Not Modified`: Cached copy on client is still fresh; saves bandwidth.
- **4xx (Client Error):** The request contains bad syntax or cannot be fulfilled:
  - `400 Bad Request`: Malformed JSON or invalid parameter syntax.
  - `401 Unauthorized`: Authentication required (missing or invalid API token).
  - `403 Forbidden`: Authenticated, but user lacks permission to access resource.
  - `404 Not Found`: Target URL does not exist on server.
  - `429 Too Many Requests`: Rate limit exceeded.
- **5xx (Server Error):** Server encountered an unhandled error:
  - `500 Internal Server Error`: Unhandled panic or exception in server code.
  - `502 Bad Gateway`: Reverse proxy (e.g., NGINX) received invalid response from backend.
  - `503 Service Unavailable`: Server overloaded or undergoing maintenance.

## How Go uses it

Go's standard library `net/http` package is considered one of the finest in computer science:
- `http.HandleFunc("/api/hello", handler)` registers URL route handlers.
- `http.ListenAndServe(":8080", nil)` starts a production-grade multi-threaded HTTP server.
- `http.Get("https://api.example.com")` issues an outgoing HTTP GET request.
- Unlike in Node.js or Python, Go's standard `net/http` server is production-ready out of the box — powers high-scale production systems at Google, Cloudflare, and Docker without requiring external web application servers.

## Go example

The example starts a local HTTP server with two routes (`/` and `/api/hello`), sends a real client GET request, and inspects the resulting HTTP status code, response headers, and body:

```go
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== HTTP Request / Response Lifecycle Preview ===")

	// 1. Setup local HTTP router with two endpoints
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Welcome to The Go Engineer HTTP Server!")
	})

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Hello from Go!", "status": "ok"}`)
	})

	// 2. Bind local listener on ephemeral port (port 0)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("Failed to bind listener: %v\n", err)
		return
	}
	defer listener.Close()

	server := &http.Server{Handler: mux}
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Shutdown(context.Background())

	serverURL := fmt.Sprintf("http://%s/api/hello", listener.Addr().String())
	fmt.Printf("Server listening actively on: %s\n\n", serverURL)

	// 3. Issue client HTTP GET request
	client := &http.Client{Timeout: 3 * time.Second}
	fmt.Printf("[Client] Sending GET request to: %s\n", serverURL)

	resp, err := client.Get(serverURL)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 4. Deconstruct HTTP Response
	fmt.Println("\n=== Deconstructing HTTP Response Wire Data ===")
	fmt.Printf("HTTP Protocol Version:  %s\n", resp.Proto)
	fmt.Printf("Status Line:            %s (Code: %d)\n", resp.Status, resp.StatusCode)
	fmt.Println("Headers:")
	for key, values := range resp.Header {
		fmt.Printf("  %s: %s\n", key, values[0])
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("\nResponse Body Payload:\n%s\n", string(bodyBytes))
	fmt.Println("HTTP Lifecycle complete: Client Request -> Server Route -> Status Code + Headers + Body.")
}
```

## Step-by-step execution

1. The client opens a TCP socket to the server's IP and port.
2. The client writes formatted text down the socket: `GET /api/hello HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n`.
3. The server reads the stream, splits at `\r\n\r\n` to parse headers, and inspects the request line method and path.
4. The router matches `/api/hello` and executes the registered handler function.
5. The handler sets response headers (`Content-Type: application/json`), writes status code `200`, and writes the JSON payload bytes.
6. The client reads the incoming response status line (`HTTP/1.1 200 OK`), parses headers, and consumes the response body stream.
7. Both sides close the socket or keep it alive for the next request.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Leaking HTTP response bodies | Forgetting that `resp.Body` is an open network stream | Always call `defer resp.Body.Close()` immediately after checking `err == nil` on HTTP requests |
| Writing headers after calling `w.Write()` | The first call to `w.Write()` automatically sends `200 OK` and headers to the wire | Always set headers with `w.Header().Set()` and `w.WriteHeader()` *before* writing the body |
| Confusing `401 Unauthorized` with `403 Forbidden` | Both deny access | `401` means "you are not logged in"; `403` means "you are logged in, but not allowed here" |
| Ignoring Content-Type headers | Sending JSON with default `text/plain` | Always set `Content-Type: application/json` so clients know how to parse the payload |

## Debugging walkthrough

**Scenario: A client sends a POST request with JSON to an API endpoint, but the server responds with `415 Unsupported Media Type` or fails to parse the body.**
- **Cause:** The client omitted the `Content-Type: application/json` request header. The server does not know how to deserialize the payload.
- **Diagnosis:** Run `curl -v -X POST http://localhost:8080/api/users -d '{"name":"alex"}'` and inspect the request headers in verbose output.
- **Resolution:** Explicitly supply the header: `curl -H "Content-Type: application/json" -X POST ...`.

**Scenario: A Go client makes HTTP requests in a tight loop and exhausts system file descriptors.**
- **Cause:** The program forgot `resp.Body.Close()`. The underlying TCP socket connection is never released back to the connection pool.
- **Diagnosis:** Monitor open sockets with `netstat` or `lsof`; see thousands of sockets stuck in `CLOSE_WAIT`.
- **Resolution:** Add `defer resp.Body.Close()` and ensure `io.Copy(io.Discard, resp.Body)` is called if you don't need the body.

## Production notes

Modern production systems design APIs around clear HTTP semantics. REST APIs use standard verbs and return proper error status codes (`404` for missing items, `409 Conflict` for duplicate emails, `422 Unprocessable Entity` for validation errors). Logging status codes in access logs enables alerting dashboards (e.g., alert when 5xx errors exceed 0.5% of total requests).

## Performance implications

- Each HTTP request has header overhead (~500–1000 bytes).
- HTTP/1.1 requires a new request per connection or sequential pipeline; HTTP/2 and HTTP/3 multiplex hundreds of requests concurrently over a single TCP/QUIC connection.
- Reading `resp.Body` lazily streams bytes as they arrive from the network; reading large files with `io.ReadAll()` loads the entire file into RAM, risking memory spikes.

## Tests / verification

1. Test an HTTP request directly using terminal `curl`:
   - `curl -i https://httpbin.org/get`
   - Observe the printed raw status line (`HTTP/2 200`), response headers (`content-type`, `date`), and JSON body.
2. Test error status codes:
   - `curl -i https://httpbin.org/status/404` (Observe `404 Not Found`)
   - `curl -i https://httpbin.org/status/500` (Observe `500 Internal Server Error`)
3. Run the Go HTTP example program:
   - `go run ./curriculum/modules/01-computers-terminal-git-web/15-http-request-and-response-preview`
   - Observe the live deconstructed HTTP headers and JSON payload.

## Review questions

1. What are the four structural parts of a raw HTTP request on the wire?
2. What is the fundamental difference between HTTP methods `GET` and `POST`?
3. What is the meaning of HTTP status codes in the 2xx, 4xx, and 5xx categories?
4. What is the difference between status code `401 Unauthorized` and `403 Forbidden`?
5. Why must Go developers always call `resp.Body.Close()` when making client HTTP requests?

## NEXT UP

Complete the [Shell + Git Lab](../../projects/shell-git-lab/README.md) project, then take the [Module 01 assessment](../../assessments/shell-git-lab/README.md).
