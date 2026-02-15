# Unfollow Tracker - Project Structure

## Overview

A SaaS application for Instagram users to track unfollowers. Built with Go for high-performance data processing, PostgreSQL for reliable storage, and a modern web frontend.

---

## 📁 Root Directory Structure

```
unfollow-tracker/
├── 📁 cmd/                    # Application entry points
│   ├── api/                   # REST API server
│   │   └── main.go
│   ├── fetcher/               # Scraping worker service
│   │   └── main.go
│   ├── scheduler/             # Cron-like job scheduler
│   │   └── main.go
│   └── migrator/              # Database migration runner
│       └── main.go
│
├── 📁 internal/               # Private application code
│   ├── api/                   # API layer
│   │   ├── handlers/          # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── accounts.go
│   │   │   ├── unfollows.go
│   │   │   └── health.go
│   │   ├── middleware/        # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── ratelimit.go
│   │   │   └── logging.go
│   │   ├── routes.go
│   │   └── server.go
│   │
│   ├── domain/                # Core business logic (pure Go)
│   │   ├── account.go         # Account entity
│   │   ├── follower.go        # Follower entity (hashed)
│   │   ├── snapshot.go        # Snapshot entity
│   │   ├── unfollow.go        # Unfollow entity
│   │   └── errors.go          # Domain errors
│   │
│   ├── service/               # Business logic services
│   │   ├── account_service.go
│   │   ├── diff_service.go    # Chunked diffing logic
│   │   ├── scan_service.go    # Scan orchestration
│   │   └── notification_service.go
│   │
│   ├── repository/            # Data access layer
│   │   ├── postgres/
│   │   │   ├── account_repo.go
│   │   │   ├── snapshot_repo.go
│   │   │   └── unfollow_repo.go
│   │   └── redis/
│   │       ├── follower_cache.go
│   │       ├── rate_limiter.go
│   │       └── queue.go
│   │
│   ├── fetcher/               # Scraping/fetching logic
│   │   ├── instagram/
│   │   │   ├── scraper.go     # Public profile scraping
│   │   │   ├── parser.go      # HTML/JSON parsing
│   │   │   └── ratelimit.go   # Platform-specific limits
│   │   ├── proxy/
│   │   │   ├── pool.go        # Proxy rotation
│   │   │   └── health.go      # Proxy health checks
│   │   └── worker.go          # Fetch worker logic
│   │
│   ├── queue/                 # Queue management
│   │   ├── producer.go
│   │   ├── consumer.go
│   │   └── jobs.go            # Job definitions
│   │
│   └── config/                # Configuration
│       └── config.go
│
├── 📁 pkg/                    # Public/shared packages
│   ├── hasher/                # Follower ID hashing
│   │   └── hasher.go
│   ├── chunker/               # Chunk management for diffing
│   │   └── chunker.go
│   └── httputil/              # HTTP utilities
│       └── client.go
│
├── 📁 migrations/             # Database migrations (SQL)
│   ├── 001_create_accounts.up.sql
│   ├── 001_create_accounts.down.sql
│   ├── 002_create_snapshots.up.sql
│   ├── 002_create_snapshots.down.sql
│   ├── 003_create_unfollows.up.sql
│   └── 003_create_unfollows.down.sql
│
├── 📁 web/                    # Frontend SaaS application
│   ├── 📁 src/
│   │   ├── 📁 components/
│   │   │   ├── common/        # Shared UI components
│   │   │   │   ├── Button.jsx
│   │   │   │   ├── Card.jsx
│   │   │   │   ├── Modal.jsx
│   │   │   │   └── Loading.jsx
│   │   │   ├── layout/
│   │   │   │   ├── Navbar.jsx
│   │   │   │   ├── Sidebar.jsx
│   │   │   │   └── Footer.jsx
│   │   │   ├── auth/
│   │   │   │   ├── LoginForm.jsx
│   │   │   │   └── RegisterForm.jsx
│   │   │   ├── dashboard/
│   │   │   │   ├── StatCards.jsx
│   │   │   │   ├── UnfollowList.jsx
│   │   │   │   ├── TrendChart.jsx
│   │   │   │   └── AccountHealth.jsx
│   │   │   └── accounts/
│   │   │       ├── AddAccountModal.jsx
│   │   │       ├── AccountCard.jsx
│   │   │       └── AccountList.jsx
│   │   │
│   │   ├── 📁 pages/
│   │   │   ├── Landing.jsx        # Marketing landing page
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   ├── Dashboard.jsx      # Main dashboard
│   │   │   ├── Accounts.jsx       # Manage tracked accounts
│   │   │   ├── Unfollows.jsx      # Detailed unfollow history
│   │   │   ├── Settings.jsx
│   │   │   └── Pricing.jsx
│   │   │
│   │   ├── 📁 hooks/
│   │   │   ├── useAuth.js
│   │   │   ├── useAccounts.js
│   │   │   └── useUnfollows.js
│   │   │
│   │   ├── 📁 services/
│   │   │   ├── api.js             # Axios/fetch wrapper
│   │   │   ├── auth.js
│   │   │   ├── accounts.js
│   │   │   └── unfollows.js
│   │   │
│   │   ├── 📁 store/              # State management (Zustand/Redux)
│   │   │   ├── authStore.js
│   │   │   └── accountStore.js
│   │   │
│   │   ├── 📁 styles/
│   │   │   ├── index.css          # Global styles
│   │   │   ├── variables.css      # CSS variables/tokens
│   │   │   └── components/        # Component styles
│   │   │
│   │   ├── 📁 utils/
│   │   │   ├── formatters.js
│   │   │   └── validators.js
│   │   │
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── router.jsx
│   │
│   ├── public/
│   │   ├── favicon.ico
│   │   └── assets/
│   │
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── .env.example
│
├── 📁 deploy/                 # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile.api
│   │   ├── Dockerfile.fetcher
│   │   ├── Dockerfile.scheduler
│   │   ├── Dockerfile.migrator
│   │   └── Dockerfile.web
│   ├── docker-compose.yml     # Local development
│   ├── docker-compose.prod.yml
│   └── kubernetes/            # K8s manifests (optional)
│       ├── api-deployment.yaml
│       ├── fetcher-deployment.yaml
│       └── ingress.yaml
│
├── 📁 scripts/                # Utility scripts
│   ├── setup.sh               # Initial setup
│   ├── dev.sh                 # Run dev environment
│   └── seed.sh                # Seed test data
│
├── 📁 docs/                   # Documentation
│   ├── api.md                 # API documentation
│   ├── architecture.md        # System architecture
│   └── deployment.md          # Deployment guide
│
├── 📁 plans/                  # (existing) Planning docs
│   ├── strategies_to_fetch.md
│   ├── tech_stack.md
│   └── project_structure.md   # This file
│
├── .env.example               # Environment template
├── .gitignore
├── go.mod                     # Go modules
├── go.sum
├── Makefile                   # Build commands
└── README.md
```

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                           FRONTEND (Web SaaS)                        │
│                    React/Vite + Modern CSS                           │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                            API SERVICE                               │
│                         (Go + chi/net-http)                          │
│    ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│    │   Auth   │  │ Accounts │  │Unfollows │  │  Health/Status   │   │
│    └──────────┘  └──────────┘  └──────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
         │                                              │
         ▼                                              ▼
