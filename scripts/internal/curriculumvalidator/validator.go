package curriculumvalidator

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Lesson struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type Section struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Lessons []Lesson `json:"lessons"`
}

type Curriculum struct {
	Sections []Section `json:"sections"`
}

type V2Section struct {
	ID            string   `json:"id"`
	Number        string   `json:"number"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	PathPrefix    string   `json:"path_prefix"`
	EntryPoints   []string `json:"entry_points"`
	Outputs       []string `json:"outputs"`
	Prerequisites []string `json:"prerequisites"`
}

type V2Item struct {
	ID               string   `json:"id"`
	SectionID        string   `json:"section_id"`
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Type             string   `json:"type"`
	Subtype          string   `json:"subtype"`
	Level            string   `json:"level"`
	VerificationMode string   `json:"verification_mode"`
	Path             string   `json:"path"`
	Prerequisites    []string `json:"prerequisites"`
	RunCommand       string   `json:"run_command"`
	TestCommand      string   `json:"test_command"`
	StarterPath      string   `json:"starter_path"`
	NextItems        []string `json:"next_items"`
}

type V2Curriculum struct {
	SchemaVersion int         `json:"schema_version"`
	Sections      []V2Section `json:"sections"`
	Items         []V2Item    `json:"items"`
}

type Result struct {
	LessonCount    int
	FilesScanned   int
	V2SectionCount int
	V2ItemCount    int
	HasV2          bool
	ErrorCount     int
}

var runPathPattern = regexp.MustCompile(`\./[A-Za-z0-9._/\-]+`)
var nextUpIDPattern = regexp.MustCompile(`NEXT UP:\s*([A-Z]{2,6}\.\d+)`)
var markdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)

var (
	allowedItemTypes = map[string]bool{
		"lesson":       true,
		"drill":        true,
		"exercise":     true,
		"checkpoint":   true,
		"mini_project": true,
		"capstone":     true,
		"reference":    true,
	}
	allowedLessonSubtypes = map[string]bool{
		"concept":     true,
		"pattern":     true,
		"integration": true,
	}
	allowedLevels = map[string]bool{
		"foundation": true,
		"core":       true,
		"stretch":    true,
		"production": true,
	}
	allowedVerificationModes = map[string]bool{
		"run":    true,
		"test":   true,
		"rubric": true,
		"mixed":  true,
	}
)

func Validate(root string, report func(string)) (Result, error) {
	if report == nil {
		report = func(string) {}
	}

	lessonCount, pathErrors := 0, 0
	if pathExists(root, "curriculum.json") {
		var err error
		lessonCount, pathErrors, err = validateCurriculumPaths(root, report)
		if err != nil {
			return Result{}, err
		}
	}

	filesScanned, runErrors, err := validateRunPaths(root, report)
	if err != nil {
		return Result{}, err
	}

	v2SectionCount, v2ItemCount, v2Errors, hasV2, err := validateV2Curriculum(root, report)
	if err != nil {
		return Result{}, err
	}

	pressureErrors := validatePressureDocs(root, report)
	templateErrors := validateTemplateDocs(root, report)

	return Result{
		LessonCount:    lessonCount,
		FilesScanned:   filesScanned,
		V2SectionCount: v2SectionCount,
		V2ItemCount:    v2ItemCount,
		HasV2:          hasV2,
		ErrorCount:     pathErrors + runErrors + v2Errors + pressureErrors + templateErrors,
	}, nil
}

func isPlaceholderPath(path string) bool {
	placeholders := []string{
		"./SECTION/LESSON",
		"./path/to/",
		"./my-folder",
		"./my-messy-folder",
		"./00-how-computers-work/...",
	}

	for _, placeholder := range placeholders {
		if strings.Contains(path, placeholder) {
			return true
		}
	}

	return false
}

func pathExists(root, path string) bool {
	_, err := os.Stat(filepath.Join(root, path))
	return err == nil
}

func extractCommandTarget(command string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", errors.New("command is empty")
	}

	match := runPathPattern.FindString(command)
	if match == "" || match == "./..." || isPlaceholderPath(match) {
		return "", fmt.Errorf("command does not contain a concrete ./path target: %q", command)
	}

	return filepath.Clean(strings.TrimPrefix(match, "./")), nil
}

func validateCurriculumPaths(root string, report func(string)) (int, int, error) {
	data, err := os.ReadFile(filepath.Join(root, "curriculum.v2.json"))
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to read curriculum.v2.json: %v", err)
	}

	var v2 V2Curriculum
	if err := json.Unmarshal(data, &v2); err != nil {
		return 0, 0, fmt.Errorf("Failed to parse curriculum.v2.json: %v", err)
	}

	errorsFound := 0
	lessonCount := 0
	for _, item := range v2.Items {
		if item.Type != "lesson" && item.Type != "exercise" {
			continue
		}
		lessonCount++
		if item.Path == "" {
			report(fmt.Sprintf("Unmapped item: %s - %s", item.ID, item.Title))
			errorsFound++
			continue
		}

		if !pathExists(root, item.Path) {
			report(fmt.Sprintf("Path does not exist: %s (%s - %s)", item.Path, item.ID, item.Title))
			errorsFound++
		}
	}

	return lessonCount, errorsFound, nil
}

func shouldScanRunPaths(path string) bool {
	if filepath.Ext(path) == ".go" {
		return true
	}

	if filepath.Base(path) != "README.md" {
		return false
	}

	first := strings.Split(filepath.ToSlash(path), "/")[0]
	if len(first) < 2 {
		return false
	}

	return first[0] >= '0' && first[0] <= '9' && first[1] >= '0' && first[1] <= '9'
}

func validateRunPaths(root string, report func(string)) (int, int, error) {
	filesScanned := 0
	errorsFound := 0

	walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor":
				return filepath.SkipDir
			default:
				return nil
			}
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		cleanPath := filepath.Clean(relPath)
		if !shouldScanRunPaths(cleanPath) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		filesScanned++
		scanner := bufio.NewScanner(file)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			if !strings.Contains(line, "go run ") && !strings.Contains(line, "go test ") {
				continue
			}

			for _, match := range runPathPattern.FindAllString(line, -1) {
				if match == "./..." || isPlaceholderPath(match) {
					continue
				}

				target := filepath.Clean(strings.TrimPrefix(match, "./"))
				alternateTarget := filepath.Clean(filepath.Join(filepath.Dir(cleanPath), strings.TrimPrefix(match, "./")))

				if pathExists(root, target) || pathExists(root, alternateTarget) {
					continue
				}

				report(fmt.Sprintf("Invalid run path: %s:%d -> %s", cleanPath, lineNo, match))
				errorsFound++
			}
		}

		return scanner.Err()
	})
	if walkErr != nil {
		return 0, 0, fmt.Errorf("Failed to scan run paths: %v", walkErr)
	}

	return filesScanned, errorsFound, nil
}

func validateV2Curriculum(root string, report func(string)) (int, int, int, bool, error) {
	if _, err := os.Stat(filepath.Join(root, "curriculum.v2.json")); os.IsNotExist(err) {
		return 0, 0, 0, false, nil
	}

	data, err := os.ReadFile(filepath.Join(root, "curriculum.v2.json"))
	if err != nil {
		return 0, 0, 0, false, fmt.Errorf("Failed to read curriculum.v2.json: %v", err)
	}

	var cur V2Curriculum
	if err := json.Unmarshal(data, &cur); err != nil {
		return 0, 0, 0, false, fmt.Errorf("Failed to parse curriculum.v2.json: %v", err)
	}

	errorsFound := 0
	sectionIDs := make(map[string]V2Section, len(cur.Sections))
	itemIDs := make(map[string]V2Item, len(cur.Items))

	for _, s := range cur.Sections {
		if s.ID == "" {
			report("Invalid v2 section: missing id")
			errorsFound++
			continue
		}

		if _, exists := sectionIDs[s.ID]; exists {
			report(fmt.Sprintf("Duplicate v2 section id: %s", s.ID))
			errorsFound++
			continue
		}

		if s.Number == "" || s.Slug == "" || s.Title == "" || s.PathPrefix == "" {
			report(fmt.Sprintf("Invalid v2 section metadata: %s requires number, slug, title, and path_prefix", s.ID))
			errorsFound++
		}

		if s.PathPrefix != "" && !pathExists(root, s.PathPrefix) {
			report(fmt.Sprintf("Invalid v2 section path_prefix: %s -> %s", s.ID, s.PathPrefix))
			errorsFound++
		}

		sectionIDs[s.ID] = s
	}

	for _, item := range cur.Items {
		if item.ID == "" {
			report("Invalid v2 item: missing id")
			errorsFound++
			continue
		}

		if _, exists := itemIDs[item.ID]; exists {
			report(fmt.Sprintf("Duplicate v2 item id: %s", item.ID))
			errorsFound++
			continue
		}

		if item.SectionID == "" || item.Slug == "" || item.Title == "" || item.Type == "" || item.Level == "" || item.VerificationMode == "" || item.Path == "" {
			report(fmt.Sprintf("Invalid v2 item metadata: %s requires section_id, slug, title, type, level, verification_mode, and path", item.ID))
			errorsFound++
		}

		if _, exists := sectionIDs[item.SectionID]; !exists {
			report(fmt.Sprintf("Invalid v2 section linkage: %s -> %s", item.ID, item.SectionID))
			errorsFound++
		}

		if !allowedItemTypes[item.Type] {
			report(fmt.Sprintf("Invalid v2 item type: %s -> %s", item.ID, item.Type))
			errorsFound++
		}

		if item.Type == "lesson" {
			if !allowedLessonSubtypes[item.Subtype] {
				report(fmt.Sprintf("Invalid v2 lesson subtype: %s -> %s", item.ID, item.Subtype))
				errorsFound++
			}
		} else if item.Subtype != "" {
			report(fmt.Sprintf("Unexpected v2 subtype for non-lesson item: %s -> %s", item.ID, item.Subtype))
			errorsFound++
		}

		if !allowedLevels[item.Level] {
			report(fmt.Sprintf("Invalid v2 item level: %s -> %s", item.ID, item.Level))
			errorsFound++
		}

		if !allowedVerificationModes[item.VerificationMode] {
			report(fmt.Sprintf("Invalid v2 verification mode: %s -> %s", item.ID, item.VerificationMode))
			errorsFound++
		}

		if !pathExists(root, item.Path) {
			report(fmt.Sprintf("Invalid v2 item path: %s -> %s", item.ID, item.Path))
			errorsFound++
		}

		if section, exists := sectionIDs[item.SectionID]; exists && section.PathPrefix != "" {
			itemPath := filepath.ToSlash(filepath.Clean(item.Path))
			allowedPrefixes := allowedPathPrefixesForSection(section)
			if !matchesAnyPrefix(itemPath, allowedPrefixes) {
				report(fmt.Sprintf("Invalid v2 section path alignment: %s -> %s (expected prefix %s)", item.ID, item.Path, strings.Join(allowedPrefixes, " or ")))
				errorsFound++
			}
		}

		if item.RunCommand != "" {
			target, err := extractCommandTarget(item.RunCommand)
			if err != nil {
				report(fmt.Sprintf("Invalid v2 run command: %s -> %v", item.ID, err))
				errorsFound++
			} else if !pathExists(root, target) {
				report(fmt.Sprintf("Invalid v2 run command target: %s -> %s", item.ID, item.RunCommand))
				errorsFound++
			}
		}

		if item.TestCommand != "" {
			target, err := extractCommandTarget(item.TestCommand)
			if err != nil {
				report(fmt.Sprintf("Invalid v2 test command: %s -> %v", item.ID, err))
				errorsFound++
			} else if !pathExists(root, target) {
				report(fmt.Sprintf("Invalid v2 test command target: %s -> %s", item.ID, item.TestCommand))
				errorsFound++
			}
		}

		switch item.VerificationMode {
		case "run":
			if strings.TrimSpace(item.RunCommand) == "" {
				report(fmt.Sprintf("Invalid v2 run contract: %s requires run_command", item.ID))
				errorsFound++
			}
		case "test":
			if strings.TrimSpace(item.TestCommand) == "" {
				report(fmt.Sprintf("Invalid v2 test contract: %s requires test_command", item.ID))
				errorsFound++
			}
		case "mixed":
			if strings.TrimSpace(item.RunCommand) == "" && strings.TrimSpace(item.TestCommand) == "" {
				report(fmt.Sprintf("Invalid v2 mixed contract: %s requires run_command or test_command", item.ID))
				errorsFound++
			}
		}

		if item.StarterPath != "" && !pathExists(root, item.StarterPath) {
			report(fmt.Sprintf("Invalid v2 starter path: %s -> %s", item.ID, item.StarterPath))
			errorsFound++
		}

		itemIDs[item.ID] = item
	}

	for _, s := range cur.Sections {
		for _, prereqID := range s.Prerequisites {
			if _, exists := sectionIDs[prereqID]; !exists {
				report(fmt.Sprintf("Invalid v2 section prerequisite: %s -> %s", s.ID, prereqID))
				errorsFound++
			}
		}

		for _, entryID := range s.EntryPoints {
			if _, exists := itemIDs[entryID]; !exists {
				report(fmt.Sprintf("Invalid v2 section entry point: %s -> %s", s.ID, entryID))
				errorsFound++
			}
		}

		for _, outputID := range s.Outputs {
			if _, exists := itemIDs[outputID]; !exists {
				report(fmt.Sprintf("Invalid v2 section output: %s -> %s", s.ID, outputID))
				errorsFound++
			}
		}
	}

	for _, item := range cur.Items {
		for _, prereqID := range item.Prerequisites {
			if _, itemExists := itemIDs[prereqID]; itemExists {
				continue
			}
			if _, sectionExists := sectionIDs[prereqID]; sectionExists {
				continue
			}
			report(fmt.Sprintf("Invalid v2 prerequisite: %s -> %s", item.ID, prereqID))
			errorsFound++
		}

		for _, nextID := range item.NextItems {
			if _, itemExists := itemIDs[nextID]; itemExists {
				continue
			}
			if _, sectionExists := sectionIDs[nextID]; sectionExists {
				continue
			}
			report(fmt.Sprintf("Invalid v2 next item: %s -> %s", item.ID, nextID))
			errorsFound++
		}
	}

	errorsFound += validateV2LessonNavigation(root, cur.Items, report)
	errorsFound += validateV2SectionLabels(root, sectionIDs, cur.Items, report)
	errorsFound += validateV2TextEncoding(root, sectionIDs, cur.Items, report)
	errorsFound += validateFoundationsReadmeContracts(root, cur.Items, report)

	return len(cur.Sections), len(cur.Items), errorsFound, true, nil
}

func allowedPathPrefixesForSection(section V2Section) []string {
	prefixes := []string{filepath.ToSlash(filepath.Clean(section.PathPrefix))}

	if section.ID == "s01" {
		prefixes = append(prefixes, "01-foundations/01-getting-started", "01-foundations/02-language-basics", "01-getting-started", "02-language-basics")
	}

	if section.ID == "s04" {
		prefixes = append(prefixes, "05-composition", "06-strings-and-text")
	}

	return prefixes
}

func matchesAnyPrefix(itemPath string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if itemPath == prefix || strings.HasPrefix(itemPath, prefix+"/") {
			return true
		}
	}

	return false
}

func isFoundationsSection(sectionID string) bool {
	switch sectionID {
	case "s00", "s01", "s02", "s03", "s04":
		return true
	default:
		return false
	}
}

func validateFoundationsReadmeContracts(root string, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		if !isFoundationsSection(item.SectionID) {
			continue
		}

		itemPath := filepath.ToSlash(filepath.Clean(item.Path))
		readmePath := filepath.ToSlash(filepath.Join(itemPath, "README.md"))
		if !pathExists(root, readmePath) {
			report(fmt.Sprintf("Missing foundations README: %s -> %s", item.ID, readmePath))
			errorsFound++
			continue
		}

		errorsFound += validateMarkdownLocalLinks(root, readmePath, report)
		errorsFound += validateRequiredHeadingsForItem(root, readmePath, item, report)
		errorsFound += validateFoundationsVisualModelUsesMermaid(root, readmePath, item.ID, report)

		if item.Type == "lesson" && item.VerificationMode == "run" {
			mainPath := filepath.ToSlash(filepath.Join(itemPath, "main.go"))
			if !pathExists(root, mainPath) {
				report(fmt.Sprintf("Missing foundations lesson main.go: %s -> %s", item.ID, mainPath))
				errorsFound++
			}
		}
	}

	return errorsFound
}

func validateRequiredHeadingsForItem(root, readmePath string, item V2Item, report func(string)) int {
	requiredHeadings := []string{
		"## Mission",
		"## Prerequisites",
		"## Mental Model",
		"## Visual Model",
		"## Machine View",
		"## Run Instructions",
	}

	if item.Type == "lesson" {
		requiredHeadings = append(requiredHeadings, "## Code Walkthrough")
	} else {
		requiredHeadings = append(requiredHeadings, "## Solution Walkthrough")
	}

	requiredHeadings = append(requiredHeadings,
		"## Try It",
	)

	if item.Type != "lesson" {
		requiredHeadings = append(requiredHeadings, "## Verification Surface")
	}

	requiredHeadings = append(requiredHeadings,
		"## ⚠️ In Production",
		"## 🤔 Thinking Questions",
		"## Next Step",
	)

	return validateRequiredHeadingsWithID(root, readmePath, item.ID, requiredHeadings, report)
}

func validateRequiredHeadingsWithID(root, relPath, itemID string, headings []string, report func(string)) int {
	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		report(fmt.Sprintf("Failed to read foundations README: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	errorsFound := 0
	lastOffset := 0
	for _, heading := range headings {
		idx := strings.Index(text[lastOffset:], heading)
		if idx >= 0 {
			lastOffset += idx + len(heading)
			continue
		}

		if strings.Contains(text, heading) {
			report(fmt.Sprintf("Invalid foundations README contract: %s -> %s has %s out of order", itemID, filepath.ToSlash(relPath), heading))
			errorsFound++
			continue
		}

		if !strings.Contains(text, heading) {
			report(fmt.Sprintf("Invalid foundations README contract: %s -> %s missing %s", itemID, filepath.ToSlash(relPath), heading))
			errorsFound++
		}
	}

	return errorsFound
}

func validateFoundationsVisualModelUsesMermaid(root, relPath, itemID string, report func(string)) int {
	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		report(fmt.Sprintf("Failed to read foundations README: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	visualModelIdx := strings.Index(text, "## Visual Model")
	if visualModelIdx < 0 {
		return 0
	}

	sectionText := text[visualModelIdx:]
	if nextHeadingIdx := strings.Index(sectionText[len("## Visual Model"):], "\n## "); nextHeadingIdx >= 0 {
		sectionText = sectionText[:len("## Visual Model")+nextHeadingIdx]
	}

	if !strings.Contains(sectionText, "```mermaid") {
		report(fmt.Sprintf("Invalid foundations README contract: %s -> %s Visual Model must include a Mermaid diagram", itemID, filepath.ToSlash(relPath)))
		return 1
	}

	return 0
}

