# Track A: Package Design

## Mission

This track teaches you how to organize a Go codebase so package boundaries reinforce the design
instead of hiding problems behind generic folders and global helpers.

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

## Track Milestone

`PD.3` is the current package-design output.

If you can explain:

- why package names should describe a domain instead of a grab-bag
- why `internal/` is a compiler-enforced boundary, not just a folder convention
- why a layout should grow with the codebase instead of being over-designed on day one

then the package-design part of Section 14 is doing its job.

## Next Step

After `PD.3`, continue to the [Structured Logging track](../structured-logging) or back to the
[Section 14 overview](../README.md).
