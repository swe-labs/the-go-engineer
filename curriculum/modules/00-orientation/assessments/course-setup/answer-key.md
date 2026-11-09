# Answer Key

## Part 1 — Environment verification

1. **Expected:** output like `git version 2.4X.X` (any current version is acceptable).
2. **Expected:** exit code `0`, meaning the process completed successfully without
   errors (Lesson 04: zero is the universal success contract; non-zero signals failure).
3. **Expected:** two non-empty lines — the configured committer name and email. If one
   was blank, the fix is `git config --global user.name "Your Name"` and
   `git config --global user.email "you@example.com"`.
4. **Expected:** the editor name plus its verification command, for example VS Code via
   `code --version`, or confirmation that the IDE launches. Any real editor is valid.

## Part 2 — Workspace and evidence

5. **Expected:** the curriculum repository stays a clean, authoritative reference
   (your "textbook"); personal code, experiments, and project implementations live in
   your own workspace so nothing in the curriculum is modified or committed by accident.
6. **Expected:** Bash `mkdir -p ~/the-go-engineer-workspace`; PowerShell
   `mkdir ~/the-go-engineer-workspace` (PowerShell's `mkdir` creates parent directories;
   it does not accept `-p`).
7. **Expected:** `setup.txt` records three things — the Git version, the committer
   name, and the committer email (produced by `git --version`, `git config user.name`,
   `git config user.email`).
8. **Expected:** `>` creates (or overwrites) the file with the first command's output;
   `>>` appends subsequent output without destroying what is already there.

## Part 3 — Curriculum contract

9. **Expected:** zero magic means no required step is hidden — every command, file, and
   behavior is explained mechanically. A setup-process example: running
   `git config --global user.email ...` is explicit, verifiable (re-run the read
   command), and its effect (non-empty output) is observable.
10. **Expected:** Absorb (Lesson) -> Practice (Guided Lab) -> Synthesize (Project) ->
    Verify (Assessment Checkpoint). Artifacts: lesson README, guided drills, the
    finished project, the checkpoint submission.
11. **Expected:** 1) Expected behavior, 2) Actual observed behavior, 3) Exact command
    and working directory, 4) Full verbatim error output, 5) Environment context
    (OS, tool versions, shell).
12. **Expected:** a completed 5-field case file, for example — Expected: a non-empty
    email line; Actual: blank output; Command & path: `git config user.email` run from
    `~`; Full error/output: empty result (no error text); Environment: OS, Git version,
    shell. Hypothesis: identity was never configured globally; verification:
    `git config --global user.email "you@example.com"` then re-run the read command.
