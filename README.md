![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/migrieren.svg?style=shield)](https://circleci.com/gh/alexfalkowski/migrieren)
[![codecov](https://codecov.io/gh/alexfalkowski/migrieren/graph/badge.svg?token=R2OD8WIKD0)](https://codecov.io/gh/alexfalkowski/migrieren)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/migrieren)](https://goreportcard.com/report/github.com/alexfalkowski/migrieren)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/migrieren.svg)](https://pkg.go.dev/github.com/alexfalkowski/migrieren)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🧭 Migrieren

Migrieren is a Go service that runs database schema migrations through a gRPC API with an HTTP RPC façade.

The service wraps [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) so callers can ask it to migrate a named database to a target version without embedding migration logic into every application.

## 🧰 What the service does

- Exposes RPCs to migrate configured databases, apply all pending migrations,
  and inspect migration status.
- Looks up a logical database name in configuration.
- Reads the migration source URL and database URL through the service filesystem abstraction.
- Executes migrations with `golang-migrate`.
- Reports the current migration version and dirty state for a configured
  database.
- Returns migration logs and request metadata to the caller.
- Publishes HTTP and gRPC health checks plus Prometheus-style metrics when configured through the shared service runtime.

## 🚚 Supported drivers

The currently wired drivers are defined in code:

- Migration sources:
  - `file://...`
  - `github://...`
- Databases:
  - Postgres via `pgx5://...`

Database URLs use the `pgx5://` scheme in config and secrets. Internally the service rewrites that to the driver URL format expected by the Postgres migrate driver.

> [!NOTE]
> Driver support comes from the drivers imported by this service. Adding another source or database kind requires wiring the driver in code, not only changing YAML.

## ✅ Prerequisites

- Go, using the version declared in `go.mod`.
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

## 🚀 Quick start

Install dependencies, build the binary, and start the server with the checked-in development/test config:

```sh
make dep
make build
cd test
../migrieren server -config file:.config/server.yml
```

This builds `./migrieren` in the repository root. The checked-in test config is
run from `test/` because its `file:...` references point at `test/secrets/` and
its file migration sources are relative to that working directory.

> [!WARNING]
> The checked-in test config points Postgres at `localhost:5433` and OTLP tracing at `http://localhost:4318/v1/traces`. Start the local feature-test services or provide your own config before expecting migrations, health checks, and tracing to succeed.
> In the feature harness, `localhost:5433` is a nonnative fault-injection proxy
> to a backing Postgres instance on `localhost:5432`; both ports must line up
> with either the local harness or your replacement config.

If you use this repository's shared local Docker environment, start and stop it with:

```sh
make start
make stop
```

Those targets manage the shared sibling Docker checkout and may require SSH access
to that repository. If you do not use them, provide an equivalent local Postgres
database named `test` with user/password `test:test` on `localhost:5432`; the
feature harness exposes it to the service through the proxy on `localhost:5433`.

For live reload during local development:

```sh
make dev
```

`make dev` uses `air` and starts the same `server` command with `test/.config/server.yml`.

> [!TIP]
> Use `make help` to list the Make targets exposed by this checkout before reaching for direct tool commands.

## 🔎 How Migrieren resolves a migration

The API accepts a logical database name such as `postgres`. The service then:

1. Looks that name up in `migrate.databases`.
2. Reads the configured `source` value through its filesystem abstraction.
3. Reads the configured `url` value through the same abstraction.
4. Passes the resolved source URL and database URL to the core migrator.

That means `source` and `url` in the YAML config are usually references to files or other resolvable inputs, not the final literal migration/source strings themselves.

The resolver accepts the shared go-service source-string forms:

- `env:NAME`: reads the value of environment variable `NAME`. An omitted or
  unset variable fails resolution; an explicitly empty variable resolves to an
  empty value.
- `file:<path>`: reads the file at `<path>`. Relative paths are resolved from
  the server process working directory.
- any other value: used as the literal migration source URL or database URL.

## ⚙️ Configuration

Migrieren is configured through the `server` command input file. The repository includes a representative config at `test/.config/server.yml`.

### 🧩 Minimal configuration shape

At minimum, you need:

- `health.duration`
- `health.timeout`
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

`health.duration` sets how often the service evaluates its noop, online, and
per-database health registrations. `health.timeout` bounds the online check and
each database target's source validation and database ping. Migration request
timeouts are configured separately by the transport and database URL options.

### 🗄️ Migration database entries

`migrate.databases` must contain at least one entry. Each entry must have a
unique `name` and all of these fields:

- `name`: the logical database name used by the API request.
- `source`: how to resolve the migration source URL.
- `url`: how to resolve the database connection URL.

`migrate.logs` is optional. `migrate.logs.max` bounds how many migration log
lines each `Migrate`/`ApplyMigrations` response returns. Omitting it or setting
it to `0` selects the default of 100; only positive values set a custom cap, and
there is no value that disables returned logs. Negative values are invalid.
When the limit is exceeded, the oldest lines are dropped and the first returned
entry is a `migration logs truncated (showing last N of M)` marker.

In the test setup, the referenced secret files contain the values consumed by
the migrator. The GitHub source file is generated locally and is not checked
in:

```text
# test/secrets/source
file://migrations

# test/secrets/github (generated by make github-source)

# test/secrets/pg
pgx5://test:test@localhost:5433/test?sslmode=disable&x-migrations-table=migrieren_schema_migrations&x-statement-timeout=5000&x-multi-statement=true&x-multi-statement-max-size=2097152
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

### 🐙 GitHub migration sources

GitHub sources use this URL shape:

```text
github://[user:personal-access-token@]owner/repository[/path][#ref]
```

- Without user information, the repository is read without authentication.
- For an authenticated source, `make github-source` reads
  `MIGRIEREN_GITHUB_TOKEN` and writes the credential-bearing URL to the ignored
  `test/secrets/github` file. The username in that URL is ignored.
- The optional path selects a migration directory within the repository.
- The optional fragment selects a branch, tag, or commit SHA. When omitted,
  GitHub's default branch is used.

The migration directory is listed when the source opens. Keep the token in a
local or CI secret; `make features` generates the ignored source file before
running the feature suite.

### 🐘 Postgres URL options

Postgres targets use `pgx5://...` URLs. In addition to normal Postgres
connection parameters, Migrieren recognizes these migration-driver query
parameters:

- `x-migrations-table`: migration table name. When omitted, the upstream driver
  default is used.
- `x-migrations-table-quoted`: boolean. When `true`, `x-migrations-table` must
  include surrounding double quotes.
- `x-statement-timeout`: statement timeout in milliseconds for migration
  operations. A migration request deadline can lower this timeout. Health checks
  instead use `health.timeout` to bound `PingContext` and do not use this
  statement timeout.
- `x-multi-statement`: boolean enabling multi-statement migrations.
- `x-multi-statement-max-size`: byte limit for multi-statement migrations.
  Empty or non-positive values use the upstream driver default.

Malformed boolean or integer values reject the configured database URL.

The upstream `golang-migrate` pgx driver can acquire its migration advisory lock
or inspect migration metadata while constructing or using the database driver.
Those upstream paths do not consistently expose Migrieren's request context, so
request cancellation and deadlines should not be treated as a strict bound for
every advisory-lock wait or status inspection. Migrieren relies on the upstream
driver behavior here rather than maintaining a local pgx driver fork, and can
tighten this once context-aware migrate v5 driver APIs are available.

### 🔐 Authentication and access control

Migrieren can require an authenticated, authorized caller on both the gRPC API
and the HTTP façade. This reuses the shared `go-service` transport wiring:

- `transport.<http|grpc>.token` enables server-side token verification. The
  supported kinds are `jwt`, `paseto`, and `ssh`. When present, application RPCs
  require a valid `Authorization: Bearer <token>` value; the operation routes
  (`/migrieren/livez`, `/migrieren/readyz`, `/migrieren/healthz`, metrics, and
  the gRPC health service) are always exempt.
- `transport.access` enables a Casbin authorization policy that is enforced
  after verification. The policy subject is the verified token subject, the
  object is the transport-scoped method, and the action is `invoke`.

The verification audience is bound to the request so a token cannot be replayed
against another endpoint:

- HTTP audience is `"POST <path>"`, for example `POST /migrieren.v1.Service/Status`.
- gRPC audience is the full method, for example `/migrieren.v1.Service/Status`.

The transport-scoped access objects are prefixed by transport, for example
`http:POST /migrieren.v1.Service/Status` and `grpc:/migrieren.v1.Service/Status`.

The checked-in test config wires the `ssh` kind on both transports. SSH tokens
fix `sub == kid == key`, so the signing key id is the verified subject that the
policy is evaluated against:

```yaml
transport:
  access:
    model: file:secrets/access_model
    policy: file:secrets/access_policy
  http:
    address: tcp://:11000
    token:
      kind: ssh
      ssh:
        key: migrieren
        exp: 1h
        keys:
          migrieren:
            public: file:secrets/ssh_migrieren.pub
          guest:
            public: file:secrets/ssh_guest.pub
  grpc:
    address: tcp://:12000
    token:
      kind: ssh
      ssh:
        key: migrieren
        exp: 1h
        keys:
          migrieren:
            public: file:secrets/ssh_migrieren.pub
          guest:
            public: file:secrets/ssh_guest.pub
```

The `secrets/access_policy` file grants the `migrieren` subject `invoke` on every
method for both transports. The `guest` key is trusted for verification but is
absent from the policy, so a `guest`-signed token authenticates yet is denied.

The Ruby feature harness mints matching tokens with `nonnative`
(`Migrieren.auth_token`) and attaches them automatically: an HTTP `post`
override adds a route-scoped Bearer header, and a gRPC client interceptor adds
route-scoped `authorization` metadata. See `test/features/v1/transport/*/auth.feature`.

### 🧪 About the checked-in test config

`test/.config/server.yml` intentionally contains both valid and invalid database definitions:

- `postgres` and `github` are used for successful migration scenarios.
- `timeout` is used for deadline and cancellation scenarios.
- `logs` is used for bounded migration log scenarios.
- `missing_source`, `invalid_source`, `missing_url`, `invalid_url`,
  `invalid_db`, `invalid_quoted_table`, `invalid_incomplete_quoted_table`, and
  `invalid_port` exist to exercise failure paths in feature tests.

> [!IMPORTANT]
> The checked-in config is a test fixture. Use it for local development and feature coverage, but do not treat it as a healthy production config.

That is why the checked-in config is useful for development and testing, but it is not a "healthy production example" as-is.

## 🔌 API

The protobuf contract lives at `api/migrieren/v1/service.proto`.

### 🔷 gRPC

The service exposes:

- `migrieren.v1.Service/Migrate`
- `migrieren.v1.Service/ApplyMigrations`
- `migrieren.v1.Service/PlanMigrations`
- `migrieren.v1.Service/Status`
- `migrieren.v1.Service/ListDatabases`

Request:

- `database`: logical database name from config.
- `version`: target migration version as `uint64`; must be between `1` and the server-supported signed integer maximum.

Response:

- `meta`: request metadata emitted by the service runtime.
- `migration.database`: echoed database name.
- `migration.version`: echoed target version.
- `migration.logs`: in-memory migration log lines collected during execution. Returned logs are capped at a configured maximum (`migrate.logs.max`, default 100) and start with a `migration logs truncated (showing last N of M)` marker when older lines were discarded.

Conceptual request:

```protobuf
database: "postgres"
version: 1
```

### ⏫ Apply all pending migrations

`migrieren.v1.Service/ApplyMigrations` applies all pending up migrations for a
configured database and reports the resulting migration version. If the database
is already current, the request succeeds as a no-op and reports the current
version.

Request:

- `database`: logical database name from config.

Response:

- `meta`: request metadata emitted by the service runtime.
- `migration.database`: echoed database name.
- `migration.version`: resulting migration version after applying pending
  migrations.
- `migration.logs`: in-memory migration log lines collected during execution.
  Returned logs are capped at a configured maximum (`migrate.logs.max`, default
  100) and start with a `migration logs truncated (showing last N of M)` marker
  when older lines were discarded.

Conceptual request:

```protobuf
database: "postgres"
```

> [!NOTE]
> ApplyMigrations uses the upstream `golang-migrate` v4 all-pending up path.
> Request cancellation and deadlines are honored around the migration call, but
> strict cancellation still depends on upstream context-aware driver support in
> some database paths.
>
> Migrieren intentionally does not expose step-based up/down, rollback-by-step,
> all-down, or version-zero migration APIs. Use `Migrate` when the caller knows
> the target version, or `ApplyMigrations` to converge to the latest available
> up migration.

### 🧭 Migration plan

`migrieren.v1.Service/PlanMigrations` reports current status and the migration
versions that would run for a configured database without applying migration
files. It opens the configured source to discover available migration versions,
but it does not read migration bodies, execute migrations, or expose configured
source strings, database URL strings, resolved URLs, or secret values.

Request:

- `database`: logical database name from config.
- `target_version` (optional): explicit version to preview. When omitted, the
  service preserves latest-up planning. When supplied, it must be greater than
  `0`, within the server-supported signed integer range, and present in the
  migration source.

Response:

- `meta`: request metadata emitted by the service runtime.
- `plan.status`: current migration status for the database.
- `plan.latest_version`: highest migration version available from the
  configured source, or `0` when the source has no migrations.
- `plan.target_version`: effective planned version. It equals an explicit
  request target when supplied; otherwise it is the latest source version for
  an up plan or the current status version when no up migrations can apply.
- `plan.direction`: `MIGRATION_DIRECTION_UP`, `MIGRATION_DIRECTION_DOWN`, or
  `MIGRATION_DIRECTION_NONE`.
- `plan.pending_versions`: source versions that would run in execution order.
  A down plan includes the current version and excludes the target version. For
  an omitted target on a dirty database, this preserves legacy behavior by
  listing later discovered source versions while `plan.direction` remains
  `MIGRATION_DIRECTION_NONE`; it is informational rather than executable.

Conceptual request:

```protobuf
database: "postgres"
target_version: 1
```

> [!NOTE]
> PlanMigrations is informational and does not reserve the migration source or
> database for a later request. Migrieren intentionally does not expose
> step-based up/down, rollback-by-step, all-down, or version-zero migration
> APIs. Use `Migrate` for an explicit target version and `ApplyMigrations` for
> latest. An explicit plan reports the same safe migration failure as `Migrate`
> when the current database state is dirty, its requested target is absent from
> the source, or a recorded current version is absent from the source.

### 🔎 Migration status

`migrieren.v1.Service/Status` reports the current migration state for a
configured database without applying migration files.

Request:

- `database`: logical database name from config.

Response:

- `meta`: request metadata emitted by the service runtime.
- `status.database`: echoed database name.
- `status.version`: current clean or dirty migration version. It is `0` when
  `status.state` is `MIGRATION_STATE_UNAPPLIED`, and can also be `0` with
  `MIGRATION_STATE_DIRTY` when migration metadata records a dirty recovery state
  without a version. Always use `status.state` to distinguish those cases.
- `status.state`: one of `MIGRATION_STATE_UNAPPLIED`,
  `MIGRATION_STATE_CLEAN`, or `MIGRATION_STATE_DIRTY`.

Conceptual request:

```protobuf
database: "postgres"
```

> [!NOTE]
> Status is non-migrating, but strict request cancellation depends on upstream
> `golang-migrate` v4 context support. Some database-driver inspection paths do
> not accept Migrieren's request context until the upstream project provides
> context-aware driver APIs.

Migrieren intentionally does not expose a public force-version or dirty-state
repair RPC. Dirty status means an operator must investigate the failed
migration and repair the database state outside this service before any
migration metadata is changed. A remote repair endpoint would be destructive
metadata surgery that Migrieren cannot validate generically, because the service
cannot prove the real schema or data was repaired correctly for every
application.

### 📋 Database discovery

`migrieren.v1.Service/ListDatabases` reports configured logical database names
from `migrate.databases` in config order.

Request:

- Empty request.

Response:

- `meta`: request metadata emitted by the service runtime.
- `databases[].name`: configured logical database name.

The response does not include configured `source` strings, configured `url`
strings, resolved database URLs, resolved migration source URLs, or secret
values.

### 🌐 HTTP façade

The HTTP RPC façade exposes the same operation at:

- `POST /migrieren.v1.Service/Migrate`
- `POST /migrieren.v1.Service/ApplyMigrations`
- `POST /migrieren.v1.Service/PlanMigrations`
- `POST /migrieren.v1.Service/Status`
- `POST /migrieren.v1.Service/ListDatabases`

The checked-in test config requires a route-scoped Bearer value for application
RPCs. After `make dep`, run the examples below from `test/` and define this
helper once; it uses the throwaway test signing key to mint a header for the
exact route passed to it without placing the token in curl's process arguments:

```sh
http_authorization() {
  bundle exec ruby -Ilib -rnonnative -rmigrieren \
    -e 'puts "Authorization: #{Migrieren.http_authorization(ARGV.fetch(0))}"' "$1"
}
```

If your own config omits token verification and access control, the
`Authorization` header is unnecessary.

Example:

```http
POST http://localhost:11000/migrieren.v1.Service/Migrate
Content-Type: application/json
Authorization: Bearer <route-scoped-token>

{
  "database": "postgres",
  "version": 1
}
```

Copy-paste request against the local HTTP façade:

```sh
http_authorization '/migrieren.v1.Service/Migrate' |
  curl -sS -X POST http://localhost:11000/migrieren.v1.Service/Migrate \
  -H 'Content-Type: application/json' \
  --header @- \
  -d '{"database":"postgres","version":1}'
```

Copy-paste apply-all request against the local HTTP façade:

```sh
http_authorization '/migrieren.v1.Service/ApplyMigrations' |
  curl -sS -X POST http://localhost:11000/migrieren.v1.Service/ApplyMigrations \
  -H 'Content-Type: application/json' \
  --header @- \
  -d '{"database":"postgres"}'
```

Copy-paste plan request against the local HTTP façade:

```sh
http_authorization '/migrieren.v1.Service/PlanMigrations' |
  curl -sS -X POST http://localhost:11000/migrieren.v1.Service/PlanMigrations \
  -H 'Content-Type: application/json' \
  --header @- \
  -d '{"database":"postgres","target_version":1}'
```

Copy-paste status request against the local HTTP façade:

```sh
http_authorization '/migrieren.v1.Service/Status' |
  curl -sS -X POST http://localhost:11000/migrieren.v1.Service/Status \
  -H 'Content-Type: application/json' \
  --header @- \
  -d '{"database":"postgres"}'
```

Copy-paste database discovery request against the local HTTP façade:

```sh
http_authorization '/migrieren.v1.Service/ListDatabases' |
  curl -sS -X POST http://localhost:11000/migrieren.v1.Service/ListDatabases \
  -H 'Content-Type: application/json' \
  --header @- \
  -d '{}'
```

### 🚦 Error mapping

Transport behavior is intentionally simple:

- Missing or invalid credentials when token verification is enabled:
  - gRPC: `Unauthenticated`
  - HTTP: `401`
- Authenticated subject denied by the access policy:
  - gRPC: `PermissionDenied`
  - HTTP: `403`
- Migration version outside the supported range:
  - gRPC: `InvalidArgument`
  - HTTP: `400`
- Unknown database name:
  - gRPC: `NotFound`
  - HTTP: `404`
- Configuration, source, database, or migration failures:
  - gRPC: `Internal`
  - HTTP: `500`
- Request canceled by the caller:
  - gRPC: `Canceled`
  - HTTP: `499` (`Client Closed Request`, non-standard)
- Request deadline exceeded:
  - gRPC: `DeadlineExceeded`
  - HTTP: `504`

The core migrator also treats `migrate.ErrNoChange` as a successful no-op and still returns any accumulated logs.

For Migrate, ApplyMigrations, PlanMigrations, and Status requests that pass
request validation and then fail, the server adds safe failure diagnostics when
they are available. gRPC exposes them as trailers and HTTP exposes them as
response headers without changing the HTTP status code or error body shape.
This includes unknown database names, source/database setup and
reference-resolution failures, and core migration failures; invalid version
values return `InvalidArgument`/`400` before this diagnostic path runs.

- `migration-error`: one of `not_found`, `canceled`, `deadline_exceeded`,
  `invalid_config`, `invalid_migration`, or `unknown`.
- `migration-log-count`: number of migration log lines returned.
- `migration-stage`: `source` when migration-source setup or reference
  resolution failed, or `url` when database setup or database-URL reference
  resolution failed. `url` identifies the database-setup branch; it does not
  imply malformed URL syntax. Later migration or database-inspection failures
  omit this value.
- `migration-log-last`: last migration log line when logs were captured.

HTTP uses the same names as response headers:

- `Migration-Error`
- `Migration-Log-Count`
- `Migration-Stage`
- `Migration-Log-Last`

## 💓 Health and observability

When running with the shared service runtime, Migrieren exposes:

- HTTP health endpoints:
  - `/migrieren/healthz`
  - `/migrieren/livez`
  - `/migrieren/readyz`
- HTTP metrics:
  - `/migrieren/metrics`
- gRPC health checks for `migrieren.v1.Service`

The probe contracts are deliberately different:

| Surface | Checks observed | Operational meaning |
| --- | --- | --- |
| HTTP `/migrieren/livez` | `noop` | Process liveness only. |
| HTTP `/migrieren/readyz` | `noop` | Process readiness only; it does not gate on migration dependencies. |
| HTTP `/migrieren/healthz` | `online` and every configured database target | Public connectivity plus source validation and database reachability. |
| gRPC health for `migrieren.v1.Service` | `noop` | Service-process health independent of migration dependencies. |

The shared `online` check uses its default public connectivity endpoints. Use
`healthz` for dependency health, not as a liveness probe that should remain
stable during an egress or database outage.

There is an important detail in the checked-in test config:

- gRPC health for `migrieren.v1.Service` reports `SERVING`.
- HTTP `/migrieren/livez` and `/migrieren/readyz` report healthy.
- HTTP `/migrieren/healthz` is expected to be unhealthy because the test config deliberately registers invalid database entries for failure-path coverage.

Health probes use internal registration names such as `noop` and `online`.
Treat those names as reserved for health wiring rather than as migration
database names.

The sample test config also enables:

- text logging,
- Prometheus metrics, and
- OTLP tracing with `http://localhost:4318/v1/traces`.

GitHub migration sources have one health-check exception: `/migrieren/healthz`
parses a configured `github://` source but does not open the remote repository
during the health check. The remote source is opened during migration
execution, so GitHub reachability failures can still appear on a migration
request after health has reported on the configured target.

## 🛠️ Development

### 📋 Common commands

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
make proto-stale
```

What those do in this repository:

- `make dep`: installs Go dependencies, runs `go mod tidy`, vendors modules, and installs Ruby gems for `test/`.
- `make build`: builds the release binary `./migrieren`.
- `make build-test`: builds a test binary with the `features` build tag.
- `make specs`: runs Go tests with `gotestsum`, `-race`, and vendored dependencies.
- `make features`: builds the test binary and runs the Ruby/Cucumber feature suite in `test/`.
- `make benchmarks`: builds the release binary and runs the benchmark-tagged Ruby harness.
- `make lint`: lints Go, the Ruby test harness, and protobuf definitions.
- `make sec`: runs `govulncheck` and the Trivy repository scan.
- `make coverage`: creates HTML and function coverage reports under `test/reports/`.

### 🧪 Local test harness expectations

The Ruby harness under `test/` assumes:

- HTTP server on `http://localhost:11000`
- gRPC server on `localhost:12000`
- Postgres reachable by the service on `localhost:5433`
- A backing local Postgres instance on `localhost:5432` for the nonnative proxy
  and direct harness cleanup/assertion helpers

The feature harness process wiring lives in `test/nonnative.yml`.

### 🗂️ Repository layout

Key locations:

- `main.go`: CLI entrypoint.
- `internal/cmd/server.go`: registers the `server` command.
- `internal/config/config.go`: top-level config composition.
- `internal/migrate/`: core migration engine and driver wiring.
- `internal/api/migrate/`: API-facing adapter that resolves database names through config.
- `internal/api/v1/migrate/`: versioned API contract shared by HTTP and gRPC transports.
- `internal/api/v1/transport/grpc/`: gRPC server implementation.
- `internal/api/v1/transport/http/`: HTTP RPC route registration.
- `internal/health/`: health registration and database-specific health checks.
- `api/migrieren/v1/service.proto`: protobuf contract.
- `test/`: Ruby feature-test harness, migrations, and local test fixtures.

## 🧬 Protobuf workflow

The API contract is managed with `buf`.

Generate code:

```sh
make proto-generate
```

Lint and breaking-change checks:

```sh
make proto-lint
make proto-breaking
make proto-stale
```

`make proto-stale` verifies generated protobuf outputs are current; CI runs it
after the breaking-change check.

Generation is configured in `api/buf.gen.yaml` and currently writes:

- Go protobuf and gRPC files into `api/`
- Ruby protobuf and gRPC files into `test/lib/`

> [!CAUTION]
> Do not hand-edit generated protobuf stubs. Update `api/migrieren/v1/service.proto` and regenerate instead.

## 🤝 Notes for contributors

- Package documentation for Go packages belongs in `doc.go`.
- The feature harness public API is primarily the Ruby helpers under `test/lib/`.
- If `go test` or `go list` starts failing because of vendoring drift, run `make dep`.
