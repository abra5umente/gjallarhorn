# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Development (runs Go backend + Vite frontend together)
./scripts/dev.sh

# Or manually:
go run .          # Backend on :8080
npm run dev       # Frontend on :5173

# Production build (creates single binary with embedded frontend)
./scripts/build.sh

# Linting
npm run lint      # ESLint for frontend
go fmt ./...      # Go formatting

# Docker
docker-compose up -d
```

## Architecture

Gjallarhorn is an uptime monitoring app that sends Pushover notifications when services go down. It deploys as a single Go binary with the React frontend embedded.

### Backend (Go/Echo)

- **main.go**: Entry point, Echo routes, embeds `dist/` folder, serves both API and static files
- **monitor.go**: `MonitorService` runs background health checks at configurable interval (`CHECK_INTERVAL` env var, default 60s). Tracks `ConsecutiveFailures` and marks services offline after 3 consecutive failures. Sends notifications on status transitions and hourly reminders for prolonged outages
- **notifications.go**: Pushover integration with three notification types (down alert, hourly reminder, recovery)
- **storage.go**: JSON file persistence to `/data/services.json` and `/data/config.json`
- **models.go**: Data models and validation (interval bounds 30-3600s)

HTTP 2xx, 3xx, and 401 are all considered healthy (401 means auth-required services don't false-trigger).

### Frontend (React/Vite/TailwindCSS)

- **src/context/**: `ServiceContext` and `NotificationContext` manage state (no Redux)
- **src/services/api.js**: Axios client with `/api` baseURL, 10s timeout
- **src/components/**: ServiceList (with bulk selection), ServiceForm, NotificationSettings, StatusBadge, Header, LoadingSpinner, BulkActionBar, BulkEditModal

UI auto-refreshes services every 30 seconds.

### API Routes

```
GET/POST         /api/services
PUT/DELETE       /api/services/:id
GET              /api/services/:id/status
POST/PUT/DELETE  /api/services/bulk     # Bulk operations (all-or-nothing)
GET/POST         /api/notifications/config
GET              /swagger/*              # Swagger UI docs
```

## Coding Style

- Go: Run `go fmt ./...`, use CamelCase types, mixedCase functions
- Frontend: ESLint config, two-space indentation, single quotes, PascalCase components
- Commits: Imperative mood, ≤72 chars, reference issues with `Refs #123`

## Testing

No automated tests exist yet. When adding:
- Go: `_test.go` files alongside sources, table-driven tests, `go test ./...`
- Frontend: Vitest/React Testing Library in `src/__tests__/`

## Configuration

Copy `env.example` to `.env` for local secrets. Key env vars:
- `PORT` (default 8080)
- `CHECK_INTERVAL` (default 60s) - health check interval in seconds
- `PUSHOVER_USER_KEY`, `PUSHOVER_APP_TOKEN`, `PUSHOVER_ENABLED`

Data persists to `/data` directory (auto-created).
