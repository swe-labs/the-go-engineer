package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var approvedPioneers = []string{
	"Alan Turing",
	"Dennis Ritchie",
	"Ken Thompson",
	"Rob Pike",
	"Russ Cox",
	"Robert Griesemer",
	"Edsger W. Dijkstra",
	"Dijkstra",
	"C.A.R. Hoare",
	"Tony Hoare",
	"Donald Knuth",
	"Donald E. Knuth",
	"Brian Kernighan",
	"Brian W. Kernighan",
	"John von Neumann",
	"Linus Torvalds",
	"Maurice Wilkes",
	"Maurice V. Wilkes",
	"Doug McIlroy",
	"Fred Brooks",
	"Leslie Lamport",
	"Claude Shannon",
	"Claude E. Shannon",
	"Alan Kay",
	"Alan Perlis",
	"Alan J. Perlis",
	"David Parnas",
	"Butler Lampson",
	"Andrew Tanenbaum",
	"Andrew S. Tanenbaum",
	"Peter Landin",
	"Robin Milner",
	"Ada Lovelace",
	"George Boole",
	"Vint Cerf",
	"Dave Cheney",
	"Rebecca Stambler",
	"Martin Fowler",
	"Andrew Gerrand",
	"Eric S. Raymond",
	"Eric Raymond",
	"Roy Fielding",
}

func isApprovedPioneer(attribution string) bool {
	for _, p := range approvedPioneers {
		if strings.Contains(attribution, p) {
			return true
		}
	}
	return false
}

func checkLessonAnchorQuote(fullPath string) (hasQuote, hasAttribution, hasPioneer bool, err error) {
	file, err := os.Open(fullPath)
	if err != nil {
		return false, false, false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	lineCount := 0

	for scanner.Scan() && lineCount < 25 {
		lineCount++
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "> ") || strings.HasPrefix(line, ">\"") || strings.HasPrefix(line, "> \"") {
			hasQuote = true
			if strings.Contains(line, "—") || strings.Contains(line, "--") {
				hasAttribution = true
				if isApprovedPioneer(line) {
					hasPioneer = true
					break
				}
			}
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return hasQuote, hasAttribution, hasPioneer, scanErr
	}
	return hasQuote, hasAttribution, hasPioneer, nil
}

// validateAnchorQuotes verifies that every instructional lesson in locked modules
// opens with an authoritative anchor quote from a computing or Go pioneer.
func validateAnchorQuotes() {
	fmt.Println("=== validate-anchor-quotes ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	// Modules audited and locked into zero-magic pioneer anchor quote requirements
	// (single source of truth: curriculum/expected.json).
	lockedModules := lockedModuleSet()

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
			errorf("%s: missing files.readme_path in metadata", item.ID)
			ok = false
			continue
		}

		fullPath := resolvePath(item.Files.ReadmePath)
		hasQuote, hasAttribution, hasPioneer, err := checkLessonAnchorQuote(fullPath)
		if err != nil {
			errorf("%s: cannot open README at %s: %v", item.ID, fullPath, err)
			ok = false
			continue
		}

		if !hasQuote {
			errorf("%s (%s): missing pioneer anchor quote before ## Learning objective in %s", item.ID, item.Slug, item.Files.ReadmePath)
			missingQuotes++
			ok = false
		} else if !hasAttribution {
			errorf("%s (%s): anchor quote lacks attribution (— Author) in %s", item.ID, item.Slug, item.Files.ReadmePath)
			missingQuotes++
			ok = false
		} else if !hasPioneer {
			errorf("%s (%s): anchor quote author is not a recognized computing pioneer in %s", item.ID, item.Slug, item.Files.ReadmePath)
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
