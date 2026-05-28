package main

import "strings"

func ValidateReadmeContracts(cur *Curriculum, report *Report) {
	check := "readme-contracts"
	contracts := cur.Files["readme.contracts.json"]
	rawContracts := list(contracts, "contracts")
	if len(rawContracts) == 0 {
		// Some versions store contracts in a map/object. Accept either, but require the document exists.
		if _, ok := contracts["contract_definitions"]; !ok {
			report.Warn(check, "readme.contracts.json", "no contracts array found; ensure contract definitions are explicit")
		}
	}
	// Every referenced contract ID should be known or follow the accepted v3 naming convention.
	known := map[string]bool{}
	for _, rv := range rawContracts {
		c, ok := rv.(map[string]any)
		if ok {
			known[str(c, "contract_id")] = true
		}
	}
	isAcceptable := func(id string) bool {
		if id == "" {
			return false
		}
		if known[id] {
			return true
		}
		return strings.HasSuffix(id, ".v3") || strings.Contains(id, "module.v3") || strings.Contains(id, "lesson.v3") || strings.Contains(id, "project.v3")
	}
	for _, mod := range cur.Modules {
		rc := obj(mod, "readme_contract")
		cid := str(rc, "contract_id")
		if !isAcceptable(cid) {
			report.Error(check, str(mod, "id"), "unknown readme_contract id %q", cid)
		}
		if str(rc, "documentation_mode") == "" {
			report.Error(check, str(mod, "id"), "readme_contract.documentation_mode is required")
		}
	}
	for _, item := range cur.Items {
		rc := obj(item, "readme_contract")
		cid := str(rc, "contract_id")
		if !isAcceptable(cid) {
			report.Error(check, str(item, "id"), "unknown readme_contract id %q", cid)
		}
		if str(rc, "documentation_mode") == "" {
			report.Error(check, str(item, "id"), "readme_contract.documentation_mode is required")
		}
		if v, ok := rc["required"].(bool); ok && !v {
			report.Error(check, str(item, "id"), "readme_contract.required must be true")
		}
	}
	for _, p := range cur.Projects {
		rc := obj(p, "readme_contract")
		cid := str(rc, "contract_id")
		if !isAcceptable(cid) {
			report.Error(check, str(p, "id"), "unknown readme_contract id %q", cid)
		}
	}
	for _, a := range cur.Assessments {
		rc := obj(a, "readme_contract")
		cid := str(rc, "contract_id")
		if cid != "" && !isAcceptable(cid) {
			report.Error(check, str(a, "id"), "unknown readme_contract id %q", cid)
		}
	}
}
