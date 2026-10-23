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
	ErrEmptyQuery    = errors.New("search query cannot be empty")
	ErrInvalidRegex  = errors.New("invalid regex pattern")
	ErrFileNotFound  = errors.New("log file not found")
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
	IsMatch    bool // true if target match, false if context line
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
	if strings.TrimSpace(cfg.Query) == "" && cfg.MinLevel == LevelAll {
		return SearchResult{}, ErrEmptyQuery
	}

	var matchFn func(line string) bool
	if cfg.Query != "" {
		if cfg.IsRegex {
			flags := ""
			if !cfg.CaseSensitive {
				flags = "(?i)"
			}
			re, err := regexp.Compile(flags + cfg.Query)
			if err != nil {
				return SearchResult{}, fmt.Errorf("%w: %v", ErrInvalidRegex, err)
			}
			matchFn = func(line string) bool {
				return re.MatchString(line)
			}
		} else {
			target := cfg.Query
			if !cfg.CaseSensitive {
				target = strings.ToLower(target)
			}
			matchFn = func(line string) bool {
				cmp := line
				if !cfg.CaseSensitive {
					cmp = strings.ToLower(cmp)
				}
				return strings.Contains(cmp, target)
			}
		}
	} else {
		matchFn = func(line string) bool { return true }
	}

	targetMinPriority := levelPriority(string(cfg.MinLevel))

	scanner := bufio.NewScanner(r)
	// Support long log lines up to 1MB
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var result SearchResult
	lineNum := 0

	// Ring buffer for ContextBefore lines
	beforeCapacity := cfg.ContextBefore
	if beforeCapacity < 0 {
		beforeCapacity = 0
	}
	beforeRing := make([]MatchedLine, 0, beforeCapacity)

	remainingAfterContext := 0

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		lineNum++
		result.TotalScanned++
		line := scanner.Text()

		// Extract level if specified
		passesLevel := true
		if targetMinPriority > 0 {
			lvl := extractLogLevel(line)
			if levelPriority(lvl) < targetMinPriority {
				passesLevel = false
			}
		}

		isMatch := passesLevel && matchFn(line)

		if isMatch {
			result.TotalMatches++

			// Flush accumulated context before
			for _, b := range beforeRing {
				result.Lines = append(result.Lines, b)
			}
			beforeRing = beforeRing[:0]

			result.Lines = append(result.Lines, MatchedLine{
				LineNumber: lineNum,
				Content:    line,
				IsMatch:    true,
			})
			remainingAfterContext = cfg.ContextAfter
		} else if remainingAfterContext > 0 {
			result.Lines = append(result.Lines, MatchedLine{
				LineNumber: lineNum,
				Content:    line,
				IsMatch:    false,
			})
			remainingAfterContext--
		} else if beforeCapacity > 0 {
			if len(beforeRing) >= beforeCapacity {
				beforeRing = beforeRing[1:]
			}
			beforeRing = append(beforeRing, MatchedLine{
				LineNumber: lineNum,
				Content:    line,
				IsMatch:    false,
			})
		}
	}

	return result, scanner.Err()
}

func extractLogLevel(line string) string {
	upper := strings.ToUpper(line)
	for _, lvl := range []string{"ERROR", "WARN", "INFO", "DEBUG", "FATAL"} {
		if strings.Contains(upper, "["+lvl+"]") || strings.Contains(upper, "LEVEL="+lvl) || strings.Contains(upper, "\""+lvl+"\"") {
			return lvl
		}
	}
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
