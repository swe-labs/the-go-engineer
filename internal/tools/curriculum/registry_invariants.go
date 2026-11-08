package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// allowedDocumentationModes is the closed set of documentation modes carried by
// readme_contract.documentation_mode. The value is authoring intent, not a free
// text note: a typo silently changes how a document is treated by tooling, so
// every value must be declared here.
var allowedDocumentationModes = map[string]bool{
	"assessment":          true,
	"companion-overview":  true,
	"elective-overview":   true,
	"explanation":         true,
	"flagship":            true,
	"flagship-checkpoint": true,
	"how-to":              true,
	"module-overview":     true,
	"project":             true,
	"reference":           true,
	"repository-root":     true,
	"tutorial":            true,
}

// allowedContentFlags is the closed set of content_contract keys that may be
// referenced by a "when:<flag>" file policy token.
var allowedContentFlags = map[string]bool{
	"readme_required":             true,
	"runnable_required":           true,
	"tests_required":              true,
	"visual_model_required":       true,
	"machine_view_required":       true,
	"common_mistakes_required":    true,
	"production_notes_required":   true,
	"review_questions_required":   true,
	"portfolio_artifact_required": true,
}

// isKebabToken reports whether s is a non-empty lowercase kebab-case token. It
// is used for registry keys and role labels that must remain machine-stable.
func isKebabToken(s string) bool {
	if s == "" || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			continue
		}
		return false
	}
	return true
}

func validCodePolicy(p string) bool {
	switch p {
	case "", "optional", "required", "forbidden":
		return true
	}
	return false
}

func validSeverityTier(t string) bool {
	switch t {
	case "error", "warning", "skip":
		return true
	}
	return false
}

