package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("=== 1. DNS Domain Name Resolution ===")
	targetDomain := "example.com"
	ips, err := net.LookupHost(targetDomain)
	if err != nil {
		fmt.Printf("DNS resolution for %s failed (%v); using 127.0.0.1 for local demo\n", targetDomain, err)
		ips = []string{"127.0.0.1"}
	} else {
		fmt.Printf("Successfully resolved DNS for %q to IP addresses:\n", targetDomain)
		for i, ip := range ips {
			fmt.Printf("  [%d] IP: %s\n", i+1, ip)
		}
	}

	fmt.Println("\n=== 2. TCP Server Socket Binding & Ephemeral Ports ===")
	// Port "0" requests the OS kernel to allocate an available dynamic ephemeral port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("Failed to bind listener: %v\n", err)
		return
	}
	defer listener.Close()

	serverAddress := listener.Addr().String()
	fmt.Printf("Server actively listening on: %s\n", serverAddress)

	// Concurrency preview: accept connection in background
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Printf("[Server Engine] Accepted incoming TCP connection from Remote Client: %s\n", conn.RemoteAddr())
	}()

	fmt.Println("\n=== 3. Client TCP Handshake & Connection Inspection ===")
	clientConn, err := net.DialTimeout("tcp", serverAddress, 2*time.Second)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer clientConn.Close()

	fmt.Printf("[Client Engine] Connected successfully!\n")
	fmt.Printf("  Local Client Socket (Ephemeral Port): %s\n", clientConn.LocalAddr())
	fmt.Printf("  Remote Target Server Socket (Host:Port): %s\n", clientConn.RemoteAddr())

	<-serverDone
	fmt.Println("\nSummary: DNS resolved host -> TCP established connection -> Ports routed traffic.")
}
