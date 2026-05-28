#!/usr/bin/env python3
"""Audit zero-magic curriculum metadata and optional learner-facing repository files.

This script is dependency-free. It validates metadata references, detects placeholders,
checks final path architecture, and optionally verifies real README/code/test/asset files.
"""

from __future__ import annotations

import argparse
import json
import re
from dataclasses import dataclass, asdict
from pathlib import Path
from typing import Any, Dict, Iterable, List, Optional, Tuple

METADATA_FILES = [
    "path.core.json",
    "path.electives.json",
    "projects.json",
    "assessments.json",
    "crossrefs.json",
    "concepts.json",
    "failures.json",
    "readme.contracts.json",
    "migration.v2-to-v3.json",
    "workspace.json",
]

REQUIRED_TOP_LEVEL_DIRS = ["metadata", "curriculum", "tools", "docs"]
RECOMMENDED_TOOLS_DIRS = ["validate", "generate", "audit", "migrate", "authoring"]
FORBIDDEN_DIR_NAMES = {"codex", "ai", "chatgpt"}

REQUIRED_ZM_TEXT = [
    "problem_solved",
    "why_it_exists",
    "mental_model",
    "under_the_hood",
    "how_go_uses_it",
    "real_world_usage",
    "proof_of_understanding",
]
REQUIRED_ZM_ARRAYS = ["beginner_mistakes", "failure_modes"]
GENERIC_REASON_PATTERNS = [
    r"both .* build understanding within",
    r"related concept:\s*$",
    r"requires understanding of\s*\.?$",
    r"this module builds on concepts from the previous module",
    r"hands-on practice\.?$",
]
PLACEHOLDER_PATTERNS = [
    r"\[TODO[:\]]",
    r"\bTBD[:\]]",
    r"lorem ipsum",
    r"this lesson explains what problem",
    r"one step in the learner'?s path",
    r"mechanically without understanding",
]
README_REQUIRED_HEADINGS = [
    "learning objective",
    "why this matters",
    "mental model",
    "under the hood",
    "how go uses it",
    "common mistakes",
    "debugging walkthrough",
    "production notes",
    "practice",
    "review questions",
]


@dataclass
class Finding:
    severity: str
    code: str
    target: str
    message: str
    path: Optional[str] = None


@dataclass
class AuditResult:
    metadata_dir: str
    repo_root: str
    strict_repository: bool
    counts: Dict[str, int]
    findings: List[Finding]

    @property
    def errors(self) -> int:
        return sum(1 for f in self.findings if f.severity == "error")

    @property
    def warnings(self) -> int:
        return sum(1 for f in self.findings if f.severity == "warning")


def find_metadata_dir(repo: Path) -> Path:
    for rel in ["metadata", "."]:
        candidate = repo / rel
        if (candidate / "path.core.json").exists() and (candidate / "path.electives.json").exists():
            return candidate
    legacy = repo / "curriculum"
    if (legacy / "path.core.json").exists() and (legacy / "path.electives.json").exists():
        return legacy
    raise SystemExit("cannot locate metadata directory containing path.core.json and path.electives.json")


def load_json(path: Path, findings: List[Finding]) -> Any:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError:
        findings.append(Finding("error", "missing_metadata_file", path.name, f"missing {path.name}", str(path)))
        return None
    except json.JSONDecodeError as exc:
        findings.append(Finding("error", "invalid_json", path.name, f"invalid JSON: {exc}", str(path)))
        return None


def as_list(value: Any) -> List[Any]:
    return value if isinstance(value, list) else []


def text_has_pattern(text: str, patterns: Iterable[str]) -> Optional[str]:
    for pat in patterns:
        if re.search(pat, text, flags=re.I):
            return pat
    return None


def add(findings: List[Finding], severity: str, code: str, target: str, message: str, path: Optional[str] = None) -> None:
    findings.append(Finding(severity, code, target, message, path))


def collect(metadata: Dict[str, Any]) -> Tuple[List[Dict[str, Any]], List[Dict[str, Any]], List[Dict[str, Any]], List[Dict[str, Any]]]:
    core = metadata.get("path.core.json") or {}
    electives = metadata.get("path.electives.json") or {}
    modules = as_list(core.get("modules")) + as_list(electives.get("modules"))
    items = as_list(core.get("items")) + as_list(electives.get("items"))
    projects = as_list((metadata.get("projects.json") or {}).get("projects"))
    assessments = as_list((metadata.get("assessments.json") or {}).get("assessments"))
    return modules, items, projects, assessments


