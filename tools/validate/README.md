# Validation

Validation is split into two layers.

## `validate/curriculum/`

Strict Go validation used by CI and release gates.

It validates:

- metadata JSON parsing
- module and item graph integrity
- cross-reference targets
- concept ownership and reinforcement
- project and assessment bindings
- failure taxonomy coverage
- README contract definitions
- canonical typed curriculum paths
- repository file existence and content quality when strict mode is enabled

## `validate/repository/`

Python validation helpers for focused checks on learner-facing files.

These tools are useful during authoring and local repair work. They should support the Go validator, not replace it.
