.PHONY: vendor

include bin/build/make/service.mak

# Build release binary.
build:
	go build -race -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=latest'" -mod vendor -o migrieren main.go

# Build test binary.
build-test:
	go test -race -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=latest'" -mod vendor -c -tags features -covermode=atomic -o migrieren -coverpkg=./... github.com/alexfalkowski/migrieren

sanitize-coverage:
	bin/quality/go/cov

# Get the HTML coverage for go.
html-coverage: sanitize-coverage
	go tool cover -html test/reports/final.cov

# Get the func coverage for go.
func-coverage: sanitize-coverage
	go tool cover -func test/reports/final.cov

# Send coveralls data.
goveralls: sanitize-coverage
	goveralls -coverprofile=test/reports/final.cov -service=circle-ci -repotoken=zl0TVeSjn3TgnoATsUhpQycpFScnOoyji

# Run go security checks.
go-sec:
	gosec -quiet -exclude-dir=test -exclude=G104 ./...

# Run security checks.
sec: go-sec

# Release to docker hub.
docker:
	bin/build/docker/push migrieren

# Start the environment.
start:
	bin/build/docker/env start

# Stop the environment.
stop:
	bin/build/docker/env stop
