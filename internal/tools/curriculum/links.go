package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	markdownLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	inlineCodeRegex   = regexp.MustCompile("`[^`]*`")
)

// validateMarkdownLinks crawls all Markdown files in curriculum/ and verifies
// that 100% of internal relative hyperlinks resolve to real physical files on disk.
func validateMarkdownLinks() {
	fmt.Println("=== validate-markdown-links ===")
	ok := true

	totalLinks := 0
	brokenLinks := 0

	err := filepath.Walk(curriculumDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			errorf("cannot open markdown file %s: %v", path, err)
			ok = false
			return nil
		}
		defer file.Close()

		dir := filepath.Dir(path)
		scanner := bufio.NewScanner(file)
		lineNo := 0
		inCodeBlock := false

		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			if strings.HasPrefix(trimmed, "```") {
				inCodeBlock = !inCodeBlock
				continue
			}
			if inCodeBlock {
				continue
			}

			// Strip inline backtick code spans to avoid false positives on Go generics
			lineNoCode := inlineCodeRegex.ReplaceAllString(line, "")

			matches := markdownLinkRegex.FindAllStringSubmatch(lineNoCode, -1)
			for _, m := range matches {
				if len(m) < 3 {
					continue
				}
				target := strings.TrimSpace(m[2])

				// Ignore external, anchor-only, mailto, and virtual mock routes
				if strings.HasPrefix(target, "http://") ||
					strings.HasPrefix(target, "https://") ||
					strings.HasPrefix(target, "#") ||
					strings.HasPrefix(target, "mailto:") ||
					strings.HasPrefix(target, "file://") ||
					strings.HasPrefix(target, "/wiki/") {
					continue
				}

				// Strip anchor fragment and query params
				cleanTarget := target
				if idx := strings.IndexAny(cleanTarget, "#?"); idx != -1 {
					cleanTarget = cleanTarget[:idx]
				}
				cleanTarget = strings.TrimSpace(cleanTarget)
				if cleanTarget == "" {
					continue
				}

				totalLinks++

				resolved := filepath.Clean(filepath.Join(dir, cleanTarget))
				if _, err := os.Stat(resolved); err != nil {
					relSrc, _ := filepath.Rel(curriculumDir, path)
					errorf("%s:%d -> '%s' (target does not exist on disk: %s)", relSrc, lineNo, target, resolved)
					brokenLinks++
					ok = false
				}
			}
		}
		return nil
	})

	if err != nil {
		errorf("error walking markdown files in %s: %v", curriculumDir, err)
		ok = false
	}

	fmt.Printf("  Verified %d relative markdown links across curriculum (%d broken).\n", totalLinks, brokenLinks)
	if ok && brokenLinks == 0 {
		fmt.Println("  All physical markdown links resolved successfully.")
	} else {
		fmt.Println("  Markdown link validation failed with broken links.")
	}
}
