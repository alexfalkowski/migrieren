# Accepted design for the runtime

Part of the [accepted-design index](../accepted-design.md).

## Required runtime configuration

Top-level config validation requires the `Health` and `Migrate` pointers and
the embedded `*go-service/v3/config.Config`. Missing sections fail validation.
Keep `healthConfig`, `migrateConfig`, and `decorateConfig` as direct projections.

[Runtime configuration](../../../internal/config/config.go) ·
[Configuration scenarios](../../../test/features/config/config.feature)

## Database telemetry setup errors

The migration path returns errors from `telemetry.Open` and
`telemetry.RegisterDBStatsMetrics`. The ping path returns `telemetry.Open`
errors and skips DB stats registration. Consider changing these paths to
`runtime.Must` only if a concrete service startup invariant requires fail-fast
behavior.

[Database driver and ping](../../../internal/migrate/database/database.go)

## Health probes

`noop` and `online` are reserved probe names by convention. Database-name
validation leaves enforcement to configuration authors. Revisit enforcement
if there is concrete evidence of configuration misuse or the owner asks to
enforce the reservation.

HTTP `healthz` observes `online` and the per-database checks. gRPC health for
`migrieren.v1.Service` observes `noop`. The default feature configuration
therefore expects failing HTTP `healthz` from intentionally invalid database
entries and healthy gRPC health.

[Probe registration](../../../internal/health/health.go) ·
[Feature configuration](../../../test/.config/server.yaml) ·
[Operator guidance](../../../README.md)
