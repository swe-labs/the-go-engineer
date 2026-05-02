// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Composition
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Building complex types using named-field composition.
//   - The "has-a" relationship vs the "is-a" (inheritance) model.
//   - Accessing inner component methods through explicit dot notation.
//   - Designing reusable domain components (GPSLocation, ContactInfo).
//
// WHY THIS MATTERS:
//   - Go deliberately omits class-based inheritance to avoid the tight
//     coupling and complexity of deep type hierarchies. Composition
//     allows developers to build robust systems by combining small,
//     focused building blocks. This "Zero-Magic" approach ensures
//     that ownership and behavior are always explicit, making the
//     codebase easier to audit and refactor.
//
// RUN:
//   go run ./04-types-design/16-composition
//
// KEY TAKEAWAY:
//   - Composition enables flexible type construction through explicit ownership.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import "fmt"

//
//   - Explicit vs Implicit behavior in type construction.
//   - Memory layout implications of nested structs.
//   - Selector resolution for named fields.
//   - Designing reusable domain components.
//
// TECHNICAL RATIONALE:
//   Go deliberately omits class-based inheritance to avoid the tight
//   coupling and fragility inherent in deep type hierarchies (e.g., the
//   fragile base class problem). Composition allows developers to build
//   complex types by combining small, independent units of state and
//   behavior. This ensures that ownership is always explicit via named
//   fields, making the execution flow predictable and the codebase
//   easier to reason about at scale.
//

// --- COMPONENT TYPES (small, focused, reusable) ---

// GPSLocation (Struct): encapsulates geographic coordinates and metadata as a reusable value type.
type GPSLocation struct {
	Latitude  float64
	Longitude float64
	Label     string
}

// String (Method) implements the fmt.Stringer interface for formatted coordinate output.
// GPSLocation.String (Method): implements the fmt.Stringer interface for formatted coordinate output.
func (g GPSLocation) String() string {
	return fmt.Sprintf("%s (%.4f, %.4f)", g.Label, g.Latitude, g.Longitude)
}

// DistanceLabel (Method) provides domain-specific hemispheric classification based on latitude.
// GPSLocation.DistanceLabel (Method): provides domain-specific hemispheric classification based on latitude.
func (g GPSLocation) DistanceLabel() string {
	if g.Latitude > 0 {
		return "Northern Hemisphere"
	}
	return "Southern Hemisphere"
}

// ContactInfo (Struct): encapsulates communication metadata into a reusable component.
type ContactInfo struct {
	Phone string
	Email string
}

// Summary (Method) generates a single-line string representation of contact details.
// ContactInfo.Summary (Method): generates a single-line string representation of contact details.
func (c ContactInfo) Summary() string {
	return fmt.Sprintf("phone: %s | email: %s", c.Phone, c.Email)
}

// Warehouse (Struct): demonstrates named-field composition by aggregating location and contact components.
type Warehouse struct {
	ID       int
	Name     string
	Capacity int
	Location GPSLocation // Composition: Warehouse HAS a GPSLocation
	Contact  ContactInfo // Composition: Warehouse HAS ContactInfo
}

// PrintDetails (Method) orchestrates component methods to output a full entity summary.
// Warehouse.PrintDetails (Method): orchestrates component methods to output a full entity summary.
func (w Warehouse) PrintDetails() {
	fmt.Printf("  Warehouse #%d: %s\n", w.ID, w.Name)
	fmt.Printf("     Capacity: %d items\n", w.Capacity)
	// Accessing component methods explicitly via dot notation.
	fmt.Printf("     Location: %s (%s)\n", w.Location, w.Location.DistanceLabel())
	fmt.Printf("     Contact:  %s\n", w.Contact.Summary())
}

// DeliveryRoute (Struct): demonstrates structural reuse of the GPSLocation component across different domains.
type DeliveryRoute struct {
	RouteID     string
	Origin      GPSLocation
	Destination GPSLocation
}

// Describe (Method) provides a text-based summary of the route's path.
// DeliveryRoute.Describe (Method): provides a text-based summary of the route's path.
func (r DeliveryRoute) Describe() {
	fmt.Printf("  Route %s: %s -> %s\n", r.RouteID, r.Origin.Label, r.Destination.Label)
}

func main() {
	fmt.Println("=== Composition: Explicit Component Ownership ===")
	fmt.Println()

	// 1. Component Initialization.
	// We instantiate the outer type by providing concrete values for its
	// inner component fields.
	fmt.Println("--- Constructing Composed Entities ---")
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
	westCoast.PrintDetails()
	fmt.Println()

	// 2. Component Reuse.
	// The GPSLocation component can be extracted and reused in a
	// completely different type (DeliveryRoute) without any changes
	// to the component itself.
	fmt.Println("--- Reusing Shared Components ---")
	route := DeliveryRoute{
		RouteID: "RT-001",
		Origin:  westCoast.Location,
		Destination: GPSLocation{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Label:     "New York, NY",
		},
	}
	route.Describe()

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CO.2 -> 04-types-design/17-embedding")
	fmt.Println("Run    : go run ./04-types-design/17-embedding")
	fmt.Println("Current: CO.1 (composition)")
	fmt.Println("---------------------------------------------------")
}
