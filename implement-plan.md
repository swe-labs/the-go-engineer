# The Go Engineer v3 — Curriculum Metadata Finalization and Full Instructional Implementation

## Role

Act as the principal curriculum architect, senior Go engineer, instructional designer, technical editor, assessment designer, and repository quality engineer responsible for bringing **The Go Engineer v3** from its current metadata/scaffold state to a coherent, technically accurate, professionally authored, testable learning curriculum.

Repository: `https://github.com/swe-labs/the-go-engineer`

Target branch: `release/v3`

Primary scope:

- `curriculum/`
- `curriculum/modules/`
- All authoritative curriculum manifests, registries, schemas, validators, and documentation that must change to keep the curriculum consistent.

This is a real curriculum design and implementation project, not a metadata beautification exercise.

## 1. Current state and non-negotiable understanding

The curriculum contains substantial structural infrastructure, manifests, lesson IDs, metadata records, some partially written lessons, examples, tests, and project scaffolds. However, **the existence of a record, file, title, description, or passing validator must not be treated as proof that the lesson or curriculum item is complete or correct**.

Many existing fields may be placeholders, generic template text, copied metadata, incomplete drafts, or content created primarily to satisfy a schema or validator. Existing lessons may also be incomplete, pedagogically weak, technically incorrect, redundant, or disconnected from their declared outcomes.

Assume nothing is correct merely because it exists.

Audit the actual repository before deciding what to preserve, revise, replace, merge, remove, or implement.

The goal is to produce a curriculum that is:

- Coherent from absolute beginner foundations to advanced professional Go backend engineering.
- Technically accurate, idiomatic, secure, maintainable, and relevant to real engineering work.
- Designed around observable learner capabilities rather than lesson-count inflation.
- Pedagogically sequenced, with prerequisites that reflect actual knowledge dependencies.
- Rich in meaningful practice, debugging, testing, implementation, reflection, and assessment.
- Consistent between metadata and physical instructional files.
- Verifiable through runnable code, automated tests, explicit acceptance criteria, and reproducible commands where applicable.
- Maintainable by future contributors without duplicating hundreds of manually synchronized rules.

Do not claim the curriculum is complete, stable, production-ready, industry-leading, or empirically validated unless the repository contains sufficient evidence to support that claim.

## 2. Repository documents to inspect

Start with the following files and their actual contents on `release/v3`:

- `curriculum/README.md`
- `curriculum/workspace.json`
- `curriculum/schema.v3.json`
- `curriculum/path.core.json`
- `curriculum/path.electives.json`
- `curriculum/projects.json`
- `curriculum/assessments.json`
- `curriculum/concepts.json`
- `curriculum/failures.json`
- `curriculum/crossrefs.json`
- `curriculum/readme.contracts.json`
- `curriculum/legacy/migration.v2-to-v3.json`
- `curriculum/contributor-authoring-guide.md`
- `curriculum/zero-magic-authoring-spec.md`
- `curriculum/quality-rubric.md`
- `curriculum/GRADUATION_CRITERIA.md`
- `curriculum/PILOT_TESTING_PROTOCOL.md`
- `curriculum/companion-frontend/README.md`
- All existing curriculum tooling under `internal/tools/curriculum/`
- All actual lesson, lab, project, assessment, starter, solution, test, and supporting files beneath `curriculum/modules/` and the relevant elective directories.

Inspect relevant root-level documentation, tests, scripts, CI workflows, and existing project code when they affect curriculum contracts.

Do not rely on the README as the sole representation of the curriculum. Do not rely on JSON manifests as the sole representation of actual teaching content.

The workspace manifest says the split curriculum files are authoritative. Verify that this is consistent with the implementation and preserve that policy unless a documented architectural reason justifies changing it.

## 3. Rules for protecting the repository

Before editing:

1. Confirm the current branch and working-tree status.
2. Record the current commit SHA.
3. Inspect existing uncommitted changes. Never overwrite, reset, discard, or reformat unrelated user work.
4. Identify the actual validation commands and how the project is built and tested.
5. Create a baseline report containing the current failures, warnings, file counts, and known inconsistencies.

