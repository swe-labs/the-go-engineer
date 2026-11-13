package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The fixture curriculum is deliberately tiny but exercises every manifest
// shape the validators read: two modules (one locked, one terminal), a lesson,
// an assessment item, a project record with no code artifacts, a concept, and
// a failure taxonomy. Every baseline fixture must pass every validator before
// its failure case is allowed to corrupt exactly one invariant.

const fixtureLastUpdated = "2026-01-01"

func fixtureHeader(documentType string) string {
	return `{
  "schema_version": "3.1.0",
  "document_type": "` + documentType + `",
  "curriculum_version": "3.0.0",
  "status": "stable",
  "last_updated": "` + fixtureLastUpdated + `"`
}

func writeCurriculumFixture(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "curriculum")

	files := map[string]string{
		"expected.json": fixtureHeader("expected_inventory") + `,
  "locked_modules": ["module-00"],
  "inventory": {
    "core_items": 2,
    "core_modules": 2,
    "elective_items": 0,
    "elective_modules": 0,
    "total_items": 2,
    "total_modules": 2,
    "assessments": 1,
    "projects": 1,
    "concepts": 1,
    "starter_parity_projects": 0,
    "opslane_milestones": 0,
    "path_project_items": 0,
    "path_assessment_items": 1,
    "shared_readme_bindings": 0
  }
}
`,
		"workspace.json": fixtureHeader("workspace_manifest") + `
}
`,
		"path.core.json": fixtureHeader("core_path") + `,
  "modules": [
    {
      "id": "module-00",
      "number": 0,
      "title": "Alpha Module",
      "phase": "foundation",
      "path": "00-alpha",
      "status": "stable",
      "learning_goal": "Learn the fixture concept.",
      "order": 0,
      "required": true,
      "portfolio_output": false,
      "prerequisites": [],
      "entry_item_ids": ["core-00-01"],
      "terminal_item_ids": ["core-00-01"],
      "tags": ["alpha"],
      "readme_contract": {
        "contract_id": "module.v3",
        "documentation_mode": "module-overview",
        "required": true
      },
      "cognitive_load": "low",
      "recommended_break_after": false,
      "contains_foundational_hard_concepts": false,
      "pacing": "gentle"
    },
    {
      "id": "module-01",
      "number": 1,
      "title": "Beta Module",
      "phase": "foundation",
      "path": "01-beta",
      "status": "stable",
      "learning_goal": "Finish the fixture path.",
      "order": 1,
      "required": true,
      "portfolio_output": false,
      "prerequisites": [],
      "entry_item_ids": [],
      "terminal_item_ids": [],
      "tags": ["beta"],
      "readme_contract": {
        "contract_id": "module.v3",
        "documentation_mode": "module-overview",
        "required": true
      },
      "cognitive_load": "moderate",
      "recommended_break_after": false,
      "contains_foundational_hard_concepts": false,
      "pacing": "steady"
    }
  ],
  "items": [
    {
      "content_contract": {
        "common_mistakes_required": false,
        "machine_view_required": false,
        "portfolio_artifact_required": false,
        "production_notes_required": false,
        "readme_required": true,
        "review_questions_required": false,
        "runnable_required": true,
        "tests_required": true,
        "visual_model_required": false
      },
      "crossrefs": {
        "builds_on": [],
        "preview_only": [],
        "related": [],
        "reinforced_in": []
      },
      "difficulty": "beginner",
      "estimated_minutes": 30,
      "files": {
        "assets_dir": "modules/00-alpha/lessons/01-alpha/assets",
        "diagram_paths": [],
        "main_path": "modules/00-alpha/lessons/01-alpha/main.go",
        "readme_path": "modules/00-alpha/lessons/01-alpha/README.md",
        "solution_path": "",
        "starter_path": "",
        "test_path": "modules/00-alpha/lessons/01-alpha/main_test.go"
      },
      "id": "core-00-01",
      "learning_objective": "Explain the fixture concept.",
      "module_id": "module-00",
      "next_item_ids": [],
      "order": 1,
      "phase": "foundation",
      "prerequisites": [],
      "proof": {
        "assessment_id": "assessment-module-00",
        "expected_artifact": "A short explanation of the fixture concept.",
        "mastery_checks": ["Can explain the fixture concept."],
        "practice_task": "Explain the fixture concept in one paragraph."
      },
      "readme_contract": {
        "contract_id": "lesson.v3",
        "documentation_mode": "tutorial",
        "required": true
      },
      "slug": "alpha-lesson",
      "source_legacy_ids": [],
      "status": "ready",
      "subtype": "concept",
      "tags": ["alpha", "concept"],
      "title": "Alpha lesson",
      "type": "lesson",
      "verification": {
        "expected_output": "Program prints a fixture greeting.",
        "manual_steps": ["Read the README."],
        "mode": "manual",
        "run_command": "go run ./main.go",
        "test_command": "go test ./..."
      },
      "zero_magic": {
        "beginner_mistakes": ["Assuming the fixture is production data."],
        "how_go_uses_it": "The fixture uses ordinary Go syntax.",
        "mental_model": "The fixture is a small, complete example.",
        "problem_solved": "The fixture concept removes ambiguity.",
        "proof_of_understanding": "The learner explains the fixture concept.",
        "real_world_usage": "Fixtures keep validation deterministic.",
        "under_the_hood": "The test writes files to a temporary directory.",
        "why_it_exists": "Validators need a small known-good corpus."
      },
      "zero_magic_status": "reviewed"
    },
    {
      "crossrefs": {
        "builds_on": [],
        "preview_only": [],
        "related": [],
        "reinforced_in": []
      },
      "difficulty": "beginner",
      "estimated_minutes": 30,
      "files": {
        "assets_dir": "",
        "diagram_paths": [],
        "main_path": "",
        "readme_path": "modules/00-alpha/assessments/checkpoint/README.md",
        "solution_path": "",
        "starter_path": "",
        "test_path": ""
      },
      "id": "core-00-02",
      "learning_objective": "Demonstrate the fixture concept.",
      "module_id": "module-00",
      "next_item_ids": [],
      "order": 2,
      "phase": "foundation",
      "prerequisites": [],
      "readme_contract": {
        "contract_id": "assessment.v3",
        "documentation_mode": "assessment",
        "required": true
      },
      "slug": "alpha-checkpoint",
      "source_legacy_ids": [],
      "status": "ready",
      "subtype": "assessment",
      "tags": ["alpha", "assessment"],
      "title": "Alpha checkpoint",
      "type": "assessment",
      "verification": {
        "expected_output": "Assessment evidence is complete.",
        "manual_steps": ["Submit the fixture artifact."],
        "mode": "manual"
      }
    }
  ]
}
`,
		"path.electives.json": fixtureHeader("elective_path") + `,
  "modules": [],
  "items": []
}
`,
		"projects.json": fixtureHeader("projects") + `,
  "projects": [
    {
      "id": "project-alpha",
      "module_id": "module-00",
      "slug": "alpha",
      "title": "Alpha project",
      "type": "project",
      "prerequisites": [],
      "deliverables": ["A fixture artifact."],
      "verification": {
        "run_command": "make build"
      },
      "rubric": {
        "criteria": []
      },
      "assessment_id": "assessment-module-00"
    }
  ]
}
`,
		"assessments.json": fixtureHeader("assessments") + `,
  "assessments": [
    {
      "id": "assessment-module-00",
      "module_id": "module-00",
      "title": "Alpha checkpoint",
      "type": "checkpoint",
      "criteria": [
        {
          "name": "Conceptual clarity",
          "weight": 100
        }
      ],
      "passing_score": 70,
      "target_ids": ["core-00-01"],
      "evidence_required": ["A working fixture artifact."],
      "retake_policy": "Retake after revision."
    }
  ]
}
`,
		"crossrefs.json": fixtureHeader("crossrefs") + `,
  "crossrefs": {
    "references": []
  }
}
`,
		"concepts.json": fixtureHeader("concept_ownership_registry") + `,
  "concepts": [
    {
      "concept": "fixture concept",
      "canonical_owner": "core-00-01",
      "preview_locations": [],
      "reinforcement_locations": []
    }
  ]
}
`,
		"failures.json": fixtureHeader("operational_failure_taxonomy") + `,
  "failure_categories": [],
  "required_coverage": {}
}
`,
		"readme.contracts.json": `{
  "schema_version": "3.1.0",
  "document_type": "readme_contracts",
  "curriculum_version": "3.0.0",
  "status": "stable",
  "last_updated": "` + fixtureLastUpdated + `",
  "contracts": [
    {
      "id": "root.v3",
      "applies_to": "repository-root",
      "purpose": "Fixture root contract.",
      "shape": "overview",
      "selector": {
        "doc": "repository-root"
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["Root overview"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "module.v3",
      "applies_to": "core-module",
      "purpose": "Fixture module contract.",
      "shape": "overview",
      "selector": {
        "doc": "module",
        "aliases": ["elective-module.v3"]
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["Overview"],
      "optional_headings": ["Extra"],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "lesson.v3",
      "applies_to": "core-lesson",
      "purpose": "Fixture lesson contract.",
      "shape": "standard-lesson",
      "selector": {
        "type": ["lesson"],
        "subtype_not": ["reference"]
      },
      "code_policy": "optional",
      "file_policy": {
        "main": "when:runnable_required",
        "readme": "required",
        "test": "when:tests_required"
      },
      "required_headings": ["Learning objective"],
      "optional_headings": ["Visual model"],
      "heading_aliases": {},
      "subtype_modifiers": {
        "concept": {
          "description": "Core concept lesson."
        }
      }
    },
    {
      "id": "elective-lesson.v3",
      "applies_to": "elective-lesson",
      "purpose": "Fixture elective-lesson contract.",
      "shape": "reference",
      "selector": {
        "type": ["lesson"],
        "subtype": ["reference"]
      },
      "code_policy": "required",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["When you need this"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "project.v3",
      "applies_to": "project",
      "purpose": "Fixture project contract.",
      "shape": "project",
      "selector": {
        "type": ["project"],
        "subtype_not": ["flagship-checkpoint"]
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required",
        "solution": "when:portfolio_artifact_required",
        "starter": "when:portfolio_artifact_required"
      },
      "required_headings": ["Project overview"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "assessment.v3",
      "applies_to": "assessment",
      "purpose": "Fixture assessment contract.",
      "shape": "assessment",
      "selector": {
        "type": ["assessment"],
        "subtype_not": ["milestone-assessment"]
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["Purpose"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "flagship.v3",
      "applies_to": "flagship",
      "purpose": "Fixture flagship contract.",
      "shape": "flagship",
      "selector": {
        "doc": "flagship-module"
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["Flagship mission"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "flagship-checkpoint.v3",
      "applies_to": "flagship-checkpoint",
      "purpose": "Fixture flagship-checkpoint contract.",
      "shape": "flagship-checkpoint",
      "selector": {
        "type": ["checkpoint"],
        "subtype": ["flagship-checkpoint"]
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["Mission"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    },
    {
      "id": "companion-stage.v3",
      "applies_to": "companion-stage",
      "purpose": "Fixture companion-stage contract.",
      "shape": "companion",
      "selector": {
        "doc": "companion-stage"
      },
      "code_policy": "optional",
      "file_policy": {
        "readme": "required"
      },
      "required_headings": ["Stage mission"],
      "optional_headings": [],
      "heading_aliases": {},
      "subtype_modifiers": {}
    }
  ],
  "severity_policy": {
    "description": "Tiered severity for fixture tests.",
    "default_content_shape": "warning",
    "content_shape_enforcement": "enforce",
    "by_status": {
      "draft": "skip",
      "review": "warning",
      "ready": "error"
    },
    "rollout_note": "Fixture policy: skip, warn, and block tiers are active.",
    "always_error": ["duplicate-id"],
    "locked_modules_ref": "curriculum/expected.json#/locked_modules",
    "locked_modules_note": "Fixture locked modules live in expected.json."
  },
  "validation_rules": [
    "The generated catalog must be wrapped in curriculum:generated:start and curriculum:generated:end markers.",
    "Structural defects such as duplicate-id always block."
  ]
}
`,
		"schema.v3.json": `{
  "$id": "https://example.invalid/fixture.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Fixture curriculum schema"
}
`,
		"README.md": `<!-- curriculum:generated:start -->
# Fixture Curriculum

## Root overview

Fixture root for validator tests.
<!-- curriculum:generated:end -->
`,
		"modules/00-alpha/README.md": `# Alpha Module

## Overview

The locked fixture module.

[Assessment](assessments/checkpoint/README.md)

## Start here

Start with the lesson.
`,
		"modules/01-beta/README.md": `# Beta Module

## Overview

The terminal fixture module.

[Catalog](../../README.md)

## Completion

The fixture path is complete.
`,
		"modules/00-alpha/lessons/01-alpha/README.md": "# Alpha lesson\n\n" +
			"> \"Simplicity is prerequisite for reliability.\" — Edsger W. Dijkstra\n\n" +
			"## Learning objective\n\nExplain the fixture concept.\n\n" +
			"```go\npackage main\n\nfunc main() {}\n```\n\n" +
			"## NEXT UP\n\nThis is the terminal lesson in the module.\n",
		"modules/00-alpha/lessons/01-alpha/main.go":      "package main\n\nfunc main() {}\n",
		"modules/00-alpha/lessons/01-alpha/main_test.go": "package main\n\nimport \"testing\"\n\nfunc TestAlpha(t *testing.T) {}\n",
		"modules/00-alpha/assessments/checkpoint/README.md": "# Alpha checkpoint\n\n" +
			"## Purpose\n\n" +
			"Verify that the learner can explain the fixture concept and show the artifact. " +
			"The next module link is `../01-beta/README.md` so the exit gate can resolve it. " +
			"This paragraph exists to keep the assessment README above the minimum size.\n",
	}

	for rel, content := range files {
		full := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("create fixture directory for %s: %v", rel, err)
		}
		writeFixture(t, full, content)
	}

	// The repository-root contract targets the README beside curriculum/, not
	// the generated catalog inside it.
	writeFixture(t, filepath.Join(filepath.Dir(dir), "README.md"), `# Fixture Repository

## Root overview

The repository-level fixture README.
`)

	for _, rel := range []string{
		"modules/00-alpha/lessons/01-alpha/assets",
		"modules/00-alpha/projects/alpha",
	} {
		if err := os.MkdirAll(filepath.Join(dir, filepath.FromSlash(rel)), 0o755); err != nil {
			t.Fatalf("create fixture directory %s: %v", rel, err)
		}
	}
	return dir
}

