"""Add missing files field to core-10-25."""
import json

with open('curriculum/path.core.json', 'r', encoding='utf-8') as f:
    core = json.load(f)

for item in core['items']:
    if item.get('id') == 'core-10-25':
        item['files'] = {
            "assets_dir": "10-security/25-refresh-tokens/assets",
            "diagram_paths": [],
            "main_path": "10-security/25-refresh-tokens/main.go",
            "readme_path": "10-security/25-refresh-tokens/README.md",
            "solution_path": "",
            "starter_path": "",
            "test_path": "10-security/25-refresh-tokens/main_test.go"
        }
        item['source_legacy_ids'] = []
        item['tags'] = ["security", "auth", "jwt"]
        item['documentation_mode'] = "tutorial"
        item['readme_status'] = "scaffolded"
        item['readme_contract'] = {
            "contract_id": "lesson.v3",
            "documentation_mode": "tutorial",
            "required": True
        }
        print(f"Added missing fields to core-10-25")
        break

with open('curriculum/path.core.json', 'w', encoding='utf-8') as f:
    json.dump(core, f, indent=2, ensure_ascii=False)
print("Written")
