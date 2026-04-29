# GitHub Automation Tool

This guide documents the maintainer automation script for repository labels, issues, and setup tasks.

## Scope

The script supports:

- label creation and updates
- bulk issue creation
- dry-run previews
- repository setup checks

Use it for maintainer setup work, not normal lesson authoring.

## The Go Engineer Conventions

Issue titles use:

```text
[TYPE] short description
```

Allowed types:

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

Examples:

```text
[FEAT] add SY.1 sync.Mutex lesson
[FIX] correct SQL injection example in DB.3
[DOCS] align roadmap with v2.1 final status
[TEST] add TE.4 benchmark coverage
[CHORE] update repository labels
```

Every issue should include:

- section ID, if applicable
- lesson IDs, if applicable
- expected files
- validation plan
- risk notes

## Prerequisites

- Python 3.8+
- GitHub token with required repository permissions
- Python packages: `requests`, `pyyaml`

Install dependencies:

```bash
pip install requests pyyaml
```

## Authentication

Set a token:

```bash
export GITHUB_TOKEN=your_token_here
```

Or use a local `.env` file that is not committed.

## Usage

```bash
python github-automation.py --labels-only
python github-automation.py --issues-only
python github-automation.py --all
python github-automation.py --dry-run --all
```

## Configuration

Use JSON or YAML configuration.

Example:

```yaml
owner: rasel9t6
repo: the-go-engineer

labels:
  - name: lesson
    color: "0366D6"
    description: "Lesson or tutorial work"
  - name: bug
    color: "D73A49"
    description: "Incorrect behavior or broken content"

issues:
  - title: "[FEAT] add SY.1 sync.Mutex lesson"
    body: |
      ## Scope
      Section: s07
      Lesson IDs: SY.1

      ## Validation
      - go test ./...
      - go run ./scripts/validate_curriculum.go
    labels:
      - lesson
```

## Safety Rules

- Run `--dry-run` before applying large changes.
- Do not commit tokens or `.env` files.
- Keep generated issue bodies aligned with `CONTRIBUTING.md`.
- Use bracketed issue types consistently.
- Do not generate issues that imply architecture changes without maintainer approval.

## Troubleshooting

### Repository not found

Check owner/repo names and token permissions.

### Bad credentials

Check that `GITHUB_TOKEN` is set and not expired.

### Label already exists

Use update mode if available or adjust the config.

### Rate limiting

Split large operations into smaller batches.
