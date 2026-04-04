#!/usr/bin/env python3
"""
GitHub Automation Tool
Universal script for creating labels, issues, and discussions on ANY GitHub repository

Features:
- ✅ Create custom labels with colors and descriptions
- ✅ Create issues with labels, milestones, and assignees
- ✅ Create discussions (GitHub Discussions API)
- ✅ Support multiple repositories
- ✅ Config-driven (JSON or YAML)
- ✅ Dry-run mode for safety
- ✅ Works from anywhere with just a token and config file

Usage:
    python3 github-automation.py --config config.json --token ghp_xxxx
    python3 github-automation.py --config config.yaml --token $GH_TOKEN --dry-run
    python3 github-automation.py --config config.json --token $GH_TOKEN --issues-only

Requirements:
    pip3 install requests pyyaml

Author: GitHub Automation Tool
License: MIT
"""

import argparse
import json
import os
import sys
from typing import Dict, List, Optional, Tuple, Any
import requests
from datetime import datetime

try:
    import yaml
    YAML_AVAILABLE = True
except ImportError:
    YAML_AVAILABLE = False

# ============================================================================
# CONSTANTS
# ============================================================================

API_BASE = "https://api.github.com"
DEFAULT_CONFIG_PATHS = [
    "github-config.json",
    "github-config.yaml",
    ".github-config.json",
    ".github-config.yaml",
]

# ============================================================================
# CONFIGURATION LOADER
# ============================================================================

class ConfigLoader:
    """Load configuration from JSON or YAML files"""
    
    @staticmethod
    def load(config_path: str) -> Dict[str, Any]:
        """Load config from file"""
        if not os.path.exists(config_path):
            raise FileNotFoundError(f"Config file not found: {config_path}")
        
        with open(config_path, 'r') as f:
            if config_path.endswith('.json'):
                return json.load(f)
            elif config_path.endswith('.yaml') or config_path.endswith('.yml'):
                if not YAML_AVAILABLE:
                    raise ImportError("PyYAML not installed. Install with: pip3 install pyyaml")
                return yaml.safe_load(f)
            else:
                raise ValueError(f"Unsupported config format: {config_path}")
    
    @staticmethod
    def find_default():
        """Search for default config file"""
        for path in DEFAULT_CONFIG_PATHS:
            if os.path.exists(path):
                return path
        return None

# ============================================================================
# GITHUB API CLIENT
# ============================================================================

