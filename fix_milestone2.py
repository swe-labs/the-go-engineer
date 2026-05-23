"""Milestone 2: Add 4 missing gates to workspace.json, check zero_magic_status alignment."""
import json

with open('curriculum/workspace.json', 'r', encoding='utf-8') as f:
    ws = json.load(f)

# Current gates
gates = ws.get('quality_gates', [])
print(f"Current gates ({len(gates)}): {gates}")

# Gates to add (from migration.json audit finding H-7)
new_gates = [
    'assessment-criteria-ids-match-module',
    'crossref-entity-domains-valid',
    'project-module-id-matches-reinforces-prefix',
    'zero-magic-semantic-quality',
]

added = 0
for g in new_gates:
    if g not in gates:
        gates.append(g)
        added += 1
        print(f"  Added: {g}")

print(f"Added {added} new gates. Total: {len(gates)}")

# Write
ws['quality_gates'] = gates
with open('curriculum/workspace.json', 'w', encoding='utf-8') as f:
    json.dump(ws, f, indent=2, ensure_ascii=False)
print("workspace.json written")
