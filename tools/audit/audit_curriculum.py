#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(p): return json.load(open(p, encoding="utf-8"))
def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--curriculum-dir", default="curriculum"); p.add_argument("--output", default="dist/curriculum-audit.json"); args=p.parse_args()
    md=Path(args.metadata_dir)
    core=load(md/"path.core.json"); elec=load(md/"path.electives.json"); projects=load(md/"projects.json"); assessments=load(md/"assessments.json")
    report={"core_modules":len(core.get("modules",[])),"core_items":len(core.get("items",[])),"elective_items":len(elec.get("items",[])),"projects":len(projects.get("projects",[])),"assessments":len(assessments.get("assessments",[])),"findings":[]}
    Path(args.output).parent.mkdir(parents=True, exist_ok=True); Path(args.output).write_text(json.dumps(report, indent=2)+"\n")
    print(json.dumps(report, indent=2))
if __name__ == "__main__": main()
