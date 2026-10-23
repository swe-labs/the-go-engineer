package main

import (
	"errors"
	"strings"
	"testing"
)

const sampleLogs = `2026-10-23T10:00:00Z [DEBUG] Starting background worker
2026-10-23T10:00:01Z [INFO] HTTP server listening on :8080
2026-10-23T10:00:02Z [INFO] User login initiated: alice@example.com
2026-10-23T10:00:03Z [WARN] Database pool connection spike (active: 18)
2026-10-23T10:00:04Z [ERROR] Payment gateway timeout for order #45821
2026-10-23T10:00:05Z [INFO] Retry payment queued for order #45821
2026-10-23T10:00:06Z [DEBUG] Garbage collector invoked`

func TestSearchStream_Substring(t *testing.T) {
	searcher := NewLogSearcher()

	// Case-insensitive
	res, err := searcher.SearchStream(t.Context(), strings.NewReader(sampleLogs), SearchConfig{
		Query:         "PAYMENT",
		CaseSensitive: false,
		MinLevel:      LevelAll,
	})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if res.TotalMatches != 2 {
		t.Fatalf("expected 2 matches, got %d", res.TotalMatches)
	}

	// Case-sensitive (uppercase PAYMENT should yield 0 matches since log has Payment)
	resSens, err := searcher.SearchStream(t.Context(), strings.NewReader(sampleLogs), SearchConfig{
		Query:         "PAYMENT",
		CaseSensitive: true,
		MinLevel:      LevelAll,
	})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if resSens.TotalMatches != 0 {
		t.Fatalf("expected 0 matches for case-sensitive, got %d", resSens.TotalMatches)
	}
}

func TestSearchStream_Regex(t *testing.T) {
	searcher := NewLogSearcher()

	// Matches order IDs (#\d+)
	res, err := searcher.SearchStream(t.Context(), strings.NewReader(sampleLogs), SearchConfig{
		Query:    `order #\d+`,
		IsRegex:  true,
		MinLevel: LevelAll,
	})
	if err != nil {
		t.Fatalf("regex search failed: %v", err)
	}

	if res.TotalMatches != 2 {
		t.Fatalf("expected 2 matches for regex order #\\d+, got %d", res.TotalMatches)
	}

	// Invalid regex pattern
	_, errBad := searcher.SearchStream(t.Context(), strings.NewReader(sampleLogs), SearchConfig{
		Query:   `[invalid(regex`,
		IsRegex: true,
	})
	if !errors.Is(errBad, ErrInvalidRegex) {
		t.Fatalf("expected ErrInvalidRegex, got %v", errBad)
	}
}

func TestSearchStream_MinLevelFiltering(t *testing.T) {
	searcher := NewLogSearcher()

	// Filter by LevelWarn: should match lines with [WARN] and [ERROR]
	res, err := searcher.SearchStream(t.Context(), strings.NewReader(sampleLogs), SearchConfig{
		MinLevel: LevelWarn,
	})
	if err != nil {
		t.Fatalf("level search failed: %v", err)
	}

	if res.TotalMatches != 2 {
		t.Fatalf("expected 2 warnings/errors, got %d", res.TotalMatches)
	}
}

func TestSearchStream_ContextLines(t *testing.T) {
	searcher := NewLogSearcher()

	// Search for "[ERROR]" with 1 before and 1 after line of context
	res, err := searcher.SearchStream(t.Context(), strings.NewReader(sampleLogs), SearchConfig{
		Query:         "[ERROR]",
		ContextBefore: 1,
		ContextAfter:  1,
	})
	if err != nil {
		t.Fatalf("context search failed: %v", err)
	}

	if res.TotalMatches != 1 {
		t.Fatalf("expected 1 match, got %d", res.TotalMatches)
	}

	// Total lines should be 3: line 4 (before), line 5 (match), line 6 (after)
	if len(res.Lines) != 3 {
		t.Fatalf("expected 3 total lines (1 before + 1 match + 1 after), got %d", len(res.Lines))
	}

	if res.Lines[0].LineNumber != 4 || res.Lines[0].IsMatch {
		t.Fatalf("unexpected before context line: %+v", res.Lines[0])
	}
	if res.Lines[1].LineNumber != 5 || !res.Lines[1].IsMatch {
		t.Fatalf("unexpected match line: %+v", res.Lines[1])
	}
	if res.Lines[2].LineNumber != 6 || res.Lines[2].IsMatch {
		t.Fatalf("unexpected after context line: %+v", res.Lines[2])
	}
}
