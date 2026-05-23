"""Investigate workspace.json structure and fix H-7 (4 missing gates)."""
import json

with open('curriculum/workspace.json', 'r', encoding='utf-8') as f:
    ws = json.load(f)

print("workspace.json top-level keys:", list(ws.keys()))
print()

# Check quality_gates structure
qg = ws.get('quality_gates')
if qg:
    print(f"quality_gates type: {type(qg).__name__}")
    if isinstance(qg, list):
        print(f"quality_gates has {len(qg)} entries")
        print(f"First entry keys: {list(qg[0].keys()) if qg else 'empty'}")
        for g in qg[:5]:
            print(f"  gate: {g}")
    elif isinstance(qg, dict):
        print(f"quality_gates has {len(qg)} entries")
        for k, v in list(qg.items())[:5]:
            print(f"  {k}: {v}")
    else:
        print(f"quality_gates value: {str(qg)[:200]}")
