package main

import (
	"os"
	"testing"
)

func TestLoadVulnDB(t *testing.T) {
	dbContent := `{
		"vulnerabilities": [
			{
				"id": "GO-2024-0001",
				"module": "golang.org/x/crypto",
				"version": "<0.17.0",
				"description": "Test vuln",
				"fixed_in": "0.17.0",
				"severity": "HIGH"
			}
		]
	}`
	tmpFile := "test_vuln_db.json"
	os.WriteFile(tmpFile, []byte(dbContent), 0644)
	defer os.Remove(tmpFile)

	db, err := loadVulnDB(tmpFile)
	if err != nil {
		t.Fatalf("loadVulnDB error: %v", err)
	}

	if len(db.Vulnerabilities) != 1 {
		t.Errorf("expected 1 vuln, got %d", len(db.Vulnerabilities))
	}
	if db.Vulnerabilities[0].ID != "GO-2024-0001" {
		t.Errorf("ID = %q, want GO-2024-0001", db.Vulnerabilities[0].ID)
	}
}

func TestLoadVulnDBInvalidJSON(t *testing.T) {
	tmpFile := "test_vuln_db_bad.json"
	os.WriteFile(tmpFile, []byte("not json"), 0644)
	defer os.Remove(tmpFile)

	_, err := loadVulnDB(tmpFile)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadVulnDBMissingFile(t *testing.T) {
	_, err := loadVulnDB("nonexistent.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCheckVulnerabilities(t *testing.T) {
	db := &VulnDB{
		Vulnerabilities: []VulnEntry{
			{
				ID:      "GO-2024-0001",
				Module:  "golang.org/x/crypto",
				Version: "<0.17.0",
				FixedIn: "0.17.0",
			},
			{
				ID:      "GO-2024-0002",
				Module:  "github.com/gorilla/mux",
				Version: "<1.8.1",
				FixedIn: "1.8.1",
			},
			{
				ID:      "GO-2024-0003",
				Module:  "example.com/pre",
				Version: "<0.17.0-rc.2",
				FixedIn: "0.17.0-rc.2",
			},
			{
				ID:      "GO-2024-0004",
				Module:  "example.com/pre-v",
				Version: "<v0.17.0-rc.2",
				FixedIn: "v0.17.0-rc.2",
			},
			{
				ID:      "GO-2024-0005",
				Module:  "example.com/unfixed",
				Version: ">=0.0.0",
				FixedIn: "",
			},
		},
	}

	tests := []struct {
		name    string
		modules []ModuleVersion
		wantN   int
	}{
		{
			name: "vulnerable module flagged",
			modules: []ModuleVersion{
				{Path: "golang.org/x/crypto", Version: "v0.16.0"},
			},
			wantN: 1,
		},
		{
			name: "patched version not flagged",
			modules: []ModuleVersion{
				{Path: "golang.org/x/crypto", Version: "v0.17.0"},
			},
			wantN: 0,
		},
		{
			name: "multiple vulnerabilities found",
			modules: []ModuleVersion{
				{Path: "golang.org/x/crypto", Version: "v0.16.0"},
				{Path: "github.com/gorilla/mux", Version: "v1.8.0"},
			},
			wantN: 2,
		},
		{
			name: "no vulnerable modules",
			modules: []ModuleVersion{
				{Path: "github.com/lib/pq", Version: "v1.10.9"},
			},
			wantN: 0,
		},
		{
			name: "earlier prerelease below prerelease fix is flagged",
			modules: []ModuleVersion{
				{Path: "example.com/pre", Version: "v0.17.0-rc.1"},
				{Path: "example.com/pre-v", Version: "v0.17.0-rc.1"},
			},
			wantN: 2,
		},
		{
			name: "matching prerelease or release at prerelease fix is not flagged",
			modules: []ModuleVersion{
				{Path: "example.com/pre", Version: "v0.17.0-rc.2"},
				{Path: "example.com/pre-v", Version: "v0.17.0"},
			},
			wantN: 0,
		},
		{
			name: "advisory with empty FixedIn stays flagged",
			modules: []ModuleVersion{
				{Path: "example.com/unfixed", Version: "v1.2.3"},
			},
			wantN: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkVulnerabilities(tc.modules, db)
			if len(got) != tc.wantN {
				t.Errorf("expected %d vulns, got %d", tc.wantN, len(got))
			}
		})
	}
}

func TestCheckGovulncheckInstalled(t *testing.T) {
	result := checkGovulncheckInstalled()
	_ = result
}

func TestVersionLess(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
		desc string
	}{
		{"v0.16.0", "v0.17.0", true, "older core is less"},
		{"v0.18.0", "v0.17.0", false, "newer core is NOT less"},
		{"v1.0.0-rc.1", "v1.0.0", true, "prerelease sorts before its release"},
		{"v1.0.0", "v1.0.0-rc.1", false, "release sorts after its prerelease"},
		{"v1.0.0-rc.1", "v0.17.0", false, "newer core with prerelease suffix is NOT less (regression)"},
		{"v2.0.0+incompatible", "v2.0.0", false, "build metadata ignored (regression)"},
		{"v2.0.0+build-1", "v2.0.0", false, "hyphen inside build metadata is not a prerelease"},
		{"v0.17.0", "v0.17.0", false, "equal versions are NOT less"},
		{"v0.17.0-rc.1", "v0.17.0-rc.2", true, "earlier prerelease sorts before later prerelease"},
		{"v0.17.0-rc.2", "v0.17.0-rc.1", false, "later prerelease sorts after earlier prerelease"},
		{"v0.17.0-rc.2", "v0.17.0-rc.2", false, "equal prereleases are NOT less"},
		{"v0.17.0-rc.2", "v0.17.0-rc.10", true, "numeric prerelease identifiers compare numerically"},
		{"v1.0.0-1", "v1.0.0-alpha", true, "numeric identifier sorts before non-numeric"},
		{"v1.0.0-alpha", "v1.0.0-alpha.1", true, "fewer prerelease identifiers sort before more"},
		{"v1.0.0-alpha.1", "v1.0.0-alpha.beta", true, "numeric sub-identifier sorts before alphanumeric"},
		{"v1.0.0-alpha.beta", "v1.0.0-beta", true, "lexical comparison on non-numeric identifiers"},
		{"v1.0.0-beta.2", "v1.0.0-beta.11", true, "semver spec numeric beta ordering"},
		{"v1.0.0-beta.11", "v1.0.0-rc.1", true, "semver spec lexical beta vs rc"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if got := versionLess(tc.a, tc.b); got != tc.want {
				t.Errorf("versionLess(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
