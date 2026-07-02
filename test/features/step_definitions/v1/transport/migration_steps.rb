# frozen_string_literal: true

def expect_postgres_migration(version)
  expect_postgres_migration_table(version)
  expect_postgres_accounts(version)
end

def expect_log_migration(version)
  expect(Migrieren.pg.log_migration_version).to eq(version)
end

def expect_postgres_migration_table(version)
  expect(Migrieren.pg.table?('schema_migrations')).to be(false)
  expect(Migrieren.pg.table?('migrieren_schema_migrations')).to be(true)
  expect(Migrieren.pg.migration_version).to eq(version)
end

def expect_postgres_accounts(version)
  expect(Migrieren.pg.table?('accounts')).to be(true)
  expect(Migrieren.pg.account_count).to be > 0
  expect(Migrieren.pg.column?('accounts', 'update_at')).to eq(version == 2)
end

MIGRATION_STATES = {
  1 => 'unapplied',
  2 => 'clean',
  3 => 'dirty',
  'MIGRATION_STATE_UNAPPLIED' => 'unapplied',
  'MIGRATION_STATE_CLEAN' => 'clean',
  'MIGRATION_STATE_DIRTY' => 'dirty'
}.freeze

def migration_state(state)
  MIGRATION_STATES.fetch(state) { MIGRATION_STATES.fetch(state.to_s) }
end

MIGRATION_DIRECTIONS = {
  1 => 'none',
  2 => 'up',
  'MIGRATION_DIRECTION_NONE' => 'none',
  'MIGRATION_DIRECTION_UP' => 'up'
}.freeze

def migration_direction(direction)
  MIGRATION_DIRECTIONS.fetch(direction) { MIGRATION_DIRECTIONS.fetch(direction.to_s) }
end

def migration_versions(versions)
  return Range.new(*versions.split('..').map(&:to_i)).to_a if versions.include?('..')

  versions.split(',').reject(&:empty?).map(&:to_i)
end

def empty_to_nil(value)
  value.empty? ? nil : value
end

Then('I should not see a completed timeout migration') do
  expect(Migrieren.pg.timeout_migration_version).to be_nil
end
