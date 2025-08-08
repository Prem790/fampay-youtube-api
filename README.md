# FamPay YouTube API

A service that fetches latest videos from YouTube for a given search query and stores them in a database.

## Features

- Fetches YouTube videos periodically
- Stores video data in PostgreSQL
- Caches responses using Redis
- Provides REST API endpoints
- React dashboard for video management

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and fill in your YouTube API key
3. Run setup script:
   ```bash
   make setup
   ```

## Development

### Backend
```bash
make run  # Run with live reload
```

### Frontend
```bash
make frontend-dev
```

### Docker
```bash
make docker-up    # Start all services
make docker-down  # Stop all services
```

## Testing
```bash
make test
```

## API Endpoints

- `GET /api/videos` - List all videos
- `GET /api/videos/search` - Search videos

## Dashboard

The dashboard is available at `http://localhost:3000` when running in development mode.
