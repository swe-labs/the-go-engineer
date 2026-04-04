# GitHub Automation Scripts

**Complete suite of scripts for automating GitHub label and issue creation**

> **📁 Directory Reorganization Notice**
>
> As of recent updates, automation scripts have been moved to `maintainer-scripts/` subdirectory to separate maintainer tools from curriculum scripts. This keeps the main `scripts/` directory focused on Go-based validation and development tools for learners.

## 📁 Scripts Overview

| Script | Location | Purpose | Flexibility | Use Case |
|--------|----------|---------|-------------|----------|
| **`github-automation.py`** | `maintainer-scripts/` | Universal tool (config-driven) | ⭐⭐⭐⭐⭐ Highest | Any repo, any config |
| **`setup-github-issues.py`** | `maintainer-scripts/` | The Go Engineer specific | ⭐⭐ Medium | This project only |
| **`setup-github-issues.sh`** | `maintainer-scripts/` | Shell alternative | ⭐⭐ Medium | No Python needed |
| **`validate_curriculum.go`** | `scripts/` | Go curriculum validation | ⭐⭐⭐⭐ High | Curriculum integrity |

---

---

## 🎯 Which Script Should I Use?

### Use `github-automation.py` If:

✅ You want a **reusable tool** for multiple projects
✅ You need to create **different label sets**
✅ You want **maximum flexibility**
✅ You work with **multiple GitHub repos**
✅ You want to **separate config from code**
✅ You like **JSON or YAML configs**

**Perfect for:** Managing labels/issues across organizations, automating setup for multiple projects

### Use `setup-github-issues.py` If:

✅ You only need it for **The Go Engineer project**
✅ You want **hardcoded configuration**
✅ You prefer **no external config files**
✅ You're doing **one-time setup**

**Perfect for:** Quick setup without creating separate config files

### Use `setup-github-issues.sh` If:

✅ You don't have **Python installed**
✅ You need **shell-only solution**
✅ You have `curl` and `jq` available

**Perfect for:** Minimal dependencies environment

---

## 🚀 Quick Start Guide

### Option A: Universal Tool (Recommended)

```bash
# 1. Install dependencies
pip3 install requests pyyaml

# 2. Create config file (see examples)
cp maintainer-scripts/github-config.example.json my-config.json
# Edit my-config.json with your settings

# 3. Test with dry-run
python3 maintainer-scripts/github-automation.py \
  --config my-config.json \
  --token ghp_xxxxxxxxxxxx \
  --dry-run

# 4. Execute
python3 maintainer-scripts/github-automation.py \
  --config my-config.json \
  --token ghp_xxxxxxxxxxxx
```

### Option B: The Go Engineer Specific

```bash
# 1. Install dependencies
pip3 install requests pyyaml

# 2. Run directly
python3 maintainer-scripts/setup-github-issues.py \
  --token ghp_xxxxxxxxxxxx \
  --dry-run

# 3. Execute
python3 maintainer-scripts/setup-github-issues.py \
  --token ghp_xxxxxxxxxxxx
```

### Option C: Shell Version

```bash
# 1. Verify dependencies
curl --version
jq --version

# 2. Run
bash maintainer-scripts/setup-github-issues.sh \
  --token ghp_xxxxxxxxxxxx \
  --dry-run

# 3. Execute
bash maintainer-scripts/setup-github-issues.sh \
  --token ghp_xxxxxxxxxxxx
```

---

## 📋 File Descriptions

### `maintainer-scripts/github-automation.py`

**Universal GitHub automation tool**

```
Features:
- Config-driven (JSON or YAML)
- Works with any GitHub repo
- Create labels, issues, discussions
- Dry-run mode for safety
- Environment variable support
- Auto-detect config files

Usage:
python3 maintainer-scripts/github-automation.py --config config.json --token $GH_TOKEN

Options:
--config PATH          Config file (JSON/YAML)
--token TOKEN          GitHub Personal Access Token
--dry-run              Preview changes
--labels-only          Create labels only
--issues-only          Create issues only
--discussions-only     Create discussions only
```

**When to use:** Every day, for any project

### `maintainer-scripts/setup-github-issues.py`

**The Go Engineer specific setup script**

```
Features:
- Hardcoded: 38 labels
- Hardcoded: 7 issues
- No config file needed
- Quick one-time setup

Usage:
python3 maintainer-scripts/setup-github-issues.py --token $GH_TOKEN

Options:
--token TOKEN          GitHub Personal Access Token (required)
--dry-run              Preview without creating
--labels-only          Create labels only
--issues-only          Create issues only
```

**When to use:** First-time setup of The Go Engineer

### `maintainer-scripts/setup-github-issues.sh`

**Shell script alternative**

```
Features:
- No Python dependency
- Uses curl and jq
- Similar to Python version
- Works on any Unix-like system

Usage:
bash maintainer-scripts/setup-github-issues.sh --token $GH_TOKEN

Options:
--token TOKEN          GitHub Personal Access Token (required)
--dry-run              Preview without creating
--labels-only          Create labels only
--issues-only          Create issues only
-h, --help             Show help
```

**When to use:** When Python isn't available

