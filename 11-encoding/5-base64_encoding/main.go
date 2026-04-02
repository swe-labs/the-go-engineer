// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

// ============================================================================
// Section 11: Encoding — Base64 Encoding
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What Base64 is and why it's used
//   - base64.StdEncoding (standard base64)
//   - base64.URLEncoding (URL-safe base64)
//   - Encoding/Decoding strings and bytes
//
// ANALOGY:
//   Imagine you need to mail an irregularly shaped rock (binary data) in an envelope,
//   but the post office only accepts flat letters (text).
//   Base64 is a machine that crushes the rock, turns it into a perfectly flat letter
//   made of A-Z, a-z, 0-9, +, and /. On the receiving end, another machine
//   reads the letter and reconstructs the exact same rock.
//
// ENGINEERING DEPTH:
//   Base64 is NOT encryption! It provides zero security.
//   It takes every 3 bytes (24 bits) of raw data and splits it into 4 chunks
//   of 6 bits each. Each 6-bit chunk maps to one of 64 printable ASCII characters.
//   This means Base64 expands payload size by ~33%.
//
//   Why use it? Many protocols (HTTP headers, JSON, Email) were designed for text.
//   Base64 lets you safely transport binary data (images, tokens) over text channels.
//
// RUN: go run ./11-encoding/5-base64_encoding
// ============================================================================

func main() {
	fmt.Println("=== Base64 Encoding & Decoding ===")
	fmt.Println()

	// =====================================================================
	// 1. Standard Base64 Encoding
	// =====================================================================
	fmt.Println("1️⃣  Standard Base64:")

	secretMessage := "This string contains sensitive/raw data?"

	// Encode: Takes []byte, returns string
	// StdEncoding uses '+' and '/' as the 62nd and 63rd characters
	encoded := base64.StdEncoding.EncodeToString([]byte(secretMessage))

	fmt.Printf("   Original: %q\n", secretMessage)
	fmt.Printf("   Encoded:  %s\n", encoded) // Ends with = for padding

	// Decode: Takes string, returns []byte
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal("Decode error:", err)
	}

	fmt.Printf("   Decoded:  %q\n", string(decodedBytes))
	fmt.Println()

	// =====================================================================
	// 2. Binary Data Encoding (Images, Tokens, Keys)
	// =====================================================================
	fmt.Println("2️⃣  Binary Data Encoding:")

	// Simulated raw cryptographic key or binary data
	rawKey := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0xFF}

	fmt.Printf("   Raw Bytes:   %v\n", rawKey)

	encodedKey := base64.StdEncoding.EncodeToString(rawKey)
	fmt.Printf("   Base64 Key:  %s\n", encodedKey)
	fmt.Println()

	// =====================================================================
	// 3. URL-Safe Base64
	// =====================================================================
	// Standard base64 uses '+' and '/'.
	// These characters have special meaning in URLs (space and path separator).
	// URLEncoding replaces '+' with '-' and '/' with '_' so it can be used in URLs.
	fmt.Println("3️⃣  URL-Safe Base64:")

	tokenData := "user:rasel9t6|role:admin|path:/api/v1/update"

	// Standard encoding (creates a string that breaks URLs due to '/')
	stdToken := base64.StdEncoding.EncodeToString([]byte(tokenData))

	// URL encoding (safe to put in a query parameter)
	urlToken := base64.URLEncoding.EncodeToString([]byte(tokenData))

	fmt.Printf("   Standard: %s\n", stdToken)
	fmt.Printf("   URL-Safe: %s\n", urlToken)

	// Decoding URL-safe token
	decodedToken, _ := base64.URLEncoding.DecodeString(urlToken)
	fmt.Printf("   Decoded:  %s\n", string(decodedToken))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Base64 is for transporting binary data over text-only protocols")
	fmt.Println("  - It is NOT encryption (it is easily decodable by anyone)")
	fmt.Println("  - base64.StdEncoding for JSON, Headers, General purpose")
	fmt.Println("  - base64.URLEncoding for tokens inside URLs (swaps +/ for -_)")
	fmt.Println("  - Base64 increases data size by exactly 33%")
}
