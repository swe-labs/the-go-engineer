#!/usr/bin/env python3
from pathlib import Path

def write(path, content):
    p=Path(path); p.parent.mkdir(parents=True, exist_ok=True); p.write_text(content, encoding="utf-8")

def main():
    write("curriculum/shared/README.md", "# Shared Curriculum References\n\nShared guides used across modules.\n")
    write("curriculum/shared/glossary.md", "# Glossary\n\nAdd stable definitions as concepts are implemented.\n")
    write("curriculum/shared/zero-magic-guide.md", "# Zero-Magic Guide\n\nEvery concept must explain the problem, model, mechanics, mistakes, debugging, production use, and proof.\n")
    write("curriculum/shared/debugging-guide.md", "# Debugging Guide\n\nStart from the symptom, isolate the boundary, reproduce, inspect, fix, and verify.\n")
    write("curriculum/shared/portfolio-rubric.md", "# Portfolio Rubric\n\nProjects must be runnable, tested, documented, and explain trade-offs.\n")
    print("generated shared docs")
if __name__ == "__main__": main()
