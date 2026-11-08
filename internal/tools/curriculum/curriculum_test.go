package main

import (
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		target   string
		expected int
	}{
		{
			name:     "found at start",
			slice:    []string{"id", "title", "order"},
			target:   "id",
			expected: 0,
		},
		{
			name:     "found in middle",
			slice:    []string{"id", "title", "order"},
			target:   "title",
			expected: 1,
		},
		{
			name:     "found at end",
			slice:    []string{"id", "title", "order"},
			target:   "order",
			expected: 2,
		},
		{
			name:     "not found",
			slice:    []string{"id", "title", "order"},
			target:   "nonexistent",
			expected: -1,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			target:   "any",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indexOf(tt.slice, tt.target)
			if result != tt.expected {
				t.Fatalf("expected indexOf(%v, %q) = %d, got %d", tt.slice, tt.target, tt.expected, result)
			}
		})
	}
}

func TestExtractItemKeysInOrder(t *testing.T) {
	rawJSON := []byte(`{
		"schema_version": "3.1.0",
		"items": [
			{
				"id": "item-01",
				"title": "First Item",
				"order": 1
			},
			{
				"id": "item-02",
				"order": 2,
				"title": "Second Item"
			}
		]
	}`)

	entries := extractItemKeysInOrder(rawJSON)
	if len(entries) != 2 {
		t.Fatalf("expected 2 itemKeyEntries, got %d", len(entries))
	}

	// Verify first item
	if entries[0].id != "item-01" {
		t.Errorf("entry 0: expected id 'item-01', got %q", entries[0].id)
	}
	expectedKeys0 := []string{"id", "title", "order"}
	if len(entries[0].keys) != len(expectedKeys0) {
		t.Fatalf("entry 0: expected %d keys, got %d", len(expectedKeys0), len(entries[0].keys))
	}
	for i, k := range expectedKeys0 {
		if entries[0].keys[i] != k {
			t.Errorf("entry 0 key[%d]: expected %q, got %q", i, k, entries[0].keys[i])
		}
	}

	// Verify second item (preserving explicit byte ordering)
	if entries[1].id != "item-02" {
		t.Errorf("entry 1: expected id 'item-02', got %q", entries[1].id)
	}
	expectedKeys1 := []string{"id", "order", "title"}
	if len(entries[1].keys) != len(expectedKeys1) {
		t.Fatalf("entry 1: expected %d keys, got %d", len(expectedKeys1), len(entries[1].keys))
	}
	for i, k := range expectedKeys1 {
		if entries[1].keys[i] != k {
			t.Errorf("entry 1 key[%d]: expected %q, got %q", i, k, entries[1].keys[i])
		}
	}
}

func TestResolvePath(t *testing.T) {
	// 1. Absolute path returns directly
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}
	abs := filepath.Join(wd, "any-test-file.txt")
	if resolved := resolvePath(abs); resolved != abs {
		t.Errorf("expected absolute path %q to resolve to itself, got %q", abs, resolved)
	}

	// 2. Relative path to existing file
	relFile := "path.core.json"
	resolved := resolvePath(relFile)
	if _, err := os.Stat(resolved); err != nil {
		t.Errorf("expected resolvePath(%q) to point to existing file, got %q: %v", relFile, resolved, err)
	}
}

func TestCanonicalValidationGatesExecution(t *testing.T) {
	if len(validators()) == 0 {
		t.Fatal("validator registry is empty")
	}
	for _, v := range validators() {
		errorCount = 0
		warningCount = 0
		v.Run()
		if errorCount != 0 {
			t.Errorf("%s failed with %d errors", v.Name, errorCount)
		}
	}
}

// TestRegistryResolvesEveryItem asserts that the README contract registry
// resolves every curriculum item to exactly one item-level family.
func TestRegistryResolvesEveryItem(t *testing.T) {
	reg := loadRegistry()
	if len(reg.Contracts) == 0 {
		t.Fatal("registry has no contracts")
	}

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	items := append(append([]Item{}, core.Items...), electives.Items...)
	if len(items) == 0 {
		t.Fatal("no items loaded")
	}
	for _, item := range items {
		if _, ok := reg.resolveItemContract(item); !ok {
			t.Errorf("%s (type=%s subtype=%s): did not resolve to exactly one contract family", item.ID, item.Type, item.Subtype)
		}
		if item.ReadmeContract == nil || item.ReadmeContract.ContractID == "" {
			continue
		}
		declared := item.ReadmeContract.ContractID
		if aliases := reg.aliasesOf(); aliases[declared] != "" {
			declared = aliases[declared]
		}
		if !reg.definedContractID(declared) {
			t.Errorf("%s: declared contract_id %q is not defined in the registry", item.ID, item.ReadmeContract.ContractID)
		}
	}
}

