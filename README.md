[![CircleCI](https://circleci.com/gh/alexfalkowski/migrieren.svg?style=svg)](https://circleci.com/gh/alexfalkowski/migrieren)
[![Coverage Status](https://coveralls.io/repos/github/alexfalkowski/migrieren/badge.svg?branch=master)](https://coveralls.io/github/alexfalkowski/migrieren?branch=master)

# Migrieren

Migrieren provides a way to migrate your databases.

## Rational

Migrating databases is an interesting topic with many caveats. Basically we have 2 categories:
- Migrate the database before the application is deployed.
- Migrate the database at the most convenient time.

We don't have a preferred method. We just want to provide you with the best option.

### Why a service?

Well every language or framework provides a way to migrate. Though a lot of them are tied to [ORMs](https://en.wikipedia.org/wiki/Object%E2%80%93relational_mapping), which in our experience are not the best tool. We find writing migrations in the native language of the database to be far superior.

Since you are more than likely going to use microservices we don't need to reinvent the wheel for every framework. Just use the service!

## Design

The service is based around the awesome work [migrate](https://github.com/golang-migrate/migrate). So please check that out to see how to best use it. It can support many configurations that can be easily added.

For now we support the following sources:
- File
- GitHub

For now we support the following databases:
- MySQL
- PostgreSQL

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
