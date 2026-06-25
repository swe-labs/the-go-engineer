package main

import "strings"

func validateReadmeContracts(m Metadata) ValidationResult {
	var r ValidationResult
	if len(m.RawContracts) == 0 {
		r.Errorf("readme.contracts.json is empty or missing")
		return r
	}
	raw, ok := m.RawContracts["contracts"].(map[string]any)
	if !ok || len(raw) == 0 {
		// Some versions use top-level named contracts; accept those if the file has policy fields.
		if _, hasPolicy := m.RawContracts["repository_layout_policy"]; !hasPolicy {
			r.Errorf("readme.contracts.json must define contracts or repository_layout_policy")
		}
	}
	policy, _ := m.RawContracts["repository_layout_policy"].(map[string]any)
	if policy != nil {
		if root, _ := policy["curriculum_root"].(string); root != "curriculum/" {
			r.Errorf("readme contracts curriculum_root must be curriculum/")
		}
		forbidden, _ := policy["forbidden_folder_names"].([]any)
		for _, value := range forbidden {
			name := strings.ToLower(strings.TrimSpace(value.(string)))
			if name == "codex" || name == "ai" || name == "agent" || name == "bot" || name == "llm" {
				continue
			}
		}
	}
	return r
}
