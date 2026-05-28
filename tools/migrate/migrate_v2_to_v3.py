#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--output", default="dist/v2-to-v3-migration-plan.json"); args=p.parse_args()
    migration=json.load(open(Path(args.metadata_dir)/"migration.v2-to-v3.json", encoding="utf-8"))
    report=json.load(open(Path(args.metadata_dir)/"legacy/unmapped-v2-report.json", encoding="utf-8"))
    out={"migration":migration.get("migration",{}),"coverage_summary":{k:report.get(k) for k in ["source_item_count","directly_mapped_count","policy_covered_count","unmapped_count"]}}
    Path(args.output).parent.mkdir(parents=True, exist_ok=True); Path(args.output).write_text(json.dumps(out, indent=2)+"\n")
    print(json.dumps(out["coverage_summary"], indent=2))
if __name__ == "__main__": main()
