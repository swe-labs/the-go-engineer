#!/usr/bin/env python3
import argparse, json
from pathlib import Path

HEADINGS = ["Mission", "Prerequisites", "Mental Model", "Visual Model", "Machine View", "Run Instructions", "Code Walkthrough", "Try It", "In Production", "Thinking Questions", "Next Step"]

def load_json(path):
    with open(path, encoding="utf-8") as f:
        return json.load(f)

def lesson_readme(item):
    zm = item.get("zero_magic") or {}
    title = item.get("title", item.get("id"))
    prereqs = item.get("prerequisites") or []
    run = (item.get("verification") or {}).get("run_command") or "Read the lesson and complete the practice task."
    test = (item.get("verification") or {}).get("test_command") or "No automated test is required for this lesson."
    mistakes = "\n".join(f"- {m}" for m in zm.get("beginner_mistakes", [])) or "- Confusing the surface syntax with the underlying model."
    return f"""# {title}\n\n## Mission\n\n{item.get('learning_objective') or zm.get('problem_solved') or 'Understand the concept and prove it with a small working example.'}\n\n## Prerequisites\n\n{chr(10).join(f'- {p}' for p in prereqs) if prereqs else '- None. This lesson establishes its own local context.'}\n\n## Mental Model\n\n{zm.get('mental_model') or 'Treat this concept as a small machine with inputs, state, behavior, and observable output.'}\n\n## Visual Model\n\n```text\ninput -> concept boundary -> behavior -> observable result\n```\n\n## Machine View\n\n{zm.get('under_the_hood') or 'The machine executes the code step by step. Focus on what data exists, where it moves, and what changes.'}\n\n## Run Instructions\n\n```bash\n{run}\n{test}\n```\n\n## Code Walkthrough\n\nRead the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.\n\n## Try It\n\n1. Run the command exactly as shown.\n2. Change one input.\n3. Predict the output before running again.\n4. Explain the result in your own words.\n\nCommon mistakes:\n\n{mistakes}\n\n## In Production\n\n{zm.get('real_world_usage') or 'Production systems use this idea to make behavior explicit, testable, and maintainable.'}\n\n## Thinking Questions\n\n1. What problem does this concept solve?\n2. What would break if this concept were removed?\n3. How would you test the behavior?\n\n## Next Step\n\nContinue to the next item listed in metadata: `{', '.join(item.get('next_item_ids') or ['module checkpoint'])}`.\n"""

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--metadata-dir", default="metadata")
    parser.add_argument("--only-missing", action="store_true")
    args = parser.parse_args()
    core = load_json(Path(args.metadata_dir) / "path.core.json")
    electives = load_json(Path(args.metadata_dir) / "path.electives.json")
    count = 0
    for item in core.get("items", []) + electives.get("items", []):
        if item.get("type") != "lesson":
            continue
        files = item.get("files") or {}
        readme = files.get("readme_path")
        if not readme:
            continue
        path = Path(readme)
        if args.only_missing and path.exists():
            continue
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(lesson_readme(item), encoding="utf-8")
        if files.get("main_path"):
            main_path = Path(files["main_path"])
            main_path.parent.mkdir(parents=True, exist_ok=True)
            if not main_path.exists() or not args.only_missing:
                main_path.write_text(f"package main\n\nimport \"fmt\"\n\nfunc main() {{\n\tfmt.Println(\"{item.get('title', item['id'])}\")\n}}\n", encoding="utf-8")
        if files.get("test_path"):
            test_path = Path(files["test_path"])
            test_path.parent.mkdir(parents=True, exist_ok=True)
            if not test_path.exists() or not args.only_missing:
                test_path.write_text("package main\n\nimport \"testing\"\n\nfunc TestLessonCompiles(t *testing.T) {\n\t// The lesson is validated by compilation and README-driven practice.\n}\n", encoding="utf-8")
        for key in ["starter_path", "solution_path", "assets_dir"]:
            if files.get(key):
                Path(files[key]).mkdir(parents=True, exist_ok=True)
        count += 1
    print(f"generated lessons: {count}")

if __name__ == "__main__":
    main()
