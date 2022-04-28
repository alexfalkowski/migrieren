FROM golang:1.18.1-bullseye AS build

ARG version=latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=${version}'" -a -o migrieren main.go

FROM debian:bullseye-slim

WORKDIR /

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get -y upgrade && \
    apt-get install -y --no-install-recommends \
    ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /app/migrieren /migrieren
ENTRYPOINT ["/migrieren"]
