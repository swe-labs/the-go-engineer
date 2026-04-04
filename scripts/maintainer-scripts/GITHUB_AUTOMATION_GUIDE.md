# GitHub Automation Tool

A universal Python script for automating GitHub repository setup and maintenance tasks. This tool helps maintainers quickly configure labels, create issues, and manage repository structure for educational and open-source projects.

## Features

- **Universal Configuration**: Works with any GitHub repository
- **Label Management**: Create, update, and manage issue labels
- **Issue Creation**: Bulk create issues from configuration
- **Discussion Support**: Create discussion posts (if enabled)
- **Flexible Configuration**: Support for both JSON and YAML formats
- **Dry Run Mode**: Test configurations without making changes
- **Progress Tracking**: Real-time progress bars and status updates
- **Error Handling**: Robust error handling with detailed logging

## Quick Start

### 1. Prerequisites

- Python 3.8+
- GitHub Personal Access Token with appropriate permissions
- Required Python packages: `requests`, `pyyaml`

### 2. Installation

```bash
# Install dependencies
pip install requests pyyaml

# Or using requirements file
pip install -r requirements.txt
```

### 3. Configuration

Copy the example configuration file:

```bash
cp github-config.example.yaml github-config.yaml
# OR
cp github-config.example.json github-config.json
```

Edit the configuration file with your repository details:

```yaml
owner: your-username
repo: your-repo-name
```

### 4. Authentication

Set your GitHub token as an environment variable:

```bash
export GITHUB_TOKEN=your_personal_access_token_here
```

Or create a `.env` file in the same directory:

```
GITHUB_TOKEN=your_personal_access_token_here
```

### 5. Run the Tool

```bash
# Create labels only
python github-automation.py --labels-only

# Create issues only
python github-automation.py --issues-only

# Create everything (labels + issues)
python github-automation.py --all

# Dry run (preview changes without applying)
python github-automation.py --dry-run --all
```

## Configuration Format

### Repository Settings

```yaml
owner: your-github-username
repo: your-repository-name
```

### Labels Configuration

Each label requires:
- `name`: Label name (required)
- `color`: Hex color code without # (required)
- `description`: Optional description

```yaml
labels:
  - name: bug
    color: "D73A49"
    description: "Something is broken"
  - name: feature
    color: "0366D6"
    description: "New feature request"
```

### Issues Configuration

Each issue supports:
- `title`: Issue title (required)
- `body`: Issue description in Markdown (required)
- `labels`: Array of label names to apply
- `assignees`: Array of GitHub usernames

```yaml
issues:
  - title: "[FEATURE] Add new functionality"
    body: |
      ## Description
      This issue describes the new feature.

      ## Requirements
      - [ ] Requirement 1
      - [ ] Requirement 2
    labels:
      - feature
      - priority/high
    assignees:
      - username1
      - username2
```

### Discussions Configuration

For repositories with discussions enabled:

```yaml
discussions:
  - category_id: "DIC_kwDOxxxxxx"
    title: "Welcome to our discussions!"
    body: "Feel free to ask questions and share ideas."
```

## Command Line Options

```
Usage: python github-automation.py [OPTIONS]

Options:
  --config FILE          Configuration file path (default: github-config.yaml or .json)
  --labels-only          Create/update labels only
  --issues-only          Create issues only
  --discussions-only     Create discussions only
  --all                  Create labels, issues, and discussions
  --dry-run              Preview changes without applying them
  --force                Skip confirmation prompts
  --verbose              Enable verbose logging
  --help                 Show this help message
```

## Examples

### Educational Repository Setup

For a Go learning repository like "the-go-engineer":

```yaml
owner: rasel9t6
repo: the-go-engineer

labels:
  - name: lesson
    color: "0366D6"
    description: "New lesson or tutorial"
  - name: exercise
    color: "C6E48B"
    description: "Practice exercise"
  - name: bug
    color: "D73A49"
    description: "Code issue or bug"

issues:
  - title: "[LESSON] Add concurrency patterns"
    body: "Create a new section on Go concurrency patterns"
    labels: [lesson, priority/high]
```

### Project Management Setup

```yaml
labels:
  - name: epic
    color: "3E1B6F"
    description: "Large feature spanning multiple issues"
  - name: priority/critical
    color: "B60205"
    description: "Breaks existing functionality"
  - name: status/in-progress
    color: "1E90FF"
    description: "Currently being worked on"
```

## GitHub Token Permissions

Your Personal Access Token needs these permissions:

- `repo` - Full control of private repositories
- `public_repo` - Access public repositories
- `issues` - Read/write access to issues
- `discussions` - Read/write access to discussions (if using)

## Troubleshooting

### Common Issues

1. **"Repository not found"**
   - Check owner/repo names in config
   - Verify token has access to the repository

2. **"Bad credentials"**
   - Ensure GITHUB_TOKEN is set correctly
   - Check token hasn't expired

3. **"Label already exists"**
   - Use `--force` to update existing labels
   - Or remove existing labels manually first

4. **"Discussions not enabled"**
   - Enable discussions in repository settings
   - Or remove discussions from config

### Debug Mode

Enable verbose logging for detailed error information:

```bash
python github-automation.py --verbose --all
```

### Rate Limiting

GitHub API has rate limits. The tool includes automatic retry logic, but for large operations, consider:

- Use dry-run first to verify configuration
- Split large operations into smaller batches
- Wait between operations if needed

## Advanced Usage

### Custom Configuration Files

```bash
# Use a specific config file
python github-automation.py --config my-custom-config.yaml --all

# Use JSON format
python github-automation.py --config config.json --labels-only
```

### Environment Variables

```bash
# Override config values
export GITHUB_OWNER=myorg
export GITHUB_REPO=myproject

# Use different token for different operations
export GITHUB_TOKEN_LABELS=token1
export GITHUB_TOKEN_ISSUES=token2
```

### Integration with CI/CD

Add to your GitHub Actions workflow:

```yaml
- name: Setup Repository
  run: |
    pip install requests pyyaml
    python scripts/maintainer-scripts/github-automation.py --all
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with `--dry-run`
5. Submit a pull request

## License

This tool is part of the Go Engineer project. See LICENSE file for details.