class GitHubAPI:
    """Universal GitHub API client"""
    
    def __init__(self, token: str, dry_run: bool = False, verbose: bool = True):
        self.token = token
        self.dry_run = dry_run
        self.verbose = verbose
        self.headers = {
            "Authorization": f"token {token}",
            "Accept": "application/vnd.github.v3+json",
            "User-Agent": "GitHub-Automation-Tool",
        }
        self.session = requests.Session()
        self.session.headers.update(self.headers)
    
    def _log(self, message: str):
        """Log message if verbose enabled"""
        if self.verbose:
            print(message)
    
    def verify_connection(self, owner: str, repo: str) -> bool:
        """Verify API connection and repo access"""
        url = f"{API_BASE}/repos/{owner}/{repo}"
        try:
            response = self.session.get(url)
            if response.status_code == 200:
                user = response.json()["owner"]["login"]
                repo_name = response.json()["name"]
                self._log(f"✓ Authenticated and repo access verified: {user}/{repo_name}")
                return True
            else:
                self._log(f"✗ Repo access denied: {response.status_code}")
                return False
        except Exception as e:
            self._log(f"✗ Connection error: {str(e)}")
            return False
    
    def verify_user(self) -> Optional[str]:
        """Verify authentication and return username"""
        url = f"{API_BASE}/user"
        try:
            response = self.session.get(url)
            if response.status_code == 200:
                user = response.json()["login"]
                self._log(f"✓ Authenticated as: {user}")
                return user
            else:
                self._log(f"✗ Authentication failed: {response.status_code}")
                return None
        except Exception as e:
            self._log(f"✗ Connection error: {str(e)}")
            return None
    
    # ========================================================================
    # LABEL OPERATIONS
    # ========================================================================
    
    def create_label(self, owner: str, repo: str, name: str, color: str, 
                    description: str = "") -> Tuple[bool, str]:
        """Create a label"""
        url = f"{API_BASE}/repos/{owner}/{repo}/labels"
        data = {
            "name": name,
            "color": color,
        }
        if description:
            data["description"] = description
        
        if self.dry_run:
            self._log(f"  [DRY-RUN] Would create label: {name} (#{color})")
            return True, f"Label '{name}' (dry-run)"
        
        try:
            response = self.session.post(url, json=data)
            
            if response.status_code == 201:
                return True, f"✓ Label '{name}' created"
            elif response.status_code == 422:
                return True, f"⚠ Label '{name}' already exists"
            else:
                error = response.json().get("message", "Unknown error")
                return False, f"✗ Failed to create label '{name}': {error}"
        except Exception as e:
            return False, f"✗ Error creating label '{name}': {str(e)}"
    
    def get_labels(self, owner: str, repo: str) -> Dict[str, Dict]:
        """Get all labels from repo"""
        url = f"{API_BASE}/repos/{owner}/{repo}/labels"
        try:
            response = self.session.get(url)
            if response.status_code == 200:
                labels = {}
                for label in response.json():
                    labels[label["name"]] = {
                        "color": label["color"],
                        "description": label.get("description", "")
                    }
                return labels
            return {}
        except Exception:
            return {}
    
    # ========================================================================
    # ISSUE OPERATIONS
    # ========================================================================
    
    def create_issue(self, owner: str, repo: str, title: str, body: str = "",
                    labels: List[str] = None, assignees: List[str] = None,
                    milestone: Optional[int] = None) -> Tuple[bool, str, Optional[int]]:
        """Create an issue"""
        url = f"{API_BASE}/repos/{owner}/{repo}/issues"
        data = {
            "title": title,
            "body": body or "",
        }
        if labels:
            data["labels"] = labels
        if assignees:
            data["assignees"] = assignees
        if milestone:
            data["milestone"] = milestone
        
        if self.dry_run:
            self._log(f"  [DRY-RUN] Would create issue: {title}")
            return True, f"Issue '{title}' (dry-run)", None
        
        try:
            response = self.session.post(url, json=data)
            
            if response.status_code == 201:
                issue_num = response.json()["number"]
                return True, f"✓ Issue #{issue_num}: {title}", issue_num
            else:
                error = response.json().get("message", "Unknown error")
                return False, f"✗ Failed to create issue: {error}", None
        except Exception as e:
            return False, f"✗ Error creating issue: {str(e)}", None
    
    # ========================================================================
    # DISCUSSION OPERATIONS
    # ========================================================================
    
    def create_discussion(self, owner: str, repo: str, category_id: str,
                         title: str, body: str = "") -> Tuple[bool, str]:
        """Create a discussion (requires GraphQL)"""
        url = f"{API_BASE}/graphql"
        
        # Update headers for GraphQL
        headers = self.headers.copy()
        headers["Accept"] = "application/vnd.github.v3+json"
        
        query = """
        mutation($repositoryId:ID!,$categoryId:ID!,$title:String!,$body:String!) {
            createDiscussion(input:{repositoryId:$repositoryId,categoryId:$categoryId,title:$title,body:$body}) {
                discussion {
                    id
                    title
                    url
                }
            }
        }
        """
        
        variables = {
            "repositoryId": repo,
            "categoryId": category_id,
            "title": title,
            "body": body or "",
        }
        
        if self.dry_run:
            self._log(f"  [DRY-RUN] Would create discussion: {title}")
            return True, f"Discussion '{title}' (dry-run)"
        
        try:
            response = self.session.post(
                url,
                json={"query": query, "variables": variables},
                headers=headers
            )
            
            if response.status_code == 200:
                data = response.json()
                if "errors" in data:
                    return False, f"✗ Failed to create discussion: {data['errors'][0]['message']}"
                else:
                    disc_url = data["data"]["createDiscussion"]["discussion"]["url"]
                    return True, f"✓ Discussion created: {disc_url}"
            else:
                return False, f"✗ Failed to create discussion: {response.status_code}"
        except Exception as e:
            return False, f"✗ Error creating discussion: {str(e)}"

# ============================================================================
# MAIN ORCHESTRATOR
# ============================================================================

