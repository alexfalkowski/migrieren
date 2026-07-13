include bin/build/make/help.mak
include bin/build/make/grpc.mak
include bin/build/make/git.mak
include bin/build/make/claude.mak
include bin/build/make/codex.mak

# Generate the local GitHub migration source from MIGRIEREN_GITHUB_TOKEN.
github-source:
	@./scripts/github-source
