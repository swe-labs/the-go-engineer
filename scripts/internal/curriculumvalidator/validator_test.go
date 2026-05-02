package curriculumvalidator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateAcceptsValidV2Fixture(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)

	result, reports := runValidate(t, root)
	if result.ErrorCount != 0 {
		t.Fatalf("Validate returned %d errors, reports: %v", result.ErrorCount, reports)
	}
	if !result.HasV2 {
		t.Fatalf("Validate did not detect curriculum.v2.json")
	}
	if result.V2SectionCount != len(canonicalV2Sections) {
		t.Fatalf("V2SectionCount = %d, want %d", result.V2SectionCount, len(canonicalV2Sections))
	}
	if result.V2ItemCount != len(cur.Items) {
		t.Fatalf("V2ItemCount = %d, want %d", result.V2ItemCount, len(cur.Items))
	}
	if result.PlaceholderCount != 0 {
		t.Fatalf("PlaceholderCount = %d, want 0", result.PlaceholderCount)
	}
	if result.FilesScanned == 0 {
		t.Fatalf("FilesScanned = 0, want validator to scan fixture run surfaces")
	}
}

func TestValidateRejectsArchitectureDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*V2Curriculum)
		want   string
	}{
		{
			name: "wrong section order",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0], cur.Sections[1] = cur.Sections[1], cur.Sections[0]
			},
			want: "Invalid v2 architecture contract: section position 0 -> s01 (expected s00)",
		},
		{
			name: "wrong section number",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Number = "99"
			},
			want: "Invalid v2 section number: s00 -> 99 (expected 00)",
		},
		{
			name: "wrong section slug",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Slug = "machine-basics"
			},
			want: "Invalid v2 section slug: s00 -> machine-basics (expected how-computers-work)",
		},
		{
			name: "wrong section title",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Title = "Machines"
			},
			want: "Invalid v2 section title: s00 -> Machines (expected How Computers Work)",
		},
		{
			name: "wrong path prefix",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].PathPrefix = "00-machines"
			},
			want: "Invalid v2 section path_prefix: s00 -> 00-machines (expected 00-how-computers-work)",
		},
		{
			name: "wrong status",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Status = "draft"
			},
			want: "Invalid v2 section status: s00 -> draft",
		},
		{
			name: "wrong phase",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Phase = "legacy"
			},
			want: "Invalid v2 section phase: s00 -> legacy",
		},
		{
			name: "missing summary",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Summary = ""
			},
			want: "Invalid v2 section metadata: s00 requires number, slug, title, path_prefix, phase, and summary",
		},
		{
			name: "wrong entry points",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].EntryPoints = []string{"HC.2"}
			},
			want: "Invalid v2 section entry points: s00 -> HC.2 (expected HC.1)",
		},
		{
			name: "wrong outputs",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].Outputs = []string{"HC.4"}
			},
			want: "Invalid v2 section outputs: s00 -> HC.4 (expected HC.5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			cur := writeValidV2Fixture(t, root)
			tt.mutate(&cur)
			writeCurriculum(t, root, cur)

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected validation error containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsSchemaAndMetadataDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*V2Curriculum)
		want   string
	}{
		{
			name: "invalid schema version",
			mutate: func(cur *V2Curriculum) {
				cur.SchemaVersion = 2
			},
			want: "Invalid v2 schema_version: 2 (expected 1)",
		},
		{
			name: "duplicate section id",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[1].ID = cur.Sections[0].ID
			},
			want: "Duplicate v2 section id: s00",
		},
		{
			name: "duplicate item id",
			mutate: func(cur *V2Curriculum) {
				cur.Items[1].ID = cur.Items[0].ID
			},
			want: "Duplicate v2 item id: HC.1",
		},
		{
			name: "invalid section id format",
			mutate: func(cur *V2Curriculum) {
				cur.Sections[0].ID = "section-00"
			},
			want: "Invalid v2 section id format: section-00",
		},
		{
			name: "invalid item id format",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].ID = "HC-one"
			},
			want: "Invalid v2 item id format: HC-one",
		},
		{
			name: "invalid item type",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Type = "lessen"
			},
			want: "Invalid v2 item type: HC.1 -> lessen",
		},
		{
			name: "invalid item level",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Level = "advanced"
			},
			want: "Invalid v2 item level: HC.1 -> advanced",
		},
		{
			name: "invalid verification mode",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].VerificationMode = "manual"
			},
			want: "Invalid v2 verification mode: HC.1 -> manual",
		},
		{
			name: "invalid lesson subtype",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Subtype = "syntax"
			},
			want: "Invalid v2 lesson subtype: HC.1 -> syntax",
		},
		{
			name: "unexpected non-lesson subtype",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Type = "exercise"
				cur.Items[0].Subtype = "concept"
			},
			want: "Unexpected v2 subtype for non-lesson item: HC.1 -> concept",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			cur := writeValidV2Fixture(t, root)
			tt.mutate(&cur)
			writeCurriculum(t, root, cur)

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected validation error containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsCurriculumJSONContractDrift(t *testing.T) {
	tests := []struct {
		name string
		edit func(t *testing.T, root string)
		want string
	}{
		{
			name: "non canonical formatting",
			edit: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "curriculum.v2.json")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read curriculum fixture: %v", err)
				}
				edited := strings.Replace(string(data), "  \"schema_version\"", "    \"schema_version\"", 1)
				writeFile(t, root, "curriculum.v2.json", edited)
			},
			want: "Invalid v2 curriculum JSON formatting:",
		},
		{
			name: "duplicate object key",
			edit: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "curriculum.v2.json")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read curriculum fixture: %v", err)
				}
				edited := strings.Replace(string(data), "\"schema_version\": 1,", "\"schema_version\": 1,\n  \"schema_version\": 1,", 1)
				writeFile(t, root, "curriculum.v2.json", edited)
			},
			want: "Duplicate JSON key in curriculum.v2.json: $.schema_version",
		},
		{
			name: "unknown field",
			edit: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "curriculum.v2.json")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read curriculum fixture: %v", err)
				}
				edited := strings.Replace(string(data), "\"schema_version\": 1,", "\"schema_version\": 1,\n  \"legacy\": true,", 1)
				writeFile(t, root, "curriculum.v2.json", edited)
			},
			want: "Invalid v2 curriculum JSON object: $ has unknown field legacy",
		},
		{
			name: "null array",
			edit: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "curriculum.v2.json")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read curriculum fixture: %v", err)
				}
				edited := strings.Replace(string(data), "\"prerequisites\": [],", "\"prerequisites\": null,", 1)
				writeFile(t, root, "curriculum.v2.json", edited)
			},
			want: "Invalid v2 curriculum JSON array: $.items[0].prerequisites must be [] instead of null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeValidV2Fixture(t, root)
			tt.edit(t, root)

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected curriculum JSON contract report containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsCurriculumOrderingAndReferenceDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*V2Curriculum)
		want   string
	}{
		{
			name: "out of order items",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0], cur.Items[1] = cur.Items[1], cur.Items[0]
			},
			want: "Invalid v2 item order:",
		},
		{
			name: "duplicate next item",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].NextItems = []string{"HC.5", "HC.5"}
			},
			want: "Invalid v2 item HC.1 next_items: duplicate value HC.5",
		},
		{
			name: "self prerequisite",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Prerequisites = []string{"HC.1"}
			},
			want: "Invalid v2 prerequisite: HC.1 cannot reference itself",
		},
		{
			name: "invalid slug format",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Slug = "What Is A Program"
			},
			want: "Invalid v2 item slug format: HC.1 -> What Is A Program",
		},
		{
			name: "path outside section prefix",
			mutate: func(cur *V2Curriculum) {
				cur.Items[0].Path = "01-getting-started/001-hc-1"
				cur.Items[0].RunCommand = "go run ./01-getting-started/001-hc-1"
			},
			want: "Invalid v2 section path alignment: HC.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			cur := writeValidV2Fixture(t, root)
			tt.mutate(&cur)
			writeCurriculum(t, root, cur)

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected curriculum ordering/reference report containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsRootEscapingPaths(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(t *testing.T, root string, cur *V2Curriculum)
		want   string
	}{
		{
			name: "item path escapes root",
			mutate: func(t *testing.T, root string, cur *V2Curriculum) {
				cur.Items[0].Path = "../outside"
			},
			want: "Invalid v2 item path: HC.1 -> ../outside",
		},
		{
			name: "windows absolute item path is rejected cross platform",
			mutate: func(t *testing.T, root string, cur *V2Curriculum) {
				cur.Items[0].Path = "C:/outside"
			},
			want: "Invalid v2 item path: HC.1 -> C:/outside",
		},
		{
			name: "backslash item path is rejected cross platform",
			mutate: func(t *testing.T, root string, cur *V2Curriculum) {
				cur.Items[0].Path = `00-how-computers-work\001-hc-1`
			},
			want: `Invalid v2 item path: HC.1 -> 00-how-computers-work\001-hc-1`,
		},
		{
			name: "starter path escapes root",
			mutate: func(t *testing.T, root string, cur *V2Curriculum) {
				cur.Items[0].StarterPath = "../starter"
			},
			want: "Invalid v2 starter path: HC.1 -> ../starter",
		},
		{
			name: "run command target escapes root",
			mutate: func(t *testing.T, root string, cur *V2Curriculum) {
				cur.Items[0].RunCommand = "go run ./../outside"
			},
			want: "Invalid v2 run command target: HC.1 -> go run ./../outside",
		},
		{
			name: "markdown link escapes root",
			mutate: func(t *testing.T, root string, cur *V2Curriculum) {
				writeFile(t, root, filepath.Join(cur.Items[0].Path, "README.md"), validLessonReadme(cur.Items[0].RunCommand)+"\n[escape](../../../outside.md)\n")
			},
			want: "Broken local doc link:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			cur := writeValidV2Fixture(t, root)
			tt.mutate(t, root, &cur)
			writeCurriculum(t, root, cur)

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected validation error containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsFakeHeadingsInCodeFence(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)
	readme := strings.Join([]string{
		"# Fake headings",
		"",
		"```markdown",
		"## Mission",
		"## Prerequisites",
		"## Mental Model",
		"## Visual Model",
		"## Machine View",
		"## Run Instructions",
		"## Code Walkthrough",
		"## Try It",
		"## In Production",
		"## Thinking Questions",
		"## Next Step",
		"```",
	}, "\n")
	writeFile(t, root, filepath.Join(cur.Items[0].Path, "README.md"), readme)

	result, reports := runValidate(t, root)
	if result.ErrorCount == 0 {
		t.Fatalf("expected README contract errors")
	}
	requireReportContains(t, reports, "Invalid foundations README contract: HC.1")
	requireReportContains(t, reports, "missing ## Mission")
}

func TestValidateRejectsLessonSourceHeaderLevelDrift(t *testing.T) {
	t.Run("missing level header", func(t *testing.T) {
		root := t.TempDir()
		cur := writeValidV2Fixture(t, root)
		item := requireItem(t, cur, "HC.1")
		writeFile(t, root, filepath.Join(item.Path, "main.go"), "package main\n\nfunc main() {}\n")

		result, reports := runValidate(t, root)
		if result.ErrorCount == 0 {
			t.Fatalf("expected missing lesson level header report")
		}
		requireReportContains(t, reports, "Missing v2 lesson level header: HC.1")
	})

	t.Run("wrong level header", func(t *testing.T) {
		root := t.TempDir()
		cur := writeValidV2Fixture(t, root)
		item := requireItem(t, cur, "PD.1")
		writeFile(t, root, filepath.Join(item.Path, "main.go"), mainGoForItemWithLevel(item, "Core"))

		result, reports := runValidate(t, root)
		if result.ErrorCount == 0 {
			t.Fatalf("expected invalid lesson level header report")
		}
		requireReportContains(t, reports, "Invalid v2 lesson level header: PD.1")
		requireReportContains(t, reports, "expected Level: Production")
	})

	t.Run("wrong run header", func(t *testing.T) {
		root := t.TempDir()
		cur := writeValidV2Fixture(t, root)
		item := requireItem(t, cur, "HC.1")
		source := strings.Replace(mainGoForItem(item), "// RUN: "+item.RunCommand, "// RUN: go run ./wrong/path", 1)
		writeFile(t, root, filepath.Join(item.Path, "main.go"), source)

		result, reports := runValidate(t, root)
		if result.ErrorCount == 0 {
			t.Fatalf("expected invalid lesson run header report")
		}
		requireReportContains(t, reports, "Invalid v2 lesson run header: HC.1")
		requireReportContains(t, reports, "expected RUN: "+item.RunCommand)
	})
}

func TestValidateRejectsLegacyReferenceLabels(t *testing.T) {
	tests := []struct {
		name   string
		label  string
		report string
	}{
		{
			name:   "forward reference label",
			label:  "> **Forward Reference:** Continue with the next lesson.",
			report: "uses Forward Reference label",
		},
		{
			name:   "backward reference label",
			label:  "> **Backward Reference:** Remember the previous lesson.",
			report: "uses Backward Reference label",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			cur := writeValidV2Fixture(t, root)
			item := requireItem(t, cur, "HC.1")
			writeFile(t, root, filepath.Join(item.Path, "README.md"), validLessonReadme(item.RunCommand)+"\n"+tt.label+"\n")

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected legacy reference label report")
			}
			requireReportContains(t, reports, "Invalid markdown cross-reference alert:")
			requireReportContains(t, reports, tt.report)
		})
	}
}

