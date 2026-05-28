#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(path): return json.load(open(path, encoding="utf-8"))

def main():
    p=argparse.ArgumentParser(); p.add_argument("--metadata-dir", default="metadata"); p.add_argument("--curriculum-dir", default="curriculum"); p.add_argument("--output", default="dist/curriculum.v3.full.generated.json"); args=p.parse_args()
    md=Path(args.metadata_dir)
    snapshot={"document_type":"generated_full_snapshot", "metadata":{}}
    for name in ["workspace.json","schema.v3.json","path.core.json","path.electives.json","concepts.json","projects.json","assessments.json","crossrefs.json","failures.json","readme.contracts.json","migration.v2-to-v3.json","VALIDATION.metadata.json"]:
        snapshot["metadata"][name]=load(md/name)
    out=Path(args.output); out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(snapshot, indent=2, ensure_ascii=False)+"\n", encoding="utf-8")
    print(f"wrote {out}")
if __name__ == "__main__": main()
