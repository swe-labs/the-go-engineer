# Networking basics

## Learning objective

Create TCP servers and clients using Go's `net` package, understand ports, listening addresses (localhost vs 0.0.0.0), and the difference between TCP and UDP transport.

## Why this matters

Every HTTP service in Go ultimately calls `net.Listen("tcp", addr)`. The address you choose — `":8080"`, `"localhost:8080"`, or `"0.0.0.0:8080"` — determines who can connect. When you deploy to Docker, your container must listen on `0.0.0.0` to accept connections from outside the container. When you develop locally, you use `localhost` for security. Getting this wrong means "it works on my machine" but fails in production.

## Mental model

A TCP connection is a pipe between two processes. One process listens (server), the other dials (client). The server binds to a port and accepts connections. The client dials to an IP:port and establishes a socket. Data flows bidirectionally through the socket. Think of it as a phone call: the server publishes its number (IP:port), the client dials it, and once connected, both sides can talk.

## Core idea

| Concept | Description | Go type/function |
|---|---|---|
| Listen | Bind to a port and wait for connections | `net.Listen("tcp", addr)` |
| Accept | Block until a client connects | `listener.Accept()` |
| Dial | Connect to a remote server | `net.Dial("tcp", addr)` |
| Port | 16-bit number identifying the service | `1024-65535` (ephemeral) |
| localhost | Loopback interface (127.0.0.1) | Same machine only |
| 0.0.0.0 | All interfaces | Accessible from outside |

## Under the hood

`net.Listen("tcp", ":8080")` calls `syscall.Socket` (creating a TCP socket), then `syscall.Bind` (binding to port 8080 on all interfaces), then `syscall.Listen` (marking the socket as passive with a backlog queue). The backlog (default 128 on Linux) is the number of pending connections the kernel accepts before Go has called `Accept`. When `listener.Accept()` is called, Go blocks on `syscall.Accept4` which returns a new file descriptor for the client connection. Go wraps this fd in a `net.Conn` supporting `Read`, `Write`, `Close`.

`net.Dial("tcp", "127.0.0.1:8080")` calls `syscall.Socket` (creating a socket) and `syscall.Connect` (initiating the TCP three-way handshake). The handshake (SYN, SYN-ACK, ACK) happens in kernel space. Once complete, `Dial` returns a `net.Conn`.

## How Go uses it

Go's `net/http` package builds on `net.Listen` and `net.Dial`:
- `http.ListenAndServe(addr, handler)` creates a TCP listener and accepts connections forever.
- `http.ListenAndServeTLS` does the same with TLS wrapping.
- `http.Client` uses `net.Dial` (via `net.Dialer`) to connect to servers.
- `net/http/httptest` uses `net.Listen` to create test servers on random ports.

Production Go services also use the `net` package directly for:
- Custom protocol servers (memcached, Redis, gRPC builds on TCP)
- Health check endpoints on a separate port
- TCP proxies and load balancers
- UDP syslog or metrics collectors

## Go example

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("Listening on 127.0.0.1:%d\n", port)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		conn.Write([]byte("Hello, " + msg))
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "client\n")
	reply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("Got reply: %s", reply)
}
```

## Step-by-step execution

1. `net.Listen("tcp", "127.0.0.1:0")` creates a TCP socket bound to the loopback interface on an OS-assigned port (since `:0` means "pick any available port").
2. The server goroutine calls `listener.Accept()` which blocks until a client connects.
3. `net.Dial("tcp", addr)` creates a client socket and initiates the TCP handshake with the server.
4. The server's `Accept` returns the new connection.
5. The client sends `"client\n"`.
6. The server reads the line, prepends `"Hello, "`, and writes it back.
7. The client reads and prints the reply.

## Common mistakes

- **Listening on `localhost` in Docker.** `net.Listen("tcp", "localhost:8080")` binds to 127.0.0.1. Inside Docker, this means the service is only reachable from inside the container, not from other containers or the host. Always listen on `0.0.0.0:8080` or `":8080"` in containerized apps.
- **Not closing the listener.** Forgetting `listener.Close()` leaks file descriptors and prevents port reuse. Use `defer` immediately after successful `Listen`.
- **Assuming port 0 is always available.** Port 0 means "OS pick one", which is great for testing but unpredictable for production. In production, use a specific port from a configuration source.
- **Not handling `Accept` errors in a loop.** `listener.Accept()` can return temporary errors (e.g., `EINTR`). The standard pattern is `for { conn, err := listener.Accept(); if err != nil { return } }`.
- **Writing to a closed connection.** Writing to a closed TCP socket returns `io.EOF` or `net.ErrClosed`. Check errors after `Write` on persistent connections.

## Debugging walkthrough

Consider a Go HTTP server that works in development but cannot be reached from other containers in Docker Compose:

```go
http.ListenAndServe("localhost:8080", mux)
```

**Symptom:** Browser on the host connects successfully. Another container at `web:8080` gets "connection refused".

**Investigation:** Run `netstat -tlnp | grep 8080` inside the container. The output shows `127.0.0.1:8080` (loopback only). Other containers connect via Docker's bridge network, which uses the container's eth0 interface, not loopback.

**Root cause:** `localhost` binds to `127.0.0.1`, which only accepts connections from within the same network namespace. Other containers connect from a different IP.

**Fix:** Change to `":8080"` or `"0.0.0.0:8080"`:
```go
http.ListenAndServe(":8080", mux)
```

## Production notes

- Always make the listen address configurable via environment variable: `addr := os.Getenv("LISTEN_ADDR")` with default `":8080"`.
- Use `net.Listen("tcp", addr)` directly instead of `http.ListenAndServe` when you need custom keepalive settings or TLS configuration.
- Set TCP keepalive on long-lived connections: `tcpConn.SetKeepAlive(true)` and `tcpConn.SetKeepAlivePeriod(30*time.Second)`.
- For production, always specify the port explicitly to avoid conflicts. Document the port in your Dockerfile with `EXPOSE`.
- In Kubernetes, containers share the pod network namespace. The `localhost` address can reach other containers in the same pod.

## Performance implications

- Each `Accept` creates a new goroutine for the connection. For high-throughput servers, use a worker pool pattern.
- `net.Listen` with a zero port is useful for tests but has a cost: the OS must find an available port, which involves scanning the ephemeral port range.
- TCP handshake latency is typically 1-3 RTT. Reusing connections (HTTP keepalive) avoids this overhead.
- The listen backlog can fill during traffic spikes. Monitor `netstat -s` for "listen queue overflow" counters.

## Practice task

Write a function `StartEchoServer(addr string) (net.Listener, error)` that starts a TCP echo server (reads a line, writes it back). Write a function `SendMessage(addr, msg string) (string, error)` that connects, sends a message, and returns the reply. Write a `main()` that starts the server, sends a few messages, and prints the replies.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/05-networking-basics
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/05-networking-basics
```

## Review questions

1. What is the difference between listening on `":8080"` and `"localhost:8080"`?
2. Why does `net.Listen("tcp", ":0")` return a specific port number?
3. What happens if you write to a TCP connection after the remote side has closed it?
4. How does the listen backlog work, and what happens when it is full?
5. Why should `listener.Accept()` be called in a loop?

## NEXT UP

Docker basics — writing Dockerfiles for Go applications, building minimal images, and using FROM, RUN, COPY, CMD.
