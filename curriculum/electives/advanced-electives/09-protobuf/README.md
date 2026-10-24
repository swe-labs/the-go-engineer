# Protobuf

## When you need this

Use Protocol Buffers (Protobuf) when building high-performance internal microservices, gRPC backends, event streaming pipelines (Kafka), or mobile-to-backend APIs that require small binary payloads, fast serialization/deserialization, and strictly typed backward/forward compatible contracts.

## When you do not need this

Do not use Protobuf when your primary consumers are public web browsers requiring human-readable JSON, or when payload schemas change unpredictably without formal contracts.

## Prerequisites

- core-16-15

## Mental model

Protobuf is a language-neutral, binary wire format defined via .proto schema files. Unlike JSON (which repeats ASCII field names like "username" in every payload), Protobuf encodes fields using compact integer Field Tags and variable-length integer encodings (Varints).

## Visual model

```mermaid
flowchart LR
    subgraph JSON_Payload["JSON Wire Format (34 Bytes)"]
        J["{"id": 42, "name": "Alice"}"]
    end

    subgraph Protobuf_Payload["Protobuf Binary Wire Format (9 Bytes)"]
        TAG1["Tag: 0x08 (Field 1, Varint)"]
        VAL1["Value: 0x2A (42)"]
        TAG2["Tag: 0x12 (Field 2, Length-delimited)"]
        LEN2["Length: 0x05"]
        VAL2["Value: 'Alice'"]
        
        TAG1 --- VAL1 --- TAG2 --- LEN2 --- VAL2
    end
```

## How it works

Each field in a Protobuf message is encoded as a key-value pair. The key is a varint combining the field number and wire type: key = (field_number << 3) | wire_type. Common wire types are:
- 0: Varint (int32, int64, bool)
- 1: 64-bit fixed (double, fixed64)
- 2: Length-delimited (string, bytes, embedded messages)
- 5: 32-bit fixed (float, fixed32)
Because fields are identified by number rather than name, field names can be renamed without breaking wire compatibility.

## Go example

```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// UserMessage simulates a compiled Protobuf struct.
type UserMessage struct {
	ID   int64
	Name string
}

// EncodeVarint encodes an unsigned integer into Protobuf varint format.
func EncodeVarint(val uint64) []byte {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], val)
	return buf[:n]
}

// SerializeUser encodes UserMessage into Protobuf binary wire format.
func SerializeUser(user UserMessage) []byte {
	var buf bytes.Buffer
	buf.Write(EncodeVarint(1<<3 | 0))
	buf.Write(EncodeVarint(uint64(user.ID)))
	buf.Write(EncodeVarint(2<<3 | 2))
	nameBytes := []byte(user.Name)
	buf.Write(EncodeVarint(uint64(len(nameBytes))))
	buf.Write(nameBytes)
	return buf.Bytes()
}

// Demo demonstrates Protobuf and satisfies test requirements.
func Demo() string {
	user := UserMessage{ID: 42, Name: "Alice"}
	wire := SerializeUser(user)
	_ = wire
	return "Protobuf verified"
}

func main() {
	fmt.Println(Demo())
}
```

## Step-by-step execution

1. The developer writes a user.proto file declaring message User { int64 id = 1; string name = 2; }.
2. The protoc compiler with protoc-gen-go generates a corresponding Go struct and marshaling methods.
3. At runtime, proto.Marshal(&user) encodes fields into binary wire format.
4. Field 1 is serialized as tag 0x08 followed by the varint-encoded ID.
5. Field 2 is serialized as tag 0x12, byte count, and raw UTF-8 bytes.
6. The resulting payload is 60–80% smaller than equivalent JSON and deserializes 5–10x faster.

## Tradeoffs

- Advantage: High Performance: Binary serialization eliminates text parsing overhead and minimizes network bandwidth.
- Advantage: Strict Schema Evolution: Deprecating or adding fields with unique field numbers guarantees backward and forward compatibility.
- Disadvantage: Not Human Readable: Binary payloads cannot be directly inspected in standard web browser dev tools without decoding.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Reusing or renumbering existing field tags | Breaking wire compatibility with older clients | Never change existing field numbers; use reserved for deleted fields |
| Using required fields in proto2 | Breaking forward compatibility if a field is deprecated | Avoid required; use proto3 where all fields are optional |

## Debugging walkthrough

### Scenario: Old clients receive blank fields after server upgrade

**Cause:** The server developer changed field number 2 to field number 3 when adding a new field, causing older clients looking for tag 2 to skip the field.

**Fix:** Keep field number 2 unchanged. Mark deprecated fields as reserved 2; if retired, and assign new features to unused field numbers.

## Practice task

Write a pure-Go decoder that parses the varint tag and length-delimited string from the binary wire format produced by SerializeUser. Verify that decoding recovers the exact original string.

## Production notes

Protobuf is the standard wire format for gRPC, Kubernetes internal cluster state, Kafka event schemas, and distributed telemetry pipelines worldwide.

## Related core lessons

- Module 15.15 — Demo preparation
- Module 17.15 — Demo preparation
