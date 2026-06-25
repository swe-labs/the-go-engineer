# Live coding interview practice

## Learning objective

Approach live coding interviews with a structured problem-solving framework, verbalize your thought process clearly, write testable Go code under time pressure, and handle feedback gracefully.

## Why this matters

Live coding interviews are the most stressful part of the engineering interview process for most candidates. Unlike take-home projects, there is no undo button and no time to research. The pressure of someone watching you type, combined with the expectation to talk through your reasoning, creates a cognitive load that can derail even experienced engineers. A structured framework reduces that load by automating the process — so you can focus on the problem, not the panic.

## Mental model

Think of a live coding interview as pair programming with a future teammate. The interviewer knows the solution and is evaluating how you think, not whether you get the right answer. Your job is to make your thinking visible. The interviewer cannot read your mind — they can only evaluate what you say and type. A clear verbalization of your approach is worth more than a perfect solution you arrive at silently.

## Core idea

The live coding interview follows a five-step framework:

**Step 1: Clarify the problem (2-3 minutes)**
- Restate the problem in your own words.
- Ask clarifying questions: input format, edge cases, constraints.
- Confirm example inputs and expected outputs.

**Step 2: Design the approach (3-5 minutes)**
- Describe the algorithm or data structure you plan to use.
- Discuss time and space complexity upfront.
- Ask the interviewer: "Does this approach seem reasonable before I start coding?"

**Step 3: Write the code (10-15 minutes)**
- Start with the function signature and type definitions.
- Write the main logic first, then handle edge cases.
- Use descriptive variable names and Go conventions.
- Add comments for non-obvious logic.

**Step 4: Test the code (3-5 minutes)**
- Walk through the example input manually.
- Check edge cases: empty input, single element, duplicates, negative values.
- Use the Go playground or local `go test` if available.

**Step 5: Discuss improvements (2-3 minutes)**
- What would you optimize if the input were 10x larger?
- What tradeoffs did you make?
- How would you test this in production?

**Common problem patterns for Go interviews**:

| Pattern | Example | Go-specific tools |
|---|---|---|
| String manipulation | Reverse a string, anagram check | `strings`, `unicode/utf8` |
| Array/slice operations | Two sum, merge intervals | Slice operations, `sort` |
| Linked lists | Reverse, detect cycle | Pointer manipulation |
| Trees | DFS, BFS, level order | Recursion, slices as queues |
| Dynamic programming | Fibonacci, knapsack | Memoization with maps |
| Concurrency | Fan-out/fan-in, worker pool | Goroutines, channels, `sync` |

## Under the hood

The problem statement parser extracts structured information from a natural language problem description. It uses regex-based pattern matching to identify:

- The problem title (`Title: ...`)
- Example inputs and outputs (`Example 1:`)
- Constraints (`Constraints: ...`)

The test case generator produces basic test cases by matching keywords in the description to known problem types:

- "reverse" or "palindrome" generates string reversal tests.
- "sum" or "two sum" generates index pair tests.
- "sort" or "merge" generates array sorting tests.
- "fib" or "fibonacci" generates Fibonacci sequence tests.
- "fizz" or "buzz" generates FizzBuzz tests.
- "anagram" generates string comparison tests.
- Unknown problems get generic fallback cases.

This is a simplified version of what a real problem parser would do. Production coding platforms like LeetCode use curated test case databases with thousands of edge cases.

## How Go uses it

Go is a popular choice for live coding interviews because:

- The syntax is small and easy to write without memorization.
- There is no inheritance, generics complexity (before 1.18), or metaprogramming to worry about.
- The `testing` package is built in — you can write tests during the interview.
- `go fmt` means style is never a discussion point.
- The playground (play.golang.org) is a common interview tool.

Go-specific tips for live coding:

- Use `:=` for type inference to reduce typing.
- Prefer `for range` over index-based loops.
- Use `switch` for multi-condition branching.
- Use `make([]T, 0, n)` for pre-allocation when size is known.
- Handle errors explicitly; do not use `_` for critical returns.

## Go example

The problem statement parser and test case generator reads a problem description and produces basic test cases to verify your solution.

```go
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
		if strings.HasPrefix(strings.ToLower(trimmed), "constraint") ||
			strings.HasPrefix(strings.ToLower(trimmed), "constraints") {
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

	if strings.Contains(desc, "fizz") || strings.Contains(desc, "fizzbuzz") {
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
```

## Step-by-step execution

