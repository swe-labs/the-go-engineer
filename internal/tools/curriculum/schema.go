package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func validateSchema() {
	fmt.Println("=== validate-schema ===")
	ok := true

	// 1. Parse all 9 canonical active JSON manifests and verify metadata consistency
	files := []string{
		"workspace.json",
		"path.core.json",
		"path.electives.json",
		"projects.json",
		"assessments.json",
		"crossrefs.json",
		"concepts.json",
		"failures.json",
		"readme.contracts.json",
	}
	for _, f := range files {
		fullPath := fmt.Sprintf("%s/%s", curriculumDir, f)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			warnf("cannot read %s: %v", f, err)
			ok = false
			continue
		}
		if !json.Valid(data) {
			warnf("%s: invalid JSON", f)
			ok = false
			continue
		}

		var header struct {
			SchemaVersion     string `json:"schema_version"`
			CurriculumVersion string `json:"curriculum_version"`
			Status            string `json:"status"`
			LastUpdated       string `json:"last_updated"`
		}
		if err := json.Unmarshal(data, &header); err == nil {
			if header.SchemaVersion != "3.1.0" {
				warnf("%s: expected schema_version '3.1.0', got '%s'", f, header.SchemaVersion)
				ok = false
			}
			if header.CurriculumVersion != "3.0.0" {
				warnf("%s: expected curriculum_version '3.0.0', got '%s'", f, header.CurriculumVersion)
				ok = false
			}
			if header.Status != "stable" {
				warnf("%s: expected status 'stable', got '%s'", f, header.Status)
				ok = false
			}
			if header.LastUpdated != "2026-10-23" {
				warnf("%s: expected last_updated '2026-10-23', got '%s'", f, header.LastUpdated)
				ok = false
			}
		}
		fmt.Printf("  %s: valid JSON & metadata\n", f)
	}

	// Validate schema.v3.json (the JSON Schema definition)
	schemaPath := filepath.Join(curriculumDir, "schema.v3.json")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		warnf("cannot read schema.v3.json: %v", err)
		ok = false
	} else if !json.Valid(schemaData) {
		warnf("schema.v3.json: invalid JSON")
		ok = false
	} else {
		var schemaMeta struct {
			ID     string `json:"$id"`
			Schema string `json:"$schema"`
			Title  string `json:"title"`
		}
		if err := json.Unmarshal(schemaData, &schemaMeta); err != nil || schemaMeta.ID == "" || schemaMeta.Schema == "" {
			warnf("schema.v3.json: missing $id or $schema")
			ok = false
		} else {
			fmt.Printf("  schema.v3.json: valid JSON Schema definition\n")
		}
	}

	// Verify readme.contracts.json defines core contracts and anti-overengineering rules
	var contractsBundle struct {
		Contracts []struct {
			ID                   string         `json:"id"`
			AppliesTo            string         `json:"applies_to"`
			CoreRequiredHeadings []string       `json:"core_required_headings"`
			ContextualHeadings   []string       `json:"contextual_headings"`
			SubtypeProfiles      map[string]any `json:"subtype_profiles"`
		} `json:"contracts"`
	}
	readJSON("readme.contracts.json", &contractsBundle)
	contractMap := make(map[string]struct {
		CoreRequiredHeadings []string
		ContextualHeadings   []string
		SubtypeProfiles      map[string]any
	})
	for _, c := range contractsBundle.Contracts {
		contractMap[c.ID] = struct {
			CoreRequiredHeadings []string
			ContextualHeadings   []string
			SubtypeProfiles      map[string]any
		}{c.CoreRequiredHeadings, c.ContextualHeadings, c.SubtypeProfiles}
	}

	expectedContracts := []string{"root.v3", "module.v3", "lesson.v3", "elective-lesson.v3", "project.v3", "companion-stage.v3"}
	for _, cName := range expectedContracts {
		c, exists := contractMap[cName]
		if !exists {
			warnf("readme.contracts.json: missing contract '%s'", cName)
			ok = false
		} else if cName == "lesson.v3" {
			if len(c.CoreRequiredHeadings) == 0 || len(c.ContextualHeadings) == 0 {
				warnf("readme.contracts.json: lesson.v3 missing core_required_headings or contextual_headings")
				ok = false
			}
			if len(c.SubtypeProfiles) == 0 {
				warnf("readme.contracts.json: lesson.v3 missing subtype_profiles")
				ok = false
			}
		}
	}

	// 2. Read core and electives for deeper checks
	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)
	var projects ProjectsBundle
	readJSON("projects.json", &projects)
	var crossrefs CrossrefBundle
	readJSON("crossrefs.json", &crossrefs)

	// Verify inventory counts and module statuses
	if len(core.Items) != 421 {
		warnf("path.core.json: expected exactly 421 core items, got %d", len(core.Items))
		ok = false
	}
	if len(electives.Items) != 20 {
		warnf("path.electives.json: expected exactly 20 elective items, got %d", len(electives.Items))
		ok = false
	}
	if len(core.Modules) != 18 {
		warnf("path.core.json: expected exactly 18 core modules, got %d", len(core.Modules))
		ok = false
	}
	for _, m := range core.Modules {
		if m.Status != "stable" {
			warnf("path.core.json module %s: expected status 'stable', got '%s'", m.ID, m.Status)
			ok = false
		}
	}

	// 3. Check required fields on items
	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	type fieldCheck struct {
		name      string
		getter    func(Item) any
		skipTypes map[string]bool
	}
	checks := []fieldCheck{
		{"id", func(i Item) any { return i.ID }, nil},
		{"module_id", func(i Item) any { return i.ModuleID }, nil},
		{"slug", func(i Item) any { return i.Slug }, nil},
		{"title", func(i Item) any { return i.Title }, nil},
		{"type", func(i Item) any { return i.Type }, nil},
		{"subtype", func(i Item) any { return i.Subtype }, nil},
		{"status", func(i Item) any { return i.Status }, nil},
		{"difficulty", func(i Item) any { return i.Difficulty }, nil},
		{"phase", func(i Item) any { return i.Phase }, nil},
		{"learning_objective", func(i Item) any { return i.LearningObjective }, nil},
		{"estimated_minutes", func(i Item) any { return i.EstimatedMinutes }, map[string]bool{"project": true}},
		{"prerequisites", func(i Item) any { return i.Prerequisites }, nil},
		{"zero_magic", func(i Item) any { return i.ZeroMagic }, map[string]bool{"project": true, "assessment": true, "checkpoint": true, "capstone": true}},
		{"crossrefs", func(i Item) any { return i.CrossRefs }, nil},
		{"proof", func(i Item) any { return i.Proof }, map[string]bool{"project": true, "assessment": true, "checkpoint": true, "capstone": true}},
		{"files", func(i Item) any { return i.Files }, map[string]bool{"checkpoint": true}},
	}
	for _, item := range allItems {
		for _, c := range checks {
			if c.skipTypes[item.Type] {
				continue
			}
			v := c.getter(item)
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				warnf("%s: missing '%s'", item.ID, c.name)
				ok = false
			} else if rv.Kind() == reflect.String && rv.String() == "" {
				warnf("%s: '%s' is empty", item.ID, c.name)
				ok = false
			} else if rv.Kind() == reflect.Slice && rv.IsNil() {
				warnf("%s: '%s' is null (empty array expected)", item.ID, c.name)
				ok = false
			}
		}
	}

	// 4. Check zero_magic fields
	for _, item := range allItems {
		if item.ZeroMagic != nil {
			zm := item.ZeroMagic
			zmChecks := []struct {
				name  string
				value string
				arr   []any
			}{
				{"problem_solved", zm.ProblemSolved, nil},
				{"why_it_exists", zm.WhyItExists, nil},
				{"mental_model", zm.MentalModel, nil},
				{"under_the_hood", zm.UnderTheHood, nil},
				{"how_go_uses_it", zm.HowGoUsesIt, nil},
				{"real_world_usage", zm.RealWorldUsage, nil},
				{"proof_of_understanding", zm.ProofOfUnderstanding, nil},
				{"beginner_mistakes", "", zm.BeginnerMistakes},
			}
			for _, c := range zmChecks {
				if c.arr != nil {
					if len(c.arr) == 0 {
						warnf("%s: zero_magic.%s is empty", item.ID, c.name)
						ok = false
					}
				} else if c.value == "" {
					warnf("%s: zero_magic.%s is empty", item.ID, c.name)
					ok = false
				}
			}
		}
	}

	// 5. Check proof fields
	for _, item := range allItems {
		if item.Proof != nil {
			if item.Proof.PracticeTask == "" {
				warnf("%s: proof.practice_task is empty", item.ID)
				ok = false
			}
			if item.Proof.AssessmentID == "" {
				warnf("%s: proof.assessment_id is empty", item.ID)
				ok = false
			}
		}
	}

	// 6. Check canonical key ordering in core JSON using byte-level tokenization (alphabetical order)
	// Read raw file bytes for deterministic token-level key extraction
	rawBytes, err := os.ReadFile(filepath.Join(curriculumDir, "path.core.json"))
	if err == nil {
		// Find items array, extract keys for each item object in file order
		itemKeys := extractItemKeysInOrder(rawBytes)
		for _, ik := range itemKeys {
			for i := 1; i < len(ik.keys); i++ {
				if ik.keys[i] < ik.keys[i-1] {
					warnf("%s: key '%s' out of canonical alphabetical order (after '%s')", ik.id, ik.keys[i], ik.keys[i-1])
					ok = false
					break
				}
			}
		}
	} else {
		warnf("cannot read path.core.json for key ordering check: %v", err)
		ok = false
	}

	// 7. Check no duplicate item IDs across core and electives
	seen := make(map[string]string)
	for _, item := range allItems {
		if prev, dup := seen[item.ID]; dup {
			warnf("duplicate item ID '%s' (in %s and %s)", item.ID, prev, item.SourceLegacyIDs)
			ok = false
		}
		if len(item.SourceLegacyIDs) > 0 {
			seen[item.ID] = item.SourceLegacyIDs[0]
		} else {
			seen[item.ID] = item.ModuleID
		}
	}

	// 8. Check module IDs are referenced by items consistently
	moduleIDs := make(map[string]*Module)
	for i := range core.Modules {
		m := &core.Modules[i]
		moduleIDs[m.ID] = m
	}
	for i := range electives.Modules {
		m := &electives.Modules[i]
		moduleIDs[m.ID] = m
	}
	for _, item := range allItems {
		if _, ok := moduleIDs[item.ModuleID]; !ok {
			warnf("%s: references unknown module_id '%s'", item.ID, item.ModuleID)
			ok = false
		}
	}

	// 9. Check assessment ids match valid assessments
	type Assessment struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	var assessments AssessmentsBundle
	readJSON("assessments.json", &assessments)
	assessIDs := make(map[string]bool)
	for _, a := range assessments.Assessments {
		assessIDs[a.ID] = true
	}
	if len(assessIDs) == 0 {
		// Try parsing raw
		var rawAssessments map[string]interface{}
		readJSON("assessments.json", &rawAssessments)
		fmt.Printf("  assessments keys: %v\n", keysOfMap(rawAssessments))
	}

	for _, item := range allItems {
		if item.Proof != nil && item.Proof.AssessmentID != "" {
			if !assessIDs[item.Proof.AssessmentID] {
				warnf("%s: proof.assessment_id '%s' not found in assessments.json", item.ID, item.Proof.AssessmentID)
				ok = false
			}
		}
	}

	if ok {
		fmt.Println("  All schema checks passed.")
	} else {
		fmt.Println("  Schema checks completed with warnings/errors.")
	}
}

