# Web preview: client, server, DNS, and ports

## Learning objective

Trace the client-server architecture of the Web, resolving domain names via DNS and connecting over TCP ports.

## Why this matters

Every web application and backend API depends on four foundational network layers: finding the server's IP address (DNS), establishing a reliable connection (TCP), delivering traffic to the correct application (Ports), and exchanging formatted messages (HTTP). Without understanding these layers, network timeouts, port conflicts, DNS propagation delays, and firewall issues feel like inexplicable magic.

> [!NOTE]
> **Zero-Magic Reassurance:**
> In the Go example below, you will see server plumbing like `go func()`, `net.Listen()`, and `chan struct{}`. These concurrency and networking primitives will be taught step-by-step in Modules 06 and 14. For now, focus on the mechanical network concepts: how DNS translates human names to IP addresses, how ports demultiplex services, and how the TCP 3-way handshake establishes a reliable socket stream.

## Mental model

The web operates on a client-server architecture layered like a postal delivery system:
1. **The Client (The Sender):** A browser, terminal `curl` command, or microservice initiating a request.
2. **DNS (The Phonebook):** Translates human-readable domain names (`example.com`) into numerical machine IP addresses (`93.184.216.34`).
3. **TCP/IP (The Delivery Highway):** Breaks data into packets, routes them across internet routers (IP), and guarantees reliable, ordered, error-checked delivery (TCP).
4. **Ports (The Apartment Numbers):** A single server host with one IP address can run many different services (web server, database, SSH). A 16-bit port number (0–65535) tells the OS kernel which specific process should receive the incoming network packets.

## Visual model

```mermaid
sequenceDiagram
    autonumber
    actor Client as Browser / HTTP Client
    participant DNS as DNS Resolver (Port 53)
    participant Server as Web Server Host (IP: 93.184.216.34)

    Note over Client, DNS: 1. Name Resolution Layer (DNS over UDP 53)
    Client->>DNS: Query: A record for "example.com"
    DNS-->>Client: Resolved IP: 93.184.216.34

    Note over Client, Server: 2. Transport Layer: TCP 3-Way Handshake
    Client->>Server: SYN (Target Port 443; Client Seq: 100)
    Server-->>Client: SYN-ACK (Server Seq: 500; Ack: 101)
    Client->>Server: ACK (Connection Established; Ack: 501)

    Note over Client, Server: 3. Application Layer Protocol (HTTP / TLS Stream)
    Client->>Server: GET / HTTP/1.1 (Host: example.com)
    Server-->>Client: HTTP/1.1 200 OK (HTML Payload)
```

## Core idea

The internet does not route traffic by domain names; it routes packets exclusively by numerical IP addresses. DNS acts as the universal distributed translation directory. Once the IP address is known, the client initiates a TCP handshake to a designated port number, opening a bidirectional byte stream between client and server processes.

## Under the hood

### 1. The DNS Hierarchy
DNS is a globally distributed, hierarchical database:
- **Root Nameservers (`.`):** 13 root server clusters worldwide direct queries to Top-Level Domain (TLD) servers.
- **TLD Nameservers:** Manage top-level domains like `.com`, `.org`, `.net`, `.io`.
- **Authoritative Nameservers:** The designated servers that hold the official DNS records for a specific domain (e.g., `ns1.google.com`).
- **Common Record Types:**
  - `A`: Maps a domain name to an IPv4 address (e.g., `93.184.216.34`).
  - `AAAA`: Maps a domain name to a 128-bit IPv6 address.
  - `CNAME`: Canonical Name (an alias pointing to another domain name).
  - `MX`: Mail Exchange records for routing email.

### 2. The TCP 3-Way Handshake
Before any data can travel between client and server, TCP establishes a reliable connection using three synchronized packets:
1. **`SYN` (Synchronize):** Client chooses a random Initial Sequence Number (ISN) and sends a packet with the SYN flag set to server port 80/443.
2. **`SYN-ACK` (Synchronize-Acknowledge):** Server acknowledges client's sequence number and sends its own random sequence number.
3. **`ACK` (Acknowledge):** Client acknowledges server's sequence number. The TCP socket is now established and data transmission begins.

### 3. Port Number Architecture
A port is an unsigned 16-bit integer, spanning `0` to `65535`:
- **Well-Known Ports (`0`–`1023`):** Reserved for standard privileged services. On Unix systems, processes must have root privileges to bind to them:
  - Port `80`: HTTP (Unencrypted Web)
  - Port `443`: HTTPS (TLS Encrypted Web)
  - Port `22`: SSH (Secure Shell)
  - Port `53`: DNS (Domain Name Service)
- **Registered Ports (`1024`–`49151`):** Typically used for application services and databases:
  - Port `3306`: MySQL
  - Port `5432`: PostgreSQL
  - Port `6379`: Redis
  - Port `8080`: Common alternative HTTP development port
- **Dynamic / Ephemeral Ports (`49152`–`65535`):** Automatically allocated by the client OS kernel for the client's end of the connection, reclaimed when the connection closes.

## How Go uses it

