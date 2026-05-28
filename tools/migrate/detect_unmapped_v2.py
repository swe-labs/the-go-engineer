#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(p): return json.load(open(p, encoding="utf-8"))
def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); args=p.parse_args()
    report=load(Path(args.metadata_dir)/"legacy/unmapped-v2-report.json")
    print(json.dumps({"source_item_count":report.get("source_item_count"),"unmapped_count":report.get("unmapped_count"),"unmapped_items":report.get("unmapped_items",[])}, indent=2))
    raise SystemExit(1 if report.get("unmapped_count") else 0)
if __name__ == "__main__": main()