class GitHubAutomation:
    """Orchestrate GitHub automation tasks"""
    
    def __init__(self, config: Dict[str, Any], token: str, dry_run: bool = False):
        self.config = config
        self.token = token
        self.dry_run = dry_run
        self.api = GitHubAPI(token, dry_run=dry_run)
    
    def validate_config(self) -> bool:
        """Validate config structure"""
        required = ["owner", "repo"]
        for field in required:
            if field not in self.config:
                print(f"✗ Config error: Missing required field '{field}'")
                return False
        
        return True
    
    def run(self, labels_only: bool = False, issues_only: bool = False,
           discussions_only: bool = False) -> bool:
        """Execute automation"""
        owner = self.config["owner"]
        repo = self.config["repo"]
        
        # Header
        mode = "[DRY-RUN] " if self.dry_run else ""
        print(f"\n{'='*70}")
        print(f"  {mode}GitHub Automation for {owner}/{repo}")
        print(f"{'='*70}\n")
        
        # Verify connection
        print("Step 1: Verifying GitHub API connection...")
        if not self.api.verify_user():
            print("\n✗ Authentication failed. Check your token.")
            return False
        
        if not self.api.verify_connection(owner, repo):
            print(f"\n✗ Cannot access repo {owner}/{repo}. Check permissions.")
            return False
        print()
        
        # Create labels
        if not issues_only and not discussions_only:
            self._create_labels(owner, repo)
        
        # Create issues
        if not labels_only and not discussions_only:
            self._create_issues(owner, repo)
        
        # Create discussions
        if not labels_only and not issues_only:
            self._create_discussions(owner, repo)
        
        # Summary
        print("="*70)
        if self.dry_run:
            print("  ✓ Dry-run completed successfully!")
            print("  Run without --dry-run to actually create resources.")
        else:
            print("  ✓ Automation completed successfully!")
            print(f"  Visit: https://github.com/{owner}/{repo}/issues")
        print("="*70)
        print()
        
        return True
    
    def _create_labels(self, owner: str, repo: str):
        """Create labels from config"""
        if "labels" not in self.config or not self.config["labels"]:
            return
        
        labels = self.config["labels"]
        print(f"Step 2: Creating labels...")
        print(f"  Total labels to create: {len(labels)}\n")
        
        created = 0
        for label in labels:
            name = label.get("name")
            color = label.get("color")
            description = label.get("description", "")
            
            if not name or not color:
                print(f"  ⚠ Skipping invalid label (missing name or color)")
                continue
            
            success, message = self.api.create_label(owner, repo, name, color, description)
            print(f"  {message}")
            if success:
                created += 1
        
        print(f"\n  Summary: {created}/{len(labels)} created/verified\n")
    
    def _create_issues(self, owner: str, repo: str):
        """Create issues from config"""
        if "issues" not in self.config or not self.config["issues"]:
            return
        
        issues = self.config["issues"]
        print(f"Step 3: Creating issues...")
        print(f"  Total issues to create: {len(issues)}\n")
        
        created = 0
        for i, issue in enumerate(issues, 1):
            title = issue.get("title")
            body = issue.get("body", "")
            labels = issue.get("labels", [])
            assignees = issue.get("assignees", [])
            milestone = issue.get("milestone")
            
            if not title:
                print(f"  ⚠ Issue #{i}: Missing title, skipping")
                continue
            
            print(f"  Creating issue #{i}: {title}")
            success, message, issue_num = self.api.create_issue(
                owner, repo, title, body, labels, assignees, milestone
            )
            print(f"    {message}")
            
            if success:
                created += 1
        
        print(f"\n  Summary: {created}/{len(issues)} created\n")
    
    def _create_discussions(self, owner: str, repo: str):
        """Create discussions from config"""
        if "discussions" not in self.config or not self.config["discussions"]:
            return
        
        discussions = self.config["discussions"]
        print(f"Step 4: Creating discussions...")
        print(f"  Total discussions to create: {len(discussions)}\n")
        
        created = 0
        for i, discussion in enumerate(discussions, 1):
            category_id = discussion.get("category_id")
            title = discussion.get("title")
            body = discussion.get("body", "")
            
            if not category_id or not title:
                print(f"  ⚠ Discussion #{i}: Missing category_id or title, skipping")
                continue
            
            print(f"  Creating discussion #{i}: {title}")
            success, message = self.api.create_discussion(
                owner, repo, category_id, title, body
            )
            print(f"    {message}")
            
            if success:
                created += 1
        
        print(f"\n  Summary: {created}/{len(discussions)} created\n")

# ============================================================================
# CLI
# ============================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Universal GitHub Automation Tool",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 github-automation.py --config config.json --token ghp_xxxx
  python3 github-automation.py --config config.yaml --token $GH_TOKEN --dry-run
  python3 github-automation.py --token ghp_xxxx  # Auto-detect config
  python3 github-automation.py --config config.json --token ghp_xxxx --issues-only

Config File Format:
  {
    "owner": "rasel9t6",
    "repo": "the-go-engineer",
    "labels": [
      {"name": "epic", "color": "3E1B6F", "description": "Large feature"},
      ...
    ],
    "issues": [
      {"title": "Issue title", "body": "...", "labels": ["epic"], ...},
      ...
    ],
    "discussions": [
      {"category_id": "DIC_xxx", "title": "...", "body": "..."},
      ...
    ]
  }
        """
    )
    
    parser.add_argument(
        "--config",
        help="Config file path (JSON or YAML). Auto-detected if not provided."
    )
    
    parser.add_argument(
        "--token",
        required=True,
        help="GitHub Personal Access Token"
    )
    
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Preview changes without creating anything"
    )
    
    parser.add_argument(
        "--labels-only",
        action="store_true",
        help="Only create labels"
    )
    
    parser.add_argument(
        "--issues-only",
        action="store_true",
        help="Only create issues"
    )
    
    parser.add_argument(
        "--discussions-only",
        action="store_true",
        help="Only create discussions"
    )
    
    args = parser.parse_args()
    
    # Find config file
    config_path = args.config
    if not config_path:
        config_path = ConfigLoader.find_default()
        if config_path:
            print(f"ℹ Auto-detected config: {config_path}\n")
        else:
            print("✗ No config file found. Provide with --config")
            sys.exit(1)
    
    # Load config
    try:
        config = ConfigLoader.load(config_path)
    except Exception as e:
        print(f"✗ Config error: {e}")
        sys.exit(1)
    
    # Validate config
    automation = GitHubAutomation(config, args.token, dry_run=args.dry_run)
    if not automation.validate_config():
        sys.exit(1)
    
    # Run automation
    success = automation.run(
        labels_only=args.labels_only,
        issues_only=args.issues_only,
        discussions_only=args.discussions_only
    )
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
