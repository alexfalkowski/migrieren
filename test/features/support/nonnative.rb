# frozen_string_literal: true

Nonnative.configure do |config|
  config.load_file('nonnative.yml')
end

BeforeAll do
  Migrieren.pg.verify_destroy
  expect_destroyed_database
end

Before('@clean') do
  Migrieren.pg.destroy
  expect_destroyed_database
end

After('@clean') do
  Migrieren.pg.destroy
  expect_destroyed_database
end

def expect_destroyed_database
  return if Migrieren.pg.destroyed?

  raise 'expected Migrieren.pg.destroy to remove managed tables'
end
