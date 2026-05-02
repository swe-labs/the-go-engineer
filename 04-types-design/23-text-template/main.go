// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Text Templates
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining and parsing text templates using 'text/template'.
//   - Utilizing dot notation (.), pipelines, and actions (if, range).
//   - Hydrating templates with exported struct fields and maps.
//   - Executing templates efficiently into io.Writer destinations.
//
// WHY THIS MATTERS:
//   - Templates decouple presentation logic from domain data, which is
//     essential for generating dynamic emails, configuration files, or
//     web pages. The 'text/template' package provides a safe,
//     performance-oriented mechanism for text generation that leverages
//     reflection to hydrate placeholders. Mastering templates ensures
//     that your text-processing logic remains modular and maintainable
//     as requirements evolve.
//
// RUN:
//   go run ./04-types-design/23-text-template
//
// KEY TAKEAWAY:
//   - Templates decouple presentation from data through reflection-based hydration.
// ============================================================================

// Commercial use is prohibited without permission.

package main

//

import (
	"fmt"
	"os"
	"text/template"
)

// TECHNICAL RATIONALE:
//   Templates provide a declarative mechanism for generating text output,
//   effectively decoupling the data structure from its visual
//   representation. The 'text/template' engine utilizes reflection
//   to map exported struct fields or map keys to placeholders at
//   runtime. By parsing templates once into an Abstract Syntax Tree
//   (AST) and executing them against an io.Writer, Go ensures
//   high-performance text generation with minimal memory overhead
//   compared to manual string concatenation.
//

// EmailData (Struct) aggregates domain data into exported fields for reflection-based template hydration.
type EmailData struct {
	RecipientName string
	SenderName    string
	Subject       string
	Body          string
	Items         []string // demo a loop
	UnreadCount   int
}

func main() {
	fmt.Println("=== Text Templates: Decoupled Presentation ===")
	fmt.Println()

	// 1. Template Definition.
	// Placeholders are denoted by double curly braces {{ }}. The dot (.)
	// refers to the data object passed during execution.
	const emailTemplate = `
Subject: {{ .Subject }}

Hello {{ .RecipientName }},

{{ .Body }}

{{ if .Items }}
Priority Items:
{{ range .Items }}
  - {{ . }}
{{ end }}
{{ end }}

{{ if gt .UnreadCount 0 }}
Status: You have {{ .UnreadCount }} unread messages.
{{ else }}
Status: Inbox clear.
{{ end }}

Regards,
{{ .SenderName }}
`

	// 2. Compilation (Parsing).
	// tmpl (template.Template) encapsulates the parsed Abstract Syntax Tree (AST) of the template for efficient execution.
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		fmt.Printf("Fatal: failed to parse template: %v\n", err)
		return
	}

	// 3. Data Context (Hydration).
	// Fields MUST be exported for the reflection-based template engine to
	// access them at runtime.
	data := EmailData{
		RecipientName: "Alice",
		SenderName:    "Opslane Automation",
		Subject:       "System Status Update",
		Body:          "The scheduled maintenance was successful.",
		Items:         []string{"DB Migration", "Cache Warmup"},
		UnreadCount:   2,
	}

	// 4. Execution.
	// Execute hydrates the template AST with data and writes the result
	// to an io.Writer (in this case, os.Stdout).
	fmt.Println("--- Generated Output ---")
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Printf("Fatal: failed to execute template: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ST.6 -> 04-types-design/24-config-parser-project")
	fmt.Println("Run    : go run ./04-types-design/24-config-parser-project")
	fmt.Println("Current: ST.5 (text-template)")
	fmt.Println("---------------------------------------------------")
}
