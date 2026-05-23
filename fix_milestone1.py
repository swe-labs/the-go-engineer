"""Milestone 1: Fix preview_only crossrefs, assessment_id mismatches, display_text."""
import json, copy

# ── Load ──
with open('curriculum/path.core.json', 'r', encoding='utf-8') as f:
    core = json.load(f)

with open('curriculum/crossrefs.json', 'r', encoding='utf-8') as f:
    xrefs = json.load(f)

items = core['items']
item_map = {it['id']: it for it in items}
refs = xrefs['crossrefs']['references']

# ════════════════════════════════════════════════════════════════
# 1. Fix proof.assessment_id mismatches (63 items)
# ════════════════════════════════════════════════════════════════
print("=== Fixing assessment_id mismatches ===")
fixes = 0
# Map known wrong assessment_ids to their correct module mapping
# Module-08 items have assessment-module-11 -> should be assessment-module-08
# Module-10 items have assessment-module-09 -> should be assessment-module-10
# Module-15 items have assessment-module-16 -> should be assessment-module-15
assessment_fix_map = {
    'assessment-module-11': 'assessment-module-08',  # for module-08 items
}

for item in items:
    iid = item.get('id', '')
    mid = item.get('module_id', '')
    proof = item.get('proof', {})
    if not proof:
        continue
    aid = proof.get('assessment_id', '')
    expected = f"assessment-{mid}"
    
    if aid and aid != expected:
        proof['assessment_id'] = expected
        fixes += 1

print(f"Fixed {fixes} assessment_id mismatches")

# ════════════════════════════════════════════════════════════════
# 2. Add 11 missing preview_only crossrefs
# ════════════════════════════════════════════════════════════════
print("\n=== Adding missing preview_only crossrefs ===")

# Build existing set
existing = {(r['from_id'], r['to_id']) for r in refs}

new_refs_data = [
    ('core-06-19', 'core-11-22', 'Race detector', 'Learners using `go test -race` in module-06 preview a concept explained fully in Module 11.21-11.22.'),
    ('core-06-19', 'core-11-21', 'Race condition', 'Learners using `go test -race` in module-06 preview a concept explained fully in Module 11.21-11.22.'),
    ('core-07-20', 'core-11-06', 'Context cancellation', 'Context cancellation in I/O operations previews the full context lifecycle explained in Module 11.06.'),
    ('core-07-20', 'core-08-14', 'HTTP timeouts', 'Server timeouts enforced in module-08 are previewed when learners configure timeout-sensitive I/O in module-07.'),
    ('core-07-21', 'core-08-15', 'Graceful shutdown', 'Clean resource teardown in module-07 previews the full graceful shutdown pattern in module-08.'),
    ('core-08-03', 'core-11-11', 'Goroutines', 'HTTP handlers in module-08 run in goroutines. Learners preview goroutines before the full explanation in module 11.11.'),
    ('core-08-14', 'core-11-06', 'Context cancellation', 'Server timeouts in module-08 use context deadlines internally, previewed before the full context lifecycle in module 11.06-11.08.'),
    ('core-08-14', 'core-11-08', 'Context deadlines', 'Server timeouts in module-08 use context deadlines internally, previewed before the full context lifecycle in module 11.06-11.08.'),
    ('core-08-15', 'core-11-11', 'Goroutines', 'Graceful shutdown in module-08 drains goroutine-based handlers, previewed before the full goroutine explanation in module 11.11.'),
    ('core-09-23', 'core-11-06', 'Context cancellation', 'Database query cancellation in module-09 uses context, previewed before the full context lifecycle in module 11.06.'),
    ('core-11-07', 'core-08-14', 'HTTP timeouts', 'Forward reference: context deadlines in module-11 are directly applied in HTTP server timeout configuration in module-08.'),
]

added = 0
for fid, tid, label, reason in new_refs_data:
    if (fid, tid) in existing:
        print(f"  Already exists: {fid} -> {tid}")
        continue
    
    # Find labels from item titles
    from_title = item_map.get(fid, {}).get('title', fid)
    to_title = item_map.get(tid, {}).get('title', tid)
    from_mod = item_map.get(fid, {}).get('module_id', '').replace('module-', 'Module ').replace('-', '.')
    to_mod = item_map.get(tid, {}).get('module_id', '').replace('module-', 'Module ').replace('-', '.')
    
    new_ref = {
        'from_id': fid,
        'relation': 'preview_only',
        'to_id': tid,
        'target_id': tid,
        'label': f'Preview only: Explained fully in {to_mod} \u2014 {to_title}',
        'display_text': f'Preview: {to_title} \u2014 covered in {to_mod}',
        'reason': reason,
        'required': False,
    }
    refs.append(new_ref)
    added += 1
    print(f"  Added: {fid} -> {tid} ({label})")

print(f"Added {added} new preview_only refs")

# ════════════════════════════════════════════════════════════════
# 3. Populate display_text for related crossrefs
# ════════════════════════════════════════════════════════════════
print("\n=== Populating display_text for related crossrefs ===")
filled = 0
for ref in refs:
    if ref.get('relation') == 'related' and not ref.get('display_text'):
        # Use the label as display_text, stripping "Related: " prefix if present
        label = ref.get('label', '')
        if label:
            if label.startswith('Related: '):
                ref['display_text'] = label
            else:
                ref['display_text'] = f'Related: {label}'
        else:
            ref['display_text'] = 'Related concept'
        filled += 1

print(f"Filled display_text for {filled} related crossrefs")

# ════════════════════════════════════════════════════════════════
# Write
# ════════════════════════════════════════════════════════════════
with open('curriculum/path.core.json', 'w', encoding='utf-8') as f:
    json.dump(core, f, indent=2, ensure_ascii=False)
print("\npath.core.json written")

with open('curriculum/crossrefs.json', 'w', encoding='utf-8') as f:
    json.dump(xrefs, f, indent=2, ensure_ascii=False)
print("crossrefs.json written")

# ── Verify ──
print("\n=== Verification ===")
with open('curriculum/path.core.json', 'r', encoding='utf-8') as f:
    v = json.load(f)
mismatches = 0
for item in v['items']:
    mid = item.get('module_id', '')
    aid = item.get('proof', {}).get('assessment_id', '')
    if aid and aid != f'assessment-{mid}':
        mismatches += 1
print(f"Remaining assessment_id mismatches: {mismatches}")

with open('curriculum/crossrefs.json', 'r', encoding='utf-8') as f:
    vx = json.load(f)
px_count = sum(1 for r in vx['crossrefs']['references'] if r.get('relation') == 'preview_only')
empty_dt = sum(1 for r in vx['crossrefs']['references'] if r.get('relation') == 'related' and not r.get('display_text'))
print(f"Preview_only refs total: {px_count}")
print(f"Related refs with empty display_text: {empty_dt}")
