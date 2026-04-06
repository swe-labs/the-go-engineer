#!/usr/bin/env python3
"""
GitHub Automation Tool
Universal script for creating labels, milestones, issues, and discussions.

Usage:
    python3 github-automation.py --config config.json --token ghp_xxxx
    python3 github-automation.py --config config.yaml --token $GH_TOKEN --dry-run
    python3 github-automation.py --config config.json --token $GH_TOKEN --milestones-only

Requirements:
    pip3 install requests pyyaml
"""

import argparse
import json
import os
import sys
from typing import Any, Dict, List, Optional, Tuple

import requests

try:
    import yaml

    YAML_AVAILABLE = True
except ImportError:
    YAML_AVAILABLE = False


API_BASE = "https://api.github.com"
DEFAULT_CONFIG_PATHS = [
    "github-config.json",
    "github-config.yaml",
    ".github-config.json",
    ".github-config.yaml",
]


class ConfigLoader:
    """Load configuration from JSON or YAML files."""

    @staticmethod
    def load(config_path: str) -> Dict[str, Any]:
        if not os.path.exists(config_path):
            raise FileNotFoundError(f"Config file not found: {config_path}")

        with open(config_path, "r", encoding="utf-8") as f:
            if config_path.endswith(".json"):
                return json.load(f)
            if config_path.endswith(".yaml") or config_path.endswith(".yml"):
                if not YAML_AVAILABLE:
                    raise ImportError("PyYAML not installed. Install with: pip3 install pyyaml")
                return yaml.safe_load(f)
            raise ValueError(f"Unsupported config format: {config_path}")

    @staticmethod
    def find_default() -> Optional[str]:
        for path in DEFAULT_CONFIG_PATHS:
            if os.path.exists(path):
                return path
        return None


