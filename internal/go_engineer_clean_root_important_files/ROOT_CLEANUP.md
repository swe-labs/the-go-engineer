# Root cleanup manifest

## Kept at root

- `README.md`
- `go.mod`
- `go.sum`
- `.gitignore`
- `LICENSE`
- `Makefile`
- `.github/workflows/`

## Moved out of root

Legacy v2 documents were moved to `docs/legacy-v2/`.

The v2 registry was moved to `metadata/legacy/curriculum.v2.json`.

## Deleted/excluded

- `.env` — excluded because environment files and tokens must never be committed.
- `AGENTS.md` — excluded because reusable maintainer/assistant behavior now belongs in the packaged Skill and neutral `tools/authoring/` workflows.
- Duplicate root-level v2 docs — no longer active root contracts for the v3 architecture.

## Active root policy

The root should remain small and should only contain project entry points, release-critical files, and standard repository configuration.
