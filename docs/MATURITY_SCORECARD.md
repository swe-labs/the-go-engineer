# The Go Engineer Maturity Scorecard

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

- [ ] Docker build validation passes
- [ ] health/readiness/liveness semantics documented
- [ ] migration policy and checksums enforced
- [ ] observability export path tested

## Education

- [ ] learner checkpoints documented
- [ ] lesson feedback path visible
- [ ] proof surfaces executable from docs

## Consistency

- [ ] curriculum validator passes
- [ ] `go mod tidy` produces no diff
- [ ] formatting gate clean
- [ ] OpenAPI surface aligns with handler routes

## Community

- [ ] issue templates and PR template available
- [ ] code of conduct and contributing guide visible
- [ ] label taxonomy documented
- [ ] CODEOWNERS in place
