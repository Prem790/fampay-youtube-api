# FamPay YouTube API - Full Stack Application

A scalable YouTube video fetching and search application built with **Go**, **React**, **MongoDB**, and **Docker**. This project continuously fetches the latest videos from YouTube, stores them in a database, and provides a beautiful dashboard to search and view videos.

## 🎯 Project Overview

### Features
- **🔄 Real-time YouTube Integration**: Continuously fetches latest videos every 10 seconds
- **🔍 Dual Search Modes**: Search stored videos or live YouTube search
- **📱 Responsive Dashboard**: Modern React UI with Tailwind CSS
- **⚡ Smart Pagination**: Efficient pagination with sorting options
- **🔑 Multiple API Keys**: Auto-rotation when quota exhausted
- **🐳 Fully Dockerized**: Complete containerized setup
- **📊 Performance Optimized**: MongoDB indexes, Redis caching

### Tech Stack
- **Backend**: Go (Gin framework), MongoDB, Redis
- **Frontend**: React, Vite, Tailwind CSS, React Query
- **Database**: MongoDB Atlas (cloud)
- **Caching**: Redis
- **Infrastructure**: Docker, Docker Compose

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │   Go Backend    │    │   MongoDB       │
│   (Port 5173)   │◄──►│   (Port 8080)   │◄──►│   (Atlas)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Redis Cache   │
                       │   (Port 6380)   │
                       └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- **Docker & Docker Compose** (for containerized setup)
