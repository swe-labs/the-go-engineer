package main

import (
	"fmt"
	"sort"
)

func ValidateGraph(cur *Curriculum, report *Report) {
	validateModuleGraph(cur, report)
	validateItemGraph(cur, report)
	validateModuleReachability(cur, report)
}

func validateModuleGraph(cur *Curriculum, report *Report) {
	check := "module-graph"
	for _, mod := range cur.Modules {
		id := str(mod, "id")
		for _, pid := range stringsList(mod, "prerequisites") {
			if _, ok := cur.ModuleByID[pid]; !ok {
				report.Error(check, id, "unknown prerequisite module %q", pid)
			}
			if pid == id {
				report.Error(check, id, "module cannot depend on itself")
			}
		}
		for _, eid := range stringsList(mod, "entry_item_ids") {
			item, ok := cur.ItemByID[eid]
			if !ok {
				report.Error(check, id, "entry item %q does not exist", eid)
				continue
			}
			if str(item, "module_id") != id {
				report.Error(check, id, "entry item %q belongs to %s", eid, str(item, "module_id"))
			}
		}
		for _, tid := range stringsList(mod, "terminal_item_ids") {
			item, ok := cur.ItemByID[tid]
			if !ok {
				report.Error(check, id, "terminal item %q does not exist", tid)
				continue
			}
			if str(item, "module_id") != id {
				report.Error(check, id, "terminal item %q belongs to %s", tid, str(item, "module_id"))
			}
		}
	}
	// Module DAG cycle check.
	edges := map[string][]string{}
	for _, mod := range cur.Modules {
		edges[str(mod, "id")] = stringsList(mod, "prerequisites")
	}
	detectCycles(edges, check, report)
}

func validateItemGraph(cur *Curriculum, report *Report) {
	check := "item-graph"
	moduleOrder := moduleOrderMap(cur)
	for _, item := range cur.Items {
		id := str(item, "id")
		for _, pid := range stringsList(item, "prerequisites") {
			p, ok := cur.ItemByID[pid]
			if !ok {
				report.Error(check, id, "unknown prerequisite item %q", pid)
				continue
			}
			if moduleOrder[str(p, "module_id")] > moduleOrder[str(item, "module_id")] {
				report.Error(check, id, "prerequisite %q is in a later module", pid)
			}
		}
		for _, nid := range stringsList(item, "next_item_ids") {
			if _, ok := cur.ItemByID[nid]; !ok {
				report.Error(check, id, "unknown next_item_id %q", nid)
			}
			if nid == id {
				report.Error(check, id, "item cannot be next_item_id of itself")
			}
		}
		proof := obj(item, "proof")
		aid := str(proof, "assessment_id")
		if aid != "" {
			if _, ok := cur.AssessmentByID[aid]; !ok {
				report.Error(check, id, "proof.assessment_id %q not found", aid)
			}
		}
		if stringsList(item, "next_item_ids") == nil {
			report.Error(check, id, "next_item_ids must be an array, not null")
		}
		if stringsList(item, "prerequisites") == nil {
			report.Error(check, id, "prerequisites must be an array, not null")
		}
	}
	edges := map[string][]string{}
	for _, item := range cur.Items {
		edges[str(item, "id")] = stringsList(item, "next_item_ids")
	}
	detectCycles(edges, check, report)
}

func validateModuleReachability(cur *Curriculum, report *Report) {
	check := "reachability"
	itemsByModule := map[string][]map[string]any{}
	for _, item := range cur.Items {
		itemsByModule[str(item, "module_id")] = append(itemsByModule[str(item, "module_id")], item)
	}
	for _, mod := range cur.Modules {
		mid := str(mod, "id")
		if len(itemsByModule[mid]) == 0 {
			report.Error(check, mid, "module has no items")
		}
		reachable := map[string]bool{}
		queue := append([]string{}, stringsList(mod, "entry_item_ids")...)
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			if reachable[id] {
				continue
			}
			reachable[id] = true
			item, ok := cur.ItemByID[id]
			if !ok {
				continue
			}
			for _, next := range stringsList(item, "next_item_ids") {
				if !reachable[next] {
					queue = append(queue, next)
				}
			}
		}
		for _, item := range itemsByModule[mid] {
			if !reachable[str(item, "id")] {
				report.Error(check, str(item, "id"), "not reachable from module entry items %s", formatList(stringsList(mod, "entry_item_ids")))
			}
		}
	}
}

func moduleOrderMap(cur *Curriculum) map[string]int {
	out := map[string]int{}
	for _, mod := range cur.Modules {
		out[str(mod, "id")] = num(mod, "number")
	}
	return out
}

func detectCycles(edges map[string][]string, check string, report *Report) {
	visited := map[string]bool{}
	stack := map[string]bool{}
	path := []string{}
	var dfs func(string)
	dfs = func(id string) {
		if stack[id] {
			report.Error(check, id, "cycle detected: %s -> %s", fmt.Sprint(path), id)
			return
		}
		if visited[id] {
			return
		}
		visited[id] = true
		stack[id] = true
		path = append(path, id)
		for _, next := range edges[id] {
			if _, exists := edges[next]; exists {
				dfs(next)
			}
		}
		path = path[:len(path)-1]
		stack[id] = false
	}
	keys := make([]string, 0, len(edges))
	for k := range edges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		dfs(k)
	}
}
