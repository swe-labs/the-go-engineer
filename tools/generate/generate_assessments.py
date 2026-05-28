#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--only-missing", action="store_true"); args=p.parse_args()
    data=json.load(open(Path(args.metadata_dir)/"assessments.json", encoding="utf-8"))
    count=0
    for assessment in data.get("assessments", []):
        files=assessment.get("files") or {}
        readme=files.get("readme_path")
        if not readme: continue
        base=Path(readme).parent; base.mkdir(parents=True, exist_ok=True)
        if not args.only_missing or not Path(readme).exists():
            Path(readme).write_text(f"# {assessment.get('title')}\n\n## Purpose\n\nValidate mastery for the target curriculum items.\n\n## Targets\n\n{chr(10).join('- `'+t+'`' for t in assessment.get('target_ids', []))}\n\n## Evidence\n\nSubmit answers, code evidence, test output, and explanation where required.\n", encoding="utf-8")
        defaults={"questions_path":"# Questions\n\nAnswer each question with concrete reasoning and evidence.\n", "answer_key_path":"# Answer Key\n\nMaintainers should fill this with expected reasoning and grading notes.\n", "rubric_path":"# Rubric\n\nGrade using the assessment criteria in metadata.\n"}
        for key, content in defaults.items():
            rel=files.get(key)
            if rel and (not args.only_missing or not Path(rel).exists()):
                Path(rel).parent.mkdir(parents=True, exist_ok=True); Path(rel).write_text(content, encoding="utf-8")
        if files.get("assets_dir"): Path(files["assets_dir"]).mkdir(parents=True, exist_ok=True)
        count+=1
    print(f"generated assessments: {count}")
if __name__ == "__main__": main()
