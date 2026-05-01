# PD.1 Naming Conventions

## Mission

Master the "Go Way" of naming packages. Learn how to choose concise, descriptive, and domain-focused names that make your code self-documenting and avoid common pitfalls like "util," "common," or "helper."

## Prerequisites

- None (Foundational concept).

## Mental Model

Think of Package Naming as **Labeling Boxes in a Warehouse**.

1. **The Purpose**: A box labeled "Tools" is useless. A box labeled "Screwdrivers" is clear.
2. **The Context**: You don't label the screwdrivers inside the box as "ScrewdriverA," "ScrewdriverB." They are just "Philips," "Flathead."
3. **The Result**: When you use them, you say `screwdrivers.Philips`, which is perfect. You don't want to say `tools.PhilipsScrewdriver` (stuttering).

## Visual Model

| ❌ Bad Name | ✅ Good Name | Why? |
| --- | --- | --- |
| `util`, `common` | `json`, `path`, `email` | Utilities should be grouped by their actual domain. |
| `interfaces` | `provider`, `service` | Don't name a package after a technical construct. |
| `my_cool_pkg` | `coolpkg` | Go package names should be lowercase, single words. |

## Machine View

- **Import Path**: The package name is the last part of the import path.
- **Identifier Selection**: Go identifiers are scoped by package. The name of the package should provide context for the names inside it.
- **Lowercase**: Always use lowercase. Avoid underscores or camelCase.

## Run Instructions

```bash
# Run the demo to see how naming impacts readability
go run ./09-architecture/01-package-design/1-naming
```

## Code Walkthrough

### Stuttering Example
Shows why `user.UserService` is redundant compared to `user.Service`.

### Generic Package Example
Shows how a `util` package becomes a "junk drawer" that makes dependency management difficult.

## Try It

1. Look at `main.go`. Identify three instances of "stuttering" in the variable or type names.
2. Refactor the code to remove the stuttering.
3. Try to split the `util` package into two domain-specific packages (e.g., `math` and `stringutil`).

## In Production
**Naming is a signal of design quality.** If you find it hard to name a package, it's often a sign that the package is doing too many unrelated things. A well-named package has a "Single Responsibility."

## Thinking Questions
1. Why are `util` and `common` considered "code smells" in Go?
2. Should a package name be plural (`users`) or singular (`user`)?
3. How does package naming affect how you write documentation (GoDoc)?

## Next Step

Next: `PD.2` -> `09-architecture/01-package-design/2-visibility`

Open `09-architecture/01-package-design/2-visibility/README.md` to continue.
