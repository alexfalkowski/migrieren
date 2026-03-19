![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/migrieren.svg?style=shield)](https://circleci.com/gh/alexfalkowski/migrieren)
[![codecov](https://codecov.io/gh/alexfalkowski/migrieren/graph/badge.svg?token=R2OD8WIKD0)](https://codecov.io/gh/alexfalkowski/migrieren)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/migrieren)](https://goreportcard.com/report/github.com/alexfalkowski/migrieren)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/migrieren.svg)](https://pkg.go.dev/github.com/alexfalkowski/migrieren)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Migrieren

Migrieren is a Go service that runs database schema migrations through a gRPC API with an HTTP RPC façade.

The service wraps [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) so callers can ask it to migrate a named database to a target version without embedding migration logic into every application.

## What the service does

- Exposes a single RPC: `migrieren.v1.Service/Migrate`.
- Looks up a logical database name in configuration.
- Reads the migration source URL and database URL through the service filesystem abstraction.
- Executes the migration with `golang-migrate`.
- Returns migration logs and request metadata to the caller.
- Publishes HTTP and gRPC health checks plus Prometheus-style metrics when configured through the shared service runtime.

## Supported drivers

The currently wired drivers are defined in code:

- Migration sources:
  - `file://...`
  - `github://...`
- Databases:
  - Postgres via `pgx5://...`

Database URLs use the `pgx5://` scheme in config and secrets. Internally the service rewrites that to the driver URL format expected by the Postgres migrate driver.

## Prerequisites

- Go `1.26.0` or newer.
- Ruby and Bundler for the feature-test harness under `test/`.
- The `bin/` git submodule. Most `make` targets delegate into scripts under that submodule.

Initialize the submodule before relying on `make`:

```sh
git submodule sync
git submodule update --init
```

If you hit vendoring errors such as "inconsistent vendoring", refresh dependencies with:

```sh
make dep
```

## Quick start

Install dependencies, build the binary, and start the server with the checked-in development/test config:

```sh
make dep
make build
./migrieren server -i file:test/.config/server.yml
```

This builds `./migrieren` in the repository root.

For live reload during local development:

```sh
make dev
```

`make dev` uses `air` and starts the same `server` command with `test/.config/server.yml`.

## How Migrieren resolves a migration

The API accepts a logical database name such as `postgres`. The service then:

1. Looks that name up in `migrate.databases`.
2. Reads the configured `source` value through its filesystem abstraction.
3. Reads the configured `url` value through the same abstraction.
4. Passes the resolved source URL and database URL to the core migrator.

That means `source` and `url` in the YAML config are usually references to files or other resolvable inputs, not the final literal migration/source strings themselves.

## Configuration

Migrieren is configured through the `server` command input file. The repository includes a representative config at `test/.config/server.yml`.

### Minimal configuration shape

At minimum, you need:

- `migrate.databases`
- `transport.http.address`
- `transport.grpc.address`

Example:

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
transport:
  http:
    address: tcp://:11000
    timeout: 5s
  grpc:
    address: tcp://:12000
    timeout: 5s
```

### Migration database entries

Each configured database entry has:

- `name`: the logical database name used by the API request.
- `source`: how to resolve the migration source URL.
- `url`: how to resolve the database connection URL.

In the checked-in test setup, the referenced secret files contain the actual values consumed by the migrator:

```text
# test/secrets/source
file://migrations

# test/secrets/github
github://alexfalkowski/app-config/test/migrations

# test/secrets/pg
pgx5://test:test@localhost:5433/test?sslmode=disable
```

So a config entry like this:

```yaml
migrate:
  databases:
    - name: postgres
      source: file:secrets/source
      url: file:secrets/pg
```

ultimately migrates Postgres using the `file://migrations` source and the `pgx5://...` database URL resolved from those files.

### About the checked-in test config

`test/.config/server.yml` intentionally contains both valid and invalid database definitions:

- `postgres` and `github` are used for successful migration scenarios.
- `missing_source`, `invalid_source`, `missing_url`, `invalid_url`, `invalid_db`, and `invalid_port` exist to exercise failure paths in feature tests.

That is why the checked-in config is useful for development and testing, but it is not a "healthy production example" as-is.

## API

The protobuf contract lives at `api/migrieren/v1/service.proto`.

### gRPC

The service exposes:

- `migrieren.v1.Service/Migrate`

Request:

- `database`: logical database name from config.
- `version`: target migration version as `uint64`.

Response:

- `meta`: request metadata emitted by the service runtime.
- `migration.database`: echoed database name.
- `migration.version`: echoed target version.
- `migration.logs`: in-memory migration log lines collected during execution.

Conceptual request:

```protobuf
database: "postgres"
version: 1
```

### HTTP façade

The HTTP RPC façade exposes the same operation at:

- `POST /migrieren.v1.Service/Migrate`

Example:

```http
POST http://localhost:11000/migrieren.v1.Service/Migrate
Content-Type: application/json

{
  "database": "postgres",
  "version": 1
}
```

### Error mapping

Transport behavior is intentionally simple:

- Unknown database name:
  - gRPC: `NotFound`
  - HTTP: `404`
- Configuration, source, database, or migration failures:
  - gRPC: `Internal`
  - HTTP: `500`

The core migrator also treats `migrate.ErrNoChange` as a successful no-op and still returns any accumulated logs.

## Health and observability

When running with the shared service runtime, Migrieren exposes:

- HTTP health endpoints:
  - `/healthz`
  - `/livez`
  - `/readyz`
- HTTP metrics:
  - `/metrics`
- gRPC health checks for `migrieren.v1.Service`

There is an important detail in the checked-in test config:

- gRPC health for `migrieren.v1.Service` reports `SERVING`.
- HTTP `/livez` and `/readyz` report healthy.
- HTTP `/healthz` is expected to be unhealthy because the test config deliberately registers invalid database entries for failure-path coverage.

The sample test config also enables:

- text logging,
- Prometheus metrics, and
- OTLP tracing with `http://localhost:4318/v1/traces`.

## Development

### Common commands

Use `make help` to list available targets. Common ones are:

```sh
make dep
make build
make build-test
make specs
make features
make benchmarks
make lint
make format
make sec
make coverage
make proto-generate
make proto-lint
make proto-breaking
```

What those do in this repository:

- `make dep`: installs Go dependencies, runs `go mod tidy`, vendors modules, and installs Ruby gems for `test/`.
- `make build`: builds the release binary `./migrieren`.
- `make build-test`: builds a test binary with the `features` build tag.
- `make specs`: runs Go tests with `gotestsum`, `-race`, and vendored dependencies.
- `make features`: builds the test binary and runs the Ruby/Cucumber feature suite in `test/`.
- `make benchmarks`: builds the release binary and runs the benchmark-tagged Ruby harness.
- `make lint`: lints Go, the Ruby test harness, and protobuf definitions.
- `make sec`: runs `govulncheck`.
- `make coverage`: creates HTML and function coverage reports under `test/reports/`.

### Local test harness expectations

The Ruby harness under `test/` assumes:

- HTTP server on `http://localhost:11000`
- gRPC server on `localhost:12000`
- Postgres reachable on `localhost:5433`

The feature harness process wiring lives in `test/nonnative.yml`.

### Repository layout

Key locations:

- `main.go`: CLI entrypoint.
- `internal/cmd/server.go`: registers the `server` command.
- `internal/config/config.go`: top-level config composition.
- `internal/migrate/`: core migration engine and driver wiring.
- `internal/api/migrate/`: transport-facing adapter that resolves database names through config.
- `internal/api/v1/transport/grpc/`: gRPC server implementation.
- `internal/api/v1/transport/http/`: HTTP RPC route registration.
- `internal/health/`: health registration and database-specific health checks.
- `api/migrieren/v1/service.proto`: protobuf contract.
- `test/`: Ruby feature-test harness, migrations, and local test fixtures.

## Protobuf workflow

The API contract is managed with `buf`.

Generate code:

```sh
make proto-generate
```

Lint and breaking-change checks:

```sh
make proto-lint
make proto-breaking
```

Generation is configured in `api/buf.gen.yaml` and currently writes:

- Go protobuf and gRPC files into `api/`
- Ruby protobuf and gRPC files into `test/lib/`

Do not hand-edit generated protobuf stubs. Update `api/migrieren/v1/service.proto` and regenerate instead.

## Notes for contributors

- Package documentation for Go packages belongs in `doc.go`.
- The feature harness public API is primarily the Ruby helpers under `test/lib/`.
- If `go test` or `go list` starts failing because of vendoring drift, run `make dep`.

## Changelog

See `CHANGELOG.md` for release notes.
