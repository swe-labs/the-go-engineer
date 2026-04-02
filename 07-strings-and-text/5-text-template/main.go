// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 7: Strings & Text — Text Templates
// Level: Intermediate
// ============================================================================
//
// RUN: go run ./07-strings-and-text/5-text-template
// ============================================================================

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type EmailData struct {
	RecipientName string
	SenderName    string
	Subject       string
	Body          string
	Items         []string // demo a loop
	UnreadCount   int
}

func main() {

	fmt.Println("--- Text template example ---")

	emailTemplate := `
Subject: {{ .Subject }}

{{.Body}}

{{if .Items}}
   Related Items:
{{range .Items}}
	- {{.}}
{{end}}
{{end}}

{{if gt .UnreadCount 0}}
You have {{.UnreadCount}} unreads.
{{else}}
You have no messages
{{end}}


- Thanks
{{.SenderName}}
`
	// 1. Template Parsing
	// `New("name")` allocates a new text template in memory.
	// `.Parse()` compiles the raw string into an Abstract Syntax Tree (AST).
	// This parsing only needs to happen ONCE during application startup!
	tmpl, err := template.New("email-message").Parse(emailTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err.Error())
		os.Exit(1)
	}

	// 2. Data Hydration
	// The struct fields MUST be exported (Capitalized) so the template engine
	// (which uses the `reflect` package) can read them at runtime!
	data := EmailData{
		RecipientName: "Alice",
		SenderName:    "Bob's Auto-Responder",
		Subject:       "Your Weekly Update",
		Body:          "Here is the update you requested. We hope you find it useful.",
		Items:         []string{"Report A", "Document B", "Summary C"},
		UnreadCount:   0,
	}

	// 3. The strings.Builder
	// We use strings.Builder instead of string concatenation (+) because strings are
	// immutable in Go. Concatenating 10,000 strings allocates 10,000 new memory blocks.
	// Builder acts as an expandable byte buffer.
	var output strings.Builder

	// 4. Execution
	// `tmpl.Execute` dynamically walks the AST, hydrates the variables using reflection
	// on `data`, and streams the resulting bytes directly into the `output` builder.
	err = tmpl.Execute(&output, data)
	if err != nil {
		fmt.Println("Error executing template:", err.Error())
		os.Exit(1)
	}

	fmt.Println(strings.ToUpper(output.String()))
}