def item_kind_folder(item: Dict[str, Any]) -> str:
    item_type = str(item.get("type", "lesson")).lower()
    subtype = str(item.get("subtype", "")).lower()
    if item_type == "lab" or subtype == "lab":
        return "labs"
    if item_type == "project" or subtype == "project":
        return "projects"
    if item_type in {"assessment", "checkpoint"} or subtype in {"assessment", "checkpoint"}:
        return "assessments"
    return "lessons"


def is_canonical_content_path(rel: str, item: Optional[Dict[str, Any]] = None) -> bool:
    if not rel:
        return True
    allowed_prefix = rel.startswith("curriculum/modules/") or rel.startswith("curriculum/electives/")
    if not allowed_prefix:
        return False
    if item is None:
        return True
    folder = item_kind_folder(item)
    return f"/{folder}/" in rel or rel.endswith(f"/{folder}")


def is_canonical_module_path(path: str, elective: bool = False) -> bool:
    if elective:
        return path.startswith("curriculum/electives/")
    return path.startswith("curriculum/modules/")


def audit_repo_architecture(repo: Path, findings: List[Finding], strict_repository: bool) -> None:
    for d in REQUIRED_TOP_LEVEL_DIRS:
        if not (repo / d).exists():
            sev = "error" if strict_repository else "warning"
            add(findings, sev, "missing_top_level_dir", d, f"recommended top-level directory '{d}/' is missing", d)
    for bad in FORBIDDEN_DIR_NAMES:
        matches = [p for p in repo.rglob(bad) if p.is_dir() and ".git" not in p.parts]
        for path in matches:
            add(findings, "error", "forbidden_tool_branded_dir", bad, f"folder name '{bad}' is not allowed; use responsibility-based names", str(path.relative_to(repo)))
    tools = repo / "tools"
    if tools.exists():
        for d in RECOMMENDED_TOOLS_DIRS:
            if not (tools / d).exists():
                add(findings, "warning", "missing_tools_subdir", d, f"tools/{d}/ is recommended for maintainability", f"tools/{d}")


