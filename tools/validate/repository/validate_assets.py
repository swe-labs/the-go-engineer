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

def all_items(metadata_dir, strict=True, root=None):
    core, electives, _, _ = load_metadata(metadata_dir)
    modules_path = {}
    for m in core.get("modules", []):
        modules_path[m["id"]] = m.get("path", "")
    for m in electives.get("modules", []):
        modules_path[m["id"]] = m.get("path", "")
    for items in [core.get("items", []), electives.get("items", [])]:
        for item in items:
            if not strict and root and item.get("module_id") in modules_path:
                mod_path = root / modules_path[item["module_id"]]
                if not mod_path.exists():
                    continue
            yield item

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
    parser.add_argument("--strict", action=argparse.BooleanOptionalAction, default=True)
    args = parser.parse_args()
    root = Path.cwd()
    errors = []
    for item in all_items(args.metadata_dir, strict=args.strict, root=root):
        files = item.get("files") or {}
        assets_dir = files.get("assets_dir")
        if assets_dir and not (root / assets_dir).exists():
            errors.append(f"{item['id']}: assets_dir missing at {assets_dir}")
        for diagram in files.get("diagram_paths") or []:
            if not (root / diagram).exists():
                errors.append(f"{item['id']}: diagram missing at {diagram}")
    fail(errors)

if __name__ == "__main__":
    main()