Do not run destructive Git commands.

Do not change branches, rewrite history, force-push, publish a release, or create a release tag.

Do not make a huge unreviewable rewrite of the entire curriculum in one operation.

Keep changes focused and reviewable. Work in small logical batches, run the relevant checks after each batch, and maintain a persistent implementation plan so work can continue safely across sessions.

Do not claim to have run commands, reviewed files, verified behavior, or passed tests unless you actually did so.

## 4. Required workflow: phase-gated implementation

Follow the phases below in order. Each phase must produce concrete artifacts and have explicit acceptance criteria.

Do not begin by writing all the lessons. First establish what the curriculum is supposed to teach and how its records relate.

### Phase 0 — Full repository and curriculum audit

Perform a systematic inventory of all curriculum metadata and physical instructional artifacts.

#### 0.1 Metadata inventory

For every core and elective item, collect at least:

- Stable item ID.
- Module ID and module title.
- Order and phase.
- Item type.
- Current title and description.
- Learning objective.
- Declared prerequisites.
- Estimated duration.
- Difficulty.
- Tags and concepts.
- Content contract.
- Assessment and project bindings.
- Cross-references.
- Declared physical paths.
- Current implementation status.
- Existing runnable example and test paths.
- Evidence of placeholder or template data.
- Any related validation or rendering rules.

Include projects, labs, checkpoints, assessments, concept ownership, operational failure coverage, migration records, and graduation requirements where applicable.

Do not silently omit large manifests or modules because they are inconvenient to inspect. Build scripts or structured reports to inspect the full data set when needed.

#### 0.2 Physical artifact inventory

Inspect the real files, not only their paths in JSON.

Determine whether each item has:

- A meaningful README or equivalent primary lesson.
- Actual instructional prose rather than headings and template text.
- Clear prerequisites and learning objectives.
- A useful mental model and visual or textual execution model where appropriate.
- Correct explanations, worked examples, and expected outputs.
- Meaningful guided practice and independent practice.
- Common mistakes and debugging guidance.
- Review questions or other checks for understanding.
- A practical artifact or runnable example when justified.
- Tests where appropriate.
- A real assessment and explicit mastery criteria.
- Accurate commands and working internal links.

Distinguish genuinely incomplete files from legitimate non-runnable conceptual lessons. Do not demand Go code or automated tests for every item if they do not fit the instructional purpose.

#### 0.3 Placeholder and semantic quality audit

Search the entire curriculum for:

- Generic titles and descriptions.
- Repeated learning objectives with superficial wording changes.
- Text such as “TODO”, “TBD”, “fill this in”, “This module builds on concepts from the previous module”, and similar placeholders.
- Incomplete sentences and blank concepts in cross-reference reasons.
- Copy-pasted content that does not match the lesson.
- Repeated or contradictory examples.
- Unjustified or impossible durations.
- Placeholder assessments and generic rubrics.
- Broken or invented file references.
- Missing prerequisites.
- Prerequisites that point forward or require unintroduced knowledge.
- Incorrect project/module/assessment relationships.
- Claims of implementation or stability unsupported by real files.
- Metadata fields that exist only to satisfy schema validation.
- Contradictions between a lesson's stated objective and its actual content.

Use semantic judgment as well as text searches. A generic sentence may be a placeholder even when it does not contain a known placeholder keyword.

#### 0.4 Consistency audit

Reconcile:

- Module numbering, names, paths, and status.
- Core and elective membership.
- Lesson IDs and physical files.
- Prerequisite and progression graphs.
- Project and assessment references.
- Concept ownership and coverage.
- Cross-reference targets and human-readable labels.
- Manifest counts and README counts.
- Workload estimates and pacing claims.
- Graduation criteria and actual assessment evidence.
- Schema version and curriculum release version.
- Migration mappings and current item identities.
- Declared quality gates and the checks that actually exist.
- Generated documentation and authoritative sources.

