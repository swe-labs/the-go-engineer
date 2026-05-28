# Go Engineer curriculum automation.
#
# This Makefile is intentionally small at the root and delegates real work to
# tools/. It is safe for local use and CI. Release targets are strict by design.

SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

GO ?= go
PYTHON ?= python3

METADATA_DIR ?= metadata
CURRICULUM_DIR ?= curriculum
TOOLS_DIR ?= tools
DOCS_DIR ?= docs
DIST_DIR ?= dist

CURRICULUM_VALIDATOR := ./$(TOOLS_DIR)/validate/curriculum
REPOSITORY_VALIDATORS := $(TOOLS_DIR)/validate/repository

STRICT_REPOSITORY ?= 0

.PHONY: help
help:
	@echo "Go Engineer curriculum commands"
	@echo ""
	@echo "Setup and health:"
	@echo "  make doctor                 Check required local tools and repository layout"
	@echo "  make ensure-dirs            Create required empty top-level directories"
	@echo "  make check-root             Check root cleanliness and forbidden files/folders"
	@echo "  make check-secrets          Check that secret-like files are not tracked"
	@echo ""
	@echo "Validation:"
	@echo "  make validate-metadata      Validate metadata graph, contracts, and references"
	@echo "  make validate-repository    Validate learner-facing curriculum files"
	@echo "  make validate-release       Run the full release-quality validation gate"
	@echo "  make validate-all           Alias for validate-release"
	@echo "  make validate-python        Compile all Python tooling"
	@echo "  make validate-docs          Check required maintainer docs exist"
	@echo ""
	@echo "Go tooling:"
	@echo "  make test-tools             Run Go tests for validator tooling"
	@echo "  make test                   Run all Go tests"
	@echo "  make vet                    Run go vet"
	@echo "  make fmt-tools              Format Go validator tooling"
	@echo "  make fmt-check              Check Go formatting without changing files"
	@echo "  make tidy                   Run go mod tidy and ensure go.mod/go.sum stay clean"
	@echo ""
	@echo "Generation:"
	@echo "  make generate-lessons       Generate lesson files from metadata"
	@echo "  make generate-labs          Generate lab files from metadata"
	@echo "  make generate-modules       Generate module README files"
	@echo "  make generate-projects      Generate project README files"
	@echo "  make generate-assessments   Generate assessment files"
	@echo "  make generate-shared-docs   Generate shared learner docs"
	@echo "  make generate-all           Run all generators"
	@echo ""
	@echo "Audit:"
	@echo "  make audit-curriculum       Audit metadata and repository coherence"
	@echo "  make audit-modules          Audit module quality"
	@echo "  make audit-lessons          Audit lesson quality"
	@echo "  make audit-labs             Audit lab quality"
	@echo "  make audit-projects         Audit project quality"
	@echo "  make audit-assessments      Audit assessment quality"
	@echo "  make audit-completion       Produce completion report"
	@echo "  make audit-all              Run all available audits"
	@echo ""
	@echo "Migration:"
	@echo "  make migration-report       Generate v2-to-v3 migration report"
	@echo "  make detect-unmapped-v2     Detect unmapped legacy v2 items"
	@echo ""
	@echo "Release:"
	@echo "  make release-artifacts      Generate dist/ artifacts"
	@echo "  make clean-dist             Remove generated dist artifacts"
	@echo "  make clean                  Remove local generated caches and dist artifacts"

.PHONY: ensure-dirs
ensure-dirs:
	@mkdir -p \
		$(METADATA_DIR) \
		$(CURRICULUM_DIR) \
		$(TOOLS_DIR)/validate/curriculum \
		$(TOOLS_DIR)/validate/repository \
		$(TOOLS_DIR)/generate \
		$(TOOLS_DIR)/audit \
		$(TOOLS_DIR)/migrate \
		$(TOOLS_DIR)/authoring \
		$(DOCS_DIR) \
		$(DIST_DIR)
	@touch $(DIST_DIR)/.gitkeep

.PHONY: doctor
doctor: ensure-dirs
	@command -v $(GO) >/dev/null || { echo "missing Go: $(GO)"; exit 1; }
	@command -v $(PYTHON) >/dev/null || { echo "missing Python: $(PYTHON)"; exit 1; }
	@test -f go.mod || { echo "missing go.mod"; exit 1; }
	@test -f Makefile || { echo "missing Makefile"; exit 1; }
	@test -d $(METADATA_DIR) || { echo "missing $(METADATA_DIR)/"; exit 1; }
	@test -d $(CURRICULUM_DIR) || { echo "missing $(CURRICULUM_DIR)/"; exit 1; }
	@test -d $(TOOLS_DIR) || { echo "missing $(TOOLS_DIR)/"; exit 1; }
	@test -d $(DOCS_DIR) || { echo "missing $(DOCS_DIR)/"; exit 1; }
	@test -d $(DIST_DIR) || { echo "missing $(DIST_DIR)/"; exit 1; }
	@echo "doctor passed"

