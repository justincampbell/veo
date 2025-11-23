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

## Features

### Authentication & Configuration
- [x] Bearer token authentication via environment variable
- [x] Default club via `VEO_CLUB` environment variable
- [ ] OAuth login flow
- [ ] Configuration file support

### List Recordings
- [x] View all recordings with pagination support
  - [x] Default: first page (20 items)
  - [x] `--page N` for specific page
  - [x] `--all` to fetch all pages
- [x] Table output with title, slug, duration, and date
- [x] JSON output format (`--json`)
- [x] Dynamic terminal width detection
- [x] Display in local timezone

### Get Recording Details
- [x] Get match details by ID or `latest`
- [x] Show score, teams, and metadata
- [x] JSON output format (`--json`)
- [x] Generate share URLs
- [x] Generate highlights URLs

### Update Recordings
- [ ] Update match metadata (command exists but needs implementation)
- [ ] Update team sides/colors
- [ ] Set/update scores

## Contributing

This project tracks work via GitHub Issues. See the [issue list](https://github.com/justincampbell/veo/issues) for open tasks and feature requests.

## API Documentation

See [`docs/api.md`](docs/api.md) for detailed API documentation discovered through reverse engineering.
