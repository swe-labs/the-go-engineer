package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// validateInlinePaths verifies that repository paths cited inside inline code spans in
// curriculum markdown resolve to real files or directories on disk. The relative
// markdown-link crawler only inspects [text](target) targets; prose may also cite paths
// inside backticks (for example, practice tasks that direct learners to open
// `curriculum/modules/…/README.md`). Such citations previously escaped every gate — a
// lesson once sent learners to two lesson files that did not exist.
//
// A span is treated as a repository path when, after stripping an optional leading
// "curriculum/", it starts with one of the top-level curriculum directories and ends
// either with a known file extension or with "/" (a directory citation).
func validateInlinePaths() {
	fmt.Println("=== validate-inline-paths ===")
	ok := true
	violations := 0
	filesChecked := 0

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

		filesChecked++
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 64*1024)
		scanner.Buffer(buf, 1024*1024)
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

			for _, span := range inlineCodeSpans(line) {
				rel, valid := curriculumRelativePath(span)
				if !valid {
					continue
				}
				if !resolvesUnderAncestor(curriculumDir, path, rel) {
					relPath, _ := filepath.Rel(curriculumDir, path)
					errorf("%s:%d cites a repository path that does not exist on disk: `%s`", relPath, lineNo, span)
					violations++
					ok = false
				}
			}
		}
		if scanErr := scanner.Err(); scanErr != nil {
			relPath, _ := filepath.Rel(curriculumDir, path)
			errorf("%s: scanner error while reading inline paths: %v", relPath, scanErr)
			ok = false
		}
		return nil
	})

	if err != nil {
		errorf("error walking markdown files for inline paths: %v", err)
		ok = false
	}

	fmt.Printf("  Scanned %d markdown files for inline repository-path citations (%d violations).\n", filesChecked, violations)
	if ok && violations == 0 {
		fmt.Println("  All cited repository paths resolve to real files or directories.")
	} else {
		fmt.Println("  Inline-path validation failed with violations.")
	}
}

// inlineCodeSpans returns the contents of single-backtick inline code spans on a line,
// skipping fenced blocks (handled by the caller).
func inlineCodeSpans(line string) []string {
	var spans []string
	for i := 0; i < len(line); {
		if line[i] != '`' {
			i++
			continue
		}
		end := strings.IndexByte(line[i+1:], '`')
		if end < 0 {
			break
		}
		spans = append(spans, line[i+1:i+1+end])
		i += end + 2
	}
	return spans
}

// resolvesUnderAncestor reports whether rel resolves to an existing file when
// joined with the citing file's directory or any of its ancestors up to
// curriculumDir. Module READMEs cite sibling paths relative to the module root
// (e.g. `assessments/checkpoint/README.md` from modules/15/README.md), while
// lessons cite curriculum-prefixed paths — both must resolve.
func resolvesUnderAncestor(curriculumRoot, citingFile, rel string) bool {
	dir, err := filepath.Rel(curriculumRoot, filepath.Dir(citingFile))
	if err != nil {
		dir = "."
	}
	candidate := filepath.FromSlash(rel)
	for {
		if _, err := os.Stat(filepath.Join(curriculumRoot, dir, candidate)); err == nil {
			return true
		}
		if dir == "." {
			return false
		}
		dir = filepath.Dir(dir)
	}
}

// curriculumRelativePath reports whether an inline code span looks like a repository
// path citation, and if so returns the slash path relative to curriculumDir.
func curriculumRelativePath(span string) (string, bool) {
	span = strings.TrimSpace(span)
	span = strings.TrimPrefix(span, "curriculum/")
	span = strings.TrimSuffix(span, "/")
	if span == "" || strings.ContainsAny(span, "{}") {
		return "", false
	}

	topDirs := []string{"modules/", "assessments/", "projects/", "companion-frontend/", "electives/"}
	topLevel := false
	for _, d := range topDirs {
		if strings.HasPrefix(span, d) {
			topLevel = true
			break
		}
	}
	if !topLevel {
		return "", false
	}

	if strings.ContainsAny(span, "\\\n\t") {
		return "", false
	}
	for _, part := range strings.Split(span, "/") {
		if part == "" || part == "." || part == ".." {
			return "", false
		}
	}

	if !strings.Contains(span, ".") {
		// Directory citation without extension (e.g. .../projects/tooling-failure-lab).
		return span, true
	}
	exts := []string{".md", ".json", ".svg", ".txt", ".sh", ".ps1", ".yaml", ".yml", ".go"}
	for _, ext := range exts {
		if strings.HasSuffix(span, ext) {
			return span, true
		}
	}
	return "", false
}
