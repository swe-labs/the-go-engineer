package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// lessonIssue is one content-shape issue found by the lesson validator.
type lessonIssue struct {
	Sev severity
	Msg string
}

func validateLessons() {
	fmt.Println("=== validate-lessons ===")

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	allItems := make([]Item, 0, len(core.Items)+len(electives.Items))
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	// Statuses that should have content on disk.
	contentStatuses := map[string]bool{
		"ready": true, "published": true, "live": true,
		"review": true, "complete": true,
	}

	// Map every item's resolved README path to its id, so a NEXT UP link can be
	// validated against the item's declared next_item_ids.
	idByReadme := make(map[string]string, len(allItems))
	for _, it := range allItems {
		if it.Files != nil && it.Files.ReadmePath != "" {
			idByReadme[filepath.Clean(resolvePath(it.Files.ReadmePath))] = it.ID
		}
	}

	// 1. Check lesson source directory structure (structural: always blocks).
	warnedMissing := make(map[string]bool)
	missing := func(path, msg string) {
		if !warnedMissing[path] {
			errorf("lesson path '%s': %s", path, msg)
			warnedMissing[path] = true
		}
	}

	for _, item := range allItems {
		if item.Files == nil || !contentStatuses[item.Status] {
			continue
		}
		check := func(path, label string) {
			if path != "" && !pathExists(resolvePath(path)) {
				missing(path, label+" not found")
			}
		}
		check(item.Files.ReadmePath, "README")
		check(item.Files.MainPath, "main source")
		check(item.Files.StarterPath, "starter code")
		check(item.Files.SolutionPath, "solution")
		check(item.Files.TestPath, "test file")
	}

	var issues []lessonIssue

	// 2. Check verification.mode consistency (content-shape).
	// automated => a run or test command must be declared
	// manual    => at least one manual step must be declared
	// mixed     => a run/test command or at least one manual step
	modeChecked := 0
	for _, item := range allItems {
		if item.Verification == nil {
			continue
		}
		modeChecked++
		sev := contentShapeSeverity(item)
		hasCmd := strings.TrimSpace(item.Verification.RunCommand) != "" || strings.TrimSpace(item.Verification.TestCommand) != ""
		hasSteps := len(item.Verification.ManualSteps) > 0
		switch item.Verification.Mode {
		case "automated":
			if !hasCmd {
				issues = append(issues, lessonIssue{sev, fmt.Sprintf("%s: verification.mode=automated but no run_command/test_command declared", item.ID)})
			}
		case "manual":
			if !hasSteps {
				issues = append(issues, lessonIssue{sev, fmt.Sprintf("%s: verification.mode=manual but no manual_steps declared", item.ID)})
			}
		case "mixed":
			if !hasCmd && !hasSteps {
				issues = append(issues, lessonIssue{sev, fmt.Sprintf("%s: verification.mode=mixed but neither a command nor manual_steps is declared", item.ID)})
			}
		default:
			issues = append(issues, lessonIssue{sev, fmt.Sprintf("%s: verification.mode '%s' is not 'automated', 'manual', or 'mixed'", item.ID, item.Verification.Mode)})
		}
	}

	// 3. Check README NEXT UP footers (content-shape).
	nextChecked := 0
	for _, item := range allItems {
		if item.Files == nil || item.Files.ReadmePath == "" {
			continue
		}
		if item.Type != "lesson" && item.Type != "lab" {
			continue
		}
		fullPath := resolvePath(item.Files.ReadmePath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		nextChecked++
		sev := contentShapeSeverity(item)

		section, found := nextUpSection(string(data))
		if !found {
			issues = append(issues, lessonIssue{sev, fmt.Sprintf("%s: missing '## NEXT UP' section in README", item.ID)})
			continue
		}
		if len(item.NextItemIDs) == 0 {
			continue
		}
		if !nextUpReferencesValid(section, item, fullPath, idByReadme) {
			issues = append(issues, lessonIssue{sev, fmt.Sprintf("%s: NEXT UP section does not link to a declared next item %v", item.ID, item.NextItemIDs)})
		}
	}

	reportLessonIssues(issues, modeChecked, nextChecked)

	// 4. Count machine role comments in source files (informational).
	machineRoleCount := 0
	for _, item := range allItems {
		if item.Files == nil || item.Files.MainPath == "" {
			continue
		}
		data, err := os.ReadFile(resolvePath(item.Files.MainPath))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//go:role") || strings.HasPrefix(trimmed, "/*go:role") {
				machineRoleCount++
			}
		}
	}
	fmt.Printf("  Machine role comments found: %d\n", machineRoleCount)
}

const maxLessonPrint = 40

// reportLessonIssues prints lesson issues (errors first), counting all of them
// while printing at most maxLessonPrint per severity.
func reportLessonIssues(issues []lessonIssue, modeChecked, nextChecked int) {
	errTotal, warnTotal := 0, 0
	printedErr, printedWarn := 0, 0
	for _, is := range issues {
		switch is.Sev {
		case sevError:
			errTotal++
			errorCount++
			if printedErr < maxLessonPrint {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", is.Msg)
				printedErr++
			}
		case sevWarning:
			warnTotal++
			warningCount++
			if printedWarn < maxLessonPrint {
				fmt.Fprintf(os.Stderr, "WARN: %s\n", is.Msg)
				printedWarn++
			}
		}
	}
	fmt.Printf("  Checked verification.mode on %d items and NEXT UP on %d lessons (%d errors, %d advisory).\n",
		modeChecked, nextChecked, errTotal, warnTotal)
	if warnTotal > maxLessonPrint {
		fmt.Printf("  (advisory output truncated at %d; %d total)\n", maxLessonPrint, warnTotal)
	}
	fmt.Println("  All lesson checks passed.")
}

// nextUpSection returns the text of the NEXT UP section (from the heading to the
// next H2 or EOF) and whether the heading was found.
func nextUpSection(content string) (string, bool) {
	lines := strings.Split(content, "\n")
	start := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "## ") {
			continue
		}
		if headingMatches(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), "NEXT UP") {
			start = i
			break
		}
	}
	if start < 0 {
		return "", false
	}
	var sb strings.Builder
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			break
		}
		sb.WriteString(lines[i])
		sb.WriteString("\n")
	}
	return sb.String(), true
}

// nextUpReferencesValid reports whether the NEXT UP section links to (or names)
// one of the item's declared next items.
func nextUpReferencesValid(section string, item Item, readmePath string, idByReadme map[string]string) bool {
	baseDir := filepath.Dir(readmePath)

	// 1. Follow markdown links and match them to a next item's README.
	for _, target := range extractMarkdownTargets(section) {
		target = strings.SplitN(target, "#", 2)[0]
		if target == "" || strings.Contains(target, "://") {
			continue
		}
		resolved := filepath.Clean(filepath.Join(baseDir, filepath.FromSlash(target)))
		if id, ok := idByReadme[resolved]; ok {
			for _, nid := range item.NextItemIDs {
				if id == nid {
					return true
				}
			}
		}
	}

	// 2. Fallback: the section names a next item id directly.
	for _, nid := range item.NextItemIDs {
		if strings.Contains(section, nid) {
			return true
		}
	}
	return false
}

// extractMarkdownTargets returns the link targets of every [label](target) in text.
func extractMarkdownTargets(text string) []string {
	var out []string
	for {
		open := strings.Index(text, "](")
		if open < 0 {
			break
		}
		rest := text[open+2:]
		close := strings.Index(rest, ")")
		if close < 0 {
			break
		}
		out = append(out, strings.TrimSpace(rest[:close]))
		text = rest[close+1:]
	}
	return out
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
