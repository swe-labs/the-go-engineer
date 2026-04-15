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
	writeValidPressureDocs(t, root)
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
	writeValidPressureDocs(t, root)
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
	writeValidPressureDocs(t, root)
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
	writeValidPressureDocs(t, root)
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
	writeValidPressureDocs(t, root)
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
	writeValidPressureDocs(t, root)
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

func TestValidateRejectsRubricSurfaceMissingRequiredHeading(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "docs/stages/expert-layer/tasks/review-db6-repository-boundary.md", `# DB.6 Repository Boundary Review

## Mission

Review the surface.

## Type

- review task

## Level

- core

## Prerequisites

- item

## Task

1. do the review

## Evidence

- show evidence

## Common Weak Answers

- weak answer

## Next Step

next
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
	if !containsReport(reports, "Invalid rubric/checkpoint surface headings: docs/stages/expert-layer/tasks/review-db6-repository-boundary.md missing ## Rubric") {
		t.Fatalf("expected missing-rubric-heading error in reports: %v", reports)
	}
}

func TestValidateRejectsBrokenFlagshipCheckpointLink(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "docs/stages/flagship-project/checkpoint-guidance.md", `# Flagship Project Checkpoint Guidance

Use the checkpoint set:

- [Checkpoint set index](./checkpoints/README.md)
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
	if !containsReport(reports, "Missing required pressure-doc link: docs/stages/flagship-project/checkpoint-guidance.md -> ./slices/README.md") {
		t.Fatalf("expected missing-checkpoint-link error in reports: %v", reports)
	}
}

func TestValidateRejectsBrokenTemplateDocLink(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "docs/templates/README.md", `# Templates

- [Roadmap](./ADVANCED_CONTENT_ROADMAP.md)
`)
	writeFile(t, root, "docs/templates/ADVANCED_CONTENT_ROADMAP.md", `# Roadmap

- [Missing](./DOES_NOT_EXIST.md)
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
	if !containsReport(reports, "Broken local doc link: docs/templates/ADVANCED_CONTENT_ROADMAP.md -> ./DOES_NOT_EXIST.md") {
		t.Fatalf("expected broken-template-link error in reports: %v", reports)
	}
}

func TestValidateAcceptsFoundationsReadmeContract(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "01-foundations/05-functions-and-errors",
      "entry_points": ["FE.1"],
      "outputs": ["FE.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.1",
      "section_id": "s04",
      "slug": "functions-basics",
      "title": "Functions Basics",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "01-foundations/05-functions-and-errors/1-functions-basics",
      "prerequisites": [],
      "run_command": "go run ./01-foundations/05-functions-and-errors/1-functions-basics",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "01-foundations/05-functions-and-errors/1-functions-basics")
	writeFile(t, root, "01-foundations/05-functions-and-errors/1-functions-basics/README.md", `# FE.1

## Mission

mission

## Prerequisites

- none

## Mental Model

mental model

## Visual Model

diagram

## Machine View

machine view

## Run Instructions

go run ./01-foundations/05-functions-and-errors/1-functions-basics

## Code Walkthrough

walkthrough

## Try It

1. try

## Next Step

next
`)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 0 {
		t.Fatalf("expected 0 validation errors, got %d with reports %v", result.ErrorCount, reports)
	}
}

func TestValidateRejectsFoundationsReadmeMissingMachineView(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "01-foundations/05-functions-and-errors",
      "entry_points": ["FE.1"],
      "outputs": ["FE.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.1",
      "section_id": "s04",
      "slug": "functions-basics",
      "title": "Functions Basics",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "01-foundations/05-functions-and-errors/1-functions-basics",
      "prerequisites": [],
      "run_command": "go run ./01-foundations/05-functions-and-errors/1-functions-basics",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "01-foundations/05-functions-and-errors/1-functions-basics")
	writeFile(t, root, "01-foundations/05-functions-and-errors/1-functions-basics/README.md", `# FE.1

## Mission

mission

## Prerequisites

- none

## Mental Model

mental model

## Visual Model

diagram

## Run Instructions

go run ./01-foundations/05-functions-and-errors/1-functions-basics

## Code Walkthrough

walkthrough

## Try It

1. try

## Next Step

next
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
	if !containsReport(reports, "Invalid foundations README contract: FE.1 -> 01-foundations/05-functions-and-errors/1-functions-basics/README.md missing ## Machine View") {
		t.Fatalf("expected missing-machine-view error in reports: %v", reports)
	}
}

