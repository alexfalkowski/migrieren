FROM golang:1.20.5-bullseye AS build

ARG version=latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags="-X 'github.com/alexfalkowski/migrieren/cmd.Version=${version}'" -a -o migrieren main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/migrieren /migrieren
ENTRYPOINT ["/migrieren"]
