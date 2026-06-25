# Response writing

## Learning objective

Write HTTP responses using `w.Write`, `w.WriteHeader`, `json.NewEncoder(w)`, set the `Content-Type` header correctly, and implement streaming responses for large payloads.

## Why this matters

The response is what the client sees. A correct response has the right status code, the right headers, and the right body — in that order. Getting the order wrong silently breaks the response. Go's `ResponseWriter` enforces protocol correctness at the interface level, but you must understand the rules to avoid subtle bugs. Writing responses efficiently also matters: streaming a large JSON array instead of buffering it saves megabytes of memory under load.

## Mental model

Think of `ResponseWriter` as a three-step pipeline: first you set headers (via `w.Header().Set`), then you set the status code (via `w.WriteHeader`), then you write the body (via `w.Write` or `json.NewEncoder(w).Encode`). Once you enter the body phase, the header and status are locked. You cannot change them. The pipeline is one-way: headers → status → body. If you call `w.Write` before `w.WriteHeader`, Go automatically sends 200 OK as the status.

## Core idea

`http.ResponseWriter` is an interface with three methods:

```go
type ResponseWriter interface {
    Header() http.Header    // Get/set response headers (before WriteHeader)
    Write([]byte) (int, error)  // Write body; if WriteHeader not called, sends 200 OK
    WriteHeader(int)        // Set status code (must be called before Write)
}
```

Writing a JSON response follows this pattern:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(data)  // calls w.Write internally
```

For streaming, you can write to `w` incrementally. Go's server uses chunked transfer encoding automatically when the content length is unknown.

## Under the hood

`ResponseWriter` wraps a `bufio.Writer` attached to the TCP connection. When you call `w.Header().Set`, you modify a `http.Header` map. When you call `w.WriteHeader(code)`, the server writes the status line and all buffered headers to the wire. After that, `w.Write` writes directly to the TCP stream. If you use `json.NewEncoder(w)`, it writes JSON bytes into the same stream — Go uses chunked encoding if no `Content-Length` was set. You can set `Content-Length` explicitly via `w.Header().Set("Content-Length", strconv.Itoa(len(data)))` to avoid chunked encoding.

## How Go uses it

```go
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming not supported", http.StatusInternalServerError)
        return
    }
    for i := 0; i < 10; i++ {
        fmt.Fprintf(w, "event %d\n", i)
        flusher.Flush()
        time.Sleep(time.Second)
    }
}
```

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Event struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func singleEventHandler(w http.ResponseWriter, r *http.Request) {
	event := Event{
		ID:        1,
		Message:   "hello",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, event)
}

func eventsListHandler(w http.ResponseWriter, r *http.Request) {
	events := []Event{
		{ID: 1, Message: "start", Timestamp: time.Now().Format(time.RFC3339)},
		{ID: 2, Message: "processing", Timestamp: time.Now().Format(time.RFC3339)},
		{ID: 3, Message: "complete", Timestamp: time.Now().Format(time.RFC3339)},
	}
	writeJSON(w, http.StatusOK, events)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming not supported"})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	for i := 0; i < 5; i++ {
		fmt.Fprintf(w, "data: message %d at %s\n\n", i, time.Now().Format(time.RFC3339))
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /event", singleEventHandler)
	mux.HandleFunc("GET /events", eventsListHandler)
	mux.HandleFunc("GET /stream", streamHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client requests `GET /event`. `singleEventHandler` creates an `Event` struct.
2. `writeJSON` sets `Content-Type: application/json` on the header map.
3. `w.WriteHeader(http.StatusOK)` writes `HTTP/1.1 200 OK` plus headers to the wire.
4. `json.NewEncoder(w).Encode(event)` marshals the struct and writes the JSON bytes.
5. For `GET /stream`, the handler sets `Content-Type: text/event-stream` and writes a line, then calls `Flush()` to push it to the client immediately.
6. The loop continues, writing and flushing 5 events with 500ms delays.
7. After the handler returns, the connection is available for reuse.

## Common mistakes

- Setting `Content-Type` after `WriteHeader`. Headers must be set before the status code is sent. Setting them after is silently ignored.
- Calling `w.WriteHeader` twice. The first call sends the status; the second is logged but ignored. Remove duplicate calls.
- Using `json.Marshal` and then `w.Write` instead of `json.NewEncoder(w)`. Both work, but `NewEncoder` writes directly to the response without allocating an intermediate buffer.
- Forgetting to handle the `Content-Type` for binary data. Use `application/octet-stream` for arbitrary bytes and `Content-Disposition` for file downloads.
- Writing to the response after the handler returns. This is a no-op but indicates a logic error (e.g., a goroutine outliving the request).

## Debugging walkthrough

A handler returns an empty JSON response `{}` instead of the expected data:

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := db.QueryRow("SELECT ...").Scan(&user.Name)
    if err != nil {
        writeJSON(w, 500, map[string]string{"error": "db error"})
        return  // Missing return! Code continues.
    }
    writeJSON(w, 200, user)
}
```

**Symptom**: Always returns `{}` even with valid data.

**Root cause**: The `w.WriteHeader(500)` and `json.Encode` execute, but the function doesn't return — it continues to `writeJSON(w, 200, user)` at the bottom. The second `WriteHeader(200)` is ignored, but `json.NewEncoder(w).Encode(user)` writes after the error JSON, corrupting the output.

**Fix**: Add `return` after the error handling block.

## Production notes

- Always set `Content-Type` explicitly before `WriteHeader`. Never rely on `http.DetectContentType` — it guesses based on the first 512 bytes and can produce incorrect results for JSON.
- For large JSON arrays, use a streaming JSON encoder: write `[`, encode each item, write `,` between items, then write `]`. This avoids buffering the entire array in memory.
- Use `http.Server` with `WriteTimeout` to prevent slow clients from holding connections open while you stream a response.
- For HTTP/2 server push, use `http.Pusher`. Check `w.(http.Pusher)` and call `Push` before writing the response.

## Performance implications

- `json.NewEncoder(w)` writes directly to the TCP buffer — no intermediate allocation. For large payloads, this is significantly more memory-efficient than `json.Marshal`.
- Setting `Content-Length` header avoids chunked encoding, which adds ~50 bytes overhead per response. Pre-compute length for small-to-medium payloads.
- `http.Flusher.Flush()` forces a write to the TCP socket. Frequent flushes (every row in a stream) increase system calls. Batch where possible.
- Writing to `ResponseWriter` from multiple goroutines is not safe. Synchronize writes or use a single writer goroutine per request.

## Practice task

Write a handler `reportHandler` that returns a JSON report as a stream. The report contains a header object, then an array of 100 items (generated in a loop), then a footer object. Use `json.NewEncoder(w)` to write each part directly. Set `Content-Type` and return 200. Test with `httptest.NewRecorder` to verify the complete JSON output is valid.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/07-response-writing
go test ./curriculum/modules/08-http-rest-apis/lessons/07-response-writing
```

## Review questions

1. What happens if you call `w.WriteHeader(404)` after `w.Write([]byte("hello"))`?
2. What is the difference between `json.NewEncoder(w).Encode(v)` and `json.Marshal(v)` followed by `w.Write(b)`?
3. How does Go decide whether to use chunked transfer encoding?
4. Why is it important to call `return` after `http.Error` in a handler?
5. What interface must `ResponseWriter` implement to support streaming (Server-Sent Events)?

## NEXT UP

API error shape — designing a consistent, structured JSON error format that clients can parse programmatically.