The physical module layout currently includes Modules 00–16 and 18, while Module 17 is represented as an elective track in metadata. Verify that the model is intentional and consistently represented everywhere.

#### 0.5 Audit deliverables

Create or update a persistent plan such as `curriculum/IMPLEMENTATION_PLAN.md`, choosing the most appropriate existing planning location if one already exists.

Include:

- Audit scope and date.
- Baseline commit SHA.
- Inventory counts derived from actual data.
- Placeholder counts and representative examples.
- Missing or incomplete instructional artifacts.
- Structural and semantic inconsistencies.
- Curriculum design questions that require a decision.
- Recommended changes grouped by priority.
- Proposed phase order and module batches.
- Risks, dependencies, and acceptance criteria.
- A progress checklist that can be resumed across sessions.

Use these classifications:

- **P0 — Structural blocker:** IDs, prerequisites, architecture, or data contracts prevent a coherent curriculum.
- **P1 — Core educational defect:** Missing foundations, major duplication, incorrect sequencing, inaccurate instruction, or absent essential assessment.
- **P2 — Content-quality defect:** Weak explanation, inadequate examples, generic metadata, or incomplete practice.
- **P3 — Improvement:** Nonessential polish or optional enrichment.

Do not start mass metadata or content rewriting until the audit identifies the intended scope and the structural blockers.

### Phase 1 — Finalize the curriculum blueprint

Establish the curriculum-wide educational design before individually polishing every lesson.

#### 1.1 Define the learner and graduate

Document:

- Assumed starting knowledge.
- Target learner.
- Intended graduate capabilities.
- Core technical boundaries.
- What the curriculum explicitly does not promise.
- Required evidence for graduating.
- The role of optional electives.
- The role of the frontend companion track.

Maintain the intended zero-to-professional progression without assuming that reading lessons automatically produces professional competence.

#### 1.2 Review module architecture

For every module, finalize:

- Canonical number, ID, name, and directory.
- Purpose and scope.
- Measurable module-level outcomes.
- Prerequisite modules and knowledge.
- Main topics and explicit boundaries.
- Why the module occurs at its current point.
- Practical outputs.
- Assessment strategy.
- Approximate workload.
- Dependencies on projects and other modules.
- Topics intentionally deferred to later modules or electives.

Identify missing foundations, redundant modules, overloaded modules, misplaced advanced concepts, and gaps in the progression.

Do not add modules or lessons simply to increase curriculum size. Merge, remove, or defer content when there is a defensible educational reason. Preserve useful existing material and stable identities where possible.

#### 1.3 Validate the overall learning progression

Review the progression across:

- Computer, terminal, Git, and web foundations.
- Go setup and programming fundamentals.
- Functions, errors, data semantics, types, interfaces, and packages.
- Testing, debugging, and refactoring.
- CLI applications, files, JSON, and configuration.
- HTTP and APIs.
- SQL, PostgreSQL, migrations, and transactions.
- Authentication, authorization, and security.
- Context, lifecycle, concurrency, and resource ownership.
- Observability and diagnostics.
- Performance and memory engineering.
- Architecture and distributed systems.
- Containers, CI/CD, deployment, and recovery.
- Flagship Opslane.
- Portfolio and interview readiness.
- Optional advanced electives.
- Frontend companion skills.

These are audit areas, not an instruction to add every topic indiscriminately. Validate their placement, depth, and prerequisites against the actual curriculum and target graduate.

Consider whether foundational algorithmic reasoning, generics fundamentals, recovery exercises, negative security tests, browser end-to-end testing, and other essential competencies have appropriate coverage. Avoid forcing advanced or niche topics into the beginner core without a clear reason.

#### 1.4 Freeze structural decisions

Once reviewed, record the canonical architecture and decisions in the appropriate documentation.

Freeze, unless a genuine defect is discovered:

- Stable IDs.
- Module boundaries.
- Core/elective membership.
- Intended lesson order.
- Main prerequisite relationships.
- Canonical directory conventions.
- Metadata ownership rules.

