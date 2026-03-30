# Content Authoring Guide

This guide explains how to write effective discussion summaries for SocialPulse.

## File Format

Summaries can be written in YAML or JSON. YAML is recommended for readability.

### File Naming

Use descriptive, URL-friendly names:

```
content/summaries/
├── 2026-03-25-ai-coding-slowdown.yaml
├── 2026-03-24-where-are-ai-apps.yaml
└── 2026-03-23-claude-productivity.yaml
```

Pattern: `YYYY-MM-DD-slug.yaml`

## Summary Structure

Every summary has three main sections:

```yaml
schema_version: "1.0.0"

meta:
  # Source metadata

article:
  # Article summary

discussion:
  # Discussion analysis
```

## Meta Section

The meta section identifies the source:

```yaml
meta:
  article_url: "https://example.com/original-article"
  discussion_url: "https://news.ycombinator.com/item?id=12345678"
  platform: hackernews  # or 'reddit'
  article_date: 2026-03-25T00:00:00Z
  summary_date: 2026-03-26T00:00:00Z
  article_author: "Author Name"
  discussion_comment_count: 150
  item_id: "12345678"  # Platform-specific ID
```

### Platform-Specific Fields

| Platform | Required Fields |
|----------|-----------------|
| HackerNews | `item_id` |
| Reddit | `item_id`, `subreddit` |

## Article Section

Summarize the source article:

```yaml
article:
  title: "The Original Article Title"

  thesis: |
    The main argument in 200-400 characters. This should capture
    the core message that readers will remember.

  key_arguments:
    - "First major point with evidence"
    - "Second point addressing counterarguments"
    - "Third point with implications"

  examples:
    - "Concrete example from the article"
    - "Another supporting case study"

  prescription: "What the author recommends doing"

  tags:
    - ai-coding
    - software-engineering
    - productivity

  sentiment: cautious
```

### Writing the Thesis

The thesis should:

- Capture the author's main argument
- Be understandable without reading the article
- Be 200-400 characters
- Avoid jargon unless necessary

**Good:**
> AI coding agents encourage developers to move fast without understanding, producing brittle code. True craftsmanship requires slowing down and maintaining ownership.

**Too vague:**
> The author discusses AI and coding.

### Key Arguments

List 3-5 main points:

- Each should be self-contained
- 100-200 characters each
- Use active voice
- Include evidence when possible

### Tags

Use lowercase, hyphenated tags:

- `ai-coding` not `AI Coding`
- Be specific but not too narrow
- Aim for 3-7 tags per article

### Sentiment Values

| Value | Use When |
|-------|----------|
| `optimistic` | Article is hopeful about the subject |
| `cautious` | Article raises concerns but sees potential |
| `pessimistic` | Article is negative or critical |
| `neutral` | Article is informational without strong opinion |
| `provocative` | Article is intentionally controversial |

## Discussion Section

Analyze the community response:

```yaml
discussion:
  sentiment: divided

  personas:
    - name: "The Pragmatist"
      description: "Practitioners focused on what works"
      stance: nuanced
      prevalence: dominant
      percentage: 35
      core_argument: "Theory is nice but execution matters"
      quotes:
        - text: "Actual quote from the discussion"
          author: "username"

  tangents:
    - topic: "Related Discussion Topic"
      description: "What commenters said about this tangent"

  consensus_points:
    - "Something most people agreed on"

  open_questions:
    - "Question that remained unresolved"
```

### Discussion Sentiment

| Value | Use When |
|-------|----------|
| `supportive` | Discussion mostly agrees with article |
| `divided` | Strong opinions on both sides |
| `critical` | Discussion mostly disagrees |
| `exploratory` | Discussion explores tangents more than the thesis |

## Personas

See [Writing Personas](personas.md) for detailed guidance.

## Quality Checklist

Before committing a summary:

- [ ] Schema version is "1.0.0"
- [ ] All URLs are valid and accessible
- [ ] `item_id` matches the discussion URL
- [ ] Thesis is 200-400 characters
- [ ] 3-5 key arguments listed
- [ ] At least 2 personas defined
- [ ] Each persona has 1-2 quotes
- [ ] Quotes are verbatim (use `fetch-quotes` to verify)
- [ ] Tags are lowercase and hyphenated
- [ ] Sentiment values are valid enums

## Validation

Always validate before committing:

```bash
socialpulse validate
```

Then verify against the platform API:

```bash
socialpulse verify --quotes
```
