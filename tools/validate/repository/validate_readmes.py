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


def check_readme(path, owner):
    errors = []
    text = path.read_text(encoding="utf-8")
    last = -1
    for heading in REQUIRED_README_HEADINGS:
        idx = text.find(heading)
        if idx < 0:
            errors.append(f"{owner}: missing {heading}")
        elif idx < last:
            errors.append(f"{owner}: heading out of order {heading}")
        else:
            last = idx
    lower = text.lower()
    for token in FORBIDDEN_TOKENS:
        if token in lower:
            errors.append(f"{owner}: README contains forbidden token {token!r}")
    if "```" not in text:
        errors.append(f"{owner}: README should include a fenced command/code block")
    return errors

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--metadata-dir", default="metadata")
    parser.add_argument("--curriculum-dir", default="curriculum")
    args = parser.parse_args()
    root = Path.cwd()
    errors = []
    for item in all_items(args.metadata_dir):
        rel = (item.get("files") or {}).get("readme_path")
        if not rel:
            errors.append(f"{item['id']}: missing files.readme_path")
            continue
        path = root / rel
        if not path.exists():
            errors.append(f"{item['id']}: README missing at {rel}")
            continue
        errors.extend(check_readme(path, item["id"]))
    fail(errors)

if __name__ == "__main__":
    main()
