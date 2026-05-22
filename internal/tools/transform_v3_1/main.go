package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Module struct {
	ID              string   `json:"id"`
	Number          int      `json:"number"`
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Phase           string   `json:"phase"`
	Path            string   `json:"path"`
	Status          string   `json:"status"`
	LearningGoal    string   `json:"learning_goal"`
	Summary         string   `json:"summary"`
	Order           int      `json:"order"`
	Required        bool     `json:"required"`
	PortfolioOutput bool     `json:"portfolio_output"`
	Prerequisites   []string `json:"prerequisites"`
	EntryItemIDs    []string `json:"entry_item_ids"`
	TerminalItemIDs []string `json:"terminal_item_ids"`
	Tags            []string `json:"tags"`
	ReadmeStatus    string   `json:"readme_status"`
	ReadmeContract              any      `json:"readme_contract"`
	SourceLegacyIDs             []string `json:"source_legacy_section_ids"`
	CognitiveLoad               string   `json:"cognitive_load"`
	RecommendedBreakAfter       bool     `json:"recommended_break_after"`
	ContainsFoundationalHardConcepts bool `json:"contains_foundational_hard_concepts"`
	Pacing                      string   `json:"pacing"`
}

type Item struct {
	ID                   string         `json:"id"`
	ModuleID             string         `json:"module_id"`
	Slug                 string         `json:"slug"`
	Title                string         `json:"title"`
	Type                 string         `json:"type"`
	Subtype              string         `json:"subtype"`
	Status               string         `json:"status"`
	Difficulty           string         `json:"difficulty"`
	Phase                string         `json:"phase"`
	Order                int            `json:"order"`
	EstimatedMinutes     int            `json:"estimated_minutes"`
	LearningObjective    string         `json:"learning_objective"`
	RequiredPriorKnowledge []string     `json:"required_prior_knowledge"`
	Prerequisites        []string       `json:"prerequisites"`
	NextItemIDs          []string       `json:"next_item_ids"`
	ZeroMagic            any            `json:"zero_magic"`
	CrossRefs            any            `json:"crossrefs"`
	Proof                any            `json:"proof"`
	ContentContract      any            `json:"content_contract"`
	Verification         any            `json:"verification"`
	Files                any            `json:"files"`
	SourceLegacyIDs      []string       `json:"source_legacy_ids"`
	Tags                 []string       `json:"tags"`
	DocumentationMode    string         `json:"documentation_mode"`
	ReadmeStatus         string         `json:"readme_status"`
	ZeroMagicStatus      string         `json:"zero_magic_status"`
	ReadmeContract       any            `json:"readme_contract"`
}

