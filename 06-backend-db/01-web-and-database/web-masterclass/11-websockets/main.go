// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - WebSockets
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to upgrade a standard HTTP connection to a persistent WebSocket.
//   - How to handle full-duplex, real-time communication in Go.
//   - How to build a simple Echo Server using 'gorilla/websocket'.
//   - The importance of the "Message Loop" pattern.
//
// WHY THIS MATTERS:
//   - Traditional HTTP is "Pull-based"-the client must ask for data.
//     WebSockets are "Push-based"-the server can send data to the
//     client at any time. This is essential for building real-time
//     features like chat, live sports scores, or stock tickers.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/11-websockets
//
// KEY TAKEAWAY:
//   - Go's lightweight goroutines make it perfect for handling
//     thousands of concurrent WebSocket connections.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Stage 06: Web Masterclass - WebSockets
//
//   - Upgrading: Switching protocols (HTTP -> WS)
//   - Reading/Writing: Bidirectional streams
//   - Connection Lifecycle: Handling disconnects
//
// ENGINEERING DEPTH:
//   When you "Upgrade" a connection, the HTTP handler hands over
//   control of the underlying TCP socket to the WebSocket library.
//   Because Go uses lightweight goroutines, each WebSocket
//   connection typically lives in its own dedicated goroutine.
//   This allows Go to scale to tens of thousands of simultaneous
//   live connections on a single server, which would crash many
//   other language runtimes.

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// For learning, allow all origins. In production, be strict!
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	mux := http.NewServeMux()

	// 1. Serve a simple HTML page with JavaScript to test the WebSocket
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h2>WebSocket Echo Test</h2>
				<input id="input" type="text" placeholder="Type something...">
				<button onclick="send()">Send</button>
				<pre id="output"></pre>
				<script>
					const socket = new WebSocket("ws://localhost:8090/ws");
					socket.onmessage = (e) => {
						document.getElementById("output").innerText += "\n" + e.data;
					};
					function send() {
						const msg = document.getElementById("input").value;
						socket.send(msg);
					}
				</script>
			</body>
			</html>
		`)
	})

	// 2. The WebSocket endpoint
	mux.HandleFunc("GET /ws", handleWebSocket)

	fmt.Println("=== Web Masterclass: WebSockets ===")
	fmt.Println("  🚀 Server starting on http://localhost:8090")
	fmt.Println()
	fmt.Println("  1. Open http://localhost:8090 in your browser.")
	fmt.Println("  2. Type a message and hit Send.")
	fmt.Println("  3. Watch the server echo it back instantly!")

	log.Fatal(http.ListenAndServe(":8090", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: Section 07 - Concurrency")
	fmt.Println("Current: MC.11 (websockets)")
	fmt.Println("Previous: MC.10 (comments)")
	fmt.Println("---------------------------------------------------")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 3. Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("  [WS] New client connected.")

	// 4. Enter the message loop
	for {
		// Read message from browser
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("  [WS] Client disconnected.")
			break
		}

		log.Printf("  [WS] Received: %s", message)

		// Echo message back to browser
		err = conn.WriteMessage(mt, []byte("Server Echo: "+string(message)))
		if err != nil {
			log.Println("  [WS] Write failed.")
			break
		}
	}
}
