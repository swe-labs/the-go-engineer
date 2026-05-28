package main

import (
	"fmt"
	"strings"
)

func ValidateSchema(cur *Curriculum, report *Report) {
	check := "schema"
	for _, name := range metadataFiles {
		f := cur.Files[name]
		// schema.v3.json is a JSON Schema document; it uses $schema/$id/title rather than curriculum metadata keys.
		if name == "schema.v3.json" {
			if str(f, "$schema") == "" && str(f, "title") == "" {
				report.Error(check, name, "schema document must declare $schema or title")
			}
			continue
		}
		if str(f, "schema_version") == "" {
			report.Error(check, name, "missing schema_version")
		}
		if str(f, "document_type") == "" {
			report.Error(check, name, "missing document_type")
		}
		if str(f, "curriculum_version") == "" {
			report.Error(check, name, "missing curriculum_version")
		}
	}
	ValidateRepositoryArchitecture(cur, report)
	validateDuplicateIDs(cur, report)
	validateModuleRequiredFields(cur, report)
	validateItemRequiredFields(cur, report)
}

func ValidateRepositoryArchitecture(cur *Curriculum, report *Report) {
	check := "repository-architecture"
	core := cur.Files["path.core.json"]
	elect := cur.Files["path.electives.json"]
	coreStruct := stringsList(core, "repository_structure")
	electStruct := stringsList(elect, "repository_structure")
	requireEntry := func(entries []string, want string, source string) {
		if !stringIn(entries, want) {
			report.Error(check, source, "repository_structure missing %q", want)
		}
	}
	requireEntry(coreStruct, "metadata/path.core.json", "path.core.json")
	requireEntry(coreStruct, "modules/00-orientation", "path.core.json")
	requireEntry(coreStruct, "electives/advanced-electives", "path.core.json")
	requireEntry(coreStruct, "tools/", "path.core.json")
	requireEntry(coreStruct, "assets/", "path.core.json")
	requireEntry(electStruct, "electives/advanced-electives", "path.electives.json")
	for _, m := range cur.CoreModules {
		if str(m, "slug") == "advanced-electives" {
			report.Error(check, str(m, "id"), "advanced-electives must not be in core modules")
		}
	}
	foundElective := false
	for _, m := range cur.ElectiveModules {
		if str(m, "slug") == "advanced-electives" {
			foundElective = true
		}
	}
	if !foundElective {
		report.Error(check, "path.electives.json", "advanced-electives module missing from electives")
	}
}

func validateDuplicateIDs(cur *Curriculum, report *Report) {
	check := "unique-ids"
	seen := map[string]string{}
	scan := func(kind string, items []map[string]any, key string) {
		for _, x := range items {
			id := str(x, key)
			if id == "" {
				report.Error(check, kind, "missing %s", key)
				continue
			}
			if prev, ok := seen[id]; ok {
				report.Error(check, id, "duplicate id in %s and %s", prev, kind)
				continue
			}
			seen[id] = kind
		}
	}
	scan("items", cur.Items, "id")
	scan("modules", cur.Modules, "id")
	scan("projects", cur.Projects, "id")
	scan("assessments", cur.Assessments, "id")
	scan("concepts", cur.Concepts, "concept")
}