func TestValidateRejectsStandardsDrift(t *testing.T) {
	tests := []struct {
		name string
		edit func(string) string
		want string
	}{
		{
			name: "code standards level taxonomy omits production",
			edit: func(text string) string {
				return strings.Replace(text, "Level: Foundation | Core | Production | Stretch", "Level: Foundation | Core | Stretch", 1)
			},
			want: "Invalid level taxonomy:",
		},
		{
			name: "verification command uses stale coverage flag spelling",
			edit: func(text string) string {
				return strings.Replace(text, "go test -coverprofile=coverage.out ./...", "go test -coverprofile coverage.out ./...", 1)
			},
			want: "Invalid verification command:",
		},
		{
			name: "public code standards reference maintainer-only agents doc",
			edit: func(text string) string {
				return text + "\nSee AGENTS.md for details.\n"
			},
			want: "references maintainer-only AGENTS.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeValidV2Fixture(t, root)
			writeFile(t, root, "CODE-STANDARDS.md", tt.edit(validCodeStandards()))

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected standards drift report containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsMarkdownCrossReferenceContractDrift(t *testing.T) {
	tests := []struct {
		name  string
		block string
		want  string
	}{
		{
			name: "curriculum reference uses unsupported alert type",
			block: strings.Join([]string{
				"> [!WARNING]",
				"> Review [HC.1 What is a Program?](./README.md) before continuing.",
			}, "\n"),
			want: "uses [!WARNING] for curriculum reference",
		},
		{
			name: "curriculum reference lacks README link",
			block: strings.Join([]string{
				"> [!NOTE]",
				"> HC.1 introduces this idea.",
			}, "\n"),
			want: "references a curriculum ID without a clickable README.md link",
		},
		{
			name:  "legacy cross-reference heading",
			block: "## Forward Reference\n\nContinue later.",
			want:  "Invalid markdown cross-reference heading:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			cur := writeValidV2Fixture(t, root)
			item := requireItem(t, cur, "HC.1")
			writeFile(t, root, filepath.Join(item.Path, "README.md"), validLessonReadme(item.RunCommand)+"\n"+tt.block+"\n")

			result, reports := runValidate(t, root)
			if result.ErrorCount == 0 {
				t.Fatalf("expected cross-reference contract report containing %q", tt.want)
			}
			requireReportContains(t, reports, tt.want)
		})
	}
}

func TestValidateRejectsGoSourceMarkdownAlertSyntax(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)
	item := requireItem(t, cur, "HC.1")
	source := mainGoForItem(item) + "\n// [!NOTE] Markdown alerts belong in README files.\n"
	writeFile(t, root, filepath.Join(item.Path, "main.go"), source)

	result, reports := runValidate(t, root)
	if result.ErrorCount == 0 {
		t.Fatalf("expected Go source markdown alert report")
	}
	requireReportContains(t, reports, "Invalid Go source cross-reference comment: ")
}

func TestValidateRejectsMissingMachineRoleComment(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)
	item := requireItem(t, cur, "HC.1")
	source := mainGoForItem(item) + "\nfunc helper() {}\n"
	writeFile(t, root, filepath.Join(item.Path, "main.go"), source)

	result, reports := runValidate(t, root)
	if result.ErrorCount == 0 {
		t.Fatalf("expected missing Machine Role comment report")
	}
	requireReportContains(t, reports, "Missing Machine Role comment:")
	requireReportContains(t, reports, "helper")
}