A freeze is not permission to preserve an incorrect design. Structural changes later require a documented reason and impact analysis.

**Phase 1 acceptance:** A coherent module map, documented graduate outcomes, resolved structural blockers, and an approved progression model. No mass lesson writing is necessary yet.

### Phase 2 — Finalize all lesson and instructional-item metadata

Now audit and author the metadata for the complete intended curriculum, before attempting to implement every lesson.

Work through the core path and electives using the canonical manifests. Do not confuse completing metadata with completing lesson content.

#### 2.1 Author meaningful metadata

For each lesson, lab, checkpoint, and other instructional item, provide fields appropriate to its type.

At minimum, establish:

- **Title:** Specific, clear, and recognizable. Avoid vague or generic names.
- **Description or summary:** Explain what the learner will study and why it matters.
- **Learning objective:** State observable behavior the learner should be able to demonstrate.
- **Prerequisites:** Name only knowledge genuinely needed before starting.
- **Difficulty:** Reflect the actual conceptual and practical demands.
- **Estimated duration:** Include the intended learning activities, not just reading time.
- **Instructional role:** Explain whether the item introduces, develops, applies, integrates, or assesses knowledge.
- **Concept coverage:** Link to the concepts actually taught.
- **Content contract:** Specify the sections and artifacts that the lesson genuinely needs.
- **Practice expectations:** Define meaningful learner activity.
- **Assessment evidence:** Explain what proves the objective was achieved.
- **Next steps and cross-references:** Link to real, relevant items with specific reasons.
- **Physical paths:** Match the canonical file layout.
- **Project bindings:** Match actual project scope and progression.
- **Status:** Reflect actual implementation maturity, not intention.

Retain additional schema-required fields where appropriate. Review every field's semantics before populating it.

Do not copy the same learning objective, practice task, mastery check, or cross-reference explanation across unrelated lessons.

#### 2.2 Make objectives measurable

Prefer objectives such as:

“Given a function that returns a wrapped error, identify how to inspect the error chain and write a test that verifies the underlying error with `errors.Is`.”

Avoid objectives such as:

“Understand errors in Go.”

The first describes observable evidence. The second is too vague to assess reliably.

Objectives must match the lesson's actual scope and expected learner level.

#### 2.3 Make content contracts instructional, not mechanical

Define the required content sections based on the learning task.

A conceptual lesson might need a mental model, a worked trace, common misconceptions, and review questions.

A programming lesson might also need runnable examples, exercises, tests, and explicit verification commands.

A debugging lab needs a reproducible failure, diagnostic steps, evidence of root cause, and a verified fix.

A project needs a specification, constraints, acceptance criteria, staged tasks, tests, and a rubric.

Do not require every possible section for every item. Do not let a template dictate content that has no instructional purpose.

#### 2.4 Repair all curriculum relationships

For each item, verify:

- Prerequisites exist and form a valid learning progression.
- Dependencies are not accidentally circular.
- Concepts are owned and taught in a meaningful location.
- Cross-references target existing IDs and explain the actual relationship.
- Assessment criteria evaluate the declared objectives.
- Projects reinforce concepts the learner has already encountered.
- Projects do not require unexplained future concepts.
- Electives do not accidentally block core completion.
- Physical file paths and starter/solution/test paths are valid.
- IDs remain stable unless a change is explicitly justified.

Do not fabricate a relationship to make a validator pass.

#### 2.5 Metadata is complete only when it is meaningful

A field is not complete merely because it is nonempty or has the right data type.

Before accepting a record, ask:

1. Is the value specific to this lesson?
2. Is it technically correct?
3. Does it describe what the learner will actually do?
4. Does it fit the module's scope and level?
5. Is it consistent with the prerequisite graph?
6. Can the learning outcome be assessed?
7. Will the physical lesson be able to satisfy the metadata honestly?

Where the schema lacks useful constraints, improve the relevant validator or quality check instead of lowering the standard.

#### 2.6 Review global registries and assessments