type Bundle struct {
	SchemaVersion       string   `json:"schema_version"`
	DocumentType        string   `json:"document_type"`
	CurriculumVersion   string   `json:"curriculum_version"`
	LastUpdated         string   `json:"last_updated"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	ArchitecturalDecision any    `json:"architectural_decision"`
	RepositoryStructure []string `json:"repository_structure"`
	Modules             []Module `json:"modules"`
	Items               []Item   `json:"items"`
}

func main() {
	// Find curriculum directory
	wd, _ := os.Getwd()
	curriculumDir := filepath.Join(wd, "curriculum")
	if _, err := os.Stat(curriculumDir); err != nil {
		// Walk up
		for dir := wd; dir != "."; dir = filepath.Dir(dir) {
			candidate := filepath.Join(dir, "curriculum")
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				curriculumDir = candidate
				break
			}
			if dir == "/" || dir == `\` {
				break
			}
		}
	}

	path := filepath.Join(curriculumDir, "path.core.json")
	fmt.Printf("Reading %s\n", path)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
		os.Exit(1)
	}

	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		os.Exit(1)
	}

	// ====== TRANSFORMATIONS ======

	// 1. Update module-07 (CLI) - add runtime coordination
	for i := range bundle.Modules {
		m := &bundle.Modules[i]
		if m.ID == "module-07" {
			m.Title = "CLI, Files, JSON, Config, and Runtime Coordination"
			m.LearningGoal = "Build useful local programs and understand process lifecycle before servers."
			m.Summary = "Build useful local programs and understand process lifecycle before servers."
			if !contains(m.Tags, "runtime") {
				m.Tags = append(m.Tags, "runtime")
			}
			fmt.Println("  Updated module-07 title + tags")
		}
	}

	// 2. Remove module-08 (Context and Time Fundamentals)
	var keptModules []Module
	for _, m := range bundle.Modules {
		if m.ID != "module-08" {
			keptModules = append(keptModules, m)
		}
	}
	bundle.Modules = keptModules
	fmt.Println("  Removed module-08 (Context and Time Fundamentals)")

	// 3. Update shifted modules
	type shift struct {
		oldID         string
		newID         string
		number        int
		slug          string
		title         string
		pathv         string
		learningGoal  string
		summary       string
		tags          []string
		entry         []string
		terminal      []string
	}
	shifts := []shift{
		{"module-09", "module-08", 8, "http-rest-apis", "HTTP and REST APIs",
			"08-http-rest-apis", "Build the first serious backend service.", "Build the first serious backend service.",
			[]string{"backend"}, []string{"core-08-01"}, []string{"core-08-20"}},
		{"module-10", "module-09", 9, "sql-postgres-persistence", "SQL and PostgreSQL Persistence",
			"09-sql-postgres-persistence", "Teach real relational backend work.", "Teach real relational backend work.",
			[]string{"data"}, []string{"core-09-01"}, []string{"core-09-26"}},
		{"module-11", "module-10", 10, "auth-security", "Authentication, Authorization, and Security",
			"10-auth-security", "Make secure backend behavior explicit.", "Make secure backend behavior explicit.",
			[]string{"security"}, []string{"core-10-01"}, []string{"core-10-24"}},
		{"module-12", "module-11", 11, "lifecycle-context-concurrency", "Lifecycle, Context, and Concurrency",
			"11-lifecycle-context-concurrency", "Teach safe concurrent systems and request lifecycle after learners understand services.",
			"Teach safe concurrent systems and request lifecycle after learners understand services.",
			[]string{"concurrency", "lifecycle"}, []string{"core-11-01"}, []string{"core-11-28"}},
		{"module-13", "module-12", 12, "observability-diagnostics", "Observability and Diagnostics",
			"12-observability-diagnostics", "Make services observable before optimizing them.",
			"Make services observable before optimizing them.",
			[]string{"observability"}, []string{"core-12-01"}, []string{"core-12-20"}},
		{"module-14", "module-14", 14, "architecture-distributed-systems", "Architecture and Distributed Systems",
			"14-architecture-distributed-systems", "Teach architecture as tradeoff reasoning with distributed systems fundamentals.",
			"Teach architecture as tradeoff reasoning with distributed systems fundamentals.",
			[]string{"architecture", "distributed-systems"}, []string{"core-13-01"}, []string{"core-13-17"}},
		{"module-15", "module-15", 15, "docker-cicd-deployment", "Docker, CI/CD, and Deployment",
			"15-docker-cicd-deployment", "Make software shippable and repeatable.",
			"Make software shippable and repeatable.",
			[]string{"delivery"}, []string{"core-14-01"}, []string{"core-14-19"}},
		{"module-16", "module-16", 16, "portfolio-interview-readiness", "Portfolio and Interview Readiness",
			"16-portfolio-interview-readiness", "Convert competence into employability.",
			"Convert competence into employability.",
			[]string{"career"}, []string{"core-15-01"}, []string{"core-15-15"}},
	}

	// Build oldID→newID map for prerequisite updates
	moduleIDMap := make(map[string]string)
	for _, s := range shifts {
		moduleIDMap[s.oldID] = s.newID
	}

	for i := range bundle.Modules {
		m := &bundle.Modules[i]
		for _, s := range shifts {
			if m.ID == s.oldID {
				m.ID = s.newID
				m.Number = s.number
				m.Slug = s.slug
				m.Title = s.title
				m.Path = s.pathv
				m.LearningGoal = s.learningGoal
				m.Summary = s.summary
				m.Tags = s.tags
				m.EntryItemIDs = s.entry
				m.TerminalItemIDs = s.terminal
				break
			}
		}
	}

	// 4. Insert module-13 (Performance and Memory Engineering)
	perfModule := Module{
		ID:              "module-13",
		Number:          13,
		Slug:            "performance-memory-engineering",
		Title:           "Performance and Memory Engineering",
		Phase:           "performance",
		Path:            "13-performance-memory-engineering",
		Status:          "planned",
		LearningGoal:    "Make services efficient through profiling, benchmark-driven optimization, and memory understanding.",
		Summary:         "Make services efficient through profiling, benchmark-driven optimization, and memory understanding.",
		Order:           13,
		Required:        true,
		PortfolioOutput: true,
		Prerequisites:   []string{"module-12"},
		EntryItemIDs:    []string{"core-12-10"},
		TerminalItemIDs: []string{"core-12-17"},
		Tags:            []string{"performance"},
		ReadmeStatus:    "scaffolded",
		ReadmeContract: map[string]any{
			"contract_id":       "module.v3",
			"documentation_mode": "module-overview",
		},
	}
	insertAfter := -1
	for i, m := range bundle.Modules {
		if m.ID == "module-12" {
			insertAfter = i
			break
		}
	}
	if insertAfter >= 0 {
		bundle.Modules = append(bundle.Modules[:insertAfter+1],
			append([]Module{perfModule}, bundle.Modules[insertAfter+1:]...)...)
		fmt.Println("  Inserted module-13 (Performance and Memory Engineering)")
	}

	// 5. Update prerequisites for all modules
	prereqOverride := map[string][]string{
		"module-08": {"module-07"},
		"module-09": {"module-08"},
		"module-10": {"module-09"},
		"module-11": {"module-10"},
		"module-12": {"module-11"},
		"module-13": {"module-12"},
		"module-14": {"module-13"},
		"module-15": {"module-14"},
		"module-16": {"module-15"},
		"module-18": {"module-16"},
	}
	for i := range bundle.Modules {
		m := &bundle.Modules[i]
		if override, ok := prereqOverride[m.ID]; ok {
			m.Prerequisites = override
		}
	}

	// 6. Update item module_ids
	itemModuleMap := make(map[string]string)

	// HTTP: core-08-* → module-08 (these were in old module-09)
	for i := 1; i <= 20; i++ {
		itemModuleMap[fmt.Sprintf("core-08-%02d", i)] = "module-08"
	}
	// SQL: core-09-* → module-09 (were in old module-10)
	for i := 1; i <= 26; i++ {
		itemModuleMap[fmt.Sprintf("core-09-%02d", i)] = "module-09"
	}
	// Auth: core-10-* → module-10 (were in old module-11)
	for i := 1; i <= 24; i++ {
		itemModuleMap[fmt.Sprintf("core-10-%02d", i)] = "module-10"
	}
	// Context+Concurrency: core-11-* → module-11
	for i := 1; i <= 28; i++ {
		itemModuleMap[fmt.Sprintf("core-11-%02d", i)] = "module-11"
	}
	// Observability: core-12-01..09 + core-12-18..20 → module-12
	for i := 1; i <= 9; i++ {
		itemModuleMap[fmt.Sprintf("core-12-%02d", i)] = "module-12"
	}
	for _, i := range []int{18, 19, 20} {
		itemModuleMap[fmt.Sprintf("core-12-%02d", i)] = "module-12"
	}
	// Performance: core-12-10..17 → module-13
	for i := 10; i <= 17; i++ {
		itemModuleMap[fmt.Sprintf("core-12-%02d", i)] = "module-13"
	}

	changed := 0
	for i := range bundle.Items {
		item := &bundle.Items[i]
		if newMod, ok := itemModuleMap[item.ID]; ok && item.ModuleID != newMod {
			item.ModuleID = newMod
			changed++
		}
	}
	fmt.Printf("  Changed module_id for %d items\n", changed)

	// 7. Fix next_item_ids bridges
	for i := range bundle.Items {
		item := &bundle.Items[i]
		switch item.ID {
		case "core-07-21":
			item.NextItemIDs = []string{"core-08-01"}
		case "core-11-09":
			item.NextItemIDs = []string{"core-11-10"}
		case "core-11-28":
			item.NextItemIDs = []string{"core-12-01"}
		}
	}

	// 8. Write back
	bundle.SchemaVersion = "3.1.0"
	bundle.CurriculumVersion = "3.1.0-draft.1"
	bundle.LastUpdated = "2026-05-22"

	raw, err := json.MarshalIndent(bundle, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling: %v\n", err)
		os.Exit(1)
	}

	// Post-process: schema_version appears first due to field ordering
	if err := os.WriteFile(path, raw, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  Saved path.core.json (schema 3.1.0)")

	// === Verify ===
	var verify Bundle
	data2, _ := os.ReadFile(path)
	json.Unmarshal(data2, &verify)
	fmt.Println("\n=== Module Summary ===")
	for _, m := range verify.Modules {
		fmt.Printf("  %s (%d): %s | prereqs=%s | entry=%s | terminal=%s\n",
			m.ID, m.Number, m.Title, strings.Join(m.Prerequisites, ","),
			strings.Join(m.EntryItemIDs, ","), strings.Join(m.TerminalItemIDs, ","))
	}

	// Verify items by module
	fmt.Println("\n=== Items by Module ===")
	moduleItems := make(map[string][]string)
	for _, item := range verify.Items {
		moduleItems[item.ModuleID] = append(moduleItems[item.ModuleID], item.ID)
	}
	for _, m := range verify.Modules {
		if items, ok := moduleItems[m.ID]; ok {
			fmt.Printf("  %s (%d items): %s .. %s\n", m.ID, len(items), items[0], items[len(items)-1])
		}
	}

	fmt.Println("\nTransform complete. Run 'go run ./internal/tools/curriculum/ validate-graph' to verify.")
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