func TestValidateReportsBrokenLocalLinkOnlyOnce(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)
	writeFile(t, root, filepath.Join(cur.Items[0].Path, "README.md"), validLessonReadme(cur.Items[0].RunCommand)+"\n[missing](./missing.md)\n")

	result, reports := runValidate(t, root)
	if result.ErrorCount == 0 {
		t.Fatalf("expected broken local link report")
	}

	count := countReportsContaining(reports, "Broken local doc link:")
	if count != 1 {
		t.Fatalf("broken local link reports = %d, want 1; reports: %v", count, reports)
	}
}

func TestValidateHandlesRubricModeExplicitly(t *testing.T) {
	t.Run("accepts rubric with no commands when README exists", func(t *testing.T) {
		root := t.TempDir()
		cur := writeValidV2Fixture(t, root)

		result, reports := runValidate(t, root)
		if result.ErrorCount != 0 {
			t.Fatalf("Validate returned %d errors, reports: %v", result.ErrorCount, reports)
		}
		requireItem(t, cur, "OPSL.10")
	})

	t.Run("rejects rubric with no README", func(t *testing.T) {
		root := t.TempDir()
		cur := writeValidV2Fixture(t, root)
		item := requireItem(t, cur, "OPSL.10")
		if err := os.Remove(filepath.Join(root, item.Path, "README.md")); err != nil {
			t.Fatalf("remove README: %v", err)
		}

		result, reports := runValidate(t, root)
		if result.ErrorCount == 0 {
			t.Fatalf("expected rubric README proof-surface error")
		}
		requireReportContains(t, reports, "Invalid v2 rubric contract: OPSL.10 requires README proof surface")
	})
}

