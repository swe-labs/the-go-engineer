#!/usr/bin/env python3
import argparse, json, re, sys
from pathlib import Path

REQUIRED_README_HEADINGS = [
    "## Mission", "## Prerequisites", "## Mental Model", "## Visual Model",
    "## Machine View", "## Run Instructions", "## Try It", "## In Production",
    "## Thinking Questions", "## Next Step",
]
FORBIDDEN_TOKENS = ["todo", "tbd", "placeholder", "lorem ipsum", "coming soon"]

def load_json(path):
    with open(path, encoding="utf-8") as f:
        return json.load(f)

def load_metadata(metadata_dir):
    metadata_dir = Path(metadata_dir)
    core = load_json(metadata_dir / "path.core.json")
    electives = load_json(metadata_dir / "path.electives.json")
    projects = load_json(metadata_dir / "projects.json")
    assessments = load_json(metadata_dir / "assessments.json")
    return core, electives, projects, assessments

def all_items(metadata_dir):
    core, electives, _, _ = load_metadata(metadata_dir)
    return core.get("items", []) + electives.get("items", [])

def fail(errors):
    if errors:
        for error in errors:
            print(f"error: {error}")
        sys.exit(1)
    print("VALIDATION PASSED")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--metadata-dir", default="metadata")
    parser.add_argument("--curriculum-dir", default="curriculum")
    args = parser.parse_args()
    root = Path.cwd()
    errors = []
    if (root / "AGENTS.md").exists():
        errors.append("AGENTS.md must not exist at repository root")
    if (root / ".env").exists():
        errors.append(".env must not be committed")
    for forbidden in ["codex", "ai", "agent", "bot", "llm"]:
        for path in root.rglob(forbidden):
            if path.is_dir():
                errors.append(f"forbidden folder name: {path}")
    for item in all_items(args.metadata_dir):
        files = item.get("files") or {}
        for key, rel in files.items():
            if isinstance(rel, list):
                vals = rel
            else:
                vals = [rel]
            for value in vals:
                if not value:
                    continue
                if not str(value).startswith("curriculum/"):
                    errors.append(f"{item['id']}: {key} is not under curriculum/: {value}")
    fail(errors)

if __name__ == "__main__":
    main()
