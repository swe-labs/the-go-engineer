package main

import (
	"fmt"
	"strings"
)

func validateGraph() {
	fmt.Println("=== validate-graph ===")
	ok := true

	var core CoreBundle
	var electives ElectiveBundle
	readJSON("path.core.json", &core)
	readJSON("path.electives.json", &electives)
	var crossrefs CrossrefBundle
	readJSON("crossrefs.json", &crossrefs)

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	itemMap := make(map[string]*Item)
	for i := range allItems {
		item := &allItems[i]
		itemMap[item.ID] = item
	}

	// Also register project IDs as valid targets
	var projects ProjectsBundle
	readJSON("projects.json", &projects)
	for _, p := range projects.Projects {
		itemMap[p.ID] = nil // exists but not an Item
	}

	moduleMap := make(map[string]*Module)
	for i := range core.Modules {
		moduleMap[core.Modules[i].ID] = &core.Modules[i]
	}
	for i := range electives.Modules {
		moduleMap[electives.Modules[i].ID] = &electives.Modules[i]
	}

	// 1. Check all prerequisite item IDs exist
	for _, item := range allItems {
		for _, prereq := range item.Prerequisites {
			if _, exists := itemMap[prereq]; !exists {
				warnf("%s: prerequisite '%s' does not exist", item.ID, prereq)
				ok = false
			}
		}
		for _, np := range item.NextItemIDs {
			if _, exists := itemMap[np]; !exists {
				warnf("%s: next_item_id '%s' does not exist", item.ID, np)
				ok = false
			}
		}
	}

	// 2. Check all module prerequisites exist
	for _, mod := range core.Modules {
		for _, prereq := range mod.Prerequisites {
			if _, exists := moduleMap[prereq]; !exists {
				warnf("module '%s': prerequisite '%s' does not exist", mod.ID, prereq)
				ok = false
			}
		}
	}

	// 3. Check crossref targets exist
	for _, ref := range crossrefs.Crossrefs.References {
		if _, exists := itemMap[ref.FromID]; !exists {
			warnf("crossref from '%s' does not exist", ref.FromID)
			ok = false
		}
		if _, exists := itemMap[ref.TargetID]; !exists {
			warnf("crossref to '%s' does not exist", ref.TargetID)
			ok = false
		}
	}

	// 4. Check item-level crossrefs
	for _, item := range allItems {
		if item.CrossRefs == nil {
			continue
		}
		checkRefs := func(refs []CrossrefRef, category string) {
			for _, ref := range refs {
				if _, exists := itemMap[ref.TargetID]; !exists {
					warnf("%s: crossrefs.%s target '%s' does not exist", item.ID, category, ref.TargetID)
					ok = false
				}
			}
		}
		checkRefs(item.CrossRefs.BuildsOn, "builds_on")
		checkRefs(item.CrossRefs.PreviewOnly, "preview_only")
		checkRefs(item.CrossRefs.Related, "related")
		checkRefs(item.CrossRefs.ReinforcedIn, "reinforced_in")
	}

	// 5. Check for orphan items (items not reachable from any module entry point)
	for _, mod := range core.Modules {
		reachable := make(map[string]bool)
		// BFS from entry points
		queue := make([]string, 0)
		for _, e := range mod.EntryItemIDs {
			queue = append(queue, e)
		}
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			if reachable[id] {
				continue
			}
			reachable[id] = true
			item, exists := itemMap[id]
			if !exists {
				continue
			}
			for _, next := range item.NextItemIDs {
				if !reachable[next] {
					queue = append(queue, next)
				}
			}
		}
		// Check all items in this module are reachable
		for _, item := range allItems {
			if item.ModuleID != mod.ID {
				continue
			}
			if !reachable[item.ID] {
				warnf("%s (module %s): not reachable from module entry points %v", item.ID, mod.ID, mod.EntryItemIDs)
				ok = false
			}
		}
	}

	// 6. Check no cycles in item graph
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var hasCycle bool
	var dfs func(id string)
	dfs = func(id string) {
		if inStack[id] {
			warnf("cycle detected involving item '%s'", id)
			hasCycle = true
			ok = false
			return
		}
		if visited[id] {
			return
		}
		visited[id] = true
		inStack[id] = true
		item, exists := itemMap[id]
		if exists {
			for _, next := range item.NextItemIDs {
				dfs(next)
			}
		}
		inStack[id] = false
	}
	for _, item := range allItems {
		if !visited[item.ID] {
			dfs(item.ID)
		}
	}
	if !hasCycle {
		fmt.Println("  No cycles detected in item graph.")
	}

	// 7. Check assessment_id naming convention
	for _, item := range allItems {
		if item.Proof != nil && item.Proof.AssessmentID != "" {
			if !strings.HasPrefix(item.Proof.AssessmentID, "assessment-") {
				warnf("%s: assessment_id '%s' does not follow 'assessment-*' convention", item.ID, item.Proof.AssessmentID)
				ok = false
			}
		}
	}

	if ok {
		fmt.Println("  All graph checks passed.")
	} else {
		fmt.Println("  Graph checks completed with warnings/errors.")
	}
}
