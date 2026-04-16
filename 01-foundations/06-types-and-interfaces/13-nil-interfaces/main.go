// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Nil Interfaces
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The difference between nil interface and interface holding typed nil
//   - Why nil interfaces can cause panics
//   - How to properly check for nil
//
// RUN: go run ./01-foundations/06-types-and-interfaces/13-nil-interfaces
// ============================================================================

type Writer interface {
	Write(p []byte) (n int, err error)
}

type NullWriter struct{}

func (NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func process(w Writer) {
	if w == nil {
		fmt.Println("  Writer is truly nil")
		return
	}
	fmt.Println("  Writer is not nil, processing...")
	w.Write([]byte("test"))
}

func main() {
	fmt.Println("=== Nil Interfaces ===")
	fmt.Println()

	fmt.Println("--- Truly Nil Interface ---")
	var w Writer
	process(w)

	fmt.Println()
	fmt.Println("--- Interface Holding Typed Nil ---")
	var np *string = nil
	var i interface{} = np
	fmt.Printf("  Interface value: %v\n", i)
	fmt.Printf("  Is nil? %t (type is %T)\n", i == nil, i)

	fmt.Println()
	fmt.Println("--- NullWriter (Not Nil) ---")
	var nw Writer = NullWriter{}
	process(nw)

	fmt.Println()
	fmt.Println("--- Nil Pointer to Interface ---")
	var np2 *NullWriter = nil
	var i2 interface{} = np2
	fmt.Printf("  Interface value: %v\n", i2)
	fmt.Printf("  Is nil? %t (type is %T)\n", i2 == nil, i2)
	process(i2.(Writer))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Truly nil interface: both type and value are nil")
	fmt.Println("  - Interface holding typed nil: type is set, value is nil")
	fmt.Println("  - Check interface == nil before using")
	fmt.Println("  - A nil pointer stored in interface is NOT nil interface")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: TI.14 functional options")
	fmt.Println("   Current: TI.13 (nil interfaces)")
	fmt.Println("---------------------------------------------------")
}
