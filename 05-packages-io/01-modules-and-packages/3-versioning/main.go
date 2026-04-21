// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./05-packages-io/01-modules-and-packages/3-versioning
package main

import "fmt"

// ============================================================================
// Stage 05: Modules and Packages — Versioning Workshop
// Level: Intermediate → Advanced
// ============================================================================
//
// SEMANTIC VERSIONING IN GO:
//
//   v1.2.3
//   │ │ │
//   │ │ └── PATCH: bug fixes (backward compatible)
//   │ └──── MINOR: new features (backward compatible)
//   └────── MAJOR: breaking changes (not backward compatible)
//
// MAJOR VERSION RULE (v2+):
//   When a module reaches v2, the import path MUST include /v2:
//     import "github.com/example/pkg/v2"
//   This allows v1 and v2 to coexist in the same project.
//
// THE REPLACE DIRECTIVE:
//   Use `replace` in go.mod for:
//   1. Local development of a dependency
//   2. Fork overrides
//   3. Debugging dependency issues
//
//   Example in go.mod:
//     replace github.com/original/pkg => ../local-fork
//     replace github.com/original/pkg => github.com/myfork/pkg v1.0.0
//
// THE EXCLUDE DIRECTIVE:
//   Block specific versions (for example, known buggy releases):
//     exclude github.com/some/pkg v1.2.3
//
// VENDORING:
//   go mod vendor        — copy all dependencies into ./vendor/
//   go build -mod=vendor — build using vendored dependencies only
// ============================================================================

// Version represents a semantic version.
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// IsCompatible checks whether two versions share the same major version.
func (v Version) IsCompatible(other Version) bool {
	return v.Major == other.Major
}

// IsNewer compares major, then minor, then patch.
func (v Version) IsNewer(other Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch > other.Patch
}

func main() {
	fmt.Println("=== Semantic Versioning Demo ===")
	fmt.Println()

	versions := []Version{
		{1, 0, 0},
		{1, 1, 0},
		{1, 1, 1},
		{2, 0, 0},
		{2, 1, 0},
	}

	fmt.Println("Version history:")
	for i, v := range versions {
		if i == 0 {
			continue
		}

		prev := versions[i-1]
		status := "✅ compatible"
		if !v.IsCompatible(prev) {
			status = "⚠️  BREAKING"
		}

		fmt.Printf("  %s → %s  %s\n", prev, v, status)
	}

	fmt.Println()
	fmt.Println("Key takeaway:")
	fmt.Println("  In Go, v2+ modules require a /v2 suffix in the import path.")
	fmt.Println("  This allows v1 and v2 to coexist in the same binary.")
	fmt.Println()
	fmt.Println("  import \"github.com/example/pkg\"    ← v0.x or v1.x")
	fmt.Println("  import \"github.com/example/pkg/v2\" ← v2.x")
	fmt.Println()
	fmt.Println("Replace is most useful during local development or controlled fork testing.")
	fmt.Println("It should clarify dependency resolution, not hide long-term version problems.")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT STEP: Optional MP.4 build tags")
	fmt.Println("Then continue to Stage 05 when you're ready.")
	fmt.Println("Current: MP.3 (versioning workshop)")
	fmt.Println("---------------------------------------------------")
}
