#!/usr/bin/env python3
import argparse, json
from pathlib import Path

def load(path):
    return json.load(open(path, encoding="utf-8"))

def main():
    p = argparse.ArgumentParser()
    p.add_argument("--metadata-dir", default="metadata")
    p.add_argument("--only-missing", action="store_true")
    args = p.parse_args()
    items = load(Path(args.metadata_dir)/"path.core.json").get("items", []) + load(Path(args.metadata_dir)/"path.electives.json").get("items", [])
    count = 0
    for item in items:
        if item.get("type") != "lab":
            continue
        rel = (item.get("files") or {}).get("readme_path")
        if not rel: continue
        path = Path(rel)
        if args.only_missing and path.exists(): continue
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(f"# {item.get('title')}\n\n## Mission\n\n{item.get('learning_objective')}\n\n## Prerequisites\n\n{', '.join(item.get('prerequisites') or ['None'])}\n\n## Tasks\n\n1. Read the scenario.\n2. Implement the required behavior.\n3. Run validation.\n4. Record what failed and how you fixed it.\n\n## Verification\n\n```bash\n{(item.get('verification') or {}).get('test_command') or 'go test ./...'}\n```\n", encoding="utf-8")
        count += 1
    print(f"generated labs: {count}")
if __name__ == "__main__": main()
