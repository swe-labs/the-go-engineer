package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateCrossrefs(cur *Curriculum, report *Report) {
	check := "crossrefs"
	allowed := map[string]bool{"builds_on": true, "preview_only": true, "related": true, "reinforced_in": true}
	// Item-level crossrefs.
	for _, item := range cur.Items {
		id := str(item, "id")
		cr := obj(item, "crossrefs")
		if cr == nil {
			report.Error(check, id, "missing crossrefs object")
			continue
		}
		for rel, raw := range cr {
			if !allowed[rel] {
				report.Error(check, id, "unknown crossref relation %q", rel)
				continue
			}
			refs, ok := raw.([]any)
			if !ok {
				report.Error(check, id, "crossrefs.%s must be an array", rel)
				continue
			}
			for _, rv := range refs {
				r, ok := rv.(map[string]any)
				if !ok {
					report.Error(check, id, "crossrefs.%s contains non-object", rel)
					continue
				}
				target := str(r, "target_id")
				if target == "" {
					report.Error(check, id, "crossrefs.%s missing target_id", rel)
					continue
				}
				if !validID(cur, target) {
					report.Error(check, id, "crossrefs.%s target %q does not resolve", rel, target)
				}
				reason := str(r, "reason")
				checkReason(report, check, id, rel, reason)
			}
		}
	}
	// Global crossrefs.
	cross := obj(cur.Files["crossrefs.json"], "crossrefs")
	refs := list(cross, "references")
	for idx, rv := range refs {
		r, ok := rv.(map[string]any)
		if !ok {
			report.Error(check, "crossrefs.json", "references[%d] is not an object", idx)
			continue
		}
		from := str(r, "from_id")
		target := str(r, "target_id")
		if target == "" {
			target = str(r, "to_id")
		}
		rel := str(r, "relation")
		if !allowed[rel] {
			report.Error(check, fmt.Sprintf("references[%d]", idx), "invalid relation %q", rel)
		}
		if !validID(cur, from) {
			report.Error(check, fmt.Sprintf("references[%d]", idx), "from_id %q does not resolve", from)
		}
		if !validID(cur, target) {
			report.Error(check, fmt.Sprintf("references[%d]", idx), "target_id/to_id %q does not resolve", target)
		}
		checkReason(report, check, fmt.Sprintf("references[%d]", idx), rel, str(r, "reason"))
	}
}

func checkReason(report *Report, check, entity, rel, reason string) {
	if len(strings.TrimSpace(reason)) < 30 {
		report.Error(check, entity, "%s reason is too short", rel)
	}
	for _, pat := range genericReasonPatterns {
		if pat.MatchString(reason) {
			report.Error(check, entity, "%s reason is generic: %q", rel, reason)
		}
	}
	// Require reasons to communicate a pedagogical relation.
	lowered := strings.ToLower(reason)
	anchors := []string{"because", "requires", "extends", "contrasts", "reinforces", "previews", "applies", "connects", "uses", "builds"}
	hasAnchor := false
	for _, a := range anchors {
		if strings.Contains(lowered, a) {
			hasAnchor = true
			break
		}
	}
	if !hasAnchor {
		report.Warn(check, entity, "%s reason should state explicit pedagogical relationship", rel)
	}
}

func ValidateConcepts(cur *Curriculum, report *Report) {
	check := "concepts"
	seen := map[string]bool{}
	for _, c := range cur.Concepts {
		name := str(c, "concept")
		if name == "" {
			report.Error(check, "concepts.json", "concept entry missing concept name")
			continue
		}
		if seen[name] {
			report.Error(check, name, "duplicate concept")
		}
		seen[name] = true
		owner := str(c, "canonical_owner")
		if !validID(cur, owner) {
			report.Error(check, name, "canonical_owner %q does not resolve", owner)
		}
		if len(stringsList(c, "reinforcement_locations")) == 0 {
			report.Error(check, name, "must have at least one reinforcement location")
		}
		for _, loc := range stringsList(c, "preview_locations") {
			if !validID(cur, loc) {
				report.Error(check, name, "preview_location %q does not resolve", loc)
			}
		}
		for _, loc := range stringsList(c, "reinforcement_locations") {
			if !validID(cur, loc) {
				report.Error(check, name, "reinforcement_location %q does not resolve", loc)
			}
		}
	}
}

func ValidateZeroMagic(cur *Curriculum, report *Report) {
	check := "zero-magic"
	requiredStrings := []string{"problem_solved", "why_it_exists", "mental_model", "under_the_hood", "how_go_uses_it", "real_world_usage", "proof_of_understanding"}
	requiredArrays := []string{"beginner_mistakes", "execution_timeline", "failure_modes", "hidden_magic_checks", "performance_implications", "step_by_step_execution"}
	for _, item := range cur.Items {
		id := str(item, "id")
		zm := obj(item, "zero_magic")
		if zm == nil {
			report.Error(check, id, "missing zero_magic block")
			continue
		}
		for _, key := range requiredStrings {
			value := str(zm, key)
			if value == "" {
				report.Error(check, id, "zero_magic.%s missing", key)
				continue
			}
			if len(value) < 40 {
				report.Error(check, id, "zero_magic.%s is too short for world-class explanation", key)
			}
		}
		for _, key := range requiredArrays {
			if len(list(zm, key)) == 0 {
				report.Error(check, id, "zero_magic.%s missing or empty", key)
			}
		}
		// A world-class lesson needs more than one beginner mistake and failure mode unless it is a pure orientation lesson.
		if str(item, "phase") != "orientation" {
			if len(list(zm, "beginner_mistakes")) < 2 {
				report.Warn(check, id, "should include at least two beginner mistakes")
			}
			if len(list(zm, "failure_modes")) < 2 {
				report.Warn(check, id, "should include at least two failure modes")
			}
		}
		raw, _ := json.Marshal(zm)
		scanPlaceholderText(report, check, id, string(raw))
	}
}

