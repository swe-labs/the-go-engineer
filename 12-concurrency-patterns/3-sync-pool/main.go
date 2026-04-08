// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bytes"
	"fmt"
	"sync"
)

// ============================================================================
// Section 12: Concurrency Patterns � sync.Pool
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - sync.Pool: reuse temporary objects to reduce GC pressure
//   - The correct Get() ? use ? reset ? Put() lifecycle
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
//   temporary � never assume an object from the pool is clean or that a Put()
//   object will be returned by the next Get().
//
//   PRODUCTION RULE: Reset the object (buf.Reset(), clear the struct) before
//   Put(). Otherwise the next caller gets stale data � a serious security bug
//   if the buffer contains HTTP headers or auth tokens.
//
// RUN: go run ./12-concurrency-patterns/3-sync-pool
// ============================================================================

var bufPool = sync.Pool{
	New: func() any {
		buf := bytes.NewBuffer(make([]byte, 0, 4096))
		return buf
	},
}

func GetBuffer() *bytes.Buffer {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func PutBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 64*1024 {
		return
	}
	buf.Reset()
	bufPool.Put(buf)
}

func buildHTTPResponseWithPool(status int, body string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	fmt.Fprintf(buf, "HTTP/1.1 %d OK\r\n", status)
	fmt.Fprintf(buf, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(buf, "\r\n%s", body)

	return buf.String()
}

func buildHTTPResponseWithoutPool(status int, body string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HTTP/1.1 %d OK\r\n", status)
	fmt.Fprintf(&buf, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(&buf, "\r\n%s", body)
	return buf.String()
}

type RequestContext struct {
	UserID    string
	RequestID string
	Tags      []string
	Headers   map[string]string
}

func (r *RequestContext) Reset() {
	r.UserID = ""
	r.RequestID = ""
	r.Tags = r.Tags[:0]
	for k := range r.Headers {
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

	return fmt.Sprintf("processed request %s for user %s", rc.RequestID, rc.UserID)
}

func main() {
	resp := buildHTTPResponseWithPool(200, `{"ok":true}`)
	fmt.Println("Response via pool:")
	fmt.Println(resp)

	fmt.Println("\nProcessed:", processRequest("usr_42", "req_001"))
	fmt.Println("Processed:", processRequest("usr_99", "req_002"))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("?? NEXT UP: TE.1 unit testing")
	fmt.Println("   Current: CP.3 (sync.Pool)")
	fmt.Println("---------------------------------------------------")
}
