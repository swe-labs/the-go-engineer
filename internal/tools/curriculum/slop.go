package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// slopRule defines a regex pattern and explanation for banned AI/filler prose.
type slopRule struct {
	Name    string
	Pattern *regexp.Regexp
	Reason  string
}

var slopRules = []slopRule{
	{
		Name:    "delve-into",
		Pattern: regexp.MustCompile(`(?i)\bdelve\s+(?:into|deeper|in)\b`),
		Reason:  "banned AI-filler transition; replace with direct action verb (explain, inspect, trace)",
	},
	{
		Name:    "seamless",
		Pattern: regexp.MustCompile(`(?i)\bseamless(?:ly)?\b`),
		Reason:  "vague superlative; state the exact mechanism or invariant that operates without failure",
	},
	{
		Name:    "cutting-edge",
		Pattern: regexp.MustCompile(`(?i)\bcutting-edge\b`),
		Reason:  "unearned marketing superlative; state toolchain version, date, or measured capability",
	},
	{
		Name:    "state-of-the-art",
		Pattern: regexp.MustCompile(`(?i)\bstate-of-the-art\b`),
		Reason:  "unearned marketing superlative; cite specific benchmarks or standards",
	},
	{
		Name:    "game-changer",
		Pattern: regexp.MustCompile(`(?i)\bgame-changer\b`),
		Reason:  "marketing hype; state the measured engineering improvement",
	},
	{
		Name:    "revolutionize",
		Pattern: regexp.MustCompile(`(?i)\brevolutionize(?:s|d)?\b`),
		Reason:  "marketing hype; describe the concrete architectural change",
	},
	{
		Name:    "unlock-potential",
		Pattern: regexp.MustCompile(`(?i)\bunlock(?:ing|s|ed)?\s+the\s+full\b`),
		Reason:  "vague promotional cliché; state the specific capability the learner gains",
	},
	{
		Name:    "important-to-note",
		Pattern: regexp.MustCompile(`(?i)\b(?:it is|it's)\s+important\s+to\s+note\b`),
		Reason:  "empty throat-clearing transition; state the consequence directly",
	},
	{
		Name:    "worth-noting",
		Pattern: regexp.MustCompile(`(?i)\b(?:it is|it's)\s+worth\s+noting\b`),
		Reason:  "empty throat-clearing transition; state the consequence directly",
	},
	{
		Name:    "needless-to-say",
		Pattern: regexp.MustCompile(`(?i)\bneedless\s+to\s+say\b`),
		Reason:  "unnecessary filler; delete the sentence or make the consequence explicit",
	},
	{
		Name:    "in-simple-terms",
		Pattern: regexp.MustCompile(`(?i)\bin\s+simple\s+terms\b`),
		Reason:  "banned by readme.contracts.json Rule 15; explain the physical mechanism directly",
	},
	{
		Name:    "not-only-but-also",
		Pattern: regexp.MustCompile(`(?i)\bnot only\b.*?\bbut also\b`),
		Reason:  "convoluted construction; replace with 'and' or split into two sentences",
	},
}

// validateAntiSlop scans all Markdown prose across the curriculum for banned
// AI-generated filler, lexical tells, and unearned superlatives.
func validateAntiSlop() {
	fmt.Println("=== validate-anti-slop ===")
	ok := true
	checkedFiles := 0
	violationCount := 0

	err := filepath.WalkDir(curriculumDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		// Skip authoring guides, checklists, and rubrics that document the rules themselves
		baseName := strings.ToLower(d.Name())
		if strings.Contains(baseName, "authoring") ||
			strings.Contains(baseName, "rubric") ||
			strings.Contains(baseName, "checklist") {
			return nil
		}

		checkedFiles++
		data, err := os.ReadFile(path)
		if err != nil {
			errorf("cannot read markdown file %s: %v", path, err)
			ok = false
			return nil
		}

		relPath, _ := filepath.Rel(curriculumDir, path)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		lineNo := 0
		inCode := false

		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			if strings.HasPrefix(trimmed, "```") {
				inCode = !inCode
				continue
			}
			// Skip code blocks and historical pioneer blockquotes
			if inCode || strings.HasPrefix(trimmed, ">") {
				continue
			}

			for _, rule := range slopRules {
				if match := rule.Pattern.FindString(line); match != "" {
					violationCount++
					ok = false
					errorf("%s:%d [%s] contains banned slop phrase '%s' (%s)",
						relPath, lineNo, rule.Name, match, rule.Reason)
				}
			}
		}

		return nil
	})

	if err != nil {
		errorf("failed to walk curriculum directory: %v", err)
		ok = false
	}

	fmt.Printf("  Scanned %d markdown files for anti-slop lexical compliance (%d violations).\n", checkedFiles, violationCount)
	if ok && violationCount == 0 {
		fmt.Println("  Strict Anti-Slop Prose Invariant verified across all curriculum documents.")
	} else {
		fmt.Println("  Anti-slop validation failed with violations reported above.")
	}
}