Inspect and revise the concept registry, failure taxonomy, cross-references, projects, assessments, and migration data.

Keep these registries only where they provide real value. Remove redundant or unused records when appropriate, with reference updates and a documented rationale.

Do not invent coverage relationships solely to make a coverage report appear complete.

#### 2.7 Metadata phase deliverables

Produce:

- A finalized canonical module and lesson catalog.
- Corrected core and elective manifests.
- Consistent project and assessment bindings.
- Corrected concept and failure mappings.
- Valid and meaningful cross-references.
- A field-level placeholder and semantic-quality report.
- Updated curriculum documentation or generated output where appropriate.
- Validation evidence.

Do not label metadata as finalized if critical decisions remain unresolved. Record legitimate future content-dependent adjustments explicitly.

**Phase 2 acceptance:** Every planned instructional item has meaningful, internally consistent metadata and a clear implementation contract. All structural references resolve. Metadata-quality checks pass, or documented exceptions explain precisely what remains and why.

### Phase 3 — Implement the real instructional content, one module at a time

After the curriculum blueprint and metadata catalog are established, implement the physical curriculum in prerequisite order.

Do not attempt to author all modules in one enormous operation.

Start with Module 00, complete it to the agreed standard, then progress through the dependency graph. Prefer small lesson batches inside a module when that improves reviewability. Prioritize prerequisites before dependent lessons.

#### 3.1 Begin each module with a readiness check

Before writing its lessons:

- Read the module's approved outcomes and scope.
- Read every lesson's metadata and content contract.
- Review relevant prerequisites and earlier lessons.
- Review project and assessment relationships.
- Identify required technical verification.
- Check for duplicated material and missing foundations.
- Confirm that the module's metadata is ready to implement.

If the intended content exposes a genuine architectural flaw, document the issue and make the smallest justified metadata change. Do not blindly implement a bad specification.

#### 3.2 Write real instructional material

Every lesson must teach the topic, not merely announce it.

Use a clear instructional progression appropriate to the learner:

1. State the problem and why the topic matters.
2. Establish the prerequisite knowledge.
3. Explain the mental model.
4. Show how the system or language actually behaves.
5. Work through concrete examples.
6. Explain important edge cases and failure modes.
7. Provide guided practice.
8. Provide independent practice where appropriate.
9. Explain common mistakes and diagnostic techniques.
10. Verify understanding using meaningful questions, exercises, or tests.
11. Connect the topic to later engineering work.
12. Provide the next learning step.

This is a guide, not a mandatory identical template. Use the structure that best teaches the specific subject.

Write for a learner who may be encountering the idea for the first time. Define terminology before depending on it. Explain why code works rather than asking the learner to memorize it.

Avoid unexplained jumps in difficulty, unexplained abbreviations, large code dumps, and advanced abstractions introduced before their prerequisites.

#### 3.3 Maintain technical accuracy

For Go content:

- Follow the Go version and language semantics supported by the repository.
- Prefer idiomatic, maintainable Go.
- Explain language behavior precisely.
- Verify API signatures, standard-library behavior, and version-sensitive claims.
- Use the official documentation and reliable technical references when verification is needed.
- Do not invent compiler behavior, runtime guarantees, security properties, or performance results.
- Distinguish specification guarantees from implementation details.
- Avoid teaching implementation-specific internals as universal rules.
- Use concurrency examples that have explicit ownership, cancellation, synchronization, and error-handling semantics.
- Explain security boundaries and failure behavior, not just successful requests.
- Explain performance claims with measurement or defensible reasoning rather than assumptions.

The curriculum's zero-magic philosophy should make behavior understandable. It must not become an excuse to burden beginners with irrelevant runtime internals before they have the foundations needed to understand them.

#### 3.4 Make examples reproducible

When a lesson requires runnable code:

