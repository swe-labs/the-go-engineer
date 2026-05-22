package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Assessment struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Type           string   `json:"type"`
	Status         string   `json:"status"`
	TargetIDs      []string `json:"target_ids"`
	Items          []any    `json:"-"`
}

type AssessmentsBundle struct {
	SchemaVersion       string `json:"schema_version"`
	DocumentType        string `json:"document_type"`
	CurriculumVersion   string `json:"curriculum_version"`
	LastUpdated         string `json:"last_updated"`
	Assessments         []any  `json:"assessments"`
}

func main() {
	wd, _ := os.Getwd()
	curDir := filepath.Join(wd, "curriculum")
	if _, err := os.Stat(curDir); err != nil {
		for dir := wd; dir != "."; dir = filepath.Dir(dir) {
			candidate := filepath.Join(dir, "curriculum")
			if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
				curDir = candidate
				break
			}
		}
	}

	path := filepath.Join(curDir, "assessments.json")
	fmt.Printf("Reading %s\n", path)
	data, _ := os.ReadFile(path)

	// Parse as raw to manipulate
	var rawMap map[string]any
	json.Unmarshal(data, &rawMap)

	assessments := rawMap["assessments"].([]any)

	// Remove old module-08 assessment (Context and Time)
	var keep []any
	for _, a := range assessments {
		m := a.(map[string]any)
		if m["id"].(string) == "assessment-module-08" {
			fmt.Printf("  Removed %s\n", m["id"])
			continue
		}
		keep = append(keep, a)
	}
	assessments = keep

	// Rename/update assessments
	type rename struct {
		oldID string
		newID string
		title string
	}
	renames := []rename{
		{"assessment-module-09", "assessment-module-08", "Module 08 Assessment — HTTP and REST APIs"},
		{"assessment-module-10", "assessment-module-09", "Module 09 Assessment — SQL and PostgreSQL Persistence"},
		{"assessment-module-11", "assessment-module-10", "Module 10 Assessment — Authentication, Authorization, and Security"},
		{"assessment-module-12", "assessment-module-11", "Module 11 Assessment — Lifecycle, Context, and Concurrency"},
		{"assessment-module-14", "assessment-module-14", "Module 14 Assessment — Architecture and Distributed Systems"},
		{"assessment-module-15", "assessment-module-15", "Module 15 Assessment — Docker, CI/CD, and Deployment"},
		{"assessment-module-16", "assessment-module-16", "Module 16 Assessment — Portfolio and Interview Readiness"},
		{"assessment-module-17", "assessment-module-17", "Module 17 Assessment — Electives and Advanced Engineering"},
	}

	for i, a := range assessments {
		m := a.(map[string]any)
		for _, r := range renames {
			if m["id"].(string) == r.oldID {
				m["id"] = r.newID
				m["title"] = r.title
				fmt.Printf("  Updated %s → %s\n", r.oldID, r.newID)
				_ = i
			}
		}
	}

	// Split old module-13 assessment (Obs+Perf) into module-12 (Obs) + module-13 (Perf)
	var splitKeep []any
	for _, a := range assessments {
		m := a.(map[string]any)
		if m["id"].(string) == "assessment-module-13" {
			// Create module-12 assessment (Observability)
			obs := cloneMap(m)
			obs["id"] = "assessment-module-12"
			obs["title"] = "Module 12 Assessment — Observability and Diagnostics"
			obs["target_ids"] = []any{}
			for _, tid := range m["target_ids"].([]any) {
				id := tid.(string)
				if strings.HasPrefix(id, "core-12-") {
					sub := id[8:]
					num := 0
					fmt.Sscanf(sub, "%d", &num)
					if num <= 9 || num >= 18 {
						obs["target_ids"] = append(obs["target_ids"].([]any), id)
					}
				}
			}
			fmt.Printf("  Created %s with %d targets\n", obs["id"], len(obs["target_ids"].([]any)))
			splitKeep = append(splitKeep, obs)

			// Create module-13 assessment (Performance)
			perf := cloneMap(m)
			perf["id"] = "assessment-module-13"
			perf["title"] = "Module 13 Assessment — Performance and Memory Engineering"
			perf["target_ids"] = []any{}
			for _, tid := range m["target_ids"].([]any) {
				id := tid.(string)
				if strings.HasPrefix(id, "core-12-") {
					sub := id[8:]
					num := 0
					fmt.Sscanf(sub, "%d", &num)
					if num >= 10 && num <= 17 {
						perf["target_ids"] = append(perf["target_ids"].([]any), id)
					}
				}
			}
			fmt.Printf("  Created %s with %d targets\n", perf["id"], len(perf["target_ids"].([]any)))
			splitKeep = append(splitKeep, perf)
			fmt.Printf("  Removed old %s\n", m["id"])
		} else {
			splitKeep = append(splitKeep, a)
		}
	}
	assessments = splitKeep

	rawMap["assessments"] = assessments
	rawMap["schema_version"] = "3.1.0"
	rawMap["curriculum_version"] = "3.1.0-draft.1"

	out, _ := json.MarshalIndent(rawMap, "", "    ")
	os.WriteFile(path, out, 0644)
	fmt.Println("  Saved assessments.json")

	// ===== Fix item assessment_id references =====
	fmt.Println("\n=== Fixing item assessment_id references ===")
	corePath := filepath.Join(curDir, "path.core.json")
	coreData, _ := os.ReadFile(corePath)

	assessmentIDMap := map[string]string{
		"assessment-module-08": "assessment-module-11", // old context/time → merged
		"assessment-module-09": "assessment-module-08", // HTTP shift
		"assessment-module-10": "assessment-module-09", // SQL shift
		"assessment-module-11": "assessment-module-10", // Auth shift
		"assessment-module-12": "assessment-module-11", // Concurrency → Lifecycle
		"assessment-module-13": "assessment-module-12", // Obs+Perf → Obs
	}

	content := string(coreData)
	fixedCount := 0
	for oldID, newID := range assessmentIDMap {
		count := strings.Count(content, fmt.Sprintf(`"assessment_id": "%s"`, oldID))
		if count > 0 {
			content = strings.ReplaceAll(content, fmt.Sprintf(`"assessment_id": "%s"`, oldID), fmt.Sprintf(`"assessment_id": "%s"`, newID))
			fixedCount += count
			fmt.Printf("  Replaced %d refs: %s → %s\n", count, oldID, newID)
		}
	}

	// Old module-08 items (core-11-01..09) now in module-11 → assessment-module-11
	// The map above handles this: assessment-module-08 → assessment-module-11

	os.WriteFile(corePath, []byte(content), 0644)
	fmt.Printf("  Fixed %d total assessment_id references in path.core.json\n", fixedCount)

	fmt.Println("\nDone. Run validator to verify.")
}

func cloneMap(m map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		result[k] = v
	}
	return result
}
