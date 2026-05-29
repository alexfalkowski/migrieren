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
