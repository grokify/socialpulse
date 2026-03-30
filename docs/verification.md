# Verification

SocialPulse includes anti-hallucination tools to verify content accuracy against platform APIs.

## Why Verification Matters

When using LLMs to generate discussion summaries, hallucinations can occur:

- **Fake URLs** - Links that don't exist
- **Wrong metadata** - Incorrect comment counts, dates, titles
- **Fabricated quotes** - Text attributed to users who never wrote it
- **Invented usernames** - Authors that don't exist

Verification catches these issues before publishing.

## Verification Workflow

```bash
# 1. Validate structure
socialpulse validate

# 2. Verify against platform APIs
socialpulse verify

# 3. Check quotes (slower, fetches comments)
socialpulse verify --quotes

# 4. Auto-fix metadata discrepancies
socialpulse verify --fix

# 5. Fetch real quotes for personas
socialpulse fetch-quotes
```

## What Gets Verified

### Metadata Verification

| Field | Check | Auto-Fix |
|-------|-------|----------|
| `item_id` | Exists in API | No |
| `article_url` | Matches submission | Yes |
| `discussion_url` | Correct format | Yes |
| `title` | Matches submission | Yes |
| `comment_count` | Within 10% of API | Yes |

### Quote Verification

| Check | Description |
|-------|-------------|
| Author exists | Username appears in discussion comments |
| Text similarity | Quote text matches actual comment (>30% overlap) |

## Commands

### socialpulse verify

Basic verification against platform APIs:

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
```

### socialpulse verify --quotes

Fetches comments to verify persona quotes:

```bash
socialpulse verify --quotes --verbose
```

```
  2026-03-25-slowing-down.yaml:
    Fetching comments to verify 10 quotes...
    Fetched 458 comments from 312 authors
    - Quote text not found: "This exact phrase..." by "user123" (best match: 15% similarity)
  2026-03-25-slowing-down.yaml: issues found
```

### socialpulse verify --fix

Auto-fixes metadata discrepancies:

```bash
socialpulse verify --fix
```

This updates:

- `article_url` - Corrected to API value
- `discussion_url` - Corrected format
- `title` - Updated to match
- `comment_count` - Updated to current value

### socialpulse fetch-quotes

Replaces quotes with verified ones from actual comments:

```bash
socialpulse fetch-quotes --verbose
```

See [fetch-quotes command](cli/fetch-quotes.md) for details.

## Platform Support

| Platform | Metadata | Quotes | API |
|----------|----------|--------|-----|
| HackerNews | Full | Full | Firebase API |
| Reddit | Planned | Planned | Reddit API |

## Best Practices

### 1. Verify Early and Often

```bash
# After creating a summary
socialpulse validate
socialpulse verify
```

### 2. Use Fetch-Quotes for New Summaries

Don't manually copy quotes. Use the tool:

```bash
socialpulse fetch-quotes --dry-run --verbose
# Review matches
socialpulse fetch-quotes
```

### 3. CI/CD Integration

Add verification to your pipeline:

```yaml
# GitHub Actions
- name: Validate and verify
  run: |
    socialpulse validate
    socialpulse verify
```

### 4. Regular Re-verification

Comment counts change over time. Periodically update:

```bash
socialpulse verify --fix
git diff
git commit -m "chore: update comment counts"
```

## Troubleshooting

### "Item not found"

The `item_id` doesn't exist in the platform API.

**Causes:**

- Incorrect ID
- Submission was deleted
- Private/hidden submission

**Solution:** Verify the discussion URL is accessible and extract the correct ID.

### Low Quote Similarity

Quotes don't match any comment text.

**Causes:**

- Quote was paraphrased, not verbatim
- Comment was edited after summarization
- Quote was hallucinated by LLM
- Comment was deleted

**Solution:** Use `fetch-quotes` to get verified quotes.

### "Author not found"

The quoted username doesn't appear in the discussion.

**Causes:**

- Username was hallucinated
- User deleted their account
- Comment was deleted

**Solution:** Use `fetch-quotes` to get quotes from real users.

### Rate Limiting

For large sites with many summaries:

```bash
# Process one file at a time
for f in content/summaries/*.yaml; do
  socialpulse verify --quotes -c socialpulse.yaml
  sleep 5
done
```

## Understanding Similarity Scores

Quote verification uses word overlap similarity:

```
similarity = matching_words / total_words_in_quote
```

| Score | Interpretation |
|-------|----------------|
| >0.7 | Strong match, likely verbatim |
| 0.3-0.7 | Partial match, may be paraphrased |
| <0.3 | Weak match, likely different text |

The threshold for "found" is 0.3 (30%).
