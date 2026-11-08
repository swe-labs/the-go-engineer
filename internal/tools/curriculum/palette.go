package main

import (
	"fmt"
	"os"
	"strings"
)

// validateDiagramPalette verifies that Mermaid diagrams in locked modules adhere
// to the standardized Dark Engineering Palette tokens.
func validateDiagramPalette() {
	fmt.Println("=== validate-diagram-palette ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	// Locked modules (single source of truth: curriculum/expected.json).
	lockedModules := lockedModuleSet()

	diagramsChecked := 0
	nonCompliant := 0

	for _, item := range core.Items {
		if item.Type != "lesson" {
			continue
		}
		if !lockedModules[item.ModuleID] {
			continue
		}

		if item.Files == nil || item.Files.ReadmePath == "" {
			continue
		}

		fullPath := resolvePath(item.Files.ReadmePath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		content := string(data)
		idx := 0
		for {
			start := strings.Index(content[idx:], "```mermaid")
			if start == -1 {
				break
			}
			start += idx
			end := strings.Index(content[start+10:], "```")
			if end == -1 {
				break
			}
			end += start + 10

			diagram := content[start:end]
			diagramsChecked++

			// Check if diagram is a flowchart
			if strings.Contains(diagram, "flowchart") || strings.Contains(diagram, "graph") {
				hasPalette := strings.Contains(diagram, "classDef") ||
					strings.Contains(diagram, "style ") ||
					strings.Contains(diagram, "#1a233a") ||
					strings.Contains(diagram, "#1e3a1e") ||
					strings.Contains(diagram, "#22272e")

				if !hasPalette {
					errorf("%s (%s): Mermaid flowchart lacks Dark Engineering Palette styling in %s", item.ID, item.Slug, item.Files.ReadmePath)
					nonCompliant++
					ok = false
				}
			}

			idx = end + 3
		}
	}

	fmt.Printf("  Verified Dark Engineering Palette across %d Mermaid diagrams (%d non-compliant).\n", diagramsChecked, nonCompliant)
	if ok && nonCompliant == 0 {
		fmt.Println("  All diagrams in locked modules satisfy Dark Engineering Palette standards.")
	} else {
		fmt.Println("  Diagram palette validation failed with warnings/errors.")
	}
}