// validateRegistryFields enforces the modeled contract-registry invariants:
// role labels, heading allowlists, subtype modifiers, aliases, file policies,
// severity policy completeness, and the authoritative locked-module pointer.
func validateRegistryFields(reg *Registry, exp *ExpectedBundle) bool {
	ok := true

	appliesToOwners := make(map[string]string)
	for _, c := range reg.Contracts {
		if c.ID == "" {
			errorf("readme.contracts.json: contract with an empty id")
			ok = false
			continue
		}
		if !isKebabToken(c.AppliesTo) {
			errorf("readme.contracts.json: %s has invalid applies_to '%s' (expected a lowercase kebab-case role)", c.ID, c.AppliesTo)
			ok = false
		} else if prev, dup := appliesToOwners[c.AppliesTo]; dup {
			errorf("readme.contracts.json: applies_to '%s' is declared by both '%s' and '%s'", c.AppliesTo, prev, c.ID)
			ok = false
		} else {
			appliesToOwners[c.AppliesTo] = c.ID
		}

		if c.Shape == "" {
			errorf("readme.contracts.json: %s has an empty shape", c.ID)
			ok = false
		}
		if len(c.RequiredHeadings) == 0 {
			errorf("readme.contracts.json: %s has no required headings", c.ID)
			ok = false
		}
		if !validCodePolicy(c.CodePolicy) {
			errorf("readme.contracts.json: %s has unknown code_policy '%s'", c.ID, c.CodePolicy)
			ok = false
		}
		if c.Selector.Doc == "" && len(c.Selector.Type) == 0 {
			errorf("readme.contracts.json: %s selector must declare doc or type", c.ID)
			ok = false
		}
		if c.Selector.Doc != "" && len(c.Selector.Type) > 0 {
			errorf("readme.contracts.json: %s selector must not mix doc and type selectors", c.ID)
			ok = false
		}

		required := make(map[string]bool, len(c.RequiredHeadings))
		for _, h := range c.RequiredHeadings {
			n := normalizeHeading(h)
			if n == "" {
				errorf("readme.contracts.json: %s has an empty required heading", c.ID)
				ok = false
				continue
			}
			if required[n] {
				errorf("readme.contracts.json: %s declares required heading '%s' more than once", c.ID, h)
				ok = false
			}
			required[n] = true
		}
		optional := make(map[string]bool, len(c.OptionalHeadings))
		for _, h := range c.OptionalHeadings {
			n := normalizeHeading(h)
			if n == "" {
				errorf("readme.contracts.json: %s has an empty optional heading", c.ID)
				ok = false
				continue
			}
			if optional[n] {
				errorf("readme.contracts.json: %s declares optional heading '%s' more than once", c.ID, h)
				ok = false
				continue
			}
			if required[n] {
				errorf("readme.contracts.json: %s declares heading '%s' as both required and optional", c.ID, h)
				ok = false
			}
			optional[n] = true
		}

		for variant, target := range c.HeadingAliases {
			nv := normalizeHeading(variant)
			nt := normalizeHeading(target)
			if nv == "" || nt == "" {
				errorf("readme.contracts.json: %s has an empty heading alias", c.ID)
				ok = false
				continue
			}
			if !required[nt] && !optional[nt] {
				errorf("readme.contracts.json: %s alias '%s' targets '%s', which is not a declared heading", c.ID, variant, target)
				ok = false
			}
		}

		for subtype, mod := range c.SubtypeModifiers {
			if strings.TrimSpace(subtype) == "" {
				errorf("readme.contracts.json: %s has a subtype modifier with an empty name", c.ID)
				ok = false
			}
			if strings.TrimSpace(mod.Description) == "" {
				errorf("readme.contracts.json: %s subtype '%s' has no description", c.ID, subtype)
				ok = false
			}
			if !validCodePolicy(mod.CodePolicy) {
				errorf("readme.contracts.json: %s subtype '%s' has unknown code_policy '%s'", c.ID, subtype, mod.CodePolicy)
				ok = false
			}
			modRequired := make(map[string]bool, len(mod.RequiredHeadings))
			for _, h := range mod.RequiredHeadings {
				n := normalizeHeading(h)
				if n == "" {
					errorf("readme.contracts.json: %s subtype '%s' has an empty required heading", c.ID, subtype)
					ok = false
					continue
				}
				if !required[n] && !optional[n] {
					errorf("readme.contracts.json: %s subtype '%s' requires heading '%s', which is not declared by the family", c.ID, subtype, h)
					ok = false
				}
				modRequired[n] = true
			}
			for _, h := range mod.ForbiddenHeadings {
				n := normalizeHeading(h)
				if n == "" {
					errorf("readme.contracts.json: %s subtype '%s' has an empty forbidden heading", c.ID, subtype)
					ok = false
					continue
				}
				if !required[n] && !optional[n] {
					errorf("readme.contracts.json: %s subtype '%s' forbids heading '%s', which is not declared by the family", c.ID, subtype, h)
					ok = false
				}
				if modRequired[n] {
					errorf("readme.contracts.json: %s subtype '%s' both requires and forbids heading '%s'", c.ID, subtype, h)
					ok = false
				}
			}
		}

		for key, token := range c.FilePolicy {
			switch key {
			case "readme", "main", "test", "starter", "solution":
			default:
				errorf("readme.contracts.json: %s has unknown file_policy key '%s'", c.ID, key)
				ok = false
				continue
			}
			switch {
			case token == "required" || token == "optional":
			case strings.HasPrefix(token, "when:"):
				flag := strings.TrimPrefix(token, "when:")
				if !allowedContentFlags[flag] {
					errorf("readme.contracts.json: %s file_policy[%s] references unknown content flag '%s'", c.ID, key, flag)
					ok = false
				}
			default:
				errorf("readme.contracts.json: %s file_policy[%s] has unknown token '%s'", c.ID, key, token)
				ok = false
			}
		}
	}

	ok = validateSeverityPolicy(reg, exp) && ok
	return ok
}

