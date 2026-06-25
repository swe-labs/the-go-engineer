package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Config struct {
	Root          string
	MetadataDir   string
	CurriculumDir string
	Strict        bool
}

type ValidationResult struct {
	Errors   []string
	Warnings []string
}

func (r *ValidationResult) Errorf(format string, args ...any) {
	r.Errors = append(r.Errors, fmt.Sprintf(format, args...))
}

func (r *ValidationResult) Warnf(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

func (r *ValidationResult) Merge(other ValidationResult) {
	r.Errors = append(r.Errors, other.Errors...)
	r.Warnings = append(r.Warnings, other.Warnings...)
}

func (r ValidationResult) OK() bool { return len(r.Errors) == 0 }

func (r ValidationResult) Print() {
	for _, warning := range r.Warnings {
		fmt.Printf("warning: %s\n", warning)
	}
	for _, err := range r.Errors {
		fmt.Printf("error: %s\n", err)
	}
	fmt.Printf("errors: %d\nwarnings: %d\n", len(r.Errors), len(r.Warnings))
}

type Module struct {
	ID                               string   `json:"id"`
	Number                           int      `json:"number"`
	Slug                             string   `json:"slug"`
	Title                            string   `json:"title"`
	Phase                            string   `json:"phase"`
	Path                             string   `json:"path"`
	Status                           string   `json:"status"`
	LearningGoal                     string   `json:"learning_goal"`
	Summary                          string   `json:"summary"`
	Order                            int      `json:"order"`
	Required                         bool     `json:"required"`
	PortfolioOutput                  bool     `json:"portfolio_output"`
	Prerequisites                    []string `json:"prerequisites"`
	EntryItemIDs                     []string `json:"entry_item_ids"`
	TerminalItemIDs                  []string `json:"terminal_item_ids"`
	SourceLegacySectionIDs           []string `json:"source_legacy_section_ids"`
	CognitiveLoad                    string   `json:"cognitive_load"`
	RecommendedBreakAfter            bool     `json:"recommended_break_after"`
	ContainsFoundationalHardConcepts bool     `json:"contains_foundational_hard_concepts"`
	Pacing                           string   `json:"pacing"`
	Tags                             []string `json:"tags"`
	ReadmeStatus                     string   `json:"readme_status"`
	ReadmeContract                   any      `json:"readme_contract"`
}

type FileRefs struct {
	ReadmePath   string   `json:"readme_path"`
	MainPath     string   `json:"main_path"`
	TestPath     string   `json:"test_path"`
	StarterPath  string   `json:"starter_path"`
	SolutionPath string   `json:"solution_path"`
	AssetsDir    string   `json:"assets_dir"`
	DiagramPaths []string `json:"diagram_paths"`
}

type Verification struct {
	RunCommand  string `json:"run_command"`
	TestCommand string `json:"test_command"`
	RaceCommand string `json:"race_command"`
}

type Proof struct {
	PracticeTask string `json:"practice_task"`
	AssessmentID string `json:"assessment_id"`
}

type ZeroMagic struct {
	ProblemSolved           string `json:"problem_solved"`
	WhyItExists             string `json:"why_it_exists"`
	MentalModel             string `json:"mental_model"`
	UnderTheHood            string `json:"under_the_hood"`
	HowGoUsesIt             string `json:"how_go_uses_it"`
	RealWorldUsage          string `json:"real_world_usage"`
	ProofOfUnderstanding    string `json:"proof_of_understanding"`
	BeginnerMistakes        any    `json:"beginner_mistakes"`
	ExecutionTimeline       any    `json:"execution_timeline"`
	PerformanceImplications any    `json:"performance_implications"`
	DebuggingChecklist      any    `json:"debugging_checklist"`
}

type CrossRefs struct {
	BuildsOn     []Reference `json:"builds_on"`
	PreviewOnly  []Reference `json:"preview_only"`
	ReinforcedIn []Reference `json:"reinforced_in"`
	Related      []Reference `json:"related"`
}

type Reference struct {
	ID       string `json:"id"`
	TargetID string `json:"target_id"`
	Reason   string `json:"reason"`
	Type     string `json:"type"`
	FromID   string `json:"from_id"`
	ToID     string `json:"to_id"`
}

type Item struct {
	ID                     string       `json:"id"`
	ModuleID               string       `json:"module_id"`
	Slug                   string       `json:"slug"`
	Title                  string       `json:"title"`
	Type                   string       `json:"type"`
	Subtype                string       `json:"subtype"`
	Status                 string       `json:"status"`
	Difficulty             string       `json:"difficulty"`
	Phase                  string       `json:"phase"`
	Order                  int          `json:"order"`
	EstimatedMinutes       int          `json:"estimated_minutes"`
	LearningObjective      string       `json:"learning_objective"`
	RequiredPriorKnowledge []string     `json:"required_prior_knowledge"`
	Prerequisites          []string     `json:"prerequisites"`
	NextItemIDs            []string     `json:"next_item_ids"`
	ZeroMagic              *ZeroMagic   `json:"zero_magic"`
	CrossRefs              CrossRefs    `json:"crossrefs"`
	Proof                  *Proof       `json:"proof"`
	ContentContract        any          `json:"content_contract"`
	Verification           Verification `json:"verification"`
	Files                  FileRefs     `json:"files"`
	SourceLegacyIDs        []string     `json:"source_legacy_ids"`
	Tags                   []string     `json:"tags"`
	DocumentationMode      string       `json:"documentation_mode"`
	ReadmeStatus           string       `json:"readme_status"`
	ZeroMagicStatus        string       `json:"zero_magic_status"`
	ReadmeContract         any          `json:"readme_contract"`
}

type PathBundle struct {
	SchemaVersion       string   `json:"schema_version"`
	DocumentType        string   `json:"document_type"`
	CurriculumVersion   string   `json:"curriculum_version"`
	LastUpdated         string   `json:"last_updated"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	RepositoryStructure []string `json:"repository_structure"`
	Modules             []Module `json:"modules"`
	Items               []Item   `json:"items"`
}

type Project struct {
	ID              string       `json:"id"`
	ModuleID        string       `json:"module_id"`
	Slug            string       `json:"slug"`
	Title           string       `json:"title"`
	Status          string       `json:"status"`
	TargetIDs       []string     `json:"target_ids"`
	Prerequisites   []string     `json:"prerequisites"`
	Reinforces      []string     `json:"reinforces"`
	AssessmentID    string       `json:"assessment_id"`
	Files           FileRefs     `json:"files"`
	Verification    Verification `json:"verification"`
	Rubric          any          `json:"rubric"`
	ReadmeStatus    string       `json:"readme_status"`
	SourceLegacyIDs []string     `json:"source_legacy_ids"`
}

type ProjectsBundle struct {
	SchemaVersion     string    `json:"schema_version"`
	DocumentType      string    `json:"document_type"`
	CurriculumVersion string    `json:"curriculum_version"`
	LastUpdated       string    `json:"last_updated"`
	Status            string    `json:"status"`
	Projects          []Project `json:"projects"`
}

type Assessment struct {
	ID           string            `json:"id"`
	ModuleID     string            `json:"module_id"`
	Title        string            `json:"title"`
	Type         string            `json:"type"`
	Status       string            `json:"status"`
	TargetIDs    []string          `json:"target_ids"`
	Files        map[string]string `json:"files"`
	Criteria     []map[string]any  `json:"criteria"`
	Rubric       map[string]any    `json:"rubric"`
	ReadmeStatus string            `json:"readme_status"`
}

type AssessmentsBundle struct {
	SchemaVersion     string       `json:"schema_version"`
	DocumentType      string       `json:"document_type"`
	CurriculumVersion string       `json:"curriculum_version"`
	LastUpdated       string       `json:"last_updated"`
	Status            string       `json:"status"`
	Assessments       []Assessment `json:"assessments"`
}

type Concept struct {
	Concept                string   `json:"concept"`
	CanonicalOwner         string   `json:"canonical_owner"`
	PreviewLocations       []string `json:"preview_locations"`
	ReinforcementLocations []string `json:"reinforcement_locations"`
}

type ConceptsBundle struct {
	SchemaVersion     string    `json:"schema_version"`
	DocumentType      string    `json:"document_type"`
	CurriculumVersion string    `json:"curriculum_version"`
	LastUpdated       string    `json:"last_updated"`
	Status            string    `json:"status"`
	Concepts          []Concept `json:"concepts"`
}

type CrossrefBundle struct {
	SchemaVersion     string `json:"schema_version"`
	DocumentType      string `json:"document_type"`
	CurriculumVersion string `json:"curriculum_version"`
	LastUpdated       string `json:"last_updated"`
	Status            string `json:"status"`
	Crossrefs         struct {
		References []Reference `json:"references"`
	} `json:"crossrefs"`
}

type FailureCategory struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ModuleIDs   []string `json:"module_ids"`
	Modules     []string `json:"modules"`
}

type FailuresBundle struct {
	SchemaVersion     string              `json:"schema_version"`
	DocumentType      string              `json:"document_type"`
	CurriculumVersion string              `json:"curriculum_version"`
	LastUpdated       string              `json:"last_updated"`
	Status            string              `json:"status"`
	FailureCategories []FailureCategory   `json:"failure_categories"`
	RequiredCoverage  map[string][]string `json:"required_coverage"`
}

type Metadata struct {
	Core         PathBundle
	Electives    PathBundle
	Projects     ProjectsBundle
	Assessments  AssessmentsBundle
	Concepts     ConceptsBundle
	Crossrefs    CrossrefBundle
	Failures     FailuresBundle
	RawContracts map[string]any
	RawMigration map[string]any
	RawWorkspace map[string]any
}

func discoverRoot(explicit string) (string, error) {
	if explicit != "" {
		return filepath.Abs(explicit)
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if fileExists(filepath.Join(dir, "metadata", "path.core.json")) {
			return dir, nil
		}
		if parent := filepath.Dir(dir); parent == dir {
			break
		}
	}
	return "", fmt.Errorf("could not find repository root containing metadata/path.core.json")
}

func loadMetadata(cfg Config) (Metadata, ValidationResult) {
	var m Metadata
	var r ValidationResult
	read := func(name string, target any) {
		path := filepath.Join(cfg.MetadataDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			r.Errorf("cannot read %s: %v", name, err)
			return
		}
		if err := json.Unmarshal(data, target); err != nil {
			r.Errorf("cannot parse %s: %v", name, err)
		}
	}
	read("path.core.json", &m.Core)
	read("path.electives.json", &m.Electives)
	read("projects.json", &m.Projects)
	read("assessments.json", &m.Assessments)
	read("concepts.json", &m.Concepts)
	read("crossrefs.json", &m.Crossrefs)
	read("failures.json", &m.Failures)
	read("readme.contracts.json", &m.RawContracts)
	read("migration.v2-to-v3.json", &m.RawMigration)
	read("workspace.json", &m.RawWorkspace)
	return m, r
}

func allItems(m Metadata) []Item {
	items := append([]Item{}, m.Core.Items...)
	items = append(items, m.Electives.Items...)
	return items
}

func allModules(m Metadata) []Module {
	modules := append([]Module{}, m.Core.Modules...)
	modules = append(modules, m.Electives.Modules...)
	return modules
}

func itemIDs(m Metadata) map[string]Item {
	result := map[string]Item{}
	for _, item := range allItems(m) {
		result[item.ID] = item
	}
	return result
}

func moduleIDs(m Metadata) map[string]Module {
	result := map[string]Module{}
	for _, module := range allModules(m) {
		result[module.ID] = module
	}
	return result
}

func projectIDs(m Metadata) map[string]Project {
	result := map[string]Project{}
	for _, project := range m.Projects.Projects {
		result[project.ID] = project
	}
	return result
}

func assessmentIDs(m Metadata) map[string]Assessment {
	result := map[string]Assessment{}
	for _, assessment := range m.Assessments.Assessments {
		result[assessment.ID] = assessment
	}
	return result
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isStableStatus(status string) bool {
	switch status {
	case "stable", "published", "ready":
		return true
	default:
		return false
	}
}

func canonicalModulePath(path string) bool {
	return strings.HasPrefix(path, "curriculum/modules/") || strings.HasPrefix(path, "curriculum/electives/")
}

func canonicalContentPath(path string) bool {
	if path == "" {
		return true
	}
	if !strings.HasPrefix(path, "curriculum/modules/") && !strings.HasPrefix(path, "curriculum/electives/") {
		return false
	}
	return strings.Contains(path, "/lessons/") || strings.Contains(path, "/labs/") || strings.Contains(path, "/projects/") || strings.Contains(path, "/assessments/")
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
