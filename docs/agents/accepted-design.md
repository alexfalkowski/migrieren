# Accepted design

This index records Migrieren's owner-accepted design choices and tradeoffs.
The area files carry the same weight as [AGENTS.md](../../AGENTS.md).

Read the relevant files before proposing a change or recording a finding or
audit entry. Revisit a decision when the owner asks or new evidence changes
its documented contract or tradeoff.

- [Migration service](accepted-design/migration-service.md) covers execution
  modes, plan output, dirty-state recovery, and the CLI boundary.
- [Runtime](accepted-design/runtime.md) covers configuration, database
  telemetry setup errors, and health probes.
- [Project workflow](accepted-design/project-workflow.md) covers feature
  coverage, CI services, Docker releases, and dependency changes.
