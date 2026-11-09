# Distributed Metrics Platform

A small Go HTTP service for ingesting and storing application metrics (counters, gauges, histograms), backed by PostgreSQL and exposing itself for scraping by Prometheus/Grafana.

## What it does

Clients POST metric data points to the service. Each metric is validated, normalized, and persisted to PostgreSQL. The service also exposes a Prometheus endpoint reporting its own operational metrics (ingest counts, errors, latency), so the ingestion pipeline itself can be monitored via Prometheus + Grafana.

## Architecture

```
cmd/server        entrypoint: runs DB migrations, starts HTTP server
internal/handlers HTTP layer (/ingest)
internal/services business logic: validation + orchestration of ingestion
internal/utils    normalization/validation helpers (names, labels, timestamps)
internal/models   Metric domain type
internal/storage  Postgres connection pool + persistence (pgx)
internal/telemetry Prometheus instrumentation (counters/histograms for the service itself)
migrations        SQL schema migrations (golang-migrate)
```

### Request flow

1. `POST /ingest` → `handlers.IngestHandler` decodes JSON into a `Metric` and does basic request validation.
2. `services.Ingest` validates metric type/value, normalizes the name/labels, fills in a timestamp if missing, and records ingest duration/outcome metrics.
3. `storage.Save` writes the metric row into Postgres via a pgx connection pool.

## HTTP Endpoints

| Endpoint    | Method | Description                                   |
|-------------|--------|------------------------------------------------|
| `/ingest`   | POST   | Accepts a JSON-encoded metric and stores it   |
| `/metrics`  | GET    | Prometheus-format metrics for this service    |
| `/healthz`  | GET    | Liveness/readiness check (returns `200 ok`)   |

### Metric payload

```json
{
  "name": "http_requests",
  "value": 1,
  "type": "counter",
  "labels": {"route": "/api/foo"},
  "timestamp": 1731000000
}
```

- `name` is required and lower-cased/trimmed.
- `type` must be one of `counter`, `gauge`, `histogram`.
- `value` must be a finite number (no NaN/Inf).
- `labels` keys/values are trimmed; keys are lower-cased.
- `timestamp` is a Unix timestamp; if omitted/zero, the server fills in "now".

## Data model

Metrics are stored in a single `metrics` table (see `migrations/000001_create_metrics_table.up.sql`):

- `id` (bigserial), `name`, `value` (float8), `type`, `labels` (jsonb), `timestamp` (bigint)
- Indexed on `(name, timestamp DESC)` for efficient lookups by metric name.

## Observability stack

The service emits its own metrics via `internal/telemetry`:

- `metrics_ingested_total` — counter of successful ingests
- `metrics_ingest_errors_total` — counter of failed ingests
- `ingest_duration_seconds` — histogram of ingest latency, labeled by `outcome`

`docker-compose.yml` wires up a full local observability stack:

- **app** — this service (port 8080)
- **db** — PostgreSQL 18 (port 5432)
- **prometheus** — scrapes the app's `/metrics` endpoint (port 9090, config in `prometheus.yml`)
- **grafana** — dashboards over Prometheus data (port 3000, provisioning in `grafana/provisioning`)

## Running locally

```bash
docker-compose up --build
```

This starts the app, Postgres, Prometheus, and Grafana. Database migrations run automatically on app startup (`runMigrations` in `cmd/server/main.go`).

Then:

```bash
curl -X POST localhost:8080/ingest \
  -H 'Content-Type: application/json' \
  -d '{"name":"http_requests","value":1,"type":"counter"}'

curl localhost:8080/metrics
curl localhost:8080/healthz
```

### Environment variables

- `DATABASE_URL` — Postgres connection string, e.g. `postgres://user:pass@host:5432/dbname?sslmode=disable`

## Testing

```bash
go test -v -race ./...
```

Requires a reachable Postgres instance (see `DATABASE_URL`) with migrations applied — CI spins up a Postgres service container for this (`.github/workflows/continuous-integration.yml`).

## Deployment

- **Docker** — multi-stage `Dockerfile` builds a static Go binary on Alpine and ships it with the `migrations` directory.
- **Kubernetes** — manifests under `k8s/` (namespace, configmap, secret example, deployment, service) run the container with DB connection details split across a ConfigMap (host/port/db/sslmode) and a Secret (user/password), reassembled into `DATABASE_URL` via env var interpolation. Includes startup/readiness/liveness probes against `/healthz`.
- **CI/CD** — GitHub Actions:
  - `continuous-integration.yml`: runs migrations + tests against a Postgres service container on every push to `main`, then builds the binary.
  - `deploy.yml`: deploys the service (see workflow for target, e.g. Google Cloud Run).

## Tech stack

- Go 1.26
- [pgx/v5](https://github.com/jackc/pgx) — PostgreSQL driver/connection pooling
- [golang-migrate](https://github.com/golang-migrate/migrate) — schema migrations
- [prometheus/client_golang](https://github.com/prometheus/client_golang) — metrics instrumentation
- Prometheus + Grafana — metrics scraping and visualization
