#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--output", default="dist/migration-report.json"); args=p.parse_args()
    md=Path(args.metadata_dir); migration=json.load(open(md/"migration.v2-to-v3.json", encoding="utf-8")); unmapped=json.load(open(md/"legacy/unmapped-v2-report.json", encoding="utf-8"))
    report={"status":"passed" if unmapped.get("unmapped_count")==0 else "failed", "migration_summary":migration.get("migration",{}).get("legacy_item_coverage_summary",{}), "unmapped_count":unmapped.get("unmapped_count"), "source_sha256":unmapped.get("source_sha256")}
    Path(args.output).parent.mkdir(parents=True, exist_ok=True); Path(args.output).write_text(json.dumps(report, indent=2)+"\n")
    print(json.dumps(report, indent=2))
if __name__ == "__main__": main()
