package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ValidateMetadata(cfg Config) ValidationResult {
	var result ValidationResult
	result.Merge(validateJSONFiles(cfg))
	metadata, loadResult := loadMetadata(cfg)
	result.Merge(loadResult)
	if !loadResult.OK() {
		return result
	}
	result.Merge(validateSchemaBasics(metadata))
	result.Merge(validateGraph(metadata))
	result.Merge(validateItems(metadata))
	result.Merge(validateProjects(metadata))
	result.Merge(validateAssessments(metadata))
	result.Merge(validateConcepts(metadata))
	result.Merge(validateCrossrefs(metadata))
	result.Merge(validateFailures(metadata))
	result.Merge(validateReadmeContracts(metadata))
	result.Merge(validateMigration(metadata))
	return result
}

func validateJSONFiles(cfg Config) ValidationResult {
	var r ValidationResult
	required := []string{
		"workspace.json", "schema.v3.json", "path.core.json", "path.electives.json", "concepts.json",
		"projects.json", "assessments.json", "crossrefs.json", "failures.json", "readme.contracts.json",
		"migration.v2-to-v3.json", "VALIDATION.metadata.json", "legacy/curriculum.v2.json",
		"legacy/curriculum.v2.lock.json", "legacy/unmapped-v2-report.json", "legacy/migration-notes.md",
	}
	for _, rel := range required {
		path := filepath.Join(cfg.MetadataDir, rel)
		if !fileExists(path) {
			r.Errorf("missing metadata file %s", rel)
			continue
		}
		if strings.HasSuffix(rel, ".json") {
			data, err := os.ReadFile(path)
			if err != nil {
				r.Errorf("cannot read %s: %v", rel, err)
				continue
			}
			if !json.Valid(data) {
				r.Errorf("invalid JSON in %s", rel)
			}
		}
	}
	return r
}

func validateSchemaBasics(m Metadata) ValidationResult {
	var r ValidationResult
	for name, value := range map[string]string{
		"path.core.document_type":      m.Core.DocumentType,
		"path.electives.document_type": m.Electives.DocumentType,
		"projects.document_type":       m.Projects.DocumentType,
		"assessments.document_type":    m.Assessments.DocumentType,
		"concepts.document_type":       m.Concepts.DocumentType,
		"crossrefs.document_type":      m.Crossrefs.DocumentType,
		"failures.document_type":       m.Failures.DocumentType,
	} {
		if strings.TrimSpace(value) == "" {
			r.Errorf("%s is empty", name)
		}
	}
	if len(m.Core.Modules) == 0 || len(m.Core.Items) == 0 {
		r.Errorf("core path must contain modules and items")
	}
	if len(m.Electives.Modules) == 0 || len(m.Electives.Items) == 0 {
		r.Errorf("electives path must contain modules and items")
	}
	return r
}

func validateItems(m Metadata) ValidationResult {
	var r ValidationResult
	seen := map[string]bool{}
	moduleID := moduleIDs(m)
	assessmentID := assessmentIDs(m)
	for _, item := range allItems(m) {
		if item.ID == "" {
			r.Errorf("item has empty id")
		}
		if seen[item.ID] {
			r.Errorf("duplicate item id %s", item.ID)
		}
		seen[item.ID] = true
		if item.ModuleID == "" || moduleID[item.ModuleID].ID == "" {
			r.Errorf("%s references unknown module_id %s", item.ID, item.ModuleID)
		}
		if strings.TrimSpace(item.Title) == "" {
			r.Errorf("%s title is empty", item.ID)
		}
		if strings.TrimSpace(item.Type) == "" {
			r.Errorf("%s type is empty", item.ID)
		}
		if !isStableStatus(item.Status) {
			r.Errorf("%s status must be stable/ready/published, got %q", item.ID, item.Status)
		}
		if item.ZeroMagicStatus != "" && item.ZeroMagicStatus != "golden" {
			r.Errorf("%s zero_magic_status must be golden", item.ID)
		}
		if item.ReadmeStatus != "" && item.ReadmeStatus != "golden" {
			r.Errorf("%s readme_status must be golden", item.ID)
		}
		validateZeroMagic(item, &r)
		validateFiles(item.ID, item.Files, &r)
		if item.Proof != nil && item.Proof.AssessmentID != "" {
			if _, ok := assessmentID[item.Proof.AssessmentID]; !ok {
				r.Errorf("%s proof.assessment_id %s does not exist", item.ID, item.Proof.AssessmentID)
			}
		}
	}
	return r
}

