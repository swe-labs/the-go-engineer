# Go Engineer Curriculum

Welcome to the learner-facing part of Go Engineer.

This directory contains the actual curriculum: module guides, lesson explanations, runnable Go examples, tests, labs, projects, assessments, diagrams, and shared references. Metadata in `metadata/` defines the learning graph; this directory is where the learner experiences that graph.

## How to use this directory

Start with the first module and move in order unless you already know the prerequisites.

```text
curriculum/
├── modules/     required core curriculum
├── electives/   optional advanced curriculum
└── shared/      reusable learner references
```

Each module follows the same shape:

```text
curriculum/modules/{module}/
├── README.md
├── lessons/
├── labs/
├── projects/
├── assessments/
└── assets/
```

That structure is intentional. Lessons teach, labs practice, projects prove integrated ability, and assessments check mastery.

## Learning contract

This curriculum uses a zero-magic learning system.

That means:

- every concept is explained before it is required
- every command is tied to a reason
- every exercise has an observable proof surface
- every project has a quality standard
- every assessment maps back to taught material
- every shortcut is either explained or removed

You are not expected to guess what invisible experience a senior engineer already has. The curriculum makes those assumptions visible.

## How to move through a module

For each module:

1. Read the module `README.md`.
2. Complete lessons in order.
3. Run the example code where provided.
4. Run tests where provided.
5. Complete labs without looking at the solution first.
6. Finish any module project.
7. Complete the checkpoint assessment.
8. Review mistakes before moving on.

A lesson is complete when you can explain it, run it, modify it, and recognize common failure modes.

## Required local tools

You will eventually need:

- Go
- a terminal
- Git
- a code editor
- a browser
- Docker and PostgreSQL in later modules

Module 00 does not assume you already know these tools. It teaches how to use the repository safely before technical depth begins.

## Where to begin

Start here:

```text
curriculum/modules/00-orientation/README.md
```

Then continue to:

```text
curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/README.md
```

## Proof mindset

Do not treat reading as completion.

For every lesson, ask:

- What did I learn?
- What command or artifact proves it?
- What mistake would show I do not understand it yet?
- How would I explain this to another beginner?
- How will this matter in a professional Go codebase?

That habit is the foundation of the whole curriculum.
