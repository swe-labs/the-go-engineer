package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateLessons() {
	fmt.Println("=== validate-lessons ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	allBundles := []struct {
		modules []Module
		items   []Item
	}{
		{core.Modules, core.Items},
		{electives.Modules, electives.Items},
	}

	// Statuses that should have content on disk
	contentStatuses := map[string]bool{
		"ready": true, "published": true, "live": true,
		"review": true, "complete": true,
	}

	// 1. Check lesson source directory structure
	warnedMissing := make(map[string]bool)
	warn := func(path string, msg string) {
		if !warnedMissing[path] {
			warnf("lesson path '%s': %s", path, msg)
			warnedMissing[path] = true
			ok = false
		}
	}

	for _, bundle := range allBundles {
		for _, mod := range bundle.modules {
			for _, item := range bundle.items {
				if item.ModuleID != mod.ID {
					continue
				}
				if item.Files == nil {
					continue
				}
				if !contentStatuses[item.Status] {
					continue
				}

				// Check readme exists
				if item.Files.ReadmePath != "" {
					rp := filepath.Join(curriculumDir, "..", item.Files.ReadmePath)
					if !pathExists(rp) {
						warn(item.Files.ReadmePath, "README not found")
					}
				}

				// Check main source file exists
				if item.Files.MainPath != "" {
					mp := filepath.Join(curriculumDir, "..", item.Files.MainPath)
					if !pathExists(mp) {
						warn(item.Files.MainPath, "main source not found")
					}
				}

				// Check starter code exists
				if item.Files.StarterPath != "" {
					sp := filepath.Join(curriculumDir, "..", item.Files.StarterPath)
					if !pathExists(sp) {
						warn(item.Files.StarterPath, "starter code not found")
					}
				}

				// Check solution exists
				if item.Files.SolutionPath != "" {
					sol := filepath.Join(curriculumDir, "..", item.Files.SolutionPath)
					if !pathExists(sol) {
						warn(item.Files.SolutionPath, "solution not found")
					}
				}

				// Check tests exist
				if item.Files.TestPath != "" {
					tp := filepath.Join(curriculumDir, "..", item.Files.TestPath)
					if !pathExists(tp) {
						warn(item.Files.TestPath, "test file not found")
					}
				}
			}
		}
	}

	// 2. Check RUN: headers in source files
	for _, bundle := range allBundles {
		for _, mod := range bundle.modules {
			for _, item := range bundle.items {
				if item.ModuleID != mod.ID {
					continue
				}
				if item.Files == nil || item.Files.MainPath == "" {
					continue
				}
				if item.Verification == nil {
					continue
				}
				mp := filepath.Join(curriculumDir, "..", item.Files.MainPath)
				data, err := os.ReadFile(mp)
				if err != nil {
					continue
				}
				content := string(data)
				runHeader := fmt.Sprintf("RUN: %s", item.Verification.RunCommand)
				if !strings.Contains(content, runHeader) {
					warnf("%s (%s): RUN header mismatch — expected '%s'", item.ID, item.Files.MainPath, runHeader)
					ok = false
				}
			}
		}
	}

	// 3. Check README NEXT UP footers
	for _, bundle := range allBundles {
		for _, mod := range bundle.modules {
			for _, item := range bundle.items {
				if item.ModuleID != mod.ID {
					continue
				}
				if item.Files == nil || item.Files.ReadmePath == "" {
					continue
				}
				rp := filepath.Join(curriculumDir, "..", item.Files.ReadmePath)
				data, err := os.ReadFile(rp)
				if err != nil {
					continue
				}
				content := string(data)
				lines := strings.Split(content, "\n")

				// Find the NEXT UP footer
				var hasFooter bool
				for _, line := range lines {
					if strings.Contains(line, "NEXT UP:") {
						hasFooter = true

						// Check it references a valid next item
						if len(item.NextItemIDs) > 0 {
							// Build a lookup from item IDs to titles
							titleByID := make(map[string]string)
							for _, i := range bundle.items {
								titleByID[i.ID] = i.Title
							}
							var found bool
							for _, nid := range item.NextItemIDs {
								if title, ok := titleByID[nid]; ok {
									if strings.Contains(line, title) || strings.Contains(line, nid) {
										found = true
										break
									}
								}
							}
							if !found {
								warnf("%s: NEXT UP footer may reference wrong item (expected one of %v)", item.ID, item.NextItemIDs)
								ok = false
							}
						}
						break
					}
				}
				if !hasFooter && item.Type == "lesson" {
					warnf("%s: missing NEXT UP footer in README", item.ID)
					ok = false
				}
			}
		}
	}

	// 4. Check Machine Role comments in source files (for v2.1 compliance)
	machineRoleCount := 0
	for _, bundle := range allBundles {
		for _, mod := range bundle.modules {
			for _, item := range bundle.items {
				if item.ModuleID != mod.ID {
					continue
				}
				if item.Files == nil || item.Files.MainPath == "" {
					continue
				}
				mp := filepath.Join(curriculumDir, "..", item.Files.MainPath)
				data, err := os.ReadFile(mp)
				if err != nil {
					continue
				}
				content := string(data)
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if strings.HasPrefix(trimmed, "//go:role") || strings.HasPrefix(trimmed, "/*go:role") {
						machineRoleCount++
					}
				}
			}
		}
	}
	fmt.Printf("  Machine role comments found: %d\n", machineRoleCount)

	if ok {
		fmt.Println("  All lesson checks passed.")
	} else {
		fmt.Println("  Lesson checks completed with warnings/errors.")
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
