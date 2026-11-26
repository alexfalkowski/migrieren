# Accepted design for the project workflow

Part of the [accepted-design index](../accepted-design.md).

## Feature testing and the GitHub fixture

Ruby/Cucumber is the primary automated coverage for service behavior. Assess
coverage through those scenarios when deciding whether package-local Go tests
are needed.

The regular HTTP and gRPC feature scenarios exercise the `github://` source
driver against the shared remote fixture. Keep that coverage in the default
`make features` gate. Raise concrete fixture failures, missing credentials, or
upstream behavior changes when there is current evidence.

[Feature suite](../../../test/features) ·
[Feature configuration](../../../test/.config/server.yaml) ·
[GitHub source setup](../../../scripts/github-source)

## CI auxiliary services

CircleCI owns the auxiliary Mimir service used by `build-service`. Revisit
that choice when there is an observed failure or incompatibility.

[CircleCI configuration](../../../.circleci/config.yml)

## Docker publication

The architecture-specific jobs build and push images before the manifest job
publishes the multi-architecture manifest. The workflow allows automatic
reruns. Raise a reliability finding when there is evidence of a current
failure mode such as tag drift, an unscanned published digest, harmful rerun
side effects, or a rollback or reproducibility incident.

[CircleCI release workflow](../../../.circleci/config.yml)

## Dependency changes

Use the repository's Make targets to update dependencies and the shared
`review-pr` workflow to prepare review PRs. That workflow stages all local
changes before committing, including dependency files produced by the Make
targets. Raise a reliability finding when there is evidence of dependency
drift that this workflow failed to capture.

[Repository Makefile](../../../Makefile) ·
[Shared review PR workflow](../../../bin/skills/review-pr/SKILL.md)
