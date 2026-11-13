# Scripts & Tooling

This directory contains utility scripts for validating, generating, and maintaining The Go Engineer curriculum and toolchains.

---

## 1. Curriculum Validation

### `scripts/validate_curriculum.go`

Delegates to the canonical Go graph and schema engine in `internal/tools/curriculum`. Supports argument forwarding to run either the full verification suite or any individual gate.

```bash
# Run all quality gates (default)
go run ./scripts/validate_curriculum.go

# Run a specific quality gate
go run ./scripts/validate_curriculum.go validate-schema
go run ./scripts/validate_curriculum.go validate-graph
go run ./scripts/validate_curriculum.go validate-zero-magic
go run ./scripts/validate_curriculum.go validate-markdown-links
go run ./scripts/validate_curriculum.go validate-sese-navigation
```

### The Authoritative Quality Gates

| Command | Target Gate | Description |
| :--- | :--- | :--- |
| `validate-schema` | Schema 3.1.0 Integrity | Validates all 27 JSON manifests (8 top-level + 19 per-module `manifest.json` files), required properties, canonical alphabetical key ordering, and anti-overengineering contracts. |
| `validate-graph` | DAG Acyclicity | Detects circular dependencies and validates all module entry/terminal paths. |
| `validate-lessons` | Lesson Quality | Enforces zero machine role comments (`//go:role`) and validates `NEXT UP` footers. |
| `validate-prerequisite-closure` | Pedagogical Progression | Verifies no hidden prerequisite leaks or future-concept violations across module boundaries. |
| `validate-project-binding` | Project Placement | Verifies that all 29 projects bind to valid curriculum modules, anchors, deliverables, and rubrics. |
| `validate-no-placeholder-zm` | Zero Placeholders | Strictly scans all `zero_magic` blocks to ensure zero boilerplate or placeholder phrases. |
| `validate-assessment-completeness` | Assessment Rigor | Verifies all 48 assessments include non-generic, weighted criteria and passing scores. |
| `validate-unique-tags` | Tag Taxonomy | Enforces tag consistency across all core and elective items. |
| `validate-crossrefs` | Cross-Reference Graph | Validates bidirectional integrity across 2,228 cross-reference edges. |
| `validate-concept-ownership` | Concept Ownership | Ensures all 69 canonical concepts map to unambiguous owners without conceptual leaks. |
| `validate-cognitive-load` | Cognitive Pacing | Verifies cognitive density ratings and recommended break pacing. |
| `validate-operational-failures` | Failure Taxonomy | Verifies 100% mathematical symmetry across all 13 operational failure categories. |
| `validate-failure-engineering` | Operational Traps | Ensures critical production modules (08–15, 18) include concrete failure examples. |
| `validate-zero-magic` | Exemplar Certification | Validates that all 356 golden lessons strictly adhere to the zero-magic authoring spec. |
| `validate-physical-artifacts` | README Artifact Parity | Verifies physical existence and minimum byte sizes of all 447 instructional READMEs on disk. |
| `validate-assessment-artifacts` | Assessment Suites | Verifies complete 4-file assessment suites across all checkpoints and 14 Opslane milestones. |
| `validate-starter-parity` | Code AST Parity | Compiles and validates Go AST correctness across all 27 project starters and solutions. |
| `validate-markdown-links` | Link Crawler | Crawls all Markdown files, extracts relative links/images, verifies physical existence on disk. |
| `validate-inline-paths` | Link Crawler | Crawls all Markdown files, extracts relative links/images, verifies physical existence on disk. |
| `validate-sese-navigation` | SESE Flow | Enforces front door (`Start here` -> Lesson 01), prohibits escape hatches, verifies assessment exit gates. |
| `validate-project-scripts` | Script Parity | Verifies cross-platform parity (`verify.sh` and `verify.ps1`) for script-backed projects. |
| `validate-no-emojis` | Professional Prose | Scans all Markdown files, enforcing zero non-code emojis across prose, headers, lists, and tables. |
| `validate-anchor-quotes` | Pioneer Anchors | Enforces authoritative computing pioneer and Go designer anchor quotes with valid attribution at the top of lessons in locked modules. |
| `validate-diagram-palette` | Diagram Styling & Syntax | Enforces Dark Engineering Palette classes/color tokens on flowcharts and verifies matching block delimiters (`subgraph`/`loop`/`alt`/`opt` vs `end`). |
| `validate-readme-contracts` | Documentation Contract | Resolves each item's document contract (family + subtype modifiers) from `readme.contracts.json` and verifies required headings, sequential heading order, code policy, and file layout. Structural defects always block; content-shape defects in locked modules block, and other items follow `severity_policy.by_status` (draft/scaffolded skip, review/reviewed warn, ready-and-later block). `--strict` promotes every remaining advisory to an error for the release gate. |
| `validate-anti-slop` | Anti-Slop Prose Invariant | Scans all Markdown prose across the curriculum for banned AI filler, lexical tells, unearned marketing superlatives, and empty transition phrases. |

---

## 2. Curriculum Catalog Generator

The authoritative generator for [`curriculum/README.md`](../curriculum/README.md) is built directly into the Go curriculum toolchain (`internal/tools/curriculum/catalog.go`). It computes quick stats, module breakdown tables, workload pacing models, project lists, assessment registries, and cross-reference distributions directly from the canonical JSON manifests without requiring external runtimes like Node.js.

```bash
# Generate curriculum/README.md via toolchain
go run ./internal/tools/curriculum generate-catalog

# Or via the script forwarder
go run ./scripts/validate_curriculum.go generate-catalog
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

# 5. Run full curriculum verification suite
go run ./scripts/validate_curriculum.go
```
