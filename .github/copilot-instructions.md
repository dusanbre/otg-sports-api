# OTG Sport API - AI Coding Agent Instructions

## Architecture Overview

**Hybrid Stack:** Go backend with TypeScript/Drizzle for schema management only.
- **Go** (`main.go`): CLI-based application using Cobra (serve, sync, apikey commands)
- **TypeScript** (Node.js): Schema definition and migrations via Drizzle ORM
- **PostgreSQL**: Primary data store with connection pooling

**Core Components:**
1. `cmd/`: Cobra CLI commands (serve, sync, apikey)
2. `internal/api/`: REST API server with Chi router
3. `internal/goalserve/`: HTTP client for GoalServe API (rate-limited: 1 req/sec)
4. `internal/database/`: Singleton DB wrapper with Squirrel query builder
5. `internal/services/`: Business logic (sport-specific sync services)

**Supported Sports:**
- **Soccer**: `SoccerSyncService` → `soccer_matches` table
- **Basketball**: `BasketballSyncService` → `basketball_matches` table

**Data Flow:**
```
GoalServe API → Client (rate-limited) → SoccerSyncService     → soccer_matches
                                      → BasketballSyncService → basketball_matches
                ↓
        Fetch today + future 7 days → Upsert logic → sport-specific tables
                ↓
REST API ← Chi Router ← API Key Auth ← Rate Limiting ← Client Request
```

## CLI Commands

```bash
# Start the REST API server
go run main.go serve --port 8080

# Run the data sync scheduler
go run main.go sync

# API key management
go run main.go apikey create --name "My App" --sports soccer,basketball
go run main.go apikey list
go run main.go apikey revoke <id>
```

## REST API

**Base URL:** `/api/v1`
**Authentication:** API key via `X-API-Key` header or `Authorization: Bearer <key>`
**Documentation:** Swagger UI at `/swagger/index.html`

### Endpoints
- `GET /health` - Health check (public)
- `GET /api/v1/soccer/matches` - List soccer matches
- `GET /api/v1/soccer/matches/{id}` - Get single match
- `GET /api/v1/soccer/matches/live` - Live matches
- `GET /api/v1/soccer/leagues` - List leagues
- `GET /api/v1/basketball/matches` - List basketball matches
- `GET /api/v1/basketball/matches/{id}` - Get single match
- `GET /api/v1/basketball/matches/live` - Live matches
- `GET /api/v1/basketball/leagues` - List leagues

## Critical Patterns

### Database Access
- **Singleton pattern** in `database.GetInstance()` - never create direct connections
- Use **Squirrel builder** for queries: `db.Builder.Select(...).From(...).Where(...)`
- PostgreSQL uses `$1, $2` placeholders via `sq.Dollar` format
- Connection pool: 25 max open, 5 max idle

Example query pattern (from `queries.go`):
```go
query := db.Builder.Select("*").From("soccer_matches").Where("match_id = ?", matchID)
sql, args, _ := query.ToSql()
db.Conn.QueryRow(sql, args...).Scan(...)
```

### Schema Management
- **Define schemas** in `migrations/schema.ts` using Drizzle
- **Generate migrations**: `npx drizzle-kit generate` (creates SQL in `migrations/drizzle/`)
- **Apply migrations**: Use Drizzle push or run SQL directly
- **Go models** (`database/models.go`) must match TypeScript schema - use `sql.Null*` types

### API Client Conventions
- **Rate limiting**: GoalServe client uses `time.Ticker` (1 req/sec) - always `<-c.rateLimiter.C`
- **JSON handling**: Handle both single object and arrays (see `GoalServeSoccerMatchesData.UnmarshalJSON`)
- **Date parsing**: Supports `02.01.2006` format; combine date+time for match scheduling
- **Error handling**: Log and continue on single match failures to avoid blocking batch sync

### Service Layer
- **Upsert pattern** in `soccer_sync.go`: Check existence → Insert or Update
- Returns `(inserted bool, error)` to track sync metrics
- Logs inserted/updated counts at end of each sync run
- Fetches 3 time windows: today, past 7 days, future 7 days

## Development Workflows

### Local Setup
```bash
# Start Postgres
docker compose up -d db

# Environment vars (.env required)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=otg
GOALSERVE_API_KEY=<your_key>
GOALSERVE_URL=https://www.goalserve.com

# Run migrations (Node.js)
npm install
npx drizzle-kit generate
npx drizzle-kit push

# Build & run Go app
go run main.go
```

### Adding New Fields
1. Update `migrations/schema.ts` (Drizzle schema)
2. Run `npx drizzle-kit generate` and `npx drizzle-kit push`
3. Update `database/models.go` (Go struct with matching types)
4. Update the sport-specific sync service (e.g., `services/soccer_sync.go` or `services/basketball_sync.go`)

### Adding a New Sport
1. Create `internal/goalserve/{sport}_models.go` with API response structs
2. Add fetch methods to `internal/goalserve/client.go` (e.g., `FetchBasketballTodayMatches`)
3. Add table schema to `migrations/schema.ts` and run migrations
4. Add Go struct to `internal/database/models.go`
5. Create `internal/services/{sport}_sync.go` with sync service
6. Add query methods to `internal/database/queries.go`
7. Register scheduler job in `main.go`

### Testing Match Sync
- Sample data: `etc/sample/soccernew.json`, `etc/sample/bsktbl_home.json`
- Manual sync: Run `go run main.go sync` (runs immediate sync on startup before scheduler)
- Scheduler runs every 1 minute (configured in `cmd/sync.go` via `gocron.DurationJob`)

### Regenerating Swagger Docs
```bash
# After modifying handler annotations
~/go/bin/swag init -g main.go -o internal/api/docs --parseDependency --parseInternal
```

## Key Files Reference

### CLI
- [cmd/root.go](cmd/root.go): Cobra root command
- [cmd/serve.go](cmd/serve.go): REST API server command
- [cmd/sync.go](cmd/sync.go): Sync scheduler command
- [cmd/apikey.go](cmd/apikey.go): API key management commands

### API
- [internal/api/server.go](internal/api/server.go): Chi router setup
- [internal/api/handlers/](internal/api/handlers/): Request handlers
- [internal/api/middleware/](internal/api/middleware/): Auth, CORS, rate limiting
- [internal/api/docs/](internal/api/docs/): Generated Swagger documentation

### Soccer
- [internal/services/soccer_sync.go](internal/services/soccer_sync.go): Soccer upsert logic
- [internal/goalserve/soccer_models.go](internal/goalserve/soccer_models.go): Soccer API response models

### Basketball
- [internal/services/basketball_sync.go](internal/services/basketball_sync.go): Basketball upsert logic
- [internal/goalserve/basketball_models.go](internal/goalserve/basketball_models.go): Basketball API response models

### Shared
- [main.go](main.go): Scheduler setup, graceful shutdown
- [internal/database/db.go](internal/database/db.go): Singleton DB connection with Squirrel
- [internal/database/models.go](internal/database/models.go): Go structs for all sports
- [internal/database/queries.go](internal/database/queries.go): Query methods for all sports
- [internal/goalserve/client.go](internal/goalserve/client.go): Rate-limited API client
- [migrations/schema.ts](migrations/schema.ts): Source of truth for DB schema

## Common Gotchas

- **Don't** use Drizzle ORM in Go code - only for schema/migrations
- **Never** skip rate limiter in `goalserve.Client` methods
- **Always** use `sql.Null*` types for nullable columns in Go structs
- **Match time parsing**: Handle both `@formatted_date` and `@date` fields (fallback logic)
- **Events field**: JSON stored as string in DB; marshal/unmarshal in Go service layer
