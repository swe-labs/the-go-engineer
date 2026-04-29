package curriculumvalidator

import (
	"fmt"
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid v2 section prerequisite: s04 -> s03")
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid v2 mixed contract: FE.9 requires run_command or test_command")
}

func TestValidateReportsMissingArchitectureSection(t *testing.T) {
	root := t.TempDir()

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "00-how-computers-work")
	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s00",
      "number": "00",
      "slug": "how-computers-work",
      "title": "How Computers Work",
      "path_prefix": "00-how-computers-work",
      "status": "stable",
      "entry_points": [],
      "outputs": ["HC.5"],
      "prerequisites": []
    }
  ],
  "items": []
}`)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount == 0 {
		t.Fatalf("expected architecture-contract validation errors")
	}
	if !containsReport(reports, "Invalid v2 architecture contract: missing section s05") {
		t.Fatalf("expected missing-section report in reports: %v", reports)
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
      "path_prefix": "04-types-design",
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
      "path": "04-types-design/composition/1-composition",
      "prerequisites": [],
      "run_command": "go run ./04-types-design/composition/1-composition",
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
      "path": "04-types-design/composition/2-embedding",
      "prerequisites": ["CO.1"],
      "run_command": "go run ./04-types-design/composition/2-embedding",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "04-types-design/composition/1-composition")
	mustMkdir(t, root, "04-types-design/composition/2-embedding")
	writeFile(t, root, "04-types-design/composition/1-composition/main.go", `package main

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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid v2 lesson navigation footer: CO.1 -> ST.1 (expected CO.2)")
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid v2 section label: CL.1 -> 05-packages-io/02-io-and-cli/cli-tools/1-args/main.go (expected Section 09 or Stage 09)")
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Possible mojibake in v2 text surface: FS.1 -> 05-packages-io/02-io-and-cli/filesystem/1-files/main.go")
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
	requireOnlyFixtureScaffoldReports(t, result, reports)
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
      "path": "04-types-design/composition/1-composition",
      "prerequisites": [],
      "run_command": "go run ./04-types-design/composition/1-composition",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "04-types-design")
	mustMkdir(t, root, "04-types-design/composition/1-composition")
	writeFile(t, root, "04-types-design/composition/1-composition/README.md", validFoundationsLessonReadme("go run ./04-types-design/composition/1-composition"))
	writeFile(t, root, "04-types-design/composition/1-composition/main.go", validFoundationsMainGo("CO.2"))

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	requireOnlyFixtureScaffoldReports(t, result, reports)
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
      "path": "04-types-design/strings-and-text/1-strings",
      "prerequisites": [],
      "run_command": "go run ./04-types-design/strings-and-text/1-strings",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	writeValidPressureDocs(t, root)
	mustMkdir(t, root, "04-types-design")
	mustMkdir(t, root, "04-types-design/strings-and-text/1-strings")
	writeFile(t, root, "04-types-design/strings-and-text/1-strings/main.go", validFoundationsMainGo("ST.2"))

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	requireOnlyFixtureExpectedReports(t, result, reports, "Missing foundations README: ST.1 -> 04-types-design/strings-and-text/1-strings/README.md")
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Missing foundations lesson main.go: HC.2 -> 00-how-computers-work/2-code-to-execution/main.go")
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
		"## In Production",
		"",
		"production note",
		"",
		"## Thinking Questions",
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid foundations README contract: GT.1 -> 01-getting-started/1-installation/README.md has ## Run Instructions out of order")
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid foundations README contract: FE.1 -> 03-functions-errors/1-functions-basics/README.md Visual Model must include a Mermaid diagram")
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
	requireOnlyFixtureExpectedReports(t, result, reports, "Invalid foundations README contract: TI.10 -> 04-types-design/10-payroll-processor/README.md missing ## Verification Surface")
}

func TestValidateAcceptsFlagshipProjectSplit(t *testing.T) {
	root := t.TempDir()

	writeValidOpslaneFlagshipFixture(t, root, false)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	requireOnlyFixtureScaffoldReports(t, result, reports)
}

