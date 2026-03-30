# Schema Reference

SocialPulse uses JSON Schema for content validation. This page documents all fields and their requirements.

## Summary Schema

The complete summary schema with all fields:

### Root Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `schema_version` | string | Yes | Schema version (currently "1.0.0") |
| `meta` | object | Yes | Source metadata |
| `article` | object | Yes | Article summary |
| `discussion` | object | Yes | Discussion analysis |

### Meta Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `article_url` | string | Yes | URL of the original article |
| `discussion_url` | string | Yes | URL of the social discussion |
| `platform` | enum | Yes | `hackernews` or `reddit` |
| `article_date` | datetime | Yes | Publication date (ISO 8601) |
| `summary_date` | datetime | Yes | Summary creation date |
| `article_author` | string | Yes | Author name |
| `discussion_comment_count` | integer | Yes | Number of comments |
| `item_id` | string | No | Platform-specific ID |
| `subreddit` | string | No | Reddit subreddit name |
| `score` | integer | No | Discussion score/points |

### Article Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Article title |
| `thesis` | string | Yes | Core argument (200-400 chars) |
| `key_arguments` | string[] | Yes | Main points (3-5 items) |
| `examples` | string[] | No | Concrete examples cited |
| `prescription` | string | No | Author's recommendation |
| `tags` | string[] | Yes | Categorical tags |
| `sentiment` | enum | Yes | Article sentiment |

#### Article Sentiment Enum

```yaml
sentiment: optimistic | cautious | pessimistic | neutral | provocative
```

### Discussion Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `sentiment` | enum | Yes | Discussion sentiment |
| `personas` | Persona[] | Yes | Commenter archetypes |
| `tangents` | Tangent[] | No | Sub-topics |
| `consensus_points` | string[] | No | Points of agreement |
| `open_questions` | string[] | No | Unresolved questions |

#### Discussion Sentiment Enum

```yaml
sentiment: supportive | divided | critical | exploratory
```

### Persona Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Persona name (e.g., "The Skeptic") |
| `description` | string | Yes | Brief description (50-100 chars) |
| `stance` | enum | Yes | Relation to thesis |
| `prevalence` | enum | Yes | How common this viewpoint is |
| `percentage` | integer | Yes | Estimated percentage (0-100) |
| `quotes` | Quote[] | Yes | Verbatim quotes |
| `core_argument` | string | Yes | Distilled argument (100-150 chars) |

#### Stance Enum

```yaml
stance: agrees | disagrees | nuanced | tangential
```

#### Prevalence Enum

```yaml
prevalence: dominant | significant | minority
```

### Quote Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `text` | string | Yes | Verbatim quote text |
| `author` | string | Yes | Username/handle |

### Tangent Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `topic` | string | Yes | Tangent topic name |
| `description` | string | Yes | Brief description |

## Complete Example

```yaml
schema_version: "1.0.0"

meta:
  article_url: "https://example.com/article"
  discussion_url: "https://news.ycombinator.com/item?id=12345678"
  platform: hackernews
  article_date: 2026-03-25T00:00:00Z
  summary_date: 2026-03-26T00:00:00Z
  article_author: "Jane Doe"
  discussion_comment_count: 150
  item_id: "12345678"
  score: 500

article:
  title: "Article Title Here"
  thesis: "The main argument of the article summarized clearly."
  key_arguments:
    - "First key point"
    - "Second key point"
    - "Third key point"
  examples:
    - "Example from the article"
  prescription: "What the author recommends"
  tags:
    - technology
    - programming
  sentiment: cautious

discussion:
  sentiment: divided
  personas:
    - name: "The Pragmatist"
      description: "Experienced practitioners focused on applicability"
      stance: nuanced
      prevalence: dominant
      percentage: 35
      core_argument: "Good in theory, execution determines success"
      quotes:
        - text: "In my experience, this works for small teams."
          author: "dev_user"
        - text: "The key is knowing when to apply it."
          author: "tech_lead"
    - name: "The Skeptic"
      description: "Questions assumptions and demands evidence"
      stance: disagrees
      prevalence: significant
      percentage: 25
      core_argument: "The evidence is anecdotal"
      quotes:
        - text: "Has anyone measured this? Show me data."
          author: "data_person"
  tangents:
    - topic: "Related Technology"
      description: "Discussion veered into comparing alternatives"
  consensus_points:
    - "Current approaches have problems"
  open_questions:
    - "How does this scale?"
```

## JSON Schema Files

The embedded JSON schemas are available in the `schema/` directory:

- `schema/summary.schema.json` - Summary validation
- `schema/digest.schema.json` - Digest validation
- `schema/site.schema.json` - Site configuration
- `schema/source-links.schema.json` - Source links

### Programmatic Access

```go
import "github.com/grokify/socialpulse/schema"

// Get embedded schema
summarySchema := schema.SummarySchemaJSON()

// Parse as map
schemaMap, err := schema.SummarySchemaMap()
```
