package main

import (
	"fmt"
	"regexp"
	"strings"
)

type ReviewCategory int

const (
	Correctness ReviewCategory = iota
	Style
	Performance
	Security
	Reliability
	Maintainability
)

func (c ReviewCategory) String() string {
	switch c {
	case Correctness:
		return "Correctness"
	case Style:
		return "Style"
	case Performance:
		return "Performance"
	case Security:
		return "Security"
	case Reliability:
		return "Reliability"
	case Maintainability:
		return "Maintainability"
	default:
		return "Unknown"
	}
}

type Check struct {
	Category    ReviewCategory
	Description string
	Passed      bool
	Detail      string
}

type ReviewResult struct {
	File    string
	Content string
	Checks  []Check
	Score   int
	Total   int
}

func NewReviewResult(file, content string) ReviewResult {
	return ReviewResult{File: file, Content: content}
}

func (rr *ReviewResult) AddCheck(cat ReviewCategory, desc string, passed bool, detail string) {
	rr.Checks = append(rr.Checks, Check{cat, desc, passed, detail})
	rr.Total++
	if passed {
		rr.Score++
	}
}

func RunCodeReview(content string) ReviewResult {
	rr := NewReviewResult("reviewed.go", content)

	hasErrorHandling := strings.Contains(content, "if err != nil")
	rr.AddCheck(Correctness, "Error handling: checks returned errors", hasErrorHandling, map[bool]string{true: "uses if err != nil pattern", false: "missing error handling"}[hasErrorHandling])

	hasContext := strings.Contains(content, "context.Background()") || strings.Contains(content, "context.WithCancel") || strings.Contains(content, "context.WithTimeout") || strings.Contains(content, "context.Context")
	rr.AddCheck(Correctness, "Context usage: proper context propagation", hasContext, map[bool]string{true: "context detected", false: "no context usage"}[hasContext])

	hasMutex := strings.Contains(content, "sync.Mutex") || strings.Contains(content, "sync.RWMutex")
	hasChan := strings.Contains(content, "chan") || strings.Contains(content, "make(chan")
	if hasMutex || hasChan {
		rr.AddCheck(Correctness, "Concurrency safety: synchronization used", true, "uses mutex or channels for concurrency")
	} else {
		rr.AddCheck(Correctness, "Concurrency safety: synchronization used", true, "no concurrent access detected")
	}

	gofmt := regexp.MustCompile(`\t`).MatchString(content)
	rr.AddCheck(Style, "Formatting: uses tabs for indentation", gofmt, map[bool]string{true: "tabs detected", false: "no tabs found, may use spaces"}[gofmt])

	hasExported := regexp.MustCompile(`^func [A-Z]`).MatchString(content)
	rr.AddCheck(Style, "Naming: exported functions use PascalCase", hasExported, map[bool]string{true: "exported functions follow convention", false: "no exported functions or naming issue"}[hasExported])

	hasGoRoutine := strings.Contains(content, "go ")
	if hasGoRoutine {
		hasWG := strings.Contains(content, "sync.WaitGroup")
		rr.AddCheck(Reliability, "Goroutine lifecycle: goroutines have synchronization", hasWG, map[bool]string{true: "uses WaitGroup for goroutine coordination", false: "goroutines without WaitGroup may leak"}[hasWG])
	} else {
		rr.AddCheck(Reliability, "Goroutine lifecycle: goroutines have synchronization", true, "no goroutines to manage")
	}

	hasInputValidation := strings.Contains(content, "len(") && (strings.Contains(content, "== 0") || strings.Contains(content, "< 0"))
	rr.AddCheck(Security, "Input validation: validates inputs before use", hasInputValidation, map[bool]string{true: "input length validation detected", false: "no input validation detected"}[hasInputValidation])

	hasSQL := strings.Contains(content, "sql.DB") || strings.Contains(content, "database/sql")
	if hasSQL {
		hasParamQuery := strings.Contains(content, "$1") || strings.Contains(content, "?") && strings.Contains(content, "Query")
		rr.AddCheck(Security, "SQL injection: uses parameterized queries", hasParamQuery, map[bool]string{true: "parameterized queries detected", false: "possible string concatenation in SQL"}[hasParamQuery])
	}

	hasAlloc := strings.Contains(content, "make(") || strings.Contains(content, "new(")
	if hasAlloc {
		hasCap := regexp.MustCompile(`make\(\[?[a-zA-Z]`).MatchString(content)
		rr.AddCheck(Performance, "Allocation: pre-allocates slices/maps with capacity", hasCap, map[bool]string{true: "pre-allocation detected", false: "no pre-allocation, may cause reallocation"}[hasCap])
	} else {
		rr.AddCheck(Performance, "Allocation: pre-allocates slices/maps with capacity", true, "no allocations to optimize")
	}

	hasDefer := strings.Contains(content, "defer ")
	rr.AddCheck(Maintainability, "Resource cleanup: uses defer for cleanup", hasDefer, map[bool]string{true: "defer usage detected", false: "no defer usage, resources may not be released"}[hasDefer])

	return rr
}

func (rr ReviewResult) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Review: %s\nScore: %d/%d (%.0f%%)\n\n", rr.File, rr.Score, rr.Total, float64(rr.Score)/float64(rr.Total)*100))
	for _, c := range rr.Checks {
		status := "PASS"
		if !c.Passed {
			status = "FAIL"
		}
		b.WriteString(fmt.Sprintf("[%s] [%s] %s\n", status, c.Category, c.Description))
		b.WriteString(fmt.Sprintf("      %s\n", c.Detail))
	}
	return b.String()
}

func main() {
	code := `package main

import (
	"context"
	"fmt"
	"sync"
)

type Server struct {
	mu     sync.Mutex
	items  map[string]string
}

func NewServer() *Server {
	return &Server{items: make(map[string]string)}
}

func (s *Server) Get(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.items[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

func main() {
	s := NewServer()
	ctx := context.Background()
	val, err := s.Get(ctx, "hello")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(val)
}
`
	result := RunCodeReview(code)
	fmt.Println(result.Summary())
}
