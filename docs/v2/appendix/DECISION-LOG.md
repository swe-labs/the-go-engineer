# V2 Decision Log

Use this file to record planning decisions that affect v2 structure, sequencing, migration, or
release strategy.

## How To Use This Log

Each decision should capture:

- date
- decision summary
- why the decision was made
- alternatives considered
- follow-up actions

## Working Decisions

### 2026-04-12: learner-facing lessons should use README-first explanation with cleaner code files

- Status: Working decision
- Why it matters: the zero-to-engineer learner promise requires more explanation than source-file
  headers and scattered inline comments can carry cleanly
- Alternatives:
  - keep code-first lessons as the default everywhere and rely on inline comments
    - make learner-facing READMEs the primary explanation surface while keeping runnable code
      smaller, clearer, and still mandatory
- Follow-up: revise the lesson spec and lesson template, add a canonical lesson README template, and
  retrofit rebuilt sections to the new contract

### 2026-04-09: beta is a learner-facing architecture redesign, not just alpha plus polish

- Status: Working decision
- Why it matters: alpha validated migration discipline, but beta must reshape the learner-facing
  system so the repo actually functions as a zero-to-engineer training system
- Alternatives:
  - keep the alpha section architecture and mostly polish it
  - use beta to regroup alpha content into a stage-based learning system with new beginner,
    production, expert, and flagship layers
- Follow-up: use `docs/v2/17-BETA-ARCHITECTURE-DECISION.md` as the baseline for beta planning and
  implementation issue design

### 2026-04-09: beta stage mapping may split alpha sections wherever learner-facing boundaries are stronger than current folder boundaries

- Status: Working decision
- Why it matters: the beta stage model cannot stay honest if alpha sections are treated as
  indivisible when they already mix beginner setup, architecture, operations, and capstone work
- Alternatives:
  - regroup only whole alpha sections and accept fuzzy stage boundaries
  - allow Section `01` and Section `14` style splits wherever the learner-facing architecture is
    materially clearer
- Follow-up: use `docs/v2/18-BETA-STAGE-MAPPING.md` as the concrete authority for `#173` and later
  beta implementation issue design

### 2026-04-09: beta exercise design separates type, difficulty, starter mode, and verification mode

- Status: Working decision
- Why it matters: exercise systems drift when one label is forced to carry format, challenge,
  scaffolding, and proof all at once
- Alternatives:
  - keep a looser alpha-style exercise model and let sections improvise
  - use separate axes so beta can scale a real exercise bank without muddy semantics
- Follow-up: use `docs/v2/19-BETA-EXERCISE-RUBRIC-SYSTEM.md` as the authority for `#174` and later
  beta practice implementation work

### 2026-04-09: beta must treat foundation, production, expert pressure, and flagship work as first-class curriculum layers

- Status: Working decision
- Why it matters: beta cannot honestly claim zero-to-engineer scope if those layers remain implied
  or optional
- Alternatives:
  - keep those concerns distributed loosely across existing stages
  - define them explicitly so implementation work has owned outcomes and clearer rollout order
- Follow-up: use `docs/v2/20-BETA-LAYER-DEFINITIONS.md` as the authority for `#175` and later beta
  implementation design

### 2026-04-07: v2 planning is ready to open prototype work

- Status: Working decision
- Why it matters: the team needs a clear gate between endless planning and prototype execution
- Alternatives:
  - delay prototype work until every open planning question is fully resolved
  - begin prototype work once the planning stack is coherent enough to guide it
- Follow-up: start the prototype phase under `#84` and use the prototype to resolve non-blocking
  design questions

### 2026-04-07: keep the current top-level section numbering and directory layout through prototype and early migration waves

- Status: Working decision
- Why it matters: early renumbering would add churn before the training system itself is proven
- Alternatives:
  - preserve the current numbering and layout first
  - redesign the top-level layout before prototype validation
- Follow-up: only revisit numbering after the prototype and first migration waves prove a clear
  learner benefit

### 2026-04-07: keep `curriculum.json` active while designing richer v2 metadata in parallel

- Status: Working decision
- Why it matters: the repo needs schema evolution without breaking the current stable line
- Alternatives:
  - replace `curriculum.json` immediately
  - keep it active and prototype a richer v2 schema in parallel
- Follow-up: prove the richer metadata model in the prototype before deciding final canonical
  storage

### 2026-04-07: keep assessment and progress repo-native through the prototype and alpha stages

- Status: Working decision
- Why it matters: platform work would dilute focus before the curriculum system is proven
- Alternatives:
  - start platform features during early v2 planning
  - keep assessment repo-native until later rollout stages
- Follow-up: revisit only after alpha if repo-native checkpoints and projects prove insufficient

### 2026-04-07: v2 should use one curriculum with three canonical learning paths

- Status: Working decision
- Why it matters: path design affects docs, metadata, checkpoints, and migration behavior
- Alternatives:
  - maintain one curriculum with route logic
  - build separate content trees for different learner types
- Follow-up: prove the path model in the prototype through one section README, one checkpoint, and
  one milestone item

### 2026-04-07: the Bridge Path requires milestone mini-projects rather than every mini-project

- Status: Working decision
- Why it matters: the Bridge Path should preserve proof without recreating the Full Path workload
- Alternatives:
  - require every mini-project
  - require only milestone mini-projects plus path-critical project work
- Follow-up: verify during the prototype that this still preserves honest readiness signals

### 2026-04-07: every v2 section should have a top-level README as a required navigation surface

