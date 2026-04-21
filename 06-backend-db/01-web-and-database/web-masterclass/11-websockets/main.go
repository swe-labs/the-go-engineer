// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./06-backend-db/01-web-and-database/web-masterclass/11-websockets
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ============================================================================
// Stage 06: Web Masterclass — WebSockets
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to upgrade a standard HTTP GET request into a persistent WebSocket.
//   - Using `gorilla/websocket`, the industry standard WS package for Go.
//   - Implementing a basic message pump (echo server).
//
// ENGINEERING DEPTH:
//   Standard HTTP is half-duplex and stateless. WebSockets provide a full-duplex,
//   persistent TCP connection over a single socket.
//   In Go, each HTTP request runs in its own Goroutine. When you upgrade the
//   connection to a WebSocket, that Goroutine "hijacks" the underlying TCP connection
//   and stays alive as long as the socket is open. This is incredibly efficient
//   in Go compared to thread-per-connection languages.
// ============================================================================

// 1. The Upgrader
// We pre-configure an Upgrader. This dictates buffer sizes and origin checks.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// For educational purposes, allow connections from any origin.
	// In production, ALWAYS verify the Origin header to prevent CSWSH attacks!
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 2. The WebSocket Handler
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected from %s", conn.RemoteAddr().String())

	// 3. The Message Pump
	// An infinite loop that reads from the socket and writes back to it.
	for {
		// Read a message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break // Exit the loop and close the connection
		}

		log.Printf("Received: %s", message)

		// Write the message back (Echo)
		err = conn.WriteMessage(messageType, []byte(fmt.Sprintf("Server Echo: %s", message)))
		if err != nil {
			log.Printf("Write failed: %v", err)
			break
		}
	}
	log.Printf("Client disconnected: %s", conn.RemoteAddr().String())
}

func main() {
	// We serve a trivial HTML page to act as the client.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Go WebSockets</title></head>
		<body>
			<h1>Gorilla WebSocket Echo</h1>
			<input type="text" id="msg" placeholder="Type a message..." />
			<button onclick="send()">Send</button>
			<ul id="log"></ul>
			<script>
				const ws = new WebSocket("ws://localhost:8080/ws");
				ws.onmessage = function(e) {
					const li = document.createElement("li");
					li.textContent = "Re: " + e.data;
					document.getElementById("log").appendChild(li);
				};
				function send() {
					const input = document.getElementById("msg");
					ws.send(input.value);
					input.value = "";
				}
			</script>
		</body>
		</html>
		`
		w.Write([]byte(html))
	})

	http.HandleFunc("/ws", echoHandler)

	fmt.Println("🚀 WebSocket Server running on http://localhost:8080")
	fmt.Println("   Open the URL in your browser to test the live echo!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