### `maintainer-scripts/github-config.example.json`

**Example configuration in JSON format**

- Shows label structure
- Shows issue structure
- Shows discussions structure
- Copy and customize for your needs

**Commands:**
```bash
cp maintainer-scripts/github-config.example.json my-config.json
# Edit my-config.json
python3 maintainer-scripts/github-automation.py --config my-config.json --token $GH_TOKEN
```

### `maintainer-scripts/github-config.example.yaml`

**Example configuration in YAML format**

- Same structure as JSON
- More human-readable
- Better for manual editing
- Comments included

**Commands:**
```bash
cp maintainer-scripts/github-config.example.yaml my-config.yaml
# Edit my-config.yaml
python3 maintainer-scripts/github-automation.py --config my-config.yaml --token $GH_TOKEN
```

### `maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md`

**Complete documentation**

- Comprehensive guide
- Configuration formats
- Best practices
- Troubleshooting
- Advanced usage
- CI/CD integration

### `scripts/validate_curriculum.go`

**Go curriculum validation tool**

```
Features:
- Validates curriculum.json structure
- Checks lesson file existence
- Verifies Go code syntax
- Generates reports

Usage:
go run scripts/validate_curriculum.go

Options:
-h, --help     Show help
-v, --verbose  Verbose output
```

**When to use:** Validate curriculum integrity, check for broken links

---

## 🔐 Getting Your GitHub Token

### Step-by-Step

1. Go to: **https://github.com/settings/tokens**
2. Click **"Generate new token"** → **"Generate new token (classic)"**
3. **Fill in:**
   - Name: `My GitHub Automation Token`
   - Expiration: 90 days (or as needed)
4. **Select scopes:**
   ```
   ✅ repo (Full control of private/public repositories)
   ✅ admin:repo_hook
   ✅ workflow
   ✅ read:discussion
   ```
