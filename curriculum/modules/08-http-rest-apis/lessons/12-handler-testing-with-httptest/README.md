# Handler testing with httptest

## Learning objective

Write table-driven tests for HTTP handlers using `httptest.NewRecorder` and `httptest.NewServer`, test middleware in isolation, and use golden files to assert exact responses.

## Why this matters

Handlers are the entry point for every HTTP request. If you cannot test them thoroughly, you ship bugs. The `net/http/httptest` package gives you a full simulation environment without real port binding -- no flaky CI, no race conditions, no firewall exceptions. Every production Go service depends on httptest for correctness.

## Mental model

Think of httptest as a recording studio. Your handler is the performer. `httptest.NewRecorder` is a tape recorder that captures the response instead of sending it over a network. `httptest.NewServer` is a full soundstage -- it boots a real HTTP server on a random port, so your client code hits a real TCP socket, but the setup and teardown are instant. You are the producer who plays test scenarios through the system and inspects the recorded output.

## Core idea

The `httptest.ResponseRecorder` implements `http.ResponseWriter` and records everything your handler writes:

- **Status code** via `recorder.Code`
- **Headers** via `recorder.Header()`
- **Body** via `recorder.Body.String()`

`httptest.NewServer` creates a real `http.Server` listening on `127.0.0.1:0` (random port). It auto-closes when the test finishes (if you call `defer server.Close()`). Use it when you need end-to-end testing including client-side `http.Client` behavior.

Golden files store expected responses on disk. Compare handler output against the golden file in your test. To update, delete the golden file and re-run the test with `-update` flag support.

## Under the hood

`httptest.NewRecorder` returns a `*ResponseRecorder` which implements:

```go
type ResponseRecorder struct {
    Code      int           // the HTTP response code
    HeaderMap http.Header   // the response headers
    Body      *bytes.Buffer // the buffered response body
    Flushed   bool          // whether Flush was called
}
```

When your handler calls `w.WriteHeader(code)`, the recorder stores it in `Code`. Header writes go to `HeaderMap`. Body writes go to `Body`. There is no real network I/O -- every call is an in-memory operation.

`httptest.NewServer` works differently. It starts a real `net.Listener` on a random port and serves requests over TCP. This exercises the full `http.Server` stack including `Content-Length` computation, chunked encoding, and connection reuse. The server's `*httptest.Server.Client()` method returns an `*http.Client` pre-configured to skip TLS verification (if using HTTPS).

## How Go uses it

The Go standard library itself uses httptest extensively. Every handler in `net/http` has corresponding tests in `net/http/httptest`. Third-party frameworks like `gorilla/mux`, `chi`, and `gin` all use httptest for their test suites.

Go's toolchain encourages testing via the `testing` package, and httptest is the natural complement for HTTP services. The `go test` command treats httptest tests identically to any other test -- no special build tags or flags needed.

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var users = []User{
	{ID: 1, Name: "Alice", Role: "admin"},
	{ID: 2, Name: "Bob", Role: "user"},
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

// --- Tests ---

func TestUsersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	usersHandler(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got []User
	json.NewDecoder(resp.Body).Decode(&got)
	resp.Body.Close()
	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}
}

func TestGreetingTableDriven(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"default", "", "Hello, world!"},
		{"named", "name=Go", "Hello, Go!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/greet?"+tt.query, nil)
			rec := httptest.NewRecorder()
			greetingHandler(rec, req)
			if got := strings.TrimSpace(rec.Body.String()); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWithServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", usersHandler)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
```

## Step-by-step execution

For `TestUsersHandler`:

1. `httptest.NewRequest` creates an `*http.Request` with method GET and path `/api/users`. No real network connection is made -- the request is just a struct in memory.
2. `httptest.NewRecorder` creates a `*ResponseRecorder` that satisfies `http.ResponseWriter`.
3. `usersHandler(rec, req)` calls the handler directly (not through a server). The handler calls `w.Header().Set(...)` which writes to `rec.HeaderMap`. Then `json.NewEncoder(w).Encode(users)` marshals the struct and writes bytes to `rec.Body`.
4. `rec.Result()` returns an `*http.Response` populated from the recorder's internal state.
5. The test decodes the body and asserts the user count.

For `TestWithServer`:

1. `httptest.NewServer(mux)` starts a real HTTP server on a random port.
2. `server.URL` contains the full base URL like `http://127.0.0.1:34567`.
3. `http.Get(server.URL + "/api/users")` makes a real TCP connection.
4. The server routes the request to `usersHandler`, which runs normally.
5. The response is sent back over the socket and `http.Get` returns.
6. `defer server.Close()` shuts down the listener when the test finishes.

