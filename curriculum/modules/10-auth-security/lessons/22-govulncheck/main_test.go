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
