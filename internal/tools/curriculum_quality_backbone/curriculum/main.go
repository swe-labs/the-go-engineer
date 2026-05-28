package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Severity string

const (
	Error Severity = "ERROR"
	Warn  Severity = "WARN"
	Info  Severity = "INFO"
)

type Finding struct {
	Severity Severity `json:"severity"`
	Check    string   `json:"check"`
	Entity   string   `json:"entity,omitempty"`
	Message  string   `json:"message"`
}

type Report struct {
	Profile          string         `json:"profile"`
	Root             string         `json:"root"`
	MetadataDir      string         `json:"metadata_dir"`
	StrictRepository bool           `json:"strict_repository"`
	Counts           map[string]int `json:"counts"`
	Findings         []Finding      `json:"findings"`
}

func NewReport(profile, root, metadataDir string, strict bool) *Report {
	return &Report{
		Profile:          profile,
		Root:             root,
		MetadataDir:      metadataDir,
		StrictRepository: strict,
		Counts:           map[string]int{},
	}
}

func (r *Report) add(sev Severity, check, entity, msg string) {
	r.Findings = append(r.Findings, Finding{Severity: sev, Check: check, Entity: entity, Message: msg})
}
func (r *Report) Error(check, entity, format string, args ...any) {
	r.add(Error, check, entity, fmt.Sprintf(format, args...))
}
func (r *Report) Warn(check, entity, format string, args ...any) {
	r.add(Warn, check, entity, fmt.Sprintf(format, args...))
}
func (r *Report) Info(check, entity, format string, args ...any) {
	r.add(Info, check, entity, fmt.Sprintf(format, args...))
}
func (r *Report) ErrorCount() int {
	n := 0
	for _, f := range r.Findings {
		if f.Severity == Error {
			n++
		}
	}
	return n
}
func (r *Report) WarnCount() int {
	n := 0
	for _, f := range r.Findings {
		if f.Severity == Warn {
			n++
		}
	}
	return n
}

func (r *Report) PrintText() {
	fmt.Printf("profile: %s\n", r.Profile)
	fmt.Printf("root: %s\n", r.Root)
	fmt.Printf("metadata_dir: %s\n", r.MetadataDir)
	fmt.Printf("strict_repository: %v\n", r.StrictRepository)
	keys := make([]string, 0, len(r.Counts))
	for k := range r.Counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, r.Counts[k])
	}
	fmt.Printf("errors: %d\n", r.ErrorCount())
	fmt.Printf("warnings: %d\n", r.WarnCount())
	for _, f := range r.Findings {
		if f.Entity != "" {
			fmt.Printf("%s [%s] %s: %s\n", f.Severity, f.Check, f.Entity, f.Message)
		} else {
			fmt.Printf("%s [%s]: %s\n", f.Severity, f.Check, f.Message)
		}
	}
	if r.ErrorCount() == 0 {
		fmt.Println("VALIDATION PASSED")
	} else {
		fmt.Println("VALIDATION FAILED")
	}
}

type Config struct {
	Root             string
	MetadataDir      string
	StrictRepository bool
	JSONOutput       bool
	FailOnWarnings   bool
}

type Curriculum struct {
	Files           map[string]map[string]any
	CoreItems       []map[string]any
	ElectiveItems   []map[string]any
	Items           []map[string]any
	CoreModules     []map[string]any
	ElectiveModules []map[string]any
	Modules         []map[string]any
	Projects        []map[string]any
	Assessments     []map[string]any
	Concepts        []map[string]any
	ItemByID        map[string]map[string]any
	ModuleByID      map[string]map[string]any
	ProjectByID     map[string]map[string]any
	AssessmentByID  map[string]map[string]any
	ConceptByName   map[string]map[string]any
}

