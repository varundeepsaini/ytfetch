# YouTube Video Fetcher

A simple app that gets YouTube videos and shows them in a nice dashboard.

## What it does

- Gets YouTube videos in the background
- Shows videos in a grid or list view
- Lets you search and filter videos
- Works on all screen sizes
- Has dark mode

## What you need

- Go 1.21 or newer
- MySQL 8.0 or newer
- Node.js 14 or newer
- A YouTube API key

## How to set up

1. Get the code:
```bash
git clone <repository-url>
cd ytfetch
```

2. Install Go packages:
```bash
go mod download
```

3. Set up your settings in `.env`:
```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ytfetch
YOUTUBE_API_KEYS=your_api_key_comma_seperated
SEARCH_QUERY=your_search_term
```

4. Create the database:
```bash
CREATE DATABASE ytfetch;
```

5. Start the backend:
```bash
go run cmd/server/main.go // use run.sh to run the server with an empty db
```

6. Start the frontend:
```bash
cd web
npm install
npm start
```

## How to use

1. Open http://localhost:3000 in your browser
2. Use the search box to find videos
3. Filter by date or channel
4. Switch between grid and list views
5. Click "Load More" to see more videos

## API Endpoints

### GET /api/videos
Returns paginated list of videos sorted by publishing date.

Options:
- `limit`: How many videos to show (default: 10)
- `cursor`: Where to start from (for loading more)

Example response:
```json
{
    "videos": [
        {
            "id": "video_id",
            "title": "Video Title",
            "description": "Video Description",
            "published_at": "2024-02-20T10:00:00Z",
            "thumbnail_url": "https://..."
        }
    ],
    "next_cursor": "2024-02-19T10:00:00Z",
    "has_more": true
}
```
