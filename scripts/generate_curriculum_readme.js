const fs = require('fs');

const core = JSON.parse(fs.readFileSync('curriculum/path.core.json', 'utf8'));
const electives = JSON.parse(fs.readFileSync('curriculum/path.electives.json', 'utf8'));
const projects = JSON.parse(fs.readFileSync('curriculum/projects.json', 'utf8'));
const assessments = JSON.parse(fs.readFileSync('curriculum/assessments.json', 'utf8'));
const crossrefs = JSON.parse(fs.readFileSync('curriculum/crossrefs.json', 'utf8'));
const concepts = JSON.parse(fs.readFileSync('curriculum/concepts.json', 'utf8'));
const failures = JSON.parse(fs.readFileSync('curriculum/failures.json', 'utf8'));

function generateFullReadme() {
  const parts = [];

  parts.push('# The Go Engineer v3 — Curriculum\n');
  parts.push('> Schema: ' + core.schema_version + ' | Curriculum: ' + core.curriculum_version + ' | Status: ' + core.status + ' | Updated: ' + core.last_updated + '\n');

  // Architecture Decision
  parts.push('## Architecture Decision\n');
  parts.push('**Decision:** Reconcile and establish the Report 1 phase-based zero-to-software-engineer path as the canonical v3 curriculum.\n');
  parts.push('- **v2 Policy:** `curriculum.v2.json` is frozen legacy import data, preserved for historical reference and backward compatibility.');
  parts.push('- **Authoritative Source of Truth:** Authoritative curriculum metadata resides in the modular JSON files: `curriculum/path.core.json`, `curriculum/path.electives.json`, `curriculum/projects.json`, `curriculum/assessments.json`, `curriculum/crossrefs.json`, `curriculum/concepts.json`, and `curriculum/failures.json`.');
  parts.push('- **Physical Curriculum:** Physical lesson READMEs, runnable examples, test suites, and exercises are organized canonically under `curriculum/modules/`.');
  parts.push('- **Core Scope:** The core curriculum is dedicated to production Go backend engineering (systems, networking, persistence, concurrency, observability, architecture, and deployment). Modern full-stack frontend engineering is defined as a companion path (semantic HTML, CSS, JavaScript fundamentals, progressive enhancement with Go templates/templ/HTMX/Tailwind).');
  parts.push('- **Capstone & Learning Path Flow:** Flagship Opslane (Module 18) directly requires the complete backend engineering and deployment runbook foundation of Module 15 (`core-15-19`). Module 16 (Portfolio and Interview Readiness) is intentionally decoupled from technical implementation so learners can construct their flagship capstone first, producing the production artifact that Module 16 packages for interviews. Module 17 (Advanced Electives) is an independent elective track and not on the mandatory capstone path.\n');

  // Quick stats
  const coreItems = core.items;
  const lessons = coreItems.filter(it => it.type === 'lesson').length;
  const labs = coreItems.filter(it => it.type === 'lab').length;
  const checkpoints = coreItems.filter(it => it.type === 'checkpoint').length;
  const capstones = coreItems.filter(it => it.type === 'capstone').length;
  const pathAssessments = coreItems.filter(it => it.type === 'assessment').length;
  const pathProjects = coreItems.filter(it => it.type === 'project').length;
  const electiveItems = electives.items.length;
  const prjRecords = projects.projects.length;
  const regAssessments = assessments.assessments.length;
  const conceptCount = concepts.concepts.length;
  const failureCategoryCount = failures.failure_categories ? failures.failure_categories.length : 13;

  parts.push('## Quick Stats\n');
  parts.push('| Metric                           | Count | Notes                                                             |');
  parts.push('| -------------------------------- | ----- | ----------------------------------------------------------------- |');
  parts.push('| Core Modules                     | 18    | Modules 00 through 16, plus Module 18 (Flagship Opslane)          |');
  parts.push('| Elective Modules                 | 1     | Module 17 (Advanced Electives)                                    |');
  parts.push('| Total Modules                    | 19    | Modules 00–18                                                     |');
  parts.push('| Core Path Items (Metadata Total) | ' + coreItems.length + '   | Total items array in `path.core.json` across all core modules     |');
  parts.push('| Instructional Lessons & Labs     | ' + (lessons + labs) + '   | ' + lessons + ' lessons + ' + labs + ' lab across Modules 00–16                         |');
  parts.push('| Flagship Opslane Milestones      | ' + (checkpoints + capstones) + '    | ' + checkpoints + ' checkpoints + ' + capstones + ' capstone defense in Module 18          |');
  parts.push('| Learner Instructional Units      | ' + (lessons + labs + checkpoints + capstones) + '   | ' + (lessons + labs) + ' instructional lessons/labs + ' + (checkpoints + capstones) + ' Opslane milestones            |');
  parts.push('| Elective Items                   | ' + electiveItems + '    | Explicitly marked as planned / unreleased draft in Module 17      |');
  parts.push('| Projects                         | ' + prjRecords + '    | 27 module projects + 1 flagship capstone defense project          |');
  parts.push('| In-Path Assessment Checkpoints   | ' + pathAssessments + '    | 16 module checkpoints + 14 Opslane checkpoints + 27 project-bound nodes |');
  parts.push('| Assessment Registry Records      | ' + regAssessments + '    | 19 module assessments + 28 project assessments in `assessments.json`|');
  parts.push('| Canonical Concepts               | ' + conceptCount + '    | Fully resolved ownership without leaks in `concepts.json`         |');
  parts.push('| Operational Categories           | ' + failureCategoryCount + '    | Mapped across 18 module coverage entries                          |\n');

  // Module Map
  parts.push('## Module Map\n');
  parts.push('| #   | Module | Phase | Status | Lessons/Labs | Milestones | Projects | Assessments | Total Path Items | Path |');
  parts.push('| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | --- |');

  const phaseMap = {
    'module-00': 'orientation',
    'module-01': 'foundations',
    'module-02': 'tooling',
    'module-03': 'go-core',
    'module-04': 'go-core',
    'module-05': 'engineering-core',
    'module-06': 'engineering-core',
    'module-07': 'cli-io',
    'module-08': 'backend',
    'module-09': 'data',
    'module-10': 'security',
    'module-11': 'concurrency',
    'module-12': 'reliability',
    'module-13': 'performance',
    'module-14': 'architecture',
    'module-15': 'delivery',
    'module-16': 'career',
    'module-18': 'flagship'
  };

  let totLessons = 0, totMilestones = 0, totProjects = 0, totAssessments = 0, totItems = 0;

  core.modules.forEach(m => {
    const mNum = parseInt(m.id.replace('module-', ''), 10);
    const mItems = coreItems.filter(it => it.module_id === m.id);
    const mLessons = mItems.filter(it => it.type === 'lesson' || it.type === 'lab').length;
    const mMilestones = mItems.filter(it => it.type === 'checkpoint' || it.type === 'capstone').length;
    const mProjects = mItems.filter(it => it.type === 'project').length;
    const mAssessments = mItems.filter(it => it.type === 'assessment').length;
    const mTotal = mItems.length;

    totLessons += mLessons;
    totMilestones += mMilestones;
    totProjects += mProjects;
    totAssessments += mAssessments;
    totItems += mTotal;

    const phase = phaseMap[m.id] || m.phase;
    const pathCol = mNum === 18 ? '`curriculum/modules/18-flagship-opslane`' : '`curriculum/' + (m.path.startsWith('modules/') ? m.path : 'modules/' + m.path) + '`';

    parts.push('| ' + String(mNum).padStart(2, '0') + '  | **' + m.title + '** | ' + phase + ' | ' + m.status + ' | ' + mLessons + ' | ' + mMilestones + ' | ' + mProjects + ' | ' + mAssessments + ' | ' + mTotal + ' | ' + pathCol + ' |');
  });

  // Electives row in module map
  const eModule = electives.modules[0];
  parts.push('| 17  | **' + eModule.title + '** | ' + eModule.phase + ' | ' + eModule.status + ' | ' + electiveItems + ' | 0 | 0 | 0 | ' + electiveItems + ' | `curriculum/path.electives.json` |');

  parts.push('| —   | **Core Subtotal** | — | — | **' + totLessons + '** | **' + totMilestones + '** | **' + totProjects + '** | **' + totAssessments + '** | **' + totItems + '** | *(Modules 00–16, 18)* |');
  parts.push('| —   | **Total Curriculum** | — | — | **' + (totLessons + electiveItems) + '** | **' + totMilestones + '** | **' + totProjects + '** | **' + totAssessments + '** | **' + (totItems + electiveItems) + '** | *(' + totItems + ' core items + ' + electiveItems + ' electives)* |\n');

  // Workload Estimates & Pacing Guide
  parts.push('## Workload Estimates & 52-Week Pacing Guide\n');
  parts.push('Curriculum completion time cannot be calculated from reading time alone. A realistic workload accounts for conceptual reading, hands-on lab exercises, independent debugging, test authoring, and capstone engineering.\n');
  parts.push('### 1. Workload by Activity Layer (Core Path)\n');
  const lessonItemsCount = lessons + labs;
  const lessonMins = coreItems.filter(it => it.type === 'lesson' || it.type === 'lab').reduce((acc, it) => acc + (it.estimated_minutes || 45), 0);
  const lessonHrs = (lessonMins / 60).toFixed(1);
  const totalCoreItems = lessonItemsCount + 19 + 8 + (checkpoints + capstones) + pathAssessments;
  const totalCoreHrs = (parseFloat(lessonHrs) + 38.0 + 16.0 + 30.0 + 57.0).toFixed(1);

  parts.push('| Activity Layer | Items | On-Paper Metadata | Realistic Learner Effort | Notes |');
  parts.push('| :--- | :---: | :---: | :---: | :--- |');
  parts.push('| **Instructional Lessons & Labs** | ' + lessonItemsCount + ' | ' + lessonHrs + ' hrs | 235 – 265 hrs | Reading, mental models, running code examples, review questions (Modules 00–16) |');
  parts.push('| **Mini-Projects** | 19 | 38.0 hrs | 40 – 75 hrs | Immediate practical reinforcement (2–4 hours per project) |');
  parts.push('| **Portfolio Projects** | 8 | 16.0 hrs | 65 – 120 hrs | Multi-component services with tests, migrations, CI/CD (8–15 hours per project) |');
  parts.push('| **Flagship Opslane Milestones** | ' + (checkpoints + capstones) + ' | 30.0 hrs | 50 – 80 hrs | Integrated production system with auth, DB, payments, worker pools, metrics (Module 18) |');
  parts.push('| **In-Path Assessments & Defense** | ' + pathAssessments + ' | 57.0 hrs | 45 – 60 hrs | Self-evaluation checklists, code reviews, and capstone oral/technical defense |');
  parts.push('| **Core Total** | **' + totalCoreItems + '** | **' + totalCoreHrs + ' hrs** | **435 – 600 hrs** | **Realistic full engineering workload across all 18 core modules** |\n');

  parts.push('### 2. Weekly Commitment Planning Models\n');
  parts.push('| Weekly Commitment | Total Weeks | Total Months | Viability & Pace Description |');
  parts.push('| :---: | :---: | :---: | :--- |');
  parts.push('| **12–15 hrs/week** | **40 – 52 weeks** | **10 – 12 mos** | **(Recommended)** Completes 435–600 hours of core backend engineering in 29–50 weeks of focused study, leaving deliberate slack for review, autonomous project debugging, and parallel companion frontend progression. |');
  parts.push('| **8–10 hrs/week** | **55 – 65 weeks** | **13 – 15 mos** | Steady progress; core concepts finish in 1 year, with portfolio and Opslane defense extending into months 13–15. |');
  parts.push('| **5 hrs/week** | **90 – 115 weeks** | **21 – 26 mos** | Casual part-time; learners reach Module 08/09 by month 12 and graduate in year 2. |\n');

  parts.push('### 3. Integrated 52-Week Schedule Framework\n');
  parts.push('| Weeks | Phase | Focus Modules | Key Practical Milestones |');
  parts.push('| :---: | :--- | :--- | :--- |');
  parts.push('| **1–3** | Foundations | Modules 00–02 | Environment verification, Git workflow, Shell + Git Lab, Tooling Failure Lab |');
  parts.push('| **4–8** | Core Go Programming | Modules 03–04 | Syntax, memory models, functions, error handling, Pricing Checkout, Order Summary Refactor |');
  parts.push('| **9–14** | Engineering Disciplines | Modules 05–07 | Types, interfaces, test-driven refactoring, CLI tools, Config Parser *(Companion: HTML/CSS)* |');
  parts.push('| **15–24** | Backend & Data Persistence | Modules 08–09 | HTTP APIs, Postgres, migrations, transactions, Book Catalog REST API, Task API *(Companion: JS/DOM)* |');
  parts.push('| **25–33** | Security, Lifecycle & Concurrency | Modules 10–11 | JWT/RBAC, goroutines, channels, worker pools, Secure Authenticated API *(Companion: templ/HTMX)* |');
  parts.push('| **34–42** | Production Engineering | Modules 12–15 | Observability (slog/OTel), profiling (pprof), Docker multi-stage, CI/CD pipelines, deployment runbook |');
  parts.push('| **43–50** | Flagship Implementation | Module 18 | Flagship Opslane: multi-tenant architecture, event bus, caching, graceful shutdown, production deployment |');
  parts.push('| **51–52** | Portfolio & Graduation Review | Module 16 | Portfolio packaging, architecture defense rehearsal, resume project bullets, mock interview prep |\n');

  // Repository structure
  parts.push('## Repository Structure\n');
  parts.push('```\ncurriculum/\n  modules/\n    00-orientation/\n    01-computers-terminal-git-web/\n    02-go-setup-tooling/\n    03-programming-fundamentals/\n    04-functions-errors-data-semantics/\n    05-types-interfaces-packages-modules/\n    06-testing-debugging-refactoring/\n    07-cli-files-json-config/\n    08-http-rest-apis/\n    09-sql-postgres-persistence/\n    10-auth-security/\n    11-lifecycle-context-concurrency/\n    12-observability-diagnostics/\n    13-performance-memory-engineering/\n    14-architecture-distributed-systems/\n    15-docker-cicd-deployment/\n    16-portfolio-interview-readiness/\n    18-flagship-opslane/\n  companion-frontend/\n    README.md\n  assessments.json\n  concepts.json\n  crossrefs.json\n  failures.json\n  path.core.json\n  path.electives.json\n  projects.json\n  readme.contracts.json\n  schema.v3.json\n  GRADUATION_CRITERIA.md\n  PILOT_TESTING_PROTOCOL.md\ninternal/tools/curriculum/\n```\n');

  parts.push('## Companion Tracks & Quality Assurance Protocols\n');
  parts.push('- **[Graduation Outcomes & Capstone Defense Rubric](GRADUATION_CRITERIA.md):** 6 observable competency outcomes (unassisted problem-solving, debugging, testing, disaster recovery, architectural defense, incident response) and formal Flagship Opslane oral defense protocol.');
  parts.push('- **[Frontend Companion Track](companion-frontend/README.md):** 4-stage progressive server-driven hypermedia track (`templ` + `HTMX` + `Tailwind` + `Alpine.js`) running parallel to backend modules and converging at Flagship Opslane.');
  parts.push('- **[Beginner Pilot Testing Protocol](PILOT_TESTING_PROTOCOL.md):** Formal empirical zero-intervention testing methodology evaluating unassisted learner velocity and cognitive load across Modules 00–05.\n');

  parts.push('---\n');

  // Detailed core catalog
  parts.push('# Core Path\n');

  core.modules.forEach(m => {
    const mNum = parseInt(m.id.replace('module-', ''), 10);
    const mItems = coreItems.filter(it => it.module_id === m.id && (it.type === 'lesson' || it.type === 'lab' || it.type === 'checkpoint' || it.type === 'capstone'));
    const phase = phaseMap[m.id] || m.phase;

    parts.push('## ' + mNum + '. ' + m.title + '\n');
    parts.push('> **Phase:** ' + phase + ' | **Status:** ' + m.status + ' | **Required:** ' + (m.required !== false ? 'True' : 'False') + '\n');
    parts.push('**Learning Goal:** ' + m.learning_goal + '\n');

    if (m.id === 'module-18') {
      parts.push('**Prerequisites:** 15. Docker, CI/CD, and Deployment (Technical milestone: `core-15-19`)\n');
      parts.push('> **Architecture Note:** Module 16 (Portfolio and Interview Readiness) is intentionally decoupled from technical implementation so learners can construct Flagship Opslane first, producing the production artifact that Module 16 packages for interviews.\n');
    } else if (m.prerequisites && m.prerequisites.length > 0) {
      const prereqTitles = m.prerequisites.map(pid => {
        const pMod = core.modules.find(x => x.id === pid);
        if (pMod) {
          const pNum = parseInt(pMod.id.replace('module-', ''), 10);
          return pNum + '. ' + pMod.title;
        }
        return pid;
      });
      parts.push('**Prerequisites:** ' + prereqTitles.join(', ') + '\n');
    }

    parts.push('| #   | Type    | Slug | Title | Minutes |');
    parts.push('| --- | ------- | ---- | ----- | ------- |');

    // Cluster annotations for consolidated sessions
    const clusterNotes = {
      'core-04-15': ' *(Workshop: Deterministic Cleanup — Part 1)*',
      'core-04-16': ' *(Workshop: Deterministic Cleanup — Part 2)*',
      'core-05-03': ' *(Workshop: Method Receivers & Mutation — Part 1)*',
      'core-05-04': ' *(Workshop: Method Receivers & Mutation — Part 2)*',
      'core-05-05': ' *(Workshop: Method Receivers & Mutation — Part 3)*',
      'core-05-12': ' *(Workshop: Dynamic Type Inspection — Part 1)*',
      'core-05-13': ' *(Workshop: Dynamic Type Inspection — Part 2)*',
      'core-06-14': ' *(Workshop: Delve Debugging Workflow — Step 1)*',
      'core-06-15': ' *(Workshop: Delve Debugging Workflow — Step 2)*',
      'core-06-16': ' *(Workshop: Delve Debugging Workflow — Step 3)*',
      'core-06-17': ' *(Workshop: Delve Debugging Workflow — Step 4)*',
      'core-07-12': ' *(Workshop: In-Memory JSON — Part 1)*',
      'core-07-13': ' *(Workshop: In-Memory JSON — Part 2)*',
      'core-07-14': ' *(Workshop: Streaming JSON — Part 1)*',
      'core-07-15': ' *(Workshop: Streaming JSON — Part 2)*',
      'core-11-01': ' *(Workshop: Time Mechanics & Scheduling — Part 1)*',
      'core-11-02': ' *(Workshop: Time Mechanics & Scheduling — Part 2)*',
      'core-11-03': ' *(Workshop: Time Mechanics & Scheduling — Part 3)*',
      'core-11-04': ' *(Workshop: Time Mechanics & Scheduling — Part 4)*'
    };

    mItems.forEach((it, idx) => {
      const itemType = it.subtype || it.type;
      const note = clusterNotes[it.id] || '';
      parts.push('| ' + (idx + 1) + '   | ' + itemType.padEnd(7, ' ') + ' | ' + it.slug + ' | **' + it.title + '**' + note + ' | ' + (it.estimated_minutes || 45) + ' |');
    });

    parts.push('');
  });

  parts.push('---\n');
  parts.push('# Advanced Electives\n');

  const m17 = electives.modules[0];
  parts.push('## 17. ' + m17.title + '\n');
  parts.push('> **Phase:** ' + m17.phase + ' | **Status:** ' + m17.status + ' | **Required:** False\n');
  parts.push('**Learning Goal:** ' + m17.learning_goal + '\n');
  parts.push('**Prerequisites:** 16. Portfolio and Interview Readiness\n');
  parts.push('| #   | Type      | Slug | Title | Minutes |');
  parts.push('| --- | --------- | ---- | ----- | ------- |');

  electives.items.forEach((it, idx) => {
    const itemType = it.subtype || it.type;
    parts.push('| ' + (idx + 1) + '   | ' + itemType.padEnd(9, ' ') + ' | ' + it.slug + ' | **' + it.title + '** | ' + (it.estimated_minutes || 45) + ' |');
  });

  parts.push('\n---\n');
  parts.push('# Projects\n');
  parts.push('| #   | Module | Title | Type | Stage | Difficulty |');
  parts.push('| --- | ------ | ----- | ---- | ----- | ---------- |');

  // Sort projects by module number order
  const sortedProjects = [...projects.projects].sort((a, b) => {
    const aMod = parseInt((a.module_id || 'module-00').replace('module-', ''), 10);
    const bMod = parseInt((b.module_id || 'module-00').replace('module-', ''), 10);
    if (aMod !== bMod) return aMod - bMod;
    return a.id.localeCompare(b.id);
  });

  sortedProjects.forEach((p, idx) => {
    const modNum = p.module_id ? p.module_id.replace('module-', 'Mod ') : '—';
    parts.push('| ' + (idx + 1) + '   | ' + modNum + ' | **' + p.title + '** | ' + p.type + ' | ' + p.stage + ' | ' + p.difficulty + ' |');
  });

  parts.push('\n---\n');
  parts.push('# Assessments\n');
  parts.push('The assessment framework evaluates mastery across end-of-module milestones and practical project implementations.\n');
  parts.push('### Module Assessments (19 Records)\n');
  parts.push('| ID | Title | Type | Passing Score |');
  parts.push('| -- | ----- | ---- | ------------- |');

  const moduleAssessments = assessments.assessments.filter(a => a.id.startsWith('assessment-module-'));
  moduleAssessments.forEach(a => {
    const score = a.passing_score ? (String(a.passing_score).endsWith('%') ? a.passing_score : a.passing_score + '%') : '80%';
    parts.push('| `' + a.id + '` | **' + a.title + '** | ' + a.type + ' | ' + score + ' |');
  });

  parts.push('\n### Project Assessments (28 Records)\n');
  parts.push('| ID | Title | Type | Passing Score |');
  parts.push('| -- | ----- | ---- | ------------- |');

  const projectAssessments = assessments.assessments.filter(a => !a.id.startsWith('assessment-module-'));
  projectAssessments.forEach(a => {
    const score = a.passing_score ? (String(a.passing_score).endsWith('%') ? a.passing_score : a.passing_score + '%') : '80%';
    parts.push('| `' + a.id + '` | **' + a.title + '** | ' + a.type + ' | ' + score + ' |');
  });

  parts.push('\n---\n');
  parts.push('# Cross-References\n');

  const refs = crossrefs.crossrefs.references;
  const relDist = {};
  refs.forEach(r => { relDist[r.relation] = (relDist[r.relation] || 0) + 1; });

  parts.push('Total cross-references: **' + refs.length + '**\n');
  parts.push('### Style Rules\n');
  parts.push('| Relation      | Template                                                     |');
  parts.push('| ------------- | ------------------------------------------------------------ |');
  parts.push('| builds_on     | Builds on: Module X.Y — Lesson Title                         |');
  parts.push('| preview_only  | Preview only: Explained fully in Module X.Y — Lesson Title   |');
  parts.push('| reinforced_in | Reinforced in: Project X — Project Title                     |');
  parts.push('| related       | Related: Module X.Y — Lesson Title                           |\n');

  parts.push('### Relation Distribution\n');
  parts.push('| Relation      | Count | Description |');
  parts.push('| ------------- | ----- | ----------- |');
  parts.push('| related       | ' + (relDist.related || 0) + '   | Bidirectional conceptual links between related lessons |');
  parts.push('| builds_on     | ' + (relDist.builds_on || 0) + '   | Strict directional pedagogical prerequisites |');
  parts.push('| reinforced_in | ' + (relDist.reinforced_in || 0) + '   | Practical application links to projects and labs |');
  parts.push('| preview_only  | ' + (relDist.preview_only || 0) + '    | Early mental-model previews developed fully later |\n');

  // Preserved migration section
  parts.push('---\n');
  parts.push('# Migration: v2 → v3\n');
  parts.push('**Source Schema:** 2\n');
  parts.push('**Target Schema:** 3.1.0\n');
  parts.push('**Default Policy:** Every legacy item must be mapped to keep, rewrite, split, merge, move-to-elective, replace, or archive. Every migrated README must also receive a legacy_readme_action: reuse, rewrite, split, merge, or archive.\n');
  parts.push('**Frozen Source:** `curriculum/legacy/curriculum.v2.json`\n');
  parts.push('### High-Confidence Moves\n');
  parts.push('| Legacy Pattern        | Outcome          | Target    | Reason                                                                    |');
  parts.push('| --------------------- | ---------------- | --------- | ------------------------------------------------------------------------- |');
  parts.push('| API.4-API.9           | move-to-elective | module-17 | Protobuf/gRPC is valuable but should not block the core beginner path.    |');
  parts.push('| TI.14-TI.15           | move-to-elective | module-17 | Complex generics and generic data structures are advanced.                |');
  parts.push('| TI.12                 | move-to-elective | module-17 | Functional options are useful but not core beginner material.             |');
  parts.push('| CP.3                  | move-to-elective | module-17 | sync.Pool is performance-specialized.                                     |');
  parts.push('| SY.2 sync.Map portion | move-to-elective | module-17 | sync.Map is specialized and should follow normal map+mutex reasoning.     |');
  parts.push('| ARCH.7                | move-to-elective | module-17 | CQRS is advanced architecture.                                            |');
  parts.push('| DEPLOY.2              | move-to-elective | module-17 | Blue/green deployment is advanced deployment strategy.                    |');
  parts.push('| TE.1-TE.11            | rewrite          | module-06 | Testing and debugging must move before backend work.                      |');
  parts.push('| DB.*                 | rewrite          | module-09 | Database path must become SQL/Postgres-first with migrations and EXPLAIN. |');
  parts.push('| SEC.*                | rewrite          | module-10 | Security needs authn/authz strategy and secure API practice.              |\n');

  parts.push('### Quality Gates\n');
  parts.push('- schema-valid\n- graph-valid\n- no-orphan-core-items\n- all-prerequisites-exist\n- no-circular-dependencies\n- item-orders-sequential\n- phase-sequence-valid\n- no-null-arrays\n- crossref-targets-exist\n- crossref-targets-in-same-bundle\n- module-prerequisites-exist\n- item-prerequisites-exist\n- project-reinforces-items-exist\n- assessment-targets-exist\n- project-assessment-exists\n- no-duplicate-id\n- no-duplicate-slug\n- no-unused-crossref-item\n- item-has-zero-magic\n- item-has-proof\n- item-has-content-contract\n- item-has-verification\n- item-has-files\n- item-files-path-valid\n- assessment-retake-policy-present\n');

  parts.push('### Validation Commands\n');
  parts.push('| Command          | Script                                              |');
  parts.push('| ---------------- | --------------------------------------------------- |');
  parts.push('| validate_all     | `go run ./internal/tools/curriculum validate-all`     |');
  parts.push('| validate_graph   | `go run ./internal/tools/curriculum validate-graph`   |');
  parts.push('| validate_lessons | `go run ./internal/tools/curriculum validate-lessons` |');
  parts.push('| validate_readmes | `go run ./internal/tools/curriculum validate-readmes` |');
  parts.push('| validate_schema  | `go run ./internal/tools/curriculum validate-schema`  |\n');

  parts.push('---\n');
  parts.push('_Generated from curriculum v' + core.curriculum_version + ' | Schema ' + core.schema_version + ' | ' + core.last_updated + '_\n');

  return parts.join('\n');
}

const out = generateFullReadme();
fs.writeFileSync('curriculum/README.md', out, 'utf8');
console.log('Successfully generated curriculum/README.md! Bytes: ' + Buffer.byteLength(out, 'utf8') + ', lines: ' + out.split('\n').length);
