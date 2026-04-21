// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bytes"
	"fmt"
	"sync"
)

// ============================================================================
// Stage 07: Concurrency Patterns - sync.Pool
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
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
// RUN: go run ./07-concurrency/02-concurrency-patterns/3-sync-pool
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
	fmt.Println("NEXT UP: CP.4 bounded pipeline exercise")
	fmt.Println("   Current: CP.3 (sync.Pool)")
	fmt.Println("---------------------------------------------------")
}
