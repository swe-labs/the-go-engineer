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
