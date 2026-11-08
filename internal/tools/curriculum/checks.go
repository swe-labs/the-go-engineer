package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validatePrerequisiteClosure() {
	fmt.Println("=== validate-prerequisite-closure ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	itemMap := make(map[string]*Item)
	for i := range allItems {
		itemMap[allItems[i].ID] = &allItems[i]
	}

	// 1. Check no future-concept violations
	// An item should not depend on concepts from a later module
	moduleOrder := make(map[string]int)
	for _, m := range core.Modules {
		moduleOrder[m.ID] = m.Number
	}
	for _, m := range electives.Modules {
		moduleOrder[m.ID] = m.Number
	}

	for _, item := range allItems {
		itemModNum := moduleOrder[item.ModuleID]
		for _, prereqID := range item.Prerequisites {
			if prereqItem, exists := itemMap[prereqID]; exists {
				prereqModNum := moduleOrder[prereqItem.ModuleID]
				if prereqModNum > itemModNum || (prereqModNum == itemModNum && prereqItem.Order >= item.Order) {
					errorf("%s: prerequisite '%s' is not earlier (module %s order %d vs %s order %d)",
						item.ID, prereqID, prereqItem.ModuleID, prereqItem.Order, item.ModuleID, item.Order)
					ok = false
				}
			}
		}
	}

	// 2. Check no hidden prerequisite leaks (prereqs that don't exist as items)
	for _, item := range allItems {
		for _, prereqID := range item.Prerequisites {
			if _, exists := itemMap[prereqID]; !exists {
				errorf("%s: prerequisite '%s' does not exist in item graph", item.ID, prereqID)
				ok = false
			}
		}
	}

	// 3. Check entry items have no prerequisites from more than 1 module back (informational)
	for _, mod := range core.Modules {
		for _, entryID := range mod.EntryItemIDs {
			if entryItem, exists := itemMap[entryID]; exists {
				if len(entryItem.Prerequisites) > 0 {
					for _, prereqID := range entryItem.Prerequisites {
						if prereqItem, exists := itemMap[prereqID]; exists {
							prereqModNum := moduleOrder[prereqItem.ModuleID]
							modNum := moduleOrder[mod.ID]
							if prereqModNum < modNum-1 {
								infof("%s (entry of %s): prerequisite '%s' from %s skips module(s) — verify intentional",
									entryID, mod.ID, prereqID, prereqItem.ModuleID)
							}
						}
					}
				}
			}
		}
	}

	// 4. Check terminal items have expected next_item_ids bridge (informational)
	for _, mod := range core.Modules {
		modNum := moduleOrder[mod.ID]
		for _, termID := range mod.TerminalItemIDs {
			if termItem, exists := itemMap[termID]; exists {
				if len(termItem.NextItemIDs) == 0 && modNum < 16 {
					infof("%s (terminal of %s): no next_item_ids bridge to next module — verify intentional",
						termID, mod.ID)
				}
			}
		}
	}

	if ok {
		fmt.Println("  All prerequisite closure checks passed.")
	} else {
		fmt.Println("  Prerequisite closure checks completed with warnings/errors.")
	}
}

func validateProjectBinding() {
	fmt.Println("=== validate-project-binding ===")
	ok := true

	var core CoreBundle
	readJSON("path.core.json", &core)

	var electives ElectiveBundle
	readJSON("path.electives.json", &electives)

	// Build combined module lookup (core + electives)
	allModules := make(map[string]bool)
	for _, m := range core.Modules {
		allModules[m.ID] = true
	}
	for _, m := range electives.Modules {
		allModules[m.ID] = true
	}

	// Read projects.json as raw to check fields
	projectsPath := filepath.Join(curriculumDir, "projects.json")
	projectsData, err := os.ReadFile(projectsPath)
	if err != nil {
		fatalf("cannot read projects.json: %v", err)
	}
	var projectsRaw map[string]any
	if err := json.Unmarshal(projectsData, &projectsRaw); err != nil {
		fatalf("cannot parse projects.json: %v", err)
	}

	projects, ok2 := projectsRaw["projects"].([]any)
	if !ok2 {
		fatalf("projects.json: 'projects' is not an array")
	}

	for _, p := range projects {
		proj, isObj := p.(map[string]any)
		if !isObj {
			errorf("projects.json: project entry is not a JSON object")
			ok = false
			continue
		}
		projID, hasID := proj["id"].(string)
		if !hasID {
			errorf("projects.json: project entry is missing a string 'id'")
			ok = false
			continue
		}

		// Check required fields
		requiredFields := []string{"id", "module_id", "title", "slug", "type", "prerequisites", "deliverables", "verification", "rubric"}
		for _, field := range requiredFields {
			if _, exists := proj[field]; !exists {
				errorf("project '%s': missing required field '%s'", projID, field)
				ok = false
			}
		}

		// Check module_id references a valid module
		if mid, exists := proj["module_id"]; exists {
			midStr, midOK := mid.(string)
			if !midOK || !strings.HasPrefix(midStr, "module-") {
				errorf("project '%s': invalid module_id format '%v'", projID, mid)
				ok = false
			} else if !allModules[midStr] {
				errorf("project '%s': module_id '%s' does not match any known module", projID, midStr)
				ok = false
			}
		}

		// Check prerequisites reference valid items
		if prereqs, exists := proj["prerequisites"]; exists {
			if prereqArr, isArr := prereqs.([]any); isArr {
				for _, prereq := range prereqArr {
					prereqID, isStr := prereq.(string)
					if !isStr {
						errorf("project '%s': prerequisite entry is not a string", projID)
						ok = false
						continue
					}
					found := false
					for _, item := range core.Items {
						if item.ID == prereqID {
							found = true
							break
						}
					}
					if !found {
						errorf("project '%s': prerequisite '%s' does not exist", projID, prereqID)
						ok = false
					}
				}
			}
		}

		// Check project has assessment
		if _, exists := proj["assessment_id"]; !exists {
			errorf("project '%s': missing assessment_id", projID)
			ok = false
		}
	}

	if ok {
		fmt.Println("  All project binding checks passed.")
	} else {
		fmt.Println("  Project binding checks completed with warnings/errors.")
	}
}

var zeroMagicPlaceholderPhrases = []string{
	"This lesson explains what problem",
	"as one step in the learner's path",
	"one step in the learner's path",
	"The lesson must explain the underlying mechanics",
	"The lesson shows how Go expresses",
	"mechanically without understanding",
	"in the context of professional Go",
}

func findZeroMagicPlaceholder(zm *ZeroMagic) (string, bool) {
	fields := []string{
		zm.ProblemSolved,
		zm.WhyItExists,
		zm.MentalModel,
		zm.UnderTheHood,
		zm.HowGoUsesIt,
		zm.RealWorldUsage,
		zm.ProofOfUnderstanding,
	}
	for _, text := range fields {
		for _, phrase := range zeroMagicPlaceholderPhrases {
			if strings.Contains(text, phrase) {
				return phrase, true
			}
		}
	}
	return "", false
}

func validateNoPlaceholderZM() {
	fmt.Println("=== validate-no-placeholder-zm ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	placeholderCount := 0
	for _, item := range allItems {
		if item.ZeroMagic == nil {
			continue
		}
		zm := item.ZeroMagic

		if phrase, found := findZeroMagicPlaceholder(zm); found {
			warnf("%s: zero_magic contains placeholder text: %q", item.ID, phrase)
			placeholderCount++
			ok = false
		}

		fields := []struct {
			name string
			val  string
		}{
			{"problem_solved", zm.ProblemSolved},
			{"why_it_exists", zm.WhyItExists},
			{"mental_model", zm.MentalModel},
			{"under_the_hood", zm.UnderTheHood},
			{"how_go_uses_it", zm.HowGoUsesIt},
			{"real_world_usage", zm.RealWorldUsage},
			{"proof_of_understanding", zm.ProofOfUnderstanding},
		}
		for _, f := range fields {
			if strings.TrimSpace(f.val) == "" {
				warnf("%s: zero_magic.%s is empty", item.ID, f.name)
				placeholderCount++
				ok = false
			}
		}
	}

	if placeholderCount == 0 {
		fmt.Println("  No placeholder zero_magic content found.")
	} else {
		fmt.Printf("  Found %d items with placeholder zero_magic content.\n", placeholderCount)
	}

	if ok {
		fmt.Println("  All zero_magic checks passed.")
	} else {
		fmt.Println("  Zero-magic checks completed with warnings/errors.")
	}
}

func validateAssessmentCompleteness() {
	fmt.Println("=== validate-assessment-completeness ===")
	ok := true

	assessPath := filepath.Join(curriculumDir, "assessments.json")
	data, err := os.ReadFile(assessPath)
	if err != nil {
		fatalf("cannot read assessments.json: %v", err)
	}

	var rawMap map[string]any
	if err := json.Unmarshal(data, &rawMap); err != nil {
		fatalf("cannot parse assessments.json: %v", err)
	}

	assessments, ok2 := rawMap["assessments"].([]any)
	if !ok2 {
		fatalf("assessments.json: 'assessments' is not an array")
	}

	for _, a := range assessments {
		assess, isObj := a.(map[string]any)
		if !isObj {
			errorf("assessments.json: assessment entry is not a JSON object")
			ok = false
			continue
		}
		assessID, hasID := assess["id"].(string)
		if !hasID {
			errorf("assessments.json: assessment entry is missing a string 'id'")
			ok = false
			continue
		}

		// Check has criteria
		criteria, hasCriteria := assess["criteria"]
		if !hasCriteria {
			errorf("%s: missing criteria", assessID)
			ok = false
			continue
		}

		criteriaArr, isArr := criteria.([]any)
		if !isArr || len(criteriaArr) == 0 {
			errorf("%s: criteria is empty or not an array", assessID)
			ok = false
			continue
		}

		// Check for generic rubric (criteria without module-specific context)
		genericPatterns := []string{"general", "overall", "rubric"}
		for i, c := range criteriaArr {
			crit, isObj := c.(map[string]any)
			if !isObj {
				errorf("%s: criteria[%d] is not a JSON object", assessID, i)
				ok = false
				continue
			}
			critName := ""
			if name, exists := crit["name"]; exists {
				if s, isStr := name.(string); isStr {
					critName = s
				} else {
					errorf("%s: criteria[%d] 'name' is not a string", assessID, i)
					ok = false
				}
			}
			for _, pattern := range genericPatterns {
				if strings.Contains(strings.ToLower(critName), pattern) {
					warnf("%s: criteria[%d] '%s' appears generic (contains '%s')", assessID, i, critName, pattern)
					ok = false
				}
			}

			// Check criteria has weight/score
			if _, hasWeight := crit["weight"]; !hasWeight {
				errorf("%s: criteria[%d] '%s' missing weight", assessID, i, critName)
				ok = false
			}
		}

		// Check has passing score
		passingScore, hasPassing := assess["passing_score"]
		if !hasPassing {
			errorf("%s: missing passing_score", assessID)
			ok = false
		} else {
			if ps, ok := passingScore.(float64); ok && ps <= 0 {
				errorf("%s: passing_score is 0 or negative", assessID)
				ok = false
			}
		}

		// Check has target_ids
		targets, hasTargets := assess["target_ids"]
		if !hasTargets {
			errorf("%s: missing target_ids", assessID)
			ok = false
		} else {
			targetArr, ok := targets.([]any)
			if !ok || len(targetArr) == 0 {
				errorf("%s: target_ids is empty", assessID)
				ok = false
			}
		}

		// Check has evidence_required
		if _, hasEvidence := assess["evidence_required"]; !hasEvidence {
			errorf("%s: missing evidence_required", assessID)
			ok = false
		}

		// Check has retake_policy
		if _, hasRetake := assess["retake_policy"]; !hasRetake {
			errorf("%s: missing retake_policy", assessID)
			ok = false
		}
	}

	if ok {
		fmt.Println("  All assessment completeness checks passed.")
	} else {
		fmt.Println("  Assessment completeness checks completed with warnings/errors.")
	}
}

func validateUniqueTags() {
	fmt.Println("=== validate-unique-tags ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	// Check module-level tags
	globalTags := make(map[string][]string)
	for _, item := range allItems {
		for _, tag := range item.Tags {
			globalTags[tag] = append(globalTags[tag], item.ID)
		}
	}

	// Check module-level tags
	for _, mod := range core.Modules {
		for _, tag := range mod.Tags {
			globalTags[tag] = append(globalTags[tag], mod.ID)
		}
	}

	for tag, users := range globalTags {
		// Flag exact duplicate tags on the same item (error).
		seen := make(map[string]bool)
		for _, user := range users {
			if seen[user] {
				errorf("tag '%s' duplicated on item '%s'", tag, user)
				ok = false
			}
			seen[user] = true
		}
		// Surface globally-shared category tags as informational. Category
		// markers (e.g. "concept", "alpha") legitimately appear on many items,
		// but maintainers should be able to see which tags are shared.
		if len(users) > 1 {
			infof("tag '%s' shared across %d item(s): %v", tag, len(users), users)
		}
	}

	if ok {
		fmt.Println("  All unique-tag checks passed.")
	} else {
		fmt.Println("  Unique-tag checks completed with warnings/errors.")
	}
}

func validateConceptOwnership() {
	fmt.Println("=== validate-concept-ownership ===")
	ok := true

	var concepts struct {
		SchemaVersion string `json:"schema_version"`
		Concepts      []struct {
			Concept                string   `json:"concept"`
			CanonicalOwner         string   `json:"canonical_owner"`
			PreviewLocations       []string `json:"preview_locations"`
			ReinforcementLocations []string `json:"reinforcement_locations"`
		} `json:"concepts"`
	}
	readJSON("concepts.json", &concepts)
	if len(concepts.Concepts) == 0 {
		errorf("concepts.json: no concept entries found")
		ok = false
	}

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	// Build valid item IDs and valid module IDs
	validIDs := make(map[string]bool)
	validModuleIDs := make(map[string]bool)
	for _, item := range allItems {
		validIDs[item.ID] = true
	}
	for _, m := range core.Modules {
		validModuleIDs[m.ID] = true
	}

	// Check each concept entry
	conceptNames := make(map[string]string)
	for _, c := range concepts.Concepts {
		// Check for duplicate concept names
		if existing, dup := conceptNames[c.Concept]; dup {
			errorf("concept '%s' is duplicated (first at canonical_owner: %s)", c.Concept, existing)
			ok = false
			continue
		}
		conceptNames[c.Concept] = c.CanonicalOwner

		// Check canonical_owner resolves
		if !validIDs[c.CanonicalOwner] {
			errorf("concept '%s': canonical_owner '%s' does not exist in item graph", c.Concept, c.CanonicalOwner)
			ok = false
		}

		// Check preview_locations resolve
		for _, loc := range c.PreviewLocations {
			if !validIDs[loc] && !validModuleIDs[loc] {
				errorf("concept '%s': preview_location '%s' does not exist", c.Concept, loc)
				ok = false
			}
		}

		// Check reinforcement_locations resolve
		for _, loc := range c.ReinforcementLocations {
			if !validIDs[loc] && !validModuleIDs[loc] {
				errorf("concept '%s': reinforcement_location '%s' does not exist", c.Concept, loc)
				ok = false
			}
		}
	}

	if ok {
		fmt.Printf("  All concept ownership checks passed (%d concepts).\n", len(concepts.Concepts))
	} else {
		fmt.Printf("  Concept ownership checks completed with warnings/errors.\n")
	}
}

func validateCognitiveLoad() {
	fmt.Println("=== validate-cognitive-load ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	// Read raw module data to check fields
	var rawCore map[string]any
	readJSON("path.core.json", &rawCore)
	var rawElectives map[string]any
	readJSON("path.electives.json", &rawElectives)

	checkModules := func(raw map[string]any) {
		mods, exists := raw["modules"].([]any)
		if !exists {
			return
		}
		for _, m := range mods {
			mm, isObj := m.(map[string]any)
			if !isObj {
				errorf("modules: module entry is not a JSON object")
				ok = false
				continue
			}
			id, hasID := mm["id"].(string)
			if !hasID {
				errorf("modules: module entry is missing a string 'id'")
				ok = false
				continue
			}

			cl, hasCL := mm["cognitive_load"]
			if !hasCL {
				errorf("%s: missing cognitive_load field", id)
				ok = false
				continue
			}
			clStr, ok := cl.(string)
			if !ok || (clStr != "low" && clStr != "moderate" && clStr != "high") {
				errorf("%s: cognitive_load must be 'low', 'moderate', or 'high', got %v", id, cl)
				ok = false
			}

			rba, hasRBA := mm["recommended_break_after"]
			if !hasRBA {
				errorf("%s: missing recommended_break_after", id)
				ok = false
			} else if _, ok := rba.(bool); !ok {
				errorf("%s: recommended_break_after must be bool", id)
				ok = false
			}

			cfhc, hasCFHC := mm["contains_foundational_hard_concepts"]
			if !hasCFHC {
				errorf("%s: missing contains_foundational_hard_concepts", id)
				ok = false
			} else if _, ok := cfhc.(bool); !ok {
				errorf("%s: contains_foundational_hard_concepts must be bool", id)
				ok = false
			}

			p, hasP := mm["pacing"]
			if !hasP {
				errorf("%s: missing pacing field", id)
				ok = false
			} else {
				pStr, ok := p.(string)
				if !ok || (pStr != "gentle" && pStr != "moderate" && pStr != "steady") {
					errorf("%s: pacing must be 'gentle', 'moderate', or 'steady', got %v", id, p)
					ok = false
				}
			}

			// Check consistency: high load should have break_after
			if clStr == "high" {
				if rba, ok := mm["recommended_break_after"].(bool); ok && !rba {
					warnf("%s: high cognitive_load should typically have recommended_break_after=true", id)
				}
			}
		}
	}

	checkModules(rawCore)
	checkModules(rawElectives)

	if ok {
		fmt.Println("  All cognitive load checks passed.")
	} else {
		fmt.Println("  Cognitive load checks completed with warnings/errors.")
	}
}

func validateZeroMagic() {
	fmt.Println("=== validate-zero-magic ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	goldenCount := 0
	for _, item := range allItems {
		if item.ZeroMagicStatus != "golden" {
			continue
		}
		goldenCount++
		if item.ZeroMagic == nil {
			errorf("%s: golden lesson missing zero_magic block", item.ID)
			ok = false
			continue
		}
		zm := item.ZeroMagic

		// Check for placeholder content that indicates un-authored zero_magic
		if phrase, found := findZeroMagicPlaceholder(zm); found {
			errorf("%s: zero_magic contains placeholder text: %q", item.ID, phrase)
			ok = false
		}

		// Required spec fields for golden lessons
		if len(zm.StepByStepExecution) == 0 {
			infof("%s: golden lesson missing step_by_step_execution — recommended by authoring spec", item.ID)
		}
		if len(zm.FailureModes) == 0 && len(zm.OperationalFailureExamples) == 0 {
			errorf("%s: golden lesson missing failure_modes or operational_failure_examples", item.ID)
			ok = false
		}
		if len(zm.HiddenMagicChecks) == 0 {
			infof("%s: golden lesson missing hidden_magic_checks — recommended by authoring spec", item.ID)
		}
		if len(zm.BeginnerMistakes) == 0 {
			errorf("%s: golden lesson missing beginner_mistakes", item.ID)
			ok = false
		}
		if zm.UnderTheHood == "" || strings.Contains(zm.UnderTheHood, "must explain") {
			errorf("%s: under_the_hood is placeholder or empty", item.ID)
			ok = false
		}
		if zm.MentalModel == "" || strings.Contains(zm.MentalModel, "one step in the learner's path") {
			errorf("%s: mental_model is placeholder or empty", item.ID)
			ok = false
		}
	}

	if ok {
		fmt.Printf("  All zero-magic checks passed (%d golden lessons verified).\n", goldenCount)
	} else {
		fmt.Println("  Zero-magic checks completed with warnings/errors.")
	}
}

func validateFailureEngineering() {
	fmt.Println("=== validate-failure-engineering ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	// Modules with required failure coverage (from failures.json)
	coverageModules := map[string]bool{
		"module-08": true, "module-09": true, "module-10": true,
		"module-11": true, "module-12": true, "module-14": true,
		"module-15": true, "module-18": true,
	}

	for _, item := range allItems {
		if !coverageModules[item.ModuleID] {
			continue
		}
		if item.ZeroMagic == nil {
			// Not all items in coverage modules need zero_magic (e.g. projects)
			continue
		}
		zm := item.ZeroMagic
		if len(zm.FailureModes) == 0 && len(zm.OperationalFailureExamples) == 0 {
			errorf("%s: failure-engineering module item missing failure_modes or operational_failure_examples", item.ID)
			ok = false
		}
	}

	if ok {
		fmt.Println("  All failure-engineering checks passed.")
	} else {
		fmt.Println("  Failure-engineering checks completed with warnings/errors.")
	}
}

func validateOperationalFailures() {
	fmt.Println("=== validate-operational-failures ===")
	ok := true

	var failures struct {
		SchemaVersion     string `json:"schema_version"`
		DocumentType      string `json:"document_type"`
		CurriculumVersion string `json:"curriculum_version"`
		FailureCategories []struct {
			Category    string   `json:"category"`
			Description string   `json:"description"`
			Modules     []string `json:"modules"`
		} `json:"failure_categories"`
		RequiredCoverage map[string][]string `json:"required_coverage"`
	}
	readJSON("failures.json", &failures)

	if failures.DocumentType != "operational_failure_taxonomy" {
		errorf("failures.json: document_type should be 'operational_failure_taxonomy', got '%s'", failures.DocumentType)
		ok = false
	}

	// Build set of valid module IDs from core path + electives
	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	validModuleIDs := make(map[string]bool)
	for _, m := range core.Modules {
		validModuleIDs[m.ID] = true
	}
	for _, m := range electives.Modules {
		validModuleIDs[m.ID] = true
	}

	// 1. Check all category names are unique
	categoryNames := make(map[string]string)
	for _, fc := range failures.FailureCategories {
		if _, dup := categoryNames[fc.Category]; dup {
			errorf("failures.json: duplicate category '%s'", fc.Category)
			ok = false
			continue
		}
		categoryNames[fc.Category] = fc.Description

		// Check module references in categories resolve
		for _, modID := range fc.Modules {
			if !validModuleIDs[modID] {
				errorf("failures.json: category '%s' references unknown module '%s'", fc.Category, modID)
				ok = false
			}
		}
	}

	// 2. Check all required_coverage entries resolve
	for modID, categories := range failures.RequiredCoverage {
		if !validModuleIDs[modID] {
			errorf("failures.json: required_coverage references unknown module '%s'", modID)
			ok = false
			continue
		}
		if len(categories) == 0 {
			errorf("failures.json: required_coverage for '%s' is empty", modID)
			ok = false
			continue
		}
		for _, cat := range categories {
			if _, exists := categoryNames[cat]; !exists {
				errorf("failures.json: required_coverage for '%s' references unknown category '%s'", modID, cat)
				ok = false
			}
		}
	}

	if ok {
		fmt.Printf("  All operational failure checks passed (%d categories, %d module coverage entries).\n",
			len(failures.FailureCategories), len(failures.RequiredCoverage))
	} else {
		fmt.Println("  Operational failure checks completed with warnings/errors.")
	}
}

func validateCrossrefs() {
	fmt.Println("=== validate-crossrefs ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)

	var crossrefs CrossrefBundle
	readJSON("crossrefs.json", &crossrefs)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	itemMap := make(map[string]*Item)
	for i := range allItems {
		itemMap[allItems[i].ID] = &allItems[i]
	}

	// Add project IDs as valid targets
	var projects ProjectsBundle
	readJSON("projects.json", &projects)
	projectIDs := make(map[string]bool)
	for _, p := range projects.Projects {
		projectIDs[p.ID] = true
	}

	// 1. Check all crossrefs.json references resolve
	for _, ref := range crossrefs.Crossrefs.References {
		if _, exists := itemMap[ref.FromID]; !exists && !projectIDs[ref.FromID] {
			errorf("crossrefs.json: from_id '%s' does not exist", ref.FromID)
			ok = false
		}
		if _, exists := itemMap[ref.TargetID]; !exists && !projectIDs[ref.TargetID] {
			errorf("crossrefs.json: target_id '%s' does not exist", ref.TargetID)
			ok = false
		}
	}

	// 2. Check item-level crossrefs resolve
	moduleOrder := make(map[string]int)
	for _, m := range core.Modules {
		moduleOrder[m.ID] = m.Number
	}
	for _, m := range electives.Modules {
		moduleOrder[m.ID] = m.Number
	}

	for _, item := range allItems {
		if item.CrossRefs == nil {
			continue
		}
		checkRefs := func(refs []CrossrefRef, category string) {
			for _, ref := range refs {
				if _, exists := itemMap[ref.TargetID]; !exists && !projectIDs[ref.TargetID] {
					errorf("%s: crossrefs.%s target '%s' does not exist", item.ID, category, ref.TargetID)
					ok = false
				}
			}
		}
		checkRefs(item.CrossRefs.BuildsOn, "builds_on")
		checkRefs(item.CrossRefs.PreviewOnly, "preview_only")
		checkRefs(item.CrossRefs.Related, "related")
		checkRefs(item.CrossRefs.ReinforcedIn, "reinforced_in")

		// 3. Check preview_only links point to valid future items
		itemModNum := moduleOrder[item.ModuleID]
		for _, ref := range item.CrossRefs.PreviewOnly {
			if targetItem, exists := itemMap[ref.TargetID]; exists && targetItem != nil {
				targetModNum := moduleOrder[targetItem.ModuleID]
				if targetModNum <= itemModNum {
					errorf("%s: preview_only target '%s' is in module %s (order %d), not in a future module (item at order %d)",
						item.ID, ref.TargetID, targetItem.ModuleID, targetModNum, itemModNum)
					ok = false
				}
			}
		}
	}

	// 4. Reconcile core item-embedded crossrefs with the central registry.
	// crossrefs.json is authoritative for core items. Electives carry their
	// crossrefs only in path.electives.json, so they are checked for target
	// existence above but are not part of the central registry.
	electiveIDs := make(map[string]bool)
	for _, item := range electives.Items {
		electiveIDs[item.ID] = true
	}
	central := make(map[string]map[string]map[string]bool)
	for _, ref := range crossrefs.Crossrefs.References {
		if itemMap[ref.FromID] == nil {
			continue // project-origin references are not embedded in Item objects
		}
		if central[ref.FromID] == nil {
			central[ref.FromID] = make(map[string]map[string]bool)
		}
		if central[ref.FromID][ref.Relation] == nil {
			central[ref.FromID][ref.Relation] = make(map[string]bool)
		}
		central[ref.FromID][ref.Relation][ref.TargetID] = true
	}
	relations := []string{"builds_on", "preview_only", "related", "reinforced_in"}
	for _, item := range allItems {
		if electiveIDs[item.ID] || item.CrossRefs == nil {
			continue
		}
		embedded := map[string][]CrossrefRef{
			"builds_on":     item.CrossRefs.BuildsOn,
			"preview_only":  item.CrossRefs.PreviewOnly,
			"related":       item.CrossRefs.Related,
			"reinforced_in": item.CrossRefs.ReinforcedIn,
		}
		for _, relation := range relations {
			centralTargets := central[item.ID][relation]
			embeddedTargets := make(map[string]bool)
			for _, ref := range embedded[relation] {
				embeddedTargets[ref.TargetID] = true
				if !centralTargets[ref.TargetID] {
					errorf("%s: crossrefs.%s target '%s' is missing from crossrefs.json", item.ID, relation, ref.TargetID)
					ok = false
				}
			}
			for targetID := range centralTargets {
				if !embeddedTargets[targetID] {
					errorf("%s: crossrefs.%s target '%s' is missing from the item manifest", item.ID, relation, targetID)
					ok = false
				}
			}
		}
	}

	if ok {
		fmt.Println("  All crossref checks passed.")
	} else {
		fmt.Println("  Crossref checks completed with warnings/errors.")
	}
}
