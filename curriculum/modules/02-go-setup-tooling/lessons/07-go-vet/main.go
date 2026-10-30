package main

import (
	"fmt"
	"strings"
)

// VetDiagnostic represents an issue flagged by static analysis.
type VetDiagnostic struct {
	CheckName string
	Severity  string
	Message   string
	Example   string
}

// DiagnoseCommonVetIssues demonstrates the primary bug patterns go vet catches.
func DiagnoseCommonVetIssues() []VetDiagnostic {
	return []VetDiagnostic{
		{
			CheckName: "printf",
			Severity:  "warning",
			Message:   "Printf format %d has arg 'hello' of wrong type string",
			Example:   `fmt.Printf("%d\n", "hello")`,
		},
		{
			CheckName: "copylocks",
			Severity:  "error",
			Message:   "assignment copies lock value to pass-by-value receiver",
			Example:   `type SafeCounter struct { mu sync.Mutex; val int }; func (c SafeCounter) Inc() {}`,
		},
		{
			CheckName: "unreachable",
			Severity:  "warning",
			Message:   "unreachable code after return or panic",
			Example:   `func foo() int { return 1; fmt.Println("never reached") }`,
		},
		{
			CheckName: "nilfunc",
			Severity:  "warning",
			Message:   "comparison of function f == nil is always false",
			Example:   `if myFunc == nil { ... } (when myFunc is a declared func)`,
		},
	}
}

// CheckFormatArgs checks whether a printf format specifier matches argument counts.
func CheckFormatArgs(format string, argCount int) error {
	specifiers := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			if i+1 < len(format) && format[i+1] == '%' {
				i++ // escaped %%
				continue
			}
			specifiers++
		}
	}

	if specifiers != argCount {
		return fmt.Errorf("format %q requires %d args, but %d were provided", format, specifiers, argCount)
	}
	return nil
}

func main() {
	fmt.Println("=== Deconstructing 'go vet' Static Analysis Engine ===")

	diagnostics := DiagnoseCommonVetIssues()
	for i, d := range diagnostics {
		fmt.Printf("[%d] Checker: %-12s | Severity: %-7s\n", i+1, d.CheckName, strings.ToUpper(d.Severity))
		fmt.Printf("    Problem: %s\n", d.Message)
		fmt.Printf("    Bad Pattern: %s\n\n", d.Example)
	}

	// Demonstration of correct code verified by go vet
	user := "Senior Gopher"
	years := 5
	fmt.Printf("Verified Clean: User %s has %d years of Go experience.\n", user, years)
	fmt.Println("Status: go vet runs in CI to reject subtle bugs before compilation/deployment.")
}
