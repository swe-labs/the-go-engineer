// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: sync.Pool
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - sync.Pool fundamentals and practical application in Go.
//
// WHY THIS MATTERS:
//   - sync.Pool provides a structured approach to writing clean Go code.
//
// RUN:
//   go run ./07-concurrency/02-concurrency-patterns/3-sync-pool
//
// KEY TAKEAWAY:
//   - sync.Pool fundamentals and practical application in Go.
// ============================================================================

package main

import (
	"bytes"
	"fmt"
	"sync"
)

// Stage 07: Concurrency Patterns - sync.Pool
//
//   - sync.Pool: reuse temporary objects to reduce GC pressure
//   - The correct Get -> use -> Reset -> Put lifecycle
//   - Why you must reset objects before Put
//   - Building a byte-buffer pool similar to real standard-library patterns
//   - How to benchmark pool impact with testing.B
//
// ENGINEERING DEPTH:
//   Go's garbage collector runs when allocated heap exceeds a threshold.
//   A service processing 10,000 requests per second may allocate tens of MB per
//   second just for temporary buffers. When GC runs to reclaim this memory, you
//   pay extra latency. sync.Pool reduces that churn by recycling short-lived
//   objects between GC cycles.
//
//   Pools are intentionally cleared during GC, so you must treat pooled objects
//   as temporary. Never assume the next Get returns your previous Put, and
//   always reset the object before handing it back.
//

// bufPool (Pool): reuses response buffers so allocation pressure stays visible in the lesson.
var bufPool = sync.Pool{
	New: func() any {
		buf := bytes.NewBuffer(make([]byte, 0, 4096))
		return buf
	},
}

// GetBuffer (Function): runs the get buffer step and keeps its inputs, outputs, or errors visible.
func GetBuffer() *bytes.Buffer {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer (Function): runs the put buffer step and keeps its inputs, outputs, or errors visible.
func PutBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 64*1024 {
		return
	}
	buf.Reset()
	bufPool.Put(buf)
}

// buildHTTPResponseWithPool (Function): runs the build http response with pool step and keeps its inputs, outputs, or errors visible.
func buildHTTPResponseWithPool(status int, body string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	fmt.Fprintf(buf, "HTTP/1.1 %d OK\r\n", status)
	fmt.Fprintf(buf, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(buf, "\r\n%s", body)

	return buf.String()
}

// buildHTTPResponseWithoutPool (Function): runs the build http response without pool step and keeps its inputs, outputs, or errors visible.
func buildHTTPResponseWithoutPool(status int, body string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HTTP/1.1 %d OK\r\n", status)
	fmt.Fprintf(&buf, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(&buf, "\r\n%s", body)
	return buf.String()
}

// RequestContext (Struct): groups the state used by the request context example boundary.
type RequestContext struct {
	UserID    string
	RequestID string
	Tags      []string
	Headers   map[string]string
}

// RequestContext.Reset (Method): applies the reset operation to receiver state at a visible boundary.
func (r *RequestContext) Reset() {
	r.UserID = ""
	r.RequestID = ""
	r.Tags = r.Tags[:0]
	for k := range r.Headers {
		delete(r.Headers, k)
	}
}

// ctxPool (Pool): reuses request contexts while preserving an explicit reset boundary.
var ctxPool = sync.Pool{
	New: func() any {
		return &RequestContext{
			Tags:    make([]string, 0, 8),
			Headers: make(map[string]string, 16),
		}
	},
}

// processRequest (Function): runs the process request step and keeps its inputs, outputs, or errors visible.
func processRequest(userID, requestID string) string {
	rc := ctxPool.Get().(*RequestContext)
	defer func() {
		rc.Reset()
		ctxPool.Put(rc)
	}()

	rc.UserID = userID
	rc.RequestID = requestID
	rc.Tags = append(rc.Tags, "api", "v2")
	rc.Headers["Authorization"] = "Bearer xyz"

	return fmt.Sprintf("processed request %s for user %s", rc.RequestID, rc.UserID)
}

func main() {
	resp := buildHTTPResponseWithPool(200, `{"ok":true}`)
	fmt.Println("Response via pool:")
	fmt.Println(resp)

	fmt.Println("\nProcessed:", processRequest("usr_42", "req_001"))
	fmt.Println("Processed:", processRequest("usr_99", "req_002"))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CP.4 -> 07-concurrency/02-concurrency-patterns/4-bounded-pipeline-exercise")
	fmt.Println("   Current: CP.3 (sync.Pool)")
	fmt.Println("---------------------------------------------------")
}
