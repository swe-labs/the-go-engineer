// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bytes"
	"fmt"
	"sync"
)

// ============================================================================
// Section 24: errgroup & sync.Pool — sync.Pool
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - sync.Pool: reuse temporary objects to reduce GC pressure
//   - The correct Get() → use → reset → Put() lifecycle
//   - Why you MUST reset objects before Put()
//   - Building a production-grade byte buffer pool (used by fmt, json, http)
//   - How to benchmark pool impact with testing.B
//
// ENGINEERING DEPTH:
//   Go's garbage collector runs when allocated heap exceeds a threshold.
//   A service processing 10,000 requests/sec may allocate 10,000 * 4096 bytes
//   = 40MB per second just for temporary buffers. When GC runs to reclaim this
//   memory, every goroutine in the program is paused for microseconds (STW:
//   Stop The World). At 99th percentile this becomes your latency spike.
//
//   sync.Pool solves this by keeping a set of pooled objects. Between GC cycles,
//   Get() returns a recycled object and Put() returns it. During GC, the pool
//   is cleared (the GC intentionally evicts pools to prevent memory leaks from
//   objects that should have been freed). Objects MUST therefore be treated as
//   temporary — never assume an object from the pool is clean or that a Put()
//   object will be returned by the next Get().
//
//   PRODUCTION RULE: Reset the object (buf.Reset(), clear the struct) before
//   Put(). Otherwise the next caller gets stale data — a serious security bug
//   if the buffer contains HTTP headers or auth tokens.
//
// RUN: go run ./24-errgroup-and-pools/3-sync-pool
// ============================================================================

// ============================================================================
// ByteBufferPool — a production-grade pooled buffer
// ============================================================================
// This is the same pattern used by encoding/json, net/http, and fmt internally.

var bufPool = sync.Pool{
	// New is called when the pool is empty. It creates a fresh object.
	// New should always allocate — don't return nil here.
	New: func() any {
		// Pre-allocate 4KB: covers most HTTP responses without reallocation.
		buf := bytes.NewBuffer(make([]byte, 0, 4096))
		return buf
	},
}

// GetBuffer returns a clean buffer from the pool.
func GetBuffer() *bytes.Buffer {
	buf := bufPool.Get().(*bytes.Buffer) // Type assertion: pool stores any
	buf.Reset()                          // CRITICAL: clear previous contents
	return buf
}

// PutBuffer returns a buffer to the pool after use.
func PutBuffer(buf *bytes.Buffer) {
	// Optional: don't pool oversized buffers — they waste pool memory.
	// A 10MB buffer from a single huge response shouldn't block a 4KB slot.
	if buf.Cap() > 64*1024 {
		return // Let GC collect it
	}
	buf.Reset() // Belt-and-suspenders reset before returning
	bufPool.Put(buf)
}

// buildHTTPResponse demonstrates the correct pool usage lifecycle.
// Compare withPool vs withoutPool in the benchmark below.
func buildHTTPResponseWithPool(status int, body string) string {
	buf := GetBuffer()   // 1. Get (possibly recycled)
	defer PutBuffer(buf) // 4. Return to pool when done (ALWAYS defer this)

	// 2. Use the buffer
	fmt.Fprintf(buf, "HTTP/1.1 %d OK\r\n", status)
	fmt.Fprintf(buf, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(buf, "\r\n%s", body)

	return buf.String() // 3. Extract result BEFORE defer fires
}

func buildHTTPResponseWithoutPool(status int, body string) string {
	var buf bytes.Buffer // New allocation every call
	fmt.Fprintf(&buf, "HTTP/1.1 %d OK\r\n", status)
	fmt.Fprintf(&buf, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(&buf, "\r\n%s", body)
	return buf.String()
}

// ============================================================================
// StructPool — pooling custom structs
// ============================================================================
// Pooling structs is more error-prone than buffers because you MUST zero out
// every field. Missing a field means previous request data bleeds into the next.

type RequestContext struct {
	UserID    string
	RequestID string
	Tags      []string
	Headers   map[string]string
}

func (r *RequestContext) Reset() {
	r.UserID = ""
	r.RequestID = ""
	r.Tags = r.Tags[:0]        // Reuse the backing array, reset length to 0
	for k := range r.Headers { // Clear the map (reuse allocation)
		delete(r.Headers, k)
	}
}

var ctxPool = sync.Pool{
	New: func() any {
		return &RequestContext{
			Tags:    make([]string, 0, 8),
			Headers: make(map[string]string, 16),
		}
	},
}

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

	// Simulate processing
	return fmt.Sprintf("processed request %s for user %s", rc.RequestID, rc.UserID)
}

func main() {
	// Demonstrate correct lifecycle
	resp := buildHTTPResponseWithPool(200, `{"ok":true}`)
	fmt.Println("Response via pool:")
	fmt.Println(resp)

	// Demonstrate struct pool
	fmt.Println("\nProcessed:", processRequest("usr_42", "req_001"))
	fmt.Println("Processed:", processRequest("usr_99", "req_002")) // Reuses same struct

	// KEY TAKEAWAY:
	// - sync.Pool: reuse temporary objects to cut GC pressure
	// - Get() → use → Reset() → Put() is the complete lifecycle
	// - ALWAYS reset before Put() — stale data is a security bug
	// - Don't pool oversized objects — they waste pool slots
	// - Use -benchmem to measure before and after: allocs/op should drop to 0
	// - Pool objects are evicted on GC — never rely on pool for caching
}
