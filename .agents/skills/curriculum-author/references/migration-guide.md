# v2 to v3 migration guide

## Migration decisions

Use one of these decisions for every v2 item:

- `keep`: content is strong and maps directly to a v3 lesson.
- `rewrite`: idea is useful but explanation/code is not zero-magic quality.
- `merge`: concept belongs inside another v3 lesson.
- `move-to-elective`: valuable but not required for the core path.
- `remove`: obsolete, duplicate, too framework-specific, or not aligned with the path.

## Required migration fields

Every migrated v3 item should preserve:

- `source_legacy_ids`
- mapping decision
- reason for decision
- target module/item/project
- final learner-facing path under `curriculum/modules/` or `curriculum/electives/`
- content-quality action required

## Migration audit questions

- Does this legacy content introduce a concept before prerequisites?
- Is it core job-readiness or advanced specialization?
- Does it duplicate another lesson?
- Does it require rewrite to remove magic?
- Does the new location have project/assessment coverage?
- Is the migration traceable from both directions?
- Does the final path use the typed layout: `lessons/`, `labs/`, `projects/`, or `assessments/`?

## Path mapping rules

Map core lessons to:

```text
curriculum/modules/{module}/lessons/{lesson}/
```

Map core labs to:

```text
curriculum/modules/{module}/labs/{lab}/
```

Map core projects to:

```text
curriculum/modules/{module}/projects/{project}/
```

Map core assessments to:

```text
curriculum/modules/{module}/assessments/{assessment}/
```

Map elective content to:

```text
curriculum/electives/{elective}/{lessons|labs|projects|assessments}/{item}/
```

Do not carry over old flat paths unless the repository explicitly runs a path migration script immediately after import.