┌──────────────────┐                        ┌──────────────────────────┐
│   PostgreSQL     │                        │         Redis            │
│   (Source of     │◄──────────────────────►│  - Follower Hash Sets    │
│    Truth)        │                        │  - Rate Limiting         │
│                  │                        │  - Job Queue (Streams)   │
│  - accounts      │                        │  - Session Cache         │
│  - snapshots     │                        └──────────────────────────┘
│  - unfollows     │                                    │
│  - users         │                                    │
└──────────────────┘                                    │
                                                        ▼
                                        ┌──────────────────────────────┐
                                        │      SCHEDULER SERVICE       │
                                        │   (Cron → Queue Producer)    │
                                        │                              │
                                        │  - chunk rotation logic      │
                                        │  - adaptive frequency        │
                                        │  - job scheduling            │
                                        └──────────────────────────────┘
                                                        │
                                                        ▼
                                        ┌──────────────────────────────┐
                                        │      FETCHER WORKERS         │
                                        │    (Scalable, Stateless)     │
                                        │                              │
                                        │  - consume from queue        │
                                        │  - scrape public profiles    │
                                        │  - proxy rotation            │
                                        │  - respect rate limits       │
                                        │  - produce diff results      │
                                        └──────────────────────────────┘
                                                        │
                                                        ▼
                                        ┌──────────────────────────────┐
                                        │       PROXY POOL             │
                                        │    (IP Rotation Layer)       │
                                        └──────────────────────────────┘
