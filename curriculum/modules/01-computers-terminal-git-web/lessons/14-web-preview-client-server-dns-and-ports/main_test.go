package main

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTCPSocketLifecycle(t *testing.T) {
	// 1. Bind to ephemeral port (port 0)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen failed: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split host/port: %v", err)
	}

	portNum, err := strconv.Atoi(portStr)
	if err != nil || portNum <= 0 {
		t.Fatalf("invalid ephemeral port allocated: %s", portStr)
	}

	// 2. Accept and echo in background goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf := make([]byte, 128)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		_, _ = conn.Write(bytes.ToUpper(buf[:n]))
	}()

	// 3. Dial from client
	clientConn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		t.Fatalf("net.DialTimeout failed: %v", err)
	}
	defer clientConn.Close()

	testMsg := []byte("ping\n")
	_, err = clientConn.Write(testMsg)
	if err != nil {
		t.Fatalf("failed writing to TCP connection: %v", err)
	}

	reply := make([]byte, 128)
	n, err := clientConn.Read(reply)
	if err != nil {
		t.Fatalf("failed reading from TCP connection: %v", err)
	}

	if !strings.Contains(string(reply[:n]), "PING") {
		t.Errorf("expected echo 'PING', got %q", string(reply[:n]))
	}

	<-done
}

func TestLocalhostResolution(t *testing.T) {
	ips, err := net.LookupHost("localhost")
	if err != nil {
		t.Logf("DNS lookup for localhost skipped/offline: %v", err)
		return
	}
	if len(ips) == 0 {
		t.Errorf("expected at least one IP for localhost")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

