# 📊 Unfollow Tracker

A privacy-conscious SaaS application to track Instagram unfollowers. Built with Go for high-performance data processing and designed with ethical considerations in mind.

## 🏗️ Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│   API       │────▶│  PostgreSQL │
│   (React)   │     │   (Go)      │     │  (Truth)    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐     ┌─────────────┐
                    │   Redis     │◀────│  Scheduler  │
                    │   (Cache)   │     │  (Go)       │
                    └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Fetchers   │ ×N (scalable workers)
                    │  (Go)       │
                    └─────────────┘
```

## ✨ Features

- **Privacy-First**: Follower IDs are hashed, never stored raw
- **Efficient Scanning**: Chunked diffing reduces API calls by 20×
- **Scalable**: Stateless fetcher workers scale horizontally
- **Ethical Design**: Delayed notifications, health scores over raw numbers
- **SaaS Ready**: User authentication, subscription plans

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Node.js 18+ (for frontend)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/utkarsh/unfollow-tracker.git
cd unfollow-tracker

# Copy environment file
cp .env.example .env

# Start infrastructure (Postgres + Redis)
make docker-up

# Run database migrations
make migrate

# Start the API server
make api

# In another terminal, start the frontend
make web-install
make web
```

### Running All Services

```bash
# Start everything with Docker Compose
make dev-all
```

## 📁 Project Structure

```
├── cmd/                    # Application entry points
│   ├── api/               # REST API server
│   ├── fetcher/           # Scraping worker
│   ├── scheduler/         # Job scheduler
│   └── migrator/          # Database migrations
├── internal/              # Private application code
│   ├── api/              # HTTP handlers & middleware
│   ├── domain/           # Core business entities
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   ├── fetcher/          # Scraping logic
│   └── queue/            # Job queue management
├── pkg/                   # Shared packages
├── migrations/            # SQL migration files
├── web/                   # React frontend
└── deploy/               # Docker & K8s configs
```

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go, chi (HTTP router) |
| Database | PostgreSQL |
| Cache/Queue | Redis |
| Frontend | React, Vite |
| Deployment | Docker, Kubernetes |

## 📚 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | User login |
| GET | `/api/accounts` | List tracked accounts |
| POST | `/api/accounts` | Add account to track |
| GET | `/api/unfollows` | Get unfollow history |
| GET | `/api/health` | Health check |

## 🔒 Privacy & Ethics

This application is designed with ethical considerations:

- **No raw data**: Follower IDs are hashed before storage
- **Delayed notifications**: No instant anxiety-triggering alerts
- **Health scores**: Focus on trends, not individual losses
- **Respect rate limits**: Sustainable scraping practices

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
