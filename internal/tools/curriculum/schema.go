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

	// 1. Parse all canonical active JSON manifests and verify metadata consistency
	exp := loadExpected()
	reg := loadRegistry()

	// The contract registry's own rules are part of the invariants: they must be
	// modeled (not data-only), aliases must resolve unambiguously, and every
	// severity override must use a known tier.
	contractIDs := make(map[string]bool)
	for _, c := range reg.Contracts {
		if contractIDs[c.ID] {
			errorf("readme.contracts.json: duplicate contract id '%s'", c.ID)
			ok = false
		}
		contractIDs[c.ID] = true
	}
	aliasOwners := make(map[string]string)
	for _, c := range reg.Contracts {
		for _, a := range c.Selector.Aliases {
			if prev, dup := aliasOwners[a]; dup {
				errorf("readme.contracts.json: alias '%s' declared by both '%s' and '%s'", a, prev, c.ID)
				ok = false
			}
			aliasOwners[a] = c.ID
			if contractIDs[a] {
				errorf("readme.contracts.json: alias '%s' collides with a contract id", a)
				ok = false
			}
		}
	}
	if len(reg.ValidationRules) == 0 {
		errorf("readme.contracts.json: validation_rules must not be empty")
		ok = false
	}
	hasMarkerRule := false
	for _, rule := range reg.ValidationRules {
		if strings.Contains(rule, "curriculum:generated:start") {
			hasMarkerRule = true
			break
		}
	}
	if !hasMarkerRule {
		errorf("readme.contracts.json: validation_rules must document the generated-catalog marker rule")
		ok = false
	}
	if !validateRegistryFields(reg, exp) {
		ok = false
	}

	files := []string{
		"expected.json",
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
			errorf("cannot read %s: %v", f, err)
			ok = false
			continue
		}
		if !json.Valid(data) {
			errorf("%s: invalid JSON", f)
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
				errorf("%s: expected schema_version '3.1.0', got '%s'", f, header.SchemaVersion)
				ok = false
			}
			if header.CurriculumVersion != "3.0.0" {
				errorf("%s: expected curriculum_version '3.0.0', got '%s'", f, header.CurriculumVersion)
				ok = false
			}
			if header.Status != "stable" {
				errorf("%s: expected status 'stable', got '%s'", f, header.Status)
				ok = false
			}
			if header.LastUpdated != exp.LastUpdated {
				errorf("%s: expected last_updated '%s' (from expected.json), got '%s'", f, exp.LastUpdated, header.LastUpdated)
				ok = false
			}
		}
		fmt.Printf("  %s: valid JSON & metadata\n", f)
	}

	// Validate schema.v3.json (the JSON Schema definition)
	schemaPath := filepath.Join(curriculumDir, "schema.v3.json")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		errorf("cannot read schema.v3.json: %v", err)
		ok = false
	} else if !json.Valid(schemaData) {
		errorf("schema.v3.json: invalid JSON")
		ok = false
	} else {
		var schemaMeta struct {
			ID     string `json:"$id"`
			Schema string `json:"$schema"`
			Title  string `json:"title"`
		}
		if err := json.Unmarshal(schemaData, &schemaMeta); err != nil || schemaMeta.ID == "" || schemaMeta.Schema == "" {
			errorf("schema.v3.json: missing $id or $schema")
			ok = false
		} else {
			fmt.Printf("  schema.v3.json: valid JSON Schema definition\n")
		}
	}

	// Verify that the canonical contract families exist in the registry.
	expectedContracts := []string{
		"root.v3",
		"module.v3",
		"lesson.v3",
		"elective-lesson.v3",
		"project.v3",
		"assessment.v3",
		"flagship.v3",
		"flagship-checkpoint.v3",
		"companion-stage.v3",
	}
	for _, cName := range expectedContracts {
		c := reg.contractByID(cName)
		if c == nil {
			errorf("readme.contracts.json: missing contract '%s'", cName)
			ok = false
			continue
		}
		if len(c.RequiredHeadings) == 0 {
			errorf("readme.contracts.json: %s missing required_headings", cName)
			ok = false
		}
		if c.Selector.Doc == "" && len(c.Selector.Type) == 0 {
			errorf("readme.contracts.json: %s missing selector", cName)
			ok = false
		}
		if c.Shape == "" {
			errorf("readme.contracts.json: %s missing shape", cName)
			ok = false
		}
	}
	if c := reg.contractByID("lesson.v3"); c != nil && len(c.SubtypeModifiers) == 0 {
		errorf("readme.contracts.json: lesson.v3 missing subtype_modifiers")
		ok = false
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

	// Verify inventory counts against expected.json (single source of truth)
	if len(core.Items) != exp.Inventory.CoreItems {
		errorf("path.core.json: expected %d core items, got %d", exp.Inventory.CoreItems, len(core.Items))
		ok = false
	}
	if len(core.Modules) != exp.Inventory.CoreModules {
		errorf("path.core.json: expected %d core modules, got %d", exp.Inventory.CoreModules, len(core.Modules))
		ok = false
	}
	if len(electives.Items) != exp.Inventory.ElectiveItems {
		errorf("path.electives.json: expected %d elective items, got %d", exp.Inventory.ElectiveItems, len(electives.Items))
		ok = false
	}
	if len(electives.Modules) != exp.Inventory.ElectiveModules {
		errorf("path.electives.json: expected %d elective modules, got %d", exp.Inventory.ElectiveModules, len(electives.Modules))
		ok = false
	}
	if total := len(core.Items) + len(electives.Items); total != exp.Inventory.TotalItems {
		errorf("expected %d total items (core+electives), got %d", exp.Inventory.TotalItems, total)
		ok = false
	}
	if total := len(core.Modules) + len(electives.Modules); total != exp.Inventory.TotalModules {
		errorf("expected %d total modules (core+electives), got %d", exp.Inventory.TotalModules, total)
		ok = false
	}

	// The path manifests and the standalone registries are parallel views. Their
	// relationship is intentionally not one-to-one (a module assessment record
	// covers many lesson items, and some project records have no path item), but
	// the shape is pinned so it cannot drift silently.
	allPathItems := append(append([]Item{}, core.Items...), electives.Items...)
	pathProjectItems, pathAssessmentItems := 0, 0
	for _, item := range allPathItems {
		switch item.Type {
		case "project":
			pathProjectItems++
		case "assessment":
			pathAssessmentItems++
		}
	}
	if pathProjectItems != exp.Inventory.PathProjectItems {
		errorf("path manifests: expected %d project items, got %d", exp.Inventory.PathProjectItems, pathProjectItems)
		ok = false
	}
	if pathAssessmentItems != exp.Inventory.PathAssessmentItems {
		errorf("path manifests: expected %d assessment items, got %d", exp.Inventory.PathAssessmentItems, pathAssessmentItems)
		ok = false
	}
	if len(projects.Projects) != exp.Inventory.Projects {
		errorf("projects.json: expected %d projects, got %d", exp.Inventory.Projects, len(projects.Projects))
		ok = false
	}
	var conceptsBundle struct {
		Concepts []struct {
			Concept string `json:"concept"`
		} `json:"concepts"`
	}
	readJSON("concepts.json", &conceptsBundle)
	if len(conceptsBundle.Concepts) != exp.Inventory.Concepts {
		errorf("concepts.json: expected %d concepts, got %d", exp.Inventory.Concepts, len(conceptsBundle.Concepts))
		ok = false
	}
	for _, m := range core.Modules {
		if m.Status != "stable" {
			errorf("path.core.json module %s: expected status 'stable', got '%s'", m.ID, m.Status)
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
				errorf("%s: missing '%s'", item.ID, c.name)
				ok = false
			} else if rv.Kind() == reflect.String && rv.String() == "" {
				errorf("%s: '%s' is empty", item.ID, c.name)
				ok = false
			} else if rv.Kind() == reflect.Slice && rv.IsNil() {
				errorf("%s: '%s' is null (empty array expected)", item.ID, c.name)
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
						errorf("%s: zero_magic.%s is empty", item.ID, c.name)
						ok = false
					}
				} else if c.value == "" {
					errorf("%s: zero_magic.%s is empty", item.ID, c.name)
					ok = false
				}
			}
		}
	}

	// 5. Check proof fields. Every proof block is a complete contract: the task,
	// the assessment binding, the artifact the learner produces, the checklist
	// that proves mastery, and the rubric that scores it.
	for _, item := range allItems {
		if item.Proof == nil {
			continue
		}
		if item.Proof.PracticeTask == "" {
			errorf("%s: proof.practice_task is empty", item.ID)
			ok = false
		}
		if item.Proof.AssessmentID == "" {
			errorf("%s: proof.assessment_id is empty", item.ID)
			ok = false
		}
		if strings.TrimSpace(item.Proof.ExpectedArtifact) == "" {
			errorf("%s: proof.expected_artifact is empty", item.ID)
			ok = false
		}
		if len(item.Proof.MasteryChecks) == 0 {
			errorf("%s: proof.mastery_checks is empty", item.ID)
			ok = false
		}
		for i, check := range item.Proof.MasteryChecks {
			if strings.TrimSpace(check) == "" {
				errorf("%s: proof.mastery_checks[%d] is empty", item.ID, i)
				ok = false
			}
		}
	}

	// 5b. Check verification metadata. The mode is a closed vocabulary, and a
	// declared command must state what a successful run looks like.
	for _, item := range allItems {
		v := item.Verification
		if v == nil {
			errorf("%s: missing 'verification'", item.ID)
			ok = false
			continue
		}
		switch v.Mode {
		case "manual", "mixed", "automated":
		default:
			errorf("%s: verification.mode '%s' is not 'automated', 'manual', or 'mixed'", item.ID, v.Mode)
			ok = false
		}
		hasCommand := strings.TrimSpace(v.RunCommand) != "" || strings.TrimSpace(v.TestCommand) != ""
		if hasCommand && strings.TrimSpace(v.ExpectedOutput) == "" {
			errorf("%s: verification.expected_output is empty while a run/test command is declared", item.ID)
			ok = false
		}
		for i, step := range v.ManualSteps {
			if strings.TrimSpace(step) == "" {
				errorf("%s: verification.manual_steps[%d] is empty", item.ID, i)
				ok = false
			}
		}
	}

	// 5c. Check the files metadata that points at directories and diagrams. The
	// README/source/test paths are enforced by the physical-artifact gate; the
	// asset fields need their own existence check or a broken path is invisible.
	for _, item := range allItems {
		if item.Files == nil {
			continue
		}
		if assetsDir := strings.TrimSpace(item.Files.AssetsDir); assetsDir != "" {
			info, err := os.Stat(resolvePath(assetsDir))
			if err != nil || !info.IsDir() {
				errorf("%s: files.assets_dir '%s' does not exist or is not a directory", item.ID, assetsDir)
				ok = false
			}
		}
		for i, diagram := range item.Files.DiagramPaths {
			if strings.TrimSpace(diagram) == "" {
				errorf("%s: files.diagram_paths[%d] is empty", item.ID, i)
				ok = false
				continue
			}
			info, err := os.Stat(resolvePath(diagram))
			if err != nil || info.IsDir() {
				errorf("%s: files.diagram_paths[%d] '%s' does not exist or is not a file", item.ID, i, diagram)
				ok = false
			}
		}
	}

	// 5d. Check readme_contract metadata on every item and module. The contract
	// id resolves against the registry (aliases included), and the documentation
	// mode must be one of the declared authoring modes.
	for _, item := range allItems {
		if !validateReadmeContractFields(reg, item.ID, item.ReadmeContract) {
			ok = false
		}
	}
	for _, m := range core.Modules {
		if !validateReadmeContractFields(reg, m.ID, m.ReadmeContract) {
			ok = false
		}
	}
	for _, m := range electives.Modules {
		if !validateReadmeContractFields(reg, m.ID, m.ReadmeContract) {
			ok = false
		}
	}

	// 5e. portfolio_output is a promise that the module produces a project or
	// capstone artifact, so at least one such item must exist in the module.
	portfolioItems := make(map[string]bool)
	for _, item := range allItems {
		if item.Type == "project" || item.Type == "capstone" {
			portfolioItems[item.ModuleID] = true
		}
	}
	for _, m := range core.Modules {
		if m.PortfolioOutput && !portfolioItems[m.ID] {
			errorf("%s: portfolio_output is true but the module has no project or capstone item", m.ID)
			ok = false
		}
	}
	for _, m := range electives.Modules {
		if m.PortfolioOutput && !portfolioItems[m.ID] {
			errorf("%s: portfolio_output is true but the module has no project or capstone item", m.ID)
			ok = false
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
					errorf("%s: key '%s' out of canonical alphabetical order (after '%s')", ik.id, ik.keys[i], ik.keys[i-1])
					ok = false
					break
				}
			}
		}
	} else {
		errorf("cannot read path.core.json for key ordering check: %v", err)
		ok = false
	}

	// 7. Check no duplicate item IDs across core and electives
	seen := make(map[string]string)
	for _, item := range allItems {
		if prev, dup := seen[item.ID]; dup {
			errorf("duplicate item ID '%s' (in %s and %s)", item.ID, prev, item.SourceLegacyIDs)
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
			errorf("%s: references unknown module_id '%s'", item.ID, item.ModuleID)
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
	if len(assessments.Assessments) != exp.Inventory.Assessments {
		errorf("assessments.json: expected %d assessments, got %d", exp.Inventory.Assessments, len(assessments.Assessments))
		ok = false
	}
	assessIDs := make(map[string]bool)
	for _, a := range assessments.Assessments {
		assessIDs[a.ID] = true
	}

	for _, item := range allItems {
		if item.Proof != nil && item.Proof.AssessmentID != "" {
			if !assessIDs[item.Proof.AssessmentID] {
				errorf("%s: proof.assessment_id '%s' not found in assessments.json", item.ID, item.Proof.AssessmentID)
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
