# Security Policy

## Supported Versions

Security fixes are prioritized for:

- `main` (active integration line)
- `release/v2` (stable v2.1.x maintenance line)

## Reporting A Vulnerability

Please do not open public issues for suspected vulnerabilities.

Use one of these private channels:

- Security advisory via GitHub Security tab
- Direct maintainer contact listed in repository profile

Include:

- affected component and version/branch
- clear reproduction steps
- impact assessment (confidentiality, integrity, availability)
- potential mitigation or patch direction if known

## Response Expectations

- Acknowledgement target: within 72 hours
- Triage target: within 7 days
- Fix window: risk-based and coordinated with maintainers

## Disclosure Process

1. Maintainers reproduce and assess severity.
2. A fix is prepared in private when needed.
3. A coordinated disclosure date is agreed.
4. Public advisory/release notes are published after remediation.

## Security Controls In CI

Current baseline controls include:

- `go vet ./...`
- race detector (`go test -race ./...`)
- `govulncheck ./...`
- dependency drift checks via `go mod tidy`
- formatting and test gates

## Scope And Intentional Limits

Opslane is production-shaped for education and documents intentional boundaries in:

- [`docs/KNOWN_LIMITATIONS.md`](./docs/KNOWN_LIMITATIONS.md)
- [`docs/security/OPSLANE_THREAT_MODEL.md`](./docs/security/OPSLANE_THREAT_MODEL.md)

Please include whether your report targets a known limitation or an unexpected defect.
