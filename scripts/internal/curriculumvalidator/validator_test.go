package curriculumvalidator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateAcceptsValidV2Fixture(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "04-functions-and-errors",
      "entry_points": ["FE.1"],
      "outputs": ["FE.9"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.1",
      "section_id": "s04",
      "slug": "functions",
      "title": "Functions",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "04-functions-and-errors/1-function",
      "prerequisites": [],
      "run_command": "go run ./04-functions-and-errors/1-function",
      "test_command": "",
      "starter_path": "",
      "next_items": ["FE.9"]
    },
    {
      "id": "FE.9",
      "section_id": "s04",
      "slug": "error-handling-project",
      "title": "Error Handling Project",
      "type": "exercise",
      "subtype": "",
      "level": "core",
      "verification_mode": "mixed",
      "path": "04-functions-and-errors/8-error-handling",
      "prerequisites": ["FE.1"],
      "run_command": "go run ./04-functions-and-errors/8-error-handling",
      "test_command": "",
      "starter_path": "04-functions-and-errors/8-error-handling/_starter",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "04-functions-and-errors/1-function")
	mustMkdir(t, root, "04-functions-and-errors/8-error-handling")
	mustMkdir(t, root, "04-functions-and-errors/8-error-handling/_starter")

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 0 {
		t.Fatalf("Validate returned %d errors, reports: %v", result.ErrorCount, reports)
	}
	if !result.HasV2 {
		t.Fatalf("expected v2 metadata to be detected")
	}
}

func TestValidateRejectsUnknownSectionPrerequisite(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "04-functions-and-errors",
      "entry_points": [],
      "outputs": [],
      "prerequisites": ["s03"]
    }
  ],
  "items": []
}`)

	mustMkdir(t, root, "04-functions-and-errors")

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 1 {
		t.Fatalf("expected 1 validation error, got %d with reports %v", result.ErrorCount, reports)
	}
	if !containsReport(reports, "Invalid v2 section prerequisite: s04 -> s03") {
		t.Fatalf("expected section prerequisite error in reports: %v", reports)
	}
}

func TestValidateRejectsMixedContractWithoutCommands(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "04-functions-and-errors",
      "entry_points": ["FE.9"],
      "outputs": ["FE.9"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.9",
      "section_id": "s04",
      "slug": "error-handling-project",
      "title": "Error Handling Project",
      "type": "exercise",
      "subtype": "",
      "level": "core",
      "verification_mode": "mixed",
      "path": "04-functions-and-errors/8-error-handling",
      "prerequisites": [],
      "run_command": "",
      "test_command": "",
      "starter_path": "04-functions-and-errors/8-error-handling/_starter",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "04-functions-and-errors/8-error-handling")
	mustMkdir(t, root, "04-functions-and-errors/8-error-handling/_starter")

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 1 {
		t.Fatalf("expected 1 validation error, got %d with reports %v", result.ErrorCount, reports)
	}
	if !containsReport(reports, "Invalid v2 mixed contract: FE.9 requires run_command or test_command") {
		t.Fatalf("expected mixed-contract error in reports: %v", reports)
	}
}

func TestValidateRejectsLessonNavigationFooterMismatch(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s06",
      "number": "06",
      "slug": "composition",
      "title": "Composition",
      "path_prefix": "06-composition",
      "entry_points": ["CO.1"],
      "outputs": ["CO.2"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "CO.1",
      "section_id": "s06",
      "slug": "composition",
      "title": "Composition",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "06-composition/06-composition-and-embedding/1-composition",
      "prerequisites": [],
      "run_command": "go run ./06-composition/06-composition-and-embedding/1-composition",
      "test_command": "",
      "starter_path": "",
      "next_items": ["CO.2"]
    },
    {
      "id": "CO.2",
      "section_id": "s06",
      "slug": "embedding",
      "title": "Embedding",
      "type": "lesson",
      "subtype": "integration",
      "level": "core",
      "verification_mode": "run",
      "path": "06-composition/06-composition-and-embedding/2-embedding",
      "prerequisites": ["CO.1"],
      "run_command": "go run ./06-composition/06-composition-and-embedding/2-embedding",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "06-composition/06-composition-and-embedding/1-composition")
	mustMkdir(t, root, "06-composition/06-composition-and-embedding/2-embedding")
	writeFile(t, root, "06-composition/06-composition-and-embedding/1-composition/main.go", `package main

// Section 06: Composition

func main() {
	println("NEXT UP: ST.1 strings")
}`)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 1 {
		t.Fatalf("expected 1 validation error, got %d with reports %v", result.ErrorCount, reports)
	}
	if !containsReport(reports, "Invalid v2 lesson navigation footer: CO.1 -> ST.1 (expected CO.2)") {
		t.Fatalf("expected lesson-navigation error in reports: %v", reports)
	}
}

func TestValidateRejectsWrongSectionLabelInV2Source(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s09",
      "number": "09",
      "slug": "io-and-cli",
      "title": "I/O and CLI",
      "path_prefix": "09-io-and-cli",
      "entry_points": ["CL.1"],
      "outputs": ["CL.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "CL.1",
      "section_id": "s09",
      "slug": "args",
      "title": "Args",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "09-io-and-cli/cli-tools/1-args",
      "prerequisites": [],
      "run_command": "go run ./09-io-and-cli/cli-tools/1-args",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "09-io-and-cli/cli-tools/1-args")
	writeFile(t, root, "09-io-and-cli/cli-tools/1-args/main.go", `package main

// Section 19: CLI Tools - Command-Line Arguments

func main() {}
`)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 1 {
		t.Fatalf("expected 1 validation error, got %d with reports %v", result.ErrorCount, reports)
	}
	if !containsReport(reports, "Invalid v2 section label: CL.1 -> 09-io-and-cli/cli-tools/1-args/main.go (expected Section 09)") {
		t.Fatalf("expected section-label error in reports: %v", reports)
	}
}

func TestValidateRejectsMojibakeInV2TextSurface(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s09",
      "number": "09",
      "slug": "io-and-cli",
      "title": "I/O and CLI",
      "path_prefix": "09-io-and-cli",
      "entry_points": ["FS.1"],
      "outputs": ["FS.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FS.1",
      "section_id": "s09",
      "slug": "files",
      "title": "Files",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "09-io-and-cli/filesystem/1-files",
      "prerequisites": [],
      "run_command": "go run ./09-io-and-cli/filesystem/1-files",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "09-io-and-cli/filesystem/1-files")
	writeFile(t, root, "09-io-and-cli/filesystem/1-files/main.go", "package main\n\n// Section 09: Filesystem\n\nfunc main() {\n\tprintln(\"\u00e2\u0153\u2026 broken text\")\n}\n")

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 1 {
		t.Fatalf("expected 1 validation error, got %d with reports %v", result.ErrorCount, reports)
	}
	if !containsReport(reports, "Possible mojibake in v2 text surface: FS.1 -> 09-io-and-cli/filesystem/1-files/main.go") {
		t.Fatalf("expected mojibake error in reports: %v", reports)
	}
}

func writeFile(t *testing.T, root, relativePath, contents string) {
	t.Helper()

	fullPath := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("failed to create parent directory for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", relativePath, err)
	}
}

func mustMkdir(t *testing.T, root, relativePath string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(root, relativePath), 0o755); err != nil {
		t.Fatalf("failed to create directory %s: %v", relativePath, err)
	}
}

func containsReport(reports []string, want string) bool {
	for _, report := range reports {
		if strings.Contains(report, want) {
			return true
		}
	}

	return false
}