func TestValidateRunPathScanStillReportsInvalidCommand(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)
	writeFile(t, root, filepath.Join(cur.Items[0].Path, "main.go"), "package main\n\nfunc main() {\n\tprintln(\"go run ./missing/path\")\n}\n")

	result, reports := runValidate(t, root)
	if result.ErrorCount == 0 {
		t.Fatalf("expected invalid run path report")
	}
	requireReportContains(t, reports, "Invalid run path:")
}

func TestValidateRejectsFlagshipChainBreak(t *testing.T) {
	root := t.TempDir()
	cur := writeValidV2Fixture(t, root)
	for i := range cur.Items {
		if cur.Items[i].ID == "OPSL.2" {
			cur.Items[i].NextItems = []string{"OPSL.4"}
			break
		}
	}
	writeCurriculum(t, root, cur)

	result, reports := runValidate(t, root)
	if result.ErrorCount == 0 {
		t.Fatalf("expected flagship chain error")
	}
	requireReportContains(t, reports, "Invalid flagship module chain: OPSL.2 must point to OPSL.3")
}

func writeValidV2Fixture(t *testing.T, root string) V2Curriculum {
	t.Helper()

	writeFile(t, root, "CODE-STANDARDS.md", validCodeStandards())

	cur := V2Curriculum{
		SchemaVersion: expectedSchemaVersion,
		Sections:      make([]V2Section, 0, len(canonicalV2Sections)),
	}

	for idx, section := range canonicalV2Sections {
		cur.Sections = append(cur.Sections, V2Section{
			ID:            section.ID,
			Number:        section.Number,
			Slug:          section.Slug,
			Title:         section.Title,
			PathPrefix:    section.PathPrefix,
			Phase:         section.Phase,
			Summary:       "Summary for " + section.ID,
			Status:        section.Status,
			EntryPoints:   append([]string(nil), section.EntryPoints...),
			Outputs:       append([]string(nil), section.Outputs...),
			Prerequisites: sectionPrerequisites(idx),
		})
		mustMkdir(t, root, section.PathPrefix)
		writeFile(t, root, filepath.Join(section.PathPrefix, "README.md"), sectionReadme(section.ID))
	}

	itemIDs := orderedFixtureItemIDs()
	for _, id := range itemIDs {
		item := fixtureItem(id)
		cur.Items = append(cur.Items, item)
		writeItemSurface(t, root, item)
	}

	writeFile(t, root, "11-flagship/01-opslane/README.md", "# Opslane\n")
	writeFile(t, root, "11-flagship/01-opslane/MODULES.md", "# Modules\n")
	writeFile(t, root, "11-flagship/01-opslane/scripts/progress.go", "package main\n\nfunc main() {}\n")

	writeCurriculum(t, root, cur)
	return cur
}

