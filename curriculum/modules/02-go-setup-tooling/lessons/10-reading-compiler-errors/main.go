package main

import (
	"fmt"
	"strconv"
	"strings"
)

// CompilerError represents a structured Go compile-time error.
type CompilerError struct {
	File     string
	Line     int
	Col      int
	Message  string
	Category string
}

// ParseCompilerError decodes the standard Go compiler diagnostic format: file:line:col: message
func ParseCompilerError(raw string) (CompilerError, error) {
	parts := strings.SplitN(raw, ":", 4)
	if len(parts) < 4 {
		return CompilerError{}, fmt.Errorf("invalid compiler error format: %s", raw)
	}

	file := strings.TrimSpace(parts[0])
	line, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return CompilerError{}, fmt.Errorf("invalid line number: %s", parts[1])
	}
	col, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return CompilerError{}, fmt.Errorf("invalid column number: %s", parts[2])
	}
	msg := strings.TrimSpace(parts[3])

	cat := "General"
	switch {
	case strings.Contains(msg, "imported and not used"):
		cat = "Unused Import"
	case strings.Contains(msg, "declared and not used"):
		cat = "Unused Variable"
	case strings.Contains(msg, "cannot use"):
		cat = "Type Mismatch"
	case strings.Contains(msg, "undefined"):
		cat = "Undefined Identifier"
	case strings.Contains(msg, "syntax error"):
		cat = "Syntax Error"
	}

	return CompilerError{
		File:     file,
		Line:     line,
		Col:      col,
		Message:  msg,
		Category: cat,
	}, nil
}

func main() {
	fmt.Println("=== Deconstructing Go Compiler Error Diagnostics ===")

	sampleErrors := []string{
		`main.go:5:2: "os" imported and not used`,
		`main.go:12:15: cannot use "42" (untyped string constant) as int value in assignment`,
		`main.go:18:5: undefined: UserProfile`,
	}

	for i, raw := range sampleErrors {
		errItem, err := ParseCompilerError(raw)
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			continue
		}

		fmt.Printf("[%d] %s Error at %s (Line %d, Col %d):\n", i+1, errItem.Category, errItem.File, errItem.Line, errItem.Col)
		fmt.Printf("    Message: %s\n\n", errItem.Message)
	}

	fmt.Println("Insight: Go compiler errors are deterministic and provide pinpoint coordinate accuracy.")
	fmt.Println("Always address the FIRST error in compiler output, as subsequent errors are often cascades.")
}
