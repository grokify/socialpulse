# SocialPulse

**Static site generator for social discussion summaries**

SocialPulse transforms article discussion summaries into a searchable, browsable website. It reads YAML/JSON summary files and generates a data dashboard-style website with persona cards, sentiment analysis, and trend visualizations.

## Why SocialPulse?

When articles spark discussions on HackerNews, Reddit, or other platforms, valuable insights emerge from the conversation. SocialPulse helps you:

- **Capture discussion insights** - Summarize key viewpoints and arguments
- **Identify personas** - Categorize commenters into archetypal perspectives
- **Track sentiment** - Monitor how discussions evolve over time
- **Verify accuracy** - Check quotes and URLs against platform APIs
- **Publish easily** - Generate static sites for GitHub Pages or any host

## Features

| Feature | Description |
|---------|-------------|
| **Static Generation** | Fast, searchable HTML from YAML/JSON content |
| **Dashboard Theme** | Information-dense layouts with ECharts visualizations |
| **Persona Analysis** | Categorize participants into archetypal viewpoints |
| **Sentiment Tracking** | Article and discussion sentiment over time |
| **Platform Support** | HackerNews integration (Reddit planned) |
| **Verification** | Anti-hallucination checks against platform APIs |
| **Live Reload** | Development server with file watching |

## Quick Example

```yaml
# content/summaries/2026-03-25-ai-coding.yaml
schema_version: "1.0.0"

meta:
  platform: hackernews
  item_id: "47517539"
  discussion_comment_count: 458

article:
  title: "Thoughts on slowing the fuck down"
  thesis: "AI coding agents encourage developers to move fast without understanding..."
  sentiment: cautious

discussion:
  sentiment: divided
  personas:
    - name: "The Pragmatic Senior"
      stance: nuanced
      prevalence: dominant
      percentage: 35
      core_argument: "AI is powerful when used deliberately, dangerous as a crutch"
```

## Getting Started

Ready to create your first site?

[Get Started :material-arrow-right:](getting-started.md){ .md-button .md-button--primary }

## Installation

=== "Go Install"

    ```bash
    go install github.com/grokify/socialpulse/cmd/socialpulse@latest
    ```

=== "From Source"

    ```bash
    git clone https://github.com/grokify/socialpulse.git
    cd socialpulse
    go build -o socialpulse ./cmd/socialpulse
    ```

## Project Status

SocialPulse is under active development. Current status:

- [x] Core site generation
- [x] HackerNews integration
- [x] Persona and sentiment analysis
- [x] Quote verification
- [x] GitHub Pages deployment
- [ ] Reddit integration
- [ ] Custom themes
- [ ] Digest generation UI
