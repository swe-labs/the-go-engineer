# GitHub Automation Tool - Complete Guide

**Universal tool for creating labels, issues, and discussions on ANY GitHub repository**

## 📋 Features

✅ **Create Labels** - Custom colors, descriptions, organized by categories
✅ **Create Issues** - With labels, assignees, milestones, detailed descriptions
✅ **Create Discussions** - GitHub Discussions support
✅ **Any Repository** - Works with any GitHub repo (yours or others)
✅ **Config-Driven** - Separate configuration from code
✅ **JSON & YAML** - Support both formats
✅ **Dry-Run Mode** - Preview changes before applying
✅ **Portable** - Run from anywhere with just a token and config file

---

## 🚀 Quick Start

### Step 1: Install Dependencies

```bash
pip3 install requests pyyaml
```

### Step 2: Create Config File

Create `github-config.json`:

```json
{
  "owner": "your-username",
  "repo": "your-repo",
  "labels": [
    {
      "name": "epic",
      "color": "3E1B6F",
      "description": "Large feature"
    },
    {
      "name": "bug",
      "color": "D73A49",
      "description": "Bug report"
    }
  ],
  "issues": [
    {
      "title": "My First Issue",
      "body": "Issue description here",
      "labels": ["epic"]
    }
  ],
  "discussions": []
}
```

### Step 3: Generate GitHub Token

1. Go to: https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes:
   - `repo` (all)
   - `admin:repo_hook`
   - `workflow`
   - `read:discussion`
4. Copy token: `ghp_xxxxxxxxxxxx`

### Step 4: Run Dry-Run (Safe Preview)

```bash
python3 github-automation.py \
  --config github-config.json \
  --token ghp_xxxxxxxxxxxx \
  --dry-run
```

### Step 5: Execute

```bash
python3 github-automation.py \
  --config github-config.json \
  --token ghp_xxxxxxxxxxxx
```

---

## 📝 Configuration File Format

### JSON Format

```json
{
  "owner": "username",
  "repo": "repository-name",
  "labels": [
    {
      "name": "label-name",
      "color": "HEXCOLOR",
      "description": "Label description"
    }
  ],
  "issues": [
    {
      "title": "Issue Title",
      "body": "Issue description\n\nCan be multi-line",
      "labels": ["label1", "label2"],
      "assignees": ["username1", "username2"],
      "milestone": null
    }
  ],
  "discussions": [
    {
      "category_id": "DIC_kwDOxxxxxx",
      "title": "Discussion Title",
      "body": "Discussion content"
    }
  ]
}
```

### YAML Format

```yaml
owner: username
repo: repository-name

labels:
  - name: epic
    color: "3E1B6F"
    description: "Large feature"
  - name: bug
    color: "D73A49"
    description: "Bug report"

issues:
  - title: "Issue Title"
    body: |
      Issue description
      with multiple lines
    labels:
      - epic
    assignees: []

discussions: []
```

---

## 💻 Command Line Usage

### Basic Usage

```bash
# Auto-detect config file
python3 github-automation.py --token ghp_xxxx

# Specify config file
python3 github-automation.py --config config.json --token ghp_xxxx

# Use environment variable for token
python3 github-automation.py --config config.yaml --token $GH_TOKEN
```

### Dry-Run (Preview Only)

```bash
python3 github-automation.py \
  --config config.json \
  --token ghp_xxxx \
  --dry-run
```

### Selective Execution

```bash
# Labels only
python3 github-automation.py --config config.json --token ghp_xxxx --labels-only

# Issues only
python3 github-automation.py --config config.json --token ghp_xxxx --issues-only

# Discussions only
python3 github-automation.py --config config.json --token ghp_xxxx --discussions-only
```

---

## 🎨 Label Best Practices

### Color Naming Convention

Use meaningful colors that convey information:

