# frozen_string_literal: true

def expect_postgres_migration(version)
  expect_postgres_migration_table(version)
  expect_postgres_accounts(version)
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

Then('I should not see a completed timeout migration') do
  expect(Migrieren.pg.timeout_migration_version).to be_nil
end
