package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// validateAnchorQuotes verifies that every instructional lesson in locked modules
// opens with an authoritative anchor quote from a computing or Go pioneer.
func validateAnchorQuotes() {
	fmt.Println("=== validate-anchor-quotes ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	// Modules audited and locked into zero-magic pioneer anchor quote requirements
	lockedModules := map[string]bool{
		"module-00": true,
		"module-01": true,
		"module-02": true,
		"module-03": true,
		"module-04": true,
	}

	checkedLessons := 0
	missingQuotes := 0

	for _, item := range core.Items {
		if item.Type != "lesson" {
			continue
		}
		if !lockedModules[item.ModuleID] {
			continue
		}

		checkedLessons++
		if item.Files == nil || item.Files.ReadmePath == "" {
			warnf("%s: missing files.readme_path in metadata", item.ID)
			ok = false
			continue
		}

		fullPath := resolvePath(item.Files.ReadmePath)
		file, err := os.Open(fullPath)
		if err != nil {
			warnf("%s: cannot open README at %s: %v", item.ID, fullPath, err)
			ok = false
			continue
		}

		scanner := bufio.NewScanner(file)
		lineCount := 0
		hasQuote := false

		for scanner.Scan() && lineCount < 25 {
			lineCount++
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "> ") || strings.HasPrefix(line, ">\"") || strings.HasPrefix(line, "> \"") {
				hasQuote = true
				break
			}
		}
		file.Close()

		if !hasQuote {
			warnf("%s (%s): missing pioneer anchor quote before ## Learning objective in %s", item.ID, item.Slug, item.Files.ReadmePath)
			missingQuotes++
			ok = false
		}
	}

	fmt.Printf("  Verified pioneer anchor quotes across %d locked lessons (%d missing).\n", checkedLessons, missingQuotes)
	if ok && missingQuotes == 0 {
		fmt.Println("  All locked lessons contain authoritative pioneer anchor quotes.")
	} else {
		fmt.Println("  Anchor quote validation failed for locked modules.")
	}
}