.PHONY: check-root
check-root:
	@if [ -e AGENTS.md ]; then \
		echo "AGENTS.md should not exist at root; use tools/authoring/ and the packaged Skill instead"; \
		exit 1; \
	fi
	@if [ -f .env ]; then \
		echo ".env exists at root. Keep local secrets untracked and outside committed artifacts."; \
		exit 1; \
	fi
	@if find . -type d \( -name codex -o -name ai -o -name agent -o -name bot -o -name llm \) \
		-not -path "./.git/*" | grep -q .; then \
		echo "forbidden tool-branded folder found; use neutral responsibility names"; \
		find . -type d \( -name codex -o -name ai -o -name agent -o -name bot -o -name llm \) -not -path "./.git/*"; \
		exit 1; \
	fi
	@echo "root check passed"

.PHONY: check-secrets
check-secrets:
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		if git ls-files | grep -E '(^|/)\.env(\.|$$)' | grep -v '\.env\.example' >/dev/null; then \
			echo "tracked .env file detected"; \
			git ls-files | grep -E '(^|/)\.env(\.|$$)' | grep -v '\.env\.example'; \
			exit 1; \
		fi; \
		if git grep -nE 'ghp_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]{20,}|AKIA[0-9A-Z]{16}' -- . ':!metadata/legacy/curriculum.v2.json' >/tmp/go-engineer-secret-scan.txt 2>/dev/null; then \
			cat /tmp/go-engineer-secret-scan.txt; \
			echo "possible committed secret detected"; \
			exit 1; \
		fi; \
	fi
	@echo "secret check passed"

.PHONY: validate-metadata
validate-metadata: check-root check-secrets
	$(GO) run $(CURRICULUM_VALIDATOR) validate-metadata

PYTHON_STRICT_FLAG=$(if $(filter 1,$(STRICT_REPOSITORY)),,--no-strict)

.PHONY: validate-repository
validate-repository: validate-metadata
ifeq ($(STRICT_REPOSITORY),1)
	$(GO) run $(CURRICULUM_VALIDATOR) validate-repository --strict-repository
else
	$(GO) run $(CURRICULUM_VALIDATOR) validate-repository
endif
	@if [ -f $(REPOSITORY_VALIDATORS)/validate_repository.py ]; then \
		$(PYTHON) $(REPOSITORY_VALIDATORS)/validate_repository.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR); \
	fi
	@if [ -f $(REPOSITORY_VALIDATORS)/validate_readmes.py ]; then \
		$(PYTHON) $(REPOSITORY_VALIDATORS)/validate_readmes.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR) $(PYTHON_STRICT_FLAG); \
	fi
	@if [ -f $(REPOSITORY_VALIDATORS)/validate_code.py ]; then \
		$(PYTHON) $(REPOSITORY_VALIDATORS)/validate_code.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR) $(PYTHON_STRICT_FLAG); \
	fi
	@if [ -f $(REPOSITORY_VALIDATORS)/validate_assets.py ]; then \
		$(PYTHON) $(REPOSITORY_VALIDATORS)/validate_assets.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR) $(PYTHON_STRICT_FLAG); \
	fi

.PHONY: validate-release
validate-release: validate-python validate-docs fmt-check test-tools
	+$(MAKE) validate-repository STRICT_REPOSITORY=1
	$(GO) run $(CURRICULUM_VALIDATOR) validate-all --strict-repository

.PHONY: validate-all
validate-all: validate-release

.PHONY: validate-python
validate-python:
	@if find $(TOOLS_DIR) -name "*.py" | grep -q .; then \
		find $(TOOLS_DIR) -name "*.py" -print0 | xargs -0 $(PYTHON) -m py_compile; \
	fi
	@echo "Python tooling validation passed"

.PHONY: validate-docs
validate-docs:
	@test -f $(DOCS_DIR)/README.md
	@test -f $(DOCS_DIR)/architecture.md
	@test -f $(DOCS_DIR)/metadata-contract.md
	@test -f $(DOCS_DIR)/content-quality-standard.md
	@test -f $(DOCS_DIR)/validation-backbone.md
	@test -f $(DOCS_DIR)/authoring-guide.md
	@test -f $(DOCS_DIR)/migration-guide.md
	@test -f $(DOCS_DIR)/release-process.md
	@test -f $(DOCS_DIR)/contributor-guide.md
	@test -f $(DOCS_DIR)/review-process.md
	@if grep -RniE '\bcodex\b|ai-branded|agent-branded' $(DOCS_DIR); then \
		echo "forbidden tool-branded wording found in docs"; \
		exit 1; \
	fi
	@echo "docs validation passed"

.PHONY: test-tools
test-tools:
	$(GO) test ./$(TOOLS_DIR)/validate/curriculum/...

