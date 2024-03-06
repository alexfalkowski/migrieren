FROM golang:1.22.1-bullseye AS build

ARG version=latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=${version}'" -a -o migrieren main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /app/migrieren /migrieren
ENTRYPOINT ["/migrieren"]
