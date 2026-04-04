// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 6: Composition & Embedding — Composition
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What composition is: building complex types from simpler ones
//   - Composition vs Inheritance (Go chose composition, and here's why)
//   - The "Has-A" relationship: a Car HAS an Engine (composition)
//   - vs the "Is-A" relationship: a Car IS a Vehicle (inheritance)
//   - Accessing methods of composed types
//   - Reusing composed types for consistency
//
// ANALOGY:
//   Think of a smartphone. A smartphone is COMPOSED OF:
//     - A Camera module (has zoom, resolution, flash)
//     - A Battery (has capacity, charge level)
//     - A Screen (has size, brightness)
//     - A GPS module (has coordinates, accuracy)
//
//   The phone doesn't INHERIT from Camera or Battery.
//   Instead, it HAS these components. Each component knows how to do
//   its own job. The phone just coordinates them. That's composition.
//
// WHY GO CHOSE COMPOSITION:
//   Inheritance creates tight coupling and fragile hierarchies
//   (the "diamond problem", "yo-yo problem", "gorilla-banana problem").
//   Composition keeps each piece independent and swappable.
//   Go has NO class inheritance. Period. Composition is the way.
//
// RUN: go run ./06-composition-and-embedding/1-composition
// ============================================================================

// --- COMPONENT TYPES (small, focused, reusable) ---

// GPSLocation represents a geographic coordinate.
// This is a self-contained component — it knows nothing about who uses it.
type GPSLocation struct {
	Latitude  float64 // Degrees north (+) or south (-)
	Longitude float64 // Degrees east (+) or west (-)
	Label     string  // Human-readable name for this location
}

// String formats the location for display.
func (g GPSLocation) String() string {
	return fmt.Sprintf("%s (%.4f, %.4f)", g.Label, g.Latitude, g.Longitude)
}

// DistanceLabel gives a rough idea of the coordinate (for demo purposes).
func (g GPSLocation) DistanceLabel() string {
	if g.Latitude > 0 {
		return "Northern Hemisphere"
	}
	return "Southern Hemisphere"
}

// ContactInfo holds communication details.
// This is another self-contained component — reusable across many types.
type ContactInfo struct {
	Phone string // Primary phone number
	Email string // Primary email address
}

// Summary formats contact info for display.
func (c ContactInfo) Summary() string {
	return fmt.Sprintf("📞 %s | 📧 %s", c.Phone, c.Email)
}

// --- COMPOSED TYPES (built from components) ---

// Warehouse is COMPOSED OF GPSLocation and ContactInfo.
// "Composed of" means it HAS these types as fields — not that it IS them.
// To access GPS methods: warehouse.Location.DistanceLabel()
// To access Contact methods: warehouse.Contact.Summary()
type Warehouse struct {
	ID       int         // Unique warehouse identifier
	Name     string      // Warehouse name (e.g., "West Coast Hub")
	Capacity int         // Maximum number of items
	Location GPSLocation // COMPOSITION: Warehouse HAS a location
	Contact  ContactInfo // COMPOSITION: Warehouse HAS contact info
}

// PrintDetails displays all warehouse information.
// Notice how we access composed type methods with dot notation:
//
//	w.Location.String()  — calling GPSLocation's method
//	w.Contact.Summary()  — calling ContactInfo's method
func (w Warehouse) PrintDetails() {
	fmt.Printf("  🏭 Warehouse #%d: %s\n", w.ID, w.Name)
	fmt.Printf("     Capacity: %d items\n", w.Capacity)
	fmt.Printf("     Location: %s (%s)\n", w.Location, w.Location.DistanceLabel())
	fmt.Printf("     Contact:  %s\n", w.Contact.Summary())
}

// DeliveryRoute uses composition to link two warehouses.
// A route has an Origin and a Destination — both are GPSLocation.
// This shows the REUSABILITY of composed types.
type DeliveryRoute struct {
	RouteID     string      // Route identifier (e.g., "RT-001")
	Origin      GPSLocation // Where the delivery starts
	Destination GPSLocation // Where the delivery ends
}

// Describe prints a summary of the delivery route.
func (r DeliveryRoute) Describe() {
	fmt.Printf("  🚛 Route %s: %s → %s\n", r.RouteID, r.Origin.Label, r.Destination.Label)
}

func main() {
	fmt.Println("=== Composition: Building Complex Types from Simple Ones ===")
	fmt.Println()

	// --- Create composed types ---
	// Each component (GPSLocation, ContactInfo) is initialized inline.
	westCoast := Warehouse{
		ID:       1,
		Name:     "West Coast Distribution Hub",
		Capacity: 50000,
		Location: GPSLocation{
			Latitude:  37.7749,
			Longitude: -122.4194,
			Label:     "San Francisco, CA",
		},
		Contact: ContactInfo{
			Phone: "+1-415-555-0100",
			Email: "west-ops@logistics.io",
		},
	}

	eastCoast := Warehouse{
		ID:       2,
		Name:     "East Coast Fulfillment Center",
		Capacity: 75000,
		Location: GPSLocation{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Label:     "New York, NY",
		},
		Contact: ContactInfo{
			Phone: "+1-212-555-0200",
			Email: "east-ops@logistics.io",
		},
	}

	// Print each warehouse's details
	fmt.Println("--- Warehouses ---")
	westCoast.PrintDetails()
	fmt.Println()
	eastCoast.PrintDetails()
	fmt.Println()

	// --- REUSING COMPONENTS ---
	// The same GPSLocation type is used in both Warehouse and DeliveryRoute.
	// We can extract locations from warehouses and reuse them.
	route := DeliveryRoute{
		RouteID:     "RT-001",
		Origin:      westCoast.Location, // Reuse warehouse's GPSLocation
		Destination: eastCoast.Location, // Reuse warehouse's GPSLocation
	}

	fmt.Println("--- Delivery Route ---")
	route.Describe()

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Composition = 'Has-A' → Warehouse HAS a GPSLocation")
	fmt.Println("  - Inheritance = 'Is-A' → Go does NOT have this")
	fmt.Println("  - Components are small, focused types with their own methods")
	fmt.Println("  - Composed types access component methods via dot: w.Location.String()")
	fmt.Println("  - Components are reusable across many different parent types")
	fmt.Println("  - Next lesson: Embedding (promoted methods for cleaner syntax)")
}
