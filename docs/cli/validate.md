# socialpulse validate

Validate all content files against JSON schemas.

## Synopsis

```bash
socialpulse validate [flags]
```

## Description

Scans all content files (summaries and digests) and validates them against the embedded JSON schemas. Reports any validation errors with details about what fields are invalid.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `socialpulse.yaml` | Path to configuration file |

## Examples

### Validate all content

```bash
socialpulse validate
```

Output:

```
Validating content files...
  2026-03-25-ai-coding.yaml: OK
  2026-03-24-where-are-apps.yaml: OK
  2026-03-23-productive-claude.yaml: OK
Summary: 3 valid, 0 invalid
```

### With validation errors

```
Validating content files...
  2026-03-25-ai-coding.yaml: OK
  2026-03-24-invalid.yaml: INVALID
    - article.sentiment: must be one of [optimistic, cautious, pessimistic, neutral, provocative]
    - discussion.personas[0].stance: required field missing
Summary: 1 valid, 1 invalid
```

## Validation Checks

### Required Fields

- `schema_version`
- `meta.platform`
- `meta.article_url`
- `meta.discussion_url`
- `article.title`
- `article.thesis`
- `article.sentiment`
- `discussion.sentiment`
- `discussion.personas` (at least one)

### Enum Values

| Field | Valid Values |
|-------|--------------|
| `meta.platform` | `hackernews`, `reddit` |
| `article.sentiment` | `optimistic`, `cautious`, `pessimistic`, `neutral`, `provocative` |
| `discussion.sentiment` | `supportive`, `divided`, `critical`, `exploratory` |
| `personas[].stance` | `agrees`, `disagrees`, `nuanced`, `tangential` |
| `personas[].prevalence` | `dominant`, `significant`, `minority` |

### Persona Requirements

Each persona must have:

- `name` - Non-empty string
- `description` - Non-empty string
- `stance` - Valid enum value
- `prevalence` - Valid enum value
- `percentage` - Integer 0-100
- `core_argument` - Non-empty string
- `quotes` - At least one quote with `text` and `author`

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All files valid |
| 1 | One or more files invalid |

## Pre-commit Hook

Use validation as a pre-commit hook:

```bash
#!/bin/bash
# .git/hooks/pre-commit
socialpulse validate || exit 1
```
