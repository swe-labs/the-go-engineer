package main

// ============================================================================
// Section 7: Strings & Text — Config Parser (Exercise)
// Level: Intermediate
// ============================================================================
//
// RUN: go run ./07-strings-and-text/6-log-parser
// ============================================================================

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func parseConfig(content string) (map[string]string, error) {
	config := make(map[string]string)

	// 1. Regex Compilation
	// regexp.MustCompile() parses the regex into a state machine byte-code.
	// The "Must" idiom in Go means: "If this fails, trigger a fatal panic."
	// We use MustCompile for hardcoded strings that are guaranteed to be valid.
	re := regexp.MustCompile(`^\s*([\w.-]+)\s*=\s*(?:'([^']*)'|"([^"]*)"|([^#\s]*))?(?:\s*#.*)?$`)

	// 2. The Scanner
	// bufio.Scanner wraps an io.Reader. Instead of loading a 10GB log file
	// directly into RAM (which would crash the server), it streams the bytes
	// from memory lazily, line by line.
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNo := 0

	// .Scan() reads bytes up to the next \n character and returns true.
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()

		// 3. String Trimming
		// Strings are immutable, so TrimSpace allocates a fresh string header
		// pointing to a sliced portion of the underlying byte array.
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines or comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue // Skip to next iteration of the loop
		}

		// 4. Submatch Extraction
		// FindStringSubmatch returns an array indexing the capture groups `(...)`.
		matches := re.FindStringSubmatch(trimmedLine)
		if matches == nil {
			fmt.Printf("Line %d: '%s' - Is Invalid\n", lineNo, line)
			continue
		}

		// matches[0] = the entire matching line
		// matches[1] = the key capture group
		key := matches[1]
		var value string

		if matches[2] != "" {
			value = matches[2]
		} else if matches[3] != "" {
			value = matches[3]
		} else {
			value = matches[4]
		}

		config[key] = value
	}

	return config, nil
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

	for k, v := range config {
		fmt.Printf("%s=%q\n", k, v)
	}

}
