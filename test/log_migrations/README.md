# Log migration fixture

This directory is a feature-test fixture for the bounded migration log scenario
in `test/features/v1/transport/grpc/api.feature`.

The `logs` database target in `test/.config/server.yml` points at this source
through `test/secrets/log_source`. It intentionally contains 40 tiny migrations
so a real `golang-migrate` run emits more than 100 log lines. That lets the
feature test verify that returned migration logs are capped and marked with
`migration logs truncated`.

These files are not product migrations.
