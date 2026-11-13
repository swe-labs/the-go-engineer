package main

// Bundle loaders centralize how the validator reads the split manifest
// layout: path.core.json / path.electives.json are thin indexes (header +
// module registry), and every item/project/assessment record lives in the
// per-module manifest.json next to the module's README.md. Gates that need
// per-record raw shapes iterate moduleManifestPaths() with their own local
// structs; gates that need whole-path views use the assembled bundles.

import (
	"path/filepath"
	"strings"
)

var (
	coreBundleCache     *CoreBundle
	electiveBundleCache *ElectiveBundle
	manifestPathsCache  []string
)

// moduleRelDir resolves a module's directory relative to curriculumDir,
// mirroring resolveModuleDir's two-candidate logic (some module paths carry
// a legacy "modules/" prefix).
func moduleRelDir(modPath string) string {
	rel := filepath.ToSlash(modPath)
	if strings.HasPrefix(rel, "modules/") {
		return rel
	}
	return filepath.ToSlash(filepath.Join("modules", rel))
}

// moduleManifestPath returns the curriculumDir-relative manifest path for a
// module. Electives live outside modules/ (their index path is already
// curriculumDir-relative, e.g. electives/advanced-electives).
func moduleManifestPath(modPath string) string {
	if strings.HasPrefix(filepath.ToSlash(modPath), "electives/") {
		return filepath.ToSlash(filepath.Join(modPath, "manifest.json"))
	}
	return filepath.ToSlash(filepath.Join(moduleRelDir(modPath), "manifest.json"))
}

// moduleManifestPaths returns every per-module manifest path in canonical
// module order (core index order, then electives).
func moduleManifestPaths() []string {
	if manifestPathsCache == nil {
		var paths []string
		for _, m := range loadCoreBundle().Modules {
			paths = append(paths, moduleManifestPath(m.Path))
		}
		for _, m := range loadElectiveBundle().Modules {
			paths = append(paths, moduleManifestPath(m.Path))
		}
		manifestPathsCache = paths
	}
	return manifestPathsCache
}

// loadCoreBundle returns the assembled core-path bundle: header + module
// registry from the index, items collected from the module manifests in
// canonical order. Treat the result as read-only.
func loadCoreBundle() CoreBundle {
	if coreBundleCache == nil {
		var index CoreBundle
		readJSON("path.core.json", &index)
		core := index
		core.Items = nil
		for _, m := range index.Modules {
			var mm struct {
				Items []Item `json:"items"`
			}
			readJSON(moduleManifestPath(m.Path), &mm)
			core.Items = append(core.Items, mm.Items...)
		}
		coreBundleCache = &core
	}
	return *coreBundleCache
}

// loadElectiveBundle returns the assembled electives-path bundle.
func loadElectiveBundle() ElectiveBundle {
	if electiveBundleCache == nil {
		var index ElectiveBundle
		readJSON("path.electives.json", &index)
		electives := index
		electives.Items = nil
		for _, m := range index.Modules {
			var mm struct {
				Items []Item `json:"items"`
			}
			readJSON(moduleManifestPath(m.Path), &mm)
			electives.Items = append(electives.Items, mm.Items...)
		}
		electiveBundleCache = &electives
	}
	return *electiveBundleCache
}
