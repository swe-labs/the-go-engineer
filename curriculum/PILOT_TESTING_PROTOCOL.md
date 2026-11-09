# The Go Engineer v3 — Beginner Pilot Testing Protocol

> **Purpose:** Formal empirical protocol to evaluate unassisted learner velocity, identify instructional friction points, and audit zero-magic integrity across the foundational phase (**Modules 00 through 05**).
> **Status:** Active Quality Assurance Protocol (Testing Framework)
> **Evaluation Scope:** Formative usability testing across novice learners and transitioning engineers. This document specifies the testing methodology, benchmarks, and friction triage lifecycle; empirical cohort results and revision logs are tracked separately per testing cycle.

---

## 1. Pilot Scope & Focus Modules

The pilot tests the foundational and engineering core tracks where cognitive friction and conceptual mental-model transitions are highest:

- **Module 00: Orientation** (Repository navigation, learning model, engineering standards)
- **Module 01: Computers, Terminal, Git, and the Web** (Command-line environment, Git workflows, process models)
- **Module 02: Go Setup and Tooling** (Go toolchain installation, `go build`, `go vet`, `gofmt`, module cache)
- **Module 03: Programming Fundamentals** (Values, types, control flow, slices, maps, pointers, algorithmic problem-solving)
- **Module 04: Functions, Errors, and Data Semantics** (Named returns, error handling as values, defer mechanics, mutation semantics)
- **Module 05: Types, Interfaces, Packages, and Modules** (Structs, method sets, interfaces, generics fundamentals, documentation)

---

## 2. Participant Cohort Profiles & Sample Size Bounds

To ensure balanced qualitative validation, test subjects are recruited across two distinct formative cohorts:

| Cohort | Profile Description | Formative Sample Size | Evaluation Objective |
| :--- | :--- | :---: | :--- |
| **Cohort A (Novice Learners)** | Individuals with basic computer literacy but zero prior programming or terminal experience. | 5 participants | Test foundational terminology, CLI setup friction, zero-magic clarity, and cognitive load pacing. |
| **Cohort B (Transitioning Engineers)** | Developers with 1–3 years experience in dynamically typed languages (Python, JavaScript, Ruby, PHP). | 5 participants | Test mental model shifts: static typing, explicit error values, pointer vs value semantics, and goroutine concepts. |

> **Sample Size & Methodology Note:** Cohorts of 5 participants per profile are designed for formative qualitative usability discovery (identifying high-impact instructional bottlenecks, misleading compiler explanations, and cognitive leaps). They serve as formative evaluation to drive targeted curriculum revisions, not as broad statistical proof of universal beginner success. Significant revisions trigger a follow-up validation cohort.

---

## 3. Observation Methodology & Ethical Rescue Protocol

1. **Unassisted Environment:** Participants work on a clean workstation with only the repository, official Go documentation, and standard terminal tools available.
2. **Think-Aloud Protocol:** Participants verbalize their thought process, confusion, and expectations as they read lessons, inspect compiler output, and write code.
3. **Observer Non-Interference:** Observers do not proactively correct syntax errors, provide leading hints, or interrupt problem-solving.
4. **Ethical Rescue Protocol (Stuck Threshold):** If a participant remains completely blocked on a single step for **> 20 minutes** despite re-reading curriculum material:
   - The observer logs a **Severity 1 Friction Bug** with exact reproduction steps.
   - The observer applies the **3-Tier Socratic Hint Ladder** (Tier 1 Nudge $\rightarrow$ Tier 2 Targeted Clue $\rightarrow$ Tier 3 Walkthrough), providing the minimal atomic guidance necessary to unblock progress without giving away complete solutions.
   - Non-interference resumes immediately.
   - *Rationale:* This ensures testing remains humane, instructive, and calibrated to observable engineering milestones rather than devolving into an endurance test.

---

## 4. Empirical Friction Log Metrics

Observers log every learning barrier using structured friction entries:

### Illustrative Example: Observer Friction Log Entry

```text
[PILOT-FRICTION-LOG: ILLUSTRATIVE EXAMPLE]
Session ID:       PILOT-2026-COHORT-A-01
Participant:      P-04
Module & Lesson:  core-03-22 (Pointers as addresses)
Timestamp:        T+01:14:20
Friction Type:    Mental Model Misconception
Observed Action:  Participant wrote `p := &42` expecting to take the address of an integer literal.
Diagnostic:       Compiler message "cannot take the address of 42" did not explain *why* literals lack storage locations.
Time to Resolve:  8 minutes (Participant re-read section on memory storage).
Intervention:     None (Resolved autonomously within 20-minute threshold).
Recommendation:   Add explicit callout box in core-03-22 contrasting addressable variables vs value literals.
Status:           Tracked for Curriculum Revision Batch 1.
```

### Core Benchmark Metrics (Targets)

1. **Time-to-First-Run (TTFR):** Minutes elapsed from opening Module 00 to successfully running `main.go` in Module 02. (*Target Benchmark:* $\le 45\text{ minutes}$).
2. **Compiler Error Diagnosis Time (CEDT):** Time taken to interpret a `go build` or `go test` failure and apply a working fix. (*Target Benchmark:* $\le 5\text{ minutes}$ per error).
3. **Module Completion Duration:** Elapsed time to complete all reading, exercises, and mini-projects per module.
4. **Pointer Semantic Failure Frequency:** Count of unintentional value copy mutations vs pointer receiver errors. (*Target Benchmark:* $\le 2$ occurrences across Module 04/05).
5. **Zero-Magic Violations Detected:** Instances where a participant encounters undefined terminology, unexplained syntax, or unmotivated standard library abstractions. (*Target Benchmark:* **0 violations tolerated**).

---

## 5. Evaluation & Quality Gate Criteria (Target Thresholds)

Before foundational modules are approved for release, cohort testing results must satisfy the following formative targets:

| Gate | Criterion | Target Threshold |
| :--- | :--- | :---: |
| **Gate 1: Autonomous Completion** | Participants completing Modules 00–05 without exceeding the 20-minute stuck threshold. | $\ge 90\%$ (at least 9 of 10 unassisted) |
| **Gate 2: Project Assessment Pass Rate** | First-attempt pass rate on `project-pricing-checkout` and `project-bank-account` unit tests. | $\ge 80\%$ |
| **Gate 3: Conceptual Defense** | Oral explanation of pointer mechanics, error handling, and interfaces using required zero-magic terminology. | $\ge 85\%$ rubric score |
| **Gate 4: Zero-Magic Compliance** | Zero undefined concepts or unexplained standard library magic reported in friction logs. | $100\%$ compliance |

---

## 6. Feedback Integration & Retesting Lifecycle

1. **Friction Triage:** Review observer logs at the conclusion of each testing session.
2. **Curriculum Patching:** When friction occurs across multiple participants, patch the corresponding lesson README, add visual memory diagrams, or clarify prerequisite concepts within 48 hours.
3. **Retest Verification:** Present revised material to a second validation cohort to confirm that the friction point has been eliminated before marking the lesson golden.
