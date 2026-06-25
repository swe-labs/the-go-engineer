# HTTP Client

## Learning objective

Use `http.Client` to make HTTP requests, configure timeouts and transport settings, manage connection pooling, and handle responses with proper error checking.

## Why this matters

Most Go services communicate with other services over HTTP. Making external API calls is a fundamental task: calling a payment gateway, fetching user data from an auth service, posting metrics to a monitoring system. The `http.Client` is your tool for all of these. Misconfiguring it causes mysterious timeouts, connection leaks, and cascading failures. Getting it right is essential for reliable distributed systems.

## Mental model

Think of `http.Client` as a phone operator in a switchboard:

- The operator (Client) manages a pool of phone lines (connections).
- When you ask them to call someone (make a request), they check if there's already an open line to that destination (keep-alive connection reuse).
- If not, they plug in a new line (new TCP connection).
- If the line is dead (timeout), they tell you the call failed.
- The operator has rules: "hang up after 90 seconds of silence" (IdleConnTimeout), "don't start a new call if the last one took too long" (Timeout).
- You can cancel a call at any time (context cancellation), and the operator will hang up and free the line.

## Core idea

`http.Client` has three key configuration areas:

### 1. Timeout

```go
client := &http.Client{
    Timeout: 10 * time.Second, // overall request timeout
}
```

`Timeout` covers the entire request-response cycle: dial, TLS handshake, headers, body. If exceeded, the request is cancelled. It's the simplest and most important setting.

### 2. Transport

```go
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout: 5 * time.Second,
    },
}
```

- `MaxIdleConns`: total idle connections in the pool across all hosts.
- `MaxIdleConnsPerHost`: idle connections kept alive for a single host.
- `IdleConnTimeout`: how long an idle connection stays in the pool.
- `DialContext`: how to establish TCP connections (with dial timeout and keep-alive).
- `TLSHandshakeTimeout`: timeout for TLS negotiation.

### 3. Request creation

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
req.Header.Set("Accept", "application/json")
resp, err := client.Do(req)
```

Always use `NewRequestWithContext` to pass a context. The context controls cancellation and deadlines. `client.Do` executes the request.

## Under the hood

When you call `client.Do(req)`:

1. The client checks the request URL, applies redirect policy, and calls `Transport.RoundTrip`.
2. The transport looks up an idle connection in the pool for the host.
3. If found, it reuses the connection. If not, it calls `DialContext` to open a new TCP connection.
4. For HTTPS, a TLS handshake occurs.
5. The request is written to the connection.
6. The response is read (headers first, then body).
7. After the response body is fully read and closed, the connection returns to the idle pool.

The `http.Client` is safe for concurrent use. Multiple goroutines can share a single client.

## How Go uses it

The default `http.DefaultClient` is used by `http.Get`, `http.Post`, etc. However, the default client has **no timeout**, which is dangerous in production. Always create a custom client.

Popular patterns:
- One shared client per service (reuse connections across requests).
- Different clients for different backends (different timeout requirements).
- `http.Client` with custom `Transport` for HTTP/2, proxy, or mTLS.

The Go standard library and ecosystem (AWS SDK, Google Cloud Client Libraries) all accept `*http.Client` for transport configuration.

## Go example

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/api/users", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return users, nil
}

func (c *Client) CreateUser(ctx context.Context, name string) (*User, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]string{"name": name})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/users", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func main() {
	client := NewClient("http://localhost:8080", 10*time.Second)
	users, _ := client.GetUsers(context.Background())
	log.Printf("Users: %+v", users)
}
```

## Step-by-step execution

Calling `client.GetUsers(ctx)` against a remote server:

1. `NewRequestWithContext` creates a GET request for `http://localhost:8080/api/users` with the given context.
2. `client.Do(req)` starts the transport round-trip.
3. Transport checks the idle connection pool for `localhost:8080`. Initially empty, so it calls `DialContext`.
4. `DialContext` opens a TCP connection to `localhost:8080` with a 5-second timeout.
5. Request is written to the connection.
6. Server responds. Transport reads the response (status, headers, body).
7. `resp.Body` is passed back. The client reads and decodes it.
8. `defer resp.Body.Close()` signals the transport that the connection can be reused.
9. Connection goes back to the idle pool.