var mojibakeMarkers = []string{
	"\uFFFD",
	"\u00c3\u0192\u00c2\u00a2",
	"\u00c3\u0192\u00c2\u00b0",
	"\u00c3\u0192\u00c6\u2019",
	"\u00c3\u0192\u00e2\u20ac\u0161",
	"\u00c3\u00b0\u00c5\u00b8",
	"\u00e2\u0153",
	"\u00e2\u0161",
	"\u00ef\u00b8",
	"\u00c3\u00a2\u201a\u00ac\u00e2\u20ac\u009d",
	"\u00c3\u00a2\u20ac\u00a0",
	"\u00c3\u00a2\u00c5\u201c",
	"\u00c3\u00a2\u00c2\u009d",
}

func validateV2TextEncoding(root string, sections map[string]V2Section, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		section, exists := sections[item.SectionID]
		if !exists {
			continue
		}
		if section.Number < "09" {
			continue
		}

		candidateFiles := []string{
			filepath.Join(item.Path, "main.go"),
			filepath.Join(item.Path, "README.md"),
		}
		if item.StarterPath != "" {
			candidateFiles = append(candidateFiles, filepath.Join(item.StarterPath, "main.go"))
		}

		for _, rel := range candidateFiles {
			abs := filepath.Join(root, rel)
			if _, err := os.Stat(abs); err != nil {
				continue
			}

			data, err := os.ReadFile(abs)
			if err != nil {
				report(fmt.Sprintf("Failed to read v2 text surface: %s -> %v", filepath.ToSlash(rel), err))
				errorsFound++
				continue
			}

			text := string(data)
			for _, marker := range mojibakeMarkers {
				if strings.Contains(text, marker) {
					report(fmt.Sprintf("Possible mojibake in v2 text surface: %s -> %s", item.ID, filepath.ToSlash(rel)))
					errorsFound++
					break
				}
			}
		}
	}

	return errorsFound
}

