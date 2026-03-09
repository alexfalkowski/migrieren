![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/migrieren.svg?style=shield)](https://circleci.com/gh/alexfalkowski/migrieren)
[![codecov](https://codecov.io/gh/alexfalkowski/migrieren/graph/badge.svg?token=R2OD8WIKD0)](https://codecov.io/gh/alexfalkowski/migrieren)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/migrieren)](https://goreportcard.com/report/github.com/alexfalkowski/migrieren)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/migrieren.svg)](https://pkg.go.dev/github.com/alexfalkowski/migrieren)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Migrieren

Migrieren is a Go service that runs database migrations via a gRPC API with an HTTP RPC facade.

The service centralizes migration execution while still using native migration assets through [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate).

- Runtime: Go
- API contract: protobuf/gRPC (`api/migrieren/v1/service.proto`)
- Test harness: Ruby + Cucumber (`test/`)

## What it does

- Exposes one API call: `migrieren.v1.Service/Migrate`
- Looks up a logical database name in config
- Resolves migration source URL and database URL from configured sources
- Runs migrations with `golang-migrate`
- Returns migration logs and response metadata
- Exposes the same behavior over gRPC and HTTP RPC

## Redis-backed locking

Migration execution is guarded by a distributed mutex (Redsync over Redis/Valkey).

- Lock scope: one lock per `(source URL, database URL)` pair
- Purpose: prevent concurrent migration attempts against the same target across instances
- Config key: `redis.url`
- URL source resolution: read via go-service FS (`file:...` indirection supported)

Operational notes:

- If lock acquisition fails, migration returns an invalid migration error (transport maps to gRPC `Internal` / HTTP `500`)
- Redis client telemetry (tracing + metrics) is enabled
- Redis maintenance notifications are explicitly disabled in client options

## Supported drivers

Driver support is wired in code under `internal/migrate/source` and `internal/migrate/database`.

- Migration source schemes: `file://`, `github://`
- Database scheme: `pgx5://` (rewritten internally to `postgres://` for driver setup)

## Configuration

Pass config to `server` via `-i`:

```bash
cd test
../migrieren server -i file:.config/server.yml
```

`file:` source paths are resolved from the process working directory. The test config uses
`file:secrets/...`, so running from `test/` makes those paths resolve to `test/secrets/...`.

### Minimal config shape

```yaml
environment: development
health:
  duration: 1s
  timeout: 1s
migrate:
  databases:
    - name: postgres
      source: file:secrets/source
      url: file:secrets/pg
redis:
  url: file:secrets/redis
transport:
  http:
    address: tcp://:11000
  grpc:
    address: tcp://:12000
```

In this repository's test setup, those source files are under `test/secrets/`:

```text
test/secrets/source -> file://migrations
test/secrets/pg     -> pgx5://test:test@localhost:5433/test?sslmode=disable
test/secrets/redis  -> redis://localhost:6379/0
```

### `migrate.databases[]`

- `name`: logical database name used in requests
- `source`: source reference that resolves to a migration source URL (`file://...`, `github://...`)
- `url`: source reference that resolves to a database URL (`pgx5://...`)

If `database` is not found at runtime, transports map the error to gRPC `NotFound` / HTTP `404`.

### `health`

- `duration`: health check interval
- `timeout`: timeout budget per check

The service registers:

- `noop` and `online` checks
- one migration checker per configured database

### `redis`

- `url`: source reference that resolves to a Redis URL (`redis://...`)
- Required for distributed migration locking

## Running locally

### First-time repository setup

Most `make` targets depend on scripts from `bin/` (git submodule):

```bash
git submodule sync
git submodule update --init
```

### Dependencies and build

```bash
make dep
make build
```

This produces `./migrieren`.

If Go reports inconsistent vendoring, run:

```bash
make dep
```

### Start server

```bash
cd test
../migrieren server -i file:.config/server.yml
```

Dev loop with hot reload:

```bash
make dev
```

Default addresses from test config:

- HTTP: `localhost:11000`
- gRPC: `localhost:12000`

## API usage

### gRPC

```bash
grpcurl -plaintext \
  -d '{"database":"postgres","version":1}' \
  localhost:12000 \
  migrieren.v1.Service/Migrate
```

### HTTP RPC facade

```bash
curl -sS \
  -H 'content-type: application/json' \
  -d '{"database":"postgres","version":1}' \
  http://localhost:11000/migrieren.v1.Service/Migrate
```

Request fields:

- `database` (string): logical database name from config
- `version` (uint64): target migration version

Response fields:

- `meta` (map): request metadata/observability attributes
- `migration.database`
- `migration.version`
- `migration.logs` (array of log strings)

Error mapping:

- Unknown `database`: gRPC `NotFound`, HTTP `404`
- Invalid config/migration/lock/driver failures: gRPC `Internal`, HTTP `500`

## Development commands

```bash
make lint
make specs
make features
make benchmarks
make proto-generate
make proto-lint
```

Run `make help` for the full command list.

## Repository layout

- `main.go`: CLI app entrypoint
- `internal/cmd`: command wiring (`server`)
- `internal/config`: root config model and DI wiring
- `internal/redis`: Redis client and Redsync wiring
- `internal/migrate`: core migration execution
- `internal/api/migrate`: transport-facing adapter (`database + version` API)
- `internal/api/v1/transport/grpc`: gRPC handlers
- `internal/api/v1/transport/http`: HTTP RPC route mapping
- `internal/health`: health registrations/checkers
- `api/migrieren/v1/service.proto`: API contract

## CI notes

CI runs Postgres, Valkey (Redis-compatible), and Mimir containers before executing lint, security, tests, and coverage.

## Changelog

See `CHANGELOG.md` for release notes and changes.
