# Accepted design for the migration service

Part of the [accepted-design index](../accepted-design.md).

## Migration execution

The execution API provides `Migrate` for an explicit positive target version
and `ApplyMigrations` for converging to the latest available up migration.
Operators rolling back choose an explicit target version and call `Migrate`.
Consider step-based up or down execution, rollback by step, all-down execution,
or version-zero targets only if the owner reopens the execution scope.

[v1 service contract](../../../api/migrieren/v1/service.proto)

## Pending migrations in plans

`PlanMigrations` lists pending source version numbers. Operators can look up
the corresponding names in the migration source. Including names would require
opening each pending migration with the source driver's `ReadUp` to get its
identifier. Revisit the tradeoff if the owner asks for names or operators show
that version-number lookup is insufficient.

[Plan response contract](../../../api/migrieren/v1/service.proto)

## Dirty-state recovery

`Status` reports dirty migration metadata. Repair stays with the operator
because the service can verify the migration metadata but cannot verify that
the underlying schema or data was correctly repaired before clearing the dirty
flag. Consider a public force or repair RPC only at the owner's request for a
dedicated design covering safeguards, auditability, documentation, and tests.

[Status contract](../../../api/migrieren/v1/service.proto) ·
[Operator guidance](../../../README.md)

## Service and CLI boundary

`server` is the only command Migrieren registers. The framework provides help
and version support. Remote execution and inspection use the gRPC API or HTTP
RPC façade. Deployments can enable token verification and Casbin authorization
on these transports. The upstream `golang-migrate` CLI covers local one-shot
execution, including deploy jobs. Revisit a Migrieren batch or client
subcommand only if the owner asks for a second entry point.

[CLI wiring](../../../main.go) · [Service usage](../../../README.md)
