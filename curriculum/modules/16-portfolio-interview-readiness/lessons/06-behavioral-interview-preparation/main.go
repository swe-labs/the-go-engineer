package main

import (
	"fmt"
	"strings"
)

type STARElement int

const (
	Situation STARElement = iota
	Task
	Action
	Result
)

func (s STARElement) String() string {
	switch s {
	case Situation:
		return "Situation"
	case Task:
		return "Task"
	case Action:
		return "Action"
	case Result:
		return "Result"
	default:
		return "Unknown"
	}
}

func (s STARElement) Keywords() []string {
	switch s {
	case Situation:
		return []string{"situation", "context", "background", "was working", "was part of", "environment", "scenario", "when"}
	case Task:
		return []string{"task", "goal", "objective", "needed to", "had to", "was responsible for", "challenge", "problem", "assignment"}
	case Action:
		return []string{"implemented", "designed", "built", "created", "developed", "led", "architected", "optimized", "refactored", "introduced", "migrated", "automated", "wrote", "configured"}
	case Result:
		return []string{"result", "outcome", "improved", "reduced", "increased", "saved", "achieved", "delivered", "shipped", "launched", "completed", "won", "earned", "recognized"}
	}
	return nil
}

type STARCheck struct {
	Element   STARElement
	Found     bool
	MatchedOn string
}

type STARResult struct {
	Text     string
	Checks   []STARCheck
	Score    int
	MaxScore int
}

func AnalyzeSTAR(response string) STARResult {
	result := STARResult{Text: response, MaxScore: 4}
	lower := strings.ToLower(response)

	elements := []STARElement{Situation, Task, Action, Result}
	for _, el := range elements {
		check := STARCheck{Element: el}
		for _, kw := range el.Keywords() {
			if strings.Contains(lower, kw) {
				check.Found = true
				check.MatchedOn = kw
				break
			}
		}
		result.Checks = append(result.Checks, check)
		if check.Found {
			result.Score++
		}
	}

	return result
}

func (r STARResult) Summary() string {
	var b strings.Builder
	b.WriteString("STAR Response Analysis\n")
	b.WriteString(fmt.Sprintf("Score: %d/4 elements detected\n\n", r.Score))
	for _, c := range r.Checks {
		status := "MISSING"
		if c.Found {
			status = "FOUND"
		}
		b.WriteString(fmt.Sprintf("  [%s] %s", status, c.Element))
		if c.Found {
			b.WriteString(fmt.Sprintf(" (matched: '%s')", c.MatchedOn))
		}
		b.WriteString("\n")
	}

	if r.Score < 4 {
		missing := []string{}
		for _, c := range r.Checks {
			if !c.Found {
				missing = append(missing, c.Element.String())
			}
		}
		b.WriteString(fmt.Sprintf("\nMissing elements: %s\n", strings.Join(missing, ", ")))
		b.WriteString("Tip: Use the STAR format to structure your response:\n")
		b.WriteString("  Situation: Set the context\n")
		b.WriteString("  Task: Describe the challenge or goal\n")
		b.WriteString("  Action: Explain what you did\n")
		b.WriteString("  Result: Share the outcome\n")
	}

	return b.String()
}

func main() {
	responses := []string{
		`When I was working on the payment processing team, we had a critical issue where transactions were timing out under high load. I implemented a circuit breaker pattern using a Go channel-based design that automatically degraded non-critical services. This reduced p99 latency by 60% and eliminated timeout errors during peak traffic.`,
		`I built a REST API. It was good.`,
	}

	for i, resp := range responses {
		fmt.Printf("=== Response %d ===\n%s\n\n", i+1, resp)
		result := AnalyzeSTAR(resp)
		fmt.Println(result.Summary())
		fmt.Println()
	}
}