// TestHeadingNormalization covers alias mapping and containment matching.
func TestHeadingNormalization(t *testing.T) {
	aliases := map[string]string{"Key principles": "Core idea", "Mission": "Learning objective"}
	found := normalizeHeadings([]string{"Key principles", "Mental model"}, aliases)

	if !found["core idea"] {
		t.Error("alias 'Key principles' did not normalize to 'Core idea'")
	}
	if !found["mental model"] {
		t.Error("'Mental model' not recorded")
	}
	if !headingSatisfied("Core idea", found) {
		t.Error("Core idea should be satisfied via alias")
	}
	if headingSatisfied("NEXT UP", found) {
		t.Error("NEXT UP should not be satisfied")
	}
	if !headingMatches("Grading Rubric", "Rubric") {
		t.Error("token match failed for 'Grading Rubric' -> 'Rubric'")
	}
	if !headingMatches("7. 4-Tier Assessment Rubric", "Rubric") {
		t.Error("token match failed for enumerated '7. 4-Tier Assessment Rubric' -> 'Rubric'")
	}
	if headingMatches("Submission checklist", "Mission") {
		t.Error("'Submission checklist' must not satisfy 'Mission' (substring false positive)")
	}
	if headingMatches("Subtasks", "Tasks") {
		t.Error("'Subtasks' must not satisfy 'Tasks' (substring false positive)")
	}
	if headingMatches("Repurpose notes", "Purpose") {
		t.Error("'Repurpose notes' must not satisfy 'Purpose' (substring false positive)")
	}
	if headingMatches("Permissions", "Mission") {
		t.Error("'Permissions' must not satisfy 'Mission' (substring false positive)")
	}
}

// TestLockedModulesFromExpected asserts the locked-module set is sourced from
// expected.json, not hardcoded.
func TestLockedModulesFromExpected(t *testing.T) {
	locked := lockedModuleSet()
	if len(locked) == 0 {
		t.Fatal("lockedModuleSet is empty")
	}
	for _, id := range []string{"module-00", "module-01", "module-02", "module-03", "module-04"} {
		if !locked[id] {
			t.Errorf("expected %s to be locked", id)
		}
	}
	if locked["module-05"] {
		t.Error("module-05 should not be locked")
	}
}

// TestDocumentContractResolution verifies document-level contract families
// resolve to the expected documents, including the flagship module and the
// elective-module alias mapping back to module.v3.
func TestDocumentContractResolution(t *testing.T) {
	reg := loadRegistry()

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	if got := canonicalContractID(reg, "elective-module.v3"); got != "module.v3" {
		t.Errorf("elective-module.v3 should canonicalize to module.v3, got %q", got)
	}

	find := func(mods []Module, id string) *Module {
		for i := range mods {
			if mods[i].ID == id {
				return &mods[i]
			}
		}
		return nil
	}

	m18 := find(core.Modules, "module-18")
	if m18 == nil || m18.ReadmeContract == nil {
		t.Fatal("module-18 readme_contract is missing")
	}
	fam := reg.contractByID(canonicalContractID(reg, m18.ReadmeContract.ContractID))
	if fam == nil || fam.Selector.Doc != "flagship-module" {
		t.Fatalf("module-18 should resolve to the flagship document family, got %+v", fam)
	}

	m17 := find(electives.Modules, "module-17")
	if m17 == nil || m17.ReadmeContract == nil {
		t.Fatal("module-17 readme_contract is missing")
	}
	fam = reg.contractByID(canonicalContractID(reg, m17.ReadmeContract.ContractID))
	if fam == nil || fam.Selector.Doc != "module" {
		t.Fatalf("module-17 should resolve to the module document family, got %+v", fam)
	}
}

