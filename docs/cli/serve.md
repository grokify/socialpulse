# socialpulse serve

Run a local development server with live reload.

## Synopsis

```bash
socialpulse serve [flags]
```

## Description

Starts a local development server that watches for changes to content files and automatically rebuilds the site when changes are detected.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `socialpulse.yaml` | Path to configuration file |
| `-p, --port` | `8000` | Port to serve on |
| `-H, --host` | `127.0.0.1` | Host to bind to |
| `--no-watch` | `false` | Disable file watching |

## Examples

### Default usage

```bash
socialpulse serve
```

Opens at [http://127.0.0.1:8000](http://127.0.0.1:8000)

### Custom port

```bash
socialpulse serve -p 3000
```

### Bind to all interfaces

```bash
socialpulse serve -H 0.0.0.0
```

Useful for accessing from other devices on your network.

### Disable file watching

```bash
socialpulse serve --no-watch
```

Serves the site without rebuilding on changes.

## Behavior

### File Watching

The server watches these directories for changes:

- `content/summaries/` - Article summaries
- `content/digests/` - Periodic digests

When a `.yaml` or `.json` file changes, the server:

1. Rebuilds the affected pages
2. Logs the rebuild to the console

!!! note
    Changes to `socialpulse.yaml` require restarting the server.

### Static Assets

Static assets (CSS, JS, images) are served from the theme's embedded assets. Custom theme overrides are not yet supported.

## Troubleshooting

### Port already in use

```
Error: listen tcp 127.0.0.1:8000: bind: address already in use
```

Use a different port:

```bash
socialpulse serve -p 8080
```

### Permission denied on port 80

Ports below 1024 require root privileges. Use a higher port number.
