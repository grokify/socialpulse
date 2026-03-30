# Getting Started

This guide walks you through creating your first SocialPulse site.

## Prerequisites

- Go 1.21 or later
- Git (for deployment)

## Installation

Install the SocialPulse CLI:

```bash
go install github.com/grokify/socialpulse/cmd/socialpulse@latest
```

Verify the installation:

```bash
socialpulse --help
```

## Create a New Site

Use the `new` command to scaffold a site:

```bash
socialpulse new my-discussions
cd my-discussions
```

This creates:

```
my-discussions/
├── socialpulse.yaml      # Site configuration
├── content/
│   ├── summaries/        # Article summaries (YAML/JSON)
│   └── digests/          # Periodic aggregations
└── site/                 # Generated output (gitignored)
```

## Configure Your Site

Edit `socialpulse.yaml`:

```yaml
site:
  title: "My Discussion Summaries"
  description: "Curated insights from technical discussions"
  base_url: "https://username.github.io/my-discussions"

theme:
  name: default

content:
  summaries_dir: content/summaries
  digests_dir: content/digests

build:
  output_dir: site
```

## Add Your First Summary

Create a summary file in `content/summaries/`:

```yaml
# content/summaries/2026-03-25-example.yaml
schema_version: "1.0.0"

meta:
  article_url: "https://example.com/great-article"
  discussion_url: "https://news.ycombinator.com/item?id=12345678"
  platform: hackernews
  article_date: 2026-03-25T00:00:00Z
  summary_date: 2026-03-26T00:00:00Z
  article_author: "Jane Doe"
  discussion_comment_count: 150
  item_id: "12345678"

article:
  title: "A Great Article About Technology"
  thesis: "The main argument presented by the author, summarized in 200-400 characters for quick comprehension."
  key_arguments:
    - "First major point with supporting evidence"
    - "Second point addressing common objections"
    - "Third point with practical implications"
  examples:
    - "Concrete example cited in the article"
  prescription: "The author's recommended action or solution"
  tags:
    - technology
    - software-engineering
  sentiment: cautious

discussion:
  sentiment: divided
  personas:
    - name: "The Pragmatist"
      description: "Experienced practitioners focused on real-world applicability"
      stance: nuanced
      prevalence: dominant
      percentage: 35
      core_argument: "Good idea in theory, but implementation details determine success"
      quotes:
        - text: "I've tried this approach and it works well for small teams."
          author: "experienced_dev"
        - text: "The key is knowing when to apply these principles."
          author: "tech_lead_42"

    - name: "The Skeptic"
      description: "Questions assumptions and demands evidence"
      stance: disagrees
      prevalence: significant
      percentage: 25
      core_argument: "The evidence presented is anecdotal at best"
      quotes:
        - text: "Has anyone actually measured the impact? I'd like to see data."
          author: "data_driven"

  consensus_points:
    - "Most agree the current approach has problems"
    - "Everyone values developer productivity"

  open_questions:
    - "How does this scale to larger organizations?"
    - "What are the long-term maintenance implications?"
```

## Validate Your Content

Check that your summary is valid:

```bash
socialpulse validate
```

You should see:

```
Validating content files...
  2026-03-25-example.yaml: OK
Summary: 1 valid, 0 invalid
```

## Start the Development Server

Run the development server:

```bash
socialpulse serve
```

Open [http://127.0.0.1:8000](http://127.0.0.1:8000) in your browser.

The server watches for changes and automatically rebuilds when you edit content files.

## Build for Production

Generate the static site:

```bash
socialpulse build
```

The output is in the `site/` directory, ready for deployment.

## Deploy to GitHub Pages

Deploy directly from the CLI:

```bash
socialpulse gh-deploy
```

This builds the site and pushes to the `gh-pages` branch.

## Next Steps

- [CLI Reference](cli/index.md) - Learn all available commands
- [Content Authoring](content/authoring.md) - Write effective summaries
- [Writing Personas](content/personas.md) - Create insightful persona analysis
- [Verification](verification.md) - Verify quotes against platform APIs
