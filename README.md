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

```bash
# List videos
veo list

# Update video metadata
veo update <video-id>
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

## Project Goals

This CLI aims to provide the following functionality:
- List videos/recordings/matches (using API's terminology)
- Update video metadata: title, opponent, home/away, match type
- Update sides (for automatic goal detection and score tracking)
- Get sharing URL with timestamp from first kickoff
- Get highlights URL (automatically generated)

## Development Plan

### Phase 1: API Discovery (Current)
1. Export HAR file from browser while interacting with veo.co
2. Analyze HAR to identify:
   - Base URL and endpoints
   - Authentication mechanism
   - Request/response structures
   - Field names and data types

### Phase 2: Implementation
1. Update `internal/api/client.go` with actual endpoints
2. Update `internal/models/models.go` with real data structures
3. Implement commands in `internal/commands/`
4. Add integration tests

### Phase 3: Polish
1. Add configuration file support
2. Improve error messages
3. Add output formatting (JSON, table, etc.)
4. Documentation