// writeFixture writes content to path, failing the test on error.
func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// TestInspectGoArtifacts uses temp-dir fixtures to pin the starter/solution
// inspection contract: production and test counts, syntax-error reporting, and
// the unreadable-directory error path.
func TestInspectGoArtifacts(t *testing.T) {
	fset := token.NewFileSet()

	t.Run("valid package with test", func(t *testing.T) {
		dir := t.TempDir()
		writeFixture(t, filepath.Join(dir, "main.go"), "package main\n\nfunc main() {}\n")
		writeFixture(t, filepath.Join(dir, "main_test.go"), "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}\n")

		goFiles, testFiles, issues, err := inspectGoArtifacts(dir, fset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if goFiles != 2 || testFiles != 1 {
			t.Fatalf("got goFiles=%d testFiles=%d, want 2/1", goFiles, testFiles)
		}
		if len(issues) != 0 {
			t.Fatalf("unexpected parse issues: %v", issues)
		}
	})

	t.Run("missing test suite", func(t *testing.T) {
		dir := t.TempDir()
		writeFixture(t, filepath.Join(dir, "main.go"), "package main\n\nfunc main() {}\n")

		goFiles, testFiles, issues, err := inspectGoArtifacts(dir, fset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if goFiles != 1 || testFiles != 0 {
			t.Fatalf("got goFiles=%d testFiles=%d, want 1/0", goFiles, testFiles)
		}
		if len(issues) != 0 {
			t.Fatalf("unexpected parse issues: %v", issues)
		}
	})

	t.Run("syntax error reported", func(t *testing.T) {
		dir := t.TempDir()
		writeFixture(t, filepath.Join(dir, "broken.go"), "package main\n\nfunc main( {\n")

		_, _, issues, err := inspectGoArtifacts(dir, fset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Fatalf("expected 1 parse issue, got %d: %v", len(issues), issues)
		}
	})

	t.Run("subdirectories ignored", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "nested"), 0o755); err != nil {
			t.Fatal(err)
		}
		writeFixture(t, filepath.Join(dir, "nested", "deep.go"), "package nested\n")

		goFiles, _, _, err := inspectGoArtifacts(dir, fset)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if goFiles != 0 {
			t.Fatalf("nested .go files should be ignored, got goFiles=%d", goFiles)
		}
	})

	t.Run("unreadable directory returns error", func(t *testing.T) {
		if _, _, _, err := inspectGoArtifacts(filepath.Join(t.TempDir(), "does-not-exist"), fset); err == nil {
			t.Fatal("expected error for a missing directory")
		}
	})
}

// TestLocateProjectDir pins the slug lookup so a dropped or misspelled project
// directory can never again be silently skipped.
func TestLocateProjectDir(t *testing.T) {
	rootA := t.TempDir()
	rootB := t.TempDir()
	got := filepath.Join(rootB, "found-here")
	if err := os.MkdirAll(got, 0o755); err != nil {
		t.Fatal(err)
	}

	if match := locateProjectDir([]string{rootA, rootB}, "found-here"); match != got {
		t.Fatalf("expected %q, got %q", got, match)
	}
	if match := locateProjectDir([]string{rootA, rootB}, "missing"); match != "" {
		t.Fatalf("expected empty result for missing slug, got %q", match)
	}
}

// TestScanReadmeHeadings fixtures the README heading scanner used by the
// document-contract walker.
func TestScanReadmeHeadings(t *testing.T) {
	data := []byte("# Title\n\n## Mission\n\n```go\nfunc main() {}\n```\n\n## Grading Rubric\n")
	headings, hasGo := scanReadmeHeadings(data)
	if !hasGo {
		t.Error("expected a Go code fence to be detected")
	}
	want := []string{"Mission", "Grading Rubric"}
	if len(headings) != len(want) {
		t.Fatalf("got headings %v, want %v", headings, want)
	}
	for i := range want {
		if headings[i] != want[i] {
			t.Errorf("headings[%d]=%q, want %q", i, headings[i], want[i])
		}
	}

	if _, hasGo := scanReadmeHeadings([]byte("## Mission\n\n```text\nnot go\n```\n")); hasGo {
		t.Error("a non-Go fence must not be detected as Go code")
	}
}
