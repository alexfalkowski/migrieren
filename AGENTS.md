# AGENTS.md

This repository is a Go service (with a small Ruby test/client harness) for running database migrations via a gRPC/HTTP API. It uses a `bin/` git submodule for build tooling.

## Quick start

1. **Init submodule tooling** (required for most `make` targets):

   ```sh
   git submodule sync
   git submodule update --init
   # or: make submodule
   ```

2. **Install deps / vendor** (the Makefiles commonly run with `-mod vendor`):

   ```sh
   make dep
   ```

3. **Build and run** (CLI server command is `server`):

   ```sh
   make build
   ./migrieren server -i file:test/.config/server.yml
   ```

## Essential commands

Run `make help` to see the full list of targets (help output is generated from Makefile comments).

### Dependencies / cleanup

- `make dep`: Go + Ruby deps (Go: download/tidy/vendor; Ruby: `make -C test dep`).
- `make clean`: cleans downloaded deps via `bin/build/go/clean`.
- `make clean-dep`: cleans Go caches and Ruby deps.
- `make clean-lint`: clears golangci-lint cache.

### Build / dev

- `make build`: builds `./migrieren` (adds tags `netgo`).
- `make build-test`: builds a *test binary* with build tag `features`.
- `make dev`: runs `air` with `make dep build` and executes:
  
  ```sh
  ./migrieren server -i file:test/.config/server.yml
  ```

### Test

- `make specs`: Go tests via `gotestsum` (writes JUnit + coverage into `test/reports/`).
- `make features`: runs Ruby feature tests (`make -C test features`).
- `make benchmarks`: runs Ruby benchmarks (`make -C test benchmarks`).
- `make coverage`: generates HTML + func coverage reports from `test/reports/final.cov`.

Notes:
- `main_test.go` is guarded by `//go:build features`.
- Go test invocations in Makefiles typically use `-race -vet=off -mod vendor`.

### Lint / format / security

- `make lint`: Go lint + Ruby lint + proto lint.
- `make fix-lint`: attempts to auto-fix lint issues (Go/Ruby/proto).
- `make format`: formats Go + Ruby + proto.
- `make sec`: runs `govulncheck -test ./...`.
- `make trivy-repo`: Trivy repo scan (scripted via `bin/`).

### Protobuf / API

The API contract lives under `api/` and is managed with `buf`.

Top-level convenience targets:
- `make proto-generate`: `make -C api generate` (runs `buf generate`).
- `make proto-breaking`: `make -C api breaking` (runs `buf breaking ...`).
- `make proto-lint`: `make -C api lint`.
- `make proto-format`: `make -C api format`.

Generation configuration (`api/buf.gen.yaml`) uses remote plugins:
- Go protobuf → output into `api/` (source-relative)
- Go gRPC → output into `api/` (source-relative)
- Ruby protobuf + gRPC → output into `test/lib/`

### Local environment helpers

- `make start` / `make stop`: starts/stops a Docker-based dev environment via `bin/build/docker/env`.

CI (CircleCI) runs with containers including Postgres (`localhost:5432`) and Grafana Mimir (`localhost:9009`), and waits for them before running `make features`/`make benchmarks`.

## Repo layout

- `main.go`: entrypoint; wires CLI command(s) using `go-service/v2/cli`.
- `internal/cmd/`: registers the `server` command and composes DI modules.
- `internal/config/`: root service config type; embeds `go-service/v2/config.Config`.
- `internal/migrate/`: core migration logic built on `github.com/golang-migrate/migrate/v4`.
  - `internal/migrate/database/`: database driver wiring (currently supports `pgx5://...`).
  - `internal/migrate/source/`: migration source wiring (file/github drivers imported).
  - `internal/migrate/telemetry/logger/`: captures migration logs into memory.
- `internal/api/`:
  - `internal/api/migrate/`: transport-facing migrator (reads sources/urls from config + filesystem).
  - `internal/api/v1/transport/grpc/`: gRPC server and handlers.
  - `internal/api/v1/transport/http/`: HTTP RPC routing to gRPC handlers.
- `api/`: protobuf definitions and generated Go stubs.
- `test/`: Ruby-based feature tests + generated Ruby protobuf stubs.
  - `test/.config/server.yml`: dev/test server configuration used by `make dev` and features.

## Code patterns and conventions (observed)

### Dependency injection / modules

Modules are composed with `github.com/alexfalkowski/go-service/v2/di` (Fx-style):
- Each subsystem exposes a `var Module = di.Module(...)` (e.g. `internal/migrate/module.go`, `internal/api/v1/module.go`).
- Constructors use `di.Constructor(...)`.
- Registrations (e.g. gRPC service registration, HTTP routes) use `di.Register(...)`.

### Error handling

- Domain errors are declared as package vars (e.g. `internal/migrate/migrate.go`) and mapped in transports.
- gRPC layer maps "not found" to `codes.NotFound`, everything else to `codes.Internal` (`internal/api/v1/transport/grpc/grpc.go`).
- Migration failures attach attributes into request `meta` (e.g. `meta.WithAttribute(ctx, "migrateError", meta.Error(err))`).

### Formatting / lint

- `.editorconfig` specifies:
  - Go: tabs, indent size 4
  - Makefiles: tabs
- `.golangci.yml`:
  - Enables `default: all` with an explicit disable list.
  - Excludes generated `.pb*` files from lint/format.
  - Sets `lll.line-length: 140`.

### Build tags

Some tests/targets rely on the `features` build tag (notably `make build-test` and `main_test.go`).

## Testing (Ruby feature harness)

Ruby test harness lives in `test/` and uses:
- `cucumber` features under `test/features/**`.
- Generated Ruby protobuf stubs under `test/lib/migrieren/v1/*`.
- A shared helper module in `test/lib/migrieren.rb` that builds HTTP and gRPC clients.

Feature files include `@startup` and `@clean` tags (see `test/features/**`).

## CI notes (CircleCI)

The main CI job (`.circleci/config.yml`) runs roughly:
- `make clean && make dep`
- `make lint`
- `make proto-breaking`
- `make sec`
- `make trivy-repo`
- `make features`
- `make benchmarks`
- `make analyse`
- `make coverage && make codecov-upload`

If you add/modify build steps, ensure they fit this flow.

## Common gotchas

- **Submodule dependency**: most Make targets call scripts under `bin/`; ensure `bin/` submodule is initialized.
- **Vendoring**: multiple Go targets run with `-mod vendor`. Run `make dep` after changing Go deps.
- **Config-driven sources/urls**: migration `source` and DB `url` are loaded via `go-service/v2/os.FS` from config values like `file:secrets/pg` (see `test/.config/server.yml`).
- **DB scheme**: DB driver currently expects `pgx5://...` and rewrites to `postgres://...` internally (`internal/migrate/database/database.go`).