var metadataFiles = []string{
	"path.core.json",
	"path.electives.json",
	"projects.json",
	"assessments.json",
	"crossrefs.json",
	"concepts.json",
	"failures.json",
	"readme.contracts.json",
	"migration.v2-to-v3.json",
	"workspace.json",
	"schema.v3.json",
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	rootFlag := fs.String("root", "", "repository root containing metadata/ or curriculum/")
	metadataFlag := fs.String("metadata-dir", "", "metadata directory; defaults to root/metadata, root/curriculum, or root")
	strictFlag := fs.Bool("strict-repository", false, "require every declared README/code/test/asset file to exist and pass content checks")
	jsonFlag := fs.Bool("json", false, "emit JSON report")
	failWarnFlag := fs.Bool("fail-on-warnings", false, "return non-zero when warnings exist")
	_ = fs.Parse(os.Args[2:])

	root, metadataDir, err := resolvePaths(*rootFlag, *metadataFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(2)
	}
	cfg := Config{Root: root, MetadataDir: metadataDir, StrictRepository: *strictFlag, JSONOutput: *jsonFlag, FailOnWarnings: *failWarnFlag}

	report := NewReport(cmd, cfg.Root, cfg.MetadataDir, cfg.StrictRepository)
	cur, loadErr := LoadCurriculum(cfg.MetadataDir, report)
	if loadErr != nil {
		report.Error("load", "metadata", loadErr.Error())
		finish(report, cfg)
		return
	}

	switch cmd {
	case "validate-metadata":
		runMetadataChecks(cur, report)
	case "validate-repository", "validate-lessons", "validate-content":
		ValidateRepositoryContent(cur, report, cfg)
	case "validate-all":
		runMetadataChecks(cur, report)
		ValidateRepositoryContent(cur, report, cfg)
	case "validate-schema":
		ValidateSchema(cur, report)
	case "validate-graph":
		ValidateGraph(cur, report)
	case "validate-crossrefs":
		ValidateCrossrefs(cur, report)
	case "validate-concept-ownership":
		ValidateConcepts(cur, report)
	case "validate-assessment-completeness", "validate-assessments":
		ValidateAssessments(cur, report)
	case "validate-project-binding", "validate-projects":
		ValidateProjects(cur, report)
	case "validate-zero-magic", "validate-no-placeholder-zm":
		ValidateZeroMagic(cur, report)
		ValidateNoPlaceholders(cur, report)
	case "validate-cognitive-load":
		ValidateCognitiveLoad(cur, report)
	case "validate-operational-failures", "validate-failure-engineering":
		ValidateFailures(cur, report)
	case "validate-readme-contracts":
		ValidateReadmeContracts(cur, report)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		usage()
		os.Exit(2)
	}
	finish(report, cfg)
}

func runMetadataChecks(cur *Curriculum, report *Report) {
	ValidateSchema(cur, report)
	ValidateGraph(cur, report)
	ValidateCrossrefs(cur, report)
	ValidateConcepts(cur, report)
	ValidateProjects(cur, report)
	ValidateAssessments(cur, report)
	ValidateZeroMagic(cur, report)
	ValidateNoPlaceholders(cur, report)
	ValidateCognitiveLoad(cur, report)
	ValidateFailures(cur, report)
	ValidateReadmeContracts(cur, report)
}

func finish(report *Report, cfg Config) {
	if cfg.JSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(report)
	} else {
		report.PrintText()
	}
	if report.ErrorCount() > 0 || (cfg.FailOnWarnings && report.WarnCount() > 0) {
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: go run ./internal/tools/curriculum <command> [--root .] [--metadata-dir metadata] [--strict-repository] [--json]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  validate-metadata              Validate JSON metadata only; should pass for the 100% metadata package")
	fmt.Println("  validate-repository            Validate actual README/code/test/asset files; use --strict-repository in CI")
	fmt.Println("  validate-all                   Validate metadata and repository content")
	fmt.Println("  validate-schema                Validate JSON shape and required fields")
	fmt.Println("  validate-graph                 Validate module/item graph and reachability")
	fmt.Println("  validate-crossrefs             Validate cross-reference targets and semantic reasons")
	fmt.Println("  validate-concept-ownership     Validate concept registry ownership and reinforcement")
	fmt.Println("  validate-projects              Validate project quality, bindings, rubrics")
	fmt.Println("  validate-assessments           Validate assessment targets, rubrics, evidence")
	fmt.Println("  validate-zero-magic            Validate zero-magic authored content")
	fmt.Println("  validate-readme-contracts      Validate README contract definitions")
}

func resolvePaths(rootArg, metadataArg string) (string, string, error) {
	wd, _ := os.Getwd()
	root := rootArg
	if root == "" {
		found, err := findRepoRoot(wd)
		if err != nil {
			return "", "", err
		}
		root = found
	}
	root, _ = filepath.Abs(root)
	metadataDir := metadataArg
	if metadataDir == "" {
		candidates := []string{filepath.Join(root, "metadata"), filepath.Join(root, "curriculum"), root}
		for _, c := range candidates {
			if fileExists(filepath.Join(c, "path.core.json")) {
				metadataDir = c
				break
			}
		}
	}
	if metadataDir == "" {
		return root, "", errors.New("could not find metadata directory containing path.core.json")
	}
	if !filepath.IsAbs(metadataDir) {
		metadataDir = filepath.Join(root, metadataDir)
	}
	metadataDir, _ = filepath.Abs(metadataDir)
	return root, metadataDir, nil
}

