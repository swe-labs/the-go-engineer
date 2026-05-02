// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Embedding
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Utilizing struct embedding to promote fields and methods.
//   - Understanding the mechanics of anonymous fields vs named fields.
//   - Managing field shadowing and explicit name resolution.
//   - Leveraging embedding for rapid type construction and interface satisfaction.
//
// WHY THIS MATTERS:
//   - While standard composition is explicit, it can lead to verbose
//     code when a parent type needs to expose many behaviors of its
//     internal components. Embedding provides a concise mechanism
//     for "inheriting" method sets without the complexity of class
//     hierarchies. Understanding exactly how promotion works is
//     critical for building clean APIs that avoid naming collisions
//     and maintain the "Zero-Magic" transparency of the Go type system.
//
// RUN:
//   go run ./04-types-design/17-embedding
//
// KEY TAKEAWAY:
//   - Embedding promotes component methods to the outer type's method set.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import "fmt"

//
//   - Struct embedding as a mechanism for method set promotion.
//   - Implicit vs Explicit field access in nested structures.
//   - Selector resolution rules and field shadowing.
//   - Satisfying interfaces through embedded types.
//
// TECHNICAL RATIONALE:
//   While standard composition is explicit, it can lead to verbose
//   code when a parent type needs to expose many behaviors of its
//   internal components. Embedding provides a concise mechanism
//   for "promoting" method sets to the outer type without the complexity
//   of class hierarchies. Understanding exactly how promotion works
//   at the compiler level is critical for building clean APIs that
//   avoid naming collisions and maintain the "Zero-Magic" transparency
//   of the Go type system.
//

// --- COMPONENT TYPES ---

// Dimensions (Struct) provides physical measurement state as a component for embedding.
// Dimensions (Struct): (Struct) provides physical measurement state as a component for embedding.
type Dimensions struct {
	Width  float64
	Height float64
	Weight float64
}

// IsPortable (Method) implements portability logic based on the weight field.
// Dimensions.IsPortable (Method): (Method) implements portability logic based on the weight field.
func (d Dimensions) IsPortable() bool {
	return d.Weight < 2000
}

// FormFactor (Method) classifies structural size based on surface area calculation.
// Dimensions.FormFactor (Method): (Method) classifies structural size based on surface area calculation.
func (d Dimensions) FormFactor() string {
	area := d.Width * d.Height
	switch {
	case area < 100:
		return "compact"
	case area < 500:
		return "standard"
	default:
		return "large"
	}
}

// Battery (Struct) models power state as a component for structural promotion.
// Battery (Struct): (Struct) models power state as a component for structural promotion.
type Battery struct {
	CapacityMAh int
	ChargeLevel int
}

// IsLow (Method) determines if the charge level requires immediate attention.
// Battery.IsLow (Method): (Method) determines if the charge level requires immediate attention.
func (b Battery) IsLow() bool {
	return b.ChargeLevel < 20
}

// Status (Method) returns a human-readable summary of battery health.
// Battery.Status (Method): (Method) returns a human-readable summary of battery health.
func (b Battery) Status() string {
	if b.IsLow() {
		return fmt.Sprintf("LOW %d%% (charge soon!)", b.ChargeLevel)
	}
	return fmt.Sprintf("OK %d%%", b.ChargeLevel)
}

// Laptop (Struct) demonstrates method and field promotion through struct embedding.
// Laptop (Struct): (Struct) demonstrates method and field promotion through struct embedding.
type Laptop struct {
	Brand string
	Model string
	Dimensions
	Battery
}

// Describe (Method) generates a system summary by accessing both local and promoted fields.
// Laptop.Describe (Method): (Method) generates a system summary by accessing both local and promoted fields.
func (l Laptop) Describe() {
	fmt.Printf("  Laptop: %s %s\n", l.Brand, l.Model)
	// Width, FormFactor(), and Status() are all promoted.
	fmt.Printf("     Size: %.0f x %.0f cm, %.0fg (%s)\n",
		l.Width, l.Height, l.Weight, l.FormFactor())
	fmt.Printf("     Portable: %t\n", l.IsPortable())
	fmt.Printf("     Battery: %s\n", l.Status())
}

// Tablet (Struct) demonstrates field shadowing where outer fields hide embedded fields of the same name.
// Tablet (Struct): (Struct) demonstrates field shadowing where outer fields hide embedded fields of the same name.
type Tablet struct {
	Battery
	ChargeLevel string // Shadows Battery.ChargeLevel
}

func main() {
	fmt.Println("=== Struct Embedding: Method & Field Promotion ===")
	fmt.Println()

	// 1. Promoted Initialization.
	// We instantiate the outer type using the embedded type name as the field key.
	fmt.Println("--- Constructing Embedded Entities ---")
	macbook := Laptop{
		Brand: "Apple",
		Model: "MacBook Pro 14\"",
		Dimensions: Dimensions{
			Width:  31.3,
			Height: 22.1,
			Weight: 1600,
		},
		Battery: Battery{
			CapacityMAh: 7000,
			ChargeLevel: 85,
		},
	}

	// 2. Direct Access (Promotion).
	// We can access 'Width' and 'Status()' directly on 'macbook'
	// without traversing the '.Dimensions' or '.Battery' fields.
	macbook.Describe()
	fmt.Println()

	// 3. Name Shadowing.
	// When the outer type defines a field that exists in an embedded type,
	// the outer field "shadows" the inner one.
	fmt.Println("--- Field Shadowing Mechanics ---")
	t := Tablet{
		Battery:     Battery{CapacityMAh: 8000, ChargeLevel: 50},
		ChargeLevel: "half-charged", // Shadows Battery.ChargeLevel
	}

	// Tablet's ChargeLevel (string) wins.
	fmt.Printf("  t.ChargeLevel (Outer): %q\n", t.ChargeLevel)
	// Battery's ChargeLevel (int) is still accessible via explicit path.
	fmt.Printf("  t.Battery.ChargeLevel (Inner): %d\n", t.Battery.ChargeLevel)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CO.3 -> 04-types-design/18-bank-account-project")
	fmt.Println("Run    : go run ./04-types-design/18-bank-account-project")
	fmt.Println("Current: CO.2 (embedding)")
	fmt.Println("---------------------------------------------------")
}