func orderedFixtureItemIDs() []string {
	seen := map[string]bool{}
	var ids []string
	add := func(id string) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}

	for _, section := range canonicalV2Sections {
		if section.ID == "s11" {
			continue
		}
		for _, id := range section.EntryPoints {
			add(id)
		}
		for _, id := range section.Outputs {
			add(id)
		}
	}
	for i := 1; i <= 10; i++ {
		add(fmt.Sprintf("OPSL.%d", i))
	}

	return ids
}

func fixtureItem(id string) V2Item {
	sectionID := sectionIDForItemID(id)
	path := fixtureItemPath(sectionID, id)
	item := V2Item{
		ID:               id,
		SectionID:        sectionID,
		Slug:             strings.ToLower(strings.ReplaceAll(id, ".", "-")),
		Title:            id + " Fixture",
		Type:             "lesson",
		Subtype:          "concept",
		Level:            "core",
		Status:           "implemented",
		VerificationMode: "run",
		Path:             path,
		Prerequisites:    []string{},
		RunCommand:       "go run ./" + path,
		NextItems:        []string{},
	}

	if isFoundationsSection(sectionID) {
		item.Level = "foundation"
	}
	if sectionID == "s09" || sectionID == "s10" {
		item.Level = "production"
	}
	if sectionID == "s11" {
		item.Type = "checkpoint"
		item.Subtype = ""
		item.Level = "stretch"
		item.VerificationMode = "rubric"
		item.RunCommand = ""
		item.TestCommand = ""
		number := itemNumber(id)
		if number < 10 {
			item.NextItems = []string{fmt.Sprintf("OPSL.%d", number+1)}
		} else {
			item.Type = "capstone"
		}
	}

	return item
}

