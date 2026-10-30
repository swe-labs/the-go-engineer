# Module 00 — Orientation

## Mission

This module teaches you how to use the curriculum itself.

Before you write serious Go code, you need to know how the repository is organized, what "zero magic" means, how lessons and labs work, how assessments judge progress, and what kind of portfolio evidence you will build.

The goal is not to memorize folder names. The goal is to understand the learning system so every later module feels navigable and fair.

## Tracks in this curriculum

The curriculum organizes your engineering journey into three coordinated tracks:

1. **Core Backend Path (`curriculum/modules/`):** The primary 18-module progression from computing fundamentals and systems programming to concurrent worker pools, distributed architectures, and Flagship Opslane.
2. **Advanced Electives (`curriculum/electives/`):** Deep dives into gRPC streaming, complex generics constraints, CQRS, and advanced distributed patterns.
3. **[Frontend Companion Track](../../companion-frontend/README.md) (`curriculum/companion-frontend/`):** A 4-stage progressive server-driven hypermedia track (`templ` + `HTMX` + `Tailwind CSS` + `Alpine.js`, with `Datastar` comparison) running parallel to Modules 07–15 and converging at Flagship Opslane.

## Who this module is for

This module is for:

- complete beginners who have never used a programming curriculum before
- self-taught learners who want a clear path instead of random tutorials
- experienced developers who want to understand the repository contract
- contributors who need to know how learner-facing content is structured

## What you will build

You will build orientation artifacts:

- a repository navigation checklist
- a zero-magic explanation
- a lesson-cycle map
- a local run/test habit
- a debugging question template
- a personal roadmap
- a portfolio plan

These artifacts are small, but they matter. They become your operating system for the rest of the curriculum.

## Prerequisites

None.

You should be able to:

- open files in a text editor
- read Markdown
- use basic terminal commands such as `cd`, `ls`, and `pwd`

If those commands are unfamiliar, continue anyway. Module 01 and Module 02 will build the terminal and tooling foundation more deeply.

## Concept map

```text
Repository map
    ↓
Zero-magic contract
    ↓
Lesson/lab/project/checkpoint loop
    ↓
Run/test/verify habit
    ↓
Starter/solution workflow
    ↓
Assessment workflow
    ↓
Debugging question workflow
    ↓
Job-ready definition
    ↓
Roadmap navigation
    ↓
Portfolio plan
```

## Lessons

| Order | Type   | Path                                                           | Outcome                                                      |
| ----: | ------ | -------------------------------------------------------------- | ------------------------------------------------------------ |
|     1 | Lesson | `lessons/01-how-to-use-this-repository/`                       | Navigate metadata and learner-facing files.                  |
|     2 | Lesson | `lessons/02-what-zero-magic-means/`                            | Explain the no-hidden-assumptions teaching contract.         |
|     3 | Lesson | `lessons/03-how-lessons-exercises-projects-and-checkpoints-work/` | Practice the full learning loop.                             |
|     4 | Lesson | `lessons/04-how-to-run-code-in-this-repository/`               | Run examples and tests confidently.                          |
|     5 | Lesson | `lessons/05-how-starter-folders-work/`                         | Use starter and solution folders correctly.                  |
|     6 | Lesson | `lessons/06-how-assessments-work/`                             | Understand checkpoints, rubrics, and retakes.                |
|     7 | Lesson | `lessons/07-how-to-ask-good-debugging-questions/`              | Turn confusion into useful debugging evidence.               |
|     8 | Lesson | `lessons/08-what-job-ready-means/`                             | Define job readiness through evidence, not confidence alone. |
|     9 | Lesson | `lessons/09-how-to-use-the-roadmap/`                           | Plan the journey through modules and prerequisites.          |
|    10 | Lesson | `lessons/10-how-to-build-a-portfolio-from-this-curriculum/`    | Convert curriculum projects into portfolio proof.            |

## Labs

Module 00 includes one guided lab:

```text
lessons/03-how-lessons-exercises-projects-and-checkpoints-work/
```

The lab asks you to trace a learning cycle from lesson to proof.

## Projects

Module 00 includes one foundation project:

- [`projects/course-setup-verification/`](./projects/course-setup-verification/README.md) — Course Setup Verification: verify your Go environment, Git installation, IDE configuration, and local build verification.

While full-scale software portfolio systems begin in Module 01, this first operational project establishes your personal proof repository and environment baseline.

## Assessments

Complete the two Module 00 assessments:

1. [`assessments/course-setup/`](./assessments/course-setup/README.md) — Course Setup Verification Assessment (evaluates `project-course-setup`).
2. [`assessments/checkpoint/`](./assessments/checkpoint/README.md) — Module 00 Checkpoint (verifies curriculum navigation, zero-magic contracts, and learning cycle readiness).

## Common failure modes

| Failure                        | Symptom                                                         | Fix                                                        |
| ------------------------------ | --------------------------------------------------------------- | ---------------------------------------------------------- |
| Treating reading as completion | You can quote text but cannot run, change, or explain anything. | Use the proof task at the end of every lesson.             |
| Skipping orientation           | Later folders and validators feel mysterious.                   | Finish this module before jumping ahead.                   |
| Copying commands blindly       | A command fails and you do not know what it was meant to prove. | Read the Run Instructions and Machine View before running. |
| Looking at solutions too early | You finish faster but learn less.                               | Attempt the starter first, then compare.                   |
| Avoiding assessments           | Weak spots remain invisible.                                    | Use assessments as feedback, not punishment.               |

## Completion checklist

You are ready for Module 01 when you can:

- [ ] explain the difference between `curriculum/`, `internal/tools/`, `web/`, and `.agents/`
- [ ] find a lesson README from its metadata path
- [ ] explain zero magic in your own words
- [ ] run at least one lesson example
- [ ] run at least one test command
- [ ] use a starter folder without editing the solution first
- [ ] answer the checkpoint questions
- [ ] write a short personal roadmap and portfolio plan

## Next module

After this module, continue to:

```text
curriculum/modules/01-computers-terminal-git-web/
```

That module teaches the computing, terminal, Git, and web foundations you will use throughout the rest of the curriculum.