def audit_metadata(metadata_dir: Path, repo: Path, strict_repository: bool) -> AuditResult:
    findings: List[Finding] = []
    metadata: Dict[str, Any] = {}
    for name in METADATA_FILES:
        path = metadata_dir / name
        data = load_json(path, findings)
        if data is not None:
            metadata[name] = data

    audit_repo_architecture(repo, findings, strict_repository)

    modules, items, projects, assessments = collect(metadata)
    module_ids = {m.get("id") for m in modules if m.get("id")}
    item_ids = {i.get("id") for i in items if i.get("id")}
    project_ids = {p.get("id") for p in projects if p.get("id")}
    assessment_ids = {a.get("id") for a in assessments if a.get("id")}
    concept_names = {c.get("concept") for c in as_list((metadata.get("concepts.json") or {}).get("concepts")) if c.get("concept")}
    valid_targets = set(module_ids) | set(item_ids) | set(project_ids) | set(assessment_ids) | concept_names

    for label, records in [("module", modules), ("item", items), ("project", projects), ("assessment", assessments)]:
        seen: Dict[str, str] = {}
        for rec in records:
            rid = rec.get("id")
            if not rid:
                add(findings, "error", "missing_id", label, f"{label} missing id")
            elif rid in seen:
                add(findings, "error", "duplicate_id", rid, f"duplicate {label} id")
            else:
                seen[rid] = label

    items_by_module: Dict[str, List[Dict[str, Any]]] = {}
    for item in items:
        items_by_module.setdefault(item.get("module_id"), []).append(item)

    for mod in modules:
        mid = mod.get("id", "<missing>")
        elective = mod.get("phase") == "elective" or mid == "module-17"
        path = str(mod.get("path") or "")
        if path and not is_canonical_module_path(path, elective=elective):
            add(findings, "error", "noncanonical_module_path", mid, "module path must start with curriculum/modules/ or curriculum/electives/", path)
        for prereq in as_list(mod.get("prerequisites")):
            if prereq not in module_ids:
                add(findings, "error", "bad_module_prerequisite", mid, f"unknown module prerequisite {prereq}")
        for eid in as_list(mod.get("entry_item_ids")) + as_list(mod.get("terminal_item_ids")):
            if eid not in item_ids:
                add(findings, "error", "bad_module_boundary", mid, f"boundary item {eid} does not exist")
        count = len(items_by_module.get(mid, []))
        if mod.get("required", False) and count == 0:
            add(findings, "error", "empty_required_module", mid, "required module has no items")
        if count > 30:
            add(findings, "warning", "oversized_module", mid, f"module has {count} items; consider splitting")

    for item in items:
        iid = item.get("id", "<missing>")
        if item.get("module_id") not in module_ids:
            add(findings, "error", "bad_item_module", iid, f"unknown module_id {item.get('module_id')}")
        for prereq in as_list(item.get("prerequisites")):
            if prereq not in item_ids:
                add(findings, "error", "bad_item_prerequisite", iid, f"unknown prerequisite {prereq}")
        for nid in as_list(item.get("next_item_ids")):
            if nid not in item_ids:
                add(findings, "error", "bad_next_item", iid, f"unknown next_item_id {nid}")
        if item.get("zero_magic_status") != "golden":
            add(findings, "error", "non_golden_zero_magic_status", iid, f"zero_magic_status is {item.get('zero_magic_status')}")
        if item.get("readme_status") != "golden":
            add(findings, "error", "non_golden_readme_status", iid, f"readme_status is {item.get('readme_status')}")
        zm = item.get("zero_magic") or {}
        for key in REQUIRED_ZM_TEXT:
            val = zm.get(key)
            if not isinstance(val, str) or not val.strip():
                add(findings, "error", "missing_zero_magic_field", iid, f"zero_magic.{key} is missing or empty")
        if not any(isinstance(zm.get(key), list) and len(zm.get(key)) > 0 for key in ["failure_modes", "operational_failure_examples"]):
            add(findings, "error", "missing_failure_coverage", iid, "zero_magic needs failure_modes or operational_failure_examples")
        if not isinstance(zm.get("beginner_mistakes"), list) or len(zm.get("beginner_mistakes")) == 0:
            add(findings, "error", "missing_zero_magic_array", iid, "zero_magic.beginner_mistakes is missing or empty")
        raw_item = json.dumps(item, ensure_ascii=False)
        pat = text_has_pattern(raw_item, PLACEHOLDER_PATTERNS)
        if pat:
            add(findings, "error", "placeholder_text", iid, f"placeholder-like text matched pattern {pat}")
        proof = item.get("proof") or {}
        if proof.get("assessment_id") and proof.get("assessment_id") not in assessment_ids:
            add(findings, "error", "bad_item_assessment", iid, f"unknown assessment_id {proof.get('assessment_id')}")
        for key, rel in (item.get("files") or {}).items():
            if key.endswith("path") or key.endswith("dir"):
                if rel and not is_canonical_content_path(str(rel), item):
                    add(findings, "error", "noncanonical_content_path", iid, f"{key} must use typed curriculum/modules or curriculum/electives layout", str(rel))
        for group in ["builds_on", "preview_only", "related", "reinforced_in"]:
            for ref in as_list((item.get("crossrefs") or {}).get(group)):
                tid = ref.get("target_id")
                if tid and tid not in valid_targets:
                    add(findings, "error", "bad_item_crossref", iid, f"{group} target {tid} does not exist")
                reason = str(ref.get("reason", ""))
                pat = text_has_pattern(reason, GENERIC_REASON_PATTERNS)
                if pat:
                    add(findings, "error", "generic_crossref_reason", iid, f"{group} reason is generic: {reason}")

    for ass in assessments:
        aid = ass.get("id", "<missing>")
        for tid in as_list(ass.get("target_ids")):
            if tid not in valid_targets:
                add(findings, "error", "bad_assessment_target", aid, f"target_id {tid} does not exist")
        criteria = as_list(ass.get("criteria"))
        if not criteria:
            add(findings, "error", "missing_assessment_criteria", aid, "assessment has no criteria")
        weights = [c.get("weight") for c in criteria if isinstance(c, dict) and isinstance(c.get("weight"), (int, float))]
        if weights and round(sum(weights), 2) not in (1.0, 100.0):
            add(findings, "warning", "assessment_weight_sum", aid, f"criteria weights sum to {sum(weights)}")

    for proj in projects:
        pid = proj.get("id", "<missing>")
        if proj.get("module_id") not in module_ids:
            add(findings, "error", "bad_project_module", pid, f"unknown module_id {proj.get('module_id')}")
        if proj.get("assessment_id") and proj.get("assessment_id") not in assessment_ids:
            add(findings, "error", "bad_project_assessment", pid, f"unknown assessment_id {proj.get('assessment_id')}")
        for key in ["prerequisites", "prerequisite_item_ids", "required_concepts", "reinforcement_targets"]:
            for tid in as_list(proj.get(key)):
                if tid not in valid_targets:
                    add(findings, "error", "bad_project_reference", pid, f"{key} references {tid}, which does not exist")
        anchor = proj.get("placement_anchor_item_id")
        if anchor and anchor not in item_ids:
            add(findings, "error", "bad_project_anchor", pid, f"placement_anchor_item_id {anchor} does not exist")
        if not as_list(proj.get("deliverables")):
            add(findings, "error", "missing_project_deliverables", pid, "project has no deliverables")

    for concept in as_list((metadata.get("concepts.json") or {}).get("concepts")):
        name = concept.get("concept", "<missing>")
        owner = concept.get("canonical_owner")
        if owner not in valid_targets:
            add(findings, "error", "bad_concept_owner", name, f"canonical_owner {owner} does not exist")
        if not as_list(concept.get("reinforcement_locations")):
            add(findings, "error", "unreinforced_concept", name, "concept has no reinforcement_locations")
        for key in ["preview_locations", "reinforcement_locations"]:
            for loc in as_list(concept.get(key)):
                if loc not in valid_targets:
                    add(findings, "error", "bad_concept_location", name, f"{key} references {loc}, which does not exist")

    for ref in as_list(((metadata.get("crossrefs.json") or {}).get("crossrefs") or {}).get("references")):
        fid = ref.get("from_id")
        tid = ref.get("target_id") or ref.get("to_id")
        if fid and fid not in valid_targets:
            add(findings, "error", "bad_global_crossref_from", str(fid), "from_id does not exist")
        if tid and tid not in valid_targets:
            add(findings, "error", "bad_global_crossref_to", str(tid), "target_id/to_id does not exist")
        reason = str(ref.get("reason", ""))
        pat = text_has_pattern(reason, GENERIC_REASON_PATTERNS)
        if pat:
            add(findings, "error", "generic_global_crossref_reason", f"{fid}->{tid}", f"generic reason: {reason}")

    if strict_repository:
        audit_repository_files(repo, modules, items, findings)

    counts = {
        "modules": len(modules),
        "items": len(items),
        "projects": len(projects),
        "assessments": len(assessments),
        "concepts": len(concept_names),
        "metadata_files_loaded": len(metadata),
    }
    return AuditResult(str(metadata_dir), str(repo), strict_repository, counts, findings)


