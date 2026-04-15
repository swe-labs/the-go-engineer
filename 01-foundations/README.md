# 01 Foundations

## Mission

This folder is the true starting point for the new curriculum architecture.

Its job is not to rush learners into backend patterns or clever abstractions.
Its job is to build the first reliable mental model of:

- how Go programs run
- how decisions and repetition work
- how data is stored and moved
- how logic becomes reusable and failure-aware

## Internal Order

The foundations layer is intentionally ordered like this:

1. `01-getting-started`
2. `02-language-basics`
3. `03-control-flow`
4. `04-data-structures`
5. `05-functions-and-errors`

That order matters.
Each section should prepare the learner for the next one without hiding key ideas behind future
concepts.

## README-First Rule

Every learner-facing lesson in Foundations should work like this:

1. read the lesson `README.md`
2. open `main.go`
3. run the code
4. change the code yourself
5. attempt the starter or milestone when one exists

The README is the primary teaching surface.
The code is still mandatory, but it should stay clean enough that the learner can focus on what the
README already prepared them to see.

## Zero-Magic Rule

Foundations must stay honest.

That means:

- no helper-function design before functions are taught
- no heavy struct-based modeling before types and design
- no package architecture before engineering core
- no production concerns smuggled into beginner lessons

Small previews are allowed when they help the learner move forward, but they must be labeled
clearly and explained later in the proper section.

## Current Priority

The current rebuild line is moving backward to finish the front of Foundations:

- `01-getting-started`
- `02-language-basics`

That work matters because the later canonical sections are already in place:

- `03-control-flow`
- `04-data-structures`
- `05-functions-and-errors`

## Next Step

Once `01-getting-started` and `02-language-basics` are rebuilt, the whole foundations layer will
share one learner contract from the very first lesson onward:

- README first
- clean code second
- zero-magic sequencing
- one obvious next step at the end of every lesson