- Use the repository's established directory and package conventions.
- Include all necessary source files and dependencies.
- Specify the working directory and exact commands.
- Include expected output or verification criteria.
- Provide meaningful tests.
- Verify that the example compiles and behaves as described.
- Ensure starter code and reference solutions match the exercise.
- Avoid unexplained dependencies on files from other lessons.
- Avoid tests that pass without checking the intended behavior.
- Ensure failure examples reproduce the described failure.

Use the repository's existing starter/solution/test architecture where appropriate. Do not create a second competing architecture without a documented reason.

For conceptual lessons that do not need runnable code, provide other meaningful evidence of understanding instead.

#### 3.5 Make exercises progressively demanding

Use a deliberate mix of:

- Recall and explanation.
- Code reading and prediction.
- Execution tracing.
- Small modifications.
- Independent implementation.
- Debugging.
- Test writing.
- Refactoring.
- Failure injection.
- Integration tasks.
- Design and tradeoff explanation.

Do not turn every lesson into a miniature production system. Keep the challenge proportional to the learning objective.

Exercises should require the learner to demonstrate understanding rather than copy a complete solution.

#### 3.6 Assessments must prove the objective

Every assessment should have:

- A clear task or prompt.
- Explicit constraints.
- Observable acceptance criteria.
- Meaningful success and failure cases.
- Appropriate hints or remediation guidance.
- A rubric that distinguishes partial understanding from mastery.
- Verification instructions where applicable.

Avoid generic rubrics such as “the solution works” or “the learner understands the concept.”

Do not use one generic assessment to claim mastery of unrelated lessons. Ensure assessment criteria map to the actual learning outcomes.

#### 3.7 Projects must integrate previously taught skills

For each project:

- Define the problem and intended user or system behavior.
- State constraints and requirements.
- Identify prerequisite knowledge.
- Provide staged milestones.
- Define acceptance criteria.
- Specify required tests and evidence.
- Include relevant failure cases.
- Explain the expected artifact.
- Provide a rubric with observable levels of performance.
- Ensure the project is feasible for the intended learner.
- Ensure it reinforces concepts already taught.

Do not make a project depend on substantial unintroduced functionality merely because it creates a more impressive repository.

#### 3.8 Finish each module completely

Before moving to the next module:

- Complete its planned lessons and labs to their declared contracts.
- Verify runnable examples and tests where required.
- Complete its assessments and project bindings.
- Update metadata if implementation revealed a legitimate mismatch.
- Validate internal links and paths.
- Review prerequisites and cross-references.
- Run all relevant validators.
- Perform a semantic quality review.
- Update the persistent progress plan with completed items and unresolved issues.

A module is not complete because its directory contains files.

### Phase 4 — Strengthen validation and quality controls

Review existing curriculum tooling before creating new tools.

The repository's workspace lists validation commands including:

- `go run ./internal/tools/curriculum validate-schema`
- `go run ./internal/tools/curriculum validate-graph`
- `go run ./internal/tools/curriculum validate-lessons`
- `go run ./internal/tools/curriculum validate-prerequisite-closure`
- `go run ./internal/tools/curriculum validate-project-binding`
- `go run ./internal/tools/curriculum validate-no-placeholder-zm`
- `go run ./internal/tools/curriculum validate-assessment-completeness`
- `go run ./internal/tools/curriculum validate-unique-tags`
- `go run ./internal/tools/curriculum validate-crossrefs`
- `go run ./internal/tools/curriculum validate-project-graph`
- `go run ./internal/tools/curriculum validate-all`

Verify these commands against the actual repository and use the appropriate available commands. Do not assume that every listed command exists or behaves as advertised without checking.

Where practical, add or improve automated checks for:

- Missing and duplicate IDs.
- Missing, invalid, or cyclic prerequisites.
- Invalid physical paths.
- Unresolved cross-references.
- Placeholder titles and descriptions.
- Generic or incomplete cross-reference reasons.
- Objectives and assessment criteria that are missing or mismatched.
- Incorrect project-to-module relationships.
- Missing required lesson sections.
- Missing runnable examples or tests where required.
- Starter/solution/test parity.
- Stale module counts and generated documentation.
- Incorrect core/elective classifications.
- Schema-version compatibility.
- Orphaned records and unused registry entries.
- Inconsistent metadata and physical lesson status.

