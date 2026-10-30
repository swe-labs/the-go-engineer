# Scripts & Tooling

This directory contains utility scripts for validating, generating, and maintaining The Go Engineer curriculum and toolchains.

---

## 1. Curriculum Validation

### `scripts/validate_curriculum.go`

Delegates to the canonical Go graph and schema engine in `internal/tools/curriculum`. Supports argument forwarding to run either the full 17-gate verification suite or any individual gate.

```bash
# Run all 17 quality gates (default)
go run ./scripts/validate_curriculum.go

# Run a specific quality gate
go run ./scripts/validate_curriculum.go validate-schema
go run ./scripts/validate_curriculum.go validate-graph
go run ./scripts/validate_curriculum.go validate-zero-magic
go run ./scripts/validate_curriculum.go validate-physical-artifacts
```

### The 17 Authoritative Quality Gates

| Command | Target Gate | Description |
| :--- | :--- | :--- |
| `validate-schema` | Schema 3.1.0 Integrity | Validates all 10 JSON manifests, required properties, canonical alphabetical key ordering, and anti-overengineering contracts. |
| `validate-graph` | DAG Acyclicity | Detects circular dependencies and validates all module entry/terminal paths. |
| `validate-lessons` | Lesson Quality | Enforces zero machine role comments (`//go:role`) and validates `NEXT UP` footers. |
| `validate-prerequisite-closure` | Pedagogical Progression | Verifies no hidden prerequisite leaks or future-concept violations across module boundaries. |
| `validate-project-binding` | Project Placement | Verifies that all 28 projects bind to valid curriculum modules, anchors, deliverables, and rubrics. |
| `validate-no-placeholder-zm` | Zero Placeholders | Strictly scans all `zero_magic` blocks to ensure zero boilerplate or placeholder phrases. |
| `validate-assessment-completeness` | Assessment Rigor | Verifies all 47 assessments include non-generic, weighted criteria and passing scores. |
| `validate-unique-tags` | Tag Taxonomy | Enforces tag consistency across all core and elective items. |
| `validate-crossrefs` | Cross-Reference Graph | Validates bidirectional integrity across 2,089 cross-reference edges. |
| `validate-concept-ownership` | Concept Ownership | Ensures all 62 canonical concepts map to unambiguous owners without conceptual leaks. |
| `validate-cognitive-load` | Cognitive Pacing | Verifies cognitive density ratings and recommended break pacing. |
| `validate-operational-failures` | Failure Taxonomy | Verifies 100% mathematical symmetry across all 13 operational failure categories. |
| `validate-failure-engineering` | Operational Traps | Ensures critical production modules (08–15, 18) include concrete failure examples. |
| `validate-zero-magic` | Exemplar Certification | Validates that all 92 golden lessons strictly adhere to the zero-magic authoring spec. |
| `validate-physical-artifacts` | README Artifact Parity | Verifies physical existence and minimum byte sizes of all 441 instructional READMEs on disk. |
| `validate-assessment-artifacts` | Assessment Suites | Verifies complete 4-file assessment suites across all checkpoints and 14 Opslane milestones. |
| `validate-starter-parity` | Code AST Parity | Compiles and validates Go AST correctness across all 27 project starters and solutions. |

---

## 2. Curriculum Catalog Generator

### `scripts/generate_curriculum_readme.js`

Authoritative generator for [`curriculum/README.md`](file:///c:/Users/Rasel/Desktop/the-go-engineer/curriculum/README.md). Computes all quick stats, module breakdown tables, workload pacing models, project lists, assessment registries, and cross-reference distributions directly from the canonical JSON manifests.

```bash
node scripts/generate_curriculum_readme.js
```

---

## 3. Toolchain Unit Testing

The curriculum engine in `internal/tools/curriculum` is tested using standard Go table-driven unit tests:

```bash
go test -v ./internal/tools/curriculum/...
```

---

## 4. Full Pre-Flight CI-Equivalent Verification

Run this sequence before committing or opening pull requests:

```bash
# 1. Compile and vet all Go packages
go build ./...
go vet ./...

# 2. Format check (excluding project starters)
unformatted=$(gofmt -l . | grep -v '_starter' || true); test -z "$unformatted" || (echo "$unformatted" && exit 1)

# 3. Dependencies check
go mod tidy
git diff --exit-code -- go.mod go.sum

# 4. Run toolchain and companion frontend test suites
go test -v ./internal/tools/curriculum/...
go test -v ./curriculum/companion-frontend/...

# 5. Run full 17-gate curriculum verification suite
go run ./scripts/validate_curriculum.go
```

---

## 5. Maintainer Automation Scripts

Maintainer tools for repository labels, milestone sync, and issue automation live in `scripts/maintainer-scripts/`.

See [maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md](./maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md).
