// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./08-modules-and-packages/3-versioning/_starter
package main

import "fmt"

// ============================================================================
// MP.3 Starter: Versioning Workshop
// ============================================================================
//
// TODO:
// 1. Implement String so versions print as vMAJOR.MINOR.PATCH.
// 2. Implement IsCompatible using the major-version rule.
// 3. Implement IsNewer by comparing major, then minor, then patch.
// 4. Read version_test.go and make all tests pass.
// 5. Keep the printed guidance about /v2 imports and replace clear and accurate.
// ============================================================================

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return "TODO"
}

func (v Version) IsCompatible(other Version) bool {
	return false
}

func (v Version) IsNewer(other Version) bool {
	return false
}

func main() {
	versions := []Version{
		{1, 0, 0},
		{1, 1, 0},
		{1, 1, 1},
		{2, 0, 0},
	}

	fmt.Println("=== Versioning Workshop Starter ===")
	fmt.Println()

	for i, v := range versions {
		if i == 0 {
			continue
		}

		prev := versions[i-1]
		status := "TODO: decide compatibility"
		if v.IsCompatible(prev) {
			status = "compatible"
		}

		fmt.Printf("  %s -> %s  %s\n", prev, v, status)
	}

	fmt.Println()
	fmt.Println("Explain in your final solution:")
	fmt.Println("  - why v2+ modules require /v2 in the import path")
	fmt.Println("  - when replace is useful in local development")
}