func writeItemSurface(t *testing.T, root string, item V2Item) {
	t.Helper()

	mustMkdir(t, root, item.Path)
	if item.SectionID == "s11" {
		writeFile(t, root, filepath.Join(item.Path, "README.md"), flagshipReadme(item))
		return
	}

	writeFile(t, root, filepath.Join(item.Path, "README.md"), validLessonReadme(item.RunCommand))
	writeFile(t, root, filepath.Join(item.Path, "main.go"), mainGoForItem(item))
}

func validLessonReadme(runCommand string) string {
	return strings.Join([]string{
		"# Lesson",
		"",
		"## Mission",
		"Mission text.",
		"",
		"## Prerequisites",
		"None.",
		"",
		"## Mental Model",
		"Model.",
		"",
		"## Visual Model",
		"```mermaid",
		"graph TD",
		"    A[\"input\"] --> B[\"program\"]",
		"```",
		"",
		"## Machine View",
		"Machine.",
		"",
		"## Run Instructions",
		"```bash",
		runCommand,
		"```",
		"",
		"## Code Walkthrough",
		"Walkthrough.",
		"",
		"## Try It",
		"Try it.",
		"",
		"## In Production",
		"Production.",
		"",
		"## Thinking Questions",
		"Questions.",
		"",
		"## Next Step",
		"This fixture has no next item.",
		"",
	}, "\n")
}

