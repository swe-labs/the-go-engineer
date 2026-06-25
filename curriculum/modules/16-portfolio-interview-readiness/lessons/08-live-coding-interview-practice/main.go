package main

import (
	"fmt"
	"regexp"
	"strings"
)

type TestCase struct {
	Input    string
	Expected string
}

type Problem struct {
	Title       string
	Description string
	Examples    []string
	Constraints []string
	TestCases   []TestCase
}

func ParseProblem(text string) Problem {
	p := Problem{Description: text}

	titleRe := regexp.MustCompile(`(?i)(?:problem|title):\s*(.+?)(?:\n|$)`)
	if m := titleRe.FindStringSubmatch(text); len(m) > 1 {
		p.Title = strings.TrimSpace(m[1])
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), "example") {
			p.Examples = append(p.Examples, trimmed)
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "constraint") || strings.HasPrefix(strings.ToLower(trimmed), "constraints") {
			p.Constraints = append(p.Constraints, trimmed)
		}
	}

	if p.Title == "" {
		p.Title = "Unnamed Problem"
	}

	return p
}

type TestCaseGenerator struct{}

func NewTestCaseGenerator() *TestCaseGenerator {
	return &TestCaseGenerator{}
}

func (g *TestCaseGenerator) GenerateSimpleCases(p Problem) []TestCase {
	cases := []TestCase{}

	desc := strings.ToLower(p.Description)

	if strings.Contains(desc, "reverse") || strings.Contains(desc, "palindrome") {
		cases = append(cases,
			TestCase{Input: "hello", Expected: "olleh"},
			TestCase{Input: "racecar", Expected: "racecar"},
			TestCase{Input: "a", Expected: "a"},
			TestCase{Input: "", Expected: ""},
		)
	}

	if strings.Contains(desc, "sort") || strings.Contains(desc, "merge") {
		cases = append(cases,
			TestCase{Input: "[3,1,2]", Expected: "[1,2,3]"},
			TestCase{Input: "[1]", Expected: "[1]"},
			TestCase{Input: "[]", Expected: "[]"},
		)
	}

	if strings.Contains(desc, "sum") || strings.Contains(desc, "add") || strings.Contains(desc, "two sum") {
		cases = append(cases,
			TestCase{Input: "[2,7,11,15], target=9", Expected: "[0,1]"},
			TestCase{Input: "[3,2,4], target=6", Expected: "[1,2]"},
			TestCase{Input: "[3,3], target=6", Expected: "[0,1]"},
		)
	}

	if strings.Contains(desc, "fib") || strings.Contains(desc, "fibonacci") {
		cases = append(cases,
			TestCase{Input: "n=0", Expected: "0"},
			TestCase{Input: "n=1", Expected: "1"},
			TestCase{Input: "n=5", Expected: "5"},
			TestCase{Input: "n=10", Expected: "55"},
		)
	}

	if strings.Contains(desc, "fizz") || strings.Contains(desc, "buzz") || strings.Contains(desc, "fizzbuzz") {
		cases = append(cases,
			TestCase{Input: "n=1", Expected: "1"},
			TestCase{Input: "n=3", Expected: "Fizz"},
			TestCase{Input: "n=5", Expected: "Buzz"},
			TestCase{Input: "n=15", Expected: "FizzBuzz"},
		)
	}

	if strings.Contains(desc, "anagram") {
		cases = append(cases,
			TestCase{Input: `s="listen", t="silent"`, Expected: "true"},
			TestCase{Input: `s="hello", t="world"`, Expected: "false"},
			TestCase{Input: `s="", t=""`, Expected: "true"},
		)
	}

	if len(cases) == 0 {
		cases = append(cases,
			TestCase{Input: "example input 1", Expected: "expected output 1"},
			TestCase{Input: "example input 2", Expected: "expected output 2"},
		)
	}

	return cases
}

func (p Problem) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Problem: %s\n", p.Title))
	b.WriteString(fmt.Sprintf("Description: %s\n", p.Description))

	if len(p.Examples) > 0 {
		b.WriteString("\nExamples:\n")
		for _, e := range p.Examples {
			b.WriteString(fmt.Sprintf("  %s\n", e))
		}
	}
	if len(p.Constraints) > 0 {
		b.WriteString("\nConstraints:\n")
		for _, c := range p.Constraints {
			b.WriteString(fmt.Sprintf("  %s\n", c))
		}
	}

	return b.String()
}

func (p Problem) TestSummary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Generated %d test cases:\n\n", len(p.TestCases)))
	for i, tc := range p.TestCases {
		b.WriteString(fmt.Sprintf("Case %d:\n", i+1))
		b.WriteString(fmt.Sprintf("  Input:    %s\n", tc.Input))
		b.WriteString(fmt.Sprintf("  Expected: %s\n", tc.Expected))
	}
	return b.String()
}

func main() {
	problemText := `Title: Two Sum
Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.
Example 1: Input: nums = [2,7,11,15], target = 9, Output: [0,1]
Example 2: Input: nums = [3,2,4], target = 6, Output: [1,2]
Constraints: 2 <= nums.length <= 10^4, -10^9 <= nums[i] <= 10^9`

	problem := ParseProblem(problemText)
	fmt.Print(problem.Summary())

	generator := NewTestCaseGenerator()
	problem.TestCases = generator.GenerateSimpleCases(problem)
	fmt.Print(problem.TestSummary())
}