On the next call to the same host, the transport reuses the idle connection directly (no dial), making the request faster.

## Common mistakes

- Using `http.DefaultClient` (no timeout, no connection limits). Always use a custom client.
- Forgetting to close `resp.Body`. This leaks connections and eventually exhausts the port pool.
- Creating a new `http.Client` for every request. This defeats connection pooling. Reuse clients.
- Not passing a context. Without context, you cannot cancel requests on shutdown or set deadlines.
- Ignoring non-2xx status codes. Always check `resp.StatusCode` and return meaningful errors.
- Not wrapping errors with context. Use `fmt.Errorf("...: %w", err)` for debuggable errors.

## Debugging walkthrough

Calls to an external API hang for 60 seconds before failing:

```go
client := &http.Client{} // no Timeout set
```

**Symptom**: The request blocks for 60+ seconds, then returns a connection timeout from the OS.

**Investigation**: No `Timeout` is set on the client. The default `http.Transport` has no timeouts either. The OS TCP timeout (default ~60s) eventually fires.

**Fix**: Set a reasonable client timeout:

```go
client := &http.Client{Timeout: 10 * time.Second}
```

Also configure transport timeouts:

```go
client.Transport = &http.Transport{
    DialContext: (&net.Dialer{
        Timeout: 5 * time.Second,
    }).DialContext,
    TLSHandshakeTimeout: 5 * time.Second,
}
```

## Production notes

- Create one `http.Client` per external service with appropriate timeouts and transport settings.
- Monitor client metrics: request duration, error rate, connection pool size.
- Use `http.Client` with `Timeout` set to something reasonable (5-30s depending on the service).
- Configure `MaxIdleConnsPerHost` to avoid starving other hosts when one host is very active.
- Use `context.WithTimeout` for per-request deadlines, overriding the client-level timeout.
- For service mesh environments, ensure the transport respects proxy settings via `HTTP_PROXY`.

## Performance implications

- Connection reuse is critical for performance. Establishing a new TCP connection takes 1-3 round trips (TCP + TLS). Reusing an idle connection avoids this entirely.
- `MaxIdleConns` should be tuned to your expected concurrency. Too low and connections are created/destroyed frequently. Too high and idle connections waste memory.
- Each idle connection uses ~3 KB of memory. 1000 idle connections = ~3 MB.
- HTTP/2 multiplexes multiple requests over a single connection, reducing the need for many idle connections.
- The dial timeout protects against slow DNS and network failures.

## Practice task

1. Write an HTTP client for a `Task` API with endpoints: `GET /tasks`, `POST /tasks`, `GET /tasks/{id}`.
2. Configure the client with: 10s timeout, 50 max idle connections, 5 per host, 30s idle timeout.
3. Implement methods that take `context.Context` and return proper Go types.
4. Write tests using `httptest.NewServer` that validate: successful requests, error responses, timeouts, and context cancellation.
5. Verify the transport configuration in a test.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/21-http-client -v
```

Tests cover client methods with mock servers, error handling, timeout detection, transport configuration verification, and context cancellation.

## Review questions

1. What is the purpose of `http.Client.Timeout`? What parts of the request lifecycle does it cover?
2. Why should you reuse an `http.Client` instead of creating a new one per request?
3. What is the role of `http.Transport` in connection pooling? What does `MaxIdleConnsPerHost` control?
4. Why must `resp.Body` always be closed? What happens if it isn't?
5. How does context cancellation interact with an HTTP request in progress?

## NEXT UP

Congratulations on completing Module 08! You now understand HTTP, building REST APIs with net/http, and production server concerns. Next up: Module 09 -- SQL, PostgreSQL, and Persistence.
