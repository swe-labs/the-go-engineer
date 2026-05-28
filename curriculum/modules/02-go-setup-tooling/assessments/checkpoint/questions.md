# Module 02 Checkpoint Questions

## Part 1 — Installation and environment

1. What does `go version` prove?
2. What does `command -v go` or `which go` prove?
3. What is `GOROOT`?
4. What is `GOPATH`?
5. What are `GOMODCACHE` and `GOCACHE` used for?

## Part 2 — Running and building

6. What are the four important parts of a Hello World program?
7. What does `go run` do before executing the program?
8. Why can `go run` fail before your program starts?
9. What does `go build` produce for a command package?
10. Why might you use `go build -o`?

## Part 3 — Testing and formatting

11. What does `go test` compile?
12. What makes a useful test failure message?
13. Why is `gofmt` not optional in Go?
14. What does `gofmt -l` show?

## Part 4 — Inspection tools

15. What kind of problem can `go vet` catch?
16. Why does a clean `go vet` result not prove program correctness?
17. What does `go doc` read?
18. Give one useful `go doc` command.

## Part 5 — Editor tooling

19. What does `gopls` do?
20. Why should terminal verification still matter when the editor looks green?

## Part 6 — Reading failures

21. How should you read a compiler error?
22. How is a runtime error different from a compiler error?
23. What should you identify in a test failure?
24. Why should you fix the first compiler error first?

## Part 7 — Module root

25. What is `go.mod`?
26. What does `go env GOMOD` show?
27. Why do Go commands behave differently outside a module?
28. What is one symptom of module root confusion?

## Part 8 — Checklist

29. Write your tooling checklist for starting work on a new machine or branch.
30. Which command in the checklist feels least familiar, and how will you practice it?