Use deterministic checks for structural rules. Use targeted semantic review for instructional quality. Do not pretend simple keyword matching can prove a lesson is pedagogically excellent.

Keep the validator maintainable. Prefer a small number of well-defined rules over dozens of overlapping checks that are difficult to understand.

### Phase 5 — Curriculum-wide integration and reconciliation

After the module-by-module implementation, audit the entire curriculum again.

Verify:

- The learning path is coherent from beginning to end.
- No essential prerequisite is missing.
- No lesson is unnecessarily duplicated.
- Similar topics have distinct and justified instructional roles.
- Advanced content appears only when the learner has the necessary foundations.
- The core and electives are clearly separated.
- The frontend companion track aligns with the overall progression.
- Projects and assessments build cumulatively.
- The flagship project uses previously taught concepts appropriately.
- Every claimed competency has appropriate assessment evidence.
- Workload estimates reflect real learning and project effort.
- Graduation criteria are observable and defensible.
- Pilot-testing documentation does not claim results that have not been collected.
- README counts are generated from or reconciled with authoritative records.
- No obsolete status or release claim remains.
- No placeholder or template data is misrepresented as completed work.

Update the curriculum README, workspace metadata, graduation criteria, and other documentation only after the underlying facts have been verified.

Distinguish planned, scaffolded, authored, verified, and empirically tested content. A validator passing does not establish that learners can complete the curriculum or achieve its intended outcomes.

## 5. Standards for metadata and content quality

Apply the following standards throughout the project.

### Titles

Specific, concise, technically accurate, and distinct enough to communicate the lesson's purpose.

### Descriptions

Explain the topic, its relevance, and the learner's activity. Avoid generic summaries that could be pasted into any lesson.

### Learning objectives

Observable, appropriately scoped, and assessable.

### Prerequisites

Necessary, minimal, explicit, and supported by the actual teaching sequence.

### Durations

Realistic estimates that distinguish reading, guided practice, independent work, debugging, and project effort when useful. Do not invent precision unsupported by the task.

### Cross-references

Accurate targets and specific reasons explaining what is reused, extended, contrasted, or reinforced.

### Content

Correct, coherent, beginner-accessible, substantial enough to teach the topic, and aligned with the metadata.

### Examples

Runnable when appropriate, technically correct, reproducible, and useful for learning.

### Exercises

Purposeful, progressively demanding, and aligned with the objective.

### Assessments

Evidence-based, specific, reproducible, and tied to explicit acceptance criteria.

### Status

A truthful description of actual completion, not an aspiration or a way to make a validation gate pass.

## 6. Technical and architectural constraints

- Preserve the existing schema and manifest architecture where it remains sound.
- Use `curriculum/schema.v3.json` and the documented schema-version policy as the baseline. Do not confuse schema version with curriculum release version.
- Do not silently rewrite the curriculum's versioning policy.
- Do not manually edit generated files if authoritative source files and generators are intended to control them.
- Do not create duplicate sources of truth.
- Do not change stable IDs merely to make numbering look cleaner.
- If an ID must change, trace every reference and document the reason.
- Preserve useful historical and migration data, correcting it only when required for consistency.
- Do not delete an existing lesson, project, assessment, or registry record without reviewing its educational value and dependencies.
- Do not add a framework, dependency, tool, or technology solely because it is popular.
- Do not make optional advanced material a hidden prerequisite for the core.
- Keep frontend companion scope intentional and aligned with the main curriculum.
- Keep technical claims proportionate to evidence.
- Do not replace missing educational content with verbose boilerplate.

## 7. Workload and curriculum scope

The existing README proposes a 52-week framework and estimates 435–600 hours for the core path. Treat those numbers as claims to audit, not immutable truths.

Recalculate workload using the actual lesson activities, labs, projects, assessments, debugging, and capstone requirements.

Explain the difference between:

