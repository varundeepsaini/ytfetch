# YouTube Video Fetcher API

A Go-based API service that fetches and serves YouTube videos based on search queries in a paginated format.

## Features

- Asynchronous background fetching of YouTube videos
- Paginated API response sorted by publishing date
- Support for multiple YouTube API keys
- Efficient database storage with proper indexing
- Scalable architecture

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- YouTube Data API v3 key(s)

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── video.go
│   ├── repository/
│   │   └── video_repository.go
│   ├── service/
│   │   └── youtube_service.go
│   └── handlers/
│       └── video_handler.go
├── pkg/
│   └── youtube/
│       └── client.go
├── go.mod
├── go.sum
└── README.md
```

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd ytfetch
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=ytfetch
export YOUTUBE_API_KEYS=key1,key2,key3
export SEARCH_QUERY=your_search_query
```

4. Create MySQL database:
```sql
CREATE DATABASE ytfetch;
```

5. Run the server:
```bash
go run cmd/server/main.go
```

## API Endpoints

### GET /api/videos
Returns paginated list of videos sorted by publishing date.

Query Parameters:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)

Response:
```json
{
    "videos": [
        {
            "id": "string",
            "title": "string",
            "description": "string",
            "published_at": "datetime",
            "thumbnail_url": "string"
        }
    ],
    "total": "integer",
    "page": "integer",
    "limit": "integer"
}
```

## Configuration

The service can be configured using environment variables:

- `DB_HOST`: MySQL host (default: localhost)
- `DB_PORT`: MySQL port (default: 3306)
- `DB_USER`: MySQL user
- `DB_PASSWORD`: MySQL password
- `DB_NAME`: MySQL database name
- `YOUTUBE_API_KEYS`: Comma-separated list of YouTube API keys
- `SEARCH_QUERY`: Default search query for video fetching
- `FETCH_INTERVAL`: Interval in seconds for background fetching (default: 10)
