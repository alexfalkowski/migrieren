# AGENTS.md

This repository is a Go service (plus a Ruby feature-test harness) that runs database migrations via a gRPC API and an HTTP RPC façade.

Build/test automation is primarily driven by `make`, with most targets implemented in a required `bin/` git submodule.

## 0) First check

If `bin/` is missing, most `make` targets will fail.

```sh
git submodule sync
git submodule update --init
# or: make submodule
```

Note: `.gitmodules` points `bin` at `git@github.com:alexfalkowski/bin.git` (SSH URL). You need SSH access/keys for submodule init.

## 1) Project type

- **Primary**: Go service (`go.mod` module `github.com/alexfalkowski/migrieren`, `go 1.25.0`).
- **API**: Protobuf/gRPC in `api/`, managed by `buf`.
- **Integration/feature tests**: Ruby + Cucumber in `test/`.

## 2) Essential commands (observed)

Run `make help` for the full list (help is generated from Makefile comments).

### Dependencies

- Install/refresh deps (includes vendoring):

  ```sh
  make dep
  ```

  Observed behavior from Makefiles:
  - Go: `go mod download`, `go mod tidy`, `go mod vendor`.
  - Ruby (in `test/`): `bundler check || bundler install` with `bundler config set path vendor/bundle`.

### Build

- Build the release binary in repo root:

  ```sh
  make build
  ```

  Produces `./migrieren`.

- Build a test binary with build tag `features`:

  ```sh
  make build-test
  ```

### Run

- The CLI entrypoint is `main.go` and registers a `server` command (`internal/cmd/server.go`).

Example run using the repo’s dev/test config:

```sh
./migrieren server -i file:test/.config/server.yml
```

There is also a dev helper:

```sh
make dev
```

(Observed: `dev` uses `air` and runs `make dep build` before launching `./migrieren server -i file:test/.config/server.yml`.)

### Test

- Go tests (via `gotestsum`, outputs to `test/reports/`):

  ```sh
  make specs
  ```

- Ruby feature tests (Cucumber):

  ```sh
  make features
  # runs: make -C test features
  ```

- Benchmarks (Ruby harness):

  ```sh
  make benchmarks
  ```

- Coverage report generation:

  ```sh
  make coverage
  ```

Notes (observed in Makefiles / `main_test.go`):
- Many Go test invocations use `-race -vet=off -mod vendor`.
- `main_test.go` is guarded by `//go:build features`.

### Lint / format / security

- Lint everything:

  ```sh
  make lint
  ```

  Observed lint sources:
  - Go: `bin/build/go/lint run` (golangci-lint wrapper) + `bin/build/go/fa`.
  - Ruby: `bundler exec rubocop`.
  - Protos: `buf lint`.

- Auto-fix where possible:

  ```sh
  make fix-lint
  ```

- Format everything:

  ```sh
  make format
  ```

- Security scan:

  ```sh
  make sec
  ```

  Observed: `govulncheck -show verbose -test ./...`.

- Repo scan (Trivy):

  ```sh
  make trivy-repo
  ```

### Protobuf / API (buf)

The protobuf contract is in `api/migrieren/v1/service.proto`.

From repo root:

```sh
make proto-generate
make proto-breaking
make proto-lint
make proto-format
```

Directly in `api/`:

```sh
make -C api generate
make -C api breaking
make -C api lint
make -C api format
```

Observed generation behavior (`api/buf.gen.yaml`):
- Go protobuf + Go gRPC stubs output into `api/` (source-relative).
- Ruby protobuf + Ruby gRPC stubs output into `test/lib/`.

## 3) CI / workflow notes (observed)

CircleCI (`.circleci/config.yml`) does roughly:
- init submodules
- `make dep`
- `make lint`
- `make proto-breaking`
- `make sec`
- `make trivy-repo`
- `make features`
- `make benchmarks`
- `make analyse`
- `make coverage` then `make codecov-upload`

CI uses containers including:
- Postgres on `localhost:5432`
- Grafana Mimir on `localhost:9009`

