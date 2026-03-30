# CLI Reference

SocialPulse provides a command-line interface for managing discussion summary sites.

## Usage

```bash
socialpulse [command] [flags]
```

## Available Commands

| Command | Description |
|---------|-------------|
| [`new`](new.md) | Create a new site scaffold |
| [`serve`](serve.md) | Run local development server |
| [`build`](build.md) | Build static HTML site |
| [`validate`](validate.md) | Validate content against schemas |
| [`verify`](verify.md) | Verify URLs and quotes against platform APIs |
| [`fetch-quotes`](fetch-quotes.md) | Fetch real quotes from discussions |
| [`gh-deploy`](gh-deploy.md) | Deploy to GitHub Pages |

## Global Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Help for any command |

## Configuration

Most commands accept a `-c, --config` flag to specify the configuration file:

```bash
socialpulse build -c my-config.yaml
```

The default configuration file is `socialpulse.yaml` in the current directory.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid input, build failure, etc.) |

## Examples

### Typical Workflow

```bash
# Create a new site
socialpulse new my-site
cd my-site

# Add content to content/summaries/

# Validate content
socialpulse validate

# Start development server
socialpulse serve

# Verify against platform APIs
socialpulse verify --quotes

# Build and deploy
socialpulse build
socialpulse gh-deploy
```

### CI/CD Pipeline

```bash
# In GitHub Actions or similar
socialpulse validate
socialpulse verify
socialpulse build -d ./public
# Deploy ./public to your hosting provider
```
