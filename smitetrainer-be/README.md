# smitetrainer-be

Go backend for estimated dragon HP series from Riot Match-V5 data.

## 1) First-time setup (Windows)

1. Install Go 1.22+ from https://go.dev/dl/
2. Open a new PowerShell window.
3. Check Go is available:
   - `go version`
4. If `go` is not found, use one of these:
   - Current shell only: `$env:Path += ';C:\Program Files\Go\bin'`
   - Then verify again: `go version`

## 2) Configure `.env`

Create/update `smitetrainer-be/.env`:

```env
RIOT_API_KEY="YOUR_RIOT_KEY"

# memory (default) or redis
CACHE_BACKEND="memory"

# shared/persistent cache (used when CACHE_BACKEND=redis)
REDIS_ADDR="127.0.0.1:6379"
REDIS_USERNAME=""
REDIS_PASSWORD=""
REDIS_DB="0"
```

## 3) Run with Docker Compose (API + Redis)

From `smitetrainer-be`:

`docker compose up --build`

This starts:

- API on `http://localhost:8080`
- Redis on `localhost:6379`

Stop it with:

`docker compose down`

To also remove Redis data volume:

`docker compose down -v`

## 4) Run backend without Docker (optional)

From `smitetrainer-be`:

`go run ./cmd/server`

If `go` is still not on PATH:

`& 'C:\Program Files\Go\bin\go.exe' run ./cmd/server`

Server starts on `http://localhost:8080` unless `PORT` is set.

## 5) Try the API

Health check:

`Invoke-RestMethod http://localhost:8080/healthz`

Dragon HP endpoint:

`Invoke-RestMethod "http://localhost:8080/v1/matches/EUW1_1234567890/dragons/2/hp?tickMs=500&windowMs=20000"`

## Environment variables

- `RIOT_API_KEY` (required)
- `PORT` (default: `8080`)
- `RIOT_HTTP_TIMEOUT_SECONDS` (default: `8`)
- `RIOT_MAX_ATTEMPTS` (default: `5`)
- `CACHE_BACKEND` (`memory` or `redis`, default: `memory`)
- `CACHE_TTL_SECONDS` (default: `300`)
- `CACHE_MAX_ENTRIES` (default: `256`)
- `REDIS_ADDR` (default: `127.0.0.1:6379`)
- `REDIS_USERNAME` (default: empty)
- `REDIS_PASSWORD` (default: empty)
- `REDIS_DB` (default: `0`)
- `REDIS_DIAL_TIMEOUT_SECONDS` (default: `3`)
- `REDIS_IO_TIMEOUT_SECONDS` (default: `3`)
- `DEFAULT_TICK_MS` (default: `500`)
- `DEFAULT_WINDOW_MS` (default: `20000`)
- `DEFAULT_BURST_MS` (default: `1200`)
- `FIGHT_RADIUS` (default: `2500`)

## Endpoint

`GET /v1/matches/{matchId}/dragons/{dragonIndex}/hp`

Query params:

- `tickMs` (optional, default `500`)
- `windowMs` (optional, default `20000`)