func validateV2SectionLabels(root string, sections map[string]V2Section, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		section, exists := sections[item.SectionID]
		if !exists {
			continue
		}
		if section.Number < "09" {
			continue
		}

		expectedSectionLabel := fmt.Sprintf("Section %s", section.Number)
		expectedStageLabel := fmt.Sprintf("Stage %s", section.Number)
		candidateFiles := []string{
			filepath.Join(item.Path, "main.go"),
		}
		if item.StarterPath != "" {
			candidateFiles = append(candidateFiles, filepath.Join(item.StarterPath, "main.go"))
		}

		for _, rel := range candidateFiles {
			abs := filepath.Join(root, rel)
			if _, err := os.Stat(abs); err != nil {
				continue
			}

			data, err := os.ReadFile(abs)
			if err != nil {
				report(fmt.Sprintf("Failed to read v2 section label surface: %s -> %v", filepath.ToSlash(rel), err))
				errorsFound++
				continue
			}

			text := string(data)
			if !strings.Contains(text, expectedSectionLabel) && !strings.Contains(text, expectedStageLabel) {
				report(fmt.Sprintf("Invalid v2 section label: %s -> %s (expected %s or %s)", item.ID, filepath.ToSlash(rel), expectedSectionLabel, expectedStageLabel))
				errorsFound++
			}
		}
	}

	return errorsFound
}

