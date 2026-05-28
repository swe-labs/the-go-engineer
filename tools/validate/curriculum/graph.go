package main

import "strings"

func validateGraph(m Metadata) ValidationResult {
	var r ValidationResult
	modules := moduleIDs(m)
	items := itemIDs(m)
	projects := projectIDs(m)

	for _, module := range allModules(m) {
		if module.ID == "" {
			r.Errorf("module has empty id")
		}
		if !canonicalModulePath(module.Path) {
			r.Errorf("%s path must start with curriculum/modules/ or curriculum/electives/, got %s", module.ID, module.Path)
		}
		if !isStableStatus(module.Status) {
			r.Errorf("%s status must be stable/ready/published, got %q", module.ID, module.Status)
		}
		if module.ReadmeStatus != "" && module.ReadmeStatus != "golden" {
			r.Errorf("%s readme_status must be golden", module.ID)
		}
		for _, prereq := range module.Prerequisites {
			if _, ok := modules[prereq]; !ok {
				r.Errorf("%s unknown prerequisite module %s", module.ID, prereq)
			}
		}
		for _, id := range module.EntryItemIDs {
			if _, ok := items[id]; !ok {
				r.Errorf("%s unknown entry_item_id %s", module.ID, id)
			}
		}
		for _, id := range module.TerminalItemIDs {
			if _, ok := items[id]; !ok {
				r.Errorf("%s unknown terminal_item_id %s", module.ID, id)
			}
		}
	}

	for _, item := range allItems(m) {
		for _, prereq := range item.Prerequisites {
			if _, ok := items[prereq]; ok {
				continue
			}
			if _, ok := modules[prereq]; ok {
				continue
			}
			r.Errorf("%s unknown prerequisite %s", item.ID, prereq)
		}
		for _, next := range item.NextItemIDs {
			if _, ok := items[next]; ok {
				continue
			}
			if _, ok := projects[next]; ok {
				continue
			}
			r.Errorf("%s unknown next_item_id %s", item.ID, next)
		}
		checkRefs := func(kind string, refs []Reference) {
			for _, ref := range refs {
				id := ref.TargetID
				if id == "" {
					id = ref.ID
				}
				if id == "" {
					id = ref.ToID
				}
				if id == "" {
					continue
				}
				if _, ok := items[id]; ok {
					continue
				}
				if _, ok := projects[id]; ok {
					continue
				}
				r.Errorf("%s crossrefs.%s references unknown target %s", item.ID, kind, id)
			}
		}
		checkRefs("builds_on", item.CrossRefs.BuildsOn)
		checkRefs("preview_only", item.CrossRefs.PreviewOnly)
		checkRefs("reinforced_in", item.CrossRefs.ReinforcedIn)
		checkRefs("related", item.CrossRefs.Related)
	}

	counts := map[string]int{}
	for _, item := range allItems(m) {
		counts[item.ModuleID]++
	}
	for _, module := range allModules(m) {
		if counts[module.ID] == 0 && !strings.Contains(module.Path, "/electives/") {
			r.Errorf("%s has zero items", module.ID)
		}
		if counts[module.ID] > 35 {
			r.Errorf("%s has %d items; split or justify module size", module.ID, counts[module.ID])
		}
	}
	return r
}

func validateCrossrefs(m Metadata) ValidationResult {
	var r ValidationResult
	items := itemIDs(m)
	projects := projectIDs(m)
	generic := []string{"build understanding", "both ", "related to"}
	for _, ref := range m.Crossrefs.Crossrefs.References {
		if _, ok := items[ref.FromID]; !ok {
			if _, ok := projects[ref.FromID]; !ok {
				r.Errorf("crossref unknown from_id %s", ref.FromID)
			}
		}
		target := ref.TargetID
		if target == "" {
			target = ref.ToID
		}
		if _, ok := items[target]; !ok {
			if _, ok := projects[target]; !ok {
				r.Errorf("crossref unknown target_id %s", target)
			}
		}
		reason := strings.TrimSpace(strings.ToLower(ref.Reason))
		if len(reason) < 24 {
			r.Errorf("crossref %s -> %s reason is too short", ref.FromID, target)
		}
		for _, phrase := range generic {
			if strings.Contains(reason, phrase) {
				r.Errorf("crossref %s -> %s has generic reason %q", ref.FromID, target, ref.Reason)
			}
		}
	}
	return r
}
