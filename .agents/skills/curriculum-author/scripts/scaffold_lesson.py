#!/usr/bin/env python3
"""Create missing learner-facing curriculum files from metadata.

The script creates high-quality scaffolds, not final lesson content. It uses the
final typed curriculum layout:

  curriculum/modules/{module}/{lessons|labs|projects|assessments}/{item}/
  curriculum/electives/{elective}/{lessons|labs|projects|assessments}/{item}/

Existing files are not overwritten unless --force is used.
"""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path
from typing import Any, Dict, List, Tuple


def find_metadata_dir(repo: Path) -> Path:
    for rel in ["metadata", "."]:
        candidate = repo / rel
        if (candidate / "path.core.json").exists() and (candidate / "path.electives.json").exists():
            return candidate
    legacy = repo / "curriculum"
    if (legacy / "path.core.json").exists() and (legacy / "path.electives.json").exists():
        return legacy
    raise SystemExit("cannot locate metadata directory containing path.core.json and path.electives.json")


def load_items_and_modules(metadata_dir: Path) -> Tuple[Dict[str, Dict[str, Any]], Dict[str, Dict[str, Any]]]:
    items: Dict[str, Dict[str, Any]] = {}
    modules: Dict[str, Dict[str, Any]] = {}
    for name in ["path.core.json", "path.electives.json"]:
        data = json.loads((metadata_dir / name).read_text(encoding="utf-8"))
        for module in data.get("modules", []):
            modules[module["id"]] = module
        for item in data.get("items", []):
            items[item["id"]] = item
    return items, modules


def slugify(text: str) -> str:
    words = re.findall(r"[a-z0-9]+", text.lower())
    return "-".join(words) or "item"


def item_folder_name(item: Dict[str, Any]) -> str:
    order = item.get("order")
    slug = item.get("slug") or slugify(item.get("title", item.get("id", "item")))
    if isinstance(order, int) and order > 0:
        return f"{order:02d}-{slug}"
    return slug


def type_folder(item: Dict[str, Any]) -> str:
    item_type = str(item.get("type", "lesson")).lower()
    subtype = str(item.get("subtype", "")).lower()
    if item_type == "lab" or subtype == "lab":
        return "labs"
    if item_type == "project" or subtype == "project":
        return "projects"
    if item_type in {"assessment", "checkpoint"} or subtype in {"assessment", "checkpoint"}:
        return "assessments"
    return "lessons"


def canonical_module_base(module: Dict[str, Any]) -> str:
    path = str(module.get("path") or module.get("slug") or "")
    slug = module.get("slug") or path.split("/")[-1]
    number = module.get("number")
    if isinstance(number, int) and not re.match(r"^\d{2}-", str(slug)):
        module_dir = f"{number:02d}-{slug}"
    else:
        module_dir = str(slug)
    if str(module.get("phase")) == "elective" or str(module.get("id")) == "module-17":
        return f"curriculum/electives/{slug if slug != 'advanced-electives' else 'advanced-electives'}"
    if path.startswith("curriculum/modules/") or path.startswith("curriculum/electives/"):
        return path.rstrip("/")
    if path.startswith("modules/"):
        return f"curriculum/{path}".rstrip("/")
    if path.startswith("electives/"):
        return f"curriculum/{path}".rstrip("/")
    if re.match(r"^\d{2}-", path):
        return f"curriculum/modules/{path}".rstrip("/")
    return f"curriculum/modules/{module_dir}".rstrip("/")


def canonical_files(item: Dict[str, Any], modules: Dict[str, Dict[str, Any]]) -> Dict[str, str]:
    files = dict(item.get("files") or {})
    module = modules.get(item.get("module_id"), {})
    base = f"{canonical_module_base(module)}/{type_folder(item)}/{item_folder_name(item)}"
    defaults = {
        "readme_path": f"{base}/README.md",
        "main_path": f"{base}/main.go",
        "test_path": f"{base}/main_test.go",
        "starter_path": f"{base}/_starter",
        "solution_path": f"{base}/_solution",
        "assets_dir": f"{base}/assets",
    }
    # If metadata paths are already canonical, respect them. Otherwise derive from final layout.
    result: Dict[str, str] = {}
    for key, default in defaults.items():
        raw = str(files.get(key) or "")
        if raw.startswith("curriculum/modules/") or raw.startswith("curriculum/electives/"):
            result[key] = raw
        elif raw:
            result[key] = default
        else:
            result[key] = default
    return result


def write_if_allowed(path: Path, content: str, force: bool) -> bool:
    path.parent.mkdir(parents=True, exist_ok=True)
    if path.exists() and not force:
        return False
    path.write_text(content, encoding="utf-8")
    return True


def title_to_identifier(title: str) -> str:
    words = re.findall(r"[A-Za-z0-9]+", title)
    if not words:
        return "Lesson"
    return "".join(w.capitalize() for w in words[:4])


def bullet_list(values: List[Any], fallback: str) -> str:
    if not values:
        return f"- {fallback}"
    return "\n".join(f"- {v}" for v in values)


