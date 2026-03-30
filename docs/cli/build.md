# socialpulse build

Build the static HTML site from content files.

## Synopsis

```bash
socialpulse build [flags]
```

## Description

Reads all summaries and digests from the content directories, then generates HTML pages using the configured theme. The output is a complete static site ready for deployment.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `socialpulse.yaml` | Path to configuration file |
| `-d, --output` | From config | Output directory (overrides config) |

## Examples

### Default build

```bash
socialpulse build
```

Outputs to the directory specified in `socialpulse.yaml` (default: `site/`).

### Custom output directory

```bash
socialpulse build -d ./public
```

### With custom config

```bash
socialpulse build -c production.yaml -d ./dist
```

## Generated Structure

```
site/
├── index.html              # Dashboard home page
├── articles/
│   ├── 2026-03-25-example.html
│   └── ...
├── digests/
│   ├── weekly-2026-03-23.html
│   └── ...
├── tags/
│   ├── technology.html
│   ├── programming.html
│   └── ...
├── css/
│   └── dashboard.css
└── js/
    └── charts.js
```

## Pages Generated

| Page Type | Description |
|-----------|-------------|
| **Index** | Dashboard with metrics, recent articles, tag cloud |
| **Article** | Individual summary with personas, quotes, analysis |
| **Tag** | Articles grouped by tag |
| **Digest** | Periodic aggregation (weekly, monthly, etc.) |

## Build Process

1. Load configuration from `socialpulse.yaml`
2. Scan content directories for `.yaml` and `.json` files
3. Parse and validate each content file
4. Generate HTML using embedded templates
5. Copy static assets (CSS, JS)
6. Write all files to output directory

## CI/CD Integration

```yaml
# GitHub Actions example
- name: Build site
  run: |
    go install github.com/grokify/socialpulse/cmd/socialpulse@latest
    socialpulse build -d ./public

- name: Deploy to Pages
  uses: peaceiris/actions-gh-pages@v3
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
    publish_dir: ./public
```
