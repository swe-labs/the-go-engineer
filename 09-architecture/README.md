# 09 Architecture & Security

## Mission

This stage expands package design into architecture patterns and security engineering so system shape and trust boundaries are taught together. It focuses on how to structure large Go applications.

By the end of this stage, a learner should be able to:

- design sensible package layouts using `cmd`, `internal`, and `pkg`
- implement clean architecture patterns like dependency injection and adapters
- identify and mitigate common web vulnerabilities
- secure APIs with proper authentication, authorization, and TLS
- understand how trust boundaries shape code structure

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `PD.1-3` | Package Design | Introduce package boundaries, naming conventions, and project layouts. |
| `ARCH.1-9` | Architecture Patterns | Cover the major service-shape and layering decisions teams face. |
| `SEC.1-11` | Security | Turn boundary safety into a concrete engineering track. |

## Why This Stage Exists Now

The learner already knows:

- how to write correct, concurrent, and tested Go code
- how to profile for performance

That is enough to start asking engineering questions like:

- how do we organize a codebase with 100,000 lines of code?
- how do we keep business logic separate from HTTP routing?
- how do we ensure attackers cannot manipulate our database or access private data?

## Suggested Learning Flow

1. `PD.1-PD.3` establish the foundational rules of Go package design.
2. `ARCH.1-ARCH.9` expand those rules into application-wide architecture.
3. `SEC.1-SEC.11` add the necessary security boundaries to that architecture.

## Next Step

After this section, continue to [10 Production Operations](../10-production).
