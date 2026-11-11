package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
			if vuln.Module != mod.Path {
				continue
			}
			// Only flag versions BELOW the fix (production code should use
			// golang.org/x/mod/semver). An advisory with no published fix must
			// stay flagged: comparing against an empty FixedIn would parse as
			// v0.0.0 and silently suppress it.
			if vuln.FixedIn == "" || versionLess(mod.Version, "v"+strings.TrimPrefix(vuln.FixedIn, "v")) {
				found = append(found, vuln)
			}
		}
	}
	return found
}

// versionLess compares two semver strings (vMAJOR.MINOR.PATCH[-PRERELEASE]).
// Production code should use golang.org/x/mod/semver.
func versionLess(a, b string) bool {
	parse := func(v string) ([3]int, string) {
		v = strings.TrimPrefix(v, "v")
		// Build metadata (+incompatible) does not affect precedence.
		if i := strings.IndexByte(v, '+'); i >= 0 {
			v = v[:i]
		}
		core, pre, _ := strings.Cut(v, "-")
		var out [3]int
		parts := strings.SplitN(core, ".", 3)
		for i, p := range parts {
			fmt.Sscanf(p, "%d", &out[i])
		}
		return out, pre
	}
	pa, aPre := parse(a)
	pb, bPre := parse(b)
	for i := range pa {
		if pa[i] != pb[i] {
			return pa[i] < pb[i]
		}
	}
	if aPre == "" || bPre == "" {
		return aPre != "" && bPre == ""
	}
	return comparePrerelease(aPre, bPre) < 0
}

// comparePrerelease orders dot-separated prerelease identifiers per SemVer 2.0.0.
func comparePrerelease(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")
	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		ia, ib := partsA[i], partsB[i]
		aNum, bNum := isASCIIDigits(ia), isASCIIDigits(ib)
		if aNum && bNum {
			ta, tb := strings.TrimLeft(ia, "0"), strings.TrimLeft(ib, "0")
			if len(ta) != len(tb) {
				if len(ta) < len(tb) {
					return -1
				}
				return 1
			}
			if ta != tb {
				return strings.Compare(ta, tb)
			}
			continue
		}
		if aNum != bNum {
			if aNum {
				return -1
			}
			return 1
		}
		if ia != ib {
			return strings.Compare(ia, ib)
		}
	}
	if len(partsA) < len(partsB) {
		return -1
	}
	if len(partsA) > len(partsB) {
		return 1
	}
	return 0
}

func isASCIIDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
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
