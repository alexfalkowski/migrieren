.PHONY: vendor

include bin/build/make/service.mak

# Build release binary.
build:
	go build -race -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=latest'" -mod vendor -o migrieren main.go

# Build test binary.
build-test:
	go test -race -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=latest'" -mod vendor -c -tags features -covermode=atomic -o migrieren -coverpkg=./... github.com/alexfalkowski/migrieren

# Release to docker hub.
docker:
	bin/build/docker/push migrieren
