# HTTP request and response preview

## Learning objective

Understand and apply HTTP request and response preview in the context of professional Go software engineering.

## Why this matters

Applications need a standard way to request and transfer data over the web. HTTP provides a universal protocol with well-defined methods (GET, POST, PUT, DELETE), status codes, headers, and content negotiation — understood by every web server, browser, and API client.

## Mental model

HTTP is a request-response protocol where a client sends a message (method + path + headers + body) and the server replies with another message (status code + headers + body). It is stateless — each request is independent, with no memory of previous requests.

## Core idea

HTTP is a text-based protocol built on TCP. A client sends a request with a method (GET, POST), a path, headers, and an optional body. The server responds with a status code (200 OK, 404 Not Found), headers, and a body. Every web interaction follows this pattern.

## Under the hood

An HTTP/1.1 request looks like: `GET /path HTTP/1.1\r\nHost: example.com\r\n\r\n`. The server reads line by line until the empty line (`\r\n\r\n`), then reads Content-Length bytes for the body. Status codes are three-digit integers: 1xx (informational), 2xx (success), 3xx (redirection), 4xx (client error), 5xx (server error). Content-Type and Content-Length headers determine how the client interprets the response body.

## How Go uses it

Go's `net/http` package provides a complete HTTP client and server implementation. `http.HandleFunc` registers routes, `http.Request` represents incoming requests, `http.ResponseWriter` sends responses. The standard library's HTTP server is production-grade — used by Docker, Kubernetes, Prometheus, and many other projects.

## Go example

The example program starts an HTTP server on port 8080 with two handlers (`/` and `/api/hello`), then makes a GET request to the server and prints the response status, headers, and body.

## Step-by-step execution

1. A client creates an HTTP request with a method (GET), a URL (`https://api.example.com/users`), and optional headers (Authorization).
2. The client opens a TCP connection to the server and sends the request as text following the HTTP protocol format.
3. The server receives the request, parses it, and routes it to the registered handler for the `/users` path.
4. The handler processes the request, generates a response with a status code (200 OK) and a JSON body, and sends it back.
5. The client receives the response, reads the status code and body, closes the connection (or reuses it for keep-alive), and processes the data.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Thinking HTTP is the same as "the internet" | Confusing one application protocol with the entire network stack | HTTP is one protocol on top of TCP/IP — there are many others |
| Believing responses always contain HTML | Not understanding Content-Type negotiation | Responses can be JSON, XML, images, binary, or any MIME type |
| Assuming HTTP is always client-initiated | Not knowing about SSE and WebSockets | Server-sent events and WebSockets enable server-initiated communication |

## Debugging walkthrough

**Scenario: A POST request receives a 404 Not Found.**
- **Cause:** The URL path does not match any route registered on the server.
- **Diagnosis:** Check the server's route registration and compare with the request URL.
- **Resolution:** Ensure the request URL matches exactly, including trailing slashes and path parameters.

**Scenario: A client receives 403 Forbidden when accessing a resource.**
- **Cause:** The client lacks proper authentication or authorization.
- **Diagnosis:** Check the request headers for Authorization or Cookie.
- **Resolution:** Include proper credentials (API key, JWT token, session cookie) in the request headers.

## Production notes

HTTP is the foundation of the web — every API (REST, GraphQL), every webpage, every image load, and every AJAX request uses HTTP. REST APIs, microservices, cloud services, and serverless functions all communicate over HTTP.

## Performance implications

- HTTP/1.1 allows multiple requests per connection (keep-alive), reducing connection overhead from ~1 RTT per request to ~1 RTT per connection.
- HTTP/2 multiplexes multiple requests over a single TCP connection, eliminating head-of-line blocking.
- Each HTTP request adds 200-500 bytes of header overhead — minimize requests and use caching headers for performance.

## Practice task

Write a Go program that starts an HTTP server on port 8080 with a single handler, then write a separate Go program that sends a GET request to `http://localhost:8080` and prints the response body and status code.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/15-http-request-and-response-preview/
```

## Review questions

1. What are the five components of an HTTP request?
2. What do HTTP status codes 200, 301, 400, 404, and 500 mean?
3. Use `curl` to make a GET request and a POST request. What is the difference in output?
4. What is the difference between HTTP/1.1 and HTTP/2?
5. What does the `Content-Type` header specify?

## NEXT UP

[Module 02 — Go Basics](../../../02-go-basics/README.md)