def audit_repository_files(repo: Path, modules: List[Dict[str, Any]], items: List[Dict[str, Any]], findings: List[Finding]) -> None:
    for mod in modules:
        mid = mod.get("id", "<missing>")
        mpath = str(mod.get("path") or "")
        if mpath:
            readme = repo / mpath / "README.md"
            if not readme.exists():
                add(findings, "error", "missing_module_readme", mid, "module README.md is missing", str(readme.relative_to(repo)))
    for item in items:
        iid = item.get("id", "<missing>")
        files = item.get("files") or {}
        contract = item.get("content_contract") or {}
        required_paths: List[Tuple[str, Optional[str]]] = []
        if contract.get("readme_required") or files.get("readme_path"):
            required_paths.append(("readme", files.get("readme_path")))
        if contract.get("runnable_required") or files.get("main_path"):
            required_paths.append(("main", files.get("main_path")))
        if contract.get("tests_required") or files.get("test_path"):
            required_paths.append(("test", files.get("test_path")))
        if files.get("starter_path"):
            required_paths.append(("starter", files.get("starter_path")))
        if files.get("solution_path"):
            required_paths.append(("solution", files.get("solution_path")))
        if files.get("assets_dir") and contract.get("visual_model_required"):
            required_paths.append(("assets", files.get("assets_dir")))
        for kind, rel in required_paths:
            if not rel:
                add(findings, "error", "missing_file_path", iid, f"{kind} path is required but missing")
                continue
            path = repo / rel
            if not path.exists():
                add(findings, "error", "missing_repo_file", iid, f"required {kind} path does not exist", rel)
                continue
            if kind == "readme":
                audit_readme(iid, path, item, findings)
            elif kind == "test" and path.is_file():
                text = path.read_text(encoding="utf-8", errors="ignore")
                if "func Test" not in text:
                    add(findings, "error", "weak_test_file", iid, "test file contains no Go Test function", rel)
            elif kind == "main" and path.is_file() and path.suffix == ".go":
                text = path.read_text(encoding="utf-8", errors="ignore")
                if "package " not in text:
                    add(findings, "error", "weak_go_file", iid, "go file has no package declaration", rel)
        for diag in as_list(files.get("diagram_paths")):
            if not (repo / diag).exists():
                add(findings, "error", "missing_diagram", iid, f"diagram path does not exist: {diag}", diag)


