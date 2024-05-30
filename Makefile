include bin/build/make/grpc.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: client-diagram server-diagram

client-diagram:
	$(MAKE) package=client create-diagram

server-diagram:
	$(MAKE) package=server create-diagram