```json
{
  "labels": [
    {
      "name": "epic",
      "color": "3E1B6F",
      "description": "Large feature"
    },
    {
      "name": "feature",
      "color": "0366D6",
      "description": "New feature"
    },
    {
      "name": "bug",
      "color": "D73A49",
      "description": "Bug fix"
    },
    {
      "name": "documentation",
      "color": "0075CA",
      "description": "Documentation"
    },
    {
      "name": "priority/critical",
      "color": "B60205",
      "description": "Urgent"
    },
    {
      "name": "priority/high",
      "color": "FF6B6B",
      "description": "Important"
    },
    {
      "name": "size/small",
      "color": "9BCD9B",
      "description": "Quick task"
    },
    {
      "name": "status/ready",
      "color": "C6E48B",
      "description": "Ready to work"
    },
    {
      "name": "status/in-progress",
      "color": "1E90FF",
      "description": "Being worked on"
    },
    {
      "name": "status/done",
      "color": "228B22",
      "description": "Completed"
    }
  ]
}
```

---

## 🎯 Issue Templates

### Epic Issue

```json
{
  "title": "[EPIC] Large Project",
  "body": "## Overview\nDescribe the epic\n\n## Sub-Tasks\n- [ ] Task 1\n- [ ] Task 2\n\n## Success Criteria\n✓ Criterion 1",
  "labels": ["epic", "priority/high"]
}
```

### Feature Issue

```json
{
  "title": "[FEATURE] New Capability",
  "body": "## Description\nWhat needs to be done?\n\n## Acceptance Criteria\n- [ ] Criterion 1\n- [ ] Criterion 2\n\n## Testing\n- [ ] Unit tests\n- [ ] Integration tests",
  "labels": ["feature", "priority/high", "size/medium"]
}
```

### Bug Issue

```json
{
  "title": "[BUG] Issue description",
  "body": "## Description\nWhat's broken?\n\n## Expected Behavior\nWhat should happen\n\n## Actual Behavior\nWhat actually happens\n\n## Steps to Reproduce\n1. Step 1\n2. Step 2",
  "labels": ["bug", "priority/high"]
}
```

### Documentation Issue

```json
{
  "title": "[DOCS] Missing documentation",
  "body": "## Overview\nWhat needs documentation?\n\n## Requirements\n- [ ] Requirement 1\n- [ ] Requirement 2\n\n## Success Criteria\n✓ Documented and reviewed",
  "labels": ["documentation"]
}
```

---

## 🔧 Environment Setup

### Save Token Safely

```bash
# On macOS/Linux
export GH_TOKEN="ghp_xxxxxxxxxxxx"

# Use in script
python3 github-automation.py --config config.json --token $GH_TOKEN
```

### Auto-Detect Config File

The script looks for these in order:
1. `github-config.json`
2. `github-config.yaml`
3. `.github-config.json`
4. `.github-config.yaml`

If found, you don't need `--config`:

```bash
python3 github-automation.py --token $GH_TOKEN
```

---

## 📂 Project Structure

```
scripts/
├── github-automation.py          # Main tool (reusable)
├── setup-github-issues.py        # Original specific script (for reference)
├── setup-github-issues.sh        # Shell alternative
├── github-config.example.json    # Example JSON config
├── github-config.example.yaml    # Example YAML config
└── README.md                     # This file
```

---

## 🔐 Security Considerations

### Token Management

```bash
# ✅ DO: Use environment variables
export GH_TOKEN="ghp_xxxx"
python3 github-automation.py --token $GH_TOKEN

# ✅ DO: Keep tokens in .env (gitignored)
# .env file (add to .gitignore)
GH_TOKEN=ghp_xxxx

# Load and use
source .env
python3 github-automation.py --token $GH_TOKEN

# ✅ DO: Use secrets in CI/CD
# .github/workflows/automation.yaml
env:
  GH_TOKEN: ${{ secrets.GH_TOKEN }}
run: python3 scripts/github-automation.py --token $GH_TOKEN

# ❌ DON'T: Hardcode token in config
# DON'T: Commit token to git
# DON'T: Share token in messages
```

### Token Permissions

Minimum required scopes:

```
✅ repo (full control)
✅ admin:repo_hook (webhook management)
✅ workflow (GitHub Actions)
✅ read:discussion (discussions)
```

---

## 🐛 Troubleshooting

### Error: "ModuleNotFoundError: No module named 'requests'"

```bash
pip3 install requests pyyaml
```

### Error: "Authentication failed"