func findRepoRoot(start string) (string, error) {
	dir, _ := filepath.Abs(start)
	for {
		if fileExists(filepath.Join(dir, "metadata", "path.core.json")) || fileExists(filepath.Join(dir, "curriculum", "path.core.json")) || fileExists(filepath.Join(dir, "path.core.json")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("could not locate repository root; pass --root or --metadata-dir")
}

func LoadCurriculum(metadataDir string, report *Report) (*Curriculum, error) {
	cur := &Curriculum{
		Files:          map[string]map[string]any{},
		ItemByID:       map[string]map[string]any{},
		ModuleByID:     map[string]map[string]any{},
		ProjectByID:    map[string]map[string]any{},
		AssessmentByID: map[string]map[string]any{},
		ConceptByName:  map[string]map[string]any{},
	}
	for _, name := range metadataFiles {
		path := filepath.Join(metadataDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return cur, fmt.Errorf("missing metadata file %s: %w", name, err)
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			return cur, fmt.Errorf("invalid JSON in %s: %w", name, err)
		}
		cur.Files[name] = m
	}
	cur.CoreModules = objectList(cur.Files["path.core.json"], "modules")
	cur.ElectiveModules = objectList(cur.Files["path.electives.json"], "modules")
	cur.Modules = append(append([]map[string]any{}, cur.CoreModules...), cur.ElectiveModules...)
	cur.CoreItems = objectList(cur.Files["path.core.json"], "items")
	cur.ElectiveItems = objectList(cur.Files["path.electives.json"], "items")
	cur.Items = append(append([]map[string]any{}, cur.CoreItems...), cur.ElectiveItems...)
	cur.Projects = objectList(cur.Files["projects.json"], "projects")
	cur.Assessments = objectList(cur.Files["assessments.json"], "assessments")
	cur.Concepts = objectList(cur.Files["concepts.json"], "concepts")
	for _, item := range cur.Items {
		cur.ItemByID[str(item, "id")] = item
	}
	for _, mod := range cur.Modules {
		cur.ModuleByID[str(mod, "id")] = mod
	}
	for _, p := range cur.Projects {
		cur.ProjectByID[str(p, "id")] = p
	}
	for _, a := range cur.Assessments {
		cur.AssessmentByID[str(a, "id")] = a
	}
	for _, c := range cur.Concepts {
		cur.ConceptByName[str(c, "concept")] = c
	}
	report.Counts["metadata_files"] = len(metadataFiles)
	report.Counts["core_modules"] = len(cur.CoreModules)
	report.Counts["elective_modules"] = len(cur.ElectiveModules)
	report.Counts["core_items"] = len(cur.CoreItems)
	report.Counts["elective_items"] = len(cur.ElectiveItems)
	report.Counts["projects"] = len(cur.Projects)
	report.Counts["assessments"] = len(cur.Assessments)
	report.Counts["concepts"] = len(cur.Concepts)
	return cur, nil
}

func objectList(m map[string]any, key string) []map[string]any {
	raw, _ := m[key].([]any)
	out := make([]map[string]any, 0, len(raw))
	for _, v := range raw {
		if mm, ok := v.(map[string]any); ok {
			out = append(out, mm)
		}
	}
	return out
}

func str(m map[string]any, key string) string   { v, _ := m[key].(string); return strings.TrimSpace(v) }
func boolVal(m map[string]any, key string) bool { v, _ := m[key].(bool); return v }
func num(m map[string]any, key string) int {
	switch v := m[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}
func list(m map[string]any, key string) []any { v, _ := m[key].([]any); return v }
func stringsList(m map[string]any, key string) []string {
	raw := list(m, key)
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
func obj(m map[string]any, key string) map[string]any { v, _ := m[key].(map[string]any); return v }
func fileExists(path string) bool                     { info, err := os.Stat(path); return err == nil && !info.IsDir() }
func dirExists(path string) bool                      { info, err := os.Stat(path); return err == nil && info.IsDir() }
func validID(cur *Curriculum, id string) bool {
	if id == "" {
		return false
	}
	if _, ok := cur.ItemByID[id]; ok {
		return true
	}
	if _, ok := cur.ModuleByID[id]; ok {
		return true
	}
	if _, ok := cur.ProjectByID[id]; ok {
		return true
	}
	if _, ok := cur.ConceptByName[id]; ok {
		return true
	}
	return false
}
func stringIn(xs []string, target string) bool {
	for _, x := range xs {
		if x == target {
			return true
		}
	}
	return false
}
func isCoreItem(id string) bool {
	return strings.HasPrefix(id, "core-") || strings.HasPrefix(id, "opslane-")
}
func isElectiveItem(id string) bool { return strings.HasPrefix(id, "elective-") }

var genericReasonPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^both .+ build understanding within`),
	regexp.MustCompile(`(?i)^this project reinforces .+ through hands-on practice\.?$`),
	regexp.MustCompile(`(?i)requires understanding of\s*\.?$`),
	regexp.MustCompile(`(?i)^related concept:\s*builds on understanding from\s*\.?$`),
	regexp.MustCompile(`(?i)^this module builds on concepts from the previous module\.?$`),
}
var placeholderPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(todo|tbd|fixme|lorem ipsum|coming soon|placeholder|scaffolded)\b`),
	regexp.MustCompile(`(?i)this lesson explains what problem`),
	regexp.MustCompile(`(?i)one step in the learner'?s path`),
	regexp.MustCompile(`(?i)mechanically without understanding`),
	regexp.MustCompile(`(?i)in the context of professional go software engineering\.?$`),
}
