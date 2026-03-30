# socialpulse verify

Verify source URLs and quotes against platform APIs.

## Synopsis

```bash
socialpulse verify [flags]
```

## Description

Verifies that article and discussion URLs in summary files match the actual data from platform APIs (e.g., HackerNews Firebase API). This helps detect hallucinated or incorrect URLs in content files.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `socialpulse.yaml` | Path to configuration file |
| `--fix` | `false` | Automatically fix discrepancies in YAML files |
| `--quotes` | `false` | Verify persona quotes exist in actual comments |
| `-v, --verbose` | `false` | Show detailed output |

## Examples

### Basic verification

```bash
socialpulse verify
```

Output:

```
Verifying 6 summary files against platform APIs...

  2026-03-25-slowing-down.yaml: OK
  2026-03-25-claude-stats.yaml: issues found
    - comment_count:
        File: 150
        API:  212
  2026-03-24-ai-apps.yaml: OK

Summary: 5 OK, 1 with discrepancies, 0 errors

Run with --fix to automatically correct discrepancies.
```

### Auto-fix discrepancies

```bash
socialpulse verify --fix
```

Updates YAML files with correct values from the API.

### Verify quotes

```bash
socialpulse verify --quotes --verbose
```

Fetches actual comments and verifies that quoted text exists:

```
Verifying 6 summary files against platform APIs...

  2026-03-25-slowing-down.yaml:
    Fetching comments to verify 10 quotes...
    Fetched 458 comments from 312 authors
  2026-03-25-slowing-down.yaml: OK (10 quotes verified)
```

## Checks Performed

### URL Verification

| Check | Description |
|-------|-------------|
| `article_url` | Matches URL in HN submission |
| `discussion_url` | Correct format for item ID |
| `item_id` | Exists in platform API |

### Metadata Verification

| Check | Description |
|-------|-------------|
| `title` | Matches submission title |
| `comment_count` | Within 10% of API value |

### Quote Verification (with `--quotes`)

| Check | Description |
|-------|-------------|
| Author exists | Username appears in discussion |
| Text similarity | Quote text matches actual comment (>30% word overlap) |

## Discrepancy Handling

### Auto-fix (`--fix`)

These fields can be auto-fixed:

- `article_url`
- `discussion_url`
- `title`
- `comment_count`

### Manual Review Required

Quote verification issues require manual review:

```
- Quote author not found: "nonexistent_user" in persona "The Skeptic"
- Quote text not found: "This exact text..." by "real_user" (best match: 15% similarity)
```

Use `socialpulse fetch-quotes` to fetch real quotes.

## Platform Support

| Platform | Verification | Quote Check |
|----------|--------------|-------------|
| HackerNews | Full | Full |
| Reddit | Planned | Planned |

## Rate Limiting

The verify command respects API rate limits:

- HackerNews: 50ms delay between requests
- Maximum 500 comments fetched per article

## Troubleshooting

### "Item not found"

The `item_id` doesn't exist. Check that:

- The ID is correct
- The submission hasn't been deleted

### Low quote similarity

Quotes may have been:

- Edited after summarization
- Paraphrased instead of verbatim
- Hallucinated by the LLM

Use `socialpulse fetch-quotes` to fetch verified quotes.
