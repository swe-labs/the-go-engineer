package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// validateReadmeContracts verifies that all instructional lessons in locked modules
// strictly satisfy the 6 non-negotiable core headings declared in readme.contracts.json:
// 1. ## Learning objective
// 2. ## Why this matters
// 3. ## Mental model
// 4. ## Common mistakes
// 5. ## Review questions
// 6. ## NEXT UP
func validateReadmeContracts() {
	fmt.Println("=== validate-readme-contracts ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	lockedModules := map[string]bool{
		"module-00": true,
		"module-01": true,
		"module-02": true,
	}

	requiredHeadings := []string{
		"Learning objective",
		"Why this matters",
		"Mental model",
		"Common mistakes",
		"Review questions",
		"NEXT UP",
	}

	lessonsChecked := 0
	contractViolations := 0

	for _, item := range core.Items {
		if item.Type != "lesson" {
			continue
		}
		if !lockedModules[item.ModuleID] {
			continue
		}

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

		lessonsChecked++
		foundHeadings := make(map[string]bool)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "## ") {
				continue
			}
			h2Text := strings.TrimSpace(strings.TrimPrefix(line, "## "))

			for _, req := range requiredHeadings {
				if strings.EqualFold(h2Text, req) || strings.HasPrefix(strings.ToLower(h2Text), strings.ToLower(req)) {
					foundHeadings[req] = true
				}
			}
		}
		file.Close()

		for _, req := range requiredHeadings {
			if !foundHeadings[req] {
				warnf("%s (%s): missing required heading '## %s' in %s", item.ID, item.Slug, req, item.Files.ReadmePath)
				contractViolations++
				ok = false
			}
		}
	}

	fmt.Printf("  Verified 6-heading core contracts across %d locked lessons (%d violations).\n", lessonsChecked, contractViolations)
	if ok && contractViolations == 0 {
		fmt.Println("  All locked lessons strictly comply with README contract specifications.")
	} else {
		fmt.Println("  README contract validation failed with violations.")
	}
}
