# Module 00 Checkpoint Answer Key

Use this after attempting the checkpoint.

## Part 1

1. `curriculum/` houses both authoritative manifest blueprints (`path.core.json`, `projects.json`, etc.) and the physical instructional content (`modules/`, `electives/`, and companion tracks).
2. `internal/tools/curriculum/` serves as the automated quality compiler, running 17 quality gates covering schemas, dependency cycles, zero-magic compliance, and AST parity.
3. `web/` hosts the learner progression portal and graph visualizer, while `.agents/` stores domain skills and governance rubrics.
4. Manifest blueprints declare the learning graph; strict validation against `schema.v3.json` guarantees zero dangling references, no cyclic dependencies, and complete filesystem parity.
5. Module: `08-http-rest-apis`. Lesson: `03-handler-lifecycle`.

## Part 2

6. Zero magic means no required assumption is hidden. Concepts, commands, code behavior, failures, and proof are explained explicitly.
7. Example: telling learners to run `go test` without explaining what a test file is, what `Test...` means, or what passing proves.
8. It should give a short local explanation and name where the concept is taught fully later.

## Part 3

9. A lesson teaches. A lab gives guided practice. A project proves integrated ability. An assessment checks mastery.
10. Reading does not prove you can run, modify, debug, or explain the idea.
11. After making a serious attempt with the starter and recording what you tried.
12. Return to the relevant lesson, fill the gap, revise the answer, and retake that part.

## Part 4

13. `go run` compiles and executes a package. `go test` compiles test files and runs test functions that verify behavior.
14. A good debugging question includes expected behavior, actual behavior, exact command, exact error, environment, recent changes, and smallest reproduction.
15. It usually means the command was run from the wrong directory or the path is wrong.

## Part 5

16. Later modules depend on earlier skills, so order is based on prerequisites, not preference.
17. Strong evidence includes tests, clear README, architecture diagram, deployment notes, tradeoff explanation, screenshots/logs, and clean commits.
18. Good answers name concrete schedule, proof tasks, and how the learner will handle weak spots.
