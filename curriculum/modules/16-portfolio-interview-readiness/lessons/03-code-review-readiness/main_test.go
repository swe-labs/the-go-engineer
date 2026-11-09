package main

import (
	"strings"
	"testing"
)

func TestRunCodeReview_GoodCode(t *testing.T) {
	code := `package main
import (
	"context"
	"sync"
)
type Store struct { mu sync.RWMutex; data map[string]string }
func NewStore() *Store { return &Store{data: make(map[string]string)} }
func (s *Store) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok { return "", fmt.Errorf("not found") }
	return val, nil
}`
	result := RunCodeReview(code)
	if result.Score < 5 {
		t.Errorf("expected good code to score >= 5, got %d/%d", result.Score, result.Total)
	}
}

func TestRunCodeReview_ErrorHandlingCheck(t *testing.T) {
	codeNoErr := `package main
func foo() { println("hello") }`
	result := RunCodeReview(codeNoErr)
	for _, c := range result.Checks {
		if c.Category == Correctness && strings.Contains(c.Description, "Error handling") {
			if c.Passed {
				t.Errorf("expected error handling check to fail for code without error handling")
			}
		}
	}
}

func TestRunCodeReview_AllCategoriesPresent(t *testing.T) {
	code := `package main
func main() {}`
	result := RunCodeReview(code)
	categories := make(map[ReviewCategory]bool)
	for _, c := range result.Checks {
		categories[c.Category] = true
	}
	expected := []ReviewCategory{Correctness, Style, Performance, Security, Reliability, Maintainability}
	for _, cat := range expected {
		if !categories[cat] {
			t.Errorf("missing category %v", cat)
		}
	}
}

func TestRunCodeReview_DeferDetection(t *testing.T) {
	code := `package main
func foo() { defer cleanup() }`
	result := RunCodeReview(code)
	for _, c := range result.Checks {
		if c.Category == Maintainability && strings.Contains(c.Description, "Resource cleanup") {
			if !c.Passed {
				t.Errorf("expected defer usage to be detected")
			}
		}
	}
}

func TestRunCodeReview_ContextDetection(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"has context", `func foo(ctx context.Context) {}`, true},
		{"no context", `func foo() {}`, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code := "package main\n" + tc.content
			result := RunCodeReview(code)
			var found bool
			for _, c := range result.Checks {
				if c.Category == Correctness && strings.Contains(c.Description, "Context") {
					found = c.Passed
					break
				}
			}
			if found != tc.expected {
				t.Errorf("context detection: got %v, want %v", found, tc.expected)
			}
		})
	}
}
