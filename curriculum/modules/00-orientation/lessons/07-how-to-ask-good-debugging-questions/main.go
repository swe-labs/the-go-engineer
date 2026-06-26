package main

import (
	"fmt"
	"strings"
)

type DebugQuestion struct {
	Expected    string
	Actual      string
	Command     string
	Error       string
	Environment string
}

func scoreQuestion(q DebugQuestion) int {
	score := 0
	if strings.TrimSpace(q.Expected) != "" {
		score++
	}
	if strings.TrimSpace(q.Actual) != "" {
		score++
	}
	if strings.TrimSpace(q.Command) != "" {
		score++
	}
	if strings.TrimSpace(q.Error) != "" {
		score++
	}
	if strings.TrimSpace(q.Environment) != "" {
		score++
	}
	return score
}

func main() {
	vague := DebugQuestion{
		Expected:    "",
		Actual:      "it doesn't work",
		Command:     "",
		Error:       "some error",
		Environment: "",
	}

	good := DebugQuestion{
		Expected:    "the server should return 200 OK on GET /health",
		Actual:      "the server returns 503 Service Unavailable",
		Command:     "curl -v http://localhost:8080/health",
		Error:       "HTTP/1.1 503 Service Unavailable",
		Environment: "Windows 11, Go 1.25, running locally with go run",
	}

	fmt.Println("Debugging Question Quality Score")
	fmt.Println("=================================")
	fmt.Println()
	fmt.Println("Vague question score:", scoreQuestion(vague), "/ 5")
	fmt.Println("  - Expected:      ", formatField(vague.Expected))
	fmt.Println("  - Actual:        ", formatField(vague.Actual))
	fmt.Println("  - Command:       ", formatField(vague.Command))
	fmt.Println("  - Error:         ", formatField(vague.Error))
	fmt.Println("  - Environment:   ", formatField(vague.Environment))
	fmt.Println()
	fmt.Println("High-quality question score:", scoreQuestion(good), "/ 5")
	fmt.Println("  - Expected:      ", formatField(good.Expected))
	fmt.Println("  - Actual:        ", formatField(good.Actual))
	fmt.Println("  - Command:       ", formatField(good.Command))
	fmt.Println("  - Error:         ", formatField(good.Error))
	fmt.Println("  - Environment:   ", formatField(good.Environment))
}

func formatField(s string) string {
	if s == "" {
		return "(empty)"
	}
	return s
}
