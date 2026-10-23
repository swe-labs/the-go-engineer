# gRPC streaming

## When you need this

Unary RPCs (single request, single response) are insufficient for real-time data feeds, large result sets, or long-lived connections. gRPC streaming solves this with three streaming patterns: server-streaming (one request, many responses), client-streaming (many requests, one response), and bidirectional streaming (many requests, many responses) — all over a single HTTP/2 connection.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-10

## Mental model

A gRPC stream is a series of messages exchanged over a single HTTP/2 stream. Server-streaming: the client sends one request and receives multiple responses sent by the server. Client-streaming: the client sends multiple requests and the server sends one response. Bidirectional: both sides send multiple messages independently — the stream is full-duplex.

## How it works

gRPC streams use HTTP/2 frames. The client opens a new HTTP/2 stream for each RPC. For server-streaming, the server sends multiple DATA frames (each containing a serialized protobuf message) followed by a HEADERS frame with the response status. The client reads messages from the stream with Recv() until io.EOF. For client-streaming, the client sends multiple DATA frames and the server sends a single response. Bidirectional streams interleave DATA frames from both sides.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: gRPC streaming")
}
```

## Tradeoffs

- Server streaming reduces per-message overhead compared to unary: one HTTP/2 stream handles all messages instead of one stream per message
- Client-side buffering: client must buffer received messages or process them incrementally — unbuffered receive+process limits throughput to one message per round-trip time
- gRPC's default max message size (4MB) applies per-message in a stream — large messages in a stream trigger the same limit as unary
- Stream-level flow control prevents OOM but adds latency under load — each window update adds a round-trip before more data flows

## Common mistakes

- Not handling io.EOF on server-streaming Recv() — io.EOF signals the end of the stream and is not an error; forgetting to check for io.EOF causes the client to treat stream end as a failure.
- Blocking Send() on a slow client — if the client is not reading fast enough, the server's Send() blocks, potentially backing up to the entire gRPC server. Use a buffer or flow control.
- Calling Send() and Recv() from the same goroutine in bidirectional streaming — the methods block, causing a deadlock. Use separate goroutines for send and receive.
- Not setting deadlines on streaming RPCs — a streaming RPC that hangs indefinitely holds resources on both client and server. Set per-RPC or per-stream deadlines.
- Using too many concurrent streams — each stream consumes goroutines and memory on both ends. Limit concurrent streams with gRPC's MaxConcurrentStreams server option.

## Practice task

The learner must implement a server-streaming RPC that sends 10 messages at 1-second intervals, a client that receives them with Recv() until io.EOF, and a bidirectional streaming RPC where the client sends 5 messages and the server echoes each one back.

## Production notes

gRPC streaming is used for real-time features: log streaming (server-streaming), file uploads with progress (client-streaming), chat and collaboration (bidirectional), event feeds (server-streaming), and large dataset pagination (server-streaming). Go services use server-streaming for API endpoints that return paginated results without page boundaries.

## Related core lessons

- Module 16.10 — gRPC fundamentals
- Module 17.09 — Protobuf
