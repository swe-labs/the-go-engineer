# The Go Engineer Maturity Scorecard

> **Note:** This scorecard is a living maturity checklist. Checked items represent currently verified surfaces; unchecked items represent deliberate next hardening targets.

## Correctness

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go test -race ./...`
- [ ] fuzz suites pass for critical parser/middleware paths

## Security

- [ ] `govulncheck ./...`
- [ ] CodeQL workflow enabled
- [ ] Dependabot enabled
- [ ] `SECURITY.md` published
- [ ] Opslane threat model documented

## Documentation

- [ ] README and curriculum registry aligned
- [ ] ADR set present and current
- [ ] known limitations and glossary maintained
- [ ] learning path and feedback loop documented

## Operations

- [x] Docker build validation passes
- [x] health/readiness/liveness semantics documented
- [x] migration policy and checksums enforced
- [x] observability export path tested

## Education

- [ ] learner checkpoints documented
- [x] lesson feedback path visible
- [ ] proof surfaces executable from docs

## Consistency

- [x] curriculum validator passes
- [x] `go mod tidy` produces no diff
- [x] formatting gate clean
- [x] OpenAPI surface aligns with handler routes

## Community

- [x] issue templates and PR template available
- [x] code of conduct and contributing guide visible
- [x] label taxonomy documented
- [x] CODEOWNERS in place
