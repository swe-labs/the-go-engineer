#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--only-missing", action="store_true"); args=p.parse_args()
    data=json.load(open(Path(args.metadata_dir)/"projects.json", encoding="utf-8"))
    count=0
    for project in data.get("projects", []):
        rel=(project.get("files") or {}).get("readme_path")
        if not rel: continue
        path=Path(rel)
        if args.only_missing and path.exists(): continue
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(f"# {project.get('title')}\n\n## Mission\n\nBuild a portfolio-quality project that proves the target module outcomes.\n\n## Deliverables\n\n- Working implementation\n- Tests or verification evidence\n- README with run instructions\n- Notes on trade-offs and failure modes\n\n## Verification\n\n```bash\n{(project.get('verification') or {}).get('test_command') or 'go test ./...'}\n```\n\n## Rubric\n\nUse the project rubric in metadata as the grading source of truth.\n", encoding="utf-8")
        for key in ["starter_path", "solution_path", "assets_dir"]:
            if (project.get("files") or {}).get(key): Path(project["files"][key]).mkdir(parents=True, exist_ok=True)
        count+=1
    print(f"generated project READMEs: {count}")
if __name__ == "__main__": main()
