package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// isEmojiCharacter tests if a rune falls into Unicode emoji or decorative symbol ranges.
func isEmojiCharacter(r rune) bool {
	// Unicode Emoji presentation and symbols (includes 0x1F449, 0x1F680, 0x1F4A1, 0x1F3AF)
	if r >= 0x1F300 && r <= 0x1FAFF {
		return true
	}
	// Miscellaneous symbols and dingbats (0x2600-0x27BF): flag every rune. Arrows
	// (U+2190-0x21FF, e.g. the -> used throughout prose) live outside this range and
	// stay allowed. The previous implementation only flagged six hardcoded dingbats,
	// which let U+270F (pencil) through a module README.
	if r >= 0x2600 && r <= 0x27BF {
		return true
	}
	return false
}

// validateNoEmojis scans all Markdown files in curriculum/ to enforce the strict
// No-Emojis Invariant across technical prose, headings, tables, and lists.
// Code blocks (fenced or inline) teaching Unicode runes or CLI indicators are excluded.
func validateNoEmojis() {
	fmt.Println("=== validate-no-emojis ===")
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

			// Strip inline backticks so code examples of UTF-8 runes are allowed
			lineNoInline := inlineCodeRegex.ReplaceAllString(line, "")

			for i := 0; i < len(lineNoInline); {
				r, size := utf8.DecodeRuneInString(lineNoInline[i:])
				if isEmojiCharacter(r) {
					relPath, _ := filepath.Rel(curriculumDir, path)
					errorf("%s:%d contains forbidden emoji '%c' (U+%04X): %s", relPath, lineNo, r, r, trimmed)
					violations++
					ok = false
					break
				}
				i += size
			}
		}
		if scanErr := scanner.Err(); scanErr != nil {
			relPath, _ := filepath.Rel(curriculumDir, path)
			errorf("%s: scanner error while reading emojis: %v", relPath, scanErr)
			ok = false
		}
		return nil
	})

	if err != nil {
		errorf("error walking markdown files for emojis: %v", err)
		ok = false
	}

	fmt.Printf("  Scanned %d markdown files for emojis (%d violations).\n", filesChecked, violations)
	if ok && violations == 0 {
		fmt.Println("  Strict No-Emojis Invariant verified across all prose.")
	} else {
		fmt.Println("  No-Emojis validation failed with violations.")
	}
}
