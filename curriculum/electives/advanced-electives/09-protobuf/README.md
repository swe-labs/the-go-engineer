# Protobuf

## When you need this

JSON and other text-based serialization formats are slow to parse, produce large payloads, and lack schema enforcement across services. Protobuf solves this by providing a compact binary format with a schema definition language, code generation for multiple languages, and backward-compatible field evolution.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-08

## Mental model

Protobuf is a contract-first serialization format. You define the message structure in a .proto file (the contract), and the protoc compiler generates Go structs with marshal/unmarshal methods. Producers and consumers share the contract — they do not share code or runtime dependencies. Fields are numbered (field 1, field 2) instead of named — this enables backward-compatible field addition and removal.

## How it works

Protobuf encoding uses a binary wire format with varint encoding for integers, length-delimited strings/bytes, and fixed-width encoding for floats/doubles. Each field is identified by its field number and wire type (varint, 64-bit, length-delimited, start-group, 32-bit). The Go protoc plugin (protoc-gen-go) generates Go structs with protobuf struct tags and Marshal/Unmarshal methods that use the google.golang.org/protobuf/proto package for encoding and decoding.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Protobuf")
}
```

## Tradeoffs

- Protobuf encoding is faster and produces smaller output than JSON (no field names in output, binary format)
- Generated marshal/unmarshal code uses zero-reflection — field assignments are direct struct operations
- Protobuf decoding allocates — each message and sub-message is heap-allocated; pooled buffers can reduce GC pressure
- The reflection-based proto adapters (proto.Marshal/Unmarshal for dynamic types) are ~10x slower than generated code paths

## Common mistakes

- Reusing field numbers after removing a field — the old field number is reserved and should not be reused (use reserved keyword in .proto).
- Changing field types after deployment — a int32 field cannot be changed to string; use a new field number and deprecate the old one.
- Not using proto3 syntax — proto2 supports required fields and custom defaults, which complicate backward compatibility.
- Storing protobuf binary data in a text-based format (JSON database) — protobuf binary is opaque and cannot be searched or indexed.
- Importing .proto files with relative paths that break when the project structure changes — use protoc -I with absolute or module-relative import paths.

## Practice task

The learner must write a .proto file with at least 2 message types and 1 service definition, generate Go code using protoc and protoc-gen-go, marshal and unmarshal a protobuf message in Go, and demonstrate backward compatibility by adding a new field to the proto file without breaking existing messages.

## Production notes

Protobuf is used in every gRPC service and many event-driven systems. Google uses protobuf for all internal RPC. Go microservices use protobuf for service-to-service communication, event schemas in message queues (Kafka, NATS), and configuration file formats. Protobuf is the standard schema format for Kubernetes API resources.

## Related core lessons

- Module 16.08 — Mockery
- Module 17.07 — go generate
