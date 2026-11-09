# Rubric

This rubric scores the Course Setup Verification assessment. It matches the three
criteria defined by the project itself (`projects/course-setup-verification/README.md`).

| Dimension | Weight | Pass criteria |
|---|---|---|
| Toolchain verification | 40% | Terminal, Git, and editor each confirmed with a real command output pasted into the answers; exit-code meaning explained correctly. |
| Git identity configured | 30% | `git config user.name` and `git config user.email` both return non-empty values; the configuration command is stated. |
| Workspace established | 30% | The dedicated workspace directory was created outside the curriculum repository; `setup.txt` records Git version, committer name, and email; the workspace rationale (Lesson 01) is explained in the learner's own words. |

## Scoring guide

- Score each dimension 0-5, then weight: Toolchain x0.40, Identity x0.30, Workspace x0.30.
- **Passing standard: 80% overall with all three dimensions satisfied** (no dimension at
  zero). A missing `setup.txt` or an unconfigured Git identity fails the assessment
  regardless of the numeric score.

## Reflection questions (9-12) are graded inside the dimensions

- Questions 9-10 score under Toolchain verification (conceptual grounding).
- Questions 5-6, 8 score under Workspace established.
- Questions 11-12 score under Git identity configured (applies the Lesson 07 debugging
  protocol to a real failure).

**Self-grade before submitting.** Compare your answers with `answer-key.md`.