5. **Click "Generate token"**
6. **COPY immediately** (you won't see it again): `ghp_xxxxxxxxxxxx`
7. **Save securely** (don't commit to git!)

### Secure Storage

```bash
# Option 1: Environment variable
export GH_TOKEN="ghp_xxxxxxxxxxxx"
python3 github-automation.py --token $GH_TOKEN

# Option 2: .env file (add to .gitignore)
echo "GH_TOKEN=ghp_xxxxxxxxxxxx" > .env
source .env
python3 github-automation.py --token $GH_TOKEN

# Option 3: GitHub CLI
gh auth login
# Then use: python3 github-automation.py --token $(gh auth token)
```

---

## 💾 Installation

### Prerequisites

```bash
# Python 3.8+
python3 --version

# Git (optional, for cloning)
git --version

# curl (only for shell version)
curl --version

# jq (only for shell version)
jq --version
```

### Setup

```bash
# Clone repo (if needed)
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer

# Navigate to scripts
cd scripts

# Install Python dependencies (for maintainer scripts)
pip3 install requests pyyaml

# Make scripts executable (optional)
chmod +x maintainer-scripts/github-automation.py maintainer-scripts/setup-github-issues.py maintainer-scripts/setup-github-issues.sh

# For Go validation script
go version  # Verify Go is installed
```

---

## 🎯 Common Workflows

### Workflow 1: Setup New Project

```bash
# 1. Create config from example
cp maintainer-scripts/github-config.example.json my-project.json

# 2. Edit with your labels and issues
nano my-project.json

# 3. Test
python3 maintainer-scripts/github-automation.py \
  --config my-project.json \
  --token $GH_TOKEN \
  --dry-run

# 4. Create
python3 maintainer-scripts/github-automation.py \
  --config my-project.json \
  --token $GH_TOKEN
```

### Workflow 2: Update Existing Project

```bash
# 1. Add labels to existing config
# (Edit github-config.json or my-config.json)

# 2. Run again (existing labels skipped, new ones added)
python3 maintainer-scripts/github-automation.py \
  --config my-config.json \
  --token $GH_TOKEN
```

### Workflow 3: Create Issues Only

```bash
# 1. Use config with existing labels
python3 maintainer-scripts/github-automation.py \
  --config my-config.json \
  --token $GH_TOKEN \
  --issues-only
```

### Workflow 4: Multiple Projects

```bash
# Create separate configs
configs/
├── project-a.json
├── project-b.json
└── project-c.json

# Run for each
python3 maintainer-scripts/github-automation.py --config configs/project-a.json --token $GH_TOKEN
python3 maintainer-scripts/github-automation.py --config configs/project-b.json --token $GH_TOKEN
python3 maintainer-scripts/github-automation.py --config configs/project-c.json --token $GH_TOKEN
```

---

## 🔧 Configuration Guide

### Minimal Config

```json
{
  "owner": "your-username",
  "repo": "your-repo",
  "labels": [],
  "issues": [],
  "discussions": []
}
```

### Label Structure

```json
{
  "name": "label-name",
  "color": "RRGGBB",
  "description": "Label description"
}
```

- **name**: Label name (required)
- **color**: Hex color without # (required)
- **description**: Optional description

### Issue Structure

```json
{
  "title": "Issue Title",
  "body": "Issue body (markdown supported)",
  "labels": ["label1", "label2"],
  "assignees": ["username1"],
  "milestone": null
}
```

- **title**: Issue title (required)
- **body**: Markdown-formatted body
- **labels**: Array of label names
- **assignees**: Array of GitHub usernames
- **milestone**: Milestone ID or null

---

## 🧪 Testing

### Dry-Run (Always Safe)

```bash
python3 github-automation.py \
  --config config.json \
  --token $GH_TOKEN \
  --dry-run
```

Output will show:
- `[DRY-RUN]` prefix on all actions
- No actual changes made
- Safe to review output

### Selective Testing

```bash
# Test labels only
python3 github-automation.py --config config.json --token $GH_TOKEN --labels-only --dry-run

# Test issues only
python3 github-automation.py --config config.json --token $GH_TOKEN --issues-only --dry-run
```

---

## 📊 Examples

### Example 1: Simple Setup

**config.json:**
```json
{
  "owner": "myname",
  "repo": "myproject",
  "labels": [
    {"name": "bug", "color": "D73A49"},
    {"name": "feature", "color": "0366D6"}
  ],
  "issues": [
    {"title": "First Issue", "labels": ["bug"]}
  ],
  "discussions": []
}
```

**Run:**
```bash
python3 github-automation.py --config config.json --token $GH_TOKEN
```

### Example 2: Comprehensive Setup

See `github-config.example.json` for comprehensive example with:
- 13 labels across categories
- 4 detailed issues
- Markdown bodies with formatting

---

## ⚙️ Advanced Features

### Environment Variables

```bash
# Use in scripts
export GH_TOKEN="ghp_xxxx"
export CONFIG="my-config.json"

# Then use
python3 github-automation.py --config $CONFIG --token $GH_TOKEN
```

### Auto-Detect Config

Place config file with standard name in current directory:
- `github-config.json`
- `github-config.yaml`
- `.github-config.json`
- `.github-config.yaml`

Then run without `--config`:
```bash
python3 github-automation.py --token $GH_TOKEN
```

### CI/CD Integration

See `GITHUB_AUTOMATION_GUIDE.md` for GitHub Actions example

---

## 🐛 Troubleshooting

### "Command not found"

```bash
# Make executable
chmod +x maintainer-scripts/github-automation.py

# Or run with python
python3 maintainer-scripts/github-automation.py --config config.json --token $GH_TOKEN
```

### "ModuleNotFoundError: No module named 'requests'"

```bash
pip3 install requests pyyaml
```

### "Authentication failed"

```bash
# Verify token is correct
# Check token hasn't expired
# Regenerate token: https://github.com/settings/tokens
```

### "Cannot access repo"

```bash
# Verify owner/repo names
# Check token has repo access
# Test: curl -H "Authorization: token $GH_TOKEN" https://api.github.com/repos/owner/repo
```

---

## 📚 Documentation

For comprehensive documentation, see:

- **Setup Guide**: GITHUB_AUTOMATION_GUIDE.md
- **Examples**: github-config.example.json, github-config.example.yaml
- **References**: Original scripts (for understanding implementation)

---

## 🤝 Keep or Delete?

### ✅ KEEP

- `maintainer-scripts/github-automation.py` - Reusable for any project
- `maintainer-scripts/github-config.example.json` - Reference documentation
- `maintainer-scripts/github-config.example.yaml` - Reference documentation
- `maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md` - Complete guide
- `scripts/validate_curriculum.go` - Go curriculum validation
- This README.md

### 🗑️ OPTIONAL (Can Delete)

- `maintainer-scripts/setup-github-issues.py` - Original specific script (duplicates automation.py)
- `maintainer-scripts/setup-github-issues.sh` - Alternative (if using Python)
- `maintainer-scripts/github-config.example.*` - If you keep templates elsewhere

### 🚀 Recommendations

1. **Keep `maintainer-scripts/github-automation.py`** - It's universal and reusable
2. **Keep `scripts/validate_curriculum.go`** - Essential for curriculum integrity
3. **Delete `maintainer-scripts/setup-github-issues.py`** - Superseded by automation.py
4. **Keep examples** - Useful as reference
5. **Keep documentation** - Essential for future use

---

## 📄 License

MIT - Use freely for any purpose

---

## 🎯 Quick Reference Card

```bash
# Installation
pip3 install requests pyyaml

# Set token
export GH_TOKEN="ghp_xxxxxxxxxxxx"

# Create config
cp maintainer-scripts/github-config.example.json my-config.json
nano my-config.json

# Test
python3 maintainer-scripts/github-automation.py --config my-config.json --token $GH_TOKEN --dry-run

# Execute
python3 maintainer-scripts/github-automation.py --config my-config.json --token $GH_TOKEN

# Specific targets
python3 maintainer-scripts/github-automation.py --config my-config.json --token $GH_TOKEN --labels-only
python3 maintainer-scripts/github-automation.py --config my-config.json --token $GH_TOKEN --issues-only

# View help
python3 maintainer-scripts/github-automation.py -h

# Validate curriculum (Go tool)
go run scripts/validate_curriculum.go
```

---

**Need help?** See `maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md` for complete documentation.
