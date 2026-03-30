# socialpulse fetch-quotes

Fetch actual quotes from discussions to replace hallucinated ones.

## Synopsis

```bash
socialpulse fetch-quotes [flags]
```

## Description

Fetches real comments from HackerNews discussions and matches them to personas based on keyword analysis. This replaces potentially hallucinated quotes with verified, real quotes from actual commenters.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `socialpulse.yaml` | Path to configuration file |
| `--dry-run` | `false` | Show matches without updating files |
| `-v, --verbose` | `false` | Show detailed matching output |
| `--max-quotes` | `2` | Maximum quotes per persona |

## Examples

### Preview matches (dry-run)

```bash
socialpulse fetch-quotes --dry-run --verbose
```

Output:

```
Fetching quotes for 6 summary files...

  2026-03-25-slowing-down.yaml:
    Fetching comments from HN...
    Fetched 458 comments
    [The Pragmatic Senior] Found: @simonw (score: 4.20, keywords: [context, coding, agent])
      "Useful context here is that the author wrote Pi, which is..."
    [The Pragmatic Senior] Found: @forgeties79 (score: 3.80, keywords: [digital, cameras])
      "The thing is though it all still feels so rudderless..."
    [The Speed Defender] Found: @drzaiusx11 (score: 3.60, keywords: [productivity, developer])
      "The productivity gains are somewhat real in a sense..."
    [The Pragmatic Senior] 2 quotes matched
    [The Speed Defender] 2 quotes matched
    (dry-run, not saved)
```

### Update files with real quotes

```bash
socialpulse fetch-quotes
```

### Limit quotes per persona

```bash
socialpulse fetch-quotes --max-quotes 3
```

## Matching Algorithm

The algorithm scores comments based on keyword relevance to each persona:

### 1. Keyword Extraction

Keywords are extracted from:

- Persona `description`
- Persona `core_argument`
- Stance-specific terms (e.g., "agree", "disagree")

### 2. Comment Scoring

Each comment is scored:

```
score = sum(keyword_length / 5.0) for each matched keyword
```

Adjustments:

| Condition | Multiplier |
|-----------|------------|
| 20-200 words | 1.2x bonus |
| <10 words | 0.5x penalty |
| >300 words | 0.7x penalty |

### 3. Quote Selection

Top N scoring comments are selected for each persona, where N = `--max-quotes`.

## Stance Keywords

| Stance | Additional Keywords |
|--------|---------------------|
| `agrees` | agree, right, correct, exactly, yes, true |
| `disagrees` | disagree, wrong, incorrect, but, however, no |
| `nuanced` | depends, context, both, nuance, tradeoff |
| `tangential` | related, tangent, also, reminds, similar |

## Quote Cleaning

Before storing quotes:

1. **Remove quoted parent text** - Lines starting with `>` are stripped
2. **Truncate long quotes** - Limited to ~300 characters at sentence boundary
3. **Normalize whitespace** - Multiple spaces collapsed

## Best Practices

### Review Before Committing

Always review fetched quotes:

```bash
socialpulse fetch-quotes --dry-run --verbose
# Review output
socialpulse fetch-quotes
git diff content/summaries/
```

### Verify After Fetching

```bash
socialpulse fetch-quotes
socialpulse verify --quotes
```

### Manual Curation

The algorithm finds relevant quotes, but manual curation may improve quality:

- Select quotes that best represent the persona
- Ensure quotes are substantive, not just agreement/disagreement
- Check that attribution is correct