func indexOf(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}

type itemKeyEntry struct {
	id   string
	keys []string
}

// extractItemKeysInOrder reads raw JSON bytes and extracts the keys of each item
// object in the items array, preserving file byte order (deterministic).
func extractItemKeysInOrder(data []byte) []itemKeyEntry {
	var result []itemKeyEntry
	dec := json.NewDecoder(strings.NewReader(string(data)))

	// Must start with top-level object '{'
	t, err := dec.Token()
	if err != nil {
		return nil
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return nil
	}

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		key, ok := tok.(string)
		if !ok {
			continue
		}
		if key != "items" {
			skipValue(dec)
			continue
		}

		// 'items' value must start with '['
		arrTok, err := dec.Token()
		if err != nil {
			break
		}
		if delim, ok := arrTok.(json.Delim); !ok || delim != '[' {
			break
		}

		// Iterate each item in items array
		for dec.More() {
			itemTok, err := dec.Token()
			if err != nil {
				break
			}
			if delim, ok := itemTok.(json.Delim); !ok || delim != '{' {
				skipValue(dec)
				continue
			}

			entry := itemKeyEntry{}
			for dec.More() {
				kTok, err := dec.Token()
				if err != nil {
					break
				}
				k, ok := kTok.(string)
				if !ok {
					continue
				}
				entry.keys = append(entry.keys, k)
				if k == "id" {
					idTok, err := dec.Token()
					if err == nil {
						entry.id = fmt.Sprintf("%v", idTok)
					}
				} else {
					skipValue(dec)
				}
			}
			// Consume closing '}'
			dec.Token()
			result = append(result, entry)
		}
		// Consume closing ']'
		dec.Token()
	}

	return result
}

// skipValue skips one JSON value (object, array, string, number, bool, null).
func skipValue(dec *json.Decoder) {
	t, err := dec.Token()
	if err != nil {
		return
	}
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			for dec.More() {
				dec.Token() // key
				skipValue(dec)
			}
			dec.Token() // closing }
		case '[':
			for dec.More() {
				skipValue(dec)
			}
			dec.Token() // closing ]
		}
	}
}

func keysOfMap(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
