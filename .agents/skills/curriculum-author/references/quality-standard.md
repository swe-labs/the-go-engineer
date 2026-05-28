# World-class curriculum quality standard

## Zero-magic standard

A learner should never be asked to accept behavior as magic. Every lesson must explain:

- the problem the concept solves
- why the concept exists
- a concrete analogy or mental model
- where the analogy stops being true
- what happens under the hood
- how Go expresses the concept
- how the concept fails in real systems
- how to debug mistakes
- how to prove understanding
- how the concept appears in production

## Beginner-to-expert progression

Write for a true beginner without diluting professional depth:

1. Start from observable behavior.
2. Introduce one new concept at a time.
3. Show a small example.
4. Trace the mechanism step-by-step.
5. Break it intentionally and debug the failure.
6. Apply it in a realistic task.
7. Connect it to production practice.
8. Ask review questions that require explanation, application, debugging, and tradeoff reasoning.

## Explanation requirements

Use this pattern for difficult ideas:

- **Plain-language definition**: one paragraph.
- **Analogy**: concrete, not cute; it must preserve the important constraints.
- **Boundary of analogy**: state what the analogy does not cover.
- **Mechanical model**: what the runtime/compiler/OS/database/network actually does.
- **Code model**: minimal Go example.
- **Failure model**: what goes wrong and why.
- **Debugging model**: symptom → hypothesis → investigation → root cause → fix.
- **Production model**: how teams use it in real services.

## Repository quality requirements

Generated content must follow the final learner-facing layout:

- module overviews under `curriculum/modules/{module}/README.md`
- lessons under `curriculum/modules/{module}/lessons/{lesson}/`
- labs under `curriculum/modules/{module}/labs/{lab}/`
- projects under `curriculum/modules/{module}/projects/{project}/`
- assessments under `curriculum/modules/{module}/assessments/{assessment}/`
- electives under `curriculum/electives/{elective}/...`
- tools under `tools/`, not `internal/tools/` unless the repository explicitly keeps internal Go packages
- authoring procedures under `tools/authoring/`, not tool-branded folder names

## Code quality requirements

Generated code must be idiomatic Go:

- small functions with clear names
- explicit error handling
- context where cancellation/timeouts matter
- table-driven tests where suitable
- deterministic tests with no external network dependency unless explicitly required
- `gofmt` clean
- no global mutable state unless the lesson teaches why it is dangerous
- comments explain intent, not every line
- examples should be minimal enough to learn from, but realistic enough to transfer to production

## Project quality requirements

A project is complete only if it has:

- realistic problem statement
- milestones
- deliverables
- verification commands/manual checks
- rubric with weighted criteria
- failure cases
- documentation expectations
- portfolio narrative
- extension ideas for stronger learners
- explicit links to reinforced lessons and concepts

## Assessment quality requirements

An assessment is complete only if it tests:

- conceptual explanation
- code reading
- debugging
- implementation
- tradeoff reasoning
- production/failure awareness

Avoid recall-only assessments. Every module assessment should contain at least one applied task and one debugging/tradeoff question.
