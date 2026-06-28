# Pull requests and code review

## Learning objective

Understand and apply pull requests and code review in the context of professional Go software engineering.

## Why this matters

Writing code alone leads to blind spots — the author knows what the code should do and may not notice edge cases, readability issues, or design flaws. Code review provides a second (and third, fourth) pair of eyes to catch issues early.

## Mental model

Code review is pair debugging at scale — the author writes the solution, and reviewers validate assumptions, catch blind spots, and ensure consistency. The goal is not to find every bug but to share context, raise the team's collective code quality bar, and spread knowledge.

## Core idea

A pull request is a proposal to merge changes from one branch into another, with attached discussion, review comments, and automated checks. Code review is the process of examining those changes to catch defects, improve design, and share context before the code reaches production.

## Under the hood

GitHub stores PR reviews as events linked to specific commits. Review comments on the "Files changed" tab reference specific lines in specific diffs. "Request changes" blocks merge via branch protection. "Approve" signals the PR is ready to merge. The merge queue can enforce that status checks pass before allowing the merge button.

## How Go uses it

Go code review follows strict conventions documented in "CodeReviewComments" on the Go wiki. Reviewers check for proper error handling, correct use of interfaces, test coverage, and idiomatic Go style. The `go vet` and `gofmt` tools automate style and correctness checks.

## Go example

The example program simulates a code review process: a diff with changed lines, a reviewer who inspects each change and leaves comments, and a merge decision based on review outcomes.

## Step-by-step execution

1. Open the PR on GitHub — review the "Files Changed" tab to see the full diff.
2. Read each changed file starting from the most critical (business logic) to the least (tests, config).
3. For each change, ask: is this correct? Is it clear? Is it tested? Could it break anything?
4. Leave inline comments for specific lines, and a summary comment with overall feedback.
5. The author addresses feedback, and the cycle repeats until approval — then the PR is merged.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Submitting a massive PR with 50+ files changed | Combining multiple changes into one PR | Break large changes into smaller, incremental PRs (under 300 lines) |
| Taking review feedback personally | Confusing code critique with personal criticism | Remember: review is about code quality, not the author |
| Merging without CI checks passing | Impatience or ignoring red status checks | Never merge until all required status checks pass |

## Debugging walkthrough

**Scenario: A code review cycle takes 5+ days because the PR is too large.**
- **Cause:** Large PRs are harder to review; reviewers procrastinate on big changes.
- **Diagnosis:** Check the PR's line count and file count.
- **Resolution:** Break the change into focused PRs under 300 lines changed. Each PR should do one thing.

**Scenario: A developer resolves a review comment with a force-push, losing discussion history.**
- **Cause:** Force-pushing after rebasing changes commit hashes, disconnecting comments from code.
- **Diagnosis:** Review comments reference old commit hashes that no longer exist.
- **Resolution:** Add new commits to address feedback instead of rebasing. If rebasing is necessary, only rebase before review, not after.

## Production notes

Every professional software team uses code review. Google mandates at least one review before any code is merged. Open-source projects like Go, Kubernetes, and VS Code require reviews from maintainers. Code review is the primary quality gate in modern software development.

## Performance implications

- Each review cycle adds 1-24 hours to delivery time — smaller, focused PRs reduce this latency.
- Thorough code review catches 60-70% of defects before they reach production, saving 10x the cost of fixing them later.
- Mechanical style issues should be automated (formatters, linters) so reviewers can focus on logic and design.

## Tests / verification

1. Open a pull request on GitHub (or use an existing open PR on a public repo like `golang/go`). Read the diff in the "Files changed" tab.
2. Leave at least three inline comments: one suggesting a functional improvement, one asking for clarification, and one praising a well-written section.
3. Practice the review workflow: request changes, then approve. Note how the PR status changes after each action.

Although the [Shell + Git Lab](../projects/shell-git-lab/README.md) is a local repository exercise, the code review skills you learn here — writing clear descriptions, reading diffs, leaving constructive feedback — apply directly when you share your project on GitHub and open your first real PR.

## Review questions

1. What information should a good PR description include?
2. How should you respond to a code review comment requesting a change?
3. What is the difference between a blocking and a non-blocking review comment?
4. Why should style feedback be automated rather than manual?
5. What does "approve" mean in a code review context?

## NEXT UP

[Lesson 14: Web preview — client, server, DNS, and ports](../14-web-preview-client-server-dns-and-ports/README.md)
