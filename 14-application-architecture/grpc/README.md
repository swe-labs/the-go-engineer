# Section 26: gRPC

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Protocol Buffers | Beginner | **Critical** | Schema-first API definition |
| Unary RPC | Beginner | **Critical** | Request/response like HTTP |
| Server streaming | Intermediate | High | Real-time data push |
| Client streaming | Advanced | High | Batch upload, chunked writes |
| Bidirectional streaming | Advanced | High | Chat, live collaboration |
| Interceptors | Advanced | **Critical** | Middleware for gRPC |

## Engineering Depth

gRPC is Go's native RPC framework, built by Google. It uses Protocol Buffers (proto3) as its IDL — a schema language that generates type-safe client and server code in any language. A Go client can seamlessly call a Python server, a Rust server, or a Java server with zero manual serialization code.

**Why gRPC over REST:**
- **2-10x faster**: Binary framing (HTTP/2) + protobuf encoding vs text JSON
- **Streaming**: First-class bidirectional streaming over a single TCP connection
- **Type safety**: The proto schema IS the contract — no drift between client and server
- **Code generation**: Client stubs generated automatically, no manual HTTP client

**Proto3 → Go generation workflow:**
```bash
# Install tools (one-time setup)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
brew install protobuf  # or apt install protobuf-compiler

# Generate Go code from .proto
protoc --go_out=. --go-grpc_out=. proto/service.proto
```

## Contents

| Directory | Topic | Level |
|-----------|-------|-------|
| `proto/` | Service definition files | Beginner |
| `gen/` | Generated Go code (committed so you can run without protoc) | — |
| `1-unary/` | Simple request/response RPC | Beginner |
| `2-streaming/` | Server-side streaming and interceptors | Intermediate |

## How to Run

```bash
# Terminal 1: start the server
go run ./26-grpc/1-unary/server

# Terminal 2: run the client
go run ./26-grpc/1-unary/client

# Terminal 1: start the streaming server
go run ./26-grpc/2-streaming/server

# Terminal 2: run the streaming client
go run ./26-grpc/2-streaming/client
```

## Adding to go.mod

```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

## References

- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [Protocol Buffers Language Guide](https://protobuf.dev/programming-guides/proto3/)
- [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc)