- Estimated instructional time.
- Independent practice.
- Project implementation.
- Review and remediation.
- Optional elective time.
- Parallel frontend companion work.

Avoid promising that every beginner will become job-ready in a fixed number of weeks. Preserve the intended ambitious learning target while presenting workload and outcomes honestly.

## 8. Work management and progress reporting

Maintain the implementation plan throughout the work.

For each completed batch, report:

- Files changed.
- Lessons or records completed.
- Structural decisions made.
- Placeholder records removed or rewritten.
- Validation commands actually run and their results.
- Example/test execution evidence.
- Remaining defects.
- Any unresolved decision that affects later modules.
- The next logical batch.

Use a clear progress state such as:

- Not audited.
- Audited.
- Metadata drafted.
- Metadata reviewed.
- Content in progress.
- Content authored.
- Technically verified.
- Pedagogically reviewed.
- Module complete.

Use the repository's existing status conventions if they can represent these states accurately. Do not force new fields into the schema without considering backward compatibility and the cost of maintaining them.

Do not mark a lesson complete merely because it was edited.

When a batch exposes a curriculum-wide defect, update the implementation plan and fix the shared cause instead of applying inconsistent local workarounds.

## 9. Decision-making policy

Proceed autonomously on routine, reversible decisions that are supported by the repository and the curriculum's stated goals.

When a structural change could affect many IDs, modules, projects, assessments, or the learning progression:

1. Document the current behavior.
2. Explain the defect.
3. Propose the smallest defensible correction.
4. Assess affected records and compatibility.
5. Implement only when the correct choice is clear.
6. Record the decision and validate its consequences.

Do not repeatedly ask for approval on ordinary wording, clear placeholder replacements, or straightforward lesson improvements.

Do not invent user preferences or unsupported educational requirements to fill gaps. If a consequential choice cannot be resolved from repository evidence, record the decision needed and its alternatives in the implementation plan, then continue with independent work.

## 10. Required execution order for this assignment

Begin now.

### First execution batch

1. Verify the repository state and record the baseline.
2. Inspect the authoritative manifests, all module directories, actual physical lesson files, curriculum rules, and validation tooling.
3. Produce the Phase 0 inventory and prioritized audit.
4. Identify and resolve the highest-impact Phase 1 curriculum architecture decisions.
5. Update the persistent implementation plan.
6. Run the available baseline validation and record actual results.

Do not start mass content generation before this audit.

### Second execution batch

After the audit, finalize the curriculum blueprint and begin Phase 2 metadata remediation in manageable batches. Fix the highest-impact structural and semantic defects first. Keep all manifests consistent as you work.

### Subsequent execution batches

Complete metadata across the intended curriculum, then implement real content one module at a time in prerequisite order. Validate each batch, update progress, and continue from the persistent plan.

Do not stop after producing a plan if implementation can safely proceed. Do not attempt to finish the entire curriculum in one unreviewable response or commit.

If the session ends, leave the repository and implementation plan in a state where the next session can continue without repeating work.

## 11. Definition of success

This assignment is successful when the repository has progressed from scaffold-oriented metadata and incomplete instructional artifacts toward a coherent, genuinely authored, verifiable curriculum.

Success is not:

- A larger number of lessons.
- Nonempty JSON fields.
- Passing schema validation alone.
- Generic text added to every README.
- A green validator achieved by weakening requirements.
- Inflated completion statistics.
- A release label unsupported by actual implementation.

Success means:

- Every planned item has a clear educational purpose.
- Metadata accurately describes the intended lesson.
- Physical content fulfills its contract.
- Prerequisites and relationships form a coherent learning graph.
- Examples, exercises, projects, and assessments reinforce the objectives.
- Technical behavior can be verified where applicable.
- Quality checks catch real defects.
- Curriculum-wide claims match the evidence.
- Future contributors can maintain the system reliably.

**Start with Phase 0 now. Inspect the actual repository, produce the inventory, and establish the implementation plan before rewriting the curriculum. Then proceed through the phases in order, preserving user work and reporting verified results.**
