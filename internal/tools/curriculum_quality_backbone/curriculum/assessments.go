package main

import "strings"

func ValidateAssessments(cur *Curriculum, report *Report) {
	check := "assessments"
	for _, a := range cur.Assessments {
		id := str(a, "id")
		if id == "" {
			report.Error(check, "assessments.json", "assessment missing id")
			continue
		}
		if str(a, "status") != "stable" {
			report.Error(check, id, "status must be stable")
		}
		if str(a, "readme_status") != "golden" {
			report.Error(check, id, "readme_status must be golden")
		}
		if len(stringsList(a, "target_ids")) == 0 {
			report.Error(check, id, "target_ids cannot be empty")
		}
		for _, tid := range stringsList(a, "target_ids") {
			if !validID(cur, tid) {
				report.Error(check, id, "target_id %q does not resolve", tid)
			}
		}
		score := num(a, "passing_score")
		if score < 70 || score > 100 {
			report.Error(check, id, "passing_score must be between 70 and 100")
		}
		if len(list(a, "evidence_required")) == 0 {
			report.Error(check, id, "evidence_required cannot be empty")
		}
		if len(list(a, "manual_review_questions")) == 0 {
			report.Error(check, id, "manual_review_questions cannot be empty")
		}
		if str(a, "retake_policy") == "" {
			report.Error(check, id, "retake_policy is required")
		}
		criteria := list(a, "criteria")
		if len(criteria) == 0 {
			report.Error(check, id, "criteria cannot be empty")
		}
		total := 0
		for i, rv := range criteria {
			c, ok := rv.(map[string]any)
			if !ok {
				report.Error(check, id, "criteria[%d] is not an object", i)
				continue
			}
			name := str(c, "name")
			if name == "" {
				name = str(c, "criterion")
			}
			if len(name) < 4 {
				report.Error(check, id, "criteria[%d] needs a specific name", i)
			}
			if strings.Contains(strings.ToLower(name), "general") || strings.Contains(strings.ToLower(name), "overall") {
				report.Error(check, id, "criteria[%d] appears generic: %q", i, name)
			}
			weight := num(c, "weight")
			if weight <= 0 {
				report.Error(check, id, "criteria[%d] missing positive weight", i)
			}
			total += weight
		}
		if len(criteria) > 0 && total != 100 {
			report.Error(check, id, "criteria weights must sum to 100, got %d", total)
		}
	}
}

func ValidateProjects(cur *Curriculum, report *Report) {
	check := "projects"
	for _, p := range cur.Projects {
		id := str(p, "id")
		if id == "" {
			report.Error(check, "projects.json", "project missing id")
			continue
		}
		required := []string{"title", "slug", "type", "stage", "difficulty", "status", "summary", "module_id", "placement_anchor_item_id", "prerequisites", "prerequisite_item_ids", "required_concepts", "reinforcement_targets", "deliverables", "tasks", "assessment_id", "rubric", "verification", "portfolio", "readme_status", "readme_contract"}
		for _, key := range required {
			if _, ok := p[key]; !ok {
				report.Error(check, id, "missing %s", key)
			}
		}
		if str(p, "status") != "stable" {
			report.Error(check, id, "status must be stable")
		}
		if str(p, "readme_status") != "golden" {
			report.Error(check, id, "readme_status must be golden")
		}
		if _, ok := cur.ModuleByID[str(p, "module_id")]; !ok {
			report.Error(check, id, "module_id %q does not resolve", str(p, "module_id"))
		}
		anchor := str(p, "placement_anchor_item_id")
		if _, ok := cur.ItemByID[anchor]; !ok {
			report.Error(check, id, "placement_anchor_item_id %q does not resolve", anchor)
		}
		for _, key := range []string{"prerequisites", "prerequisite_item_ids", "reinforces"} {
			for _, target := range stringsList(p, key) {
				if !validID(cur, target) {
					report.Error(check, id, "%s target %q does not resolve", key, target)
				}
			}
		}
		for _, key := range []string{"required_concepts", "reinforcement_targets"} {
			for _, target := range stringsList(p, key) {
				if !validID(cur, target) {
					report.Error(check, id, "%s target %q does not resolve", key, target)
				}
			}
		}
		if len(stringsList(p, "deliverables")) < 3 {
			report.Error(check, id, "must have at least 3 deliverables")
		}
		if len(list(p, "tasks")) < 2 {
			report.Error(check, id, "must have at least 2 tasks")
		}
		verification := obj(p, "verification")
		if verification == nil {
			report.Error(check, id, "missing verification block")
		} else if str(verification, "mode") == "" {
			report.Error(check, id, "verification.mode is required")
		}
		aid := str(p, "assessment_id")
		if aid == "" {
			report.Error(check, id, "missing assessment_id")
		} else if _, ok := cur.AssessmentByID[aid]; !ok {
			report.Error(check, id, "assessment_id %q does not resolve", aid)
		}
		rubric := list(p, "rubric")
		if len(rubric) == 0 {
			report.Error(check, id, "rubric cannot be empty")
		}
		total := 0
		for i, rv := range rubric {
			r, ok := rv.(map[string]any)
			if !ok {
				report.Error(check, id, "rubric[%d] is not an object", i)
				continue
			}
			if str(r, "criterion") == "" && str(r, "name") == "" {
				report.Error(check, id, "rubric[%d] missing criterion/name", i)
			}
			if len(str(r, "description")) < 30 {
				report.Warn(check, id, "rubric[%d] description should be more specific", i)
			}
			total += num(r, "weight")
		}
		if len(rubric) > 0 && total != 100 {
			report.Error(check, id, "rubric weights must sum to 100, got %d", total)
		}
		portfolio := obj(p, "portfolio")
		if portfolio == nil {
			report.Error(check, id, "missing portfolio block")
		}
	}
}