func TestValidateRejectsFoundationsExerciseMissingVerificationSurface(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "01-foundations/05-functions-and-errors",
      "entry_points": ["FE.7"],
      "outputs": ["FE.7"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.7",
      "section_id": "s04",
      "slug": "order-summary",
      "title": "Order Summary",
      "type": "exercise",
      "subtype": "",
      "level": "core",
      "verification_mode": "mixed",
      "path": "01-foundations/05-functions-and-errors/7-order-summary",
      "prerequisites": [],
      "run_command": "go run ./01-foundations/05-functions-and-errors/7-order-summary",
      "test_command": "go test ./01-foundations/05-functions-and-errors/7-order-summary",
      "starter_path": "01-foundations/05-functions-and-errors/7-order-summary/_starter",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "01-foundations/05-functions-and-errors/7-order-summary")
	mustMkdir(t, root, "01-foundations/05-functions-and-errors/7-order-summary/_starter")
	writeFile(t, root, "01-foundations/05-functions-and-errors/7-order-summary/README.md", `# FE.7

## Mission

mission

## Visual Model

diagram

## Machine View

machine view

## Run Instructions

go run ./01-foundations/05-functions-and-errors/7-order-summary

## Solution Walkthrough

walkthrough

## Try It

1. try

## Next Step

next
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
	if !containsReport(reports, "Invalid foundations README contract: FE.7 -> 01-foundations/05-functions-and-errors/7-order-summary/README.md missing ## Verification Surface") {
		t.Fatalf("expected missing-verification-surface error in reports: %v", reports)
	}
}

func TestValidateAcceptsGettingStartedReadmeContractAndSplitSectionPrefix(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.json", `{"sections":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s01",
      "number": "01",
      "slug": "core-foundations",
      "title": "Core Foundations",
      "path_prefix": "01-core-foundations",
      "entry_points": ["GT.1"],
      "outputs": ["GT.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "GT.1",
      "section_id": "s01",
      "slug": "installation",
      "title": "Installation Verification",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "01-foundations/01-getting-started/1-installation",
      "prerequisites": [],
      "run_command": "go run ./01-foundations/01-getting-started/1-installation",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "01-foundations/01-getting-started/1-installation")
	mustMkdir(t, root, "01-core-foundations")
	writeFile(t, root, "01-foundations/01-getting-started/1-installation/README.md", `# GT.1

## Mission

mission

## Mental Model

mental model

## Visual Model

diagram

## Machine View

machine view

## Run Instructions

go run ./01-foundations/01-getting-started/1-installation

## Code Walkthrough

walkthrough

## Try It

1. try

## Next Step

next
`)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount != 0 {
		t.Fatalf("expected 0 validation errors, got %d with reports %v", result.ErrorCount, reports)
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

func writeValidPressureDocs(t *testing.T, root string) {
	t.Helper()

	writeFile(t, root, "docs/templates/rubric-checkpoint-template.md", "# Template\n")
	writeFile(t, root, "docs/stages/expert-layer/README.md", "# Expert\n\n[Task index](./tasks/README.md)\n")
	writeFile(t, root, "docs/stages/expert-layer/stage-map.md", "# Stage Map\n\n[Task index](./tasks/README.md)\n")
	writeFile(t, root, "docs/stages/expert-layer/pressure-guidance.md", "# Guidance\n\n[Task index](./tasks/README.md)\n")
	writeFile(t, root, "docs/stages/expert-layer/tasks/README.md", "# Tasks\n\n- [Task](./review-db6-repository-boundary.md)\n")
	writeFile(t, root, "docs/stages/expert-layer/tasks/review-db6-repository-boundary.md", `# Review

## Mission

mission

## Type

- review task

## Level

- core

## Prerequisites

- item

## Task

1. task

## Evidence

- evidence

## Rubric

rubric

## Common Weak Answers

- weak

## Next Step

next
`)
	writeFile(t, root, "docs/stages/flagship-project/README.md", "# Flagship\n\n[Checkpoint set](./checkpoints/README.md)\n[Implementation slices](./slices/README.md)\n")
	writeFile(t, root, "docs/stages/flagship-project/stage-map.md", "# Stage Map\n\n[Checkpoint set](./checkpoints/README.md)\n[Implementation slices](./slices/README.md)\n")
	writeFile(t, root, "docs/stages/flagship-project/checkpoint-guidance.md", "# Guidance\n\n[Checkpoint set index](./checkpoints/README.md)\n[Implementation slices](./slices/README.md)\n")
	writeFile(t, root, "docs/stages/flagship-project/checkpoints/README.md", "# Checkpoints\n\n- [Foundation](./foundation-checkpoint.md)\n")
	writeFile(t, root, "docs/stages/flagship-project/checkpoints/foundation-checkpoint.md", `# Foundation

## Mission

mission

## Type

- flagship checkpoint

## Level

- foundation

## Prerequisites

- item

## Task

1. task

## Evidence

- evidence

## Rubric

rubric

## Common Weak Answers

- weak

## Next Step

[Implementation slice](../slices/foundation-slice.md)
`)
	writeFile(t, root, "docs/stages/flagship-project/slices/README.md", "# Slices\n\n- [Foundation](./foundation-slice.md)\n")
	writeFile(t, root, "docs/stages/flagship-project/slices/foundation-slice.md", "# Slice\n")
}
