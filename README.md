# SmiteTrainer

SmiteTrainer currently has two apps in one repo:

1. `smitetrainer-fe` (Vue): game UI, home screen, multiplayer reaction mode, and a training page that calls backend dragon HP API.
2. `smitetrainer-be` (Go): Riot Match-V5 data service that returns estimated dragon HP timelines.

## What The Backend Does

The Go API exposes:

1. `GET /healthz`
2. `GET /v1/matches/{matchId}/dragons/{dragonIndex}/hp`

For the dragon endpoint, it:

1. Fetches `match` and `timeline` from Riot Match-V5.
2. Finds the requested Nth elemental dragon kill event.
3. Builds an estimated dragon HP series over a window ending at kill time.
4. Returns model metadata, HP points, and basic fight markers.

Cache behavior:

1. Redis can be used as the shared cache.
2. In-memory cache is always available as fallback.
3. TTL is controlled with `CACHE_TTL_SECONDS`.

## Architecture (Mermaid)

```mermaid
flowchart LR
  U[User Browser]

  subgraph FE[smitetrainer-fe (Vue)]
    FEHome[Home / Multiplayer UI]
    FETraining[Training View\nForm + HP chart + response preview]
    FERouter[Vue Router]
  end

  subgraph BE[smitetrainer-be (Go API)]
    API[HTTP Handlers\n/healthz\n/v1/matches/:matchId/dragons/:dragonIndex/hp]
    Parser[Timeline Parser\nFind Nth dragon kill\nExtract fight markers]
    Sim[HP Simulator\nlinear + end-burst estimate]
    RC[Riot Client\nretry/backoff\nrate-limit handling]
    Cache[Cache Interface]
    Mem[(In-memory LRU)]
    Redis[(Redis)]
  end

  Riot[(Riot Match-V5 API)]

  U --> FEHome
  FEHome --> FERouter
  FEHome --> FETraining
  FETraining --> API

  API --> Parser
  API --> Sim
  API --> RC
  RC --> Cache
  Cache --> Mem
  Cache --> Redis
  RC --> Riot
```

## How To Run (Recommended: Docker for Backend)

### 1) Configure backend env

Edit `smitetrainer-be/.env` and set your key:

```env
RIOT_API_KEY="YOUR_RIOT_KEY"
CACHE_BACKEND="redis"
REDIS_ADDR="127.0.0.1:6379"
REDIS_USERNAME=""
REDIS_PASSWORD=""
REDIS_DB="0"
```

### 2) Start backend + redis

From `smitetrainer-be`:

```powershell
docker compose up --build
```

Expected:

1. API listens on `http://localhost:8080`
2. Redis listens on `localhost:6379`
3. API logs include `cache backend: redis (memory fallback enabled)`

Stop:

```powershell
docker compose down
```

Remove Redis data too:

```powershell
docker compose down -v
```

### 3) Verify backend is healthy

```powershell
Invoke-RestMethod http://localhost:8080/healthz
```

Expected JSON:

```json
{"status":"ok"}
```

### 4) Test dragon endpoint

```powershell
Invoke-RestMethod "http://localhost:8080/v1/matches/EUW1_1234567890/dragons/2/hp?tickMs=500&windowMs=20000"
```

Expected response shape includes:

1. `matchId`, `dragonIndex`, `dragonType`, `killMs`
2. `hpModel` with `isEstimated=true`
3. `series` array of `{tMs, hp, hpPct}`
4. `markers` (fight stats in analysis window)

## Run Frontend (Separate Terminal)

From `smitetrainer-fe`:

```powershell
npm install
$env:VUE_APP_API_BASE_URL="http://localhost:8080"
$env:PORT=8081
npm run serve
```

Open:

1. `http://localhost:8081/`
2. `http://localhost:8081/multiplayer`
3. `http://localhost:8081/training`

## What To Expect In The UI

1. Home view with navigation buttons and theme styling.
2. Multiplayer view is a local reaction game (best of 5 rounds).
3. Training view is wired to backend endpoint `GET /v1/matches/{matchId}/dragons/{dragonIndex}/hp`.
4. About page is implemented as project info + navigation.

## Common Failure Cases

1. `401/403` from Riot: key missing, expired, or invalid.
2. `429` from Riot: rate-limited; backend retries with backoff.
3. `404` timeline or match not found.
4. Backend starts but no data: check `.env` and Riot key first.

## Useful Backend Endpoints

1. Health: `GET /healthz`
2. Dragon HP: `GET /v1/matches/{matchId}/dragons/{dragonIndex}/hp`