def audit_readme(iid: str, path: Path, item: Dict[str, Any], findings: List[Finding]) -> None:
    text = path.read_text(encoding="utf-8", errors="ignore")
    lower = text.lower()
    if len(text.strip()) < 1200:
        add(findings, "error", "thin_readme", iid, "README is too short for world-class lesson content", str(path))
    for heading in README_REQUIRED_HEADINGS:
        if heading not in lower:
            add(findings, "error", "missing_readme_heading", iid, f"README missing section related to '{heading}'", str(path))
    if as_list(item.get("next_item_ids")) and "next up" not in lower:
        add(findings, "error", "missing_next_up", iid, "README missing NEXT UP footer", str(path))
    pat = text_has_pattern(text, PLACEHOLDER_PATTERNS)
    if pat:
        add(findings, "error", "placeholder_readme", iid, f"README contains placeholder-like text matching {pat}", str(path))
    if "```go" not in lower and (item.get("content_contract") or {}).get("runnable_required"):
        add(findings, "error", "missing_go_example", iid, "README lacks Go code block for runnable lesson", str(path))


def to_markdown(result: AuditResult) -> str:
    lines = [
        "# Curriculum Audit",
        "",
        f"Metadata directory: `{result.metadata_dir}`",
        f"Repository root: `{result.repo_root}`",
        f"Strict repository: `{result.strict_repository}`",
        "",
        "## Counts",
    ]
    for key, value in result.counts.items():
        lines.append(f"- {key}: {value}")
    lines += ["", "## Summary", f"- errors: {result.errors}", f"- warnings: {result.warnings}", ""]
    if not result.findings:
        lines.append("VALIDATION PASSED")
    else:
        lines.append("## Findings")
        for f in result.findings:
            path = f" (`{f.path}`)" if f.path else ""
            lines.append(f"- **{f.severity.upper()}** `{f.code}` `{f.target}`: {f.message}{path}")
    return "\n".join(lines) + "\n"


def main() -> None:
    parser = argparse.ArgumentParser(description="Audit a zero-magic curriculum repository")
    parser.add_argument("--repo", default=".", help="repository root")
    parser.add_argument("--strict-repository", action="store_true", help="also require actual README/code/test/asset files")
    parser.add_argument("--format", choices=["json", "markdown"], default="markdown")
    parser.add_argument("--out", help="write report to file")
    args = parser.parse_args()

    repo = Path(args.repo).resolve()
    metadata_dir = find_metadata_dir(repo)
    result = audit_metadata(metadata_dir, repo, args.strict_repository)

    if args.format == "json":
        payload = asdict(result)
        text = json.dumps(payload, indent=2)
    else:
        text = to_markdown(result)

    if args.out:
        Path(args.out).write_text(text, encoding="utf-8")
    else:
        print(text)

    raise SystemExit(1 if result.errors else 0)


if __name__ == "__main__":
    main()
