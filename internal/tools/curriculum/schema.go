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

	// 1. Parse all JSON files
	files := []string{"path.core.json", "path.electives.json", "projects.json", "assessments.json", "crossrefs.json"}
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
		fmt.Printf("  %s: valid JSON\n", f)
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
		{"zero_magic", func(i Item) any { return i.ZeroMagic }, nil},
		{"crossrefs", func(i Item) any { return i.CrossRefs }, nil},
		{"proof", func(i Item) any { return i.Proof }, nil},
		{"files", func(i Item) any { return i.Files }, nil},
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

	// 6. Check canonical key ordering in core JSON using byte-level tokenization
	canonicalKeyOrder := []string{
		"id", "module_id", "slug", "title", "type", "subtype",
		"status", "difficulty", "phase", "order", "estimated_minutes",
		"learning_objective", "required_prior_knowledge", "prerequisites",
		"next_item_ids", "zero_magic", "crossrefs", "proof",
		"content_contract", "verification", "files", "source_legacy_ids", "tags",
	}
	// Read raw file bytes for deterministic token-level key extraction
	rawBytes, err := os.ReadFile(filepath.Join(curriculumDir, "path.core.json"))
	if err == nil {
		// Find items array, extract keys for each item object in file order
		itemKeys := extractItemKeysInOrder(rawBytes)
		for _, ik := range itemKeys {
			var lastIdx int
			for _, k := range ik.keys {
				idx := indexOf(canonicalKeyOrder, k)
				if idx < lastIdx {
					warnf("%s: key '%s' out of canonical order (after position %d)", ik.id, k, lastIdx)
					ok = false
					break
				}
				lastIdx = idx
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
	// Walk through the top-level object looking for "items" key
	for {
		t, err := dec.Token()
		if err != nil {
			break
		}
		if delim, ok := t.(json.Delim); ok {
			if delim == '{' {
				// Read keys of this object
				for dec.More() {
					keyTok, err := dec.Token()
					if err != nil {
						break
					}
					key := fmt.Sprintf("%v", keyTok)
					if key == "items" {
						// Consume the ":" token
						dec.Token()
						// items should be an array
						arrTok, err := dec.Token()
						if err != nil {
							break
						}
						if arrDelim, ok := arrTok.(json.Delim); ok && arrDelim == '[' {
							for dec.More() {
								// Each item is an object
								objTok, err := dec.Token()
								if err != nil {
									break
								}
								if objDelim, ok := objTok.(json.Delim); ok && objDelim == '{' {
									entry := itemKeyEntry{}
									var objKeys []string
									for dec.More() {
										kTok, err := dec.Token()
										if err != nil {
											break
										}
										k := fmt.Sprintf("%v", kTok)
										objKeys = append(objKeys, k)
										if k == "id" {
											// Read the id value
											dec.Token() // consume ":"
											idTok, _ := dec.Token()
											entry.id = fmt.Sprintf("%v", idTok)
										} else {
											// Skip the value
											skipValue(dec)
										}
									}
									// Consume the closing }
									dec.Token()
									entry.keys = objKeys
									result = append(result, entry)
								}
							}
						}
					} else {
						skipValue(dec)
					}
				}
			}
		}
		if delim, ok := t.(json.Delim); ok && delim == '}' {
			break
		}
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
		if delim == '{' {
			for dec.More() {
				dec.Token() // key
				skipValue(dec)
			}
			dec.Token() // closing }
		} else if delim == '[' {
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