// installFixture points the package globals at dir and restores every cache and
// counter after the subtest. The validators are package-level functions with
// package-level state, so isolation is explicit rather than implicit.
func installFixture(t *testing.T, dir string) {
	t.Helper()
	prevDir := curriculumDir
	prevExpected := expectedCache
	prevRegistry := registryCache
	prevCoreBundle := coreBundleCache
	prevElectiveBundle := electiveBundleCache
	prevStrict := strictMode
	prevAdvisory := advisoryMode
	prevReport := reportPath
	prevErrors := errorCount
	prevWarnings := warningCount
	prevRecords := validationRecords
	prevValidator := currentValidator

	curriculumDir = dir
	expectedCache = nil
	registryCache = nil
	coreBundleCache = nil
	electiveBundleCache = nil
	strictMode = false
	advisoryMode = false
	reportPath = ""
	errorCount = 0
	warningCount = 0
	validationRecords = nil
	currentValidator = ""

	t.Cleanup(func() {
		curriculumDir = prevDir
		expectedCache = prevExpected
		registryCache = prevRegistry
		coreBundleCache = prevCoreBundle
		electiveBundleCache = prevElectiveBundle
		strictMode = prevStrict
		advisoryMode = prevAdvisory
		reportPath = prevReport
		errorCount = prevErrors
		warningCount = prevWarnings
		validationRecords = prevRecords
		currentValidator = prevValidator
	})
}

