#!/usr/bin/env python3
import argparse, json, shutil
from pathlib import Path

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--source-root", default="."); p.add_argument("--dry-run", action="store_true"); args=p.parse_args()
    md=Path(args.metadata_dir); core=json.load(open(md/"path.core.json", encoding="utf-8")); elec=json.load(open(md/"path.electives.json", encoding="utf-8"))
    moves=[]
    for item in core.get("items",[])+elec.get("items",[]):
        old=(item.get("source_legacy_ids") or [None])[0]
        new=(item.get("files") or {}).get("readme_path")
        if old and new: moves.append({"source_legacy_id":old,"target_readme":new})
    print(json.dumps({"planned_content_targets":len(moves),"moves":moves[:50]}, indent=2))
    if args.dry_run: return
    print("Content migration is intentionally conservative. Use --dry-run output as a review plan, then copy curated content manually or with a project-specific mapper.")
if __name__ == "__main__": main()