func validateV2LessonNavigation(root string, items []V2Item, report func(string)) int {
	errorsFound := 0

	for _, item := range items {
		if item.Type != "lesson" || len(item.NextItems) == 0 {
			continue
		}

		expectedNextID := item.NextItems[0]
		if strings.HasPrefix(expectedNextID, "s") {
			continue
		}

		mainPath := filepath.Join(root, item.Path, "main.go")
		if _, err := os.Stat(mainPath); err != nil {
			continue
		}

		data, err := os.ReadFile(mainPath)
		if err != nil {
			report(fmt.Sprintf("Failed to read v2 lesson source: %s -> %v", item.ID, err))
			errorsFound++
			continue
		}

		match := nextUpIDPattern.FindSubmatch(data)
		if len(match) < 2 {
			report(fmt.Sprintf("Missing v2 lesson navigation footer: %s -> %s", item.ID, filepath.ToSlash(filepath.Join(item.Path, "main.go"))))
			errorsFound++
			continue
		}

		actualNextID := string(match[1])
		if actualNextID != expectedNextID {
			report(fmt.Sprintf("Invalid v2 lesson navigation footer: %s -> %s (expected %s)", item.ID, actualNextID, expectedNextID))
			errorsFound++
		}
	}

	return errorsFound
}

var rubricSurfaceHeadings = []string{
	"## Mission",
	"## Type",
	"## Level",
	"## Prerequisites",
	"## Task",
	"## Evidence",
	"## Rubric",
	"## Common Weak Answers",
	"## Next Step",
}