def readme_content(item: Dict[str, Any]) -> str:
    zm = item.get("zero_magic") or {}
    next_ids = ", ".join(item.get("next_item_ids") or []) or "End of module"
    run_cmd = (item.get("verification") or {}).get("run_command") or "go run ."
    test_cmd = (item.get("verification") or {}).get("test_command") or "go test ./..."
    mistakes = zm.get("beginner_mistakes") or []
    failure_modes = zm.get("failure_modes") or zm.get("operational_failure_examples") or []
    mistake_md = bullet_list(mistakes[:6], "Add specific mistakes after implementing the lesson.")
    failure_md_lines: List[str] = []
    for fm in failure_modes[:4]:
        if isinstance(fm, dict):
            scenario = fm.get("scenario") or fm.get("failure") or "Failure scenario"
            cause = fm.get("cause") or "Cause to explain."
            fix = fm.get("fix") or fm.get("mitigation") or "Fix to implement."
            failure_md_lines.append(f"- **Scenario:** {scenario}\n  - Cause: {cause}\n  - Fix: {fix}")
        else:
            failure_md_lines.append(f"- {fm}")
    failure_md = "\n".join(failure_md_lines) or "- Add a concrete failure mode after implementation."
    steps = zm.get("step_by_step_execution") or zm.get("execution_timeline") or []
    steps_md = "\n".join(f"{i+1}. {step}" for i, step in enumerate(steps)) or "1. Add a concrete execution trace after implementation."
    perf_md = bullet_list(zm.get("performance_implications") or [], "State the performance, memory, reliability, or operational tradeoff after implementation.")
    return f"""# {item.get('title')}

## Learning objective

{item.get('learning_objective', '').strip()}

## Why this matters

{zm.get('problem_solved', '').strip()}

{zm.get('why_it_exists', '').strip()}

## Mental model

{zm.get('mental_model', '').strip()}

State where this analogy stops being true before moving on.

## Core idea

Explain the concept in plain language first, then connect it to the implementation below. Keep one primary idea per section and introduce vocabulary before using it.

## Under the hood

{zm.get('under_the_hood', '').strip()}

## How Go uses it

{zm.get('how_go_uses_it', '').strip()}

## Go example

```go
// See main.go for the runnable version of this example.
```

## Step-by-step execution

{steps_md}

## Common mistakes

{mistake_md}

## Debugging walkthrough

Use the failure modes below to debug from symptom to root cause. Do not jump straight to the fix; first identify the invariant that broke.

{failure_md}

## Production notes

{zm.get('real_world_usage', '').strip()}

## Performance implications

{perf_md}

## Practice task

{(item.get('proof') or {}).get('practice_task', zm.get('proof_of_understanding', '')).strip()}

## Tests / verification

```bash
{run_cmd}
{test_cmd}
```

## Review questions

1. Explain the concept without using the lesson wording.
2. Identify one beginner mistake and explain why it fails.
3. Debug a broken version of the example and state the root cause.
4. Describe where this appears in production Go software.

## NEXT UP

{next_ids}
"""


def go_main_content(item: Dict[str, Any]) -> str:
    ident = title_to_identifier(item.get("title", "Lesson"))
    run_cmd = (item.get("verification") or {}).get("run_command") or "go run ."
    return f'''package main

import "fmt"

// RUN: {run_cmd}

func main() {{
	fmt.Println("{item.get('title', 'lesson')}")
	fmt.Println(example{ident}())
}}

func example{ident}() string {{
	return "replace this scaffold with a minimal, focused lesson example"
}}
'''


def go_test_content(item: Dict[str, Any]) -> str:
    ident = title_to_identifier(item.get("title", "Lesson"))
    return f'''package main

import "testing"

func TestExample{ident}(t *testing.T) {{
	got := example{ident}()
	if got == "" {{
		t.Fatal("expected non-empty example output")
	}}
}}
'''


def scaffold(repo: Path, item: Dict[str, Any], modules: Dict[str, Dict[str, Any]], force: bool) -> List[str]:
    files = canonical_files(item, modules)
    created: List[str] = []
    contract = item.get("content_contract") or {}
    if files.get("readme_path"):
        path = repo / files["readme_path"]
        if write_if_allowed(path, readme_content(item), force):
            created.append(str(path.relative_to(repo)))
    if contract.get("runnable_required") or (item.get("files") or {}).get("main_path"):
        path = repo / files["main_path"]
        if write_if_allowed(path, go_main_content(item), force):
            created.append(str(path.relative_to(repo)))
    if contract.get("tests_required") or (item.get("files") or {}).get("test_path"):
        path = repo / files["test_path"]
        if write_if_allowed(path, go_test_content(item), force):
            created.append(str(path.relative_to(repo)))
    for key in ["starter_path", "solution_path", "assets_dir"]:
        rel = files.get(key)
        if rel:
            path = repo / rel
            path.mkdir(parents=True, exist_ok=True)
            marker = path / ".gitkeep"
            if not marker.exists():
                marker.write_text("", encoding="utf-8")
            created.append(str(path.relative_to(repo)) + "/")
    return created


def main() -> None:
    parser = argparse.ArgumentParser(description="Scaffold learner-facing curriculum files from metadata")
    parser.add_argument("--repo", default=".")
    parser.add_argument("--item-id", required=True)
    parser.add_argument("--force", action="store_true", help="overwrite existing files")
    parser.add_argument("--overwrite-missing-only", action="store_true", help="kept for readability; this is the default")
    args = parser.parse_args()

    repo = Path(args.repo).resolve()
    metadata_dir = find_metadata_dir(repo)
    items, modules = load_items_and_modules(metadata_dir)
    item = items.get(args.item_id)
    if not item:
        raise SystemExit(f"unknown item id: {args.item_id}")
    created = scaffold(repo, item, modules, args.force)
    print(json.dumps({"item_id": args.item_id, "created_or_touched": created}, indent=2))


if __name__ == "__main__":
    main()
