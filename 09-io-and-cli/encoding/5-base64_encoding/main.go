// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./09-io-and-cli/encoding/5-base64_encoding
package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

// ============================================================================
// Section 09: I/O and CLI — Base64 Encoding
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What base64 is and why it is used
//   - base64.StdEncoding for general text-safe transport
//   - base64.URLEncoding for URL-safe tokens
//   - Encoding and decoding strings and raw bytes
//
// ENGINEERING DEPTH:
//   Base64 is not encryption. It is a text transport format for binary data.
//   It increases payload size by roughly 33%, but it lets binary values survive
//   in systems that were designed around printable text.
// ============================================================================

func main() {
	fmt.Println("=== Base64 Encoding and Decoding ===")
	fmt.Println()

	fmt.Println("1) Standard Base64:")
	secretMessage := "This string contains sensitive/raw data?"
	encoded := base64.StdEncoding.EncodeToString([]byte(secretMessage))

	fmt.Printf("   Original: %q\n", secretMessage)
	fmt.Printf("   Encoded:  %s\n", encoded)

	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal("Decode error:", err)
	}

	fmt.Printf("   Decoded:  %q\n", string(decodedBytes))
	fmt.Println()

	fmt.Println("2) Binary Data Encoding:")
	rawKey := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0xFF}
	fmt.Printf("   Raw Bytes:   %v\n", rawKey)
	fmt.Printf("   Base64 Key:  %s\n", base64.StdEncoding.EncodeToString(rawKey))
	fmt.Println()

	fmt.Println("3) URL-Safe Base64:")
	tokenData := "user:rasel9t6|role:admin|path:/api/v1/update"
	stdToken := base64.StdEncoding.EncodeToString([]byte(tokenData))
	urlToken := base64.URLEncoding.EncodeToString([]byte(tokenData))

	fmt.Printf("   Standard: %s\n", stdToken)
	fmt.Printf("   URL-Safe: %s\n", urlToken)

	decodedToken, err := base64.URLEncoding.DecodeString(urlToken)
	if err != nil {
		log.Fatal("URL-safe decode error:", err)
	}

	fmt.Printf("   Decoded:  %s\n", string(decodedToken))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Base64 is for transporting binary data over text-only protocols")
	fmt.Println("  - It is NOT encryption (it is easily decodable by anyone)")
	fmt.Println("  - base64.StdEncoding works for JSON, headers, and general transport")
	fmt.Println("  - base64.URLEncoding avoids + and / in URL contexts")
	fmt.Println("  - Base64 increases data size by about 33%")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: EN.6 config parser")
	fmt.Println("   Current: EN.5 (Base64)")
	fmt.Println("---------------------------------------------------")
}
