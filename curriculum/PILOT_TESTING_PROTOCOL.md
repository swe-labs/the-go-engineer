# The Go Engineer v3 — Beginner Pilot Testing Protocol

> **Purpose:** Formal empirical protocol to evaluate unassisted learner velocity, identify instructional friction points, and audit zero-magic integrity across the foundational phase (**Modules 00 through 05**).  
> **Status:** Active Quality Assurance Protocol  
> **Target Outcome:** Validate that a non-Go learner can progress from installation to writing idiomatically sound, type-safe, tested Go software without human intervention.

---

## 1. Pilot Scope & Focus Modules

The pilot tests the foundational and engineering core tracks where cognitive friction is highest:

- **Module 00: Orientation** (Repository navigation, learning model, engineering standards)
- **Module 01: Computers, Terminal, Git, and the Web** (Command-line environment, Git workflows, process models)
- **Module 02: Go Setup and Tooling** (Go toolchain installation, `go build`, `go vet`, `gofmt`, module cache)
- **Module 03: Programming Fundamentals** (Values, types, control flow, slices, maps, pointers, algorithmic problem-solving)
- **Module 04: Functions, Errors, and Data Semantics** (Named returns, error handling as values, defer mechanics, mutation semantics)
- **Module 05: Types, Interfaces, Packages, and Modules** (Structs, method sets, interfaces, generics fundamentals, documentation)

---

## 2. Participant Cohort Profiles

To ensure balanced validation, test subjects must be recruited across two distinct cohorts:

| Cohort | Profile Description | Target Sample Size | Evaluation Objective |
| :--- | :--- | :---: | :--- |
| **Cohort A (Novice Learners)** | Individuals with basic computer literacy but zero prior programming or terminal experience. | 5 participants | Test foundational terminology, CLI setup friction, and zero-magic clarity. |
| **Cohort B (Transitioning Engineers)** | Developers with 1–3 years experience in dynamically typed languages (Python, JavaScript, Ruby, PHP). | 5 participants | Test mental model shifts: static typing, explicit error values, pointer semantics, and goroutine concepts. |

---

## 3. The "Zero-Intervention" Observation Methodology

1. **Unassisted Environment:** Participants work on a standard workstation with only the repository, official Go documentation, and standard terminal tools available.
2. **Think-Aloud Protocol:** Participants are encouraged to verbalize their thoughts, questions, and confusion as they read lessons, inspect compiler output, and write code.
3. **Observer Non-Interference Rule:** Observers **must not** intervene, correct syntax mistakes, or provide hints—even if the participant struggles for extended periods.
4. **Stuck Threshold (Intervention Protocol):** If a participant remains completely blocked on a single step for **> 20 minutes**, the observer documents the exact blockage, logs a **Severity 1 Friction Bug**, provides the minimum hint necessary to proceed, and resumes non-intervention.

---

## 4. Empirical Friction Log Metrics

Observers log every anomaly using the following structured metrics:

```
[PILOT-FRICTION-LOG]
Session ID:       PILOT-2026-COHORT-A-01
Participant:      P-04
Module & Lesson:  core-03-22 (Pointers as addresses)
Timestamp:        T+01:14:20
Friction Type:    Mental Model Misconception
Observed Action:  Participant wrote `p := &42` expecting to take the address of an integer literal.
Diagnostic:       Compiler message "cannot take the address of 42" did not explain *why* literals lack storage locations.
Time to Resolve:  8 minutes (Participant re-read section on memory storage).
Recommendation:   Add explicit callout box in core-03-22 contrasting addressable variables vs value literals.
```

### Core Benchmark Metrics
1. **Time-to-First-Run (TTFR):** Minutes elapsed from opening Module 00 to successfully running `main.go` in Module 02. (Benchmark: $\le 45\text{ minutes}$).
2. **Compiler Error Diagnosis Time (CEDT):** Time taken to interpret a `go build` or `go test` failure and apply a working fix. (Benchmark: $\le 5\text{ minutes}$ per error).
3. **Pointer Semantic Failure Frequency:** Count of unintentional value copy mutations vs pointer receiver errors. (Benchmark: $\le 2$ occurrences across Module 04/05).
4. **Zero-Magic Violations Detected:** Instances where a participant encounters undefined terminology or syntax not previously introduced. (Benchmark: **0 violations tolerated**).

---

## 5. Evaluation & Quality Gate Criteria

Before foundational modules are approved for release, the pilot cohort results must satisfy the following acceptance criteria:

| Gate | Criterion | Minimum Threshold |
| :--- | :--- | :---: |
| **Gate 1: Autonomous Completion** | Participants completing Modules 00–05 without exceeding the 20-minute stuck threshold. | $\ge 90\%$ |
| **Gate 2: Project Assessment Pass Rate** | First-attempt pass rate on `project-pricing-checkout` and `project-bank-account` unit tests. | $\ge 80\%$ |
| **Gate 3: Conceptual Defense** | Oral explanation of pointer mechanics, error handling, and interfaces using required terminology. | $\ge 85\%$ rubric score |
| **Gate 4: Zero-Magic Compliance** | Zero undefined concepts or unexplained standard library magic reported in friction logs. | $100\%$ compliance |

---

## 6. Feedback Integration Lifecycle

1. **Daily Triage:** Review friction logs at the end of each pilot testing day.
2. **Documentation Patching:** If two or more participants trigger the same misunderstanding, patch the corresponding lesson README or zero-magic execution timeline within 24 hours.
3. **Re-Test Verification:** Present the updated lesson to the next participant in the cohort to confirm resolution of the friction point.
