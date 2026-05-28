# Validation Changelog

## Quality backbone update

### Replaced

- Replaced legacy directory discovery that only searched for `curriculum/` with path-aware discovery for `metadata/`, `curriculum/`, or root-level JSON metadata.
- Replaced mixed transform/validator scripts with a pure validator package.
- Replaced status-limited lesson checks with strict repository validation that can require every declared README/code/test/asset file.
- Replaced generic placeholder checks that produced false positives in metadata with exact-literal metadata checks and stricter authored-content checks for real README/code files.

### Added

- `validate-metadata` profile for JSON metadata completion.
- `validate-repository --strict-repository` profile for real curriculum implementation completeness.
- JSON report output for CI.
- Robust project validation, including generated project assessments and rubric weight checks.
- Assessment quality checks for evidence, manual review questions, retake policies, passing score, target IDs, and criterion weights.
- README quality gate requiring problem framing, mental model, analogy, production context, Go code/commands, practice, review questions, and `NEXT UP` footer.
- Go code quality checks for `gofmt`, real test functions, verification commands, and placeholder code.
- Asset validation for visual model requirements.

### Metadata fixes applied

- Fixed `core-16-01` to depend on `core-15-19` instead of `elective-24`.
- Added project-specific assessment records for every project assessment ID.
- Normalized project rubrics, project verification blocks, portfolio metadata, and README contracts.
- Expanded module summaries so they are no longer duplicates of learning goals.
- Updated workspace validation commands to use the new Go validator profiles.

### Removed from the shipped tool package

- The old `main(1).go` transformation script.
- The old transform-oriented `assessments.go` behavior.
- The old assumption that project assessments can point to missing IDs.
