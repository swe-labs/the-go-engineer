# Lesson README template

Use this structure unless the repository has a stricter contract.

```markdown
# [Lesson Title]

## Learning objective
[One concrete outcome. Avoid “understand X” alone; say what the learner can do.]

## Why this matters
[Real problem solved. Connect to professional Go/software engineering.]

## Mental model
[Analogy + precise model. State where the analogy stops being true.]

## Core idea
[Teach the concept plainly. Introduce terms before using them.]

## Under the hood
[Runtime/compiler/OS/database/network mechanics.]

## How Go uses it
[How the language, standard library, runtime, or common Go practice expresses the concept.]

## Go example
```go
// minimal runnable example
```

## Step-by-step execution
1. [What happens first]
2. [What happens next]
3. [What state changes]

## Common mistakes
- Mistake: [specific beginner mistake]
  - Why it happens: [...]
  - Fix: [...]

## Debugging walkthrough
[Give a broken example, symptom, investigation, root cause, fix.]

## Production notes
[How this appears in real teams/services. Include tradeoffs and operational risks.]

## Performance implications
[State performance, memory, reliability, or operational tradeoffs.]

## Practice task
[Small task tied to proof_of_understanding.]

## Tests / verification
```bash
[commands]
```

## Review questions
1. [Explain]
2. [Apply]
3. [Debug]
4. [Tradeoff]

## NEXT UP
[Next lesson title or ID]
```

## Folder expectations

Place lessons in the final typed folder layout:

```text
curriculum/modules/{module}/lessons/{lesson}/README.md
curriculum/modules/{module}/lessons/{lesson}/main.go
curriculum/modules/{module}/lessons/{lesson}/main_test.go
curriculum/modules/{module}/lessons/{lesson}/_starter/
curriculum/modules/{module}/lessons/{lesson}/_solution/
curriculum/modules/{module}/lessons/{lesson}/assets/
```

For labs, projects, and assessments, use the corresponding typed folder: `labs/`, `projects/`, or `assessments/`.

## Quality bar

The README should be useful even without the instructor present. If a motivated beginner cannot follow it and a senior engineer cannot respect it, keep improving it.
