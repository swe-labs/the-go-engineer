package main

import (
	"fmt"
	"net"
	"os"
)

func resolveDomain(domain string) ([]string, error) {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		return nil, fmt.Errorf("lookup failed: %w", err)
	}
	return addrs, nil
}

func dialTCP(address string) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("resolve failed: %w", err)
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	return conn, nil
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

	ips, err := resolveDomain(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving domain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Resolved %s to %d IP address(es):\n", domain, len(ips))
	for _, ip := range ips {
		fmt.Printf("  %s\n", ip)
	}
	fmt.Println()

	address := net.JoinHostPort(ips[0], port)
	fmt.Printf("Dialing TCP connection to %s...\n", address)

	conn, err := dialTCP(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

	fmt.Printf("Connection established!\n")
	fmt.Printf("  Local:  %s:%d\n", localAddr.IP, localAddr.Port)
	fmt.Printf("  Remote: %s:%d\n", remoteAddr.IP, remoteAddr.Port)
	fmt.Println()

	request := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", domain)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Sent HTTP request to %s:%s\n", domain, port)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Received %d bytes in response\n", n)
	fmt.Printf("First %d bytes:\n%s\n", n, string(buf[:n]))
}
