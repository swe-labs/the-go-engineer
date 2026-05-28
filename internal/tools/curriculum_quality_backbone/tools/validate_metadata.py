import json, pathlib, re, sys
from collections import Counter, defaultdict
base = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path('/mnt/data/updated_curriculum/metadata')
errors=[]
warnings=[]
files=['path.core.json','path.electives.json','concepts.json','projects.json','assessments.json','crossrefs.json','failures.json','readme.contracts.json','migration.v2-to-v3.json','workspace.json','schema.v3.json']
data={}
for f in files:
    p=base/f
    if not p.exists():
        errors.append(f'missing file: {f}')
        continue
    try:
        data[f]=json.loads(p.read_text())
    except Exception as e:
        errors.append(f'invalid JSON {f}: {e}')
if errors:
    print('\n'.join(errors)); sys.exit(1)
core=data['path.core.json']; elect=data['path.electives.json']; concepts=data['concepts.json']; projects=data['projects.json']; ass=data['assessments.json']; cross=data['crossrefs.json']
items=core['items']+elect['items']; modules=core['modules']+elect['modules']
project_list=projects['projects']; assessments=ass['assessments']; concept_list=concepts['concepts']
item_ids=[i['id'] for i in items]; module_ids=[m['id'] for m in modules]; project_ids=[p['id'] for p in project_list]
concept_names=[c['concept'] for c in concept_list]
valid_ids=set(item_ids+module_ids+project_ids)
valid_concepts=set(concept_names)
# duplicate ids
for label, values in [('item',item_ids),('module',module_ids),('project',project_ids),('concept',concept_names)]:
    c=Counter(values)
    for v,n in c.items():
        if n>1: errors.append(f'duplicate {label}: {v}')
# core/elective architecture
if any(m['slug']=='advanced-electives' for m in core['modules']): errors.append('advanced-electives still present in core modules')
if not any(m['slug']=='advanced-electives' for m in elect['modules']): errors.append('advanced-electives missing from electives')
# statuses
for item in items:
    if item.get('status')!='stable': errors.append(f'{item["id"]}: status not stable')
    if item.get('readme_status')!='golden': errors.append(f'{item["id"]}: readme_status not golden')
    if item.get('zero_magic_status')!='golden': errors.append(f'{item["id"]}: zero_magic_status not golden')
for m in modules:
    if m.get('status')!='stable': errors.append(f'{m["id"]}: module status not stable')
    if m.get('readme_status')!='golden': errors.append(f'{m["id"]}: module readme_status not golden')
for p in project_list:
    if p.get('status')!='stable': errors.append(f'{p["id"]}: project status not stable')
    if p.get('readme_status')!='golden': errors.append(f'{p["id"]}: project readme_status not golden')
for a in assessments:
    if a.get('status')!='stable': errors.append(f'{a["id"]}: assessment status not stable')
    if a.get('readme_status')!='golden': errors.append(f'{a["id"]}: assessment readme_status not golden')
# phases key fixes
phase_expected={'module-09':'data','module-10':'security','module-12':'reliability','module-13':'performance','module-15':'delivery'}
for m in modules:
    if m['id'] in phase_expected and m.get('phase') != phase_expected[m['id']]: errors.append(f'{m["id"]}: phase {m.get("phase")} != {phase_expected[m["id"]]}')
for mid in ['module-10','module-12','module-14']:
    m=next((x for x in modules if x['id']==mid),None)
    if m and not m.get('contains_foundational_hard_concepts'):
        errors.append(f'{mid}: foundational hard concepts flag must be true')
# zero magic required fields
required=['problem_solved','why_it_exists','mental_model','under_the_hood','how_go_uses_it','beginner_mistakes','real_world_usage','proof_of_understanding','execution_timeline','failure_modes','hidden_magic_checks','performance_implications','step_by_step_execution']
for item in items:
    zm=item.get('zero_magic') or {}
    for k in required:
        if k not in zm or zm[k] in (None,'',[]): errors.append(f'{item["id"]}: zero_magic.{k} missing/empty')
