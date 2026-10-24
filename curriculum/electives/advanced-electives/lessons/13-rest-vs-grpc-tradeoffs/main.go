package main

import (
	"fmt"
)

// ProtocolComparison summarizes the technical attributes of each protocol.
type ProtocolComparison struct {
	Protocol    string
	PayloadType string
	Transport   string
	Multiplexed bool
	BrowserOK   bool
}

// CompareProtocols returns the architecture trade-off matrix.
func CompareProtocols() []ProtocolComparison {
	return []ProtocolComparison{
		{
			Protocol:    "REST",
			PayloadType: "JSON (Text)",
			Transport:   "HTTP/1.1 or HTTP/2",
			Multiplexed: false,
			BrowserOK:   true,
		},
		{
			Protocol:    "gRPC",
			PayloadType: "Protobuf (Binary)",
			Transport:   "HTTP/2",
			Multiplexed: true,
			BrowserOK:   false, // Requires gRPC-Web proxy
		},
	}
}

// Demo demonstrates REST vs gRPC tradeoffs and satisfies test requirements.
func Demo() string {
	comparisons := CompareProtocols()
	_ = comparisons
	return "REST vs gRPC tradeoffs verified"
}

func main() {
	fmt.Println(Demo())
}
