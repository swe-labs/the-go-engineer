# README Quality Rubric v1

## Scoring Scale

Each dimension is scored 1–5:

| Score | Meaning                                          |
| ----- | ------------------------------------------------ |
| 1     | Missing or misleading                            |
| 2     | Present but superficial                          |
| 3     | Adequate — meets minimum standard                |
| 4     | Strong — clear, accurate, useful                 |
| 5     | Exemplary — sets the standard for the curriculum |

---

## Dimensions

### 1. Clarity

Can a motivated learner understand the lesson without external help?

**4+ criteria:**

- No unexplained jargon
- Concepts introduced in dependency order
- Each paragraph answers an implicit "why does this matter?"
- Code examples are minimal and focused on one idea

**2–3 signals:**

- Learner must search the web to understand terminology
- Paragraphs are dense walls of text without structure

---

### 2. Operational Realism

Does the lesson reflect real production behavior?

**4+ criteria:**

- Examples use realistic data (not foo/bar/baz)
- Failure scenarios are production failures, not toy errors
- Code patterns match idiomatic Go standard library usage

**2–3 signals:**

- Examples use throwaway variables (x, y, z)
- Error handling is omitted or unrealistic
- Production tradeoffs are not mentioned

---

### 3. Magic Elimination

Does the learner explain mechanics, not just usage?

**4+ criteria:**

- Every function or concept has an under_the_hood explanation
- The learner can predict behavior before running code
- Hidden_magic_checks require mechanical reasoning

**2–3 signals:**

- Lesson explains what but not how
- Learner can complete exercises without understanding runtime behavior
- No memory model or execution model is presented

---

### 4. Debuggability

Does the lesson teach how to diagnose failures?

**4+ criteria:**

- At least one debugging walkthrough is included
- Common failure modes list concrete symptoms and root causes
- Learner practices reading error messages, stack traces, and log output

**2–3 signals:**

- Debugging is not mentioned
- Only happy-path code is shown
- Error messages are not explained

---

### 5. Production Relevance

Does the content apply to real engineering work?

**4+ criteria:**

- Production examples show real code patterns (not toy examples)
- Performance implications are discussed
- Operational considerations are addressed (deployment, monitoring, rollback)

**2–3 signals:**

- Examples are academic or artificial
- No discussion of real-world impact
- No production code patterns shown

---

### 6. Mental-Model Quality

Is the analogy accurate and durable?

**4+ criteria:**

- Mental model maps cleanly to the runtime behavior
- Model is internally consistent (no contradictions)
- Model survives edge cases (does not break under unusual conditions)

**2–3 signals:**

- Analogy is too vague to be useful ("think of it as a box")
- Analogy breaks under common edge cases
- No mental model is provided

---

### 7. Cognitive Pacing

Is the complexity appropriate for the module's position in the curriculum?

**4+ criteria:**

- Prerequisites are respected (no concepts introduced before they are taught)
- One new idea per section
- Lesson has clear stopping points (breaks, summaries)

**2–3 signals:**

- Lesson assumes knowledge not yet taught
- Multiple hard concepts introduced simultaneously
- No obvious structure or section breaks

---

### 8. Diagram Usefulness

Do diagrams clarify rather than decorate?

**4+ criteria:**

- Diagram reveals something the text cannot (memory layout, timeline, state machine)
- Diagram has a caption explaining what to observe
- Diagram is referenced from the text at the point it is needed

**2–3 signals:**

- Diagram is decorative (adds no new information)
- Diagram is not referenced from the text
- Diagram is placed before the concept is introduced

---

### 9. Code-Reading Quality

Does the lesson teach code reading, not just writing?

**4+ criteria:**

- Learner reads existing code and explains it
- Code-reading tasks precede writing tasks
- Lesson includes questions about what a given code snippet does (not just how to write it)

**2–3 signals:**

- Only writing tasks
- No code-reading exercises
- All code is provided by the lesson; none is read critically

---

### 10. Failure Realism

Are the failure scenarios credible and instructive?

**4+ criteria:**

- Failure scenarios are based on real production incidents
- Each failure describes: trigger → symptom → root cause → detection
- Lesson includes a failure exercise (find the bug, explain the crash)

**2–3 signals:**

- Failure scenarios are contrived or impossible in practice
- Only one failure mode is shown
- No detection strategy is taught

---

## Minimum Passing Score

| Phase           | Minimum per dimension | Overall |
| --------------- | --------------------- | ------- |
| Authoring draft | 2                     | 25/50   |
| Review ready    | 3                     | 35/50   |
| Golden standard | 4                     | 45/50   |

---

## Review Process

1. Self-review: author scores each dimension before submitting
2. Peer review: reviewer scores each dimension independently
3. Resolution: discrepancies >1 point are discussed and resolved
4. Final: maintainer confirms rubric scores and approves or requests revisions
