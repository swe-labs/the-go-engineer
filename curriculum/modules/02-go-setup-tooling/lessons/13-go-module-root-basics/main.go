package main

import (
	"bufio"
	"fmt"
	"strings"
)

// ModuleDeclaration models the parsed contents of a go.mod file.
type ModuleDeclaration struct {
	ModulePath string
	GoVersion  string
	Requires   []string
}

// ParseGoMod extracts essential directives from a raw go.mod file.
func ParseGoMod(content string) (ModuleDeclaration, error) {
	var mod ModuleDeclaration
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			switch fields[0] {
			case "module":
				mod.ModulePath = fields[1]
			case "go":
				mod.GoVersion = fields[1]
			case "require":
				if len(fields) >= 3 {
					mod.Requires = append(mod.Requires, fmt.Sprintf("%s %s", fields[1], fields[2]))
				}
			}
		}
	}

	if mod.ModulePath == "" {
		return mod, fmt.Errorf("missing 'module' directive")
	}

	return mod, nil
}

// ResolveImportPath demonstrates how subpackage import paths map to filesystem folders.
func ResolveImportPath(modulePath, subDir string) string {
	if subDir == "" {
		return modulePath
	}
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(modulePath, "/"), strings.TrimPrefix(subDir, "/"))
}

func main() {
	fmt.Println("=== Deconstructing 'go.mod' and Go Module Root Mechanics ===")

	sampleGoMod := `
module github.com/my-org/payment-gateway

go 1.23.0

require (
	golang.org/x/crypto v0.24.0
	github.com/google/uuid v1.6.0
)
`

	decl, err := ParseGoMod(sampleGoMod)
	if err != nil {
		fmt.Printf("Error parsing go.mod: %v\n", err)
		return
	}

	fmt.Printf("Module Root Path:   %s\n", decl.ModulePath)
	fmt.Printf("Target Go Version:  %s\n\n", decl.GoVersion)

	fmt.Println("--- Subpackage Path Resolution ---")
	internalPkg := ResolveImportPath(decl.ModulePath, "internal/processor")
	fmt.Printf("Filesystem folder: ./internal/processor\n")
	fmt.Printf("Go import path:    %q\n", internalPkg)

	fmt.Println("\nInsight: The go.mod file establishes the root boundary of the module.")
	fmt.Println("All internal packages are imported using the module path prefix, regardless")
	fmt.Println("of whether the repository has been pushed to a remote host yet.")
}
