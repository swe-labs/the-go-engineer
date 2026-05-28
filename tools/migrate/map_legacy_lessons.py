#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(p): return json.load(open(p, encoding="utf-8"))
def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--output", default="dist/legacy-map.json"); args=p.parse_args()
    md=Path(args.metadata_dir); core=load(md/"path.core.json"); elec=load(md/"path.electives.json"); projects=load(md/"projects.json")
    mapping={}
    for item in core.get("items",[])+elec.get("items",[]):
        for legacy in item.get("source_legacy_ids") or []: mapping.setdefault(legacy, []).append(item["id"])
    for project in projects.get("projects",[]):
        for legacy in project.get("source_legacy_ids") or []: mapping.setdefault(legacy, []).append(project["id"])
    Path(args.output).parent.mkdir(parents=True, exist_ok=True); Path(args.output).write_text(json.dumps(mapping, indent=2)+"\n")
    print(f"mapped legacy ids: {len(mapping)}")
if __name__ == "__main__": main()
