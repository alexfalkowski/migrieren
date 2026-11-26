# AGENTS.md

This repository is a Go service (plus a Ruby feature-test harness) that runs database migrations via a gRPC API and an HTTP RPC façade.

## Shared guidance

Use `bin/AGENTS.md` for shared skills and cross-repository defaults.

## Accepted design

Read the [accepted-design index](docs/agents/accepted-design.md) and the relevant
area files before proposing a change or recording a finding or audit entry.
Apply the recorded scope of each decision and add owner-accepted decisions to
the relevant area file as they are agreed.

## Recent session notes (keep for future sessions)

### Vendoring can break Go tooling

If you see “inconsistent vendoring” errors, run:

```sh
make dep
```

This refreshes Go modules and updates `vendor/` (and also ensures Ruby gems for `test/` are installed).

### Documentation conventions used in this repo

- **Go**
  - Package-level docs must live in a `doc.go` file in that package.
  - Prefer documenting exported constructors/types/methods with concrete:
    - inputs/outputs,
    - error behavior (sentinel errors vs underlying errors),
    - operational details (resource management, logging/telemetry),
    - examples where it clarifies usage.
  - Do not edit generated `.pb.go` files; update `.proto` comments instead, then re-generate.

- **Ruby (feature-test harness)**
  - The public “API” is primarily `test/lib/**` (helpers used by Cucumber steps).
  - Add RDoc on modules/classes and public methods describing:
    - intent (why it exists in the harness),
    - parameters and return types,
    - endpoint/route mapping for HTTP helpers,
    - any environment assumptions (localhost ports, DB URI).
  - Do not edit generated protobuf/grpc Ruby stubs; update `.proto` comments instead, then re-generate.

### Linting Ruby code

Feature-test Ruby linting is typically run via:

```sh
make -C test lint
```

Use the Makefile wiring instead of direct Ruby tool commands.

## 0) First check

If `bin/` is missing, most `make` targets will fail. Use `make submodule` once
the shared checkout is present; see `bin/AGENTS.md` for fresh-clone bootstrap
details.

## 1) Project type

- **Primary**: Go service; use `go.mod` for module path and toolchain details.
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
  - Go dependencies and vendor state are refreshed.
  - Ruby harness dependencies in `test/` are installed through the shared Ruby
    Make target.

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

Use the dev helper for local runtime work:

```sh
make dev
```

(Observed: `dev` uses the shared live-reload path and the repo dev/test config.)

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
  - Go linting through shared `bin` helpers.
  - Ruby linting through the test harness Makefile.
  - Protobuf linting through the API Makefile.

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
make proto-stale
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
- initialize submodules
- `make dep`
- `make lint`
- `make proto-breaking`
- `make proto-stale`
- `make sec`
- `make trivy-repo`
- `make features`
- `make benchmarks`
- `make analyse`
- `make coverage` then `make codecov-upload`

CI uses containers including:
- Postgres on `localhost:5432`
- Grafana Mimir on `localhost:9009`

Feature tests themselves talk to Postgres through the nonnative proxy on `localhost:5433`, which forwards to the backing Postgres container on `localhost:5432` (see `test/nonnative.yml` and `test/secrets/pg`).

## 4) Repository layout (where to look)

### Entrypoints

- `main.go`: wires CLI application.
- `internal/cmd/server.go`: registers the `server` command.
- `internal/cmd/module.go`: composes DI modules.

### Configuration

- `internal/config/config.go`: service config struct, embeds `*go-service/v3/config.Config` and includes:
  - `health` (`internal/health/config.go`)
  - `migrate` (`internal/migrate/config.go`)

A sample config used for dev/tests exists at `test/.config/server.yaml`.

### Migration core

- `internal/migrate/migrate.go`: core migration/ping logic built on `github.com/golang-migrate/migrate/v4`.
- `internal/migrate/config.go`: configured databases list + helpers.
- `internal/migrate/source/source.go`: imports migrate source drivers (`file`, `github`).
- `internal/migrate/database/database.go`: database driver wiring (observed support: `pgx5://...`).
- `internal/migrate/telemetry/logger/logger.go`: in-memory migration logger.

### API / transport

- `api/migrieren/v1/service.proto`: gRPC contract.
- `internal/api/migrate/`: API-facing migrator that resolves database names through config.
- `internal/api/v1/migrate/`: versioned API contract shared by HTTP and gRPC transports.
- `internal/api/v1/transport/grpc/`: gRPC server + handler implementation.
- `internal/api/v1/transport/http/`: HTTP routing via `go-service/v3/net/http/rpc`.

### Health

- `internal/health/health.go`: registers health checks (includes per-database checks) via `go-health/v2/server`.
- `internal/health/checker/checker.go`: checker that pings migrator.

### Ruby test harness

- `test/features/**`: Cucumber features + step definitions.
- `test/lib/migrieren.rb`: shared client helpers; constructs:
  - HTTP client to `http://localhost:11000`
  - gRPC client to `localhost:12000`
- `test/nonnative.yml`: local process/service orchestration for features; starts the service binary and proxies Postgres from `localhost:5433` to `localhost:5432`

## 5) Conventions and patterns (observed)

### Dependency injection

Uses `github.com/alexfalkowski/go-signal/v2` for dependency injection:
- Each subsystem exports `var Module = signal.Module(...)` (e.g. `internal/api/v1/module.go`, `internal/migrate/module.go`).
- Constructors registered via `signal.Constructor(...)`.
- Side-effect registrations (routes, gRPC server registration) via `signal.Register(...)`.

### Error handling

- Domain errors are exported vars in packages and returned upward (e.g. migrate errors in `internal/migrate/migrate.go`).
- gRPC transport maps “not found” to `codes.NotFound`, cancellation to `codes.Canceled`, deadlines to `codes.DeadlineExceeded`, and other migration failures to `codes.Internal` (`internal/api/v1/transport/grpc/grpc.go`).
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
- **Config-driven source/URL**: migration `source` and DB `url` are loaded via `go-service/v3/os.FS` from config values like `file:secrets/pg` (see `test/.config/server.yaml`).
- **DB URL scheme**: database driver expects `pgx5://...` and rewrites to `postgres://...` internally (`internal/migrate/database/database.go`).
- **Checked-in test config is intentionally mixed**: `test/.config/server.yaml` contains both valid and invalid database entries to exercise failure paths in feature tests.

## 7) Tooling used by Make targets

Some targets assume external tools exist on `PATH` or are provided by shared
`bin` wrappers. Prefer the Make targets as the documented command surface, and
only run direct tools for narrow diagnosis when a target reports a missing
dependency.