## Common mistakes

- Forgetting to close `resp.Body` in server tests, causing resource leaks. Always `defer resp.Body.Close()`.
- Using `httptest.NewServer` when a `NewRecorder` would suffice. NewServer is slower and unnecessary for unit-testing handler logic.
- Not calling `rec.Result().Body.Close()` -- response recorder bodies must be closed, though they consume no network resources.
- Comparing raw body strings without trimming whitespace. `fmt.Fprintf` appends no newline, but some assertions add invisible trailing spaces.
- Hard-coding server ports. Always use `server.URL` -- never guess the port.

## Debugging walkthrough

Consider this failing test:

```go
func TestBrokenGreeting(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/greet", nil)
    rec := httptest.NewRecorder()
    greetingHandler(rec, req)
    // expected "Hello, world!" but got ""
}
```

**Symptom**: Body is empty.

**Investigation**: Add `t.Logf("body=%q", rec.Body.String())` after the handler call. Output shows `""`. Check if the handler uses `r.URL.Query().Get("name")` -- it does. Check the request URL: `/greet` without query string. The handler falls through to `name == ""` and sets `name = "world"`, so `fmt.Fprintf(w, "Hello, %s!", name)` should write. Suspect the recorder isn't capturing -- but it is. The problem is that the test calls the handler directly but the request path might not match. Actually, the handler does work -- the bug is that the test never reads the body. Ensure assertions check `rec.Body.String()`.

**Root cause**: The test as stated has a logical bug -- the handler works, but the test doesn't assert body content. Always inspect the recorder body before blaming the handler.

## Production notes

- Use golden files for complex JSON/HTML responses. Store them in `testdata/*.golden`. When the API changes, update golden files intentionally.
- Never use `httptest.NewServer` in unit tests that only check handler logic. Reserve it for integration tests that exercise the full HTTP stack.
- Mock external dependencies (databases, third-party APIs) in httptest tests. The handler test should not depend on external services.
- Run tests with `go test -v -run TestName` to see individual sub-test results.
- Use `t.Parallel()` on table-driven tests to speed up execution, but not when using `httptest.NewServer` without a fresh server per sub-test.

## Performance implications

- `httptest.NewRecorder` is pure in-memory -- essentially free. Microseconds per test.
- `httptest.NewServer` allocates a TCP port and goroutine. Startup takes ~1ms. Reuse one server per test suite instead of one per test case.
- Golden file reads hit disk once per test file. The OS caches the file, so subsequent runs are fast.
- Parallel tests with shared golden files should read the file once and share the reference data via `sync.Once` or a global.

## Practice task

1. Write a handler `weatherHandler` that reads a `city` query param and returns JSON: `{"city":"London","temp_c":15}`. Return 400 if city is empty.
2. Write a table-driven test using `httptest.NewRecorder` covering: valid city, missing city, empty city.
3. Write a golden file test for the valid response.
4. Write an integration test using `httptest.NewServer` that uses `http.Get` to call the handler.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/12-handler-testing-with-httptest -v
```

Expected output includes four passing tests: `TestUsersHandler`, `TestGreetingTableDriven/default`, `TestGreetingTableDriven/named`, `TestAdminMiddleware/*`, `TestWithServer`, `TestGoldenFile`.

## Review questions

1. What is the difference between `httptest.NewRecorder` and `httptest.NewServer`? When would you use each?
2. How does `ResponseRecorder` capture the status code, headers, and body of a response?
3. Why should golden files be stored in a `testdata` directory? What does `go test` do with that directory?
4. What happens if you forget to close the body returned by `rec.Result()`? Is it a resource leak?
5. How would you test middleware that checks for an authentication header? Show the pattern using `httptest.NewRecorder`.

## NEXT UP

Body size limits -- protecting your server from oversized request payloads using `http.MaxBytesReader`.