// validateSeverityPolicy checks that the severity policy is complete, uses
// known tiers, is actually differentiated across statuses, and that every
// always-error token is documented in the validation rules.
func validateSeverityPolicy(reg *Registry, exp *ExpectedBundle) bool {
	ok := true
	sp := reg.SeverityPolicy

	if strings.TrimSpace(sp.Description) == "" {
		errorf("readme.contracts.json: severity_policy.description must not be empty")
		ok = false
	}
	switch sp.Enforcement {
	case "enforce", "advisory":
	default:
		errorf("readme.contracts.json: severity_policy.content_shape_enforcement has unknown value '%s'", sp.Enforcement)
		ok = false
	}
	if !validSeverityTier(sp.DefaultContentShape) {
		errorf("readme.contracts.json: severity_policy.default_content_shape has unknown tier '%s'", sp.DefaultContentShape)
		ok = false
	}
	if strings.TrimSpace(sp.RolloutNote) == "" {
		errorf("readme.contracts.json: severity_policy.rollout_note must not be empty")
		ok = false
	}
	if strings.TrimSpace(sp.LockedModulesNote) == "" {
		errorf("readme.contracts.json: severity_policy.locked_modules_note must not be empty")
		ok = false
	}
	if strings.TrimSpace(sp.LockedModulesRef) == "" {
		errorf("readme.contracts.json: severity_policy.locked_modules_ref must not be empty")
		ok = false
	} else if err := validateLockedModulesRef(sp.LockedModulesRef); err != nil {
		errorf("readme.contracts.json: severity_policy.locked_modules_ref %q is invalid: %v", sp.LockedModulesRef, err)
		ok = false
	}
	if len(exp.LockedModules) == 0 {
		errorf("expected.json: locked_modules must not be empty")
		ok = false
	}

	if len(sp.ByStatus) == 0 {
		errorf("readme.contracts.json: severity_policy.by_status must not be empty")
		ok = false
	}
	tiers := make(map[string]bool)
	for status, tier := range sp.ByStatus {
		if !isKebabToken(status) {
			errorf("readme.contracts.json: severity_policy.by_status has invalid status key '%s'", status)
			ok = false
		}
		if !validSeverityTier(tier) {
			errorf("readme.contracts.json: severity_policy.by_status[%s] has unknown tier '%s'", status, tier)
			ok = false
			continue
		}
		tiers[tier] = true
	}
	// A policy whose statuses all map to the same tier is not tiered at all.
	if len(tiers) < 2 {
		errorf("readme.contracts.json: severity_policy.by_status must use at least two distinct tiers")
		ok = false
	}

	if len(sp.AlwaysError) == 0 {
		errorf("readme.contracts.json: severity_policy.always_error must not be empty")
		ok = false
	}
	seen := make(map[string]bool, len(sp.AlwaysError))
	rules := strings.Join(reg.ValidationRules, "\n")
	for _, token := range sp.AlwaysError {
		if !isKebabToken(token) {
			errorf("readme.contracts.json: severity_policy.always_error token '%s' is not a lowercase kebab-case token", token)
			ok = false
			continue
		}
		if seen[token] {
			errorf("readme.contracts.json: severity_policy.always_error token '%s' is duplicated", token)
			ok = false
		}
		seen[token] = true
		if !strings.Contains(rules, token) {
			errorf("readme.contracts.json: severity_policy.always_error token '%s' is not documented in validation_rules", token)
			ok = false
		}
	}
	return ok
}

// validateLockedModulesRef resolves a file#/json-pointer reference and verifies
// that the target file exists and the referenced top-level key is present. Only
// top-level pointers are supported because that is the shape the registry uses.
func validateLockedModulesRef(ref string) error {
	parts := strings.SplitN(ref, "#", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || !strings.HasPrefix(parts[1], "/") {
		return fmt.Errorf("expected <path>#/<top-level-key>")
	}
	file := strings.TrimSpace(parts[0])
	pointer := strings.TrimPrefix(strings.TrimSpace(parts[1]), "/")
	if pointer == "" || strings.Contains(pointer, "/") {
		return fmt.Errorf("only top-level JSON pointers are supported")
	}
	data, err := os.ReadFile(resolvePath(file))
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", file, err)
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("cannot parse %s: %w", file, err)
	}
	raw, exists := root[pointer]
	if !exists {
		return fmt.Errorf("key '%s' not found in %s", pointer, file)
	}
	var list []string
	if err := json.Unmarshal(raw, &list); err != nil {
		return fmt.Errorf("key '%s' in %s is not a string array", pointer, file)
	}
	if len(list) == 0 {
		return fmt.Errorf("key '%s' in %s is empty", pointer, file)
	}
	return nil
}

// validateReadmeContractFields checks the per-document readme_contract pointer
// shared by items and modules.
func validateReadmeContractFields(reg *Registry, owner string, rc *ReadmeContract) bool {
	ok := true
	if rc == nil {
		errorf("%s: missing 'readme_contract'", owner)
		return false
	}
	if rc.ContractID == "" || !reg.definedContractID(rc.ContractID) {
		errorf("%s: readme_contract.contract_id '%s' is not defined in readme.contracts.json", owner, rc.ContractID)
		ok = false
	}
	if strings.TrimSpace(rc.DocumentationMode) == "" {
		errorf("%s: readme_contract.documentation_mode must not be empty", owner)
		ok = false
	} else if !allowedDocumentationModes[rc.DocumentationMode] {
		errorf("%s: readme_contract.documentation_mode '%s' is not a known mode", owner, rc.DocumentationMode)
		ok = false
	}
	return ok
}
