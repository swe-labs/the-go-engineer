// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 3: Collections & Pointers — Pointers
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What a pointer IS: a variable that stores a memory address
//   - The & operator (address-of) and * operator (dereference)
//   - Pass-by-value vs pass-by-reference
//   - Why pointers exist: mutation, large struct efficiency, nil-ability
//   - Pointer safety in Go (no pointer arithmetic)
//
// ENGINEERING DEPTH:
//   A pointer is just an integer holding the numeric hexadecimal address of a
//   byte in your computer's RAM (e.g. `0xc0000a6018`). In C++, you can freely
//   add or subtract from this integer (Pointer Arithmetic), which is incredibly
//   fast but is the #1 cause of catastrophic security exploits globally (buffer
//   overflows). Go explicitly bans Pointer Arithmetic. You get raw memory performance
//   without the devastating security risks.
//
// RUN: go run ./03-collections-and-pointers/4-pointers
// ============================================================================

// modifyValue receives a COPY of the int.
// Any changes to "val" inside this function do NOT affect the original.
// This is "pass by value" — Go's default behavior for all types.
func modifyValue(val int) {
	val = val * 10
	fmt.Printf("  Inside modifyValue: val = %d (this is a copy)\n", val)
}

// modifyPointer receives a POINTER to the int.
// The pointer holds the MEMORY ADDRESS of the original variable.
// Using *val (dereference) modifies the original variable.
func modifyPointer(val *int) {
	// Always check for nil before dereferencing.
	// Dereferencing a nil pointer causes a runtime PANIC (crash).
	if val == nil {
		fmt.Println("  val is nil — cannot dereference")
		return
	}
	*val = *val * 10 // *val reads/writes the value AT the address
	fmt.Printf("  Inside modifyPointer: *val = %d (original modified!)\n", *val)
}

func main() {

	// --- PASS BY VALUE (Default) ---
	fmt.Println("=== Pass by Value ===")
	num := 10
	modifyValue(num)
	fmt.Printf("  After modifyValue: num = %d (unchanged!)\n", num) // Still 10

	fmt.Println()

	// --- PASS BY POINTER (Reference) ---
	// The & operator gets the MEMORY ADDRESS of a variable.
	// &num means "the address where num is stored in memory"
	fmt.Println("=== Pass by Pointer ===")
	modifyPointer(&num)                                             // Pass the ADDRESS of num
	fmt.Printf("  After modifyPointer: num = %d (changed!)\n", num) // Now 100

	fmt.Println()

	// --- POINTER BASICS ---
	// A pointer variable stores a memory address, not a value.
	//
	// Type notation:
	//   *int   = "pointer to an int" (the type)
	//   &grade = "address of grade"  (create a pointer)
	//   *ptr   = "value at address"  (dereference — read/write the value)
	fmt.Println("=== Pointer Basics ===")
	grade := 50
	gradePtr := &grade // gradePtr now holds the memory address of grade

	fmt.Printf("  grade value:   %d\n", grade)                    // 50
	fmt.Printf("  grade address: %p\n", &grade)                   // 0xc0000b4008 (varies)
	fmt.Printf("  gradePtr:      %p (same address)\n", gradePtr)  // Same as above
	fmt.Printf("  *gradePtr:     %d (dereferenced)\n", *gradePtr) // 50

	// Modifying through the pointer changes the original
	*gradePtr = 95
	fmt.Printf("  After *gradePtr = 95: grade = %d\n", grade) // 95

	// --- WHEN TO USE POINTERS ---
	//
	// 1. MUTATION: When a function needs to modify the caller's variable
	//      func UpdateUser(u *User) { u.Name = "new" }
	//
	// 2. LARGE STRUCTS: Passing a 1MB struct by value copies 1MB to the stack.
	//      Passing a pointer copies only 8 bytes (the address).
	//      Rule of thumb: if a struct has more than 3-4 fields, use a pointer.
	//
	// 3. NIL-ABILITY: Pointers can be nil. Regular values cannot.
	//      var p *int = nil  ← valid (means "no value")
	//      var n int = nil   ← COMPILE ERROR
	//
	// 4. SHARING: When multiple parts of the code need the same data.
	//
	// SAFETY: Go has NO pointer arithmetic (unlike C/C++).
	// You cannot do p++ to move to the next memory address.
	// This eliminates an entire class of memory corruption bugs.

	// KEY TAKEAWAY:
	// - & creates a pointer (gets the address)
	// - * dereferences a pointer (accesses the value at the address)
	// - Go is pass-by-value by default. Use pointers to modify originals.
	// - Always check for nil before dereferencing.
	// - Go pointers are safe — no arithmetic, no manual memory management.
}