## 4) Repository layout (where to look)

### Entrypoints

- `main.go`: wires CLI application.
- `internal/cmd/server.go`: registers the `server` command.
- `internal/cmd/module.go`: composes DI modules.

### Configuration

- `internal/config/config.go`: service config struct, embeds `*go-service/v2/config.Config` and includes:
  - `health` (`internal/health/config.go`)
  - `migrate` (`internal/migrate/config.go`)

A sample config used for dev/tests exists at `test/.config/server.yml`.

### Migration core

- `internal/migrate/migrate.go`: core migration/ping logic built on `github.com/golang-migrate/migrate/v4`.
- `internal/migrate/config.go`: configured databases list + helpers.
- `internal/migrate/source/source.go`: imports migrate source drivers (`file`, `github`).
- `internal/migrate/database/database.go`: database driver wiring (observed support: `pgx5://...`).
- `internal/migrate/telemetry/logger/logger.go`: in-memory migration logger.

### API / transport

- `api/migrieren/v1/service.proto`: gRPC contract.
- `internal/api/migrate/`: transport-facing migrator (reads source/URL bytes using `go-service/v2/os.FS`).
- `internal/api/v1/transport/grpc/`: gRPC server + handler implementation.
- `internal/api/v1/transport/http/`: HTTP routing via `go-service/v2/net/http/rpc`.

### Health

- `internal/health/health.go`: registers health checks (includes per-database checks) via `go-health/v2/server`.
- `internal/health/checker/checker.go`: checker that pings migrator.

### Ruby test harness

- `test/features/**`: Cucumber features + step definitions.
- `test/lib/migrieren.rb`: shared client helpers; constructs:
  - HTTP client to `http://localhost:11000`
  - gRPC client to `localhost:12000`

## 5) Conventions and patterns (observed)

### Dependency injection

Uses `github.com/alexfalkowski/go-service/v2/di` with Fx-style modules:
- Each subsystem exports `var Module = di.Module(...)` (e.g. `internal/api/v1/module.go`, `internal/migrate/module.go`).
- Constructors registered via `di.Constructor(...)`.
- Side-effect registrations (routes, gRPC server registration) via `di.Register(...)`.

### Error handling

- Domain errors are exported vars in packages and returned upward (e.g. migrate errors in `internal/migrate/migrate.go`).
- gRPC transport maps a “not found” error to `codes.NotFound` and everything else to `codes.Internal` (`internal/api/v1/transport/grpc/grpc.go`).
- Migration errors are attached to request metadata via `meta.WithAttribute(...)`.

### Formatting / lint

- `.editorconfig`:
  - Go files: tabs
  - Makefiles: tabs
  - default: 2-space indentation
- `.golangci.yml`:
  - `default: all` with a curated disable list
  - excludes generated `.pb*` files
  - enables formatters including `gofmt`, `gofumpt`, `goimports`, `gci`

## 6) Gotchas (observed)

- **Submodule required**: Make targets call scripts under `bin/`; init/update submodule before running most automation.
- **Vendoring is relied on**: multiple Go targets run with `-mod vendor`; run `make dep` after changing Go deps.
- **Config-driven source/URL**: migration `source` and DB `url` are loaded via `go-service/v2/os.FS` from config values like `file:secrets/pg` (see `test/.config/server.yml`).
- **DB URL scheme**: database driver expects `pgx5://...` and rewrites to `postgres://...` internally (`internal/migrate/database/database.go`).

## 7) Tooling used by Make targets (non-exhaustive, observed)

Some targets assume these tools exist on PATH (or are provided by the `bin/` submodule wrappers):
- Go: `gotestsum`, `govulncheck`
- Proto: `buf`
- Ruby: `bundler`, `rubocop`, `cucumber`
- Dev: `air`
- Security/CI: `codecovcli`, `trivy` (invoked via `bin/`)
- Misc targets reference: `mkcert`, `dot` (Graphviz), `goda`, `gsa`, `scc` (only if you run those specific targets)
