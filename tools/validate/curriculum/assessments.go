package main

import "strings"

func validateProjects(m Metadata) ValidationResult {
	var r ValidationResult
	modules := moduleIDs(m)
	items := itemIDs(m)
	assessments := assessmentIDs(m)
	concepts := map[string]bool{}
	for _, c := range m.Concepts.Concepts {
		concepts[c.Concept] = true
	}
	seen := map[string]bool{}

	for _, project := range m.Projects.Projects {
		if project.ID == "" {
			r.Errorf("project has empty id")
		}
		if seen[project.ID] {
			r.Errorf("duplicate project id %s", project.ID)
		}
		seen[project.ID] = true
		if _, ok := modules[project.ModuleID]; !ok {
			r.Errorf("%s unknown module_id %s", project.ID, project.ModuleID)
		}
		if !isStableStatus(project.Status) {
			r.Errorf("%s status must be stable/ready/published", project.ID)
		}
		if project.Files.ReadmePath == "" || !canonicalContentPath(project.Files.ReadmePath) {
			r.Errorf("%s missing canonical readme_path", project.ID)
		}
		if project.AssessmentID != "" {
			if _, ok := assessments[project.AssessmentID]; !ok {
				r.Errorf("%s unknown assessment_id %s", project.ID, project.AssessmentID)
			}
		}
		for _, target := range project.TargetIDs {
			if _, ok := items[target]; !ok {
				r.Errorf("%s unknown target_id %s", project.ID, target)
			}
		}
		for _, prereq := range project.Prerequisites {
			if _, ok := items[prereq]; ok {
				continue
			}
			if _, ok := modules[prereq]; ok {
				continue
			}
			if _, ok := concepts[prereq]; ok {
				continue
			}
			r.Errorf("%s unknown prerequisite %s", project.ID, prereq)
		}
		if rubric, ok := project.Rubric.(map[string]any); ok {
			if err := validateWeightMap(rubric); err != "" {
				r.Errorf("%s rubric invalid: %s", project.ID, err)
			}
		}
	}
	return r
}

func validateAssessments(m Metadata) ValidationResult {
	var r ValidationResult
	items := itemIDs(m)
	modules := moduleIDs(m)
	projects := projectIDs(m)
	seen := map[string]bool{}
	for _, a := range m.Assessments.Assessments {
		if a.ID == "" {
			r.Errorf("assessment has empty id")
		}
		if seen[a.ID] {
			r.Errorf("duplicate assessment id %s", a.ID)
		}
		seen[a.ID] = true
		if !isStableStatus(a.Status) {
			r.Errorf("%s status must be stable/ready/published", a.ID)
		}
		if a.ModuleID != "" {
			if _, ok := modules[a.ModuleID]; !ok {
				r.Errorf("%s unknown module_id %s", a.ID, a.ModuleID)
			}
		}
		if len(a.TargetIDs) == 0 {
			r.Errorf("%s target_ids is empty", a.ID)
		}
		for _, target := range a.TargetIDs {
			if _, ok := items[target]; ok {
				continue
			}
			if _, ok := modules[target]; ok {
				continue
			}
			if _, ok := projects[target]; ok {
				continue
			}
			r.Errorf("%s unknown target_id %s", a.ID, target)
		}
		if path := a.Files["readme_path"]; path == "" || !canonicalContentPath(path) {
			r.Errorf("%s missing canonical files.readme_path", a.ID)
		}
		if len(a.Criteria) == 0 && a.Rubric == nil {
			r.Errorf("%s must define criteria or rubric", a.ID)
		}
	}
	return r
}

func validateWeightMap(m map[string]any) string {
	total := 0.0
	found := false
	for key, value := range m {
		if strings.Contains(strings.ToLower(key), "weight") {
			switch n := value.(type) {
			case float64:
				total += n
				found = true
			case int:
				total += float64(n)
				found = true
			}
		}
	}
	if found && (total < 99.9 || total > 100.1) {
		return "weights must sum to 100"
	}
	return ""
}
