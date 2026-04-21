package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Protobuf basics
//
// Run: go run ./06-backend-db/01-web-and-database/apis/4-protobuf-basics

func main() {
	fmt.Println("=== API.4 Protobuf basics ===")
	fmt.Println("Learn why schema-first messages exist and how protobuf defines transport contracts before code exists.")
	fmt.Println()
	fmt.Println("- Messages describe fields and their numeric tags.")
	fmt.Println("- Code generation turns one schema into many language bindings.")
	fmt.Println("- Field evolution rules protect compatibility across versions.")
	fmt.Println()
	fmt.Println("Schema mistakes are distributed mistakes because every generated client and server inherits them.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.5")
	fmt.Println("Current: API.4 (protobuf basics)")
	fmt.Println("---------------------------------------------------")
}
