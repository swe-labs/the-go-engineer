# Scripts Cleanup Notes

## Files to be Removed

The following files have been consolidated into `maintainer-scripts/` subdirectory and should be removed from this directory:

- `github-automation.py` → `maintainer-scripts/github-automation.py`
- `github-config.example.json` → `maintainer-scripts/github-config.example.json`
- `github-config.example.yaml` → `maintainer-scripts/github-config.example.yaml`
- `GITHUB_AUTOMATION_GUIDE.md` → `maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md`
- `setup-github-issues.py` → `maintainer-scripts/setup-github-issues.py`
- `setup-github-issues.sh` → `maintainer-scripts/setup-github-issues.sh`

## Rationale

These automation scripts were moved to `maintainer-scripts/` to clearly separate:
- **Curriculum & Go Tools** (scripts/): Validation scripts for learners
- **Maintainer Tools** (scripts/maintainer-scripts/): GitHub automation for contributors

This prevents learner confusion when exploring the codebase.

## Keep In This Directory

- `README.md` - Points to maintainer-scripts and validate_curriculum.go
- `validate_curriculum.go` - Go-based curriculum validation tool (for learners)

## Cleanup Command

```bash
cd scripts
rm -f github-automation.py
rm -f github-config.example.json
rm -f github-config.example.yaml
rm -f GITHUB_AUTOMATION_GUIDE.md
rm -f setup-github-issues.py
rm -f setup-github-issues.sh
```