```

---

## 📦 Key Components Explained

### 1. **API Service** (`cmd/api`)
- REST endpoints for frontend
- JWT authentication
- Rate limiting per user
- Endpoints:
  - `POST /auth/register`, `POST /auth/login`
  - `POST /accounts/track` - Add account to track
  - `GET /accounts` - List tracked accounts
  - `GET /unfollows` - Get unfollow events
  - `GET /health` - Health check

### 2. **Fetcher Service** (`cmd/fetcher`)
- Consumes jobs from Redis queue
- Scrapes public Instagram profiles
- Implements chunked scanning strategy
- Produces diffs to Redis/Postgres
- **Scale horizontally**: Run 2-4+ instances

### 3. **Scheduler Service** (`cmd/scheduler`)
- Cron-like scheduling
- Rotates chunks for each account
- Implements adaptive polling frequency
- Enqueues fetch jobs

### 4. **Diff Service** (`internal/service/diff_service.go`)
The heart of the system:
```go
// Conceptual flow
func (s *DiffService) ProcessChunk(accountID, chunkID string, currentFollowers []string) {
    currentHashes := hash(currentFollowers)
    previousHashes := s.redis.GetChunk(accountID, chunkID)
    
    unfollows := difference(previousHashes, currentHashes)
    newFollows := difference(currentHashes, previousHashes)
    
    s.persistUnfollows(accountID, unfollows)
    s.redis.UpdateChunk(accountID, chunkID, currentHashes)
}
```

---

## 🗄️ Database Schema (PostgreSQL)

```sql
-- Users (SaaS customers)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    plan VARCHAR(50) DEFAULT 'free',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tracked accounts (Instagram profiles being monitored)
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) DEFAULT 'instagram',
    username VARCHAR(255) NOT NULL,
    follower_count INTEGER,
    last_scan_at TIMESTAMPTZ,
    scan_status VARCHAR(50) DEFAULT 'pending',
    chunk_count INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, platform, username)
);

-- Follower snapshots (chunked, hashed)
CREATE TABLE snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
    chunk_id INTEGER NOT NULL,
    follower_hashes TEXT[] NOT NULL, -- Array of hashed IDs
    scanned_at TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_snapshots_account_chunk (account_id, chunk_id)
);

-- Detected unfollows
CREATE TABLE unfollows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
    follower_hash VARCHAR(64) NOT NULL, -- Hashed for privacy
    detected_at TIMESTAMPTZ DEFAULT NOW(),
    notified BOOLEAN DEFAULT FALSE,
    INDEX idx_unfollows_account (account_id),
    INDEX idx_unfollows_detected (detected_at)
);

-- Subscription/billing (for SaaS)
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    plan VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    current_period_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 🐳 Docker Compose (Development)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: unfollow_tracker
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  api:
    build:
      context: .
      dockerfile: deploy/docker/Dockerfile.api
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:postgres@postgres:5432/unfollow_tracker
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis

  fetcher:
    build:
      context: .
      dockerfile: deploy/docker/Dockerfile.fetcher
    environment:
      - DATABASE_URL=postgres://postgres:postgres@postgres:5432/unfollow_tracker
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis
    deploy:
      replicas: 2

  scheduler:
    build:
      context: .
      dockerfile: deploy/docker/Dockerfile.scheduler
    environment:
      - DATABASE_URL=postgres://postgres:postgres@postgres:5432/unfollow_tracker
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis

  web:
    build:
      context: ./web
      dockerfile: ../deploy/docker/Dockerfile.web
    ports:
      - "3000:3000"
    depends_on:
      - api

volumes:
  postgres_data:
  redis_data:
```

---

## 🚀 Getting Started (Makefile)

```makefile
.PHONY: dev build test migrate

# Run development environment
dev:
	docker-compose up -d postgres redis
	go run cmd/api/main.go &
	cd web && npm run dev

# Build all services
build:
	docker-compose build

# Run tests
test:
	go test ./...
	cd web && npm test

# Run migrations
migrate:
	go run cmd/migrator/main.go up

# Generate API docs
docs:
	swag init -g cmd/api/main.go -o docs/swagger
```

---

## 🎯 Ethical Design Considerations (Built-In)

Based on your planning docs, the architecture supports:

1. **Privacy by Design**
   - Never store raw follower IDs (only hashes)
   - No exact unfollow timestamps shown to users
   - Aggregate data where possible

2. **Gentle Notifications**
   - Delayed summaries (not instant alerts)
   - "Audience Health Score" instead of raw numbers
   - Cooldown periods between notifications

3. **Sustainable Scraping**
   - Chunked scanning reduces requests 20×
   - Adaptive polling based on account size
   - Proxy rotation to avoid bans
   - Rate limiting per platform

---

## 📝 Next Steps

1. **Initialize the project**:
   ```bash
   go mod init github.com/yourname/unfollow-tracker
   ```

2. **Set up the database schema**

3. **Implement core diff logic**

4. **Build the fetcher with proxy rotation**

5. **Create the SaaS frontend**

Would you like me to start scaffolding any specific part of this structure?
