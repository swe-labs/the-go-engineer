package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Configuration files
//
// Run: go run ./10-production/04-configuration/2-configuration-files

func main() {
	fmt.Println("=== CFG.2 Configuration files ===")
	fmt.Println("Learn how file-based config complements environment variables when the shape grows beyond a handful of keys.")
	fmt.Println()
	fmt.Println("- Structured files help when config has nested shape.")
	fmt.Println("- Parsing is not validation; do both.")
	fmt.Println("- Document precedence between file values and environment overrides.")
	fmt.Println()
	fmt.Println("Config files become liabilities when the precedence rules are unclear or when production secrets sneak into checked-in defaults.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CFG.3")
	fmt.Println("Current: CFG.2 (configuration files)")
	fmt.Println("---------------------------------------------------")
}
