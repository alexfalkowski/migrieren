# Log migration fixture

This directory is a feature-test fixture for migration log scenarios in the
HTTP and gRPC API feature suites. The apply-all scenarios expect this fixture
to converge at migration version `40`.

The `logs` database target in `test/.config/server.yml` points at this source
through `test/secrets/log_source`. It intentionally contains 40 tiny migrations
so a real `golang-migrate` run emits more than 100 log lines. That lets the
feature test verify that returned migration logs are capped and marked with
`migration logs truncated`.

These files are not product migrations.
