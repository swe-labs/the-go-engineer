// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 8: Semantic Versioning & Replace Directive
// Level: Intermediate → Advanced
// ============================================================================
//
// SEMANTIC VERSIONING IN GO:
//
//   v1.2.3
//   │ │ │
//   │ │ └── PATCH: bug fixes (backward compatible)
//   │ └──── MINOR: new features (backward compatible)
//   └────── MAJOR: breaking changes (NOT backward compatible)
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
//   Block specific versions (e.g., known buggy releases):
//     exclude github.com/some/pkg v1.2.3
//
// VENDORING:
//   go mod vendor       — copy all dependencies into ./vendor/
//   go build -mod=vendor — build using vendored dependencies only
//
//   Use Cases:
//   - Air-gapped environments (no internet access)
//   - Reproducible builds without relying on module proxy
//   - Some CI/CD pipelines require vendoring
//
// REFERENCES:
//   - https://go.dev/doc/modules/version-numbers
//   - https://semver.org/
// ============================================================================

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// IsCompatible checks if two versions are compatible (same major version)
func (v Version) IsCompatible(other Version) bool {
	return v.Major == other.Major
}

// IsNewer checks if v is newer than other
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
		{1, 1, 0}, // Minor bump: new features, backward compatible
		{1, 1, 1}, // Patch bump: bug fix
		{2, 0, 0}, // Major bump: BREAKING CHANGE
		{2, 1, 0},
	}

	fmt.Println("Version history:")
	for i, v := range versions {
		if i > 0 {
			prev := versions[i-1]
			compatible := v.IsCompatible(prev)
			emoji := "✅"
			if !compatible {
				emoji = "⚠️  BREAKING"
			}
			fmt.Printf("  %s → %s  %s\n", prev, v, emoji)
		}
	}

	fmt.Println()
	fmt.Println("Key takeaway:")
	fmt.Println("  In Go, v2+ modules require a /v2 suffix in the import path.")
	fmt.Println("  This allows v1 and v2 to coexist in the same binary.")
	fmt.Println()
	fmt.Println("  import \"github.com/example/pkg\"     ← v0.x or v1.x")
	fmt.Println("  import \"github.com/example/pkg/v2\"   ← v2.x")
}