func validCodeStandards() string {
	return strings.Join([]string{
		"# Code Quality & Style Standards",
		"",
		"## Standard Layers",
		"",
		"Machine-enforced and review-enforced rules.",
		"",
		"// Level: Foundation | Core | Production | Stretch",
		"",
		"Machine Role comments can satisfy this requirement for exported symbols.",
		"",
		"do not use legacy `Forward Reference` or `Backward Reference` labels",
		"",
		"## Curriculum Registry Standard",
		"",
		"curriculum.v2.json stays canonical.",
		"",
		"## Lesson Proof Surface",
		"",
		"One coherent proof surface.",
		"",
		"## Production-Shaped Code",
		"",
		"Production-shaped examples.",
		"",
		"```bash",
		"go test -coverprofile=coverage.out ./...",
		"```",
		"",
	}, "\n")
}

func flagshipReadme(item V2Item) string {
	if len(item.NextItems) == 0 {
		return "# Flagship\n\n## Next Step\n\nThis path is complete. Return to the section README or continue with the next project milestone.\n"
	}

	nextID := item.NextItems[0]
	nextPath := fixtureItemPath("s11", nextID)
	linkTarget, err := filepath.Rel(filepath.FromSlash(item.Path), filepath.Join(filepath.FromSlash(nextPath), "README.md"))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("# Flagship\n\n## Next Step\n\nNext: `%s` -> [`%s`](%s)\n", nextID, nextPath, filepath.ToSlash(linkTarget))
}

func mainGoForItem(item V2Item) string {
	level, ok := levelDisplayLabels[item.Level]
	if !ok {
		level = item.Level
	}

	return mainGoForItemWithLevel(item, level)
}

func mainGoForItemWithLevel(item V2Item, level string) string {
	label := ""
	switch item.SectionID {
	case "s09":
		label = "// Section 09 fixture\n"
	case "s10":
		label = "// Section 10 fixture\n"
	}

	return strings.Join([]string{
		"// Copyright (c) 2026 Rasel Hossen",
		"// Licensed under The Go Engineer License v1.0",
		"",
		"// ============================================================================",
		"// Section " + canonicalV2Sections[sectionIndex(item.SectionID)].Number + ": Fixture - " + item.Title,
		"// Level: " + level,
		"// ============================================================================",
		"//",
		"// WHAT YOU'LL LEARN:",
		"//   - How the fixture proves validator behavior.",
		"//",
		"// WHY THIS MATTERS:",
		"//   Validator fixtures need the same source contract as learner lessons.",
		"// RUN: " + item.RunCommand,
		"// ============================================================================",
		"",
		"package main",
		"",
		label + "func main() {",
		"\t// KEY TAKEAWAY:",
		"\t// - Fixture code follows the public lesson source contract.",
		"}",
		"",
	}, "\n")
}

