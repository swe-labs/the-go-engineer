package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var curriculumDir string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		curriculumDir = "curriculum"
		return
	}
	// Walk up to find the root curriculum/ directory containing path.core.json
	for dir := wd; ; {
		candidate := filepath.Join(dir, "curriculum")
		if info, err := os.Stat(filepath.Join(candidate, "path.core.json")); err == nil && !info.IsDir() {
			curriculumDir = candidate
			return
		}
		// If current directory is the curriculum directory itself
		if filepath.Base(dir) == "curriculum" {
			if info, err := os.Stat(filepath.Join(dir, "path.core.json")); err == nil && !info.IsDir() {
				curriculumDir = dir
				return
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir || dir == "." || dir == "/" || dir == `\` {
			break
		}
		dir = parent
	}
	// Fallback
	curriculumDir = filepath.Join(wd, "curriculum")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := ""
	for _, a := range os.Args[1:] {
		switch {
		case a == "--strict":
			strictMode = true
		case a == "--advisory":
			advisoryMode = true
		case a == "--report":
			reportPath = "curriculum-validation-report.json"
		case strings.HasPrefix(a, "--report="):
			reportPath = strings.TrimPrefix(a, "--report=")
		case strings.HasPrefix(a, "--"):
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n", a)
			os.Exit(2)
		default:
			if cmd == "" {
				cmd = a
			}
		}
	}
	if cmd == "" {
		fmt.Fprintln(os.Stderr, "missing command")
		os.Exit(1)
	}
	if advisoryMode && strictMode {
		fmt.Fprintln(os.Stderr, "--strict and --advisory are mutually exclusive")
		os.Exit(2)
	}
	ran := false
	for _, v := range validators() {
		if v.Name == cmd {
			v.Run()
			ran = true
			break
		}
	}
	if !ran {
		if cmd != "validate-all" {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			os.Exit(1)
		}
		for _, v := range validators() {
			v.Run()
		}
	}

	if errorCount > 0 {
		fmt.Fprintf(os.Stderr, "\nValidation failed with %d error(s).\n", errorCount)
		os.Exit(1)
	}
	if strictMode && warningCount > 0 {
		fmt.Fprintf(os.Stderr, "\nStrict release gate failed with %d warning(s).\n", warningCount)
		os.Exit(1)
	}
	if warningCount > 0 {
		fmt.Fprintf(os.Stderr, "\nAll curriculum checks passed (%d advisory warning(s); run with --strict to enforce).\n", warningCount)
		return
	}
	fmt.Println("\nAll curriculum checks passed successfully.")
}

// validator couples a CLI command name, its one-line description, and the
// function that implements it. It is the single source of truth for the
// command list, the usage text, and the `validate-all` sequence, so adding a
// gate is a one-line change instead of three parallel edits.
type validator struct {
	Name string
	Desc string
	Run  func()
}

// validators returns every curriculum validator in canonical execution order.
func validators() []validator {
	return []validator{
		{"validate-schema", "Check JSON structure, formatting, required fields", validateSchema},
		{"validate-graph", "Check prerequisites, cross-refs, module chains", validateGraph},
		{"validate-lessons", "Check lesson contracts, zero_magic, content requirements", validateLessons},
		{"validate-prerequisite-closure", "Check no hidden leaks, circular deps, future-concept violations", validatePrerequisiteClosure},
		{"validate-project-binding", "Check every project has module binding, anchor, prerequisites", validateProjectBinding},
		{"validate-no-placeholder-zm", "Check no zero_magic fields contain placeholder text", validateNoPlaceholderZM},
		{"validate-assessment-completeness", "Check assessments have module-specific criteria", validateAssessmentCompleteness},
		{"validate-unique-tags", "Check no duplicate tags across items", validateUniqueTags},
		{"validate-crossrefs", "Check crossref references resolve, no empty arrays on implemented lessons", validateCrossrefs},
		{"validate-concept-ownership", "Check concept ownership registry entries resolve, no duplicates", validateConceptOwnership},
		{"validate-cognitive-load", "Check cognitive load metadata completeness and consistency", validateCognitiveLoad},
		{"validate-operational-failures", "Check operational failure taxonomy completeness and consistency", validateOperationalFailures},
		{"validate-failure-engineering", "Check critical items have operational failure examples", validateFailureEngineering},
		{"validate-zero-magic", "Check golden lessons meet authoring spec", validateZeroMagic},
		{"validate-physical-artifacts", "Check all items have verified physical files on disk", validatePhysicalArtifacts},
		{"validate-assessment-artifacts", "Check all assessments have complete 4-file suites (README, questions, rubric, answer-key)", validateAssessmentArtifacts},
		{"validate-starter-parity", "Check project starters and solutions parse with valid Go ASTs", validateStarterParity},
		{"validate-markdown-links", "Check all relative markdown links resolve to real files on disk", validateMarkdownLinks},
		{"validate-sese-navigation", "Check Single-Entry Single-Exit navigation and exit gate integrity", validateSESENavigation},
		{"validate-project-scripts", "Check cross-platform bash/powershell project verification script parity", validateProjectScripts},
		{"validate-no-emojis", "Check strict No-Emojis Invariant across curriculum prose", validateNoEmojis},
		{"validate-anchor-quotes", "Check pioneer anchor quotes on locked module lessons", validateAnchorQuotes},
		{"validate-diagram-palette", "Check Dark Engineering Palette compliance on diagrams", validateDiagramPalette},
		{"validate-readme-contracts", "Check README contract families (headings, code policy, file policy) across all items", validateReadmeContracts},
	}
}

// printUsage renders the command list from the validator registry so the help
// text can never drift from the gates that actually run.
func printUsage() {
	fmt.Println("Usage: go run ./internal/tools/curriculum <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	width := len("validate-all")
	for _, v := range validators() {
		if len(v.Name) > width {
			width = len(v.Name)
		}
	}
	for _, v := range validators() {
		fmt.Printf("  %-*s %s\n", width, v.Name, v.Desc)
	}
	fmt.Printf("  %-*s %s\n", width, "validate-all", "Run all validators")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --strict                     Promote advisory warnings to errors (release gate)")
	fmt.Println("  --advisory                   Report structural errors only; skip content-shape checks")
	fmt.Println("  --report[=path]              Write a JSON violation report (default curriculum-validation-report.json)")
}

var errorCount int
var warningCount int

var (
	strictMode   bool
	advisoryMode bool
	reportPath   string
)

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}

// errorf reports a structural defect that always fails the build.
func errorf(format string, args ...any) {
	errorCount++
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

// warnf reports an advisory issue. It fails the build only under --strict.
func warnf(format string, args ...any) {
	warningCount++
	fmt.Fprintf(os.Stderr, "WARN: "+format+"\n", args...)
}

func infof(format string, args ...any) {
	fmt.Fprintf(os.Stdout, "INFO: "+format+"\n", args...)
}

// readJSON reads and parses a JSON file from the curriculum directory.
func readJSON(path string, v any) {
	fullPath := filepath.Join(curriculumDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		fatalf("cannot read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		fatalf("cannot parse %s: %v", path, err)
	}
}

// Module metadata structures
type Module struct {
	ID              string          `json:"id"`
	Number          int             `json:"number"`
	Slug            string          `json:"slug"`
	Title           string          `json:"title"`
	Phase           string          `json:"phase"`
	Path            string          `json:"path"`
	Status          string          `json:"status"`
	LearningGoal    string          `json:"learning_goal"`
	Summary         string          `json:"summary"`
	Order           int             `json:"order"`
	Required        bool            `json:"required"`
	PortfolioOutput bool            `json:"portfolio_output"`
	Prerequisites   []string        `json:"prerequisites"`
	EntryItemIDs    []string        `json:"entry_item_ids"`
	TerminalItemIDs []string        `json:"terminal_item_ids"`
	Tags            []string        `json:"tags"`
	ReadmeContract  *ReadmeContract `json:"readme_contract"`
}

type ZeroMagic struct {
	ProblemSolved                string `json:"problem_solved"`
	WhyItExists                  string `json:"why_it_exists"`
	MentalModel                  string `json:"mental_model"`
	UnderTheHood                 string `json:"under_the_hood"`
	HowGoUsesIt                  string `json:"how_go_uses_it"`
	RealWorldUsage               string `json:"real_world_usage"`
	ProofOfUnderstanding         string `json:"proof_of_understanding"`
	BeginnerMistakes             []any  `json:"beginner_mistakes"`
	StepByStepExecution          []any  `json:"step_by_step_execution,omitempty"`
	ExecutionTimeline            []any  `json:"execution_timeline,omitempty"`
	MemoryTimeline               []any  `json:"memory_timeline,omitempty"`
	FailureModes                 []any  `json:"failure_modes,omitempty"`
	OperationalFailureExamples   []any  `json:"operational_failure_examples,omitempty"`
	HiddenMagicChecks            []any  `json:"hidden_magic_checks,omitempty"`
	DebuggingWalkthroughs        []any  `json:"debugging_walkthroughs,omitempty"`
	DebuggingWalkthroughRequired bool   `json:"debugging_walkthrough_required,omitempty"`
	ProductionExamples           []any  `json:"production_examples,omitempty"`
	PerformanceImplications      []any  `json:"performance_implications,omitempty"`
	OperationalConsiderations    []any  `json:"operational_considerations,omitempty"`
	CodeReadingTasks             []any  `json:"code_reading_tasks,omitempty"`
	RefactoringTasks             []any  `json:"refactoring_tasks,omitempty"`
	ReviewQuestions              []any  `json:"review_questions,omitempty"`
}

type CrossrefRef struct {
	TargetID string `json:"target_id"`
	Label    string `json:"label"`
	Reason   string `json:"reason"`
	Required bool   `json:"required"`
}

type ItemCrossRefs struct {
	BuildsOn     []CrossrefRef `json:"builds_on"`
	PreviewOnly  []CrossrefRef `json:"preview_only"`
	Related      []CrossrefRef `json:"related"`
	ReinforcedIn []CrossrefRef `json:"reinforced_in"`
}

type Proof struct {
	PracticeTask     string   `json:"practice_task"`
	AssessmentID     string   `json:"assessment_id"`
	ProjectID        *string  `json:"project_id"`
	ExpectedArtifact string   `json:"expected_artifact"`
	MasteryChecks    []string `json:"mastery_checks"`
	RubricIDs        []string `json:"rubric_ids"`
}

type ReadmeContract struct {
	ContractID        string `json:"contract_id"`
	DocumentationMode string `json:"documentation_mode"`
	Required          bool   `json:"required"`
}

type ContentContract struct {
	ReadmeRequired            bool `json:"readme_required"`
	RunnableRequired          bool `json:"runnable_required"`
	TestsRequired             bool `json:"tests_required"`
	VisualModelRequired       bool `json:"visual_model_required"`
	MachineViewRequired       bool `json:"machine_view_required"`
	CommonMistakesRequired    bool `json:"common_mistakes_required"`
	ProductionNotesRequired   bool `json:"production_notes_required"`
	ReviewQuestionsRequired   bool `json:"review_questions_required"`
	PortfolioArtifactRequired bool `json:"portfolio_artifact_required"`
}

type Verification struct {
	Mode           string   `json:"mode"`
	RunCommand     string   `json:"run_command"`
	TestCommand    string   `json:"test_command"`
	RaceCommand    string   `json:"race_command"`
	ManualSteps    []string `json:"manual_steps"`
	ExpectedOutput string   `json:"expected_output"`
}

type Files struct {
	ReadmePath   string   `json:"readme_path"`
	MainPath     string   `json:"main_path"`
	StarterPath  string   `json:"starter_path"`
	SolutionPath string   `json:"solution_path"`
	TestPath     string   `json:"test_path"`
	AssetsDir    string   `json:"assets_dir"`
	DiagramPaths []string `json:"diagram_paths"`
}

type Item struct {
	ID                     string           `json:"id"`
	ModuleID               string           `json:"module_id"`
	Slug                   string           `json:"slug"`
	Title                  string           `json:"title"`
	Type                   string           `json:"type"`
	Subtype                string           `json:"subtype"`
	Status                 string           `json:"status"`
	Difficulty             string           `json:"difficulty"`
	Phase                  string           `json:"phase"`
	Order                  int              `json:"order"`
	EstimatedMinutes       int              `json:"estimated_minutes"`
	LearningObjective      string           `json:"learning_objective"`
	RequiredPriorKnowledge []string         `json:"required_prior_knowledge"`
	Prerequisites          []string         `json:"prerequisites"`
	NextItemIDs            []string         `json:"next_item_ids"`
	ZeroMagic              *ZeroMagic       `json:"zero_magic"`
	CrossRefs              *ItemCrossRefs   `json:"crossrefs"`
	Proof                  *Proof           `json:"proof"`
	ContentContract        *ContentContract `json:"content_contract"`
	Verification           *Verification    `json:"verification"`
	Files                  *Files           `json:"files"`
	ReadmeContract         *ReadmeContract  `json:"readme_contract"`
	ReadmeStatus           string           `json:"readme_status"`
	ZeroMagicStatus        string           `json:"zero_magic_status"`
	SourceLegacyIDs        []string         `json:"source_legacy_ids"`
	Tags                   []string         `json:"tags"`
}

type CoreBundle struct {
	Modules []Module `json:"modules"`
	Items   []Item   `json:"items"`
}

type ElectiveBundle struct {
	Modules []Module `json:"modules"`
	Items   []Item   `json:"items"`
}

type Project struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Order int    `json:"order"`
}

type ProjectsBundle struct {
	Projects []Project `json:"projects"`
}

type Assessment struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Order int    `json:"order"`
}

type AssessmentsBundle struct {
	Assessments []Assessment `json:"assessments"`
}

type CrossrefEntry struct {
	FromID   string `json:"from_id"`
	TargetID string `json:"target_id"`
	Relation string `json:"relation"`
	Label    string `json:"label,omitempty"`
}

type CrossrefBundle struct {
	Crossrefs struct {
		References []CrossrefEntry `json:"references"`
	} `json:"crossrefs"`
}
