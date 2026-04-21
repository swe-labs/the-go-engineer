// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 7: Strings & Text — Config Parser (Exercise)
// Level: Intermediate
// ============================================================================
//
// RUN: go run ./04-types-design/strings-and-text/6-config-parser
// TEST: go test ./04-types-design/strings-and-text/6-config-parser
// ============================================================================

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type configEntry struct {
	Key   string
	Value string
}

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
	envFileContent := `
# Application Configuration
APP_NAME="My Cool App"
APP_VERSION="1.0.2-beta" # Version with quotes
PORT=8080
DEBUG_MODE="true"
# Database Settings
DB_HOST=localhost
DB_USER = admin
DB_PASSWORD = "p@s$w Ord With Sp@ces!" # Quoted password
API_ENDPOINT = https://api.example.com/v1

# An empty value
EMPTY_KEY=
ANOTHER_KEY_NO_VALUE =`

	config, err := parseConfig(envFileContent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	summary, err := renderConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(summary)
}