```bash
# Verify token is correct
# Check token hasn't expired
# Regenerate new token from https://github.com/settings/tokens
```

### Error: "Cannot access repo"

```bash
# Check repo owner and name are correct
# Verify token has repo access permission
# Check if repo exists: https://github.com/owner/repo
```

### Error: "Config file not found"

```bash
# Specify config file explicitly
python3 github-automation.py --config /path/to/config.json --token ghp_xxxx

# Or place in current directory with auto-detected name:
# github-config.json or github-config.yaml
```

### Error: "PyYAML not installed"

```bash
pip3 install pyyaml
```

---

## 📊 Usage Examples

### Example 1: Create Labels for New Project

**config.json:**
```json
{
  "owner": "myusername",
  "repo": "myproject",
  "labels": [
    {"name": "feature", "color": "0366D6"},
    {"name": "bug", "color": "D73A49"},
    {"name": "priority/high", "color": "FF6B6B"}
  ],
  "issues": [],
  "discussions": []
}
```

**Run:**
```bash
python3 github-automation.py --config config.json --token $GH_TOKEN --labels-only
```

### Example 2: Setup Sprint Board

**config.json:**
```json
{
  "owner": "myusername",
  "repo": "myproject",
  "labels": [
    {"name": "sprint-ready", "color": "C6E48B"},
    {"name": "in-progress", "color": "1E90FF"},
    {"name": "review", "color": "9370DB"}
  ],
  "issues": [
    {"title": "Sprint Task 1", "labels": ["sprint-ready"]},
    {"title": "Sprint Task 2", "labels": ["sprint-ready"]}
  ],
  "discussions": []
}
```

### Example 3: Migrate Labels Between Repos

1. Export labels from source repo
2. Create config with same labels pointing to target repo
3. Run automation on target repo

---

## 🚀 Advanced Usage

### CI/CD Integration

**.github/workflows/automation.yaml:**
```yaml
name: GitHub Automation

on:
  workflow_dispatch:
    inputs:
      config:
        description: 'Config file'
        required: true

jobs:
  automate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.11'
      
      - name: Install dependencies
        run: pip install requests pyyaml
      
      - name: Run automation
        env:
          GH_TOKEN: ${{ secrets.GH_TOKEN }}
        run: |
          python3 scripts/github-automation.py \
            --config ${{ github.event.inputs.config }} \
            --token $GH_TOKEN
```

### Template Manager

Create multiple configs for different projects:

```
configs/
├── project-a.json
├── project-b.yaml
├── sprint-board.json
└── backlog.yaml
```

Then use as needed:

```bash
python3 github-automation.py --config configs/project-a.json --token $GH_TOKEN
python3 github-automation.py --config configs/sprint-board.json --token $GH_TOKEN
```

---

## 📖 Complete Workflow

```bash
# 1. Setup
pip3 install requests pyyaml
export GH_TOKEN="ghp_xxxxxxxxxxxx"

# 2. Create config
cat > my-config.json << 'EOF'
{
  "owner": "rasel9t6",
  "repo": "the-go-engineer",
  "labels": [...],
  "issues": [...],
  "discussions": []
}
EOF

# 3. Preview (dry-run)
python3 github-automation.py \
  --config my-config.json \
  --token $GH_TOKEN \
  --dry-run

# 4. Review output

# 5. Execute (if satisfied)
python3 github-automation.py \
  --config my-config.json \
  --token $GH_TOKEN

# 6. Verify on GitHub
# Visit: https://github.com/owner/repo/issues
```

---

## 📝 Notes

- **Idempotent**: Running twice with same config won't create duplicates
  - Existing labels are skipped
  - New issues are always created
  
- **Safe**: Use `--dry-run` to preview changes first

- **Fast**: Creates 50+ labels and issues in seconds

- **Flexible**: Works with any GitHub account or organization

---

## 🔗 Resources

- GitHub API Docs: https://docs.github.com/en/rest
- Personal Access Tokens: https://github.com/settings/tokens
- Issue Templates: https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues

---

## 📄 License

MIT - Use freely, modify as needed

## 👤 Author

GitHub Automation Tool - Universal GitHub automation for any project
