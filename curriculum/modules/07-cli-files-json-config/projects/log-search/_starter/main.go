package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

var (
	ErrEmptyQuery   = errors.New("search query cannot be empty")
	ErrInvalidRegex = errors.New("invalid regex pattern")
	ErrFileNotFound = errors.New("log file not found")
)

type LogLevel string

const (
	LevelAll   LogLevel = "ALL"
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type SearchConfig struct {
	Query         string
	IsRegex       bool
	CaseSensitive bool
	MinLevel      LogLevel
	ContextBefore int
	ContextAfter  int
}

type MatchedLine struct {
	LineNumber int
	Content    string
	IsMatch    bool
}

type SearchResult struct {
	TotalScanned int
	TotalMatches int
	Lines        []MatchedLine
}

type LogSearcher struct{}

func NewLogSearcher() *LogSearcher {
	return &LogSearcher{}
}

func levelPriority(lvl string) int {
	switch strings.ToUpper(lvl) {
	case "ERROR", "FATAL":
		return 4
	case "WARN", "WARNING":
		return 3
	case "INFO":
		return 2
	case "DEBUG", "TRACE":
		return 1
	default:
		return 0
	}
}

// SearchStream reads from an io.Reader line-by-line, matching against the SearchConfig.
func (s *LogSearcher) SearchStream(ctx context.Context, r io.Reader, cfg SearchConfig) (SearchResult, error) {
	// TODO: Compile query string or regex pattern (handling CaseSensitive flag)
	// TODO: Create bufio.Scanner to scan reader line by line (O(1) memory)
	// TODO: Use a ring buffer of capacity cfg.ContextBefore to hold preceding context lines
	// TODO: On match: flush before-ring, record match, set countdown for cfg.ContextAfter
	// TODO: Return SearchResult with TotalScanned, TotalMatches, and Lines
	return SearchResult{}, errors.New("not implemented")
}

func extractLogLevel(line string) string {
	// TODO: Detect [LEVEL] or LEVEL= in log line
	return ""
}

func main() {
	filePath := flag.String("file", "", "Log file to inspect")
	query := flag.String("query", "", "Search pattern")
	isRegex := flag.Bool("regex", false, "Interpret query as regular expression")
	caseSensitive := flag.Bool("case-sensitive", false, "Enable case sensitive matching")
	level := flag.String("level", "ALL", "Minimum log level (DEBUG, INFO, WARN, ERROR)")
	before := flag.Int("B", 0, "Context lines before match")
	after := flag.Int("A", 0, "Context lines after match")
	flag.Parse()

	var reader io.Reader
	if *filePath != "" {
		f, err := os.Open(*filePath)
		if err != nil {
			slog.Error("cannot open file", "error", err)
			os.Exit(1)
		}
		defer f.Close()
		reader = f
	} else {
		reader = os.Stdin
	}

	searcher := NewLogSearcher()
	res, err := searcher.SearchStream(context.Background(), reader, SearchConfig{
		Query:         *query,
		IsRegex:       *isRegex,
		CaseSensitive: *caseSensitive,
		MinLevel:      LogLevel(*level),
		ContextBefore: *before,
		ContextAfter:  *after,
	})
	if err != nil {
		slog.Error("search error", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Matched %d out of %d scanned lines:\n", res.TotalMatches, res.TotalScanned)
	for _, line := range res.Lines {
		prefix := "  "
		if line.IsMatch {
			prefix = "> "
		}
		fmt.Printf("%s[%4d] %s\n", prefix, line.LineNumber, line.Content)
	}
}