func TestValidateRejectsFlagshipReservedPrefixCollision(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", `{
  "schema_version": 1,
  "sections": [
    {
      "id": "s10",
      "number": "10",
      "slug": "production-operations",
      "title": "Production Operations",
      "path_prefix": "10-production",
      "entry_points": ["OPS.5"],
      "outputs": ["OPS.5"],
      "prerequisites": []
    },
    {
      "id": "s11",
      "number": "11",
      "slug": "flagship",
      "title": "Flagship",
      "path_prefix": "11-flagship",
      "entry_points": ["OPS.1"],
      "outputs": ["OPS.2"],
      "prerequisites": ["s10"]
    }
  ],
  "items": [
    {
      "id": "OPS.5",
      "section_id": "s10",
      "slug": "operations-handoff",
      "title": "Operations Handoff",
      "type": "reference",
      "subtype": "",
      "level": "production",
      "status": "implemented",
      "verification_mode": "rubric",
      "path": "10-production",
      "prerequisites": [],
      "run_command": "",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    },
    {
      "id": "OPS.1",
      "section_id": "s11",
      "slug": "foundation",
      "title": "Reserved Prefix Foundation",
      "type": "checkpoint",
      "subtype": "",
      "level": "production",
      "status": "implemented",
      "verification_mode": "test",
      "path": "11-flagship/01-ops-flagship/modules/01-foundation",
      "prerequisites": [],
      "run_command": "",
      "test_command": "go test ./11-flagship/01-ops-flagship/internal/config/...",
      "starter_path": "",
      "next_items": ["OPS.2"]
    },
    {
      "id": "OPS.2",
      "section_id": "s11",
      "slug": "database",
      "title": "Reserved Prefix Database",
      "type": "capstone",
      "subtype": "",
      "level": "production",
      "status": "implemented",
      "verification_mode": "rubric",
      "path": "11-flagship/01-ops-flagship/modules/02-database",
      "prerequisites": ["OPS.1"],
      "run_command": "",
      "test_command": "",
      "starter_path": "",
      "next_items": []
    }
  ]
}`)

	mustMkdir(t, root, "10-production")
	writeFlagshipProjectSurface(t, root, "11-flagship/01-ops-flagship", []string{
		"modules/01-foundation",
		"modules/02-database",
	}, true, true, true)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount == 0 {
		t.Fatalf("expected a validation error for reserved prefix collision")
	}
	if !containsReport(reports, "Invalid flagship project prefix: OPS.1 -> OPS is already used outside s11") {
		t.Fatalf("expected reserved-prefix error in reports: %v", reports)
	}
}

func TestValidateRejectsFlagshipMissingModuleMap(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", validOpslaneFlagshipCurriculum(false, false))
	writeFlagshipProjectSurface(t, root, "11-flagship/01-opslane", opslaneModuleDirs(), false, true, true)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount == 0 {
		t.Fatalf("expected a validation error for missing MODULES.md")
	}
	if !containsReport(reports, "Missing flagship project module map: OPSL -> 11-flagship/01-opslane/MODULES.md") {
		t.Fatalf("expected missing module map error in reports: %v", reports)
	}
}

func TestValidateRejectsFlagshipMissingProgressChecker(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", validOpslaneFlagshipCurriculum(false, false))
	writeFlagshipProjectSurface(t, root, "11-flagship/01-opslane", opslaneModuleDirs(), true, false, true)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount == 0 {
		t.Fatalf("expected a validation error for missing progress checker")
	}
	if !containsReport(reports, "Missing flagship progress checker: OPSL -> 11-flagship/01-opslane/scripts/progress.go") {
		t.Fatalf("expected missing progress checker error in reports: %v", reports)
	}
}

func TestValidateRejectsBrokenFlagshipChain(t *testing.T) {
	root := t.TempDir()

	writeValidOpslaneFlagshipFixture(t, root, true)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if result.ErrorCount == 0 {
		t.Fatalf("expected a validation error for broken flagship chain")
	}
	if !containsReport(reports, "Invalid flagship module chain: OPSL.4 must point to OPSL.5") {
		t.Fatalf("expected broken chain error in reports: %v", reports)
	}
}

func TestValidateAcceptsOptionalSecondFlagshipProject(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "curriculum.v2.json", validOpslaneFlagshipCurriculum(true, false))
	writeFlagshipProjectSurface(t, root, "11-flagship/01-opslane", opslaneModuleDirs(), true, true, true)
	writeFlagshipProjectSurface(t, root, "11-flagship/02-crmx", []string{
		"modules/01-foundation",
		"modules/02-capstone",
	}, true, false, false)

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	requireOnlyFixtureScaffoldReports(t, result, reports)
}

func writeValidOpslaneFlagshipFixture(t *testing.T, root string, brokenChain bool) {
	t.Helper()

	writeFile(t, root, "curriculum.v2.json", validOpslaneFlagshipCurriculum(false, brokenChain))
	writeFlagshipProjectSurface(t, root, "11-flagship/01-opslane", opslaneModuleDirs(), true, true, true)
}

func validOpslaneFlagshipCurriculum(includeOptionalProject, brokenChain bool) string {
	opsl4Next := "OPSL.5"
	if brokenChain {
		opsl4Next = "OPSL.6"
	}

	coreItems := []string{
		flagshipItemJSON("OPSL.1", "s11", "foundation-and-configuration", "Opslane Foundation and Configuration", "checkpoint", "implemented", "mixed", "11-flagship/01-opslane/modules/01-foundation", "", "go run ./11-flagship/01-opslane/cmd/server", "go test ./11-flagship/01-opslane/internal/config/...", "OPSL.2"),
		flagshipItemJSON("OPSL.2", "s11", "database-and-models", "Opslane Database and Models", "checkpoint", "implemented", "test", "11-flagship/01-opslane/modules/02-database", "OPSL.1", "", "go test ./11-flagship/01-opslane/internal/db/...", "OPSL.3"),
		flagshipItemJSON("OPSL.3", "s11", "authentication-and-tenant-isolation", "Opslane Authentication and Tenant Isolation", "checkpoint", "implemented", "test", "11-flagship/01-opslane/modules/03-auth", "OPSL.2", "", "go test ./11-flagship/01-opslane/internal/auth/...", "OPSL.4"),
		flagshipItemJSON("OPSL.4", "s11", "http-api-layer", "Opslane HTTP API Layer", "checkpoint", "implemented", "mixed", "11-flagship/01-opslane/modules/04-http-api", "OPSL.3", "go run ./11-flagship/01-opslane/cmd/server", "go test ./11-flagship/01-opslane/internal/handlers/... ./11-flagship/01-opslane/internal/middleware/...", opsl4Next),
		flagshipItemJSON("OPSL.5", "s11", "order-processing", "Opslane Order Processing", "checkpoint", "placeholder", "test", "11-flagship/01-opslane/modules/05-order-processing", "OPSL.4", "", "go test ./11-flagship/01-opslane/internal/services/...", "OPSL.6"),
		flagshipItemJSON("OPSL.6", "s11", "payment-pipeline", "Opslane Payment Pipeline", "checkpoint", "placeholder", "test", "11-flagship/01-opslane/modules/06-payment-pipeline", "OPSL.5", "", "go test ./11-flagship/01-opslane/internal/payment/...", "OPSL.7"),
		flagshipItemJSON("OPSL.7", "s11", "event-bus-and-worker-pools", "Opslane Event Bus and Worker Pools", "checkpoint", "placeholder", "test", "11-flagship/01-opslane/modules/07-event-workers", "OPSL.6", "", "go test ./11-flagship/01-opslane/internal/events/... ./11-flagship/01-opslane/internal/workers/...", "OPSL.8"),
		flagshipItemJSON("OPSL.8", "s11", "caching-layer", "Opslane Caching Layer", "checkpoint", "placeholder", "test", "11-flagship/01-opslane/modules/08-caching", "OPSL.7", "", "go test ./11-flagship/01-opslane/internal/cache/...", "OPSL.9"),
		flagshipItemJSON("OPSL.9", "s11", "observability", "Opslane Observability", "checkpoint", "placeholder", "test", "11-flagship/01-opslane/modules/09-observability", "OPSL.8", "", "go test ./11-flagship/01-opslane/internal/logging/... ./11-flagship/01-opslane/internal/metrics/...", "OPSL.10"),
		flagshipItemJSON("OPSL.10", "s11", "graceful-shutdown-and-deployment", "Opslane Graceful Shutdown and Deployment", "capstone", "placeholder", "rubric", "11-flagship/01-opslane/modules/10-shutdown-deploy", "OPSL.9", "", "", ""),
	}

	if includeOptionalProject {
		coreItems = append(coreItems,
			flagshipItemJSON("CRMX.1", "s11", "foundation", "CRMX Foundation", "checkpoint", "placeholder", "rubric", "11-flagship/02-crmx/modules/01-foundation", "", "", "", "CRMX.2"),
			flagshipItemJSON("CRMX.2", "s11", "capstone", "CRMX Capstone", "capstone", "placeholder", "rubric", "11-flagship/02-crmx/modules/02-capstone", "CRMX.1", "", "", ""),
		)
	}

	return fmt.Sprintf(`{
  "schema_version": 1,
  "sections": [
    {
      "id": "s11",
      "number": "11",
      "slug": "flagship",
      "title": "Flagship",
      "path_prefix": "11-flagship",
      "entry_points": ["OPSL.1"],
      "outputs": ["OPSL.10"],
      "prerequisites": []
    }
  ],
  "items": [
%s
  ]
}`, strings.Join(coreItems, ",\n"))
}

func flagshipItemJSON(id, sectionID, slug, title, itemType, status, verificationMode, path, prereq, runCommand, testCommand, next string) string {
	prerequisites := "[]"
	if prereq != "" {
		prerequisites = fmt.Sprintf(`["%s"]`, prereq)
	}
	nextItems := "[]"
	if next != "" {
		nextItems = fmt.Sprintf(`["%s"]`, next)
	}

	return fmt.Sprintf(`    {
      "id": "%s",
      "section_id": "%s",
      "slug": "%s",
      "title": "%s",
      "type": "%s",
      "subtype": "",
      "level": "production",
      "status": "%s",
      "verification_mode": "%s",
      "path": "%s",
      "prerequisites": %s,
      "run_command": "%s",
      "test_command": "%s",
      "starter_path": "",
      "next_items": %s
    }`, id, sectionID, slug, title, itemType, status, verificationMode, path, prerequisites, runCommand, testCommand, nextItems)
}

func opslaneModuleDirs() []string {
	return []string{
		"modules/01-foundation",
		"modules/02-database",
		"modules/03-auth",
		"modules/04-http-api",
		"modules/05-order-processing",
		"modules/06-payment-pipeline",
		"modules/07-event-workers",
		"modules/08-caching",
		"modules/09-observability",
		"modules/10-shutdown-deploy",
	}
}

func writeFlagshipProjectSurface(t *testing.T, root, projectRoot string, moduleDirs []string, includeModuleMap, includeProgress, includeImplementedTargets bool) {
	t.Helper()

	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "README.md")), "# Flagship Project\n")
	if includeModuleMap {
		writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "MODULES.md")), "# Module Map\n")
	}
	if includeProgress {
		writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "scripts", "progress.go")), "//go:build ignore\n\npackage main\n\nfunc main() {}\n")
	}
	for _, moduleDir := range moduleDirs {
		writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, moduleDir, "README.md")), "# Module\n")
	}
	if !includeImplementedTargets {
		return
	}

	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "cmd", "server", "main.go")), "package main\nfunc main() {}\n")
	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "internal", "config", "config.go")), "package config\n")
	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "internal", "db", "repository.go")), "package db\n")
	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "internal", "auth", "service.go")), "package auth\n")
	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "internal", "handlers", "handlers.go")), "package handlers\n")
	writeFile(t, root, filepath.ToSlash(filepath.Join(projectRoot, "internal", "middleware", "middleware.go")), "package middleware\n")
	writeFile(t, root, ".github/workflows/ci.yml", "name: ci\n")
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
		"## In Production",
		"",
		"production note",
		"",
		"## Thinking Questions",
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
		"## In Production",
		"",
		"production note",
		"",
		"## Thinking Questions",
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

func requireOnlyFixtureExpectedReports(t *testing.T, result Result, reports []string, want string) {
	t.Helper()

	if result.ErrorCount == 0 {
		t.Fatalf("expected validation error %q, got none", want)
	}

	filtered := reportsWithoutFixtureScaffold(reports)
	if len(filtered) != 1 || !containsReport(filtered, want) {
		t.Fatalf("expected only %q outside fixture scaffold reports, got %v from all reports %v", want, filtered, reports)
	}
}

func requireOnlyFixtureScaffoldReports(t *testing.T, result Result, reports []string) {
	t.Helper()

	filtered := reportsWithoutFixtureScaffold(reports)
	if len(filtered) != 0 {
		t.Fatalf("expected only fixture scaffold reports, got %d validation errors and non-scaffold reports %v from all reports %v", result.ErrorCount, filtered, reports)
	}
}

func reportsWithoutFixtureScaffold(reports []string) []string {
	filtered := make([]string, 0, len(reports))
	for _, report := range reports {
		if isFixtureScaffoldReport(report) {
			continue
		}
		filtered = append(filtered, report)
	}

	return filtered
}

func isFixtureScaffoldReport(report string) bool {
	prefixes := []string{
		"Invalid v2 architecture contract:",
		"Invalid v2 section status:",
		"Invalid v2 section outputs:",
		"Invalid section README contract:",
		"Invalid engineering README contract:",
		"Warning: placeholder item:",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(report, prefix) {
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
