package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Config validation on boot
//
// Run: go run ./10-production/04-configuration/5-config-validation-on-boot

func main() {
	fmt.Println("=== CFG.5 Config validation on boot ===")
	fmt.Println("Learn why services should reject invalid configuration before they begin handling traffic.")
	fmt.Println()
	fmt.Println("- Validate required config before serving traffic.")
	fmt.Println("- Collect configuration into one typed structure before deeper initialization.")
	fmt.Println("- Make startup errors specific enough that an operator can fix them quickly.")
	fmt.Println()
	fmt.Println("Fail-fast startup is kinder to operators because the error appears when the service launches, not minutes later under production load.")
}
