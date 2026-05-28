SHELL := /usr/bin/env bash

.PHONY: help
help:
	@echo "Go Engineer curriculum commands"
	@echo ""
	@echo "Validation:"
	@echo "  make validate-metadata      Validate metadata graph, contracts, and references"
	@echo "  make validate-repository    Validate learner-facing files strictly"
	@echo "  make validate-release       Run full release-quality validation"
	@echo "  make test-tools             Run Go tests for validation tooling"
	@echo "  make fmt-tools              Format Go validation tooling"
	@echo ""
	@echo "Artifacts:"
	@echo "  make release-artifacts      Generate dist artifacts when generators exist"
	@echo "  make clean-dist             Remove generated dist artifacts"

.PHONY: validate-metadata
validate-metadata:
	go run ./tools/validate/curriculum validate-metadata

.PHONY: validate-repository
validate-repository:
	go run ./tools/validate/curriculum validate-repository --strict-repository

.PHONY: validate-release
validate-release:
	go run ./tools/validate/curriculum validate-all --strict-repository

.PHONY: test-tools
test-tools:
	go test ./tools/validate/curriculum/...

.PHONY: fmt-tools
fmt-tools:
	gofmt -w tools/validate/curriculum

.PHONY: audit-completion
audit-completion:
	@if [ -f tools/audit/audit_completion.py ]; then \
		python3 tools/audit/audit_completion.py --metadata-dir metadata --curriculum-dir curriculum --output dist/completion-report.json; \
	else \
		echo "tools/audit/audit_completion.py not found"; \
	fi

.PHONY: release-artifacts
release-artifacts:
	mkdir -p dist
	@if [ -f tools/generate/generate_full_snapshot.py ]; then \
		python3 tools/generate/generate_full_snapshot.py --metadata-dir metadata --curriculum-dir curriculum --output dist/curriculum.v3.full.generated.json; \
	else \
		echo "tools/generate/generate_full_snapshot.py not found; skipping snapshot"; \
	fi
	@if [ -f tools/audit/audit_completion.py ]; then \
		python3 tools/audit/audit_completion.py --metadata-dir metadata --curriculum-dir curriculum --output dist/completion-report.json; \
	else \
		echo "tools/audit/audit_completion.py not found; skipping completion report"; \
	fi
	@if [ -f tools/migrate/migration_report.py ]; then \
		python3 tools/migrate/migration_report.py --metadata-dir metadata --output dist/migration-report.json; \
	else \
		echo "tools/migrate/migration_report.py not found; skipping migration report"; \
	fi
	go run ./tools/validate/curriculum validate-all --strict-repository > dist/validation-report.txt

.PHONY: clean-dist
clean-dist:
	rm -rf dist/*
	touch dist/.gitkeep
