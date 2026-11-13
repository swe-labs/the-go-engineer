package main

// Bundle loaders centralize how the validator reads the two path manifests.
// Every gate used to readJSON them independently; the loader caches the
// parsed bundles so the per-module manifest layout lands in exactly one
// place instead of ~20 call sites.

var (
	coreBundleCache     *CoreBundle
	electiveBundleCache *ElectiveBundle
)

// loadCoreBundle returns the authoritative core-path bundle. Treat the
// returned bundle as read-only: gates must not mutate shared items.
func loadCoreBundle() CoreBundle {
	if coreBundleCache == nil {
		var core CoreBundle
		readJSON("path.core.json", &core)
		coreBundleCache = &core
	}
	return *coreBundleCache
}

// loadElectiveBundle returns the authoritative electives-path bundle.
func loadElectiveBundle() ElectiveBundle {
	if electiveBundleCache == nil {
		var electives ElectiveBundle
		readJSON("path.electives.json", &electives)
		electiveBundleCache = &electives
	}
	return *electiveBundleCache
}
