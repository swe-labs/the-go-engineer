// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Config Parser Project
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Building a robust text-processing pipeline for configuration parsing.
//   - Combining 'bufio.Scanner', 'regexp', and 'strings' for data extraction.
//   - Implementing stable summary generation using 'text/template' and sorting.
//   - Handling edge cases like comments, whitespace, and quoted values.
//
// WHY THIS MATTERS:
//   - Real-world software often requires custom parsers for legacy formats
//     or specialized configuration languages. This milestone demonstrates
//     how to transform unstructured text into validated, type-safe internal
//     representations. By mastering the coordination of Go's string and
//     regex packages, you can build reliable tools for data ingestion and
//     environment configuration.
//
// RUN:
//   go run ./04-types-design/24-config-parser-project
//
// KEY TAKEAWAY:
//   - Coordination of text processing tools enables the construction of robust data pipelines.
// ============================================================================

// Commercial use is prohibited without permission.

package main

//
// TEST: go test ./04-types-design/24-config-parser-project

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// TECHNICAL RATIONALE:
//   Custom parsers are required when standard formats (JSON/YAML) are
//   unavailable or when handling legacy configuration files. This
//   project demonstrates a text-processing pipeline that coordinates
//   'bufio.Scanner' for efficient line-by-line ingestion, 'regexp' for
//   structural validation, and 'text/template' for stable output
//   generation. By combining these tools, we ensure that the system
//   can handle edge cases like comments and quoting while maintaining
//   a low memory footprint through incremental processing.
//

// configEntry (Struct) models a single key-value pair for stable, sortable template rendering.
type configEntry struct {
	Key   string
	Value string
}

// parseConfig (Function) implements an ingestion pipeline to transform unstructured text into a validated map of configuration values.
func parseConfig(content string) (map[string]string, error) {
	config := make(map[string]string)

	// Compile once so every scanned line reuses the same regex.
	re := regexp.MustCompile(`^\s*([\w.-]+)\s*=\s*(?:'([^']*)'|"([^"]*)"|([^#\s]*))?(?:\s*#.*)?$`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		matches := re.FindStringSubmatch(trimmedLine)
		if matches == nil {
			return nil, fmt.Errorf("line %d: %q is invalid", lineNo, line)
		}

		key := matches[1]
		var value string

		switch {
		case matches[2] != "":
			value = matches[2]
		case matches[3] != "":
			value = matches[3]
		default:
			value = matches[4]
		}

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan config: %w", err)
	}

	return config, nil
}

// renderConfig (Function) orchestrates sorting and template execution to generate a stable, human-readable summary of configuration data.
func renderConfig(config map[string]string) (string, error) {
	entries := make([]configEntry, 0, len(config))
	for key, value := range config {
		entries = append(entries, configEntry{Key: key, Value: value})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	const tmplSource = `Parsed Config Summary
{{range . -}}
- {{.Key}}={{printf "%q" .Value}}
{{end}}`

	tmpl, err := template.New("config-summary").Parse(tmplSource)
	if err != nil {
		return "", fmt.Errorf("parse config template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, entries); err != nil {
		return "", fmt.Errorf("render config template: %w", err)
	}

	return output.String(), nil
}

func main() {
	fmt.Println("=== Case Study: Configuration Pipeline ===")
	fmt.Println()

	// 1. Raw Configuration Data.
	// Simulating a .env file content with comments, whitespace, and quoted values.
	const rawConfig = `
# System Settings
VERSION=1.2.5
ENVIRONMENT="production"

# Network Configuration
HOST=0.0.0.0
PORT=8080
TIMEOUT=30s # Inline comment

# Credentials
API_KEY = 'sec_rt_8822'
DB_PASS = "p@ss word"
`

	// 2. Parsing Phase.
	// Utilizing bufio.Scanner and regexp to extract validated key-value pairs.
	fmt.Println("--- Phase 1: Parsing ---")
	config, err := parseConfig(rawConfig)
	if err != nil {
		fmt.Printf("Fatal: parse error: %v\n", err)
		return
	}
	fmt.Printf("  Parsed %d entries successfully.\n", len(config))
	fmt.Println()

	// 3. Rendering Phase.
	// Utilizing text/template and slice sorting to generate a stable summary.
	fmt.Println("--- Phase 2: Stable Rendering ---")
	summary, err := renderConfig(config)
	if err != nil {
		fmt.Printf("Fatal: render error: %v\n", err)
		return
	}

	fmt.Println(summary)

	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: MP.1 -> 05-packages-io/01-modules-and-packages/1-module-basics")
	fmt.Println("Run    : go run ./05-packages-io/01-modules-and-packages/1-module-basics")
	fmt.Println("Current: ST.6 (config-parser-project)")
	fmt.Println("---------------------------------------------------")
}