func validateModuleRequiredFields(cur *Curriculum, report *Report) {
	check := "module-schema"
	required := []string{"id", "number", "slug", "title", "phase", "path", "status", "learning_goal", "summary", "order", "prerequisites", "entry_item_ids", "terminal_item_ids", "tags", "readme_status", "readme_contract", "cognitive_load", "recommended_break_after", "contains_foundational_hard_concepts", "pacing"}
	allowedPhases := map[string]bool{"orientation": true, "foundations": true, "tooling": true, "go-core": true, "engineering-core": true, "cli-io": true, "data": true, "security": true, "reliability": true, "performance": true, "architecture": true, "delivery": true, "career": true, "flagship": true, "elective": true}
	allowedLoad := map[string]bool{"low": true, "moderate": true, "high": true}
	allowedPacing := map[string]bool{"gentle": true, "moderate": true, "steady": true}
	for _, m := range cur.Modules {
		id := str(m, "id")
		for _, field := range required {
			if _, ok := m[field]; !ok {
				report.Error(check, id, "missing %s", field)
			}
		}
		if !allowedPhases[str(m, "phase")] {
			report.Error(check, id, "invalid phase %q", str(m, "phase"))
		}
		if !allowedLoad[str(m, "cognitive_load")] {
			report.Error(check, id, "invalid cognitive_load %q", str(m, "cognitive_load"))
		}
		if !allowedPacing[str(m, "pacing")] {
			report.Error(check, id, "invalid pacing %q", str(m, "pacing"))
		}
		if str(m, "status") != "stable" {
			report.Error(check, id, "module status must be stable, got %q", str(m, "status"))
		}
		if str(m, "readme_status") != "golden" {
			report.Error(check, id, "module readme_status must be golden, got %q", str(m, "readme_status"))
		}
		path := str(m, "path")
		if str(m, "phase") == "elective" {
			if !strings.HasPrefix(path, "electives/") {
				report.Error(check, id, "elective module path must start with electives/, got %q", path)
			}
		} else if !strings.HasPrefix(path, "modules/") {
			report.Error(check, id, "core module path must start with modules/, got %q", path)
		}
		if strings.TrimSpace(str(m, "learning_goal")) == strings.TrimSpace(str(m, "summary")) && len(str(m, "summary")) < 120 {
			report.Info(check, id, "summary is identical to learning_goal; consider richer learner-facing summary")
		}
	}
}

func validateItemRequiredFields(cur *Curriculum, report *Report) {
	check := "item-schema"
	required := []string{"id", "module_id", "slug", "title", "type", "subtype", "status", "difficulty", "phase", "order", "estimated_minutes", "learning_objective", "required_prior_knowledge", "prerequisites", "next_item_ids", "zero_magic", "crossrefs", "proof", "content_contract", "verification", "files", "tags", "documentation_mode", "readme_status", "zero_magic_status", "readme_contract"}
	allowedTypes := map[string]bool{"lesson": true, "lab": true, "project": true, "checkpoint": true, "assessment": true, "reference": true, "capstone": true}
	for _, item := range cur.Items {
		id := str(item, "id")
		for _, field := range required {
			if _, ok := item[field]; !ok {
				report.Error(check, id, "missing %s", field)
			}
		}
		if !allowedTypes[str(item, "type")] {
			report.Error(check, id, "invalid type %q", str(item, "type"))
		}
		if str(item, "status") != "stable" {
			report.Error(check, id, "status must be stable, got %q", str(item, "status"))
		}
		if str(item, "readme_status") != "golden" {
			report.Error(check, id, "readme_status must be golden, got %q", str(item, "readme_status"))
		}
		if str(item, "zero_magic_status") != "golden" {
			report.Error(check, id, "zero_magic_status must be golden, got %q", str(item, "zero_magic_status"))
		}
		if _, ok := cur.ModuleByID[str(item, "module_id")]; !ok {
			report.Error(check, id, "unknown module_id %q", str(item, "module_id"))
		}
		if num(item, "estimated_minutes") <= 0 {
			report.Error(check, id, "estimated_minutes must be positive")
		}
		if len(str(item, "learning_objective")) < 50 {
			report.Warn(check, id, "learning_objective is short; target a measurable outcome")
		}
		files := obj(item, "files")
		expectedPrefix := "modules/"
		if isElectiveItem(id) {
			expectedPrefix = "electives/"
		}
		for _, key := range []string{"readme_path", "main_path", "starter_path", "solution_path", "test_path", "assets_dir"} {
			v := str(files, key)
			if v != "" && !strings.HasPrefix(v, expectedPrefix) {
				report.Error(check, id, "files.%s must start with %s, got %q", key, expectedPrefix, v)
			}
		}
		if !strings.HasSuffix(str(files, "readme_path"), "README.md") {
			report.Error(check, id, "files.readme_path must end with README.md")
		}
		if item["content_contract"] == nil {
			report.Error(check, id, "missing content_contract")
		}
		if item["verification"] == nil {
			report.Error(check, id, "missing verification")
		}
		if item["proof"] == nil {
			report.Error(check, id, "missing proof")
		}
	}
	report.Counts["items"] = len(cur.Items)
	report.Counts["modules"] = len(cur.Modules)
	report.Counts["module_15_items"] = countItemsInModule(cur, "module-15")
}

func countItemsInModule(cur *Curriculum, moduleID string) int {
	n := 0
	for _, item := range cur.Items {
		if str(item, "module_id") == moduleID {
			n++
		}
	}
	return n
}

func formatList(xs []string) string { return fmt.Sprintf("%v", xs) }
