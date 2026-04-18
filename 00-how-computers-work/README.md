# Section 00: How Computers Work

> **Philosophy:** You cannot write great code without understanding what the machine is actually doing. Every senior engineer you admire has a mental model of the computer underneath their code. This section builds that model.

---

## Before You Start

Most programming courses begin with syntax. This one doesn't.

Here's why: imagine trying to drive a car without understanding that pressing the accelerator burns fuel, that braking uses friction, that the engine needs oil. You could learn the mechanical steps — turn key, shift gear, press pedal — and still be a terrible driver because you have no model of what's happening underneath.

Code is the same. You can memorise syntax and still write slow, broken, dangerous programs. But engineers who understand the machine write code that is fast, reliable, and predictable — because they can *reason* about what it's doing, not just hope it works.

Work through this section slowly. These five lessons are the foundation everything else rests on.

## Section Map

| ID | Type | Surface | Core Job |
| --- | --- | --- | --- |
| `HC.1` | Lesson | [what-is-a-program](./1-what-is-a-program) | fetch-decode-execute cycle |
| `HC.2` | Lesson | [compilation-journey](./2-compilation-journey) | ast and compilation |
| `HC.3` | Lesson | [memory-basics](./3-memory-basics) | stack and heap |
| `HC.4` | Lesson | [terminal-confidence](./4-terminal-confidence) | stdout/stderr |
| `HC.5` | Lesson | [os-processes](./5-os-processes) | signals and fds |

## Checkpoint

Before moving to Section 01, you should be able to answer these without looking anything up:

**Conceptual:**
- [ ] What is the fetch-decode-execute cycle?
- [ ] What are the six basic CPU operations?
- [ ] What is the difference between compile time and runtime?
- [ ] When does Go catch type errors — at compile time or runtime?
- [ ] What is the difference between stack and heap memory?
- [ ] What is escape analysis?
- [ ] What is a process? How is it different from a program?
- [ ] What are the three default file descriptors?
- [ ] What signal does Ctrl+C send?

**Terminal:**
- [ ] Navigate to any directory and list its contents
- [ ] Find all files matching a pattern
- [ ] Search for a string inside a file
- [ ] Chain two commands so the second only runs if the first succeeds
- [ ] Read the exit code of the last command
