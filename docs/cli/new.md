# socialpulse new

Create a new SocialPulse site scaffold.

## Synopsis

```bash
socialpulse new <name> [flags]
```

## Description

Creates a new directory with the basic structure for a SocialPulse site, including configuration file and content directories.

## Arguments

| Argument | Description |
|----------|-------------|
| `name` | Name of the directory to create |

## Examples

### Create a new site

```bash
socialpulse new my-discussions
```

Creates:

```
my-discussions/
├── socialpulse.yaml
├── content/
│   ├── summaries/
│   └── digests/
│       ├── weekly/
│       ├── monthly/
│       ├── quarterly/
│       └── yearly/
└── site/
```

### Generated Configuration

The generated `socialpulse.yaml`:

```yaml
site:
  title: "my-discussions"
  description: "Discussion summaries powered by SocialPulse"
  base_url: "https://example.com"

theme:
  name: default

content:
  summaries_dir: content/summaries
  digests_dir: content/digests

build:
  output_dir: site
```

## Next Steps

After creating a site:

1. Edit `socialpulse.yaml` with your site details
2. Add summary files to `content/summaries/`
3. Run `socialpulse serve` to preview
4. Run `socialpulse build` to generate HTML