1. `ParseProblem` receives a raw problem description string. It uses regex to extract the title and splits lines to find examples and constraints.
2. If no title is found, it falls back to "Unnamed Problem."
3. `TestCaseGenerator.GenerateSimpleCases` examines the problem description for keywords.
4. For each known pattern (reverse, sum, sort, fibonacci, fizzbuzz, anagram), it returns a set of hand-crafted test cases covering typical inputs and edge cases.
5. If no pattern is recognized, two generic test cases are returned as a starting point.
6. The parsed problem and generated test cases are printed for the candidate to use when coding and verifying their solution.

## Common mistakes

- Mistake: Starting to code immediately without clarifying the problem.
  - Why it happens: Silence feels uncomfortable, so candidates fill it with typing.
  - Fix: Say "Let me make sure I understand the problem" and restate it. Ask at least one clarifying question before typing.

- Mistake: Writing code in silence. The interviewer has no idea what you are thinking.
  - Why it happens: Coding requires concentration, and talking while thinking is hard.
  - Fix: Verbalize everything: "I am writing the function signature. Now I need to handle the edge case where the input is empty. I will use a map for O(n) lookup."

- Mistake: Ignoring edge cases until the end.
  - Why it happens: Edge cases feel like an afterthought when focused on the main logic.
  - Fix: List the edge cases at the start: "Edge cases: empty input, single element, all duplicates, no solution exists."

- Mistake: Getting stuck and staying silent.
  - Why it happens: The candidate is trying to solve the problem internally.
  - Fix: Describe your stuck state: "I am considering two approaches. Option A uses sorting which gives O(n log n). Option B uses a hash map which is O(n). Let me try option B because the problem emphasizes performance."

- Mistake: Defending a wrong approach instead of incorporating feedback.
  - Why it happens: Ego or stress makes candidates double down.
  - Fix: When the interviewer hints at a better approach, say "That is a good point. Let me think about how to incorporate that."

## Debugging walkthrough

A candidate is asked to "Write a function that checks if two strings are anagrams." The candidate writes:

```go
func isAnagram(s, t string) bool {
	return sortString(s) == sortString(t)
}
```

But does not implement `sortString`. The candidate realizes the mistake, adds the sorting implementation, and tests with the generated test cases:

```
Case 1: s="listen", t="silent" -> Expected: "true"
```

The function returns `true` correctly. Edge case generates `s="", t=""` returns `true`. The candidate then mentions the hash map approach for O(n) runtime and discusses the tradeoff: sorting is O(n log n) but simpler to implement correctly.

## Production notes

In a real interview, the problem parser and test case generator are tools you use during practice, not during the interview. Prepare by:

- Solving 20-30 problems using the five-step framework.
- Practicing aloud, even when alone. Record yourself and evaluate clarity.
- Doing mock interviews with peers to simulate pressure.
- Learning the Go standard library well enough to avoid looking up syntax.

The most important skill is not writing perfect code — it is demonstrating a repeatable problem-solving process. Interviewers evaluate: Do you understand the problem? Do you consider edge cases? Do you communicate clearly? Can you accept feedback?

## Performance implications

In a live coding interview, performance concerns focus on algorithm efficiency:

- Time complexity: O(n) vs O(n log n) vs O(n^2). Discuss before coding.
- Space complexity: O(n) vs O(1). In Go, heap allocations matter for memory-constrained systems.
- Concurrency: For problems involving multiple independent computations, goroutines can improve throughput at the cost of synchronization complexity.

The generated test cases are small and run in microseconds. During the interview, the candidate might be asked: "What happens when the input size is 10 million?" This tests whether the candidate considers scalability.

## Practice task

Extend the test case generator to support a "binary search" pattern. Add detection for "binary search", "search", "find" in the description. Generate test cases for searching in a sorted array: target found (beginning, middle, end), target not found, empty array. Add corresponding unit tests.

## Tests / verification

```bash
go test ./curriculum/modules/16-portfolio-interview-readiness/lessons/08-live-coding-interview-practice
go run ./curriculum/modules/16-portfolio-interview-readiness/lessons/08-live-coding-interview-practice
```

The tests verify title extraction, fallback title, example/constraint extraction, and all six pattern types (reverse, sum, sort, fibonacci, fizzbuzz, anagram) with correct test case counts. After the practice task, binary search test cases are verified.

## Review questions

1. What are the five steps of the live coding interview framework?
2. Why should you discuss the approach with the interviewer before writing code?
3. What should you do when you realize your initial approach is wrong?
4. Why is Go a good language choice for live coding interviews?
5. What is the difference between a test case for the happy path and a test case for an edge case?

## NEXT UP

Take-home project strategy — learn how to approach take-home coding assignments with the same rigor as production projects.
