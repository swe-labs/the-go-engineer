// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ============================================================================
// Stage 05: Encoding — JSON Unmarshalling (JSON → Go)
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - json.Unmarshal: converting JSON bytes into Go structs
//   - How JSON fields map to struct fields via struct tags
//   - Handling nested JSON objects with nested structs
//   - Unmarshalling into maps for dynamic/unknown JSON
//   - What happens with extra JSON fields (ignored) and missing fields (zero value)
//   - Why you must pass a POINTER to Unmarshal
//
// ANALOGY:
//   If Marshalling is translating FROM Go TO JSON (outgoing mail),
//   then Unmarshalling is translating FROM JSON TO Go (incoming mail).
//   Your server receives JSON from an API or browser → Unmarshal → Go struct.
//
// ENGINEERING DEPTH:
//   When you call `json.Unmarshal(data, &struct)`, you MUST pass a pointer.
//   Why? Because Go is strictly "pass-by-value". If you passed the struct directly,
//   `Unmarshal` would receive a disconnected copy, populate the data, and then
//   destroy it when the function exited. By passing a pointer, `Unmarshal` can
//   dereference the memory address and directly overwrite the fields of your
//   allocated struct in the Heap/Stack.
//
// RUN: go run ./05-packages-io/02-io-and-cli/encoding/2-unmarshal
// ============================================================================

// APIResponse represents the JSON response from a weather API.
// The struct tags tell json.Unmarshal which JSON field maps to which Go field.
type APIResponse struct {
	City      string   `json:"city"`
	TempC     float64  `json:"temp_celsius"`
	Humidity  int      `json:"humidity"`
	WindKmH   float64  `json:"wind_kmh"`
	Condition string   `json:"condition"`
	Forecast  Forecast `json:"forecast"` // Nested JSON object → nested struct
}

// Forecast is a nested struct for the "forecast" JSON object.
// Nested JSON: {"forecast": {"tomorrow_high": 28, "tomorrow_low": 19}}
type Forecast struct {
	TomorrowHigh float64 `json:"tomorrow_high"`
	TomorrowLow  float64 `json:"tomorrow_low"`
	Summary      string  `json:"summary,omitempty"`
}

// Simulated JSON from a weather API.
// In real code, this would come from http.Get() → resp.Body.
var weatherJSON = `{
  "city": "Dhaka",
  "temp_celsius": 32.5,
  "humidity": 78,
  "wind_kmh": 12.3,
  "condition": "Partly Cloudy",
  "forecast": {
    "tomorrow_high": 34.0,
    "tomorrow_low": 26.5,
    "summary": "Hot and humid"
  },
  "source": "weather-api.io"
}`

// This JSON has a MISSING field and an EXTRA field:
var partialJSON = `{
  "city": "Tokyo",
  "temp_celsius": 18.0,
  "condition": "Clear"
}`

func main() {
	fmt.Println("=== JSON Unmarshalling: JSON → Go ===")
	fmt.Println()

	// =====================================================================
	// 1. Basic Unmarshalling
	// =====================================================================
	// json.Unmarshal(data []byte, v any) error
	//   - data: the JSON bytes to parse
	//   - v: a POINTER to the target variable (&variable)
	//
	// WHY A POINTER? Unmarshal needs to MODIFY your variable.
	// Without &, it gets a copy and your variable stays unchanged.
	var weather APIResponse
	err := json.Unmarshal([]byte(weatherJSON), &weather) // &weather = pointer!
	if err != nil {
		log.Fatal("Unmarshal error:", err)
	}

	fmt.Println("1️⃣  Full JSON → Struct:")
	fmt.Printf("   City:      %s\n", weather.City)
	fmt.Printf("   Temp:      %.1f°C\n", weather.TempC)
	fmt.Printf("   Humidity:  %d%%\n", weather.Humidity)
	fmt.Printf("   Wind:      %.1f km/h\n", weather.WindKmH)
	fmt.Printf("   Condition: %s\n", weather.Condition)
	fmt.Printf("   Tomorrow:  %.0f°C / %.0f°C (%s)\n",
		weather.Forecast.TomorrowHigh,
		weather.Forecast.TomorrowLow,
		weather.Forecast.Summary)
	fmt.Println()

	// =====================================================================
	// 2. Missing & Extra Fields
	// =====================================================================
	// MISSING JSON fields → struct field gets its ZERO VALUE
	// EXTRA JSON fields → silently IGNORED (no error)
	// This is a safety feature: your struct only gets what it can handle.
	var partial APIResponse
	json.Unmarshal([]byte(partialJSON), &partial)

	fmt.Println("2️⃣  Partial JSON (missing fields → zero values):")
	fmt.Printf("   City:     %s\n", partial.City)                           // "Tokyo"
	fmt.Printf("   Temp:     %.1f°C\n", partial.TempC)                      // 18.0
	fmt.Printf("   Humidity: %d (zero — not in JSON)\n", partial.Humidity)  // 0
	fmt.Printf("   Wind:     %.1f (zero — not in JSON)\n", partial.WindKmH) // 0.0
	fmt.Println()

	// =====================================================================
	// 3. Unmarshal into a map (dynamic/unknown JSON)
	// =====================================================================
	// When you don't know the JSON structure ahead of time,
	// unmarshal into map[string]any. This handles any JSON shape.
	// The values are: string, float64, bool, nil, map[string]any, []any
	var dynamic map[string]any
	json.Unmarshal([]byte(weatherJSON), &dynamic)

	fmt.Println("3️⃣  Dynamic JSON → map[string]any:")
	fmt.Printf("   city = %v (type: %T)\n", dynamic["city"], dynamic["city"])
	fmt.Printf("   temp = %v (type: %T)\n", dynamic["temp_celsius"], dynamic["temp_celsius"])
	// Note: JSON numbers always become float64 in map[string]any
	fmt.Printf("   humidity = %v (type: %T — float64, not int!)\n",
		dynamic["humidity"], dynamic["humidity"])
	fmt.Println()

	// =====================================================================
	// 4. Round-trip: Unmarshal → modify → Marshal
	// =====================================================================
	// A common pattern: receive JSON, modify data, send JSON back.
	weather.TempC = 35.0
	weather.Condition = "Sunny"

	updatedJSON, _ := json.MarshalIndent(weather, "", "  ")
	fmt.Println("4️⃣  Modified and re-marshalled:")
	fmt.Println(string(updatedJSON))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - json.Unmarshal converts JSON bytes → Go structs")
	fmt.Println("  - ALWAYS pass a pointer: json.Unmarshal(data, &target)")
	fmt.Println("  - Missing JSON fields → zero value (no error)")
	fmt.Println("  - Extra JSON fields → silently ignored (no error)")
	fmt.Println("  - Use map[string]any for unknown/dynamic JSON structures")
	fmt.Println("  - JSON numbers are ALWAYS float64 in maps (cast if needed)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: EN.3 JSON encoder (stream)")
	fmt.Println("   Current: EN.2 (JSON unmarshalling)")
	fmt.Println("---------------------------------------------------")
}
