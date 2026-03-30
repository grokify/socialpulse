# Deployment

SocialPulse generates static HTML that can be hosted anywhere. This guide covers common deployment options.

## GitHub Pages

The easiest deployment option for GitHub-hosted repositories.

### Using gh-deploy Command

```bash
socialpulse gh-deploy
```

This:

1. Builds the site
2. Commits to the `gh-pages` branch
3. Pushes to origin

### GitHub Actions (Recommended)

Automate deployment on every push:

```yaml
# .github/workflows/deploy.yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install SocialPulse
        run: go install github.com/grokify/socialpulse/cmd/socialpulse@latest

      - name: Validate content
        run: socialpulse validate

      - name: Build site
        run: socialpulse build -d ./public

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

### Configuration

Update `socialpulse.yaml` with your GitHub Pages URL:

```yaml
site:
  base_url: "https://username.github.io/repository"
```

### Custom Domain

1. Add a `CNAME` file to your repository root:
   ```
   discussions.example.com
   ```

2. Configure DNS:
   - Add a CNAME record pointing to `username.github.io`
   - Or add A records for GitHub's IPs

3. Update configuration:
   ```yaml
   site:
     base_url: "https://discussions.example.com"
   ```

## Cloudflare Pages

### Setup

1. Connect your GitHub repository to Cloudflare Pages
2. Configure build settings:
   - **Build command:** `go install github.com/grokify/socialpulse/cmd/socialpulse@latest && socialpulse build -d ./public`
   - **Build output directory:** `public`
   - **Root directory:** `/`

3. Add environment variable:
   - `GO_VERSION`: `1.21`

### Configuration

```yaml
site:
  base_url: "https://your-project.pages.dev"
```

## Netlify

### netlify.toml

```toml
[build]
  command = "go install github.com/grokify/socialpulse/cmd/socialpulse@latest && socialpulse build -d ./public"
  publish = "public"

[build.environment]
  GO_VERSION = "1.21"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

### Configuration

```yaml
site:
  base_url: "https://your-site.netlify.app"
```

## Vercel

### vercel.json

```json
{
  "buildCommand": "go install github.com/grokify/socialpulse/cmd/socialpulse@latest && socialpulse build -d ./public",
  "outputDirectory": "public",
  "framework": null
}
```

## Self-Hosted (Nginx)

### Build Locally

```bash
socialpulse build -d ./dist
```

### Nginx Configuration

```nginx
server {
    listen 80;
    server_name discussions.example.com;
    root /var/www/discussions;
    index index.html;

    location / {
        try_files $uri $uri/ =404;
    }

    # Cache static assets
    location ~* \.(css|js|png|jpg|gif|ico)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Gzip compression
    gzip on;
    gzip_types text/html text/css application/javascript;
}
```

### Deploy Script

```bash
#!/bin/bash
set -e

# Build
socialpulse build -d ./dist

# Deploy
rsync -avz --delete ./dist/ user@server:/var/www/discussions/

echo "Deployed successfully"
```

## S3 + CloudFront

### Build and Upload

```bash
# Build
socialpulse build -d ./dist

# Upload to S3
aws s3 sync ./dist s3://your-bucket-name --delete

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id YOUR_DIST_ID \
  --paths "/*"
```

### S3 Bucket Policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    }
  ]
}
```

## Docker

### Dockerfile

```dockerfile
FROM golang:1.21 AS builder

WORKDIR /app
RUN go install github.com/grokify/socialpulse/cmd/socialpulse@latest

COPY . .
RUN socialpulse build -d ./public

FROM nginx:alpine
COPY --from=builder /app/public /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Build and Run

```bash
docker build -t my-discussions .
docker run -p 8080:80 my-discussions
```

## Deployment Checklist

Before deploying:

- [ ] `socialpulse validate` passes
- [ ] `socialpulse verify` passes
- [ ] `base_url` is set correctly in `socialpulse.yaml`
- [ ] All content files are committed
- [ ] Build completes without errors
- [ ] Test locally with `socialpulse serve`

## Continuous Deployment

For any CI/CD platform:

```bash
# Install
go install github.com/grokify/socialpulse/cmd/socialpulse@latest

# Validate
socialpulse validate
socialpulse verify

# Build
socialpulse build -d ./public

# Deploy (platform-specific)
# ...
```
