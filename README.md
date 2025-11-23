# Veo CLI

Command-line interface for the Veo sports camera API.

## Installation

```bash
go install github.com/justincampbell/veo/cmd/veo@latest
```

Or build from source:

```bash
make build
```

## Usage

### Authentication

Set your Veo Bearer token as an environment variable:

```bash
export VEO_TOKEN="your-token-here"
```

You can extract this from your browser's DevTools while logged into app.veo.co (Network tab → any API request → Authorization header).

Optionally, set a default club to avoid using the `--club` flag:

```bash
export VEO_CLUB="your-club-slug"
```

### List Recordings

```bash
# List first page (20 recordings)
veo list --club your-club-slug

# Or use VEO_CLUB environment variable
veo list

# List specific page
veo list --page 2

# List all recordings (fetches all pages)
veo list --all
```

## Development

```bash
# Run tests
make test

# Build binary
make build

# Run without building
make run
```

## Current Features

- ✅ **List recordings** - View all recordings with pagination support
  - Default: first page (20 items)
  - `--page N` for specific page
  - `--all` to fetch all pages
- ✅ Bearer token authentication via environment variable
- ✅ Table output with title, slug, duration, and date

## Roadmap

See [GitHub Issues](https://github.com/justincampbell/veo/issues) for planned features and enhancements.

**Upcoming features:**
- Get match details
- Update match metadata
- Update team sides/colors
- Get highlights
- Generate sharing URLs with kickoff timestamp
- OAuth login flow
- Configuration file support
- JSON output format

## Contributing

This project tracks work via GitHub Issues. See the [issue list](https://github.com/justincampbell/veo/issues) for open tasks and feature requests.

## API Documentation

See [`docs/api.md`](docs/api.md) for detailed API documentation discovered through reverse engineering.