Go's standard library `net` package provides production-grade networking:
- `net.LookupHost(domain)` queries system DNS resolvers.
- `net.Dial("tcp", "example.com:80")` resolves DNS and executes the TCP 3-way handshake in a single call.
- `net.Listen("tcp", ":8080")` binds a TCP socket to port 8080 and listens for incoming connections.
- Go's netpoller integrates directly with OS kernel event loops (Linux `epoll`, macOS `kqueue`, Windows `IOCP`) so high-concurrency servers handle tens of thousands of connections with minimal CPU overhead.

## Go example

The example program resolves a domain name to its underlying IP addresses using DNS, spins up a local TCP server on an OS-assigned ephemeral port, establishes a client connection, and inspects the local and remote socket addresses:

```go
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
```

## Step-by-step execution

1. The client wants to fetch data from `example.com:443`.
2. The operating system checks its local DNS cache and `/etc/hosts` file.
3. If not cached, the client OS sends a UDP packet to port 53 of the configured DNS resolver (e.g., `8.8.8.8` or local router).
4. The resolver recursively queries root, TLD, and authoritative nameservers, returning the IPv4 address `93.184.216.34`.
5. The client OS selects an available ephemeral port (e.g., `51842`) and initiates the TCP 3-way handshake (`SYN` -> `SYN-ACK` -> `ACK`) to `93.184.216.34:443`.
6. Once the handshake finishes, a bidirectional socket file descriptor is active in both client and server processes.
7. Application layer protocols (HTTPS/HTTP) begin sending encrypted byte streams over this established TCP pipe.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Forgetting that domain names need DNS resolution | Assuming the computer connects to `"google.com"` directly | Remember that networks route exclusively by IP addresses; DNS lookups add network latency |
| Attempting to bind ports < 1024 without root privileges | Ports 0–1023 are restricted system ports | Use ports >= 1024 (e.g., 8080, 3000) for local application development |
| Getting "bind: address already in use" errors | Another process (or a previous un-killed instance of your app) is already listening on that port | Find and kill the conflicting PID using `lsof -i :<port>` (Linux/macOS) or `netstat -ano` (Windows) |
| Hardcoding IP addresses in code | Servers change IPs when re-provisioned or auto-scaled | Use DNS hostnames or service discovery names; avoid hardcoded IP addresses |

## Debugging walkthrough

**Scenario: A developer runs `go run main.go` to start a web server on port 8080, but gets `listen tcp :8080: bind: address already in use`.**
- **Cause:** Another process is already bound to TCP port 8080.
- **Diagnosis:**
  - Linux/macOS: Run `lsof -i :8080` to find the process name and PID.
  - Windows: Run `netstat -ano | findstr :8080` to find the holding PID.
- **Resolution:** Terminate the conflicting process (`kill <PID>` or `taskkill /F /PID <PID>`), or configure your application to use a different port (e.g., 8081).

**Scenario: A web service cannot connect to a remote database, failing with "no such host".**
- **Cause:** DNS resolution failed; the domain name does not exist, DNS servers are unreachable, or search domains are misconfigured.
- **Diagnosis:** Run `nslookup <hostname>` or `dig <hostname>` in your terminal to test DNS resolution directly.
- **Resolution:** Verify domain spelling, check network connectivity to DNS resolvers, and verify DNS server configuration.

## Production notes

In cloud-native production environments (AWS, GCP, Kubernetes), DNS and ports dictate network routing. Kubernetes uses internal DNS (`CoreDNS`) to resolve service names (e.g., `user-service.default.svc.cluster.local`), and ingress controllers forward incoming traffic from public ports 80/443 to internal container target ports. Understanding DNS caching and TTL (Time To Live) is critical: stale DNS records can cause traffic to hit decommissioned servers during zero-downtime deployments.

## Performance implications

- DNS queries take ~10–100 ms of latency; operating systems and HTTP clients cache DNS answers to prevent repeated queries.
- A TCP 3-way handshake requires one full network Round Trip Time (RTT, ~20–100 ms over the internet) before any HTTP data can be transmitted.
- HTTPS connections require an additional TLS handshake (~1–2 additional RTTs), which is why HTTP/2 and HTTP/3 support connection reuse and persistent keep-alive connections.

## Tests / verification

1. Query DNS records from your terminal:
   - POSIX: `dig example.com` or `host example.com`
   - Windows PowerShell: `Resolve-DnsName example.com`
   - Observe the returned `A` record containing the IPv4 address.
2. Inspect active listening ports on your machine:
   - Linux/macOS: `netstat -tuln` or `lsof -i -P -n | grep LISTEN`
   - Windows PowerShell: `Get-NetTCPConnection -State Listen | Select-Object -First 10 LocalAddress, LocalPort, OwningProcess`
3. Run the Go networking example:
   - `go run ./curriculum/modules/01-computers-terminal-git-web/14-web-preview-client-server-dns-and-ports`
   - Confirm DNS resolution, TCP listener binding, and ephemeral port connection.

## Review questions

1. What is the fundamental responsibility of DNS on the internet?
2. What are the three steps of the TCP 3-way handshake, and what is its purpose?
3. Why are ports needed in addition to IP addresses when connecting to a remote server?
4. What is the difference between a "Well-Known Port" and an "Ephemeral Port"?
5. What causes the common network error `bind: address already in use`, and how do you resolve it?

## NEXT UP

[Lesson 15: HTTP request and response preview](../15-http-request-and-response-preview/README.md)