- **Go 1.21+** (for local development)
- **Node.js 18+** (for local frontend development)
- **YouTube Data API v3 Keys** ([Get them here](#getting-youtube-api-keys))

## 📋 Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/Prem790/fampay-youtube-api.git
cd fampay-youtube-api
```

### 2. Get YouTube API Keys

#### How to get YouTube Data API v3 Keys:
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable **YouTube Data API v3**:
   - Go to "APIs & Services" → "Library"
   - Search for "YouTube Data API v3"
   - Click "Enable"
4. Create API Key:
   - Go to "APIs & Services" → "Credentials"
   - Click "Create Credentials" → "API Key"
   - Copy the generated key
5. **Important**: Repeat this process to create **multiple API keys** from different Google accounts

#### Why Multiple API Keys?
YouTube API has a daily quota limit of 10,000 units per key. Each search uses ~100 units, so you'll get ~100 searches per day per key. Having multiple keys ensures uninterrupted service.

### 3. Environment Configuration

#### Create Backend Environment File
Create `.env` in the root directory:

```bash
# .env file (root directory)

# Server Configuration
PORT=8080
HOST=localhost
GIN_MODE=debug

# MongoDB Configuration - Replace with your MongoDB Atlas connection
MONGODB_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@YOUR_CLUSTER.mongodb.net/fampay_youtube?retryWrites=true&w=majority
MONGODB_DATABASE=fampay_youtube

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# YouTube API Configuration - Add YOUR API keys here (comma-separated)
YOUTUBE_API_KEYS=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXX,AIzaSyYYYYYYYYYYYYYYYYYYYYYYYYYY,AIzaSyZZZZZZZZZZZZZZZZZZZZZZZZZZ

# Search Configuration
YOUTUBE_SEARCH_QUERIES=cricket,football,technology,music,gaming,news,travel,cooking,sports,entertainment
FETCH_INTERVAL=10
MAX_RESULTS_PER_QUERY=50
REGION_CODE=IN
RELEVANCE_LANGUAGE=en
```

#### Create Frontend Environment File
Create `web/.env` in the web directory:

```bash
# web/.env file (frontend directory)

VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=FamPay YouTube Dashboard
```

#### MongoDB Setup Options

**Option 1: MongoDB Atlas (Recommended)**
1. Go to [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Create a free cluster
3. Create a database user
4. Get your connection string
5. Replace `MONGODB_URI` in `.env` with your connection string

**Option 2: Local MongoDB**
```bash
# If using local MongoDB, update .env:
MONGODB_URI=mongodb://localhost:27017
```

## 🐳 Running with Docker (Recommended)

### Complete Stack with Docker Compose
```bash
# 1. Start all services
docker-compose up --build

# 2. Access the application
# - Frontend: http://localhost:5173
# - Backend API: http://localhost:8080
# - Redis: localhost:6380
```

### Services Overview
- **Backend**: Go API server with YouTube integration
- **Frontend**: React dashboard with hot reload
- **Redis**: Caching layer (separate port to avoid conflicts)

## 💻 Running Locally

### Option 1: Full Local Setup

#### Backend
```bash
# 1. Install Go dependencies
go mod download

# 2. Start Redis (if not using Docker)
# Using Docker:
docker run -d -p 6379:6379 redis:7-alpine

# Using local Redis:
# MacOS: brew install redis && brew services start redis
# Ubuntu: sudo apt install redis-server && sudo systemctl start redis

# 3. Start the backend
go run cmd/server/main.go
```

#### Frontend
```bash
# 1. Install frontend dependencies
cd web
npm install

# 2. Start the development server
npm run dev

# Frontend will be available at http://localhost:5173
```

### Option 2: Hybrid Setup (Backend in Docker, Frontend Local)
```bash
# 1. Start only backend services
docker-compose up backend redis

# 2. Run frontend locally
cd web
npm install
npm run dev
```

## 📡 API Endpoints

### Base URL: `http://localhost:8080`

| Endpoint | Method | Description |
|----------|---------|-------------|
| `/health` | GET | System health check |
| `/api/videos` | GET | Get stored videos (paginated) |
| `/api/videos/search` | GET | Search stored videos |
| `/api/videos/youtube-search` | GET | Live YouTube search |

### Example API Calls
```bash
# Health check
curl "http://localhost:8080/health"

# Get videos (paginated)
curl "http://localhost:8080/api/videos?page=1&page_size=10&sort=latest"

# Search stored videos
curl "http://localhost:8080/api/videos/search?q=cricket&page=1&page_size=5"

# Live YouTube search
curl "http://localhost:8080/api/videos/youtube-search?q=programming&page=1&page_size=5"
```

## 🔧 Configuration Options

### Backend Configuration (.env)
| Variable | Description | Example |
|----------|-------------|---------|
| `YOUTUBE_API_KEYS` | Comma-separated API keys | `key1,key2,key3` |
| `YOUTUBE_SEARCH_QUERIES` | Search terms for background fetching | `cricket,football,tech` |
| `FETCH_INTERVAL` | Seconds between API calls | `10` |
| `MAX_RESULTS_PER_QUERY` | Videos per API call | `50` |
| `REGION_CODE` | Country code for regional content | `IN` |
| `RELEVANCE_LANGUAGE` | Language preference | `en` |

### Frontend Configuration (web/.env)
| Variable | Description | Example |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_APP_TITLE` | Application title | `FamPay Dashboard` |

## 🗂️ Project Structure

```
fampay-youtube-api/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/               # HTTP request handlers
│   │   ├── middleware/             # CORS, logging middleware
│   │   └── routes/                 # Route definitions
│   ├── config/                     # Configuration loading
│   ├── models/                     # Data models
│   ├── repository/                 # Database operations
│   ├── services/                   # Business logic
│   └── worker/                     # Background YouTube fetcher
├── pkg/
│   ├── database/                   # MongoDB connection
│   └── redis/                      # Redis client
├── web/                            # React frontend
│   ├── src/
│   │   ├── components/             # React components
│   │   ├── services/               # API calls
│   │   └── utils/                  # Helper functions
│   ├── package.json
│   └── vite.config.js
├── docker/
│   ├── Dockerfile.backend          # Backend container
│   └── Dockerfile.frontend         # Frontend container
├── docker-compose.yml              # Multi-service setup
├── .env                            # Backend environment
├── go.mod                          # Go dependencies
└── README.md
```

## 🎮 Usage Guide

### Dashboard Features
1. **🏠 Home Page**: View latest stored videos
2. **🔍 Search**: Two modes available
   - **Database Search**: Search through stored videos (fast, no API usage)
   - **Live YouTube Search**: Real-time YouTube search (uses API quota)
3. **📊 Sorting**: Latest, oldest, title, channel name
4. **📄 Pagination**: Navigate through video results
5. **▶️ Video Links**: Click to open videos on YouTube

### Search Tips
- Use simple keywords like "cricket", "cooking", "technology"
- Database search supports partial matching: "tea how" matches "How to make tea?"
- Live search gets fresh results but uses API quota

## 🚨 Troubleshooting

### Common Issues

#### 1. API Quota Exhausted
```bash
# Check API status
curl "http://localhost:8080/health" | jq '.api_status'

# Solution: Add more API keys to .env
YOUTUBE_API_KEYS=key1,key2,key3,key4,key5
```

#### 2. Port Conflicts
```bash
# If ports 8080, 5173, or 6380 are in use:
# Update docker-compose.yml ports section:
ports:
  - "8081:8080"  # Change external port
```

#### 3. MongoDB Connection Issues
```bash
# Verify MongoDB Atlas connection string
# Check whitelist IP addresses in Atlas
# Ensure correct username/password
```

#### 4. Docker Build Issues
```bash
# Clean Docker cache
docker system prune -a

# Rebuild with no cache
docker-compose build --no-cache
```

#### 5. Frontend Not Loading
```bash
# Check if backend is running
curl http://localhost:8080/health

# Verify frontend .env file exists in web/ directory
# Check VITE_API_BASE_URL points to correct backend URL
```

## 🔐 Security Notes

- Never commit API keys to version control
- Use different API keys for development/production
- Implement rate limiting in production
- Use environment variables for all sensitive data

## 📈 Performance Tips

1. **API Efficiency**: System rotates through search queries (10 searches, 1 per query rotation)
2. **Database Indexes**: Optimized for common search patterns
3. **Caching**: Redis caches frequent requests
4. **Pagination**: Efficient skip/limit queries

## 🎯 FamPay Assessment Compliance

This project fully meets all FamPay requirements:

### ✅ Basic Requirements
- [x] Background YouTube API fetching (10 seconds interval)
- [x] Stored video data with required fields
- [x] Paginated GET API (sorted by published datetime DESC)
- [x] Search API (title and description)
- [x] Dockerized application
- [x] Scalable and optimized architecture

### 🏆 Bonus Features
- [x] Multiple API key support with auto-rotation
- [x] Video dashboard with filters and sorting
- [x] Optimized partial match search

## 🤝 Support

If you encounter any issues:
1. Check the troubleshooting section above
2. Verify all environment variables are set correctly
3. Ensure Docker is running properly
4. Check logs: `docker-compose logs -f`
5. Or you can reach out to me at :- pjadwani_be21@thapar.edu

---

**Happy Coding!** 🚀
