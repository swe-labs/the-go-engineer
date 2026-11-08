package main

import (
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

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

// TestHasGeneratedMarkers pins the generated-boundary contract: both markers
// must be present before a generated file is considered safe to regenerate.
func TestHasGeneratedMarkers(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{"both markers", "<!-- curriculum:generated:start -->\n# Catalog\n<!-- curriculum:generated:end -->\n", true},
		{"missing end", "<!-- curriculum:generated:start -->\n# Catalog\n", false},
		{"missing start", "# Catalog\n<!-- curriculum:generated:end -->\n", false},
		{"no markers", "# Catalog\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasGeneratedMarkers([]byte(tt.data)); got != tt.want {
				t.Fatalf("hasGeneratedMarkers() = %t, want %t", got, tt.want)
			}
		})
	}
}

// TestHeadingMatches pins whole-word, order-preserving heading matching. These
// are the false positives that substring matching used to accept.
func TestHeadingMatches(t *testing.T) {
	tests := []struct {
		name     string
		present  string
		required string
		want     bool
	}{
		{"exact", "Mission", "Mission", true},
		{"case and numbering", "7. Mission", "Mission", true},
		{"prefix allowed", "Grading rubric", "Rubric", true},
		{"mission vs submission", "Submission checklist", "Mission", false},
		{"tasks vs subtasks", "Subtasks", "Tasks", false},
		{"purpose vs repurpose", "Repurpose notes", "Purpose", false},
		{"evidence required absent suffix", "Evidence", "Evidence required", false},
		{"required is subsequence", "Evidence required for mastery", "Evidence required", true},
		{"empty required", "Mission", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := headingMatches(tt.present, tt.required); got != tt.want {
				t.Fatalf("headingMatches(%q, %q) = %t, want %t", tt.present, tt.required, got, tt.want)
			}
		})
	}
}

// TestResolveAssessmentDirName pins the manifest-first resolution rule that
// replaced the hardcoded assessment slugs.
func TestResolveAssessmentDirName(t *testing.T) {
	tests := []struct {
		name string
		rec  assessmentRecord
		want string
	}{
		{"artifact_dir wins", assessmentRecord{ID: "assessment-module-01", ArtifactDir: "shell-git-lab"}, "shell-git-lab"},
		{"slug wins", assessmentRecord{ID: "assessment-project-x", Slug: "custom-slug"}, "custom-slug"},
		{"project prefix", assessmentRecord{ID: "assessment-project-tooling-failure-lab"}, "tooling-failure-lab"},
		{"module checkpoint default", assessmentRecord{ID: "assessment-module-07"}, "checkpoint"},
		{"generic prefix", assessmentRecord{ID: "assessment-weird-format"}, "weird-format"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveAssessmentDirName(tt.rec); got != tt.want {
				t.Fatalf("resolveAssessmentDirName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestDeriveNextModuleLinks pins the metadata-derived SESE chain, including the
// 16 -> 18 jump and the terminal module having no forward link.
func TestDeriveNextModuleLinks(t *testing.T) {
	modules := []Module{
		{ID: "module-18", Order: 18, Path: "18-flagship-opslane"},
		{ID: "module-00", Order: 0, Path: "00-orientation"},
		{ID: "module-16", Order: 16, Path: "16-portfolio-interview-readiness"},
	}
	links := deriveNextModuleLinks(modules)
	if got, want := links["module-00"], "../16-portfolio-interview-readiness/README.md"; got != want {
		t.Errorf("module-00 link = %q, want %q", got, want)
	}
	if got, want := links["module-16"], "../18-flagship-opslane/README.md"; got != want {
		t.Errorf("module-16 link = %q, want %q", got, want)
	}
	if _, ok := links["module-18"]; ok {
		t.Error("terminal module must not have a forward link")
	}
}

// TestCatalogGeneratorMatchesCurriculumReadme verifies that curriculum/README.md
// is strictly in sync with the canonical JSON manifests via generateCatalog().
func TestCatalogGeneratorMatchesCurriculumReadme(t *testing.T) {
	expectedCatalog, err := generateCatalog()
	if err != nil {
		t.Fatalf("generateCatalog failed: %v", err)
	}
	readmePath := filepath.Join(curriculumDir, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("cannot read %s: %v", readmePath, err)
	}
	actual := string(data)
	if actual != expectedCatalog {
		t.Errorf("curriculum/README.md is out of sync with manifests; run 'make catalog' or 'go run ./internal/tools/curriculum generate-catalog' to update")
	}
}