func validateZeroMagic(item Item, r *ValidationResult) {
	if item.ZeroMagic == nil {
		r.Errorf("%s missing zero_magic", item.ID)
		return
	}
	checks := map[string]string{
		"problem_solved":         item.ZeroMagic.ProblemSolved,
		"why_it_exists":          item.ZeroMagic.WhyItExists,
		"mental_model":           item.ZeroMagic.MentalModel,
		"under_the_hood":         item.ZeroMagic.UnderTheHood,
		"how_go_uses_it":         item.ZeroMagic.HowGoUsesIt,
		"real_world_usage":       item.ZeroMagic.RealWorldUsage,
		"proof_of_understanding": item.ZeroMagic.ProofOfUnderstanding,
	}
	for name, value := range checks {
		if strings.TrimSpace(value) == "" {
			r.Errorf("%s zero_magic.%s is empty", item.ID, name)
		}
	}
	if !nonEmptyAny(item.ZeroMagic.BeginnerMistakes) {
		r.Errorf("%s zero_magic.beginner_mistakes is empty", item.ID)
	}
}

func validateFiles(owner string, files FileRefs, r *ValidationResult) {
	paths := map[string]string{
		"readme_path":   files.ReadmePath,
		"main_path":     files.MainPath,
		"test_path":     files.TestPath,
		"starter_path":  files.StarterPath,
		"solution_path": files.SolutionPath,
		"assets_dir":    files.AssetsDir,
	}
	for name, path := range paths {
		if path == "" {
			continue
		}
		if !canonicalContentPath(path) {
			r.Errorf("%s files.%s must use typed curriculum path, got %s", owner, name, path)
		}
	}
	for _, path := range files.DiagramPaths {
		if path != "" && !canonicalContentPath(path) {
			r.Errorf("%s diagram path must use typed curriculum path, got %s", owner, path)
		}
	}
}

func validateFailures(m Metadata) ValidationResult {
	var r ValidationResult
	modules := moduleIDs(m)
	for moduleID, categories := range m.Failures.RequiredCoverage {
		if _, ok := modules[moduleID]; !ok {
			r.Errorf("failures.required_coverage references unknown module %s", moduleID)
		}
		if len(categories) == 0 {
			r.Errorf("failures.required_coverage[%s] is empty", moduleID)
		}
	}
	known := map[string]bool{}
	for _, cat := range m.Failures.FailureCategories {
		id := cat.ID
		if id == "" {
			id = cat.Category
		}
		if id == "" {
			id = cat.Name
		}
		if id == "" {
			r.Errorf("failure category has empty id/category/name")
			continue
		}
		if known[id] {
			r.Errorf("duplicate failure category %s", id)
		}
		known[id] = true
		for _, mid := range append(cat.ModuleIDs, cat.Modules...) {
			if _, ok := modules[mid]; !ok {
				r.Errorf("failure category %s references unknown module %s", id, mid)
			}
		}
	}
	return r
}

func validateMigration(m Metadata) ValidationResult {
	var r ValidationResult
	migration, _ := m.RawMigration["migration"].(map[string]any)
	if migration == nil {
		r.Errorf("migration.v2-to-v3.json missing migration object")
		return r
	}
	if path, _ := migration["frozen_source_path"].(string); path == "" {
		r.Errorf("migration missing frozen_source_path")
	}
	summary, _ := migration["legacy_item_coverage_summary"].(map[string]any)
	if summary == nil {
		r.Errorf("migration missing legacy_item_coverage_summary")
		return r
	}
	unmapped := fmt.Sprintf("%v", summary["unmapped_count"])
	if unmapped != "0" && unmapped != "0.0" {
		r.Errorf("migration has unmapped legacy items: %s", unmapped)
	}
	return r
}

func validateConcepts(m Metadata) ValidationResult {
	var r ValidationResult
	items := itemIDs(m)
	modules := moduleIDs(m)
	projects := projectIDs(m)
	assessments := assessmentIDs(m)
	seen := map[string]bool{}
	known := func(id string) bool {
		if _, ok := items[id]; ok {
			return true
		}
		if _, ok := modules[id]; ok {
			return true
		}
		if _, ok := projects[id]; ok {
			return true
		}
		if _, ok := assessments[id]; ok {
			return true
		}
		return false
	}
	for _, concept := range m.Concepts.Concepts {
		if strings.TrimSpace(concept.Concept) == "" {
			r.Errorf("concept has empty name")
		}
		if seen[concept.Concept] {
			r.Errorf("duplicate concept %s", concept.Concept)
		}
		seen[concept.Concept] = true
		if !known(concept.CanonicalOwner) {
			r.Errorf("concept %s unknown canonical_owner %s", concept.Concept, concept.CanonicalOwner)
		}
		for _, loc := range concept.PreviewLocations {
			if !known(loc) {
				r.Errorf("concept %s unknown preview location %s", concept.Concept, loc)
			}
		}
		for _, loc := range concept.ReinforcementLocations {
			if !known(loc) {
				r.Errorf("concept %s unknown reinforcement location %s", concept.Concept, loc)
			}
		}
	}
	return r
}

func nonEmptyAny(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case string:
		return strings.TrimSpace(x) != ""
	case []any:
		return len(x) > 0
	case []string:
		return len(x) > 0
	default:
		return true
	}
}