.PHONY: test
test:
	$(GO) test ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: fmt-tools
fmt-tools:
	@if [ -d $(TOOLS_DIR)/validate/curriculum ]; then \
		gofmt -w $(TOOLS_DIR)/validate/curriculum; \
	fi

.PHONY: fmt-check
fmt-check:
	@unformatted="$$(gofmt -l $$(find . -name '*.go' -not -path './.git/*' -not -path './dist/*'))"; \
	if [ -n "$$unformatted" ]; then \
		echo "$$unformatted"; \
		echo "Go files need gofmt"; \
		exit 1; \
	fi
	@echo "Go formatting check passed"

.PHONY: tidy
tidy:
	$(GO) mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		git diff --exit-code -- go.mod go.sum; \
	fi

.PHONY: generate-lessons
generate-lessons:
	$(PYTHON) $(TOOLS_DIR)/generate/generate_lessons.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: generate-labs
generate-labs:
	$(PYTHON) $(TOOLS_DIR)/generate/generate_labs.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: generate-modules
generate-modules:
	$(PYTHON) $(TOOLS_DIR)/generate/generate_module_readmes.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: generate-projects
generate-projects:
	$(PYTHON) $(TOOLS_DIR)/generate/generate_project_readmes.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: generate-assessments
generate-assessments:
	$(PYTHON) $(TOOLS_DIR)/generate/generate_assessments.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: generate-shared-docs
generate-shared-docs:
	$(PYTHON) $(TOOLS_DIR)/generate/generate_shared_docs.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: generate-all
generate-all: generate-modules generate-lessons generate-labs generate-projects generate-assessments generate-shared-docs

.PHONY: audit-curriculum
audit-curriculum:
	$(PYTHON) $(TOOLS_DIR)/audit/audit_curriculum.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: audit-modules
audit-modules:
	$(PYTHON) $(TOOLS_DIR)/audit/audit_modules.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: audit-lessons
audit-lessons:
	$(PYTHON) $(TOOLS_DIR)/audit/audit_lessons.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: audit-labs
audit-labs:
	$(PYTHON) $(TOOLS_DIR)/audit/audit_labs.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: audit-projects
audit-projects:
	$(PYTHON) $(TOOLS_DIR)/audit/audit_projects.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: audit-assessments
audit-assessments:
	$(PYTHON) $(TOOLS_DIR)/audit/audit_assessments.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR)

.PHONY: audit-completion
audit-completion:
	@mkdir -p $(DIST_DIR)
	$(PYTHON) $(TOOLS_DIR)/audit/audit_completion.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR) --output $(DIST_DIR)/completion-report.json

.PHONY: audit-all
audit-all: audit-curriculum audit-modules audit-lessons audit-labs audit-projects audit-assessments audit-completion

.PHONY: migration-report
migration-report:
	@mkdir -p $(DIST_DIR)
	$(PYTHON) $(TOOLS_DIR)/migrate/migration_report.py --metadata-dir $(METADATA_DIR) --output $(DIST_DIR)/migration-report.json

.PHONY: detect-unmapped-v2
detect-unmapped-v2:
	$(PYTHON) $(TOOLS_DIR)/migrate/detect_unmapped_v2.py --metadata-dir $(METADATA_DIR)

.PHONY: release-artifacts
release-artifacts: validate-release
	@mkdir -p $(DIST_DIR)
	$(PYTHON) $(TOOLS_DIR)/generate/generate_full_snapshot.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR) --output $(DIST_DIR)/curriculum.v3.full.generated.json
	$(PYTHON) $(TOOLS_DIR)/audit/audit_completion.py --metadata-dir $(METADATA_DIR) --curriculum-dir $(CURRICULUM_DIR) --output $(DIST_DIR)/completion-report.json
	$(PYTHON) $(TOOLS_DIR)/migrate/migration_report.py --metadata-dir $(METADATA_DIR) --output $(DIST_DIR)/migration-report.json
	$(GO) run $(CURRICULUM_VALIDATOR) validate-all --strict-repository > $(DIST_DIR)/validation-report.txt
	@printf "# Release Notes\n\nGenerated on $$(date -u +%Y-%m-%dT%H:%M:%SZ).\n\nSee validation, completion, and migration reports in this directory.\n" > $(DIST_DIR)/release-notes.md
	@echo "release artifacts generated in $(DIST_DIR)/"

.PHONY: clean-dist
clean-dist:
	rm -rf $(DIST_DIR)/*
	touch $(DIST_DIR)/.gitkeep

.PHONY: clean
clean: clean-dist
	rm -rf .cache .pytest_cache .ruff_cache .mypy_cache
	find . -type d -name __pycache__ -prune -exec rm -rf {} +
	rm -f coverage.out coverage.html cpu.prof mem.prof trace.out