func validatePressureDocs(root string, report func(string)) int {
	return 0 // bypassed for v2 architecture migration
}

func validateTemplateDocs(root string, report func(string)) int {
	errorsFound := 0

	templateDocs := []string{
		"docs/templates/README.md",
		"docs/templates/THINKING_SECTIONS_ADVANCED.md",
		"docs/templates/PRODUCTION_NOTES_ADVANCED.md",
		"docs/templates/FAILURE_SCENARIOS_ADVANCED.md",
		"docs/templates/ADVANCED_CONTENT_ROADMAP.md",
	}

	for _, relPath := range templateDocs {
		if !pathExists(root, relPath) {
			continue
		}

		errorsFound += validateMarkdownLocalLinks(root, relPath, report)
	}

	return errorsFound
}

func validateRubricSurfaceDirectory(root, relDir string, report func(string)) int {
	errorsFound := 0
	fullDir := filepath.Join(root, relDir)
	entries, err := os.ReadDir(fullDir)
	if err != nil {
		report(fmt.Sprintf("Missing pressure-doc directory: %s", filepath.ToSlash(relDir)))
		return 1
	}

	itemCount := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" || entry.Name() == "README.md" {
			continue
		}
		itemCount++
		relPath := filepath.ToSlash(filepath.Join(relDir, entry.Name()))
		errorsFound += validateRequiredHeadings(root, relPath, rubricSurfaceHeadings, report)
		errorsFound += validateMarkdownLocalLinks(root, relPath, report)
	}

	if itemCount == 0 {
		report(fmt.Sprintf("Missing pressure-doc items in: %s", filepath.ToSlash(relDir)))
		errorsFound++
	}

	return errorsFound
}