func replaceFixture(t *testing.T, dir, rel, old, new string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(rel))
	data, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("read fixture %s: %v", rel, err)
	}
	count := strings.Count(string(data), old)
	if count != 1 {
		t.Fatalf("fixture %s: expected exactly one occurrence of %q, found %d", rel, old, count)
	}
	if err := os.WriteFile(full, []byte(strings.Replace(string(data), old, new, 1)), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", rel, err)
	}
}

func appendFixture(t *testing.T, dir, rel, extra string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(rel))
	data, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("read fixture %s: %v", rel, err)
	}
	if err := os.WriteFile(full, append(data, []byte(extra)...), 0o644); err != nil {
		t.Fatalf("append fixture %s: %v", rel, err)
	}
}

func removeFixture(t *testing.T, dir, rel string) {
	t.Helper()
	if err := os.Remove(filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
		t.Fatalf("remove fixture %s: %v", rel, err)
	}
}

func mkdirFixture(t *testing.T, dir, rel string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, filepath.FromSlash(rel)), 0o755); err != nil {
		t.Fatalf("create fixture directory %s: %v", rel, err)
	}
}

// assertValidatorClean runs one validator against the untouched fixture.
func assertValidatorClean(t *testing.T, run func()) {
	t.Helper()
	installFixture(t, writeCurriculumFixture(t))
	run()
	if errorCount != 0 || warningCount != 0 {
		t.Fatalf("baseline fixture failed: %d error(s), %d warning(s)", errorCount, warningCount)
	}
}