- Status: Working decision
- Why it matters: the current repo lacks a consistent section-entry contract for learners
- Alternatives:
  - rely mainly on the root README and directory order
  - require a top-level section README with path-aware entry guidance
- Follow-up: test the section README contract in the prototype slice before applying it widely

### 2026-04-07: keep the learning-path guide at the repo root and keep `docs/curriculum/README.md` as the public curriculum map during the prototype stage

- Status: Working decision
- Why it matters: learners need stable discovery surfaces while v2 planning is still being proven
- Alternatives:
  - move those docs immediately into a new layout
  - keep the existing public discovery locations during the prototype stage
- Follow-up: reconsider doc placement after the first learner-facing migration wave if the
  navigation system proves a better structure

### 2026-04-07: local sub-track README files are only required when a section has multiple learner-facing internal tracks or milestone surfaces

- Status: Working decision
- Why it matters: the docs system needs consistency without forcing README files where they add no
  navigation value
- Alternatives:
  - require sub-track README files everywhere
  - require them only where internal track navigation is real
- Follow-up: test this threshold in the prototype section outline

### 2026-04-07: the validator should prioritize structural and metadata truth before broader teaching checks

- Status: Working decision
- Why it matters: the repo needs trustworthy automation without pretending every curriculum quality
  concern can be parsed mechanically
- Alternatives:
  - try to automate most review logic immediately
  - keep the validator focused on paths, commands, ids, links, and required assets first
- Follow-up: expand validator scope gradually after the prototype proves which additional checks are
  worth automating

### 2026-04-07: navigation integrity should be split between validator truth and reviewer judgment

- Status: Working decision
- Why it matters: link existence and learner usefulness are different quality problems
- Alternatives:
  - try to automate all navigation quality
  - let the validator check truth while reviewers check clarity and usefulness
- Follow-up: keep refining this split during prototype validation

### 2026-04-07: rubric-verified items should use a standard README shape

- Status: Working decision
- Why it matters: rubric-based items need visible and consistent completion rules
- Alternatives:
  - let rubric items vary freely
  - require a standard README shape for rubric-based verification
- Follow-up: prove the README shape in the prototype checkpoint and project surfaces

### 2026-04-07: keep the current validator as a single growing tool through the prototype stage

- Status: Working decision
- Why it matters: splitting tooling early would add maintenance overhead before the right boundaries
  are proven
- Alternatives:
  - split validation into multiple tools now
  - keep one tool through the prototype stage and split later only if justified
- Follow-up: revisit after the first alpha wave if validator scope becomes unwieldy

### 2026-04-07: Section 04 Functions and Errors is the canonical first prototype section

- Status: Working decision
- Why it matters: the prototype needs one coherent section that can prove lesson, exercise,
  checkpoint, and mini-project behavior without infrastructure-heavy setup
- Alternatives:
  - use an earlier foundations section with less project pressure
  - use a design or application section with more complexity
  - use Section 04 as the hinge between foundations and the first mini-project
- Follow-up: implement `#85` through `#88` around the Section 04 prototype slice on `planning/v2`

### 2026-04-07: the prototype section outline uses five lessons, one drill, one guided exercise, one checkpoint, and one mini-project linkage

- Status: Working decision
- Why it matters: the prototype must be large enough to prove the system, but small enough to avoid
  becoming a stealth section migration
- Alternatives:
  - prototype too little and fail to prove the section contract
  - prototype a near-full migration of the section
  - use one compact but complete section slice with explicit practice layers
- Follow-up: use this outline as the container for `#86`, `#87`, and `#88`

### 2026-04-07: `FEP.4 Custom Errors, Wrapping, and Inspection` is the canonical prototype lesson

- Status: Working decision
- Why it matters: the prototype needs one lesson that best proves the v2 lesson contract in
  structure, tone, metadata, and production relevance
- Alternatives:
  - use `FEP.3` as a simpler concept lesson
  - use `FEP.4` as a richer pattern lesson with one bounded layered example
- Follow-up: use `FEP.4` as the reference lesson for `#86` and shape `#87` around its exit ramp

### 2026-04-07: the prototype guided exercise should use `_starter/`, but the prototype checkpoint should not

- Status: Working decision
- Why it matters: the prototype must prove the difference between supported practice and readiness
  validation
- Alternatives:
  - include scaffolding in both items and blur the boundary
  - remove scaffolding from both items and make the first synthesis step harsher than intended
  - use starter scaffolding for `FEP.E1` and require `FEP.C1` to stand on explicit pass criteria
    instead
- Follow-up: use this split in `#87` and carry it forward into later exercise and checkpoint
  templates unless a section has unusual setup burden

### 2026-04-07: prototype local path numbering should follow canonical section order rather than inherited v1 item numbering

- Status: Working decision
- Why it matters: prototype paths should describe the learner-facing v2 flow clearly instead of
  leaking numbering from the v1 source material
- Alternatives:
  - keep inherited v1-inspired numbers and accept ordering drift
  - renumber the planned prototype paths so the exercise, integration lesson, checkpoint, and
    mini-project follow the actual section order
- Follow-up: use ordered prototype paths for `FEP.E1`, `FEP.5`, `FEP.C1`, and `FEP.P1`

### 2026-04-07: the prototype mini-project should remain inside the section and may point `next_items` at the next section id

- Status: Working decision
- Why it matters: the first mini-project must prove section-level milestone behavior without
  forcing a separate projects tree or an extra schema layer
- Alternatives:
  - move the first mini-project into a separate top-level project surface
  - keep the mini-project inside Section 04 and allow section-exit navigation through the existing
    `next_items` field
- Follow-up: use this rule in `#88` and the first metadata example
