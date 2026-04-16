# Quality & Test Stage Map

## Stage Goal

This stage teaches learners how to prove correctness and measure performance instead of guessing
about either one.

## Public Stage Shape

| Track | Source Surface | Role In Stage |
| --- | --- | --- |
| Testing | [testing](../../../12-quality-and-performance/testing/) | teaches unit tests, table-driven tests, handler tests, and benchmarking |
| Profiling | [profiling](../../../12-quality-and-performance/profiling/) | teaches offline CPU profiling and live `pprof` inspection |
| Reference | [http-client-testing](../../../12-quality-and-performance/http-client-testing/) | extra testing depth outside the current public path |
| Reference | [4-testcontainers](../../../12-quality-and-performance/4-testcontainers/) | environment-shaped testing beyond the current public path |

## Live Milestone Backbone

| ID | Surface | Why It Matters |
| --- | --- | --- |
| `TE.4` | benchmarking | proves learners can measure behavior instead of guessing about speed |
| `PR.2` | live pprof endpoint | proves learners can inspect running programs and reason about hot paths |

## Recommended Order

1. `TE.1` through `TE.4`
2. `PR.1` through `PR.2`

## Reference Surfaces That Still Matter

These are still useful, but they are not the current public beta backbone:

- `http-client-testing`
- `4-testcontainers`
