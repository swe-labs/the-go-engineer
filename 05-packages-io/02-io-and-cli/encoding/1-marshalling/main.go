// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ============================================================================
// Stage 05: Encoding — JSON Marshalling (Go → JSON)
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - json.Marshal: converting Go structs to JSON byte slices
//   - json.MarshalIndent: pretty-printed JSON output
//   - Struct tags: controlling JSON field names (`json:"name"`)
//   - Tag options: omitempty, "-" (skip), string
//   - Why field names MUST be exported (uppercase) to appear in JSON
//   - time.Time serialization (automatic ISO 8601 format)
//
// ANALOGY:
//   Marshalling is like translating a Go struct into a language that
//   every system understands: JSON. Your Go server speaks Go internally,
//   but when it sends data to a browser, mobile app, or another API,
//   it translates to JSON — the universal data format of the web.
//
// ENGINEERING DEPTH:
//   Go's `encoding/json` package heavily relies on the `reflect` (Reflection)
//   package. At runtime, when you call `json.Marshal(product)`, Go dynamically
//   inspects the struct's memory layout, reads the metadata tags, and pieces
//   together a byte array. Because this happens dynamically at runtime rather
//   than strictly at compile time, JSON encoding in standard Go is slightly
//   slower than in languages with macro-based code generation (like Rust's Serde).
//   For 99% of web APIs, this overhead is utterly negligible.
//
// RUN: go run ./05-packages-io/02-io-and-cli/encoding/1-marshalling
// ============================================================================

// Product represents an item in an e-commerce catalog.
// STRUCT TAGS control how each field is serialized to JSON.
//
// TAG SYNTAX: `json:"json_field_name,option1,option2"`
//
// Common options:
//
//	`json:"name"`          → JSON field name is "name" (not "Name")
//	`json:"name,omitempty"` → Omit this field if it's the zero value
//	`json:"-"`             → NEVER include this field in JSON
//	`json:",string"`       → Encode number/bool as a JSON string
type Product struct {
	// ID will appear as "id" in JSON (lowercase, matching REST API conventions)
	ID int `json:"id"`

	// Name will appear as "name" in JSON
	Name string `json:"name"`

	// Price will appear as "price" in JSON
	Price float64 `json:"price"`

	// Description uses "omitempty": if the string is empty (""), this field
	// is completely OMITTED from the JSON output. This keeps API responses clean
	// by not sending {"description": ""} for products without descriptions.
	Description string `json:"description,omitempty"`

	// InStock will appear as "in_stock" (snake_case — common in REST APIs)
	InStock bool `json:"in_stock"`

	// CreatedAt serializes as ISO 8601 string ("2025-01-15T10:30:00Z").
	// Go's time.Time has a built-in JSON marshaller that uses RFC 3339 format.
	CreatedAt time.Time `json:"created_at"`

	// InternalSKU uses "-" tag: this field is NEVER included in JSON.
	// Use this for sensitive or internal-only data like passwords, keys, etc.
	InternalSKU string `json:"-"`

	// Tags uses "omitempty" with a slice: omitted when the slice is nil.
	// An empty initialized slice ([]string{}) WILL appear as "tags": [].
	Tags []string `json:"tags,omitempty"`
}

func main() {
	fmt.Println("=== JSON Marshalling: Go → JSON ===")
	fmt.Println()

	// --- BASIC MARSHALLING ---
	laptop := Product{
		ID:          1001,
		Name:        "ThinkPad X1 Carbon",
		Price:       1299.99,
		Description: "Ultra-lightweight business laptop with carbon fiber chassis",
		InStock:     true,
		CreatedAt:   time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC),
		InternalSKU: "SKU-TP-X1C-2025", // This will NOT appear in JSON
		Tags:        []string{"laptop", "business", "ultralight"},
	}

	// json.Marshal converts the struct to a []byte containing JSON.
	// The output is compact (no whitespace, no newlines).
	data, err := json.Marshal(laptop)
	if err != nil {
		log.Fatal("Marshal error:", err)
	}

	fmt.Println("1️⃣  json.Marshal (compact):")
	fmt.Printf("   %s\n\n", string(data))

	// --- PRETTY-PRINTED JSON ---
	// json.MarshalIndent adds newlines and indentation for readability.
	// Arguments: (value, prefix, indent)
	//   prefix: string prepended to each line (usually "")
	//   indent: string used for each indentation level (usually "  " or "	")
	prettyData, err := json.MarshalIndent(laptop, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent error:", err)
	}

	fmt.Println("2️⃣  json.MarshalIndent (pretty):")
	fmt.Println(string(prettyData))
	fmt.Println()

	// --- OMITEMPTY DEMO ---
	// Create a product WITHOUT a description or tags.
	// Fields with omitempty will be completely absent from the JSON.
	minimal := Product{
		ID:        1002,
		Name:      "USB-C Hub",
		Price:     49.99,
		InStock:   false,
		CreatedAt: time.Now(),
		// Description is "" → omitted from JSON (omitempty)
		// Tags is nil → omitted from JSON (omitempty)
	}

	minimalJSON, _ := json.MarshalIndent(minimal, "", "  ")
	fmt.Println("3️⃣  omitempty (no description, no tags):")
	fmt.Println(string(minimalJSON))
	fmt.Println()

	// --- UNEXPORTED FIELDS ---
	// CRITICAL: Only EXPORTED fields (uppercase first letter) are marshalled.
	// Unexported fields (lowercase) are INVISIBLE to encoding/json.
	type internal struct {
		Public  string `json:"public"` // ✅ Will appear in JSON
		private string // ❌ INVISIBLE to json.Marshal (lowercase)
	}

	secret := internal{Public: "visible", private: "hidden"}
	secretJSON, _ := json.Marshal(secret)
	fmt.Println("4️⃣  Unexported fields are invisible:")
	fmt.Printf("   %s\n", string(secretJSON)) // Only "public" appears
	fmt.Println("   (The 'private' field is not in the JSON!)")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - json.Marshal converts Go structs to JSON []byte")
	fmt.Println("  - json.MarshalIndent for readable output (debugging, APIs)")
	fmt.Println("  - Struct tags control field names: `json:\"my_name\"`")
	fmt.Println("  - omitempty: skip zero-value fields from output")
	fmt.Println("  - json:\"-\": never include this field (secrets, internal data)")
	fmt.Println("  - Only EXPORTED (Uppercase) fields appear in JSON")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: EN.2 JSON unmarshalling")
	fmt.Println("   Current: EN.1 (JSON marshalling)")
	fmt.Println("---------------------------------------------------")
}
