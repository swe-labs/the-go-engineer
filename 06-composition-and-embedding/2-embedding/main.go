// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 6: Composition & Embedding — Struct Embedding
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What embedding is: composition with PROMOTED methods and fields
//   - The difference between composition (named field) and embedding (anonymous field)
//   - How promoted fields work: access embedded fields directly
//   - Field shadowing: what happens when names collide
//   - Embedding for interface satisfaction
//
// ANALOGY:
//   Composition (previous lesson): "The phone HAS a camera" → phone.Camera.TakePhoto()
//   Embedding (this lesson):       "The phone HAS a camera" → phone.TakePhoto()
//
//   With embedding, the camera's methods are PROMOTED to the phone level.
//   You can call phone.TakePhoto() directly instead of phone.Camera.TakePhoto().
//   The camera is still there — you just get a shortcut.
//
// RUN: go run ./06-composition-and-embedding/2-embedding
// ============================================================================

// --- COMPONENT TYPES ---

// Dimensions represents physical measurements.
// This is a small, reusable component used by many types.
type Dimensions struct {
	Width  float64 // Width in centimeters
	Height float64 // Height in centimeters
	Weight float64 // Weight in grams
}

// IsPortable checks if the device is light enough to carry.
func (d Dimensions) IsPortable() bool {
	return d.Weight < 2000 // Less than 2kg
}

// FormFactor returns a human-readable size description.
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

// Battery represents a rechargeable battery.
type Battery struct {
	CapacityMAh int // Milliamp-hours (e.g., 4000 for a phone)
	ChargeLevel int // Current charge percentage (0-100)
}

// IsLow checks if battery needs charging.
func (b Battery) IsLow() bool {
	return b.ChargeLevel < 20
}

// Status returns a human-readable battery status.
func (b Battery) Status() string {
	if b.IsLow() {
		return fmt.Sprintf("🔴 %d%% (LOW — charge soon!)", b.ChargeLevel)
	}
	return fmt.Sprintf("🟢 %d%%", b.ChargeLevel)
}

// --- EMBEDDED TYPE ---

// Laptop EMBEDS Dimensions and Battery.
// Notice: no field names! Just the type names.
//
//	Dimensions    ← EMBEDDED (anonymous field)
//	Battery       ← EMBEDDED (anonymous field)
//	Brand string  ← named field (normal composition)
//
// PROMOTED: All fields and methods of Dimensions and Battery are
// accessible directly on Laptop:
//
//	laptop.Width          instead of  laptop.Dimensions.Width
//	laptop.IsPortable()   instead of  laptop.Dimensions.IsPortable()
//	laptop.IsLow()        instead of  laptop.Battery.IsLow()
type Laptop struct {
	Brand      string // Named field — NOT promoted
	Model      string // Named field — NOT promoted
	Dimensions        // EMBEDDED — fields and methods are promoted
	Battery           // EMBEDDED — fields and methods are promoted
}

// Describe uses both promoted fields (Width, Weight) and
// promoted methods (FormFactor, Status) without going through
// the embedded type names.
func (l Laptop) Describe() {
	fmt.Printf("  💻 %s %s\n", l.Brand, l.Model)
	fmt.Printf("     Size: %.0f × %.0f cm, %.0fg (%s)\n",
		l.Width, l.Height, l.Weight, l.FormFactor()) // Promoted from Dimensions
	fmt.Printf("     Portable: %t\n", l.IsPortable()) // Promoted from Dimensions
	fmt.Printf("     Battery: %s\n", l.Status())      // Promoted from Battery
}

// --- FIELD SHADOWING ---

// Tablet embeds Battery but also has its OWN ChargeLevel field.
// When names collide, the OUTER (parent) type wins — this is "shadowing".
type Tablet struct {
	Battery
	ChargeLevel string // Shadows Battery.ChargeLevel (different type!)
}

func main() {
	fmt.Println("=== Struct Embedding: Promoted Fields & Methods ===")
	fmt.Println()

	// --- Create a Laptop with embedded types ---
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

	// PROMOTED ACCESS: Call methods and access fields directly.
	// No need to write macbook.Dimensions.Width or macbook.Battery.Status()
	macbook.Describe()
	fmt.Println()

	// You can still access the embedded struct explicitly if needed:
	fmt.Printf("  Explicit access: macbook.Dimensions.Width = %.1f\n", macbook.Dimensions.Width)
	fmt.Printf("  Promoted access: macbook.Width = %.1f (same value!)\n", macbook.Width)
	fmt.Println()

	// --- Low battery example ---
	oldLaptop := Laptop{
		Brand:      "Lenovo",
		Model:      "ThinkPad X1",
		Dimensions: Dimensions{Width: 32.3, Height: 21.7, Weight: 1350},
		Battery:    Battery{CapacityMAh: 5200, ChargeLevel: 12},
	}
	oldLaptop.Describe()
	fmt.Println()

	// --- FIELD SHADOWING ---
	fmt.Println("--- Field Shadowing ---")
	t := Tablet{
		Battery:     Battery{CapacityMAh: 8000, ChargeLevel: 50},
		ChargeLevel: "half-charged", // Shadows Battery.ChargeLevel
	}

	// t.ChargeLevel refers to the OUTER (Tablet's) field — the string
	fmt.Printf("  t.ChargeLevel = %q (Tablet's own field)\n", t.ChargeLevel)
	// To access the shadowed inner field, use the explicit path
	fmt.Printf("  t.Battery.ChargeLevel = %d (Battery's field)\n", t.Battery.ChargeLevel)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Embedding = anonymous field → methods and fields are PROMOTED")
	fmt.Println("  - laptop.Width works instead of laptop.Dimensions.Width")
	fmt.Println("  - Outer fields SHADOW inner fields with the same name")
	fmt.Println("  - Embedding is NOT inheritance — it's syntactic sugar for composition")
	fmt.Println("  - Use embedding when you want the promoted shortcut syntax")
	fmt.Println("  - Use named fields when you want explicit access (Section 06/1)")
}
