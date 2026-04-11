# Workflow Confidence Notes

This page is for the moment when setup works "almost" correctly but you do not trust it yet.

## The First Useful Workflow

You do not need to master the whole toolchain yet.
You only need this loop:

1. open the lesson folder
2. run the lesson with `go run`
3. read the output
4. make one small change
5. run it again

That loop is enough to begin.

## If You Feel Lost

Use this recovery sequence:

1. go back to the repo root
2. run `go version`
3. rerun the hello-world lesson
4. reread the matching stage or lesson README

Returning to one known-good lesson is not failure.
It is a normal debugging move.

## Common Beginner Mistakes

### Mistake 1: Running commands from the wrong folder

Fix:

- return to the repo root
- run the command exactly as written in the lesson

### Mistake 2: Treating every error like a disaster

Fix:

- read the first useful line of the error
- check the command path
- check whether Go is installed correctly

### Mistake 3: Jumping ahead before the workflow feels stable

Fix:

- finish `GT.4 development environment`
- rerun `GT.2 hello world`
- only then move into `1 Language Fundamentals`

## Good Signs

You are gaining workflow confidence when:

- you can rerun a lesson without hesitation
- you know how to return to the repo root
- you can tell the difference between a path mistake and a code mistake
- the terminal feels like a tool, not a threat
