![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/migrieren.svg?style=shield)](https://circleci.com/gh/alexfalkowski/migrieren)
[![codecov](https://codecov.io/gh/alexfalkowski/migrieren/graph/badge.svg?token=R2OD8WIKD0)](https://codecov.io/gh/alexfalkowski/migrieren)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/migrieren)](https://goreportcard.com/report/github.com/alexfalkowski/migrieren)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/migrieren.svg)](https://pkg.go.dev/github.com/alexfalkowski/migrieren)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# Migrieren

Migrieren provides a way to migrate your databases.

## Background

Migrating databases is an interesting topic with many caveats. Basically we have 2 categories:
- Migrate the database before the application is deployed.
- Migrate the database at the most convenient time.

We don't have a preferred method. We just want to provide you with the best option.

### Why a service?

Well every language or framework provides a way to migrate. Though a lot of them are tied to [ORMs](https://en.wikipedia.org/wiki/Object%E2%80%93relational_mapping), which in our experience are not the best tool. We find writing migrations in the native language of the database to be far superior.

Since you are more than likely going to use microservices we don't need to reinvent the wheel for every framework. Just use the service!

### Migrations

There are some best practices regarding how to write effective schema migration scripts. While this service does not enforce it, you should be aware of it.

Some great information can be found in [Update your Database Schema Without Downtime](https://thorben-janssen.com/update-database-schema-without-downtime/).

## Server

The server is defined by the following [proto contract](api/migrieren/v1/service.proto). So each version of the service will have a new contract.

### Databases

This system allows you to configure many databases.

To configure we just need the have the following configuration:

```yaml
migrate:
  databases:
    -
      name: db1
      source: file://migrations
      url: path to url
    -
      name: db2
      source: file:///migrations
      url: path to url
    -
      name: db3
      source: file://migrations
      url: path to url
```

Each database has the following properties:
- A distinct name.
- The source of the migrations (file, GitHub, etc).
- The database URL (MySQL, PostgreSQL, etc).

## Health

The system defines a way to monitor all of it's dependencies.

To configure we just need the have the following configuration:

```yaml
health:
  duration: 1s (how often to check)
  timeout: 1s (when we should timeout the check)
```

## Deployment

Since we are advocating building microservices, you would normally use a [container orchestration system](https://newrelic.com/blog/best-practices/container-orchestration-explained). Here is what we recommend when using this system:
- You could have a global migration service or shard these services per [bounded context](https://martinfowler.com/bliki/BoundedContext.html).
- The client should be used as an [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/).

## Design

The service is based around the awesome work of [migrate](https://github.com/golang-migrate/migrate). So please check that out to see how to best use it. It can support many configurations that can be easily added.

## Development

If you would like to contribute, here is how you can get started.

### Structure

The project follows the structure in [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

### Dependencies

Please make sure that you have the following installed:
- [Ruby](.ruby-version)
- Golang

### Style

This project favours the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Setup

The get yourself setup, please run the following:

```sh
make setup
```

### Binaries

To make sure everything compiles for the app, please run the following:

```sh
make build-test
```

### Features

To run all the features, please run the following:

```sh
make features
```

### Changes

To see what has changed, please have a look at `CHANGELOG.md`
