package main

import (
	"strings"
	"testing"
)

func TestParseProblem_Title(t *testing.T) {
	text := "Title: Reverse a String\nGiven a string, reverse it."
	p := ParseProblem(text)
	if p.Title != "Reverse a String" {
		t.Errorf("expected title 'Reverse a String', got '%s'", p.Title)
	}
}

func TestParseProblem_FallbackTitle(t *testing.T) {
	text := "Just some description without a title"
	p := ParseProblem(text)
	if p.Title != "Unnamed Problem" {
		t.Errorf("expected fallback title, got '%s'", p.Title)
	}
}

func TestParseProblem_ExamplesAndConstraints(t *testing.T) {
	text := `Title: Test
Example 1: foo
Constraints: bar`
	p := ParseProblem(text)
	if len(p.Examples) != 1 {
		t.Errorf("expected 1 example, got %d", len(p.Examples))
	}
	if len(p.Constraints) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(p.Constraints))
	}
}

func TestGenerateSimpleCases_Reverse(t *testing.T) {
	p := Problem{Description: "Write a function to reverse a string"}
	gen := NewTestCaseGenerator()
	cases := gen.GenerateSimpleCases(p)
	if len(cases) < 3 {
		t.Errorf("expected at least 3 test cases for reverse, got %d", len(cases))
	}
}

func TestGenerateSimpleCases_TwoSum(t *testing.T) {
	p := Problem{Description: "Given an array of integers, find two numbers that add up to target"}
	gen := NewTestCaseGenerator()
	cases := gen.GenerateSimpleCases(p)
	if len(cases) < 2 {
		t.Errorf("expected at least 2 test cases for two sum, got %d", len(cases))
	}
	for _, tc := range cases {
		if !strings.Contains(tc.Expected, "[") {
			t.Errorf("expected array output format, got '%s'", tc.Expected)
		}
	}
}

func TestGenerateSimpleCases_Anagram(t *testing.T) {
	p := Problem{Description: "Check if two strings are anagrams of each other"}
	gen := NewTestCaseGenerator()
	cases := gen.GenerateSimpleCases(p)
	if len(cases) < 2 {
		t.Errorf("expected at least 2 test cases for anagram, got %d", len(cases))
	}
}

func TestGenerateSimpleCases_Fibonacci(t *testing.T) {
	p := Problem{Description: "Compute the nth Fibonacci number"}
	gen := NewTestCaseGenerator()
	cases := gen.GenerateSimpleCases(p)
	if len(cases) < 3 {
		t.Errorf("expected at least 3 test cases for fibonacci, got %d", len(cases))
	}
}

func TestGenerateSimpleCases_Fallback(t *testing.T) {
	p := Problem{Description: "Some unique problem that doesn't match known patterns"}
	gen := NewTestCaseGenerator()
	cases := gen.GenerateSimpleCases(p)
	if len(cases) != 2 {
		t.Errorf("expected exactly 2 fallback cases, got %d", len(cases))
	}
}
