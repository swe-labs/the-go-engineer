# Solution — How to use this repository

This solution shows one clear way to describe the top-level repository responsibilities.

Your wording does not need to match exactly. It should preserve the same boundaries:

- metadata describes the curriculum graph
- curriculum contains learner-facing material
- tools validate, generate, audit, and migrate
- docs explain maintainer policy
- dist contains generated release artifacts

## Run

```bash
go run ./curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/_solution
go test ./curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/_solution
```
