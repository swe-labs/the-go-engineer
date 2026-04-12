# V2 Lesson README Template

## Purpose

This document defines the canonical learner-facing README shape for beta lessons.

It exists because learner-facing lessons should not force the source file to carry all of the
explanation. The README is where we answer the learner's "what", "why", and "how" questions
without making `main.go` louder than the lesson itself.

## When To Use This Template

Use this template for:

- learner-facing lessons
- any lesson where learners benefit from line-by-line or very small chunk explanation

The explanation depth can change by stage, but the overall division of labor should stay the same:

- code file = runnable implementation
- README = teaching walkthrough

## Canonical README Shape

Every beta lesson README should usually include:

1. mission
2. prerequisites
3. mental model
4. why this lesson exists now
5. run instructions
6. code walkthrough
7. common questions or mistakes
8. production relevance
9. next step

## Canonical README Skeleton

~~~md
# Lesson Title

## Mission

One short paragraph explaining what the learner is about to understand and why it matters here.

## Prerequisites

- the prior lesson or concepts the learner should already know

## Mental Model

Name the framing idea the learner should hold while reading the code.

## Why This Lesson Exists Now

Explain why this lesson appears at this point in the section or stage and what later confusion it
prevents.

## Run Instructions

~~~bash
go run ./NN-section-name/N-lesson-name
~~~

Optional:

~~~bash
go test ./NN-section-name/N-lesson-name
~~~

## Code Walkthrough

Walk through the code line by line or in small logical chunks.

Recommended rule:

- explain each line when the step matters on its own
- only group lines when they form one inseparable statement or one tiny setup step

Useful subheadings:

### Imports And Setup

What each import or setup line is doing.

### Main Flow

Walk the learner through the example in order.

### Why This Update, Check, Loop, Or Branch Exists

Explain the non-obvious behavior and what would happen if it changed.

## Common Questions

- one likely beginner confusion
- one runtime-model confusion
- one "what if I changed this line?" answer

## Production Relevance

State where this idea appears in real Go work and what bad habit or failure mode it helps prevent.

## Next Step

Point to the next lesson, drill, exercise, or checkpoint.
~~~

## Authoring Rules

- keep the code walkthrough concrete
- answer "what", "why", and "how", not just "what"
- prefer small chunks over giant prose walls
- keep the README honest to the code that actually exists
- do not turn the README into a second unrelated textbook chapter
- do not skip the code after the docs; the code is still the required runnable proof surface

## Success Signal

This template is working when a complete beginner can:

- read the code without guessing what each part is doing
- connect the code to the section mental model
- tell why the lesson exists at this point in the curriculum
- move to the next item without feeling like the example contained hidden magic