# paths
for item in items:
    filesd=item.get('files') or {}
    pref='electives/' if item['id'].startswith('elective') else 'modules/'
    for k,v in filesd.items():
        if isinstance(v,str) and v and not v.startswith(pref): errors.append(f'{item["id"]}: files.{k} path {v} does not start with {pref}')
# refs
for item in items:
    for rel, refs in (item.get('crossrefs') or {}).items():
        for r in refs:
            tid=r.get('target_id')
            if tid and tid not in valid_ids and tid not in valid_concepts:
                errors.append(f'{item["id"]}: crossref target invalid: {tid}')
for r in cross.get('crossrefs',{}).get('references',[]):
    for k in ['from_id','target_id','to_id']:
        rid=r.get(k)
        if rid and rid not in valid_ids and rid not in valid_concepts:
            errors.append(f'global crossref {k} invalid: {rid}')
for p in project_list:
    for key in ['prerequisites','prerequisite_item_ids','reinforces','reinforcement_targets','required_concepts']:
        for t in p.get(key,[]) if isinstance(p.get(key), list) else []:
            if t not in valid_ids and t not in valid_concepts:
                errors.append(f'{p["id"]}: {key} target invalid: {t}')
    pa=p.get('placement_anchor_item_id')
    if pa and pa not in valid_ids: errors.append(f'{p["id"]}: placement anchor invalid: {pa}')
for a in assessments:
    for t in a.get('target_ids',[]):
        if t not in valid_ids and t not in valid_concepts:
            errors.append(f'{a["id"]}: target_id invalid: {t}')
# concepts coverage
for c in concept_list:
    if not c.get('canonical_owner'): errors.append(f'concept {c.get("concept")}: missing canonical_owner')
    elif c['canonical_owner'] not in valid_ids and c['canonical_owner'] not in valid_concepts: errors.append(f'concept {c["concept"]}: invalid canonical_owner {c["canonical_owner"]}')
    if not c.get('reinforcement_locations'): errors.append(f'concept {c["concept"]}: no reinforcement_locations')
# generic reasons
patterns=[r'Both .* build understanding within', r'This project reinforces .* through hands-on practice', r'requires understanding of \.$', r'Related concept:\s*builds on understanding from\s*\.$', r'This module builds on concepts from the previous module']
for item in items:
    for rel, refs in (item.get('crossrefs') or {}).items():
        for r in refs:
            reason=r.get('reason','')
            if any(re.search(pat,reason) for pat in patterns): errors.append(f'{item["id"]}: generic crossref reason: {reason}')
for r in cross.get('crossrefs',{}).get('references',[]):
    reason=r.get('reason','')
    if any(re.search(pat,reason) for pat in patterns): errors.append(f'global crossref generic reason: {reason}')
# placeholder flags exact only
placeholder_values=[]
def scan(o,path=''):
    if isinstance(o,dict):
        for k,v in o.items(): scan(v,path+'/'+str(k))
    elif isinstance(o,list):
        for i,v in enumerate(o): scan(v,path+'/'+str(i))
    elif isinstance(o,str):
        if o.strip().lower() in ('placeholder','scaffolded','todo','tbd'):
            placeholder_values.append((path,o))
for f,d in data.items(): scan(d,f)
for path,val in placeholder_values: errors.append(f'placeholder literal {val!r} at {path}')
# module sizes after split
mod_counts=Counter(i['module_id'] for i in items)
if mod_counts['module-15'] > 25: errors.append(f'module-15 still oversized: {mod_counts["module-15"]} items')
print('metadata files:', len(files))
print('core modules:', len(core['modules']))
print('elective modules:', len(elect['modules']))
print('core items:', len(core['items']))
print('elective items:', len(elect['items']))
print('projects:', len(project_list))
print('assessments:', len(assessments))
print('concepts:', len(concept_list))
print('module-15 items:', mod_counts['module-15'])
print('errors:', len(errors))
if errors:
    for e in errors[:200]: print('ERROR:',e)
    sys.exit(1)
print('VALIDATION PASSED')