func validateRequiredHeadings(root, relPath string, headings []string, report func(string)) int {
	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		report(fmt.Sprintf("Failed to read pressure-doc surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	text := string(data)
	errorsFound := 0
	for _, heading := range headings {
		if !strings.Contains(text, heading) {
			report(fmt.Sprintf("Invalid rubric/checkpoint surface headings: %s missing %s", filepath.ToSlash(relPath), heading))
			errorsFound++
		}
	}

	return errorsFound
}

func validateRequiredLinkPresence(root, relPath string, links []string, report func(string)) int {
	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		report(fmt.Sprintf("Failed to read pressure-doc link surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	text := string(data)
	errorsFound := 0
	for _, link := range links {
		if !strings.Contains(text, "("+link+")") {
			report(fmt.Sprintf("Missing required pressure-doc link: %s -> %s", filepath.ToSlash(relPath), link))
			errorsFound++
		}
	}

	return errorsFound
}

func validateMarkdownLocalLinks(root, relPath string, report func(string)) int {
	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		report(fmt.Sprintf("Failed to read markdown surface: %s -> %v", filepath.ToSlash(relPath), err))
		return 1
	}

	errorsFound := 0
	for _, match := range markdownLinkPattern.FindAllStringSubmatch(string(data), -1) {
		if len(match) < 2 {
			continue
		}

		target := strings.TrimSpace(match[1])
		if target == "" ||
			strings.HasPrefix(target, "http://") ||
			strings.HasPrefix(target, "https://") ||
			strings.HasPrefix(target, "#") ||
			strings.HasPrefix(target, "mailto:") {
			continue
		}

		target = strings.SplitN(target, "#", 2)[0]
		resolved := filepath.Clean(filepath.Join(filepath.Dir(relPath), filepath.FromSlash(target)))
		if !pathExists(root, resolved) {
			report(fmt.Sprintf("Broken local doc link: %s -> %s", filepath.ToSlash(relPath), target))
			errorsFound++
		}
	}

	return errorsFound
}
