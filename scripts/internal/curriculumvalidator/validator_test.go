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
