package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type VulnEntry struct {
	ID          string `json:"id"`
	Package     string `json:"package"`
	Module      string `json:"module"`
	Version     string `json:"version"`
	Description string `json:"description"`
	FixedIn     string `json:"fixed_in"`
	Severity    string `json:"severity"`
}

type VulnDB struct {
	Vulnerabilities []VulnEntry `json:"vulnerabilities"`
}

func loadVulnDB(path string) (*VulnDB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var db VulnDB
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}
	return &db, nil
}

type ModuleVersion struct {
	Path    string
	Version string
}

func checkVulnerabilities(modules []ModuleVersion, vulnDB *VulnDB) []VulnEntry {
	var found []VulnEntry
	for _, mod := range modules {
		for _, vuln := range vulnDB.Vulnerabilities {
			if vuln.Module == mod.Path {
				if mod.Version != "v"+vuln.FixedIn {
					found = append(found, vuln)
				}
			}
		}
	}
	return found
}

func checkGovulncheckInstalled() bool {
	_, err := exec.LookPath("govulncheck")
	return err == nil
}

func main() {
	fmt.Println("=== govulncheck Demo ===")

	vulnDB := &VulnDB{
		Vulnerabilities: []VulnEntry{
			{
				ID:          "GO-2024-0001",
				Module:      "golang.org/x/crypto",
				Version:     "<0.17.0",
				Description: "Denial of service in RSA key generation",
				FixedIn:     "0.17.0",
				Severity:    "HIGH",
			},
			{
				ID:          "GO-2024-0002",
				Module:      "github.com/gorilla/mux",
				Version:     "<1.8.1",
				Description: "Path traversal vulnerability in ServeHTTP",
				FixedIn:     "1.8.1",
				Severity:    "MEDIUM",
			},
		},
	}

	modules := []ModuleVersion{
		{Path: "golang.org/x/crypto", Version: "v0.16.0"},
		{Path: "github.com/gorilla/mux", Version: "v1.8.0"},
		{Path: "github.com/lib/pq", Version: "v1.10.9"},
	}

	fmt.Println("\nScanning modules for vulnerabilities...")
	vulnerabilities := checkVulnerabilities(modules, vulnDB)

	if len(vulnerabilities) == 0 {
		fmt.Println("No vulnerabilities found.")
	} else {
		fmt.Printf("Found %d vulnerability(ies):\n", len(vulnerabilities))
		for _, v := range vulnerabilities {
			fmt.Printf("  [%s] %s (%s)\n", v.Severity, v.ID, v.Module)
			fmt.Printf("    %s\n", v.Description)
			fmt.Printf("    Fix: update to %s\n", v.FixedIn)
		}
	}

	fmt.Println("\n=== govulncheck integration ===")
	if checkGovulncheckInstalled() {
		fmt.Println("govulncheck is installed. Run: govulncheck ./...")
	} else {
		fmt.Println("govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest")
	}

	fmt.Println("\n=== Vulnerability Database (OSV format) ===")
	fmt.Println("Go vulnerability database: https://vuln.go.dev")
	fmt.Println("OSV format: Open Source Vulnerability format")
}
