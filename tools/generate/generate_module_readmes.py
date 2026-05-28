#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(path): return json.load(open(path, encoding="utf-8"))

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--only-missing", action="store_true"); args=p.parse_args()
    bundles=[load(Path(args.metadata_dir)/"path.core.json"), load(Path(args.metadata_dir)/"path.electives.json")]
    count=0
    for bundle in bundles:
        by_module={m["id"]: [] for m in bundle.get("modules", [])}
        for item in bundle.get("items", []): by_module.setdefault(item["module_id"], []).append(item)
        for m in bundle.get("modules", []):
            path=Path(m["path"])/"README.md"
            if args.only_missing and path.exists(): continue
            path.parent.mkdir(parents=True, exist_ok=True)
            lines=[f"# {m['title']}", "", "## Goal", "", m.get("learning_goal", ""), "", "## Summary", "", m.get("summary", ""), "", "## Lessons", ""]
            for item in sorted(by_module.get(m["id"], []), key=lambda x:x.get("order",0)):
                lines.append(f"- `{item['id']}` {item['title']}")
            lines += ["", "## Completion", "", "Complete every required lesson, lab, project, and assessment for this module.", ""]
            path.write_text("\n".join(lines), encoding="utf-8")
            count+=1
    print(f"generated module READMEs: {count}")
if __name__ == "__main__": main()
