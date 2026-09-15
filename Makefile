include bin/build/make/help.mak
include bin/build/make/grpc.mak
include bin/build/make/git.mak
include bin/build/make/claude.mak
include bin/build/make/codex.mak

# Generate the local GitHub migration source from MIGRIEREN_GITHUB_TOKEN.
github-source:
	@./scripts/github-source

# Run HTTP-tagged Cucumber features.
features-http:
	@$(MAKE) features tags="@http"

# Run gRPC- and configuration-tagged Cucumber features.
features-grpc:
	@$(MAKE) features tags="@grpc or @config"

# Stage coverage profiles in an isolated CircleCI workspace partition (partition=http|grpc|benchmarks).
coverage-workspace-stage:
	@./scripts/coverage-workspace stage "$(partition)"

# Restore staged CircleCI coverage profiles to test/reports for merging.
coverage-workspace-restore:
	@./scripts/coverage-workspace restore
