# frozen_string_literal: true

module Migrieren
  ##
  # Helper for interacting with the local Postgres instance used by feature tests.
  #
  # This class is part of the Ruby feature-test harness under `test/` and is used
  # by step definitions to perform direct database setup/teardown tasks outside
  # of the Migrieren API.
  #
  # The helper connects to a local Postgres server using the `pg` gem with a
  # fixed URI intended for the test environment.
  #
  # @example Dropping test tables between scenarios
  #   Migrieren.pg.destroy
  #
  class PG
    ##
    # Creates a new Postgres helper connected to the local test database.
    #
    # Connection details are derived from a fixed URI:
    # `postgres://test:test@localhost:5432/test?sslmode=disable`.
    #
    # Notices are suppressed to keep test output clean.
    #
    # @return [Migrieren::PG]
    def initialize
      uri = URI.parse('postgres://test:test@localhost:5432/test?sslmode=disable')
      @conn = ::PG.connect(uri.hostname, uri.port, nil, nil, uri.path[1..], uri.user, uri.password)

      @conn.set_notice_processor { |_| }
    end

    ##
    # Drops tables used by the feature-test suite if they exist.
    #
    # This is intended to reset the database to a known-clean state.
    #
    # Currently, it drops:
    # - `accounts`
    # - `schema_migrations`
    #
    # @return [void]
    def destroy
      @conn.exec('DROP TABLE IF EXISTS accounts')
      @conn.exec('DROP TABLE IF EXISTS schema_migrations')
    end
  end
end
