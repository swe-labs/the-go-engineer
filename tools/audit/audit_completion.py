#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(p): return json.load(open(p, encoding="utf-8"))
def exists(rel): return bool(rel) and Path(rel).exists()

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--curriculum-dir", default="curriculum"); p.add_argument("--output", default="dist/completion-report.json"); args=p.parse_args()
    md=Path(args.metadata_dir); core=load(md/"path.core.json"); elec=load(md/"path.electives.json")
    items=core.get("items",[])+elec.get("items",[])
    complete=0; remaining=[]
    for item in items:
        files=item.get("files") or {}; required=[files.get("readme_path")]
        if files.get("main_path"): required.append(files.get("main_path"))
        if files.get("test_path"): required.append(files.get("test_path"))
        ok=all(exists(r) for r in required)
        if ok: complete+=1
        else: remaining.append({"id":item["id"],"missing":[r for r in required if not exists(r)]})
    pct=round((complete/len(items))*100,2) if items else 100
    report={"items_total":len(items),"items_complete":complete,"completion_percent":pct,"remaining":remaining}
    Path(args.output).parent.mkdir(parents=True, exist_ok=True); Path(args.output).write_text(json.dumps(report, indent=2)+"\n")
    print(json.dumps(report, indent=2))
if __name__ == "__main__": main()
