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
				if prereqModNum > itemModNum {
					warnf("%s: prerequisite '%s' is in module %s (order %d) which comes AFTER %s (order %d)",
						item.ID, prereqID, prereqItem.ModuleID, prereqModNum, item.ModuleID, itemModNum)
					ok = false
				}
			}
		}
	}

	// 2. Check no hidden prerequisite leaks (prereqs that don't exist as items)
	for _, item := range allItems {
		for _, prereqID := range item.Prerequisites {
			if _, exists := itemMap[prereqID]; !exists {
				warnf("%s: prerequisite '%s' does not exist in item graph", item.ID, prereqID)
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
		proj := p.(map[string]any)
		projID := proj["id"].(string)
		title := proj["title"].(string)

		// Check required fields
		requiredFields := []string{"id", "module_id", "title", "slug", "type", "prerequisites", "deliverables", "verification", "rubric"}
		for _, field := range requiredFields {
			if _, exists := proj[field]; !exists {
				warnf("project '%s': missing required field '%s'", projID, field)
				ok = false
			}
		}

		// Check module_id references a valid module
		if mid, exists := proj["module_id"]; exists {
			midStr, midOK := mid.(string)
			if !midOK || !strings.HasPrefix(midStr, "module-") {
				warnf("project '%s': invalid module_id format '%v'", projID, mid)
				ok = false
			} else if !allModules[midStr] {
				warnf("project '%s': module_id '%s' does not match any known module", projID, midStr)
				ok = false
			}
		}

		// Check prerequisites reference valid items
		if prereqs, exists := proj["prerequisites"]; exists {
			if prereqArr, ok := prereqs.([]any); ok {
				for _, prereq := range prereqArr {
					prereqID := prereq.(string)
					found := false
					for _, item := range core.Items {
						if item.ID == prereqID {
							found = true
							break
						}
					}
					if !found {
						warnf("project '%s': prerequisite '%s' does not exist", projID, prereqID)
						ok = false
					}
				}
			}
		}

		// Check project has assessment
		if _, exists := proj["assessment_id"]; !exists {
			warnf("project '%s': missing assessment_id", projID)
			ok = false
		}

		_ = title
	}

	if ok {
		fmt.Println("  All project binding checks passed.")
	} else {
		fmt.Println("  Project binding checks completed with warnings/errors.")
	}
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

	placeholderPhrases := []string{
		"This lesson explains what problem",
		"as one step in the learner's path",
		"The lesson must explain the underlying mechanics",
		"The lesson shows how Go expresses",
		"mechanically without understanding",
	}

	placeholderCount := 0
	for _, item := range allItems {
		if item.ZeroMagic == nil {
			continue
		}

		// Marshal to JSON string to search all fields
		zmJSON, _ := json.Marshal(item.ZeroMagic)
		zmStr := string(zmJSON)

		for _, phrase := range placeholderPhrases {
			if strings.Contains(zmStr, phrase) {
				warnf("%s: zero_magic contains placeholder text: %q", item.ID, phrase)
				placeholderCount++
				ok = false
				break
			}
		}

		// Check specific fields using reflection via JSON marshal
		var zmMap map[string]any
		json.Unmarshal(zmJSON, &zmMap)
		for _, field := range []string{"problem_solved", "why_it_exists", "mental_model", "under_the_hood", "how_go_uses_it", "real_world_usage", "proof_of_understanding"} {
			if val, exists := zmMap[field]; exists {
				if str, ok := val.(string); ok && strings.TrimSpace(str) == "" {
					warnf("%s: zero_magic.%s is empty", item.ID, field)
					placeholderCount++
					ok = false
				}
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
		assess := a.(map[string]any)
		assessID := assess["id"].(string)

		// Check has criteria
		criteria, hasCriteria := assess["criteria"]
		if !hasCriteria {
			warnf("%s: missing criteria", assessID)
			ok = false
			continue
		}

		criteriaArr, isArr := criteria.([]any)
		if !isArr || len(criteriaArr) == 0 {
			warnf("%s: criteria is empty or not an array", assessID)
			ok = false
			continue
		}

		// Check for generic rubric (criteria without module-specific context)
		genericPatterns := []string{"general", "overall", "rubric"}
		for i, c := range criteriaArr {
			crit := c.(map[string]any)
			critName := ""
			if name, exists := crit["name"]; exists {
				critName = name.(string)
			}
			for _, pattern := range genericPatterns {
				if strings.Contains(strings.ToLower(critName), pattern) {
					warnf("%s: criteria[%d] '%s' appears generic (contains '%s')", assessID, i, critName, pattern)
					ok = false
				}
			}

			// Check criteria has weight/score
			if _, hasWeight := crit["weight"]; !hasWeight {
				warnf("%s: criteria[%d] '%s' missing weight", assessID, i, critName)
				ok = false
			}
		}

		// Check has passing score
		passingScore, hasPassing := assess["passing_score"]
		if !hasPassing {
			warnf("%s: missing passing_score", assessID)
			ok = false
		} else {
			if ps, ok := passingScore.(float64); ok && ps <= 0 {
				warnf("%s: passing_score is 0 or negative", assessID)
				ok = false
			}
		}

		// Check has target_ids
		targets, hasTargets := assess["target_ids"]
		if !hasTargets {
			warnf("%s: missing target_ids", assessID)
			ok = false
		} else {
			targetArr, ok := targets.([]any)
			if !ok || len(targetArr) == 0 {
				warnf("%s: target_ids is empty", assessID)
				ok = false
			}
		}

		// Check has evidence_required
		if _, hasEvidence := assess["evidence_required"]; !hasEvidence {
			warnf("%s: missing evidence_required", assessID)
			ok = false
		}

		// Check has retake_policy
		if _, hasRetake := assess["retake_policy"]; !hasRetake {
			warnf("%s: missing retake_policy", assessID)
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
		if len(users) > 1 {
			// Tags can legitimately appear on multiple items (they're category markers)
			// But flag exact duplicates to be sure
			seen := make(map[string]bool)
			for _, user := range users {
				if seen[user] {
					warnf("tag '%s' duplicated on item '%s'", tag, user)
					ok = false
				}
				seen[user] = true
			}
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
		warnf("concepts.json: no concept entries found")
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
			warnf("concept '%s' is duplicated (first at canonical_owner: %s)", c.Concept, existing)
			ok = false
			continue
		}
		conceptNames[c.Concept] = c.CanonicalOwner

		// Check canonical_owner resolves
		if !validIDs[c.CanonicalOwner] {
			warnf("concept '%s': canonical_owner '%s' does not exist in item graph", c.Concept, c.CanonicalOwner)
			ok = false
		}

		// Check preview_locations resolve
		for _, loc := range c.PreviewLocations {
			if !validIDs[loc] && !validModuleIDs[loc] {
				warnf("concept '%s': preview_location '%s' does not exist", c.Concept, loc)
				ok = false
			}
		}

		// Check reinforcement_locations resolve
		for _, loc := range c.ReinforcementLocations {
			if !validIDs[loc] && !validModuleIDs[loc] {
				warnf("concept '%s': reinforcement_location '%s' does not exist", c.Concept, loc)
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

	type moduleLoad struct {
		ID                               string
		CognitiveLoad                    string
		RecommendedBreakAfter            bool
		ContainsFoundationalHardConcepts bool
		Pacing                           string
	}
	var allMods []moduleLoad

	// Use raw JSON to read full module fields (our struct may not have them)
	for range core.Modules {
		// We use raw JSON read separately for the full field set
	}
	_ = allMods

	// Read raw module data to check fields
	var rawCore map[string]any
	readJSON("path.core.json", &rawCore)
	var rawElectives map[string]any
	readJSON("path.electives.json", &rawElectives)

	checkModules := func(raw map[string]any, source string) {
		mods, exists := raw["modules"].([]any)
		if !exists {
			return
		}
		for _, m := range mods {
			mm := m.(map[string]any)
			id := mm["id"].(string)

			cl, hasCL := mm["cognitive_load"]
			if !hasCL {
				warnf("%s: missing cognitive_load field", id)
				ok = false
				continue
			}
			clStr, ok := cl.(string)
			if !ok || (clStr != "low" && clStr != "moderate" && clStr != "high") {
				warnf("%s: cognitive_load must be 'low', 'moderate', or 'high', got %v", id, cl)
				ok = false
			}

			rba, hasRBA := mm["recommended_break_after"]
			if !hasRBA {
				warnf("%s: missing recommended_break_after", id)
				ok = false
			} else if _, ok := rba.(bool); !ok {
				warnf("%s: recommended_break_after must be bool", id)
				ok = false
			}

			cfhc, hasCFHC := mm["contains_foundational_hard_concepts"]
			if !hasCFHC {
				warnf("%s: missing contains_foundational_hard_concepts", id)
				ok = false
			} else if _, ok := cfhc.(bool); !ok {
				warnf("%s: contains_foundational_hard_concepts must be bool", id)
				ok = false
			}

			p, hasP := mm["pacing"]
			if !hasP {
				warnf("%s: missing pacing field", id)
				ok = false
			} else {
				pStr, ok := p.(string)
				if !ok || (pStr != "gentle" && pStr != "moderate" && pStr != "steady") {
					warnf("%s: pacing must be 'gentle', 'moderate', or 'steady', got %v", id, p)
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

	checkModules(rawCore, "path.core.json")
	checkModules(rawElectives, "path.electives.json")

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

	// Golden lessons — exemplars that set the authoring standard
	goldenLessons := map[string]struct{}{
		"core-03-15": {}, // Slice length and capacity
		"core-03-16": {}, // Slice sharing and aliasing
		"core-05-08": {}, // Interfaces
		"core-05-14": {}, // Nil interfaces
		"core-11-06": {}, // Cancellation
		// Phase 2 golden lessons — module 04 fundamentals
		"core-04-08": {}, // Pointer and value mutation behavior
		"core-04-09": {}, // Errors as values
		"core-04-15": {}, // defer mechanics
		"core-04-17": {}, // panic and recover
		// Phase 2 golden lessons — module 05 type system
		"core-05-05": {}, // Receiver sets (method sets)
		// Phase 2 golden lessons — module 08 HTTP
		"core-08-03": {}, // Handler lifecycle
		"core-08-05": {}, // Request parsing
		"core-08-07": {}, // Response writing
		// Phase 2 golden lessons — module 09 databases
		"core-09-11": {}, // sql.DB as a pool
		"core-09-17": {}, // Transactions
		// Phase 2 golden lessons — module 11 concurrency
		"core-11-11": {}, // Goroutines
		"core-11-16": {}, // Channel ownership
		"core-11-23": {}, // Goroutine leaks
		// Phase 2 golden lessons — module 12 observability
		"core-12-02": {}, // Structured logging with slog
		"core-12-03": {}, // Request-scoped logging
		"core-12-04": {}, // Correlation IDs
		"core-12-05": {}, // PII redaction
		"core-12-06": {}, // Metrics
		"core-12-07": {}, // Prometheus
		"core-12-08": {}, // Tracing
		"core-12-09": {}, // OpenTelemetry
		"core-12-18": {}, // Alerting mindset
		"core-12-20": {}, // Reliability review
		// Phase 2 golden lessons — module 13 performance/memory
		"core-13-10": {}, // pprof
		"core-13-11": {}, // CPU profiling
		"core-13-12": {}, // Memory profiling
		"core-13-13": {}, // Benchmarks
		"core-13-14": {}, // Escape analysis
		"core-13-15": {}, // Memory layout
		"core-13-16": {}, // Caching basics
		"core-13-17": {}, // Cache invalidation
		// Phase 2 golden lessons — module 14 architecture & distributed systems
		"core-14-01": {}, // Why architecture exists
		"core-14-02": {}, // Package boundaries
		"core-14-03": {}, // Service layer
		"core-14-04": {}, // Repository pattern deep dive
		"core-14-05": {}, // Modular monolith
		"core-14-06": {}, // Hexagonal architecture
		"core-14-07": {}, // Domain modeling
		"core-14-08": {}, // Invariants
		"core-14-09": {}, // Event-driven basics
		"core-14-10": {}, // Queues
		"core-14-11": {}, // Retries and idempotent consumers
		"core-14-12": {}, // Payment workflow design
		"core-14-13": {}, // Multi-tenancy architecture
		"core-14-14": {}, // Caching architecture
		"core-14-15": {}, // When to split services
		"core-14-16": {}, // Microservices as a tradeoff, not a default
		"core-14-17": {}, // Architecture decision records
		// Phase 3 golden lessons — module 15 Docker, CI/CD, deployment
		"core-15-01": {}, // Linux processes
		"core-15-02": {}, // Signals
		"core-15-03": {}, // Logs
		"core-15-04": {}, // File permissions
		"core-15-05": {}, // Networking basics
		"core-15-06": {}, // Docker basics
		"core-15-07": {}, // Docker images and layers
		"core-15-08": {}, // Multi-stage builds
		"core-15-09": {}, // Docker Compose
		"core-15-10": {}, // Container health checks
		"core-15-11": {}, // GitHub Actions
		"core-15-12": {}, // CI pipeline
		"core-15-13": {}, // Release artifacts
		"core-15-14": {}, // Config in deployment
		"core-15-15": {}, // Secrets in deployment
		"core-15-16": {}, // Vulnerability scanning
		"core-15-17": {}, // One deployment target
		"core-15-18": {}, // Rollback basics
		"core-15-19": {}, // Deployment runbook
		// Phase 4 — module 16 career and portfolio
		"core-16-01": {}, // README writing
		"core-16-02": {}, // API docs
		"core-16-03": {}, // Code review workflow
		"core-16-04": {}, // Pull request etiquette
		"core-16-05": {}, // Portfolio packaging
		"core-16-06": {}, // Writing project narratives
		"core-16-07": {}, // Resume project bullets
		"core-16-08": {}, // Backend interview fundamentals
		"core-16-09": {}, // Go interview topics
		"core-16-10": {}, // SQL interview topics
		"core-16-11": {}, // Debugging interview scenarios
		"core-16-12": {}, // System design basics
		"core-16-13": {}, // Architecture defense
		"core-16-14": {}, // Capstone planning
		"core-16-15": {}, // Demo preparation
		// Critical failure-engineering items
		"core-12-01": {},
		"core-12-19": {},
		"opslane-05": {},
		"opslane-07": {},
		"opslane-13": {},
	}

	placeholderPhrases := []string{
		"mechanically without understanding",
		"one step in the learner's path",
		"in the context of professional Go",
	}

	for _, item := range allItems {
		if _, isGolden := goldenLessons[item.ID]; !isGolden {
			continue
		}
		if item.ZeroMagic == nil {
			warnf("%s: golden lesson missing zero_magic block", item.ID)
			ok = false
			continue
		}
		zm := item.ZeroMagic

		// Check for placeholder content that indicates un-authored zero_magic
		zmJSON, _ := json.Marshal(zm)
		zmStr := string(zmJSON)
		for _, phrase := range placeholderPhrases {
			if strings.Contains(zmStr, phrase) {
				warnf("%s: zero_magic contains placeholder text: %q", item.ID, phrase)
				ok = false
				break
			}
		}

		// Required spec fields for golden lessons
		if len(zm.StepByStepExecution) == 0 {
			infof("%s: golden lesson missing step_by_step_execution — recommended by authoring spec", item.ID)
		}
		if len(zm.FailureModes) == 0 && len(zm.OperationalFailureExamples) == 0 {
			warnf("%s: golden lesson missing failure_modes or operational_failure_examples", item.ID)
			ok = false
		}
		if len(zm.HiddenMagicChecks) == 0 {
			infof("%s: golden lesson missing hidden_magic_checks — recommended by authoring spec", item.ID)
		}
		if len(zm.BeginnerMistakes) == 0 {
			warnf("%s: golden lesson missing beginner_mistakes", item.ID)
			ok = false
		}
		if zm.UnderTheHood == "" || strings.Contains(zm.UnderTheHood, "must explain") {
			warnf("%s: under_the_hood is placeholder or empty", item.ID)
			ok = false
		}
		if zm.MentalModel == "" || strings.Contains(zm.MentalModel, "one step in the learner's path") {
			warnf("%s: mental_model is placeholder or empty", item.ID)
			ok = false
		}
	}

	if ok {
		fmt.Printf("  All zero-magic checks passed (%d golden lessons verified).\n", len(goldenLessons))
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
			warnf("%s: failure-engineering module item missing failure_modes or operational_failure_examples", item.ID)
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
		warnf("failures.json: document_type should be 'operational_failure_taxonomy', got '%s'", failures.DocumentType)
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
	// Allow "opslane" as a valid module key in failures taxonomy
	validModuleIDs["opslane"] = true

	// 1. Check all category names are unique
	categoryNames := make(map[string]string)
	for _, fc := range failures.FailureCategories {
		if _, dup := categoryNames[fc.Category]; dup {
			warnf("failures.json: duplicate category '%s'", fc.Category)
			ok = false
			continue
		}
		categoryNames[fc.Category] = fc.Description

		// Check module references in categories resolve
		for _, modID := range fc.Modules {
			if !validModuleIDs[modID] {
				warnf("failures.json: category '%s' references unknown module '%s'", fc.Category, modID)
				ok = false
			}
		}
	}

	// 2. Check all required_coverage entries resolve
	for modID, categories := range failures.RequiredCoverage {
		if !validModuleIDs[modID] {
			warnf("failures.json: required_coverage references unknown module '%s'", modID)
			ok = false
			continue
		}
		if len(categories) == 0 {
			warnf("failures.json: required_coverage for '%s' is empty", modID)
			ok = false
			continue
		}
		for _, cat := range categories {
			if _, exists := categoryNames[cat]; !exists {
				warnf("failures.json: required_coverage for '%s' references unknown category '%s'", modID, cat)
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
			warnf("crossrefs.json: from_id '%s' does not exist", ref.FromID)
			ok = false
		}
		if _, exists := itemMap[ref.TargetID]; !exists && !projectIDs[ref.TargetID] {
			warnf("crossrefs.json: target_id '%s' does not exist", ref.TargetID)
			ok = false
		}
	}

	// 2. Check item-level crossrefs resolve
	for _, item := range allItems {
		if item.CrossRefs == nil {
			continue
		}
		// Marshal to check all fields
		crossJSON, _ := json.Marshal(item.CrossRefs)
		var crossMap map[string]any
		json.Unmarshal(crossJSON, &crossMap)

		for _, category := range []string{"builds_on", "preview_only", "related", "reinforced_in"} {
			refs, exists := crossMap[category]
			if !exists {
				continue
			}
			refArr, ok := refs.([]any)
			if !ok {
				continue
			}
			if len(refArr) == 0 && item.Status != "planned" {
				warnf("%s: crossrefs.%s is empty (item status: %s)", item.ID, category, item.Status)
				ok = false
			}
			for _, r := range refArr {
				ref := r.(map[string]any)
				if targetID, exists := ref["target_id"]; exists {
					tid := targetID.(string)
					if _, exists := itemMap[tid]; !exists && !projectIDs[tid] {
						warnf("%s: crossrefs.%s target '%s' does not exist", item.ID, category, tid)
						ok = false
					}
				}
			}
		}
	}

	// 3. Check preview_only links point to valid future items
	moduleOrder := make(map[string]int)
	for _, m := range core.Modules {
		moduleOrder[m.ID] = m.Number
	}
	for _, item := range allItems {
		if item.CrossRefs == nil {
			continue
		}
		crossJSON, _ := json.Marshal(item.CrossRefs)
		var crossMap map[string]any
		json.Unmarshal(crossJSON, &crossMap)

		if previewRefs, exists := crossMap["preview_only"]; exists {
			previewArr, ok := previewRefs.([]any)
			if !ok {
				continue
			}
			itemModNum := moduleOrder[item.ModuleID]
			for _, r := range previewArr {
				ref := r.(map[string]any)
				if targetID, exists := ref["target_id"]; exists {
					tid := targetID.(string)
					if targetItem, exists := itemMap[tid]; exists {
						targetModNum := moduleOrder[targetItem.ModuleID]
						if targetModNum <= itemModNum {
							warnf("%s: preview_only target '%s' is in module %s (order %d), not in a future module (item at order %d)",
								item.ID, tid, targetItem.ModuleID, targetModNum, itemModNum)
							ok = false
						}
					}
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
