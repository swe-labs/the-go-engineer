package main

import "strings"

// Selector determines which documents a contract family applies to. Item-level
// families match on type/subtype; document-level families (root, module,
// flagship, companion) match on a role token carried in "doc".
type Selector struct {
	Doc        string   `json:"doc,omitempty"`
	Type       []string `json:"type,omitempty"`
	Subtype    []string `json:"subtype,omitempty"`
	SubtypeNot []string `json:"subtype_not,omitempty"`
	Aliases    []string `json:"aliases,omitempty"`
}

// SubtypeModifier layers a few extra requirements or forbiddances on top of a
// family for one specific kind of lesson (for example, orientation).
type SubtypeModifier struct {
	Description       string   `json:"description"`
	CodePolicy        string   `json:"code_policy,omitempty"`
	RequiredHeadings  []string `json:"required_headings,omitempty"`
	ForbiddenHeadings []string `json:"forbidden_headings,omitempty"`
}

// Contract is one README document family: a small invariant heading core, an
// optional allowlist, alias normalization, and code/file policies.
type Contract struct {
	ID               string                     `json:"id"`
	AppliesTo        string                     `json:"applies_to"`
	Shape            string                     `json:"shape"`
	Selector         Selector                   `json:"selector"`
	CodePolicy       string                     `json:"code_policy"`
	FilePolicy       map[string]string          `json:"file_policy"`
	RequiredHeadings []string                   `json:"required_headings"`
	OptionalHeadings []string                   `json:"optional_headings"`
	HeadingAliases   map[string]string          `json:"heading_aliases"`
	SubtypeModifiers map[string]SubtypeModifier `json:"subtype_modifiers"`
}

// SeverityPolicy is the tiered severity model shared by the content-shape
// validators.
type SeverityPolicy struct {
	Description         string            `json:"description"`
	DefaultContentShape string            `json:"default_content_shape"`
	Enforcement         string            `json:"content_shape_enforcement"`
	ByStatus            map[string]string `json:"by_status"`
	RolloutNote         string            `json:"rollout_note"`
	AlwaysError         []string          `json:"always_error"`
	LockedModulesRef    string            `json:"locked_modules_ref"`
	LockedModulesNote   string            `json:"locked_modules_note"`
}

// Registry is the in-memory form of curriculum/readme.contracts.json.
type Registry struct {
	Contracts      []Contract     `json:"contracts"`
	SeverityPolicy SeverityPolicy `json:"severity_policy"`
}

var registryCache *Registry

// loadRegistry returns the parsed README contract registry, cached after first load.
func loadRegistry() *Registry {
	if registryCache != nil {
		return registryCache
	}
	var r Registry
	readJSON("readme.contracts.json", &r)
	registryCache = &r
	return registryCache
}

// contractByID returns the contract with the given id, or nil.
func (r *Registry) contractByID(id string) *Contract {
	for i := range r.Contracts {
		if r.Contracts[i].ID == id {
			return &r.Contracts[i]
		}
	}
	return nil
}

// definedContractID reports whether id is a known contract id or a declared
// selector alias of one (for example "elective-module.v3" -> module.v3).
func (r *Registry) definedContractID(id string) bool {
	if id == "" {
		return false
	}
	if r.contractByID(id) != nil {
		return true
	}
	_, ok := r.aliasesOf()[id]
	return ok
}

// aliasesOf returns every selector alias mapped to a contract id.
func (r *Registry) aliasesOf() map[string]string {
	out := make(map[string]string)
	for i := range r.Contracts {
		for _, a := range r.Contracts[i].Selector.Aliases {
			out[a] = r.Contracts[i].ID
		}
	}
	return out
}

// resolveItemContract returns the item-level contract family for an item and
// whether exactly one family matched.
func (r *Registry) resolveItemContract(item Item) (*Contract, bool) {
	var matches []*Contract
	for i := range r.Contracts {
		c := &r.Contracts[i]
		if c.Selector.Doc != "" {
			continue // document-level family, not item-level
		}
		if !c.matchesItem(item) {
			continue
		}
		matches = append(matches, c)
	}
	if len(matches) == 1 {
		return matches[0], true
	}
	return nil, false
}

func (c *Contract) matchesItem(item Item) bool {
	s := c.Selector
	if len(s.Type) == 0 {
		return false
	}
	if !containsFold(s.Type, item.Type) {
		return false
	}
	if len(s.Subtype) > 0 && !containsFold(s.Subtype, item.Subtype) {
		return false
	}
	if len(s.SubtypeNot) > 0 && containsFold(s.SubtypeNot, item.Subtype) {
		return false
	}
	return true
}

func containsFold(list []string, v string) bool {
	for _, s := range list {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}

// -- Severity resolution ---------------------------------------------------

type severity int

const (
	sevSkip severity = iota
	sevWarning
	sevError
)

// contentShapeSeverity resolves the severity of a content-shape issue (heading,
// code policy, file policy) for an item, honoring the mode flags, the rollout
// enforcement switch, and the registry severity policy.
func contentShapeSeverity(item Item) severity {
	if advisoryMode {
		return sevSkip
	}
	if strictMode {
		return sevError
	}
	sp := loadRegistry().SeverityPolicy
	// During rollout ("advisory") content-shape issues never block; only the
	// --strict release gate promotes them. After the advisory report is clean,
	// flip content_shape_enforcement to "enforce" to apply the tiered policy.
	if sp.Enforcement != "enforce" {
		return sevWarning
	}
	if lockedModuleSet()[item.ModuleID] {
		return sevError
	}
	if s, ok := sp.ByStatus[item.Status]; ok {
		switch s {
		case "error":
			return sevError
		case "warning":
			return sevWarning
		case "skip":
			return sevSkip
		}
	}
	switch sp.DefaultContentShape {
	case "error":
		return sevError
	case "skip":
		return sevSkip
	default:
		return sevWarning
	}
}
