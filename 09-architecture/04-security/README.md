# Track SEC: Security

## Mission

Turn boundary safety into a concrete engineering track. This track covers the "Essential Defenses" for Go web applications, from basic input validation to advanced identity management and secret protection. Master the patterns that keep your users safe and your application resilient.

## Stage Ownership

This track belongs to [09 Architecture & Security](../README.md).

## Track Map

| ID | Type | Surface | Mission | Requires |
| --- | --- | --- | --- | --- |
| `SEC.1` | Lesson | [Input Validation](./1-input-validation-patterns) | Master the "Trust but Verify" principle. | entry |
| `SEC.2` | Lesson | [SQLi Prevention](./2-sql-injection-prevention) | Never concatenate SQL strings. | `SEC.1` |
| `SEC.3` | Lesson | [XSS and CSRF](./3-xss-and-csrf) | Protect the user's browser boundary. | `SEC.2` |
| `SEC.4` | Lesson | [Authentication Basics](./4-authentication-basics) | Verify "Who are you?" | `SEC.3` |
| `SEC.5` | Lesson | [JWT Implementation](./5-jwt-implementation-and-risks) | Master stateless identity tokens. | `SEC.4` |
| `SEC.6` | Lesson | [Password Hashing](./6-password-hashing) | Never store plaintext passwords. | `SEC.5` |
| `SEC.7` | Lesson | [Rate Limiting](./7-rate-limiting-patterns) | Protect resources from abuse. | `SEC.6` |
| `SEC.8` | Lesson | [TLS and HTTPS](./8-tls-and-https-in-go) | Secure data in transit. | `SEC.7` |
| `SEC.9` | Lesson | [Secrets Management](./9-secrets-management) | Keep keys out of source control. | `SEC.8` |
| `SEC.10` | Lesson | [OWASP Top 10](./10-owasp-top-10-for-go) | Use industry-standard checklists. | `SEC.9` |
| `SEC.11` | Exercise | [Secure API](./11-secure-api-exercise) | Harden a vulnerable application. | `SEC.1-10` |

## Next Step

After completing this track, you are ready for [Stage 10 Production Operations](../../10-production) or you can return to the [09 Architecture & Security overview](../README.md).
