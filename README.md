# SocialPulse

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/socialpulse/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/socialpulse/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/socialpulse/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/socialpulse/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/socialpulse/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/socialpulse/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/socialpulse
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/socialpulse
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/socialpulse
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/socialpulse
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fsocialpulse
 [loc-svg]: https://tokei.rs/b1/github/grokify/socialpulse
 [repo-url]: https://github.com/grokify/socialpulse
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/socialpulse/blob/master/LICENSE

SocialPulse is an MkDocs-like static site generator that transforms article discussion summaries into a searchable, browsable website. It reads YAML/JSON summary files and generates a data dashboard-style website with persona cards, sentiment analysis, and trend visualizations.

## Features

- ⚡ **Static Site Generation** - Build fast, searchable HTML from YAML/JSON content
- 📊 **Data Dashboard Theme** - Information-dense layouts with ECharts visualizations
- 👥 **Persona Analysis** - Categorize discussion participants into archetypal viewpoints
- 📈 **Sentiment Tracking** - Track article and discussion sentiment over time
- 🔌 **Platform Support** - HackerNews integration with Reddit planned
- ✅ **Anti-Hallucination Verification** - Verify quotes and URLs against platform APIs
- 🔄 **Live Development Server** - File watching with automatic rebuild

## Installation

```bash
go install github.com/grokify/socialpulse/cmd/socialpulse@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/socialpulse.git
cd socialpulse
go build -o socialpulse ./cmd/socialpulse
```

## Quick Start

### Create a New Site

```bash
socialpulse new my-site
cd my-site
```

This creates:

```
my-site/
├── socialpulse.yaml      # Site configuration
├── content/
│   ├── summaries/        # Article summaries (YAML/JSON)
│   └── digests/          # Periodic aggregations
└── site/                 # Generated output
```

### Add Content

Create a summary file in `content/summaries/`:

```yaml
# content/summaries/2026-03-25-example.yaml
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

article:
  title: "Example Article Title"
  thesis: "The main argument of the article in 200-400 characters."
  key_arguments:
    - "First key point made by the author"
    - "Second key point with supporting evidence"
    - "Third key point addressing counterarguments"
  tags:
    - technology
    - programming
  sentiment: cautious

discussion:
  sentiment: divided
  personas:
    - name: "The Pragmatist"
      description: "Experienced practitioners focused on real-world applicability"
      stance: nuanced
      prevalence: dominant
      percentage: 35
      core_argument: "The idea has merit but implementation details matter most"
      quotes:
        - text: "In my experience, this works well for small teams but breaks down at scale."
          author: "username123"
    - name: "The Skeptic"
      description: "Questions fundamental assumptions and asks for evidence"
      stance: disagrees
      prevalence: significant
      percentage: 25
      core_argument: "The premises are flawed and the evidence is anecdotal"
      quotes:
        - text: "Has anyone actually measured this? I'd like to see data."
          author: "datadriven"
  consensus_points:
    - "Most agree the status quo has problems"
  open_questions:
    - "How does this scale to larger organizations?"
```

### Build and Serve

```bash
# Development server with live reload
socialpulse serve

# Build static site
socialpulse build -d ./site
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `socialpulse new <name>` | Create a new site scaffold |
| `socialpulse serve` | Run development server (default: localhost:8000) |
| `socialpulse build` | Build static HTML site |
| `socialpulse validate` | Validate content against JSON schemas |
| `socialpulse verify` | Verify URLs and quotes against platform APIs |
| `socialpulse fetch-quotes` | Fetch real quotes from discussions |
| `socialpulse gh-deploy` | Deploy to GitHub Pages |

### Command Options

```bash
# Serve with custom port
socialpulse serve -p 3000 -H 0.0.0.0

# Build to custom directory
socialpulse build -d ./public

# Verify and auto-fix discrepancies
socialpulse verify --fix --quotes --verbose

# Fetch quotes with dry-run
socialpulse fetch-quotes --dry-run --verbose
```

## Configuration

Site configuration in `socialpulse.yaml`:

```yaml
site:
  title: "My Discussion Summaries"
  description: "Curated summaries of technical discussions"
  base_url: "https://example.com"

theme:
  name: default

content:
  summaries_dir: content/summaries
  digests_dir: content/digests

build:
  output_dir: site
```

## Data Model

### Summary Structure

| Field | Description |
|-------|-------------|
| `meta` | Metadata including URLs, dates, platform, comment count |
| `article` | Title, thesis, key arguments, examples, tags, sentiment |
| `discussion` | Personas, consensus points, open questions, tangents |

### Persona Fields

| Field | Description |
|-------|-------------|
| `name` | Descriptive name (e.g., "The Skeptic", "The Pragmatist") |
| `description` | Brief description of the viewpoint (50-100 chars) |
| `stance` | `agrees`, `disagrees`, `nuanced`, or `tangential` |
| `prevalence` | `dominant`, `significant`, or `minority` |
| `percentage` | Estimated percentage of participants (0-100) |
| `quotes` | Verbatim quotes with author attribution |
| `core_argument` | Distilled argument (100-150 chars) |

### Sentiment Values

**Article sentiment:** `optimistic`, `cautious`, `pessimistic`, `neutral`, `provocative`

**Discussion sentiment:** `supportive`, `divided`, `critical`, `exploratory`

## Verification

SocialPulse includes tools to verify content against platform APIs:

```bash
# Verify URLs and comment counts
socialpulse verify

# Also verify persona quotes exist in actual comments
socialpulse verify --quotes

# Auto-fix discrepancies
socialpulse verify --fix
```

The `fetch-quotes` command fetches real quotes from discussions and matches them to personas using keyword analysis:

```bash
socialpulse fetch-quotes --verbose
```

## Packages

SocialPulse provides reusable Go packages:

| Package | Description |
|---------|-------------|
| `types` | Core data structures for summaries, digests, and sites |
| `schema` | Embedded JSON schemas with validation |
| `textutil` | Text processing utilities (HTML stripping, keyword extraction) |

### Using as a Library

```go
import (
    "github.com/grokify/socialpulse/types"
    "github.com/grokify/socialpulse/schema"
    "github.com/grokify/socialpulse/textutil"
)

// Create a summary
summary := types.Summary{
    SchemaVersion: "1.0.0",
    Meta: types.Meta{
        Platform: types.PlatformHackerNews,
        // ...
    },
    // ...
}

// Strip HTML from comment text
plainText := textutil.StripHTMLAndQuotes(htmlComment)

// Extract keywords for matching
keywords := textutil.ExtractKeywords(text, 3)
```

## Contributing

Contributions are welcome! Please ensure:

1. Tests pass: `go test ./...`
2. Linting passes: `golangci-lint run`
3. Commits follow [Conventional Commits](https://www.conventionalcommits.org/)

## License

MIT License - see [LICENSE](LICENSE) for details.
