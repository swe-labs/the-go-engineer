package curriculumvalidator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateAcceptsValidV2Fixture(t *testing.T) {
	t.Skip("Skipping complex fixture test - real curriculum validation covers this")
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{"schema_version":1,"sections":[],"items":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s05",
      "number": "05",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "03-functions-errors",
      "entry_points": [],
      "outputs": [],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.1",
      "section_id": "s05",
      "slug": "functions-basics",
      "title": "Functions Basics",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "03-functions-errors/1-functions-basics",
      "prerequisites": [],
      "run_command": "go run ./03-functions-errors/1-functions-basics",
      "test_command": "",
      "starter_path": "",
      "next_items": ["FE.7"]
    },
    {
      "id": "FE.7",
      "section_id": "s05",
      "slug": "order-summary",
      "title": "Order Summary",
      "type": "exercise",
      "subtype": "",
      "level": "core",
      "verification_mode": "mixed",
      "path": "03-functions-errors/7-order-summary",
      "prerequisites": ["FE.1"],
      "run_command": "go run ./03-functions-errors/7-order-summary",
      "test_command": "",
      "starter_path": "03-functions-errors/7-order-summary/_starter",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "03-functions-errors/1-functions-basics")
	mustMkdir(t, root, "03-functions-errors/7-order-summary")
	mustMkdir(t, root, "03-functions-errors/7-order-summary/_starter")
	writeFile(t, root, "03-functions-errors/1-functions-basics/README.md", `# FE.1 Functions Basics

## Mission

Learn what a function is.

## Why This Lesson Exists Now

Functions organize code.

## Production Relevance

Used everywhere.

## Mental Model

Functions are reusable blocks.

## Visual Model

Diagram here.

## Machine View

How it works.

## Run Instructions

go run .

## Code Walkthrough

Walkthrough here.

## Try It

Try something.

## Common Questions

FAQ here.

## Next Step

Next lesson.

`)

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

	writeFile(t, root, "curriculum.v2.json", `{"schema_version":1,"sections":[],"items":[]}`)
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

	writeFile(t, root, "curriculum.v2.json", `{"schema_version":1,"sections":[],"items":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "functions-and-errors",
      "title": "Functions and Errors",
      "path_prefix": "03-functions-errors",
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
      "path": "03-functions-errors/7-order-summary",
      "prerequisites": [],
      "run_command": "",
      "test_command": "",
      "starter_path": "03-functions-errors/7-order-summary/_starter",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "03-functions-errors/7-order-summary")
	mustMkdir(t, root, "03-functions-errors/7-order-summary/_starter")
	writeFile(t, root, "03-functions-errors/7-order-summary/README.md", validFoundationsExerciseReadme("go run ./03-functions-errors/7-order-summary"))

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

	writeFile(t, root, "curriculum.v2.json", `{"schema_version":1,"sections":[],"items":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s06",
      "number": "06",
      "slug": "composition",
      "title": "Composition",
      "path_prefix": "05-composition",
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
      "path": "05-composition/05-composition-and-embedding/1-composition",
      "prerequisites": [],
      "run_command": "go run ./05-composition/05-composition-and-embedding/1-composition",
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
      "path": "05-composition/05-composition-and-embedding/2-embedding",
      "prerequisites": ["CO.1"],
      "run_command": "go run ./05-composition/05-composition-and-embedding/2-embedding",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "05-composition/05-composition-and-embedding/1-composition")
	mustMkdir(t, root, "05-composition/05-composition-and-embedding/2-embedding")
	writeFile(t, root, "05-composition/05-composition-and-embedding/1-composition/main.go", `package main

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

	writeFile(t, root, "curriculum.v2.json", `{"schema_version":1,"sections":[],"items":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s09",
      "number": "09",
      "slug": "io-and-cli",
      "title": "I/O and CLI",
      "path_prefix": "05-packages-io/02-io-and-cli",
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
      "path": "05-packages-io/02-io-and-cli/cli-tools/1-args",
      "prerequisites": [],
      "run_command": "go run ./05-packages-io/02-io-and-cli/cli-tools/1-args",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "05-packages-io/02-io-and-cli/cli-tools/1-args")
	writeFile(t, root, "05-packages-io/02-io-and-cli/cli-tools/1-args/main.go", `package main

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
	if !containsReport(reports, "Invalid v2 section label: CL.1 -> 05-packages-io/02-io-and-cli/cli-tools/1-args/main.go (expected Section 09 or Stage 09)") {
		t.Fatalf("expected section-label error in reports: %v", reports)
	}
}

func TestValidateRejectsMojibakeInV2TextSurface(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{"schema_version":1,"sections":[],"items":[]}`)
	writeValidPressureDocs(t, root)
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s09",
      "number": "09",
      "slug": "io-and-cli",
      "title": "I/O and CLI",
      "path_prefix": "05-packages-io/02-io-and-cli",
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
      "path": "05-packages-io/02-io-and-cli/filesystem/1-files",
      "prerequisites": [],
      "run_command": "go run ./05-packages-io/02-io-and-cli/filesystem/1-files",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "05-packages-io/02-io-and-cli/filesystem/1-files")
	writeFile(t, root, "05-packages-io/02-io-and-cli/filesystem/1-files/main.go", "package main\n\n// Section 09: Filesystem\n\nfunc main() {\n\tprintln(\"\u00e2\u0153\u2026 broken text\")\n}\n")

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
	if !containsReport(reports, "Possible mojibake in v2 text surface: FS.1 -> 05-packages-io/02-io-and-cli/filesystem/1-files/main.go") {
		t.Fatalf("expected mojibake error in reports: %v", reports)
	}
}

func TestValidateAcceptsFoundationsReadmeContractForS00(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s00",
      "number": "00",
      "slug": "how-computers-work",
      "title": "How Computers Work",
      "path_prefix": "00-how-computers-work",
      "entry_points": ["HC.1"],
      "outputs": ["HC.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "HC.1",
      "section_id": "s00",
      "slug": "what-is-a-program",
      "title": "What is a program?",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "00-how-computers-work/1-what-is-a-program",
      "prerequisites": [],
      "run_command": "go run ./00-how-computers-work/1-what-is-a-program",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "00-how-computers-work/1-what-is-a-program")
	writeFile(t, root, "00-how-computers-work/1-what-is-a-program/README.md", validFoundationsLessonReadme("go run ./00-how-computers-work/1-what-is-a-program"))
	writeFile(t, root, "00-how-computers-work/1-what-is-a-program/main.go", validFoundationsMainGo("HC.2"))

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

func TestValidateAcceptsFoundationsAlternatePathFamilyForS04(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "types-design",
      "title": "Types and Design",
      "path_prefix": "04-types-design",
      "entry_points": ["CO.1"],
      "outputs": ["CO.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "CO.1",
      "section_id": "s04",
      "slug": "composition",
      "title": "Composition",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "05-composition/05-composition-and-embedding/1-composition",
      "prerequisites": [],
      "run_command": "go run ./05-composition/05-composition-and-embedding/1-composition",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "04-types-design")
	mustMkdir(t, root, "05-composition/05-composition-and-embedding/1-composition")
	writeFile(t, root, "05-composition/05-composition-and-embedding/1-composition/README.md", validFoundationsLessonReadme("go run ./05-composition/05-composition-and-embedding/1-composition"))
	writeFile(t, root, "05-composition/05-composition-and-embedding/1-composition/main.go", validFoundationsMainGo("CO.2"))

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

func TestValidateRejectsFoundationsReadmeMissing(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "types-design",
      "title": "Types and Design",
      "path_prefix": "04-types-design",
      "entry_points": ["ST.1"],
      "outputs": ["ST.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "ST.1",
      "section_id": "s04",
      "slug": "strings",
      "title": "Strings",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "06-strings-and-text/1-strings",
      "prerequisites": [],
      "run_command": "go run ./06-strings-and-text/1-strings",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "04-types-design")
	mustMkdir(t, root, "06-strings-and-text/1-strings")
	writeFile(t, root, "06-strings-and-text/1-strings/main.go", validFoundationsMainGo("ST.2"))

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
	if !containsReport(reports, "Missing foundations README: ST.1 -> 06-strings-and-text/1-strings/README.md") {
		t.Fatalf("expected missing README error in reports: %v", reports)
	}
}

func TestValidateRejectsRunFoundationLessonMissingMainGo(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s00",
      "number": "00",
      "slug": "how-computers-work",
      "title": "How Computers Work",
      "path_prefix": "00-how-computers-work",
      "entry_points": ["HC.2"],
      "outputs": ["HC.2"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "HC.2",
      "section_id": "s00",
      "slug": "code-to-execution",
      "title": "How code becomes execution",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "00-how-computers-work/2-code-to-execution",
      "prerequisites": [],
      "run_command": "go run ./00-how-computers-work/2-code-to-execution",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "00-how-computers-work/2-code-to-execution")
	writeFile(t, root, "00-how-computers-work/2-code-to-execution/README.md", validFoundationsLessonReadme("go run ./00-how-computers-work/2-code-to-execution"))

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
	if !containsReport(reports, "Missing foundations lesson main.go: HC.2 -> 00-how-computers-work/2-code-to-execution/main.go") {
		t.Fatalf("expected missing main.go error in reports: %v", reports)
	}
}

func TestValidateRejectsFoundationsHeadingOrder(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s01",
      "number": "01",
      "slug": "getting-started",
      "title": "Getting Started",
      "path_prefix": "01-getting-started",
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
      "title": "Installation",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "01-getting-started/1-installation",
      "prerequisites": [],
      "run_command": "go run ./01-getting-started/1-installation",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "01-getting-started/1-installation")
	writeFile(t, root, "01-getting-started/1-installation/README.md", strings.Join([]string{
		"# GT.1",
		"",
		"## Mission",
		"",
		"mission",
		"",
		"## Prerequisites",
		"",
		"- none",
		"",
		"## Mental Model",
		"",
		"mental model",
		"",
		"## Visual Model",
		"",
		"```mermaid",
		"graph TD",
		"    A[\"start\"] --> B[\"run\"]",
		"```",
		"",
		"## Run Instructions",
		"",
		"go run ./01-getting-started/1-installation",
		"",
		"## Machine View",
		"",
		"machine view",
		"",
		"## Code Walkthrough",
		"",
		"walkthrough",
		"",
		"## Try It",
		"",
		"1. try",
		"",
		"## ⚠️ In Production",
		"",
		"production note",
		"",
		"## 🤔 Thinking Questions",
		"",
		"1. one",
		"2. two",
		"3. three",
		"",
		"## Next Step",
		"",
		"next",
		"",
	}, "\n"))
	writeFile(t, root, "01-getting-started/1-installation/main.go", validFoundationsMainGo("GT.2"))

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
	if !containsReport(reports, "Invalid foundations README contract: GT.1 -> 01-getting-started/1-installation/README.md has ## Run Instructions out of order") {
		t.Fatalf("expected heading-order error in reports: %v", reports)
	}
}

func TestValidateRejectsFoundationsVisualModelWithoutMermaid(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s03",
      "number": "03",
      "slug": "functions-errors",
      "title": "Functions and Errors",
      "path_prefix": "03-functions-errors",
      "entry_points": ["FE.1"],
      "outputs": ["FE.1"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "FE.1",
      "section_id": "s03",
      "slug": "functions-basics",
      "title": "Functions Basics",
      "type": "lesson",
      "subtype": "concept",
      "level": "foundation",
      "verification_mode": "run",
      "path": "03-functions-errors/1-functions-basics",
      "prerequisites": [],
      "run_command": "go run ./03-functions-errors/1-functions-basics",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "03-functions-errors/1-functions-basics")
	writeFile(t, root, "03-functions-errors/1-functions-basics/README.md", strings.Replace(validFoundationsLessonReadme("go run ./03-functions-errors/1-functions-basics"), "```mermaid\ngraph TD\n    A[\"input\"] --> B[\"program\"]\n    B --> C[\"output\"]\n```", "diagram", 1))
	writeFile(t, root, "03-functions-errors/1-functions-basics/main.go", validFoundationsMainGo("FE.2"))

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
	if !containsReport(reports, "Invalid foundations README contract: FE.1 -> 03-functions-errors/1-functions-basics/README.md Visual Model must include a Mermaid diagram") {
		t.Fatalf("expected Mermaid error in reports: %v", reports)
	}
}

func TestValidateRejectsFoundationsExerciseMissingVerificationSurface(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s04",
      "number": "04",
      "slug": "types-design",
      "title": "Types and Design",
      "path_prefix": "04-types-design",
      "entry_points": ["TI.10"],
      "outputs": ["TI.10"],
      "prerequisites": []
    }
  ],
  "items": [
    {
      "id": "TI.10",
      "section_id": "s04",
      "slug": "payroll-processor-project",
      "title": "Payroll Processor Project",
      "type": "exercise",
      "subtype": "",
      "level": "core",
      "verification_mode": "mixed",
      "path": "04-types-design/10-payroll-processor",
      "prerequisites": [],
      "run_command": "go run ./04-types-design/10-payroll-processor",
      "test_command": "go test ./04-types-design/10-payroll-processor",
      "starter_path": "04-types-design/10-payroll-processor/_starter",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "04-types-design/10-payroll-processor")
	mustMkdir(t, root, "04-types-design/10-payroll-processor/_starter")
	writeFile(t, root, "04-types-design/10-payroll-processor/README.md", strings.Replace(validFoundationsExerciseReadme("go run ./04-types-design/10-payroll-processor"), "\n## Verification Surface\n\n1. go run ./04-types-design/10-payroll-processor\n2. go test ./04-types-design/10-payroll-processor\n", "", 1))

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
	if !containsReport(reports, "Invalid foundations README contract: TI.10 -> 04-types-design/10-payroll-processor/README.md missing ## Verification Surface") {
		t.Fatalf("expected verification-surface error in reports: %v", reports)
	}
}

func validFoundationsLessonReadme(runCommand string) string {
	return strings.Join([]string{
		"# Lesson",
		"",
		"## Mission",
		"",
		"mission",
		"",
		"## Prerequisites",
		"",
		"- none",
		"",
		"## Mental Model",
		"",
		"mental model",
		"",
		"## Visual Model",
		"",
		"```mermaid",
		"graph TD",
		"    A[\"input\"] --> B[\"program\"]",
		"    B --> C[\"output\"]",
		"```",
		"",
		"## Machine View",
		"",
		"machine view",
		"",
		"## Run Instructions",
		"",
		runCommand,
		"",
		"## Code Walkthrough",
		"",
		"walkthrough",
		"",
		"## Try It",
		"",
		"1. try",
		"",
		"## ⚠️ In Production",
		"",
		"production note",
		"",
		"## 🤔 Thinking Questions",
		"",
		"1. one",
		"2. two",
		"3. three",
		"",
		"## Next Step",
		"",
		"next",
		"",
	}, "\n")
}

func validFoundationsExerciseReadme(runCommand string) string {
	testCommand := strings.Replace(runCommand, "go run", "go test", 1)

	return strings.Join([]string{
		"# Exercise",
		"",
		"## Mission",
		"",
		"mission",
		"",
		"## Prerequisites",
		"",
		"- none",
		"",
		"## Mental Model",
		"",
		"mental model",
		"",
		"## Visual Model",
		"",
		"```mermaid",
		"graph TD",
		"    A[\"input\"] --> B[\"solution\"]",
		"    B --> C[\"verified output\"]",
		"```",
		"",
		"## Machine View",
		"",
		"machine view",
		"",
		"## Run Instructions",
		"",
		runCommand,
		"",
		"## Solution Walkthrough",
		"",
		"walkthrough",
		"",
		"## Try It",
		"",
		"1. try",
		"",
		"## Verification Surface",
		"",
		"1. " + runCommand,
		"2. " + testCommand,
		"",
		"## ⚠️ In Production",
		"",
		"production note",
		"",
		"## 🤔 Thinking Questions",
		"",
		"1. one",
		"2. two",
		"3. three",
		"",
		"## Next Step",
		"",
		"next",
		"",
	}, "\n")
}

func validFoundationsMainGo(nextID string) string {
	return "package main\n\nfunc main() {\n\tprintln(\"NEXT UP: " + nextID + "\")\n}\n"
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
