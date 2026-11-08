# Shell + Git Lab — Rubric

This rubric scores the 14-question knowledge assessment (18 points: 10 multiple-choice
at 1 point, 4 short-answer at 2 points). Question-to-dimension mapping:

| Dimension | Questions | Points | Approx. weight |
|---|---|---|---|
| Git mental model and workflow | 2, 3, 4, 9, 11, 13 | 8 | ~45% |
| Web preview concepts (DNS, TCP, HTTP) | 6, 7, 8, 12, 14 | 7 | ~40% |
| Shell, exit codes, and pipes | 1, 5, 10 | 3 | ~15% |

## Dimension standards

### Git mental model and workflow (8 points)

| Level | Standard |
|-------|----------|
| Excellent (7-8 pts) | Correctly distinguishes working directory, staging area, and history; explains merge vs rebase; gives a correct cherry-pick/reset sequence for the wrong-branch scenario. |
| Proficient (5-6 pts) | Three-state model correct; minor imprecision on merge vs rebase or the recovery commands. |
| Developing (3-4 pts) | Names the three states but cannot trace transitions or recovery steps. |
| Not demonstrated (0-2 pts) | Git model answers are incorrect or absent. |

### Web preview concepts (7 points)

| Level | Standard |
|-------|----------|
| Excellent (6-7 pts) | Correct DNS-to-IP-to-TCP-to-HTTP sequence; correct method and status-code semantics; raw GET request written with method, path, version, and headers identified. |
| Proficient (4-5 pts) | Sequence mostly correct; omits TLS or conflates DNS stages; request format mostly right. |
| Developing (2-3 pts) | Knows individual terms (DNS, TCP, 404) but cannot order them into a request lifecycle. |
| Not demonstrated (0-1 pt) | Web answers are incorrect or absent. |

### Shell, exit codes, and pipes (3 points)

| Level | Standard |
|-------|----------|
| Excellent (3 pts) | Correct on directory creation, exit code 0 semantics, and pipe data flow. |
| Proficient (2 pts) | Two of three correct. |
| Not demonstrated (0-1 pt) | One or zero correct. |

## Scoring guide

| Total points | Result |
|--------------|--------|
| 16-18 | Exceeds expectations |
| 13-15 | Meets expectations — module passed (72%+) |
| 9-12 | Needs improvement — retake recommended |
| 0-8 | Incomplete — retake required after 48 hours |

A passing score is at least 13 of 18 points (72%). Filesystem mechanics from lesson 3
(inodes, CRLF, permissions) are verified by the capstone project's 14-check suite, not
by this knowledge assessment.
