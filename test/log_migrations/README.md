# Log migration fixture

This directory is a feature-test fixture for migration log scenarios in the
HTTP and gRPC API feature suites. The apply-all scenarios expect this fixture
to converge at migration version `40`.

The `logs` database target in `test/.config/server.yml` points at this source
through `test/secrets/log_source`. It intentionally contains 40 tiny migrations
so a real `golang-migrate` run emits far more log lines than the configured
maximum (`migrate.logs.max`, set to `20` for these suites). That lets the
feature tests verify that returned migration logs are capped at the configured
maximum and start with a `migration logs truncated (showing last N of M)` marker.

These files are not product migrations.
