# frozen_string_literal: true

module Migrieren
  ##
  # Helper for interacting with the local Postgres instance used by feature tests.
  #
  # This class is part of the Ruby feature-test harness under `test/` and is used
  # by step definitions to perform direct database setup/teardown tasks outside
  # of the Migrieren API.
  #
  # The helper connects directly to the backing local Postgres server using the
  # `pg` gem with a fixed URI intended for the test environment. This direct
  # connection uses `localhost:5432`; the service under test reaches the same
  # database through the nonnative proxy configured at `localhost:5433`.
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
    # - `migrieren_schema_migrations`
    # - `migrieren_timeout_schema_migrations`
    # - `migrieren_log_schema_migrations`
    #
    # @return [void]
    def destroy
      @conn.exec('DROP TABLE IF EXISTS accounts')
      @conn.exec('DROP TABLE IF EXISTS schema_migrations')
      @conn.exec('DROP TABLE IF EXISTS migrieren_schema_migrations')
      @conn.exec('DROP TABLE IF EXISTS migrieren_timeout_schema_migrations')
      @conn.exec('DROP TABLE IF EXISTS migrieren_log_schema_migrations')
    end

    ##
    # Seeds the managed tables, then runs {#destroy}.
    #
    # This exercises the cleanup path from the feature harness without modeling
    # it as an application feature. Callers must check {#destroyed?}, or use the
    # support hook, to assert the cleanup postcondition.
    #
    # @return [void]
    def verify_destroy
      @conn.exec('CREATE TABLE IF NOT EXISTS accounts (user_id serial PRIMARY KEY)')
      @conn.exec('CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL)')
      @conn.exec('CREATE TABLE IF NOT EXISTS migrieren_schema_migrations (version bigint NOT NULL)')
      @conn.exec('CREATE TABLE IF NOT EXISTS migrieren_timeout_schema_migrations (version bigint NOT NULL)')
      @conn.exec('CREATE TABLE IF NOT EXISTS migrieren_log_schema_migrations (version bigint NOT NULL)')

      destroy
    end

    ##
    # Checks whether all tables managed by {#destroy} are absent.
    #
    # @return [Boolean] true when cleanup left no managed tables behind
    def destroyed?
      !table?('accounts') && !table?('schema_migrations') &&
        !table?('migrieren_schema_migrations') && !table?('migrieren_timeout_schema_migrations') &&
        !table?('migrieren_log_schema_migrations')
    end

    ##
    # Checks whether a table exists in the public schema.
    #
    # @param name [String] table name
    # @return [Boolean] true when the table exists
    def table?(name)
      query = <<~SQL
        SELECT EXISTS (
          SELECT 1
          FROM information_schema.tables
          WHERE table_schema = 'public' AND table_name = $1
        )
      SQL

      @conn.exec_params(query, [name]).getvalue(0, 0) == 't'
    end

    ##
    # Checks whether a column exists on a table in the public schema.
    #
    # @param table [String] table name
    # @param column [String] column name
    # @return [Boolean] true when the column exists
    def column?(table, column)
      query = <<~SQL
        SELECT EXISTS (
          SELECT 1
          FROM information_schema.columns
          WHERE table_schema = 'public' AND table_name = $1 AND column_name = $2
        )
      SQL

      @conn.exec_params(query, [table, column]).getvalue(0, 0) == 't'
    end

    ##
    # Returns the number of rows in the accounts fixture table.
    #
    # This is intended for post-migration assertions. If `accounts` has not been
    # created yet, the underlying `pg` exception is raised.
    #
    # @return [Integer] row count
    def account_count
      @conn.exec('SELECT COUNT(*) FROM accounts').getvalue(0, 0).to_i
    end

    ##
    # Returns the current clean migration version from the configured migration table.
    #
    # This is intended for post-migration assertions. If
    # `migrieren_schema_migrations` has not been created yet, the underlying `pg`
    # exception is raised.
    #
    # @return [Integer, nil] the version, or nil when the table has no row
    def migration_version
      result = @conn.exec('SELECT version FROM migrieren_schema_migrations WHERE dirty = false LIMIT 1')

      result.ntuples.zero? ? nil : result.getvalue(0, 0).to_i
    end

    ##
    # Returns the clean migration version from the timeout migration table.
    #
    # @return [Integer, nil] the version, or nil when no clean timeout migration exists
    def timeout_migration_version
      return nil unless table?('migrieren_timeout_schema_migrations')

      result = @conn.exec('SELECT version FROM migrieren_timeout_schema_migrations WHERE dirty = false LIMIT 1')

      result.ntuples.zero? ? nil : result.getvalue(0, 0).to_i
    end

    ##
    # Returns the clean migration version from the log migration table.
    #
    # @return [Integer, nil] the version, or nil when no clean log migration exists
    def log_migration_version
      return nil unless table?('migrieren_log_schema_migrations')

      result = @conn.exec('SELECT version FROM migrieren_log_schema_migrations WHERE dirty = false LIMIT 1')

      result.ntuples.zero? ? nil : result.getvalue(0, 0).to_i
    end
  end
end
