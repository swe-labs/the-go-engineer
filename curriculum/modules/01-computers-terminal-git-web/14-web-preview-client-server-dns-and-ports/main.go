package main

import (
	"fmt"
	"os"
)

func dnsLookup(hostname string) string {
	if hostname == "example.com" {
		return "93.184.216.34"
	}
	if hostname == "localhost" {
		return "127.0.0.1"
	}
	return ""
}

func simulateConnection(address string) string {
	colonPos := 0
	for i := 0; i < len(address); i++ {
		if address[i] == ':' {
			colonPos = i
			break
		}
	}
	host := address[:colonPos]
	port := address[colonPos+1:]
	return fmt.Sprintf("Simulated connection to %s on port %s (host: %s)", address, port, host)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <domain> [port]")
		fmt.Println("Example: go run . example.com 80")
		os.Exit(1)
	}

	domain := os.Args[1]
	port := "80"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	fmt.Printf("=== Web Preview: %s ===\n", domain)
	fmt.Println()

	ip := dnsLookup(domain)
	if ip == "" {
		fmt.Printf("Could not resolve %s\n", domain)
		os.Exit(1)
	}
	fmt.Printf("Resolved %s to IP: %s\n", domain, ip)
	fmt.Println()

	address := ip + ":" + port
	fmt.Println(simulateConnection(address))
	fmt.Println()
	fmt.Println("DNS translates domain names to IP addresses.")
	fmt.Println("Ports identify specific services on a server.")
	fmt.Printf("Port %s is the standard HTTP port.\n", port)
}
