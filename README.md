![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/migrieren.svg?style=shield)](https://circleci.com/gh/alexfalkowski/migrieren)
[![codecov](https://codecov.io/gh/alexfalkowski/migrieren/graph/badge.svg?token=R2OD8WIKD0)](https://codecov.io/gh/alexfalkowski/migrieren)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/migrieren)](https://goreportcard.com/report/github.com/alexfalkowski/migrieren)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/migrieren.svg)](https://pkg.go.dev/github.com/alexfalkowski/migrieren)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Migrieren

Migrieren is a small Go service that runs database migrations via a **gRPC API** with an **HTTP RPC façade**.

It’s designed to let you centralize migrations in one place (instead of duplicating migration tooling across multiple application frameworks), while still using “native” database migrations (e.g. SQL scripts) via [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate).

- **Primary runtime:** Go service
- **API contract:** Protobuf/gRPC (see `api/migrieren/v1/service.proto`)
- **Feature tests / harness:** Ruby + Cucumber under `test/`

## Why a migration service?

Most frameworks can run migrations, but they’re often coupled to an ORM. In practice, teams frequently prefer:
- migrations written in SQL (or database-native tooling),
- migrations managed once for many services,
- a consistent operational model and observability for migrations.

Migrieren focuses on “run migration X for database Y” as an API, suitable for orchestration (Kubernetes init containers, CI jobs, etc.).

## API overview

The gRPC contract is defined here:

- `api/migrieren/v1/service.proto`

At a high level:

- `migrieren.v1.Service/Migrate` migrates a configured database to a target version and returns migration logs plus response metadata.

## Supported sources and databases

Migrieren uses `golang-migrate` under the hood. In this repo, the source/database drivers are wired internally:

- Sources: commonly `file://...` and GitHub sources (depending on driver wiring).
- Databases: Postgres via pgx (`pgx5://...`) is supported (the service rewrites internally to the expected migrate URL scheme).

Exactly which drivers are enabled is determined by the internal driver registration code.

## Configuration

Migrieren is configured via a YAML config file passed to the `server` command.

A sample development/test config exists in:

- `test/.config/server.yml`

### Migration configuration (`migrate`)

You can configure multiple databases. Each database has:
- `name`: a unique logical name referenced by the API (`database` in the request)
- `source`: how to resolve the migration source (read via the service filesystem abstraction)
- `url`: how to resolve the database connection URL (also read via the filesystem abstraction)

Example:

```/dev/null/example.server.yml#L1-28
migrate:
  databases:
    - name: db1
      # These values are read through the service filesystem abstraction.
      # Common patterns are `file:relative/path` or similar, depending on config.
      source: file:test/.migrations/db1
      url: file:test/.secrets/db1.url

    - name: db2
      source: file:test/.migrations/db2
      url: file:test/.secrets/db2.url
```

Notes:
- The service reads the contents of `source` and `url` using its filesystem abstraction (so `source` typically points to something that ultimately yields a migration source URL like `file://...`, and `url` yields a DB URL like `pgx5://...`).
- If a database name is not found in the config, the service maps that to a “not found” condition at the transport layer.

### Health configuration (`health`)

Health checks include per-database checks. Configure basic health polling:

```/dev/null/example.health.yml#L1-4
health:
  duration: 1s  # how often to check
  timeout: 1s   # per-check timeout
```

## Running the server

### Build

Most workflows are driven by `make`:

```/dev/null/example.build.sh#L1-2
make dep
make build
```

This produces a `./migrieren` binary in the repo root.

### Start the server

The CLI entrypoint registers a `server` command. Example using the repo’s dev/test configuration:

```/dev/null/example.run.sh#L1-1
./migrieren server -i file:test/.config/server.yml
```

Development hot-reload (if you have `air` available):

```/dev/null/example.dev.sh#L1-1
make dev
```

## Using the API

### gRPC (conceptual)

The service exposes `migrieren.v1.Service/Migrate`.

Request fields:
- `database`: logical name (must match a configured database entry)
- `version`: target version (uint64)

Response fields:
- `meta`: key/value metadata (used by the service for observability)
- `migration`: includes `database`, `version`, and `logs`

### HTTP façade

The HTTP façade routes RPC-like endpoints. For the v1 migrate call:

- `POST /migrieren.v1.Service/Migrate`
- JSON body: `{ "database": "...", "version": 123 }`

Example:

```/dev/null/example.http.txt#L1-9
POST http://localhost:11000/migrieren.v1.Service/Migrate
Content-Type: application/json

{
  "database": "db1",
  "version": 1
}
```

## Deployment guidance

In containerized environments (e.g. Kubernetes), common patterns are:

- Run Migrieren as a shared service per bounded context (or per environment).
- Run migrations during deploy via:
  - a CI job step, or
  - a Kubernetes init container that calls the Migrieren API before the main workload starts.

The “best” approach depends on your tolerance for coupling deploys to schema changes and your rollback strategy.

## Development

### Repository structure

This repository follows the common Go project layout. Key locations:

- `main.go`: CLI wiring
- `internal/cmd/server.go`: `server` command registration
- `api/`: protobuf contract + generation (managed by `buf`)
- `internal/migrate/`: core migration logic (wraps `golang-migrate`)
- `internal/api/v1/transport/{grpc,http}/`: gRPC + HTTP façade

### Dependencies

You’ll want:
- [Go](https://go.dev/)
- [Ruby](https://www.ruby-lang.org/en/) (for feature tests in `test/`)

### Setup

```/dev/null/example.setup.sh#L1-1
make setup
```

### Tests

Go tests:

```/dev/null/example.specs.sh#L1-1
make specs
```

Ruby feature tests (Cucumber):

```/dev/null/example.features.sh#L1-1
make features
```

### Protobuf generation

```/dev/null/example.proto.sh#L1-2
make proto-generate
make proto-lint
```

## Changelog

See `CHANGELOG.md` for release notes and changes.