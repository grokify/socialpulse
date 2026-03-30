# socialpulse gh-deploy

Deploy the site to GitHub Pages.

## Synopsis

```bash
socialpulse gh-deploy [flags]
```

## Description

Builds the static site and deploys it to the `gh-pages` branch of your repository. This is a convenient wrapper around building and pushing to GitHub Pages.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `socialpulse.yaml` | Path to configuration file |

## Prerequisites

1. **Git repository** - The site must be in a git repository
2. **Remote configured** - `origin` remote must be set
3. **Push access** - You must have push access to the repository

## Examples

### Deploy to GitHub Pages

```bash
socialpulse gh-deploy
```

This:

1. Builds the site to a temporary directory
2. Creates/updates the `gh-pages` branch
3. Pushes the built site to the branch

### Typical workflow

```bash
# Make changes to content
vim content/summaries/new-article.yaml

# Validate and verify
socialpulse validate
socialpulse verify

# Preview locally
socialpulse serve

# Deploy
socialpulse gh-deploy
```

## GitHub Pages Configuration

After deploying, configure GitHub Pages in your repository settings:

1. Go to **Settings** → **Pages**
2. Under **Source**, select `gh-pages` branch
3. Select `/ (root)` as the folder
4. Click **Save**

Your site will be available at:

```
https://<username>.github.io/<repository>/
```

## Base URL Configuration

Ensure your `socialpulse.yaml` has the correct `base_url`:

```yaml
site:
  base_url: "https://username.github.io/my-site"
```

This is used for generating absolute URLs in the site.

## Custom Domain

To use a custom domain:

1. Add a `CNAME` file to your content directory
2. Configure DNS for your domain
3. Update `base_url` in `socialpulse.yaml`

```yaml
site:
  base_url: "https://discussions.example.com"
```

## Alternative: GitHub Actions

For automated deployment, use GitHub Actions:

```yaml
# .github/workflows/deploy.yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install SocialPulse
        run: go install github.com/grokify/socialpulse/cmd/socialpulse@latest

      - name: Build site
        run: socialpulse build -d ./public

      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

## Troubleshooting

### "Not a git repository"

Initialize git and add a remote:

```bash
git init
git remote add origin https://github.com/username/repo.git
```

### "Permission denied"

Ensure you have push access to the repository. For HTTPS, you may need to configure credentials. For SSH, ensure your key is added.

### "gh-pages branch not found"

The `gh-deploy` command creates the branch if it doesn't exist.