// assertValidatorFails runs one validator against a fixture with exactly one
// targeted mutation and requires at least one diagnostic.
func assertValidatorFails(t *testing.T, run func(), mutate func(*testing.T, string)) {
	t.Helper()
	dir := writeCurriculumFixture(t)
	mutate(t, dir)
	installFixture(t, dir)
	run()
	if errorCount+warningCount == 0 {
		t.Fatalf("targeted mutation produced no diagnostic")
	}
}

// failureMutations is keyed by validator name. The test below requires one
// entry per registered validator, so adding a gate without an honest failing
// fixture is a test failure rather than an untested branch.
var failureMutations = map[string]func(*testing.T, string){
	"validate-schema": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `"title": "Alpha lesson",`, `"title": "",`)
	},
	"validate-graph": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "crossrefs.json", `"references": []`,
			`"references": [{"from_id": "core-00-01", "target_id": "core-99-99", "relation": "related"}]`)
	},
	"validate-lessons": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md", "## NEXT UP", "## Next steps")
	},
	"validate-prerequisite-closure": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `      "phase": "foundation",
      "prerequisites": [],
      "proof": {`, `      "phase": "foundation",
      "prerequisites": ["core-99-99"],
      "proof": {`)
	},
	"validate-project-binding": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "projects.json", `"module_id": "module-00",`, `"module_id": "module-99",`)
	},
	"validate-no-placeholder-zm": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `"problem_solved": "The fixture concept removes ambiguity.",`,
			`"problem_solved": "This lesson explains what problem it solves.",`)
	},
	"validate-assessment-completeness": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "assessments.json", `"passing_score": 70,`, `"passing_score": 0,`)
	},
	"validate-unique-tags": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `"tags": ["alpha", "concept"],`, `"tags": ["alpha", "alpha"],`)
	},
	"validate-crossrefs": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "crossrefs.json", `"references": []`,
			`"references": [{"from_id": "core-00-01", "target_id": "core-99-99", "relation": "related"}]`)
	},
	"validate-concept-ownership": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "concepts.json", `"canonical_owner": "core-00-01"`, `"canonical_owner": "core-99-99"`)
	},
	"validate-cognitive-load": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `"cognitive_load": "low"`, `"cognitive_load": "extreme"`)
	},
	"validate-operational-failures": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "failures.json", `"document_type": "operational_failure_taxonomy"`, `"document_type": "other"`)
	},
	"validate-failure-engineering": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `      "learning_objective": "Explain the fixture concept.",
      "module_id": "module-00",`, `      "learning_objective": "Explain the fixture concept.",
      "module_id": "module-08",`)
	},
	"validate-zero-magic": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "path.core.json", `"zero_magic_status": "reviewed"`, `"zero_magic_status": "golden"`)
	},
	"validate-physical-artifacts": func(t *testing.T, dir string) {
		removeFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md")
	},
	"validate-assessment-artifacts": func(t *testing.T, dir string) {
		removeFixture(t, dir, "modules/00-alpha/assessments/checkpoint/README.md")
	},
	"validate-starter-parity": func(t *testing.T, dir string) {
		mkdirFixture(t, dir, "modules/00-alpha/projects/alpha/_starter")
	},
	"validate-markdown-links": func(t *testing.T, dir string) {
		appendFixture(t, dir, "README.md", "\n[Broken](missing.md)\n")
	},
	"validate-inline-paths": func(t *testing.T, dir string) {
		appendFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md",
			"\nSee `curriculum/modules/00-alpha/lessons/01-alpha/missing-artifact.md` for details.\n")
	},
	"validate-sese-navigation": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "modules/00-alpha/README.md", "## Start here", "## Wrap up")
	},
	"validate-project-scripts": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "projects.json", `"run_command": "make build"`, `"run_command": "bash verify.sh"`)
	},
	"validate-no-emojis": func(t *testing.T, dir string) {
		appendFixture(t, dir, "README.md", "\nRocket: \U0001F680\n")
	},
	"validate-anchor-quotes": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md",
			"> \"Simplicity is prerequisite for reliability.\" — Edsger W. Dijkstra\n\n", "")
	},
	"validate-diagram-palette": func(t *testing.T, dir string) {
		appendFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md",
			"\n## Visual model\n\n```mermaid\nflowchart TD\n    A[\"A\"] --> B[\"B\"]\n```\n")
	},
	"validate-readme-contracts": func(t *testing.T, dir string) {
		replaceFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md", "## Learning objective", "## Goal")
	},
	"validate-anti-slop": func(t *testing.T, dir string) {
		appendFixture(t, dir, "modules/00-alpha/lessons/01-alpha/README.md",
			"\nIn simple terms, this is a cutting-edge fixture.\n")
	},
}

func TestValidatorFixtures(t *testing.T) {
	validatorsByName := make(map[string]func(), len(validators()))
	for _, v := range validators() {
		validatorsByName[v.Name] = v.Run
	}
	for name := range failureMutations {
		if _, ok := validatorsByName[name]; !ok {
			t.Errorf("failure mutation %q does not match a registered validator", name)
		}
	}
	for name := range validatorsByName {
		if _, ok := failureMutations[name]; !ok {
			t.Errorf("validator %q has no targeted failure mutation", name)
		}
	}

	for _, v := range validators() {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			t.Run("baseline", func(t *testing.T) {
				assertValidatorClean(t, v.Run)
			})
			t.Run("failure", func(t *testing.T) {
				assertValidatorFails(t, v.Run, failureMutations[v.Name])
			})
		})
	}
}
