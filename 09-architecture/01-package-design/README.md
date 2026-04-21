# Track A: Package Design

## Mission

This track teaches how to organize a Go codebase so package boundaries reinforce the design instead of hiding problems behind generic folders and global helpers.

## Stage Ownership

This track belongs to [09 Architecture & Security](../README.md).

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `PD.1` | Lesson | [naming](./1-naming) | Keeps package names readable and domain-specific. | entry |
| `PD.2` | Lesson | [visibility](./2-visibility) | Uses export rules and `internal/` to tighten boundaries. | `PD.1` |
| `PD.3` | Lesson | [project layout](./3-project-layout) | Chooses a layout that matches the real size of the system. | `PD.1`, `PD.2` |

## Suggested Order

1. Start with `PD.1` so naming rules are clear before structure gets deeper.
2. Continue to `PD.2` for export discipline and internal boundaries.
3. Finish with `PD.3` once the naming and visibility rules feel concrete.

## Next Step

After `PD.3`, continue to [Architecture Patterns](../03-architecture-patterns/1-architecture-trade-offs) or return to the [09 Architecture & Security overview](../README.md).