func ValidateNoPlaceholders(cur *Curriculum, report *Report) {
	check := "no-placeholders"
	for name, f := range cur.Files {
		raw, _ := json.Marshal(f)
		scanPlaceholderText(report, check, name, string(raw))
	}
}

func scanPlaceholderText(report *Report, check, entity, text string) {
	// Metadata may legitimately mention words like TODO in explanations about starter files.
	// For metadata, fail only on exact placeholder literal values or unmistakable filler text.
	trimmed := strings.TrimSpace(strings.ToLower(text))
	badExact := map[string]bool{"placeholder": true, "scaffolded": true, "todo": true, "tbd": true, "fixme": true, "coming soon": true, "lorem ipsum": true}
	if badExact[trimmed] {
		report.Error(check, entity, "contains placeholder literal %q", text)
		return
	}
	// Exact literal checks only for metadata. Repository content uses strict scanning.
}

func scanStrictAuthoredContent(report *Report, check, entity, text string) {
	lower := strings.ToLower(text)
	allowed := []string{"sql placeholder", "sql placeholders", "no placeholder", "placeholder content", "placeholder status"}
	for _, pat := range placeholderPatterns {
		if !pat.MatchString(text) {
			continue
		}
		allowedHit := false
		for _, a := range allowed {
			if strings.Contains(lower, a) {
				allowedHit = true
				break
			}
		}
		if !allowedHit {
			report.Error(check, entity, "contains placeholder-like text matching %q", pat.String())
		}
	}
}

func ValidateCognitiveLoad(cur *Curriculum, report *Report) {
	check := "cognitive-load"
	expectedPhase := map[string]string{"module-09": "data", "module-10": "security", "module-12": "reliability", "module-13": "performance", "module-15": "delivery"}
	hardModules := map[string]bool{"module-03": true, "module-04": true, "module-05": true, "module-09": true, "module-10": true, "module-11": true, "module-12": true, "module-13": true, "module-14": true, "module-18": true}
	for _, m := range cur.Modules {
		id := str(m, "id")
		if want, ok := expectedPhase[id]; ok && str(m, "phase") != want {
			report.Error(check, id, "phase must be %q, got %q", want, str(m, "phase"))
		}
		if hardModules[id] && !boolVal(m, "contains_foundational_hard_concepts") {
			report.Error(check, id, "must be marked contains_foundational_hard_concepts")
		}
		if str(m, "cognitive_load") == "high" && !boolVal(m, "recommended_break_after") {
			report.Warn(check, id, "high cognitive load should usually recommend a break")
		}
	}
	// Guard against future overloading.
	counts := map[string]int{}
	for _, item := range cur.Items {
		counts[str(item, "module_id")]++
	}
	for mid, count := range counts {
		if mid == "module-18" {
			continue
		}
		if count > 32 {
			report.Error(check, mid, "module is oversized with %d items; split or justify", count)
		}
	}
	if counts["module-15"] != 19 {
		report.Error(check, "module-15", "delivery module must remain split/clean at 19 items, got %d", counts["module-15"])
	}
}

func ValidateFailures(cur *Curriculum, report *Report) {
	check := "failures"
	failures := cur.Files["failures.json"]
	if str(failures, "document_type") != "operational_failure_taxonomy" {
		report.Error(check, "failures.json", "document_type must be operational_failure_taxonomy")
	}
	categories := map[string]bool{}
	for _, rv := range list(failures, "failure_categories") {
		fc, ok := rv.(map[string]any)
		if !ok {
			report.Error(check, "failures.json", "failure category must be object")
			continue
		}
		cat := str(fc, "category")
		if cat == "" {
			report.Error(check, "failures.json", "failure category missing category")
		}
		if categories[cat] {
			report.Error(check, cat, "duplicate failure category")
		}
		categories[cat] = true
		if len(str(fc, "description")) < 40 {
			report.Warn(check, cat, "description should explain production impact")
		}
		for _, mid := range stringsList(fc, "modules") {
			if _, ok := cur.ModuleByID[mid]; !ok && mid != "opslane" {
				report.Error(check, cat, "unknown module %q", mid)
			}
		}
	}
	coverage, _ := failures["required_coverage"].(map[string]any)
	for mid, raw := range coverage {
		if _, ok := cur.ModuleByID[mid]; !ok && mid != "opslane" {
			report.Error(check, mid, "required_coverage references unknown module")
		}
		arr, _ := raw.([]any)
		if len(arr) == 0 {
			report.Error(check, mid, "required_coverage cannot be empty")
		}
		for _, v := range arr {
			cat, _ := v.(string)
			if !categories[cat] {
				report.Error(check, mid, "unknown failure category %q", cat)
			}
		}
	}
}

func relationSpecificity(reason string) bool {
	words := strings.Fields(reason)
	if len(words) < 8 {
		return false
	}
	return regexp.MustCompile(`[a-zA-Z]{5,}`).MatchString(reason)
}
