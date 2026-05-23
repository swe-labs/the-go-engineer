# Fix findings from audit pass 4
# Target: crossref preview_only, assessment_id, workspace gates, display_text
# Run with: python fix_findings.py

import json
import re

# ── Load everything ──
with open('curriculum/path.core.json', 'r', encoding='utf-8') as f:
    core = json.load(f)

with open('curriculum/crossrefs.json', 'r', encoding='utf-8') as f:
    xrefs = json.load(f)

with open('curriculum/workspace.json', 'r', encoding='utf-8') as f:
    workspace = json.load(f)

item_map = {it['id']: it for it in core['items']}

# ════════════════════════════════════════════════════════════════
# H-8: Check proof.assessment_id mismatches
# ════════════════════════════════════════════════════════════════
print("=== H-8: proof.assessment_id mismatches ===")
mismatches = []
for item in core['items']:
    iid = item.get('id', '')
    mid = item.get('module_id', '')
    proof = item.get('proof', {})
    aid = proof.get('assessment_id', '')
    if aid:
        expected = f"assessment-{mid}"
        if aid != expected:
            mismatches.append((iid, mid, aid, expected))

if mismatches:
    print(f"Found {len(mismatches)} mismatches:")
    for iid, mid, aid, expected in mismatches:
        print(f"  {iid:20s} module={mid:15s} has assessment_id={aid:30s} expected={expected}")
else:
    print("All assessment_ids match their module!")

# ════════════════════════════════════════════════════════════════
# C-4: Check which of the 11 preview_only refs already exist
# ════════════════════════════════════════════════════════════════
print("\n=== C-4: Missing preview_only crossrefs ===")
audit_missing = [
    ('core-06-19', 'core-11-22', 'race detector'),
    ('core-06-19', 'core-11-21', 'race condition'),
    ('core-07-20', 'core-11-06', 'context cancellation'),
    ('core-07-20', 'core-08-14', 'HTTP timeouts'),
    ('core-07-21', 'core-08-15', 'graceful shutdown'),
    ('core-08-03', 'core-11-11', 'goroutines'),
    ('core-08-14', 'core-11-06', 'context cancellation'),
    ('core-08-14', 'core-11-08', 'context deadlines'),
    ('core-08-15', 'core-11-11', 'goroutines'),
    ('core-09-23', 'core-11-06', 'context cancellation'),
    ('core-11-07', 'core-08-14', 'HTTP timeouts — forward ref'),
]

# Build set of existing preview_only refs
existing_previews = set()
refs_list = xrefs['crossrefs']['references']
for ref in refs_list:
    if ref.get('relation') == 'preview_only':
        existing_previews.add((ref['from_id'], ref['to_id']))

still_missing = []
for fid, tid, desc in audit_missing:
    if (fid, tid) not in existing_previews:
        still_missing.append((fid, tid, desc))
    else:
        print(f"  Already exists: {fid} -> {tid} ({desc})")

print(f"Still missing: {len(still_missing)}")
for fid, tid, desc in still_missing:
    print(f"  {fid} -> {tid} ({desc})")

# ════════════════════════════════════════════════════════════════
# H-4: Count related crossrefs with empty display_text
# ════════════════════════════════════════════════════════════════
print("\n=== H-4: Related crossrefs with empty display_text ===")
empty_display = []
for ref in refs_list:
    if ref.get('relation') == 'related' and ref.get('display_text', '') == '':
        empty_display.append(ref)
print(f"Related refs with empty display_text: {len(empty_display)}")

# ════════════════════════════════════════════════════════════════
# H-7: Check workspace.json vs migration.json gates
# ════════════════════════════════════════════════════════════════
print("\n=== H-7: Quality gate divergence ===")
workspace_gates = set(g['id'] for g in workspace.get('quality_gates', []))
print(f"Workspace gates ({len(workspace_gates)}): {sorted(workspace_gates)}")

# migration gates (known from audit)
migration_only_gates = [
    'assessment-criteria-ids-match-module',
    'crossref-entity-domains-valid',
    'project-module-id-matches-reinforces-prefix',
    'zero-magic-semantic-quality',
]
for g in migration_only_gates:
    if g in workspace_gates:
        print(f"  {g}: already in workspace")
    else:
        print(f"  {g}: MISSING from workspace")
