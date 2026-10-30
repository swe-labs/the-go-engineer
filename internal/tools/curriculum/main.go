package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
		fmt.Println("Usage: go run ./internal/tools/curriculum <command>")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  validate-schema              Check JSON structure, formatting, required fields")
		fmt.Println("  validate-graph               Check prerequisites, cross-refs, module chains")
		fmt.Println("  validate-lessons             Check lesson contracts, zero_magic, content requirements")
		fmt.Println("  validate-prerequisite-closure Check no hidden leaks, circular deps, future-concept violations")
		fmt.Println("  validate-project-binding     Check every project has module binding, anchor, prerequisites")
		fmt.Println("  validate-no-placeholder-zm   Check no zero_magic fields contain placeholder text")
		fmt.Println("  validate-assessment-completeness Check assessments have module-specific criteria")
		fmt.Println("  validate-unique-tags         Check no duplicate tags across items")
		fmt.Println("  validate-crossrefs           Check crossref references resolve, no empty arrays on implemented lessons")
		fmt.Println("  validate-concept-ownership   Check concept ownership registry entries resolve, no duplicates")
		fmt.Println("  validate-cognitive-load      Check cognitive load metadata completeness and consistency")
		fmt.Println("  validate-operational-failures Check operational failure taxonomy completeness and consistency")
		fmt.Println("  validate-failure-engineering  Check critical items have operational failure examples")
		fmt.Println("  validate-zero-magic           Check golden lessons meet authoring spec")
		fmt.Println("  validate-physical-artifacts   Check all 441 items have verified physical files on disk")
		fmt.Println("  validate-assessment-artifacts Check all assessments have complete 4-file suites (README, questions, rubric, answer-key)")
		fmt.Println("  validate-starter-parity       Check project starters and solutions parse with valid Go ASTs")
		fmt.Println("  validate-all                 Run all validators")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "validate-schema":
		validateSchema()
	case "validate-graph":
		validateGraph()
	case "validate-lessons":
		validateLessons()
	case "validate-prerequisite-closure":
		validatePrerequisiteClosure()
	case "validate-project-binding":
		validateProjectBinding()
	case "validate-no-placeholder-zm":
		validateNoPlaceholderZM()
	case "validate-assessment-completeness":
		validateAssessmentCompleteness()
	case "validate-unique-tags":
		validateUniqueTags()
	case "validate-crossrefs":
		validateCrossrefs()
	case "validate-concept-ownership":
		validateConceptOwnership()
	case "validate-cognitive-load":
		validateCognitiveLoad()
	case "validate-operational-failures":
		validateOperationalFailures()
	case "validate-failure-engineering":
		validateFailureEngineering()
	case "validate-zero-magic":
		validateZeroMagic()
	case "validate-physical-artifacts":
		validatePhysicalArtifacts()
	case "validate-assessment-artifacts":
		validateAssessmentArtifacts()
	case "validate-starter-parity":
		validateStarterParity()
	case "validate-all":
		validateSchema()
		validateGraph()
		validateLessons()
		validatePrerequisiteClosure()
		validateProjectBinding()
		validateNoPlaceholderZM()
		validateAssessmentCompleteness()
		validateUniqueTags()
		validateCrossrefs()
		validateConceptOwnership()
		validateCognitiveLoad()
		validateOperationalFailures()
		validateFailureEngineering()
		validateZeroMagic()
		validatePhysicalArtifacts()
		validateAssessmentArtifacts()
		validateStarterParity()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}

	if errorCount > 0 {
		fmt.Fprintf(os.Stderr, "\nValidation failed with %d error(s).\n", errorCount)
		os.Exit(1)
	}
	fmt.Println("\nAll curriculum checks passed successfully.")
}

var errorCount int

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}

func warnf(format string, args ...any) {
	errorCount++
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
