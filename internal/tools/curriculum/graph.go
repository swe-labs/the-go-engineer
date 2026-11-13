package main

import (
	"fmt"
	"strings"
)

func validateGraph() {
	fmt.Println("=== validate-graph ===")
	ok := true

	core := loadCoreBundle()
	electives := loadElectiveBundle()

	var allItems []Item
	allItems = append(allItems, core.Items...)
	allItems = append(allItems, electives.Items...)

	itemMap := make(map[string]*Item)
	for i := range allItems {
		item := &allItems[i]
		itemMap[item.ID] = item
	}

	// Also register project IDs as valid targets
	for _, mf := range moduleManifestPaths() {
		var mm struct {
			Projects []Project `json:"projects"`
		}
		readJSON(mf, &mm)
		for _, p := range mm.Projects {
			itemMap[p.ID] = nil // exists but not an Item
		}
	}

	moduleMap := make(map[string]*Module)
	for i := range core.Modules {
		moduleMap[core.Modules[i].ID] = &core.Modules[i]
	}
	for i := range electives.Modules {
		moduleMap[electives.Modules[i].ID] = &electives.Modules[i]
	}

	coreItemIDs := make(map[string]bool)
	for _, item := range core.Items {
		coreItemIDs[item.ID] = true
	}
	moduleOrder := make(map[string]int)
	for i := range core.Modules {
		moduleOrder[core.Modules[i].ID] = core.Modules[i].Number
	}
	for i := range electives.Modules {
		moduleOrder[electives.Modules[i].ID] = electives.Modules[i].Number
	}

	// 1. Check all prerequisite item IDs exist
	for _, item := range allItems {
		for _, prereq := range item.Prerequisites {
			if _, exists := itemMap[prereq]; !exists {
				errorf("%s: prerequisite '%s' does not exist", item.ID, prereq)
				ok = false
			}
		}
		for _, np := range item.NextItemIDs {
			if _, exists := itemMap[np]; !exists {
				errorf("%s: next_item_id '%s' does not exist", item.ID, np)
				ok = false
			}
		}
	}

	// 2. Check all module prerequisites exist
	for _, mod := range core.Modules {
		for _, prereq := range mod.Prerequisites {
			if _, exists := moduleMap[prereq]; !exists {
				errorf("module '%s': prerequisite '%s' does not exist", mod.ID, prereq)
				ok = false
			}
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
					errorf("%s: crossrefs.%s target '%s' does not exist", item.ID, category, ref.TargetID)
					ok = false
				}
			}
		}
		checkRefs(item.CrossRefs.BuildsOn, "builds_on")
		checkRefs(item.CrossRefs.PreviewOnly, "preview_only")
		checkRefs(item.CrossRefs.Related, "related")
		checkRefs(item.CrossRefs.ReinforcedIn, "reinforced_in")
	}

	// 4b. builds_on must point backward. Embedded item crossrefs are the
	// single source since the central crossrefs.json registry was removed.
	buildsOn := make(map[string][]string)
	for _, item := range allItems {
		if item.CrossRefs == nil {
			continue
		}
		for _, ref := range item.CrossRefs.BuildsOn {
			buildsOn[item.ID] = append(buildsOn[item.ID], ref.TargetID)
		}
	}
	for fromID, targets := range buildsOn {
		fromItem := itemMap[fromID]
		if fromItem == nil {
			continue
		}
		fromMod := moduleOrder[fromItem.ModuleID]
		for _, targetID := range targets {
			toItem := itemMap[targetID]
			if toItem == nil {
				continue
			}
			toMod := moduleOrder[toItem.ModuleID]
			if toMod > fromMod || (toMod == fromMod && toItem.Order >= fromItem.Order) {
				errorf("%s: builds_on target '%s' is not earlier (module %d order %d vs module %d order %d)",
					fromID, targetID, toMod, toItem.Order, fromMod, fromItem.Order)
				ok = false
			}
		}
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
			if !exists || item == nil {
				continue
			}
			for _, next := range item.NextItemIDs {
				if !reachable[next] {
					queue = append(queue, next)
				}
			}
		}
		// Check all instructional items in this module are reachable
		for _, item := range allItems {
			if item.ModuleID != mod.ID {
				continue
			}
			if item.Type != "lesson" && item.Type != "lab" {
				continue
			}
			if !reachable[item.ID] {
				errorf("%s (module %s): not reachable from module entry points %v", item.ID, mod.ID, mod.EntryItemIDs)
				ok = false
			}
		}
	}

	// 6. Check no cycles in the traversal, prerequisite, and builds_on graphs.
	nextAdj := make(map[string][]string)
	prereqAdj := make(map[string][]string)
	for _, item := range allItems {
		nextAdj[item.ID] = item.NextItemIDs
		prereqAdj[item.ID] = item.Prerequisites
	}
	cycleNext := detectCycles("next_item_id", allItems, nextAdj)
	cyclePrereq := detectCycles("prerequisite", allItems, prereqAdj)
	cycleBuilds := detectCycles("builds_on", allItems, buildsOn)
	if cycleNext || cyclePrereq || cycleBuilds {
		ok = false
	} else {
		fmt.Println("  No cycles detected in next_item_id, prerequisite, or builds_on graphs.")
	}

	// 7. Check assessment_id naming convention
	for _, item := range allItems {
		if item.Proof != nil && item.Proof.AssessmentID != "" {
			if !strings.HasPrefix(item.Proof.AssessmentID, "assessment-") {
				errorf("%s: assessment_id '%s' does not follow 'assessment-*' convention", item.ID, item.Proof.AssessmentID)
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

// detectCycles reports whether the adjacency graph contains a cycle.
// It only walks nodes that exist in allItems; dangling targets are handled by
// the existence checks in validateGraph.
func detectCycles(label string, allItems []Item, adjacency map[string][]string) bool {
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	hasCycle := false

	var dfs func(id string)
	dfs = func(id string) {
		if inStack[id] {
			errorf("%s cycle detected involving item '%s'", label, id)
			hasCycle = true
			return
		}
		if visited[id] {
			return
		}
		visited[id] = true
		inStack[id] = true
		for _, next := range adjacency[id] {
			dfs(next)
		}
		inStack[id] = false
	}

	for _, item := range allItems {
		if !visited[item.ID] {
			dfs(item.ID)
		}
	}
	return hasCycle
}
