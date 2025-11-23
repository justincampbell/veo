# Veo API Documentation

This document describes the Veo API endpoints discovered through reverse engineering.

## Authentication

The API uses two authentication mechanisms:

1. **Bearer Token**: Passed in `Authorization` header
   ```
   Authorization: Bearer <token>
   ```

2. **CSRF Token**: Passed in `x-csrftoken` header (obtained from `csrftoken` cookie)
   ```
   x-csrftoken: <csrf-token>
   ```

## Base URL

```
https://app.veo.co/api/app
```

## Endpoints

### List Recordings

Lists all recordings/matches for a club.

**Request:**
```
GET /clubs/{club-slug}/recordings/?filter=own&fields=<field-list>
```

**Query Parameters:**
- `filter`: Filter type (e.g., `own`)
- `fields`: Comma-separated list of fields to return

**Useful Fields:**
- `camera`, `created`, `device`, `duration`, `identifier`, `slug`
- `title`, `team`, `url`, `thumbnail`
- `reel_url` (full game video URL)
- `processing_status`, `privacy`, `permissions`

**Response:** Array of recording objects

**Example Response:**
```json
[
  {
    "camera": "vc3-xxxxx",
    "created": "2025-11-22T20:09:36.464909+01:00",
    "duration": 3259,
    "identifier": "uuid",
    "slug": "20251116-match-team-name-hash",
    "title": "Match - Team Name",
    "url": "/matches/20251116-match-team-name-hash/",
    "thumbnail": "https://c.veocdn.com/.../thumbnail.jpg",
    "reel_url": "https://c.veocdn.com/.../reel/video.mp4"
  }
]
```

### Get Match Details

Retrieves detailed information about a specific match.

**Request:**
```
GET /matches/{match-slug-or-id}/
```

**Response:** Match object with full details

**Example Response:**
```json
{
  "id": "uuid",
  "identifier": "uuid",
  "slug": "20251116-match-team-name-hash",
  "title": "Match Title",
  "created": "2025-11-22T21:04:03.010289+01:00",
  "start": "2025-11-16T15:58:21.251456+01:00",
  "end": "2025-11-16T16:55:11.380000+01:00",
  "duration": 3410,
  "type": "match",
  "own_team_home_or_away": "home",
  "opponent_team_name": "Opponent Name",
  "opponent_club_name": "Opponent Club",
  "opponent_team_color": "blue",
  "opponent_short_name": "OPP",
  "own_team_color": "red",
  "team": "team-id",
  "info": {
    "stats": {
      "score": {
        "own": 2,
        "opponent": 1
      },
      "score_aggregated": {
        "own": 2,
        "opponent": 1
      }
    },
    "age_group": "U11"
  },
  "permissions": {...}
}
```

### Update Match

Updates match metadata.

**Request:**
```
PATCH /matches/{match-id}/
Content-Type: application/json
```

**Request Body Examples:**

Update basic info:
```json
{
  "title": "New Title - Opponent Name",
  "type": "tournament",
  "own_team_home_or_away": "away",
  "start": "2025-11-16T12:00:00"
}
```

Update team details (sides):
```json
{
  "opponent_club_name": "Opponent Club",
  "opponent_team_name": "Opponent Team",
  "opponent_team_color": "yellow",
  "opponent_short_name": "OPP",
  "own_team_color": "orange",
  "own_team_formation": "4-3-1",
  "opponent_team_formation": null,
  "team": "team-uuid"
}
```

**Response:** Updated match object

### Get Highlights

Retrieves AI-generated highlights for a match.

**Request:**
```
GET /matches/{match-slug}/highlights/?include_ai=true&fields=<field-list>
```

**Query Parameters:**
- `include_ai`: Include AI-generated highlights (boolean)
- `fields`: Comma-separated list of fields

**Useful Fields:**
- `id`, `created`, `duration`, `start`, `thumbnail`
- `is_ai_generated`, `ai_resolution`
- `videos` (contains URLs)
- `tags`, `involved_players`

**Response:** Array of highlight objects

### Get Match Videos

Retrieves video information for a match.

**Request:**
```
GET /matches/{match-slug}/videos/
```

**Response:** Array of video objects with URLs

### Get Match Periods

Retrieves period/half information (kickoff timestamps).

**Request:**
```
GET /matches/{match-slug}/periods/
```

**Response:** Array of period objects with timestamps

## Field Types

### Match Type
- `match` - Regular match
- `tournament` - Tournament game
- `training` - Training session
- `scrimmage` - Practice match

### Home/Away
- `home` - Home team
- `away` - Away team

## Notes

- Most endpoints support field selection via `fields` query parameter
- UUIDs are used as primary identifiers
- Slugs are human-readable IDs used in URLs
- The API returns ISO 8601 timestamps with timezone
