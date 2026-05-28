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
    for item in all_items(args.metadata_dir):
        files = item.get("files") or {}
        for key in ["main_path", "test_path"]:
            rel = files.get(key)
            if not rel:
                continue
            path = root / rel
            if not path.exists():
                errors.append(f"{item['id']}: missing {key} at {rel}")
                continue
            text = path.read_text(encoding="utf-8")
            if "package " not in text:
                errors.append(f"{item['id']}: {rel} missing package declaration")
            if key == "test_path" and not re.search(r"func\s+Test[A-Za-z0-9_]+\s*\(", text):
                errors.append(f"{item['id']}: {rel} has no Test* function")
            if "TODO" in text or "todo" in text.lower():
                errors.append(f"{item['id']}: {rel} contains TODO/todo")
    fail(errors)

if __name__ == "__main__":
    main()
