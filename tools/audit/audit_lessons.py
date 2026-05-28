#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(p): return json.load(open(p, encoding="utf-8"))

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--curriculum-dir", default="curriculum"); p.add_argument("--output", default="dist/lessons-audit.json"); args=p.parse_args()
    md=Path(args.metadata_dir); findings=[]
    core=load(md/"path.core.json"); elec=load(md/"path.electives.json")
    if "lessons" == "modules":
        records=core.get("modules",[])+elec.get("modules",[])
    elif "lessons" == "lessons":
        records=[i for i in core.get("items",[])+elec.get("items",[]) if i.get("type")=="lesson"]
    elif "lessons" == "labs":
        records=[i for i in core.get("items",[])+elec.get("items",[]) if i.get("type")=="lab"]
    elif "lessons" == "projects":
        records=load(md/"projects.json").get("projects",[])
    else:
        records=load(md/"assessments.json").get("assessments",[])
    report={"kind":"lessons","count":len(records),"findings":findings}
    Path(args.output).parent.mkdir(parents=True, exist_ok=True); Path(args.output).write_text(json.dumps(report, indent=2)+"\n")
    print(json.dumps(report, indent=2))
if __name__ == "__main__": main()