func sectionReadme(sectionID string) string {
	labels := canonicalSectionReadmeTracks[sectionID]
	if len(labels) == 0 {
		return "# Section\n"
	}

	return "# Section\n\n" + strings.Join(labels, "\n") + "\n"
}

func sectionPrerequisites(index int) []string {
	if index == 0 {
		return []string{}
	}
	return []string{canonicalV2Sections[index-1].ID}
}

func sectionIDForItemID(id string) string {
	if strings.HasPrefix(id, "OPSL.") {
		return "s11"
	}

	for _, section := range canonicalV2Sections {
		for _, candidate := range append(append([]string{}, section.EntryPoints...), section.Outputs...) {
			if candidate == id {
				return section.ID
			}
		}
	}

	panic("unknown fixture item id: " + id)
}

func fixtureItemPath(sectionID, id string) string {
	if sectionID == "s11" {
		return fmt.Sprintf("11-flagship/01-opslane/modules/%02d-%s", itemNumber(id), strings.ToLower(strings.ReplaceAll(id, ".", "-")))
	}

	section := canonicalV2Sections[sectionIndex(sectionID)]
	return filepath.ToSlash(filepath.Join(section.PathPrefix, fmt.Sprintf("%03d-%s", fixtureOrderNumber(id), strings.ToLower(strings.ReplaceAll(id, ".", "-")))))
}

func fixtureOrderNumber(id string) int {
	for idx, candidate := range orderedFixtureItemIDs() {
		if candidate == id {
			return idx + 1
		}
	}
	panic("unknown fixture item id: " + id)
}

func sectionIndex(sectionID string) int {
	for idx, section := range canonicalV2Sections {
		if section.ID == sectionID {
			return idx
		}
	}
	panic("unknown fixture section id: " + sectionID)
}

func itemNumber(id string) int {
	_, number, ok := splitCurriculumID(id)
	if !ok {
		panic("invalid fixture item id: " + id)
	}
	return number
}

func writeCurriculum(t *testing.T, root string, cur V2Curriculum) {
	t.Helper()

	data, err := canonicalV2CurriculumJSON(cur)
	if err != nil {
		t.Fatalf("marshal curriculum fixture: %v", err)
	}
	writeFile(t, root, "curriculum.v2.json", string(data))
}

func readCurriculum(t *testing.T, root string) V2Curriculum {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(root, "curriculum.v2.json"))
	if err != nil {
		t.Fatalf("read curriculum fixture: %v", err)
	}

	var cur V2Curriculum
	if err := json.Unmarshal(data, &cur); err != nil {
		t.Fatalf("parse curriculum fixture: %v", err)
	}
	return cur
}

func writeFile(t *testing.T, root, relativePath, contents string) {
	t.Helper()

	fullPath := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(fullPath), err)
	}
	if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func mustMkdir(t *testing.T, root, relativePath string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(root, relativePath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", relativePath, err)
	}
}

func runValidate(t *testing.T, root string) (Result, []string) {
	t.Helper()

	var reports []string
	result, err := Validate(root, func(message string) {
		reports = append(reports, message)
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	return result, reports
}

func requireReportContains(t *testing.T, reports []string, want string) {
	t.Helper()

	for _, report := range reports {
		if strings.Contains(report, want) {
			return
		}
	}
	t.Fatalf("expected report containing %q not found in %v", want, reports)
}

func countReportsContaining(reports []string, want string) int {
	count := 0
	for _, report := range reports {
		if strings.Contains(report, want) {
			count++
		}
	}
	return count
}

func requireItem(t *testing.T, cur V2Curriculum, id string) V2Item {
	t.Helper()

	for _, item := range cur.Items {
		if item.ID == id {
			return item
		}
	}
	t.Fatalf("fixture item %s not found", id)
	return V2Item{}
}
