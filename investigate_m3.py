"""Investigate module-10 and core-11-28 for Milestone 3."""
import json

with open('curriculum/path.core.json', 'r', encoding='utf-8') as f:
    core = json.load(f)

items = core['items']
item_map = {it['id']: it for it in items}

# Module-10 items
print("Module-10 items:")
for it in items:
    mid = it.get('module_id', '')
    if mid == 'module-10':
        nid = it.get('next_item_ids', [])
        o = it.get('order', '?')
        print(f"  {it['id']:20s} order={str(o)} title={it['title']} next={nid}")

# Module-10 definition
for m in core['modules']:
    if m['id'] == 'module-10':
        print(f"\nmodule-10 entry: {m.get('entry_item_ids')}")
        print(f"module-10 terminal: {m.get('terminal_item_ids')}")

# core-11-28
print("\ncore-11-28:")
it = item_map.get('core-11-28', {})
print(f"  title: {it.get('title')}")
print(f"  status: {it.get('status')}")
print(f"  subtype: {it.get('subtype')}")
print(f"  type: {it.get('type')}")
print(f"  module_id: {it.get('module_id')}")
print(f"  next_item_ids: {it.get('next_item_ids')}")
print(f"  prerequisites: {it.get('prerequisites')}")
crossrefs = it.get('crossrefs', {})
if crossrefs:
    po = crossrefs.get('preview_only', [])
    if po:
        print(f"  preview_only targets: {po}")
    bo = crossrefs.get('builds_on', [])
    if bo:
        print(f"  builds_on targets: {[e.get('target_id','') for e in bo]}")

# Module-11 items to see where core-11-28 sits
print("\nModule-11 items (25-29):")
for it in items:
    iid = it.get('id', '')
    if iid in ['core-11-27', 'core-11-28', 'core-11-29']:
        o = it.get('order', '?')
        print(f"  {iid:20s} order={str(o)} title={it['title']} next={it.get('next_item_ids',[])}")

# What does opslane-08 need for idempotency?
print("\nopslane-08:")
ops8 = item_map.get('opslane-08', {})
if ops8:
    print(f"  title: {ops8.get('title')}")
    print(f"  module_id: {ops8.get('module_id')}")