class GitHubAPI:
    """Small GitHub REST client for setup automation."""

    def __init__(self, token: str, dry_run: bool = False, verbose: bool = True):
        self.dry_run = dry_run
        self.verbose = verbose
        self.session = requests.Session()
        self.session.headers.update(
            {
                "Authorization": f"token {token}",
                "Accept": "application/vnd.github.v3+json",
                "User-Agent": "GitHub-Automation-Tool",
            }
        )

    def _log(self, message: str) -> None:
        if self.verbose:
            print(message)

    def verify_user(self) -> Optional[str]:
        url = f"{API_BASE}/user"
        try:
            response = self.session.get(url)
            if response.status_code == 200:
                user = response.json()["login"]
                self._log(f"OK Authenticated as: {user}")
                return user
            self._log(f"ERROR Authentication failed: {response.status_code}")
            return None
        except Exception as exc:
            self._log(f"ERROR Connection error: {exc}")
            return None

    def verify_connection(self, owner: str, repo: str) -> bool:
        url = f"{API_BASE}/repos/{owner}/{repo}"
        try:
            response = self.session.get(url)
            if response.status_code == 200:
                repo_json = response.json()
                self._log(f"OK Repo access verified: {repo_json['owner']['login']}/{repo_json['name']}")
                return True
            self._log(f"ERROR Repo access denied: {response.status_code}")
            return False
        except Exception as exc:
            self._log(f"ERROR Connection error: {exc}")
            return False

    def create_label(self, owner: str, repo: str, name: str, color: str, description: str = "") -> Tuple[bool, str]:
        url = f"{API_BASE}/repos/{owner}/{repo}/labels"
        data: Dict[str, Any] = {"name": name, "color": color}
        if description:
            data["description"] = description

        if self.dry_run:
            self._log(f"  [DRY-RUN] Would create label: {name} (#{color})")
            return True, f"Label '{name}' (dry-run)"

        try:
            response = self.session.post(url, json=data)
            if response.status_code == 201:
                return True, f"OK Label '{name}' created"
            if response.status_code == 422:
                return True, f"WARN Label '{name}' already exists"
            error = response.json().get("message", "Unknown error")
            return False, f"ERROR Failed to create label '{name}': {error}"
        except Exception as exc:
            return False, f"ERROR Error creating label '{name}': {exc}"

    def get_milestones(self, owner: str, repo: str) -> Dict[str, int]:
        url = f"{API_BASE}/repos/{owner}/{repo}/milestones?state=all"
        try:
            response = self.session.get(url)
            if response.status_code != 200:
                return {}
            return {milestone["title"]: milestone["number"] for milestone in response.json()}
        except Exception:
            return {}

    def create_milestone(
        self,
        owner: str,
        repo: str,
        title: str,
        description: str = "",
        due_on: Optional[str] = None,
        state: str = "open",
    ) -> Tuple[bool, str, Optional[int]]:
        url = f"{API_BASE}/repos/{owner}/{repo}/milestones"
        data: Dict[str, Any] = {"title": title, "state": state}
        if description:
            data["description"] = description
        if due_on:
            data["due_on"] = due_on

        if self.dry_run:
            self._log(f"  [DRY-RUN] Would create milestone: {title}")
            return True, f"Milestone '{title}' (dry-run)", None

        try:
            response = self.session.post(url, json=data)
            if response.status_code == 201:
                milestone_num = response.json()["number"]
                return True, f"OK Milestone '{title}' created", milestone_num
            if response.status_code == 422:
                return True, f"WARN Milestone '{title}' already exists", None
            error = response.json().get("message", "Unknown error")
            return False, f"ERROR Failed to create milestone '{title}': {error}", None
        except Exception as exc:
            return False, f"ERROR Error creating milestone '{title}': {exc}", None

    def create_issue(
        self,
        owner: str,
        repo: str,
        title: str,
        body: str = "",
        labels: Optional[List[str]] = None,
        assignees: Optional[List[str]] = None,
        milestone: Optional[int] = None,
    ) -> Tuple[bool, str, Optional[int]]:
        url = f"{API_BASE}/repos/{owner}/{repo}/issues"
        data: Dict[str, Any] = {"title": title, "body": body or ""}
        if labels:
            data["labels"] = labels
        if assignees:
            data["assignees"] = assignees
        if milestone is not None:
            data["milestone"] = milestone

        if self.dry_run:
            self._log(f"  [DRY-RUN] Would create issue: {title}")
            return True, f"Issue '{title}' (dry-run)", None

        try:
            response = self.session.post(url, json=data)
            if response.status_code == 201:
                issue_num = response.json()["number"]
                return True, f"OK Issue #{issue_num}: {title}", issue_num
            error = response.json().get("message", "Unknown error")
            return False, f"ERROR Failed to create issue: {error}", None
        except Exception as exc:
            return False, f"ERROR Error creating issue: {exc}", None

    def create_discussion(self, owner: str, repo: str, category_id: str, title: str, body: str = "") -> Tuple[bool, str]:
        url = f"{API_BASE}/graphql"
        headers = dict(self.session.headers)
        headers["Accept"] = "application/vnd.github.v3+json"
        query = """
        mutation($repositoryId:ID!,$categoryId:ID!,$title:String!,$body:String!) {
            createDiscussion(input:{repositoryId:$repositoryId,categoryId:$categoryId,title:$title,body:$body}) {
                discussion {
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
            response = self.session.post(url, json={"query": query, "variables": variables}, headers=headers)
            if response.status_code == 200:
                data = response.json()
                if "errors" in data:
                    return False, f"ERROR Failed to create discussion: {data['errors'][0]['message']}"
                discussion_url = data["data"]["createDiscussion"]["discussion"]["url"]
                return True, f"OK Discussion created: {discussion_url}"
            return False, f"ERROR Failed to create discussion: {response.status_code}"
        except Exception as exc:
            return False, f"ERROR Error creating discussion: {exc}"


class GitHubAutomation:
    """Run config-driven GitHub setup tasks."""

    def __init__(self, config: Dict[str, Any], token: str, dry_run: bool = False):
        self.config = config
        self.dry_run = dry_run
        self.api = GitHubAPI(token, dry_run=dry_run)

    def validate_config(self) -> bool:
        for field in ("owner", "repo"):
            if field not in self.config:
                print(f"ERROR Config error: Missing required field '{field}'")
                return False
        return True

    def run(
        self,
        labels_only: bool = False,
        milestones_only: bool = False,
        issues_only: bool = False,
        discussions_only: bool = False,
    ) -> bool:
        owner = self.config["owner"]
        repo = self.config["repo"]

        print(f"\n{'=' * 70}")
        prefix = "[DRY-RUN] " if self.dry_run else ""
        print(f"  {prefix}GitHub Automation for {owner}/{repo}")
        print(f"{'=' * 70}\n")

        print("Step 1: Verifying GitHub API connection...")
        if not self.api.verify_user():
            print("\nERROR Authentication failed. Check your token.")
            return False
        if not self.api.verify_connection(owner, repo):
            print(f"\nERROR Cannot access repo {owner}/{repo}. Check permissions.")
            return False
        print()

        if not issues_only and not milestones_only and not discussions_only:
            self._create_labels(owner, repo)

        milestone_map: Dict[str, int] = {}
        if not labels_only and not issues_only and not discussions_only:
            milestone_map = self._create_milestones(owner, repo)
        elif issues_only:
            milestone_map = self.api.get_milestones(owner, repo)

        if not labels_only and not milestones_only and not discussions_only:
            self._create_issues(owner, repo, milestone_map)

        if not labels_only and not milestones_only and not issues_only:
            self._create_discussions(owner, repo)

        print("=" * 70)
        if self.dry_run:
            print("  OK Dry-run completed successfully.")
            print("  Run without --dry-run to actually create resources.")
        else:
            print("  OK Automation completed successfully.")
            print(f"  Visit: https://github.com/{owner}/{repo}/issues")
        print("=" * 70)
        print()
        return True

    def _create_labels(self, owner: str, repo: str) -> None:
        labels = self.config.get("labels", [])
        if not labels:
            return

        print("Step 2: Creating labels...")
        print(f"  Total labels to create: {len(labels)}\n")
        created = 0
        for label in labels:
            name = label.get("name")
            color = label.get("color")
            description = label.get("description", "")
            if not name or not color:
                print("  WARN Skipping invalid label (missing name or color)")
                continue
            success, message = self.api.create_label(owner, repo, name, color, description)
            print(f"  {message}")
            if success:
                created += 1
        print(f"\n  Summary: {created}/{len(labels)} created/verified\n")

    def _create_milestones(self, owner: str, repo: str) -> Dict[str, int]:
        milestones = self.config.get("milestones", [])
        existing = self.api.get_milestones(owner, repo)
        if not milestones:
            return existing

        print("Step 3: Creating milestones...")
        print(f"  Total milestones to create: {len(milestones)}\n")
        created = 0
        for milestone in milestones:
            title = milestone.get("title")
            description = milestone.get("description", "")
            due_on = milestone.get("due_on")
            state = milestone.get("state", "open")
            if not title:
                print("  WARN Skipping invalid milestone (missing title)")
                continue
            success, message, number = self.api.create_milestone(owner, repo, title, description, due_on, state)
            print(f"  {message}")
            if success:
                created += 1
                if number is not None:
                    existing[title] = number

        if not self.dry_run:
            existing = self.api.get_milestones(owner, repo)
        print(f"\n  Summary: {created}/{len(milestones)} created/verified\n")
        return existing

    def _create_issues(self, owner: str, repo: str, milestone_map: Dict[str, int]) -> None:
        issues = self.config.get("issues", [])
        if not issues:
            return

        print("Step 4: Creating issues...")
        print(f"  Total issues to create: {len(issues)}\n")
        created = 0
        for idx, issue in enumerate(issues, start=1):
            title = issue.get("title")
            body = issue.get("body", "")
            labels = issue.get("labels", [])
            assignees = issue.get("assignees", [])
            milestone_value = issue.get("milestone")
            milestone_number = None
            if isinstance(milestone_value, int):
                milestone_number = milestone_value
            elif isinstance(milestone_value, str) and milestone_value:
                milestone_number = milestone_map.get(milestone_value)

            if not title:
                print(f"  WARN Issue #{idx}: Missing title, skipping")
                continue

            print(f"  Creating issue #{idx}: {title}")
            success, message, _ = self.api.create_issue(
                owner, repo, title, body, labels, assignees, milestone_number
            )
            print(f"    {message}")
            if success:
                created += 1
        print(f"\n  Summary: {created}/{len(issues)} created\n")

    def _create_discussions(self, owner: str, repo: str) -> None:
        discussions = self.config.get("discussions", [])
        if not discussions:
            return

        print("Step 5: Creating discussions...")
        print(f"  Total discussions to create: {len(discussions)}\n")
        created = 0
        for idx, discussion in enumerate(discussions, start=1):
            category_id = discussion.get("category_id")
            title = discussion.get("title")
            body = discussion.get("body", "")
            if not category_id or not title:
                print(f"  WARN Discussion #{idx}: Missing category_id or title, skipping")
                continue
            print(f"  Creating discussion #{idx}: {title}")
            success, message = self.api.create_discussion(owner, repo, category_id, title, body)
            print(f"    {message}")
            if success:
                created += 1
        print(f"\n  Summary: {created}/{len(discussions)} created\n")


def main() -> None:
    parser = argparse.ArgumentParser(description="Universal GitHub Automation Tool")
    parser.add_argument("--config", help="Config file path (JSON or YAML). Auto-detected if not provided.")
    parser.add_argument("--token", required=True, help="GitHub Personal Access Token")
    parser.add_argument("--dry-run", action="store_true", help="Preview changes without creating anything")
    parser.add_argument("--labels-only", action="store_true", help="Only create labels")
    parser.add_argument("--milestones-only", action="store_true", help="Only create milestones")
    parser.add_argument("--issues-only", action="store_true", help="Only create issues")
    parser.add_argument("--discussions-only", action="store_true", help="Only create discussions")
    args = parser.parse_args()

    config_path = args.config or ConfigLoader.find_default()
    if not config_path:
        print("ERROR No config file found. Provide with --config")
        sys.exit(1)

    try:
        config = ConfigLoader.load(config_path)
    except Exception as exc:
        print(f"ERROR Config error: {exc}")
        sys.exit(1)

    automation = GitHubAutomation(config, args.token, dry_run=args.dry_run)
    if not automation.validate_config():
        sys.exit(1)

    success = automation.run(
        labels_only=args.labels_only,
        milestones_only=args.milestones_only,
        issues_only=args.issues_only,
        discussions_only=args.discussions_only,
    )
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
